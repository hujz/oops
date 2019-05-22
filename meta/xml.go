package meta

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"oops/system"
	"os"
	"path/filepath"
	"strings"
)

type XMLSystem struct {
	XMLName xml.Name     `xml:"system"`
	Version string       `xml:"version,attr"`
	Name    string       `xml:"name,attr"`
	Server  []XMLService `xml:"service"`
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

var nameDirMapping = make(map[string]string)

func GetSystemNames() []string {
	file, _ := os.Open(dataDir)
	if file != nil {
		if names, _ := file.Readdirnames(-1); names != nil {
			for i, n := range names {
				data := readData(n, 0)
				sys := XMLSystem{}
				xml.Unmarshal(data, &sys)
				nameDirMapping[sys.Name] = n
				names[i] = sys.Name
			}
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

func GetSystemMeta(system string) *system.System {
	system = nameDirMapping[system]
	systemData := readData(system, 0)
	serverData := readData(system, 1)

	if systemData == nil || serverData == nil {
		return nil
	} else {
		xmlSystem, xmlServer := XMLSystem{}, XMLServiceList{}
		xml.Unmarshal(systemData, &xmlSystem)
		xml.Unmarshal(serverData, &xmlServer)
		return buildSystem(&xmlSystem, &xmlServer)
	}
	return nil
}

func buildSystem(xmlSystem *XMLSystem, serviceList *XMLServiceList) *system.System {
	sys := &system.System{}
	sys.Name = xmlSystem.Name
	sys.Version = xmlSystem.Version
	serviceMap := resolveService(xmlSystem, serviceList)
	resolveInstance(sys, serviceMap, xmlSystem)
	return sys
}

func resolveInstance(sys *system.System, serviceMap map[string]*system.Service, xmlSystem *XMLSystem) {
	var usedServerList []*system.Service
	for _, s := range xmlSystem.Server {
		usedServerList = append(usedServerList, getService(serviceMap, s.Name, xmlServiceIdentity(s)))
	}
	allUsedServerNames := make(map[string]string)
	allUsedServer(usedServerList, allUsedServerNames)
	usedServerList = []*system.Service{}
	for n := range allUsedServerNames {
		usedServerList = append(usedServerList, serviceMap[n])
	}
	sys.Service = usedServerList
}

func getService(serviceMap map[string]*system.Service, name, identity string) *system.Service {
	s := serviceMap[identity]
	if s == nil {
		for k, v := range serviceMap {
			if strings.HasPrefix(k, name) {
				return v
			}
		}
	}
	return s
}

func allUsedServer(ss []*system.Service, usedNames map[string]string) {
	for _, s := range ss {
		usedNames[serviceIdentity(*s)] = "1"
		for _, d := range s.Dependency {
			usedNames[serviceIdentity(*d)] = "1"
		}
		allUsedServer(s.Dependency, usedNames)
	}
}

func resolveService(sys *XMLSystem, serviceList *XMLServiceList) map[string]*system.Service {
	serviceMap := make(map[string]*system.Service)
	for _, s := range serviceList.ServiceList {
		service := system.Service{Name: s.Name, Version: s.Version, Depth: 1, Env: s.Env}
		serviceMap[xmlServiceIdentity(s)] = &service
		protocolMap := make(map[string]system.IProtocol)
		for _, p := range s.Protocol {
			protocolMap[p.Name] = system.BuildProtocol(system.Protocol{Name: p.Name, URI: p.URI})
		}
		operateMap := make(map[string]*system.Operate)
		for _, o := range s.Operate {
			operateMap[o.Name] = &system.Operate{Name: o.Name, Protocol: protocolMap[o.Protocol], Argument: o.Argument}
		}
		service.Protocol = protocolMap
		service.Operate = operateMap
	}
	resolveDependency(serviceMap, serviceList)
	return serviceMap
}

func resolveDependency(serviceMap map[string]*system.Service, xmlServiceList *XMLServiceList) {
	for _, xs := range xmlServiceList.ServiceList {
		for _, xd := range xs.Dependency {
			serviceMap[xmlServiceIdentity(xs)].Dependency = append(serviceMap[xmlServiceIdentity(xs)].Dependency, getService(serviceMap, xd.Name, xmlServiceIdentity(xd)))
		}
	}
}

func xmlServiceIdentity(service XMLService) string {
	if service.Version == "" {
		return service.Name
	}
	return service.Name + " " + service.Version
}

func serviceIdentity(service system.Service) string {
	if service.Version == "" {
		return service.Name
	}
	return service.Name + " " + service.Version
}
