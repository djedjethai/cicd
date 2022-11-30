package getter

import (
	"context"
	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/storage"
)

type Getter interface {
	Get(string) (string, error)
}

type getter struct {
	st storage.StorageRepo
	tr config.Tracer
}

func NewGetter(s storage.StorageRepo, t config.Tracer) Getter {
	return &getter{
		st: s,
		tr: t,
	}
}

func (s *getter) Get(key string) (string, error) {
	ctx, sp := s.tr.Start(context.Background(), "GetGetter")
	defer sp.End()

	value, err := s.st.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return value, nil
}
