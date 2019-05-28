package main

import (
	"fmt"
	"os"
	"testing"
)

func Test_XML(t *testing.T) {
	//file, err := os.Open("/opt/github/golang/src/oops/data/system.xml")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//data, _ := ioutil.ReadAll(file)
	//sys := System{}
	//xml.Unmarshal(data, &sys)
	//
	//fmt.Println(sys.Server[1])
	//
	//data, _ = xml.MarshalIndent(sys, "", "\t")
	//fmt.Println(string(data))
}

func Test_TestListen(t *testing.T) {
	fmt.Println(os.Getenv("HOME"))
}

func Test_Test(t *testing.T) {
	str := 2
	switch str {
	case 2:
	case 1:
		fmt.Println(11)
	default:
		fmt.Println(00)
	}
}
