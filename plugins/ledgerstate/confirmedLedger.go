package ledgerstate

import (
	"github.com/dgraph-io/badger"
	"github.com/iotaledger/goshimmer/packages/database"
	"github.com/iotaledger/goshimmer/packages/errors"
	"github.com/iotaledger/goshimmer/packages/model/ledger/address"
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

func getAddressEntryFromDatabase(addressShard ternary.Trytes) (addressEntry *address.Entry, e errors.IdentifiableError) {
	addressEntryData, err := confirmedLedgerDatabase.Get(addressShard.CastToBytes())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		} else {
			return nil, ErrDatabaseError.Derive(err, "failed to retrieve transaction")
		}
	}

	addressEntry = &address.Entry{}
	e = addressEntry.Unmarshal(addressEntryData)
	if e != nil {
		return nil, e
	}
	return addressEntry, e
}

func databaseContainsAddressEntry(addressShard ternary.Trytes) (bool, errors.IdentifiableError) {
	if contains, err := confirmedLedgerDatabase.Contains(addressShard.CastToBytes()); err != nil {
		return contains, ErrDatabaseError.Derive(err, "failed to check if the transaction exists")
	} else {
		return contains, nil
	}
}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////
