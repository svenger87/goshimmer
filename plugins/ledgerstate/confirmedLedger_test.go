package ledgerstate

import (
	"testing"

	"github.com/iotaledger/goshimmer/packages/model/ledger/address"
	"github.com/iotaledger/goshimmer/packages/model/ledger/balance"
	"github.com/iotaledger/goshimmer/packages/ternary"
	"github.com/magiconair/properties/assert"
)

func TestConfirmedLedgerDB(t *testing.T) {
	configureConfirmedLedgerDatabase(nil)

	addr := ternary.Trytes("A999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999F")
	shardMarker := ternary.Trytes("NPHTQORL9XKA")
	addressShard := address.New(addr + shardMarker)

	balanceEntries := []*balance.Entry{balance.NewValue(100, 1), balance.NewValue(100, 2)}

	addressShard.Add(balanceEntries...)

	err := storeAddressEntryInDatabase(addressShard)
	if err != nil {
		t.Error(err)
	}
	addressShardFromDB, err := getAddressEntryFromDatabase(addressShard.GetAddressShard())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, addressShardFromDB.GetAddressShard().CastToBytes(), addressShard.GetAddressShard().CastToBytes(), "AddressShard")
	assert.Equal(t, addressShardFromDB.GetBalance(), addressShard.GetBalance(), "Accumulated")
}
