package utils

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var Convert = &convertUtil{
	defaultGBKDecoder: simplifiedchinese.GBK.NewDecoder(),
}

type convertUtil struct {
	// 默认gbk解码器
	defaultGBKDecoder *encoding.Decoder
}

func (u *convertUtil) GBKToUTF8(s string) (string, error) {
	return u.defaultGBKDecoder.String(s)
}
func (u *convertUtil) GBKToUTF8Bytes(s []byte) ([]byte, error) {
	return u.defaultGBKDecoder.Bytes(s)
}
