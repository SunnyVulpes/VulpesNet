package utils

import "io"

type Header struct {
	MagicNumber int
}

type Codec interface {
	io.Closer
	ReadHeader(header *Header) error
	ReadBody(body interface{}) error
	Write(header *Header, body interface{}) error
	GetConn() io.ReadWriteCloser
}

const (
	NewSSH = iota
	RegisterService
	ReverseConn
)

//const Host = "localhost"

const Host = "47.92.71.42"
