package ledgerstate

import (
	"github.com/dgraph-io/badger"
	"github.com/iotaledger/goshimmer/packages/database"
	"github.com/iotaledger/goshimmer/packages/errors"
	"github.com/iotaledger/goshimmer/packages/model/ledger/address"
	"github.com/iotaledger/goshimmer/packages/model/value_transaction"
	"github.com/iotaledger/goshimmer/packages/node"
	"github.com/iotaledger/goshimmer/packages/ternary"
)

// region database /////////////////////////////////////////////////////////////////////////////////////////////////////

var confirmedLedgerDatabase database.Database

func configureConfirmedLedgerDatabase(plugin *node.Plugin) {
	if db, err := database.Get("confirmedLedger"); err != nil {
		panic(err)
	} else {
		confirmedLedgerDatabase = db
	}
}

func storeAddressEntryInDatabase(entry *address.Entry) errors.IdentifiableError {
	if entry.GetModified() {
		if err := confirmedLedgerDatabase.Set(entry.GetAddressShard().CastToBytes(), entry.Marshal()); err != nil {
			return ErrDatabaseError.Derive(err, "failed to store address entry")
		}

		entry.SetModified(false)
	}

	return nil
}

func getAddressEntryFromDatabase(addressShard ternary.Trytes) (*address.Entry, errors.IdentifiableError) {
	txData, err := transactionDatabase.Get(transactionHash.CastToBytes())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		} else {
			return nil, ErrDatabaseError.Derive(err, "failed to retrieve transaction")
		}
	}

	return value_transaction.FromBytes(txData), nil
}

// func databaseContainsAddressEntry(transactionHash ternary.Trytes) (bool, errors.IdentifiableError) {
// 	if contains, err := transactionDatabase.Contains(transactionHash.CastToBytes()); err != nil {
// 		return contains, ErrDatabaseError.Derive(err, "failed to check if the transaction exists")
// 	} else {
// 		return contains, nil
// 	}
// }

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////
