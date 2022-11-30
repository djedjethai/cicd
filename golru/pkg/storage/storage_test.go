package storage

import (
	"context"
	"testing"
)

func TestPut(t *testing.T) {
	ctx := context.Background()

	_ = storageT.Set(ctx, "test", "put")

	_, ok := storeT["test"]

	if !ok {
		t.Error("err in store Put() failed")
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	dt, _ := storageT.Get(ctx, "test")

	if dt != "put" {
		t.Error("err in store Get() failed")
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	_ = storageT.Delete(ctx, "test")

	_, ok := storeT["test"]

	if ok {
		t.Error("err in store Delete() failed")
	}
}
