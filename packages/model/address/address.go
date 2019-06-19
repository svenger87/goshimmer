package address

import (
	"sync"

	"github.com/iotaledger/goshimmer/packages/errors"
	"github.com/iotaledger/goshimmer/packages/model/balance"
	"github.com/iotaledger/goshimmer/packages/ternary"
)

type Entry struct {
	addressShard      ternary.Trytes
	addressShardMutex sync.RWMutex
	accumulated       int64
	accumulatedMutex  sync.RWMutex
	history           []balance.Entry
	historyMutex      sync.RWMutex
	modified          bool
}

func New(addressShard ternary.Trytes) *Entry {
	return &Entry{
		addressShard: addressShard,
		modified:     false,
	}
}

func (addressEntry *Entry) Marshal() (result []byte) {
	return
}

func (addressEntry *Entry) Unmarshal(data []byte) (err errors.IdentifiableError) {
	return
}

func getAddressFromDatabase(addressShard ternary.Trytes) (*Entry, errors.IdentifiableError) {
	return nil, nil
}

func databaseContainsAddress(addressShard ternary.Trytes) (bool, errors.IdentifiableError) {
	return false, nil
}
