package inmem

import (
	"testing"

	"github.com/boolw/go-web3/tracker/store"
)

func TestInMemoryStore(t *testing.T) {
	store.TestStore(t, func(t *testing.T) (store.Store, func()) {
		return NewInmemStore(), func() {}
	})
}
