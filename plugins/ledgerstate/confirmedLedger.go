package ledgerstate

import (
	"github.com/dgraph-io/badger"
	"github.com/iotaledger/goshimmer/packages/database"
	"github.com/iotaledger/goshimmer/packages/datastructure"
	"github.com/iotaledger/goshimmer/packages/errors"
	"github.com/iotaledger/goshimmer/packages/model/ledger/address"
	"github.com/iotaledger/goshimmer/packages/node"
	"github.com/iotaledger/goshimmer/packages/typeutils"
	"github.com/iotaledger/goshimmer/packages/unsafeconvert"
	"github.com/iotaledger/iota.go/trinary"
)

// region public api ///////////////////////////////////////////////////////////////////////////////////////////////////

func GetAddressEntry(addressShard trinary.Trytes, computeIfAbsent ...func(trinary.Trytes) *address.Entry) (result *address.Entry, err errors.IdentifiableError) {
	if cacheResult := ledgerEntryCache.ComputeIfAbsent(addressShard, func() interface{} {
		if addressEntry, dbErr := getAddressEntryFromDatabase(addressShard); dbErr != nil {
			err = dbErr

			return nil
		} else if addressEntry != nil {
			return addressEntry
		} else {
			if len(computeIfAbsent) >= 1 {
				return computeIfAbsent[0](addressShard)
			}

			return nil
		}
	}); !typeutils.IsInterfaceNil(cacheResult) {
		result = cacheResult.(*address.Entry)
	}

	return
}

func ContainsEntry(addressShard trinary.Trytes) (result bool, err errors.IdentifiableError) {
	if ledgerEntryCache.Contains(addressShard) {
		result = true
	} else {
		result, err = databaseContainsEntry(addressShard)
	}

	return
}

func StoreEntry(entry *address.Entry) {
	ledgerEntryCache.Set(entry.GetAddress()+entry.GetShard(), entry)
}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////

// region lru cache ////////////////////////////////////////////////////////////////////////////////////////////////////

var ledgerEntryCache = datastructure.NewLRUCache(LEDGER_ENTRY_CACHE_SIZE, &datastructure.LRUCacheOptions{
	EvictionCallback: onEvictEntry,
})

func onEvictEntry(_ interface{}, entry interface{}) {
	if evictedEntry := entry.(*address.Entry); evictedEntry.GetModified() {
		go func(evictedEntry *address.Entry) {
			if err := storeEntryInDatabase(evictedEntry); err != nil {
				panic(err)
			}
		}(evictedEntry)
	}
}

const (
	LEDGER_ENTRY_CACHE_SIZE = 50000
)

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////

// region database /////////////////////////////////////////////////////////////////////////////////////////////////////

var confirmedLedgerDatabase database.Database

func configureConfirmedLedgerDatabase(plugin *node.Plugin) {
	if db, err := database.Get("confirmedLedger"); err != nil {
		panic(err)
	} else {
		confirmedLedgerDatabase = db
	}
}

func storeEntryInDatabase(entry *address.Entry) errors.IdentifiableError {
	if entry.GetModified() {
		if err := confirmedLedgerDatabase.Set(unsafeconvert.StringToBytes(entry.GetAddress()+entry.GetShard()), entry.Marshal()); err != nil {
			return ErrDatabaseError.Derive(err, "failed to store address entry")
		}

		entry.SetModified(false)
	}

	return nil
}

func getAddressEntryFromDatabase(addressShard trinary.Trytes) (addressEntry *address.Entry, e errors.IdentifiableError) {
	addressEntryData, err := confirmedLedgerDatabase.Get(unsafeconvert.StringToBytes(addressShard))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		} else {
			return nil, ErrDatabaseError.Derive(err, "failed to retrieve address entry")
		}
	}

	addressEntry = &address.Entry{}
	e = addressEntry.Unmarshal(addressEntryData)
	if e != nil {
		return nil, e
	}
	return addressEntry, e
}

func databaseContainsEntry(addressShard trinary.Trytes) (bool, errors.IdentifiableError) {
	if contains, err := confirmedLedgerDatabase.Contains(unsafeconvert.StringToBytes(addressShard)); err != nil {
		return contains, ErrDatabaseError.Derive(err, "failed to check if the address entry exists")
	} else {
		return contains, nil
	}
}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////
