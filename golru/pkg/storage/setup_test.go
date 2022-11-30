package storage

import (
	"os"
	"testing"
)

var storageT StorageRepo

var storeT = make(map[string]*node)

func TestMain(m *testing.M) {

	// storageT = NewStorage(50)
	storageT = &storage{
		store: storeT,
		dll:   NewDll(50),
	}

	os.Exit(m.Run())
}
