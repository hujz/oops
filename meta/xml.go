package meta

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"oops/protocol"
	"os"
	"path/filepath"
	"strings"
)

type XMLSystem struct {
	XMLName xml.Name     `xml:"system"`
	Version string       `xml:"version,attr"`
	Name    string       `xml:"name,attr"`
	Server  []XMLService `xml:"server"`
}

type XMLServiceList struct {
	XMLName     xml.Name     `xml:"service"`
	ServiceList []XMLService `xml:"service"`
}

type XMLService struct {
	XMLName    xml.Name      `xml:"service"`
	Version    string        `xml:"version,attr"`
	Name       string        `xml:"name,attr"`
	Env        string        `xml:env,chardata`
	Operate    []XMLOperate  `xml:"operate"`
	Protocol   []XMLProtocol `xml:"protocol"`
	Dependency []XMLService  `xml:"dependency>service"`
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
	return []string{"system.xml", "service.xml"}
}

func BuildSystem(system string) *System {
	systemData := readData(system, 0)
	serverData := readData(system, 1)

	if systemData == nil || serverData == nil {
		return nil
	} else {
		xmlSystem, xmlServer := XMLSystem{}, XMLServiceList{}
		xml.Unmarshal(systemData, xmlSystem)
		xml.Unmarshal(serverData, xmlServer)
		return buildSystem(&xmlSystem, &xmlServer)
	}
	return nil
}

func buildSystem(xmlSystem *XMLSystem, serviceList *XMLServiceList) *System {
	sys := &System{}
	sys.Name = xmlSystem.Name
	sys.Version = xmlSystem.Version
	serviceMap := resolveService(xmlSystem, serviceList)
	resolveInstance(sys, serviceMap, xmlSystem)
	return sys
}

func resolveInstance(sys *System, serviceMap map[string]*Service, system *XMLSystem) {
	var usedServerList []*Service
	for _, s := range system.Server {
		usedServerList = append(usedServerList, serviceMap[xmlServiceIdentity(s)])
	}
	allUsedServerNames := make(map[string]string)
	allUsedServer(usedServerList, allUsedServerNames)
	usedServerList = []*Service{}
	for n := range allUsedServerNames {
		usedServerList = append(usedServerList, serviceMap[n])
	}
	sys.Service = usedServerList
}

func allUsedServer(ss []*Service, usedNames map[string]string) {
	for _, s := range ss {
		for _, d := range s.Dependency {
			usedNames[serviceIdentity(*d)] = "1"
		}
		allUsedServer(s.Dependency, usedNames)
	}
}

func resolveService(sys *XMLSystem, serviceList *XMLServiceList) map[string]*Service {
	serviceMap := make(map[string]*Service)
	for _, s := range serviceList.ServiceList {
		service := Service{Name: s.Name, Version: s.Version, Depth: 1, Env: s.Env}
		serviceMap[s.Name+" "+s.Version] = &service
		protocolMap := make(map[string]protocol.IProtocol)
		for _, p := range s.Protocol {
			protocolMap[p.Name] = protocol.BuildProtocol(Protocol{Name: p.Name, URI: p.URI})
		}
		operateMap := make(map[string]*Operate)
		for _, o := range s.Operate {
			operateMap[o.Name] = &Operate{Name: o.Name, Protocol: protocolMap[o.Protocol], Argument: o.Argument}
		}
		service.Protocol = protocolMap
		service.Operate = operateMap
	}
	resolveDependency(serviceMap, serviceList)
	return serviceMap
}

func resolveDependency(serviceMap map[string]*Service, xmlServiceList *XMLServiceList) {
	for _, xs := range xmlServiceList.ServiceList {
		for _, xd := range xs.Dependency {
			serviceMap[xmlServiceIdentity(xs)].Dependency = append(serviceMap[xmlServiceIdentity(xs)].Dependency, serviceMap[xmlServiceIdentity(xd)])
		}
	}
}

func xmlServiceIdentity(service XMLService) string {
	return service.Name + " " + service.Version
}
func serviceIdentity(service Service) string {
	return service.Name + " " + service.Version
}

//func fillHost(sys *System, hostList XMLHostList) {
//	hostMap := make(map[string]Host)
//	for _, h := range hostList.HostList {
//		host := Host{Name: h.Name, OS: h.OS, Family: h.Family, Virt: h.Virt, VIP: h.VIP, IP: h.IP}
//		protocolMap := make(map[string]Protocol)
//		for _, p := range h.Protocol {
//			protocolMap[p.Name] = Protocol{Name: p.Name, URI: p.URI}
//		}
//		operateMap := make(map[string]Operate)
//		for _, o := range h.Operate {
//			operateMap[o.Name] = Operate{Name: o.Name, Protocol: o.Protocol, Argument: o.Argument}
//		}
//		host.Protocol = protocolMap
//		host.Operate = operateMap
//		hostMap[h.Name] = host
//	}
//	sys.Host = hostMap
//}
