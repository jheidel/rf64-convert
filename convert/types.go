package convert

type FileHeader struct {
	Magic [4]byte
	Size  uint32
	Type  [4]byte
}

type Chunk struct {
	Type [4]byte
	Size uint32
}

type DataSize64Chunk struct {
	RiffSize    uint64
	DataSize    uint64
	SampleCount uint64
	TableLength uint32 // zero, otherwise need table array
}

type FormatChunkPayload struct {
	FormatType     uint16
	ChannelCount   uint16
	SampleRate     uint32
	BytesPerSecond uint32
	BlockAlignment uint16
	BitsPerSample  uint16
}
