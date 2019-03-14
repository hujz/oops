package oops

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"encoding/hex"
	"encoding/xml"
)

func Test_XML(t *testing.T) {
	file, err := os.Open("/opt/github/golang/src/oops/data/system.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	data, _ := ioutil.ReadAll(file)
	sys := System{}
	xml.Unmarshal(data, &sys)

	fmt.Println(sys.Server[1])

	data, _ = xml.MarshalIndent(sys, "", "\t")
	fmt.Println(string(data))
}

func Test_TestListen(t *testing.T) {
	fmt.Println(os.Getenv("HOME"))
}

func Test_Encrypt(t *testing.T) {
	p := []byte("123456781234567812345678")
	e, _ := DesEncrypt([]byte("hujz:123"), p)
	es := hex.EncodeToString(e)
	fmt.Println(es)
	data, _ := hex.DecodeString(es)
	d, _ := DesDecrypt(data, p)
	fmt.Println(string(d))
}
