package ledgerstate

import (
	"testing"

	"github.com/iotaledger/goshimmer/packages/model/ledger/address"
	"github.com/iotaledger/goshimmer/packages/model/ledger/balance"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/magiconair/properties/assert"
)

func TestConfirmedLedgerDB(t *testing.T) {
	configureConfirmedLedgerDatabase(nil)

	addr := trinary.Trytes("A9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999F")
	shardMarker := trinary.Trytes("NPHTQORL9XK")
	addressShard := address.New(addr, shardMarker)

	balanceEntries := []*balance.Entry{balance.NewValue(100, 1), balance.NewValue(100, 2)}

	addressShard.Add(balanceEntries...)

	err := storeEntryInDatabase(addressShard)
	if err != nil {
		t.Error(err)
	}
	addressShardFromDB, err := getAddressEntryFromDatabase(addressShard.GetAddress() + addressShard.GetShard())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, addressShardFromDB.GetAddress(), addressShard.GetAddress(), "Address")
	assert.Equal(t, addressShardFromDB.GetShard(), addressShard.GetShard(), "Shard")
	assert.Equal(t, addressShardFromDB.GetBalance(), addressShard.GetBalance(), "Accumulated")
}

func TestConfirmedLedgerAPI(t *testing.T) {
	configureConfirmedLedgerDatabase(nil)

	addr := trinary.Trytes("A9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999F")
	shardMarker := trinary.Trytes("NPHTQORL9XK")
	addressShard := address.New(addr, shardMarker)

	balanceEntries := []*balance.Entry{balance.NewValue(100, 1), balance.NewValue(100, 2)}

	addressShard.Add(balanceEntries...)

	StoreEntry(addressShard)

	addressShardFromCache, err := GetAddressEntry(addressShard.GetAddress() + addressShard.GetShard())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, addressShardFromCache.GetAddress(), addressShard.GetAddress(), "Address")
	assert.Equal(t, addressShardFromCache.GetShard(), addressShard.GetShard(), "Shard")
	assert.Equal(t, addressShardFromCache.GetBalance(), addressShard.GetBalance(), "Accumulated")
}
