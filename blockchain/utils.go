package blockchain

import (
	"bytes"
	"encoding/binary"
	"log"
	//"fmt"
)

// IntToHex converts an int64 to a byte array
func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// IsHashValid checks if a hash is valid
func isHashValid(hash []byte, difficulty int) bool {
	//fmt.Printf("isHashValid: %x, %d\n", hash, difficulty)
	prefix := bytes.Repeat([]byte{0}, difficulty)
	return bytes.HasPrefix(hash, prefix)
}
