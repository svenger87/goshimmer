package address

import (
	"encoding/binary"
	"strconv"
	"sync"

	"github.com/iotaledger/goshimmer/packages/errors"
	"github.com/iotaledger/goshimmer/packages/model/ledger/balance"
	"github.com/iotaledger/goshimmer/packages/typeutils"
	"github.com/iotaledger/goshimmer/packages/unsafeconvert"
	"github.com/iotaledger/iota.go/trinary"
)

type Entry struct {
	address       trinary.Trytes
	shard         trinary.Trytes
	accumulated   int64
	history       []*balance.Entry
	historyMutex  sync.RWMutex
	modified      bool
	modifiedMutex sync.RWMutex
}

func New(address, shard trinary.Trytes) *Entry {
	return &Entry{
		address: address,
		shard:   shard,
		history: make([]*balance.Entry, 0),
	}
}

func (addressEntry *Entry) GetAddress() (result trinary.Trytes) {
	result = addressEntry.address

	return
}

func (addressEntry *Entry) GetShard() (result trinary.Trytes) {
	result = addressEntry.shard

	return
}

func (addressEntry *Entry) GetBalance() (result int64) {
	addressEntry.historyMutex.RLock()
	result = addressEntry.accumulated
	addressEntry.historyMutex.RUnlock()

	return
}

func (addressEntry *Entry) GetModified() (result bool) {
	addressEntry.modifiedMutex.RLock()
	result = addressEntry.modified
	addressEntry.modifiedMutex.RUnlock()

	return
}

func (addressEntry *Entry) SetModified(modified bool) {
	addressEntry.modifiedMutex.Lock()
	addressEntry.modified = modified
	addressEntry.modifiedMutex.Unlock()
}

func (addressEntry *Entry) Add(balanceEntries ...*balance.Entry) {
	addressEntry.historyMutex.Lock()

	addressEntry.history = append(addressEntry.history, balanceEntries...)
	for _, balanceEntry := range balanceEntries {
		addressEntry.accumulated += balanceEntry.GetValue()
	}
	addressEntry.SetModified(true)

	addressEntry.historyMutex.Unlock()
}

func (addressEntry *Entry) Marshal() (result []byte) {
	result = make([]byte, MARSHALED_ENTRY_MIN_SIZE+len(addressEntry.history)*MARSHALED_BALANCE_ENTRY_SIZE)

	addressEntry.historyMutex.RLock()

	binary.BigEndian.PutUint64(result[MARSHALED_ENTRY_HISTORY_COUNT_START:MARSHALED_ENTRY_HISTORY_COUNT_END], uint64(len(addressEntry.history)))

	binary.BigEndian.PutUint64(result[MARSHALED_ENTRY_ACCUMULATED_START:MARSHALED_ENTRY_ACCUMULATED_END], uint64(addressEntry.accumulated))

	copy(result[MARSHALED_ENTRY_ADDRESS_START:MARSHALED_ENTRY_ADDRESS_END], unsafeconvert.StringToBytes(addressEntry.address))
	copy(result[MARSHALED_ENTRY_SHARD_START:MARSHALED_ENTRY_SHARD_END], unsafeconvert.StringToBytes(addressEntry.shard))

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

	addressEntry.address = trinary.Trytes(typeutils.BytesToString(data[MARSHALED_ENTRY_ADDRESS_START:MARSHALED_ENTRY_ADDRESS_END]))
	addressEntry.shard = trinary.Trytes(typeutils.BytesToString(data[MARSHALED_ENTRY_SHARD_START:MARSHALED_ENTRY_SHARD_END]))
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
