package messageCompressor_test

import (
	"compress/flate"
	"github.com/bhbosman/gomessageblock"
	"testing"
)

func TestName(t *testing.T) {
	rws := gomessageblock.NewReaderWriter()
	compress, _ := flate.NewWriter(rws, flate.DefaultCompression)
	compress.Write([]byte("sdfdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsd"))
	compress.Flush()
	println(rws.Size())

	decompress := flate.NewReader(rws)
	b := [4096]byte{}
	decompress.Read(b[:])
	println(string(b[:]))

}
