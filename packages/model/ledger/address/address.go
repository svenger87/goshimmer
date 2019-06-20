package address

import (
	"encoding/binary"
	"strconv"
	"sync"

	"github.com/iotaledger/goshimmer/packages/errors"
	"github.com/iotaledger/goshimmer/packages/model/ledger/balance"
	"github.com/iotaledger/goshimmer/packages/ternary"
	"github.com/iotaledger/goshimmer/packages/typeutils"
)

type Entry struct {
	addressShard ternary.Trytes
	accumulated  int64
	history      []*balance.Entry
	historyMutex sync.RWMutex
	modified     bool
}

func New(addressShard ternary.Trytes) *Entry {
	return &Entry{
		addressShard: addressShard,
		modified:     false,
	}
}

func (addressEntry *Entry) GetAddressShard() (result ternary.Trytes) {
	result = addressEntry.addressShard
	return
}

func (addressEntry *Entry) GetBalance() (result int64) {
	addressEntry.historyMutex.RLock()
	defer addressEntry.historyMutex.RUnlock()

	result = addressEntry.accumulated
	return
}

func (addressEntry *Entry) GetModified() bool {
	addressEntry.historyMutex.RLock()
	defer addressEntry.historyMutex.RUnlock()
	return addressEntry.modified
}

func (addressEntry *Entry) SetModified(modified bool) {
	addressEntry.modified = modified
}

func (addressEntry *Entry) Add(balanceEntries ...*balance.Entry) {
	addressEntry.historyMutex.Lock()
	defer addressEntry.historyMutex.Unlock()

	addressEntry.history = append(addressEntry.history, balanceEntries...)

	for _, balanceEntry := range balanceEntries {
		addressEntry.accumulated += balanceEntry.GetValue()
	}
	addressEntry.modified = true
}

func (addressEntry *Entry) Marshal() (result []byte) {
	result = make([]byte, MARSHALED_ENTRY_MIN_SIZE+len(addressEntry.history)*MARSHALED_BALANCE_ENTRY_SIZE)

	addressEntry.historyMutex.RLock()

	binary.BigEndian.PutUint64(result[MARSHALED_ENTRY_HISTORY_COUNT_START:MARSHALED_ENTRY_HISTORY_COUNT_END], uint64(len(addressEntry.history)))

	binary.BigEndian.PutUint64(result[MARSHALED_ENTRY_ACCUMULATED_START:MARSHALED_ENTRY_ACCUMULATED_END], uint64(addressEntry.accumulated))

	copy(result[MARSHALED_ENTRY_ADDRESS_START:MARSHALED_ENTRY_ADDRESS_END], addressEntry.addressShard.CastToBytes())

	i := 0
	for _, balanceEntry := range addressEntry.history {
		var BALANCE_START = MARSHALED_ENTRY_HISTORY_START + i*(MARSHALED_BALANCE_ENTRY_SIZE)
		var BALANCE_END = BALANCE_START + MARSHALED_BALANCE_ENTRY_SIZE

		copy(result[BALANCE_START:BALANCE_END], balanceEntry.Marshal())

		i++
	}

	addressEntry.historyMutex.RUnlock()
	return
}

func (addressEntry *Entry) Unmarshal(data []byte) (err errors.IdentifiableError) {
	dataLen := len(data)

	if dataLen <= MARSHALED_ENTRY_MIN_SIZE {
		return ErrMarshallFailed.Derive(errors.New("unmarshall failed"), "marshaled address are too short")
	}

	balancesCount := binary.BigEndian.Uint64(data[MARSHALED_ENTRY_HISTORY_COUNT_START:MARSHALED_ENTRY_HISTORY_COUNT_END])

	if dataLen < MARSHALED_ENTRY_MIN_SIZE+int(balancesCount)*MARSHALED_BALANCE_ENTRY_SIZE {
		return ErrMarshallFailed.Derive(errors.New("unmarshall failed"), "marshaled address are too short for "+strconv.FormatUint(balancesCount, 10)+" balances")
	}

	addressEntry.historyMutex.Lock()

	addressEntry.addressShard = ternary.Trytes(typeutils.BytesToString(data[MARSHALED_ENTRY_ADDRESS_START:MARSHALED_ENTRY_ADDRESS_END]))
	addressEntry.accumulated = int64(binary.BigEndian.Uint64(data[MARSHALED_ENTRY_ACCUMULATED_START:MARSHALED_ENTRY_ACCUMULATED_END]))
	addressEntry.history = make([]*balance.Entry, balancesCount)
	for i := uint64(0); i < balancesCount; i++ {
		var BALANCE_START = MARSHALED_ENTRY_HISTORY_START + i*(MARSHALED_BALANCE_ENTRY_SIZE)
		var BALANCE_END = BALANCE_START + MARSHALED_BALANCE_ENTRY_SIZE

		addressEntry.history[i] = balance.New()
		err = addressEntry.history[i].Unmarshal(data[BALANCE_START:BALANCE_END])
	}

	addressEntry.historyMutex.Unlock()
	return
}
