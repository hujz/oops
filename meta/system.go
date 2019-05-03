package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ServerStatus string

const (
	Status_Starting ServerStatus = "starting"
	Status_Running  ServerStatus = "running"
	Status_Downing  ServerStatus = "downing"
	Status_Down     ServerStatus = "down"
)

type Dependency struct {
	XMLName xml.Name `xml:"server"`
	Server  string   `xml:",chardata"`
}

type System struct {
	XMLName    xml.Name  `xml:"system"`
	Version    string    `xml:"version,attr"`
	Name       string    `xml:"name,attr"`
	Server     []Service `xml:"server"`
	LaunchNode *DependTree
}

type DependTree struct {
	Server     *Service
	Dependency []*DependTree
	RefCount   int
}

func (server *Service) Start() string {
	if server.Status == Status_Running {
		return "ok"
	}
	res := server.Invoke("start")
	if res == "ok" {
		server.Status = Status_Running
	}
	return res
}

func (server *Service) Stop() string {
	if server.Status == Status_Running {
		res := server.Invoke("stop")
		if res == "ok" {
			server.Status = Status_Down
		}
		return res
	}
	return "ok"
}

func (server *Service) State() string {
	return server.Invoke("status")
}

func (server *Service) Invoke(operate string) string {
	log.Println(server.Name, ": --> ", operate)
	var op Operate
	var pt Protocol
	for _, p := range server.Operate {
		if p.Name == operate {
			op = p
			break
		}
	}
	if op.Name == "" {
		return "not support operate"
	}
	for _, p := range server.Protocol {
		if p.Name == op.Protocol {
			pt = p
			break
		}
	}
	protocol := BuildProtocol(pt.URI)

	if err := protocol.Open(); err == nil {
		defer protocol.Close()
	}
	protocol.Invoke(op.Argument)
	return "ok"
}

func (system *System) Start() {
	system.recursion(system.LaunchNode, func(server *Service) string {
		return server.Start()
	})
}

func (system *System) Stop() {
	system.recursion(system.LaunchNode, func(server *Service) string {
		return server.Stop()
	})
}

func (system *System) Status() {
	system.recursion(system.LaunchNode, func(server *Service) string {
		return server.State()
	})
}

func (system *System) recursion(treeNode *DependTree, call func(server *Service) string) string {
	if len(treeNode.Dependency) == 0 {
		if res := call(treeNode.Server); res != "ok" {
			return res
		}
	} else {
		for n := range treeNode.Dependency {
			system.recursion(treeNode.Dependency[n], call)
		}
		if treeNode.Server == nil {
			return "ok"
		}
		if res := call(treeNode.Server); res != "ok" {
			return res
		}
	}
	return "ok"
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

func BuildSystem(system string) *System {
	systemDir := filepath.Join(dataDir, system)

	ver := filepath.Join(systemDir, "version")
	_, err := os.Stat(ver)

	fileName := "system.xml"
	if err == nil {
		data, err := ioutil.ReadFile(ver)
		if err == nil {
			fileName = string(data)
		}
	}
	fileName = filepath.Join(systemDir, fileName)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	systemCache := System{}
	xml.Unmarshal(data, &systemCache)

	resolveDepend(&systemCache)

	return &systemCache
}

func resolveDepend(system *System) *DependTree {
	serverMap := make(map[string]*Service)
	dependTreeMap := make(map[string]*DependTree)
	for n := range system.Server {
		pServer := &system.Server[n]
		serverMap[pServer.Name] = pServer
		dependTreeMap[pServer.Name] = &DependTree{Server: pServer, RefCount: 0, Dependency: make([]*DependTree, 0)}
	}
	dependencyRoot := &DependTree{Dependency: make([]*DependTree, 0)}
	for n, s := range serverMap { // every server
		dependNode := dependTreeMap[n]    // treeNode of server
		for _, ds := range s.Dependency { // every dependency of this serer
			refTreeNode := dependTreeMap[ds.Server]
			dependNode.Dependency = append(dependNode.Dependency, refTreeNode) // append dependency treeNode to this sever's treeNode
			refTreeNode.RefCount += 1                                          // refCount ++
		}
	}
	for _, t := range dependTreeMap {
		if t.RefCount == 0 {
			dependencyRoot.Dependency = append(dependencyRoot.Dependency, t)
		}
	}

	system.LaunchNode = dependencyRoot
	return dependencyRoot
}
