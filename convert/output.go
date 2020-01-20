package convert

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"time"
)

type OutputFile struct {
	file *os.File
}

func (f *OutputFile) Close() {
	f.file.Close()
}

func NewOutputFile(path string) (*OutputFile, error) {
	if path == "" {
		return nil, errors.New("output path required")
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	header := &FileHeader{
		Magic: [4]byte{'R', 'F', '6', '4'},
		Size:  0xFFFFFFFF,
		Type:  [4]byte{'W', 'A', 'V', 'E'},
	}

	if err := binary.Write(file, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	log.Printf("Created output file %q", path)
	return &OutputFile{
		file: file,
	}, nil
}

func parseFmtChunk(b []byte) (*FormatChunkPayload, error) {
	chunk := &FormatChunkPayload{}
	if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, chunk); err != nil {
		return nil, err
	}
	return chunk, nil
}

type countingReader struct {
	Reader io.Reader
	Count  uint64
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.Count += uint64(n)
	return n, err
}

func (f *OutputFile) CopyFrom(in *InputFile) error {
	// We could add 'auxi' to replicate SDRConsole's internal format.
	//   auxi, err := EncodeUTF16(TestAuxi)
	//   if err != nil {
	//   	return err
	//   }
	//   in.Chunks["auxi"] = auxi

	// Compute total expected file size
	fs := uint64(binary.Size(&FileHeader{}))
	fs += uint64(binary.Size(&Chunk{}) + binary.Size(&DataSize64Chunk{}))
	for _, chunk := range in.Chunks {
		fs += uint64(binary.Size(&Chunk{}) + len(chunk))
	}
	fs += uint64(binary.Size(&Chunk{})) // data chunk header

	fmt, err := parseFmtChunk(in.Chunks["fmt "])
	if err != nil {
		return err
	}

	log.Printf("Writing chunk \"ds64\" (template)")
	ds64h := &Chunk{
		Type: [4]byte{'d', 's', '6', '4'},
		Size: uint32(binary.Size(&DataSize64Chunk{})),
	}
	if err := binary.Write(f.file, binary.LittleEndian, ds64h); err != nil {
		return err
	}
	// Keep track of position of header so we can change it later.
	ds64Pos, err := f.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		return err
	}
	if err := binary.Write(f.file, binary.LittleEndian, &DataSize64Chunk{}); err != nil {
		return err
	}

	// Write all chunks.
	for name, chunk := range in.Chunks {
		log.Printf("Writing chunk %q", name)

		ch := &Chunk{
			Size: uint32(len(chunk)),
		}
		copy(ch.Type[:], []byte(name))

		if err := binary.Write(f.file, binary.LittleEndian, ch); err != nil {
			return err
		}
		if _, err := f.file.Write(chunk); err != nil {
			return err
		}
	}

	log.Printf("Writing chunk \"data\"")
	dc := &Chunk{
		Type: [4]byte{'d', 'a', 't', 'a'},
		Size: 0xFFFFFFFF,
	}
	if err := binary.Write(f.file, binary.LittleEndian, dc); err != nil {
		return err
	}

	log.Printf("Copying data chunk. This might take a while...")
	cr := &countingReader{
		Reader: in.DataReader,
	}
	if _, err := io.Copy(f.file, cr); err != nil {
		return err
	}
	log.Printf("Finished copy.")

	// Go back and edit ds64 chunk with corrected size information
	fs += cr.Count
	samples := cr.Count / uint64(fmt.BlockAlignment)
	durSec := float64(samples) / float64(fmt.SampleRate)
	dur := time.Duration(durSec) * time.Second
	log.Printf("WAV contains %d samples (%d bytes per sample), duration %v", samples, fmt.BlockAlignment, dur)

	log.Printf("Finalizing chunk \"ds64\"")
	if _, err := f.file.Seek(ds64Pos, os.SEEK_SET); err != nil {
		return err
	}
	ds64 := &DataSize64Chunk{
		RiffSize:    fs,
		DataSize:    cr.Count,
		SampleCount: samples,
	}
	if err := binary.Write(f.file, binary.LittleEndian, ds64); err != nil {
		return err
	}

	return nil
}
