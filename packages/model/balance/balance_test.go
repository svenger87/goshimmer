package balance

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestBalance_SettersGetters(t *testing.T) {
	balance := New()
	balance.SetValue(int64(^uint64(0) >> 1))
	balance.SetTimestamp(^uint64(0))

	assert.Equal(t, balance.GetValue(), int64(^uint64(0)>>1), "MaxValue")
	assert.Equal(t, balance.GetTimestamp(), ^uint64(0), "MaxTimestamp")
}

func TestBalance_MarshalUnmarshalGetters(t *testing.T) {
	balance := New()
	balance.SetValue(int64(^uint64(0) >> 1))
	balance.SetTimestamp(^uint64(0))

	balanceByte := balance.Marshal()
	balanceUnmarshaled := New()
	err := balanceUnmarshaled.Unmarshal(balanceByte)
	if err != nil {
		fmt.Println(err, len(balanceByte))
	}
	assert.Equal(t, balanceUnmarshaled.GetValue(), balance.GetValue(), "Value")
	assert.Equal(t, balanceUnmarshaled.GetTimestamp(), balance.GetTimestamp(), "Timestamp")
}
