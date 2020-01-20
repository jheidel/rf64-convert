package convert

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
)

type InputFile struct {
	Chunks     map[string][]byte
	DataSize   uint64
	DataReader io.Reader

	file *os.File
}

func OpenInput(path string) (*InputFile, error) {
	if path == "" {
		return nil, errors.New("input path required")
	}
	log.Printf("Loading input %q", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var header FileHeader
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	log.Printf("File is %v type %v", string(header.Magic[:]), string(header.Type[:]))

	reader := bufio.NewReader(file)

	cm := make(map[string][]byte)
	for {
		var chunk Chunk
		if err := binary.Read(reader, binary.LittleEndian, &chunk); err != nil {
			return nil, err
		}
		t := string(chunk.Type[:])
		log.Printf("Parsed chunk %q size %d", t, chunk.Size)

		if t == "data" {
			break
		}
		if chunk.Size == 0 {
			return nil, errors.New("Unexpected zero sized chunk")
		}

		cb := make([]byte, int(chunk.Size))
		_, err := reader.Read(cb)
		if err != nil {
			return nil, err
		}
		cm[t] = cb
	}

	// Get current offset into the file.
	dataOffset, err := file.Seek(0, 1)
	if err != nil {
		return nil, err
	}

	dataBytes := uint64(stat.Size()) - uint64(dataOffset)

	log.Printf("Offset from header is %d, file size is %d final is %d", dataOffset, stat.Size(), dataBytes)

	// TODO remove me.
	auxi, err := DecodeUTF16(cm["auxi"])
	if err != nil {
		return nil, err
	}
	log.Printf("auxi %q", auxi)

	return &InputFile{
		Chunks:     cm,
		DataSize:   dataBytes,
		DataReader: reader,

		file: file,
	}, nil
}

func (f *InputFile) Close() {
	f.file.Close()
}
