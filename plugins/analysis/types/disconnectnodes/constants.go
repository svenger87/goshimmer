package disconnectnodes

const (
	MARSHALLED_PACKET_HEADER = 0x04

	MARSHALLED_PACKET_HEADER_START = 0
	MARSHALLED_PACKET_HEADER_SIZE  = 1
	MARSHALLED_PACKET_HEADER_END   = MARSHALLED_PACKET_HEADER_START + MARSHALLED_PACKET_HEADER_SIZE

	MARSHALLED_SOURCE_ID_START = MARSHALLED_PACKET_HEADER_END
	MARSHALLED_SOURCE_ID_SIZE  = 20
	MARSHALLED_SOURCE_ID_END   = MARSHALLED_SOURCE_ID_START + MARSHALLED_SOURCE_ID_SIZE

	MARSHALLED_TARGET_ID_START = MARSHALLED_SOURCE_ID_END
	MARSHALLED_TARGET_ID_SIZE  = 20
	MARSHALLED_TARGET_ID_END   = MARSHALLED_TARGET_ID_START + MARSHALLED_TARGET_ID_SIZE

	MARSHALLED_TOTAL_SIZE = MARSHALLED_TARGET_ID_END
)
