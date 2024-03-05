package utils

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type JSONCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder
	enc  *json.Encoder
}

func (c *JSONCodec) ReadHeader(header *Header) error {
	return c.dec.Decode(header)
}

func (c *JSONCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *JSONCodec) Write(header *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.conn.Close()
		}
	}()
	err = c.enc.Encode(header)
	err = c.enc.Encode(body)
	if err != nil {
		log.Println("codec error: ", err)
		return err
	}
	return nil
}

func (c *JSONCodec) Close() error {
	return c.conn.Close()
}

func (c *JSONCodec) GetConn() io.ReadWriteCloser {
	return c.conn
}

func InitJSONCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JSONCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf),
	}
}
