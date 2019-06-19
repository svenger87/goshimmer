package balance

import (
	"encoding/binary"
	"sync"

	"github.com/iotaledger/goshimmer/packages/errors"
)

type Entry struct {
	value          int64
	valueMutex     sync.RWMutex
	timestamp      uint64
	timestampMutex sync.RWMutex
	modified       bool
}

func New() *Entry {
	return &Entry{}
}

// region public methods with locking //////////////////////////////////////////////////////////////////////////////////

func (balanceEntry *Entry) GetValue() (result int64) {
	balanceEntry.valueMutex.RLock()
	result = balanceEntry.value
	balanceEntry.valueMutex.RUnlock()

	return
}

func (balanceEntry *Entry) SetValue(value int64) {
	balanceEntry.valueMutex.Lock()
	balanceEntry.value = value
	balanceEntry.valueMutex.Unlock()

	balanceEntry.modified = true
}

func (balanceEntry *Entry) GetTimestamp() (result uint64) {
	balanceEntry.timestampMutex.RLock()
	result = balanceEntry.timestamp
	balanceEntry.timestampMutex.RUnlock()

	return
}

func (balanceEntry *Entry) SetTimestamp(timestamp uint64) {
	balanceEntry.timestampMutex.Lock()
	balanceEntry.timestamp = timestamp
	balanceEntry.timestampMutex.Unlock()

	balanceEntry.modified = true
}

func (balanceEntry *Entry) GetModified() bool {
	return true
}

func (balanceEntry *Entry) SetModified(modified bool) {
}

func (balanceEntry *Entry) Marshal() (result []byte) {
	result = make([]byte, MARSHALED_ENTRY_VALUE_SIZE+MARSHALED_ENTRY_TIMESTAMP_SIZE)

	balanceEntry.valueMutex.RLock()
	binary.BigEndian.PutUint64(result[MARSHALED_ENTRY_VALUE_START:MARSHALED_ENTRY_VALUE_END], uint64(balanceEntry.value))
	balanceEntry.valueMutex.RUnlock()

	balanceEntry.timestampMutex.RLock()
	binary.BigEndian.PutUint64(result[MARSHALED_ENTRY_TIMESTAMP_START:MARSHALED_ENTRY_TIMESTAMP_END], uint64(balanceEntry.timestamp))
	balanceEntry.timestampMutex.RUnlock()

	return
}

func (balanceEntry *Entry) Unmarshal(data []byte) (err errors.IdentifiableError) {
	dataLen := len(data)

	if dataLen < MARSHALED_ENTRY_VALUE_SIZE+MARSHALED_ENTRY_TIMESTAMP_SIZE {
		return ErrMarshallFailed.Derive(errors.New("unmarshall failed"), "marshaled balanceEntry are too short")
	}

	balanceEntry.valueMutex.Lock()
	balanceEntry.value = int64(binary.BigEndian.Uint64(data[MARSHALED_ENTRY_VALUE_START:MARSHALED_ENTRY_VALUE_END]))
	balanceEntry.valueMutex.Unlock()

	balanceEntry.timestampMutex.Lock()
	balanceEntry.timestamp = binary.BigEndian.Uint64(data[MARSHALED_ENTRY_TIMESTAMP_START:MARSHALED_ENTRY_TIMESTAMP_END])
	balanceEntry.timestampMutex.Unlock()

	return
}
