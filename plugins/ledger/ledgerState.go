package ledger

import (
	"sync"

	"github.com/iotaledger/goshimmer/packages/ternary"
)

type LedgerState map[ternary.Trinary]*AddressEntry

type AddressEntry struct {
	accumulated      int64
	accumulatedMutex sync.RWMutex
	history          []*BalanceEntry
	historyMutex     sync.RWMutex
}

type BalanceEntry struct {
	value          int64
	valueMutex     sync.RWMutex
	timestamp      int64
	timestampMutex sync.RWMutex
}
