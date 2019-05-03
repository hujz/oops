package meta

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Host(t *testing.T) {
	file, _ := os.Open("/opt/github/golang/src/oops/data/link/host.xml")
	data, _ := ioutil.ReadAll(file)
	host := XMLHostList{}
	err := xml.Unmarshal(data, &host)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(host)

	host.HostList[0].IP[0] = "324234"
	str, _ := xml.MarshalIndent(host, "", "\t")
	fmt.Println(string(str))
}

func Test_Service(t *testing.T) {
	file, _ := os.Open("/opt/github/golang/src/oops/data/link/service.xml")
	data, _ := ioutil.ReadAll(file)
	serviceList := &XMLServiceList{}
	err := xml.Unmarshal(data, serviceList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*serviceList)
}

func Test_Spec(t *testing.T) {
	file, _ := os.Open("/opt/github/golang/src/oops/data/spec.xml")
	data, _ := ioutil.ReadAll(file)
	serviceList := &XMLSpec{}
	err := xml.Unmarshal(data, serviceList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*serviceList)
}

func Test_Instance(t *testing.T) {
	file, _ := os.Open("/opt/github/golang/src/oops/data/link/system.xml")
	data, _ := ioutil.ReadAll(file)
	serviceList := &XMLSystem{}
	err := xml.Unmarshal(data, serviceList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*serviceList)
}