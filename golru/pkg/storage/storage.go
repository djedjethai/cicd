package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
)

var ErrorNoSuchKey = errors.New("no such key")

type StorageRepo interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Delete(context.Context, string) error
}

type storage struct {
	sync.RWMutex
	store map[string]*node
	dll   dll
}

func NewStorage(maxLgt int) StorageRepo {
	store := make(map[string]*node)
	return &storage{
		store: store,
		dll:   NewDll(maxLgt),
	}
}

func (s *storage) Set(ctx context.Context, key string, value string) error {
	tr := otel.GetTracerProvider().Tracer("try")
	_, sp := tr.Start(ctx, "SetStorage")
	defer sp.End()

	fmt.Println("is all uuiiuuiiio")

	// create node in dll
	newN, outN := s.dll.unshift(key, value)
	if outN != nil {
		// in case dll poped out the last item
		s.Lock()
		delete(s.store, outN.key)
		s.Unlock()
	}

	s.Lock()
	// add node to map
	s.store[key] = newN
	s.Unlock()

	return nil
}

func (s *storage) Get(ctx context.Context, key string) (string, error) {
	tr := otel.GetTracerProvider().Tracer("try")
	_, sp := tr.Start(ctx, "GetStorage")
	defer sp.End()

	s.RLock()
	nd, ok := s.store[key]
	s.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}

	// move the Get node(so nd) to head of dll
	ndExist := s.dll.removeNode(nd)
	if ndExist != nil {
		// re-unshift ndExist or nd(what ever they point to the same location)
		s.dll.unshiftNode(ndExist)
	}

	return nd.val, nil
}

func (s *storage) Delete(ctx context.Context, key string) error {
	tr := otel.GetTracerProvider().Tracer("try")
	_, sp := tr.Start(ctx, "DeleteStorage")
	defer sp.End()

	s.RLock()
	nd, ok := s.store[key]
	s.RUnlock()

	if ok {
		_ = s.dll.removeNode(nd)

		s.Lock()
		delete(s.store, key)
		s.Unlock()
	}

	return nil
}
