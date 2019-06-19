package balance

import (
	"sync"

	"github.com/iotaledger/goshimmer/packages/errors"
	"github.com/iotaledger/goshimmer/packages/ternary"
)

type Entry struct {
	value          int64
	valueMutex     sync.RWMutex
	timestamp      uint64
	timestampMutex sync.RWMutex
	modified       bool
}

func New(hash ternary.Trytes) *Entry {
	return &Entry{}
}

func (balanceEntry *Entry) Marshal() (result []byte) {
	return
}

func (balanceEntry *Entry) Unmarshal(data []byte) (err errors.IdentifiableError) {
	return
}
