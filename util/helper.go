package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
)

func EnDecode(key []byte, input string) ([]byte, error) {
	// cipher gives you back a block
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}

	// iv := make([]byte, aes.BlockSize)

	// _, err = io.ReadFull(rand.Reader, iv)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf(err)
	// }

	// salt is the second param
	// returns a stream
	stream := cipher.NewCTR(b, key)

	buff := &bytes.Buffer{}

	steamWriter := cipher.StreamWriter{
		S: stream,
		W: buff,
	}

	// as it writes to the buffer it simultaneously encrypts the data
	_, err = steamWriter.Write([]byte(input))
	if err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}

	return buff.Bytes(), nil

}

// EncryptWriter make a wrapper around a writer
func EncryptWriter(wtr io.Writer, key []byte) (io.Writer, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}

	stream := cipher.NewCTR(b, key)

	return cipher.StreamWriter{
		S: stream,
		W: wtr,
	}, nil
}