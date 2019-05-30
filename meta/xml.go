package meta

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"oops/system"
	"oops/util"
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
	dataDir = util.GetConfig().MetaDir
	os.Mkdir(dataDir, os.ModePerm)
}

// nameDirMapping 系统名和文件夹对应关系
var nameDirMapping map[string]string

// GetSystemNames 获取当前接管的系统名称，并缓存
func GetSystemNames() []string {
	file, _ := os.Open(dataDir)
	if file != nil {
		nameDirMapping = make(map[string]string)
		if names, _ := file.Readdirnames(-1); names != nil {
			for i, n := range names {
				xmlFileNames := getVersion(n)
				data := readData(n, xmlFileNames[0])
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

// readData 读取指定文件的配置信息
func readData(system string, xmlName string) []byte {
	systemDir := filepath.Join(dataDir, system)
	fileName := filepath.Join(systemDir, xmlName)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	return data
}

// getVersion 获取当前版本使用的配置文件文件，string[system.xml, service.xml]
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

// GetSystem 获取指定系统
func GetSystem(system string) *system.System {
	if nameDirMapping == nil {
		GetSystemNames()
		if nameDirMapping == nil {
			return nil
		}
	}

	system = nameDirMapping[system]
	xmlFileNames := getVersion(system)
	systemData := readData(system, xmlFileNames[0])
	serverData := readData(system, xmlFileNames[1])

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

// buildSystem 将xml配置数据转换成system
func buildSystem(xmlSystem *XMLSystem, serviceList *XMLServiceList) *system.System {
	sys := &system.System{}
	sys.Name = xmlSystem.Name
	sys.Version = xmlSystem.Version
	serviceMap := resolveService(xmlSystem, serviceList)
	resolveInstance(sys, serviceMap, xmlSystem)
	return sys
}

// resolveInstance 解析所有已经实例化的service
// serviceMap 所有服务映射表
// xmlSystem system中已经实例化的service
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

// 根据名字或者唯一标识获取service，identity无法获取，则查找属于name的service
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

// allUsedServer 根据依赖，将实例化的service，所有使用的服务名提取出来
func allUsedServer(ss []*system.Service, usedNames map[string]string) {
	for _, s := range ss {
		usedNames[serviceIdentity(*s)] = "1"
		for _, d := range s.Dependency {
			usedNames[serviceIdentity(*d)] = "1"
		}
		allUsedServer(s.Dependency, usedNames)
	}
}

// resolveService 将所有配置的service，转换成service映射表
func resolveService(sys *XMLSystem, serviceList *XMLServiceList) map[string]*system.Service {
	xmlServiceMap := make(map[string]*XMLService)
	for i := range serviceList.ServiceList {
		xmlServiceMap[xmlServiceIdentity(serviceList.ServiceList[i])] = &serviceList.ServiceList[i]
	}
	for i := range sys.Server {
		inst := &sys.Server[i]
		tpl := xmlServiceMap[xmlServiceIdentity(*inst)]
		if tpl == nil {
			if inst.Version == "" {
				for j, s := range serviceList.ServiceList {
					if s.Name == inst.Name {
						tpl = &serviceList.ServiceList[j]
					}
				}
			}
		}

		if tpl != nil { // copy instance to tpl
			tpl.Env = inst.Env
			if inst.Dependency != nil && len(inst.Dependency) != 0 {
				tpl.Dependency = append(tpl.Dependency, inst.Dependency...)
			}
			if inst.Operate != nil && len(inst.Operate) != 0 {
				tpl.Operate = append(tpl.Operate, inst.Operate...)
			}
			if inst.Protocol != nil && len(inst.Protocol) != 0 {
				tpl.Protocol = append(tpl.Protocol, inst.Protocol...)
			}
		} else { // add new instance to service list
			serviceList.ServiceList = append(serviceList.ServiceList, *inst)
		}
	}
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

// resolveDependency 查找并添加service的依赖
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
