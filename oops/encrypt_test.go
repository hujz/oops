package oops

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_Encrypt(t *testing.T) {
	p := []byte("123456781234567812345678")
	e, _ := DesEncrypt([]byte("hujz:123"), p)
	es := hex.EncodeToString(e)
	fmt.Println(es)
	data, _ := hex.DecodeString(es)
	d, _ := DesDecrypt(data, p)
	fmt.Println(string(d))
}
