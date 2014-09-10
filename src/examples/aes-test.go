package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
)

// verified to be identical openssl's -aes-256-cbc with:
// create 256MB random file:
// dd if=/dev/urandom of=test count=16384 bs=16384
//
// encrypt file
// openssl enc -nopad -aes-256-cfb -K 129e12fd1f9e2b31129e12fd1f9e2b31129e12fd1f9e2b31129e12fd1f9e2b31 -iv 000102030405060708090a0b0c0d0e0f -e -in test -out test-openssl.aes
//
// run go:
// $ go run aes-test.go test test-go.aes
//
// Make sure openssl and aes-test.go output matches:
// $ sha256sum test-openssl.aes test-go.aes
// 7320cc72c78ddf395964b43a03425cc6323ac700223a3210c3af6b09d58b3e67 test-openssl.aes
// 7320cc72c78ddf395964b43a03425cc6323ac700223a3210c3af6b09d58b3e67 test-go.aes
// $ cmp test1.aes test1-go.aes
// $

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

const bufferSize = 1024

func main() {
	var count int
	var size int
	readBuffer := make([]byte, bufferSize)
	cipherBuffer := make([]byte, bufferSize)

	// Setup a key that will encrypt the other text.
	key_text := []byte{0x12, 0x9e, 0x12, 0xfd, 0x1f, 0x9e, 0x2b, 0x31, 0x12, 0x9e, 0x12, 0xfd, 0x1f, 0x9e, 0x2b, 0x31, 0x12, 0x9e, 0x12, 0xfd, 0x1f, 0x9e, 0x2b, 0x31, 0x12, 0x9e, 0x12, 0xfd, 0x1f, 0x9e, 0x2b, 0x31}
	fmt.Printf("Len = %d Key= %x\n", len(key_text), key_text)
	fmt.Printf("len = %d iv = %x\n", len(commonIV), commonIV)

	// We chose our cipher type here in this case
	// we are using AES.
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
		os.Exit(-1)
	}

	plainFile, err := os.Open(os.Args[1])
	reader := bufio.NewReader(plainFile)

	cipherFile, err := os.Create(os.Args[2])
	writer := bufio.NewWriter(cipherFile)

	cfb := cipher.NewCFBEncrypter(c, commonIV)
	size = 0
	for {
		if count, err = reader.Read(readBuffer); err != nil {
			size = size + count
			break
		}
		size = size + count
		cfb.XORKeyStream(cipherBuffer[:count], readBuffer[:count])
		writer.Write(cipherBuffer[:count])
	}
	fmt.Printf("count=%d\n", count)
	fmt.Printf("size=%d\n", size)
	writer.Flush()
	plainFile.Close()
	cipherFile.Close()
}
