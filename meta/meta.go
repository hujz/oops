package meta

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type XMLSystem struct {
	XMLName  xml.Name      `xml:"system"`
	Version  string        `xml:"version,attr"`
	Name     string        `xml:"name,attr"`
	Server   []XMLService  `xml:"server"`
	Instance []XMLInstance `xml:"instance"`
}

type XMLHostList struct {
	XMLName  xml.Name  `xml:"host"`
	HostList []XMLHost `xml:"host"`
}

type XMLHost struct {
	XMLName  xml.Name      `xml:"host"`
	Name     string        `xml:"name,attr"`
	VIP      []string      `xml:"vip"`
	IP       []string      `xml:"ip"`
	OS       string        `xml:"os"`
	Family   string        `xml:"family"`
	Virt     string        `xml:"virt"`
	Protocol []XMLProtocol `xml:"protocol"`
	Operate  []XMLOperate  `xml:"operate"`
}

type XMLServiceList struct {
	XMLName     xml.Name     `xml:"service"`
	ServiceList []XMLService `xml:"service"`
}

type XMLService struct {
	XMLName    xml.Name      `xml:"service"`
	Spec       string        `xml:"spec,attr"`
	Version    string        `xml:"version,attr"`
	Name       string        `xml:"name,attr"`
	Operate    []XMLOperate  `xml:"operate"`
	Protocol   []XMLProtocol `xml:"protocol"`
	Dependency []string      `xml:"dependency>service"`
}

type XMLProtocol struct {
	XMLName xml.Name `xml:"protocol"`
	Name    string   `xml:"name,attr"`
	URI     string   `xml:",chardata"`
}

type XMLOperate struct {
	XMLName  xml.Name `xml:"operate"`
	Name     string   `xml:"name,attr"`
	Protocol string   `xml:"protocol,attr"`
	Argument string   `xml:",chardata"`
}

type XMLSpecHost struct {
	XMLName xml.Name     `xml:"host"`
	Family  string       `xml:"family,attr"`
	Operate []XMLOperate `xml:"operate"`
}
type XMLSpec struct {
	XMLName xml.Name      `xml:"spec"`
	Host    []XMLSpecHost `xml:"host"`
}

type XMLInstance struct {
	XMLName xml.Name `xml:"instance"`
	ID      string   `xml:"id,attr"`
	Service string   `xml:"service,attr"`
	Version string   `xml:"version,attr"`
	Host    []string `xml:"host"`
}

var dataDir string

func init() {
	dataDir = os.Getenv("oops_data_dir")
	if dataDir == "" {
		dataDir = os.Getenv("HOME") + "/.oops/"
	}
	os.Mkdir(dataDir, os.ModePerm)
}

func GetSystemNames() []string {
	file, _ := os.Open(dataDir)
	if file != nil {
		if names, _ := file.Readdirnames(-1); names != nil {
			return names
		}
	}
	return nil
}

func readData(system string, index int) []byte {
	systemDir := filepath.Join(dataDir, system)

	ver := getVersion(system)

	fileName := filepath.Join(systemDir, ver[index])
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	return data
}

func getVersion(system string) []string {
	systemDir := filepath.Join(dataDir, system)
	ver := filepath.Join(systemDir, "version")
	_, err := os.Stat(ver)

	if err == nil {
		data, err := ioutil.ReadFile(ver)
		if err == nil {
			str := string(data)
			return strings.Split(str, "\n")
		}
	}
	return []string{"system.xml", "service.xml", "host.xml"}
}

func BuildSystem(system string) *System {
	systemData := readData(system, 0)
	serverData := readData(system, 1)
	hostData := readData(system, 2)

	if systemData == nil || serverData == nil || hostData == nil {
		return nil
	} else {
		xmlSystem, xmlServer, xmlHost := XMLSystem{}, XMLServiceList{}, XMLHostList{}
		xml.Unmarshal(systemData, xmlSystem)
		xml.Unmarshal(serverData, xmlServer)
		xml.Unmarshal(hostData, xmlHost)
		return buildSystem(xmlSystem, xmlServer, xmlHost)
	}
	return nil
}

func buildSystem(xmlSystem XMLSystem, serviceList XMLServiceList, hostList XMLHostList) *System {
	sys := &System{}
	sys.Name = xmlSystem.Name
	sys.Version = xmlSystem.Version
	fillHost(sys, hostList)
	fillService(sys, serviceList)
	fillInstance(sys, xmlSystem)
	resolveDependency(sys)
	return sys
}

func fillInstance(sys *System, system XMLSystem) {

}

func fillService(sys *System, serviceList XMLServiceList) {
	serviceMap := make(map[string]Service)
	for _, s := range serviceList.ServiceList {
		service := Service{Name: s.Name, Version: s.Version, Spec: s.Spec}
		serviceMap[s.Name] = service
		protocolMap := make(map[string]Protocol)
		for _, p := range s.Protocol {
			protocolMap[p.Name] = Protocol{Name: p.Name, URI: p.URI}
		}
		operateMap := make(map[string]Operate)
		for _, o := range s.Operate {
			operateMap[o.Name] = Operate{Name: o.Name, Protocol: o.Protocol, Argument: o.Argument}
		}
		service.Protocol = protocolMap
		service.Operate = operateMap
	}
	sys.Service = serviceMap
}

func resolveDependency(sys *System) {

}

func fillHost(sys *System, hostList XMLHostList) {
	hostMap := make(map[string]Host)
	for _, h := range hostList.HostList {
		host := Host{Name: h.Name, OS: h.OS, Family: h.Family, Virt: h.Virt, VIP: h.VIP, IP: h.IP}
		protocolMap := make(map[string]Protocol)
		for _, p := range h.Protocol {
			protocolMap[p.Name] = Protocol{Name: p.Name, URI: p.URI}
		}
		operateMap := make(map[string]Operate)
		for _, o := range h.Operate {
			operateMap[o.Name] = Operate{Name: o.Name, Protocol: o.Protocol, Argument: o.Argument}
		}
		host.Protocol = protocolMap
		host.Operate = operateMap
		hostMap[h.Name] = host
	}
	sys.Host = hostMap
}
