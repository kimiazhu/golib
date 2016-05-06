// Description: util
// Author: ZHU HAIHUA
// Since: 2016-05-06 13:31
package util

import (
	"os"
	"math"
	"crypto/md5"
	"io"
	"encoding/hex"
)

func MD5File(filepath string) string {
	chunk := float64(8192)
	file, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	info, _ := file.Stat()
	size := info.Size()
	blocks := int64(math.Ceil(float64(size) / float64(chunk)))
	hash := md5.New()

	for i := int64(0); i < blocks; i++ {
		blocksize := int(math.Min(chunk, float64(size-int64(float64(i)*chunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
