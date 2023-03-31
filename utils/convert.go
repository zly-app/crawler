package utils

import (
	"encoding/base64"
	"reflect"
	"unsafe"

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

// string转bytes, 转换后的bytes禁止写, 否则产生运行故障
func (*convertUtil) StringToBytes(s *string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// bytes转string
func (*convertUtil) BytesToString(b []byte) *string {
	return (*string)(unsafe.Pointer(&b))
}

func (u *convertUtil) Base64Encode(text string) string {
	temp := u.StringToBytes(&text)
	bs := u.Base64EncodeBytes(temp)
	return *u.BytesToString(bs)
}

func (*convertUtil) Base64EncodeBytes(text []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(text)))
	base64.StdEncoding.Encode(buf, text)
	return buf
}
