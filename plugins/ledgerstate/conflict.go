package ledgerstate

import (
	"sync"

	"github.com/iotaledger/goshimmer/packages/errors"
	"github.com/iotaledger/goshimmer/packages/model/value_transaction"
)

// mutex to keep conflict checks atomic
var conflictCheck sync.Mutex

func IsConflicting(tx *value_transaction.ValueTransaction) (bool, errors.IdentifiableError) {
	conflictCheck.Lock()
	defer conflictCheck.Unlock()

	addrEntry, err := getAddressEntryFromDatabase(tx.GetAddress() + tx.MetaTransaction.GetShardMarker())
	if err != nil {
		return false, err
	}
	// if there is an entry, return a conflict if the balance goes negative
	if addrEntry != nil {
		return addrEntry.GetBalance()+tx.GetValue() < 0, nil
	}
	// if there is no entry, return a conflict if the value is negative
	return tx.GetValue() < 0, nil
}
