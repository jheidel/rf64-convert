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

	return &InputFile{
		Chunks:     cm,
		DataReader: reader,

		file: file,
	}, nil
}

func (f *InputFile) Close() {
	f.file.Close()
}
