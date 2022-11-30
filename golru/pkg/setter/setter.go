package setter

import (
	"context"
	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/storage"
)

type Setter interface {
	Set(string, []byte) error
}

type setter struct {
	st storage.StorageRepo
	tr config.Tracer
}

func NewSetter(s storage.StorageRepo, t config.Tracer) Setter {
	return &setter{
		st: s,
		tr: t,
	}
}

func (s *setter) Set(key string, value []byte) error {
	ctx, sp := s.tr.Start(context.Background(), "SetSetter")
	defer sp.End()

	err := s.st.Set(ctx, key, string(value))
	if err != nil {
		return err
	}
	return nil
}
