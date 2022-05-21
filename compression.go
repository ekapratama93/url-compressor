package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"

	"github.com/andybalholm/brotli"
)

func compress(url string) string {
	out := bytes.Buffer{}
	writer := brotli.NewWriterOptions(&out, brotli.WriterOptions{Quality: 11})
	writer.Write([]byte(url))
	writer.Close()
	compressedData := out.Bytes()
	return base64.StdEncoding.EncodeToString([]byte(compressedData))
}

func decompress(base64String string) string {
	decodedByte, _ := base64.StdEncoding.DecodeString(base64String)
	reader := brotli.NewReader(bytes.NewReader(decodedByte))
	res, _ := ioutil.ReadAll(reader)
	return string(res)
}
