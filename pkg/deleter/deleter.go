package deleter

import (
	"context"
	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/storage"
)

type Deleter interface {
	Delete(string) error
}

type deleter struct {
	st storage.StorageRepo
	tr config.Tracer
}

func NewDeleter(s storage.StorageRepo, t config.Tracer) Deleter {
	return &deleter{
		st: s,
		tr: t,
	}
}

func (s *deleter) Delete(key string) error {
	ctx, sp := s.tr.Start(context.Background(), "DeleteDeleter")
	defer sp.End()

	err := s.st.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
