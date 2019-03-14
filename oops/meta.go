package oops

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Dependency struct {
	XMLName xml.Name `xml:"server"`
	Server  string   `xml:",chardata"`
}

type Operate struct {
	XMLName  xml.Name `xml:"operate"`
	Name     string   `xml:"name,attr"`
	Protocol string   `xml:"protocol,attr"`
	Argument string   `xml:",chardata"`
}

type Host struct {
	XMLName xml.Name `xml:"host"`
	Vip     string   `xml:"vip"`
	IP      []string `xml:"ip"`
	Os      string   `xml:"os"`
	Via     string   `xml:"via"`
}

type Protocol struct {
	XMLName xml.Name `xml:"protocol"`
	Name    string   `xml:"name,attr"`
	URI     string   `xml:",chardata"`
}

type Server struct {
	XMLName    xml.Name     `xml:"server"`
	Type       string       `xml:"type,attr"`
	Name       string       `xml:"name,attr"`
	Operate    []Operate    `xml:"operate"`
	Protocol   []Protocol   `xml:"protocol"`
	Host       Host         `xml:"host"`
	Dependency []Dependency `xml:"dependency>server"`
}
type System struct {
	XMLName xml.Name `xml:"system"`
	Version string   `xml:"version,attr"`
	Name    string   `xml:"name,attr"`
	Server  []Server `xml:"server"`
}

type DependTree struct {
	Server     Server
	Dependency []DependTree
	Depended   []DependTree
}

var dataDir string

func init() {
	dataDir = os.Getenv("oops_data_dir")
	if dataDir == "" {
		dataDir = os.Getenv("HOME") + "/.oops/"
	}
	err := os.Mkdir(dataDir, os.ModePerm)
	fmt.Println(err)
}

func BuildSystem(system string) (*System, *DependTree, *DependTree) {
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
		return nil, nil, nil
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	systemCache := System{}
	xml.Unmarshal(data, &systemCache)

	return &systemCache, nil, nil
}

func resolveDepend(system System) *DependTree {
	serverMap := make(map[string]Server)
	dependTreeMap := make(map[string]DependTree)
	for _, s := range system.Server {
		serverMap[s.Name] = s
		dependTreeMap[s.Name] = DependTree{Server: s}
	}
	return nil
}
