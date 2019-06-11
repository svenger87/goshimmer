package removenode

const (
	MARSHALLED_PACKET_HEADER = 0x02

	MARSHALLED_PACKET_HEADER_START = 0
	MARSHALLED_PACKET_HEADER_SIZE  = 1
	MARSHALLED_PACKET_HEADER_END   = MARSHALLED_PACKET_HEADER_START + MARSHALLED_PACKET_HEADER_SIZE

	MARSHALLED_ID_START = MARSHALLED_PACKET_HEADER_END
	MARSHALLED_ID_SIZE  = 20
	MARSHALLED_ID_END   = MARSHALLED_ID_START + MARSHALLED_ID_SIZE

	MARSHALLED_TOTAL_SIZE = MARSHALLED_ID_END
)
