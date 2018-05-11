package main

import (
	"encoding/binary"
	"log"
	"time"
	"crypto/rand"
)

func cryptoReadInt64() int64 {
	if i, err := binary.ReadVarint(cryptoByteReader{}); err != nil {
		log.Println(err, "falling back to insecure seed using nanosecond time")
		return time.Now().UnixNano()
	} else {
		return i
	}
}

type cryptoByteReader struct{}

func (_ cryptoByteReader) ReadByte() (byte, error) {
	var onebyte [1]byte
	_, err := rand.Read(onebyte[:])
	return onebyte[0], err
}
