package balance

import (
	"encoding/binary"

	"github.com/iotaledger/goshimmer/packages/errors"
)

type Entry struct {
	value     int64
	timestamp uint64
}

func New() *Entry {
	return &Entry{}
}

func NewValue(value int64, timestamp uint64) *Entry {
	return &Entry{
		value:     value,
		timestamp: timestamp,
	}
}

func (balanceEntry *Entry) GetValue() (result int64) {
	result = balanceEntry.value

	return
}

func (balanceEntry *Entry) GetTimestamp() (result uint64) {
	result = balanceEntry.timestamp

	return
}

func (balanceEntry *Entry) GetModified() bool {
	return true
}

func (balanceEntry *Entry) SetModified(modified bool) {
}

func (balanceEntry *Entry) Marshal() (result []byte) {
	result = make([]byte, MARSHALED_ENTRY_VALUE_SIZE+MARSHALED_ENTRY_TIMESTAMP_SIZE)

	binary.BigEndian.PutUint64(result[MARSHALED_ENTRY_VALUE_START:MARSHALED_ENTRY_VALUE_END], uint64(balanceEntry.value))

	binary.BigEndian.PutUint64(result[MARSHALED_ENTRY_TIMESTAMP_START:MARSHALED_ENTRY_TIMESTAMP_END], uint64(balanceEntry.timestamp))

	return
}

func (balanceEntry *Entry) Unmarshal(data []byte) (err errors.IdentifiableError) {
	dataLen := len(data)

	if dataLen < MARSHALED_ENTRY_VALUE_SIZE+MARSHALED_ENTRY_TIMESTAMP_SIZE {
		return ErrMarshallFailed.Derive(errors.New("unmarshall failed"), "marshaled balanceEntry are too short")
	}

	balanceEntry.value = int64(binary.BigEndian.Uint64(data[MARSHALED_ENTRY_VALUE_START:MARSHALED_ENTRY_VALUE_END]))

	balanceEntry.timestamp = binary.BigEndian.Uint64(data[MARSHALED_ENTRY_TIMESTAMP_START:MARSHALED_ENTRY_TIMESTAMP_END])

	return
}
