package address

import (
	//"fmt"
	"fmt"
	"testing"

	"github.com/iotaledger/goshimmer/packages/model/ledger/balance"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/magiconair/properties/assert"
)

func TestAddress_SettersGetters(t *testing.T) {
	address := trinary.Trytes("A9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999F")
	shardMarker := trinary.Trytes("NPHTQORL9XK")
	addressShard := New(address + shardMarker)

	balanceEntries := []*balance.Entry{balance.NewValue(100, 1), balance.NewValue(100, 2)}

	addressShard.Add(balanceEntries...)
	assert.Equal(t, addressShard.GetAddressShard(), address+shardMarker, "AddressShard")
	assert.Equal(t, addressShard.GetBalance(), int64(200), "Accumulated")
}

func TestBalance_MarshalUnmarshalGetters(t *testing.T) {
	address := trinary.Trytes("A9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999F")
	shardMarker := trinary.Trytes("NPHTQORL9XK")
	addressShard := New(address + shardMarker)

	balanceEntries := []*balance.Entry{balance.NewValue(100, 1), balance.NewValue(100, 2)}

	addressShard.Add(balanceEntries...)

	addressShardByte := addressShard.Marshal()
	var addressShardUnmarshaled Entry
	err := addressShardUnmarshaled.Unmarshal(addressShardByte)
	if err != nil {
		fmt.Println(err, len(addressShardByte))
	}

	assert.Equal(t, addressShardUnmarshaled.GetAddressShard(), addressShard.GetAddressShard(), "AddressShard")
	assert.Equal(t, addressShardUnmarshaled.GetBalance(), addressShard.GetBalance(), "Accumulated")
}
