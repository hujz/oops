package xml

import "encoding/xml"

type System struct {
	XMLName  xml.Name   `xml:"system"`
	Version  string     `xml:"version,attr"`
	Name     string     `xml:"name,attr"`
	Server   []Service  `xml:"server"`
	Instance []Instance `xml:"instance"`
}

type HostList struct {
	XMLName  xml.Name `xml:"host"`
	HostList []Host   `xml:"host"`
}

type Host struct {
	XMLName  xml.Name   `xml:"host"`
	Name     string     `xml:"name,attr"`
	VIP      []string   `xml:"vip"`
	IP       []string   `xml:"ip"`
	OS       string     `xml:"os"`
	Family   string     `xml:"family"`
	Virt     string     `xml:"virt"`
	Protocol []Protocol `xml:"protocol"`
	Operate  []Operate  `xml:"operate"`
}

type ServiceList struct {
	XMLName     xml.Name  `xml:"service"`
	ServiceList []Service `xml:"service"`
}

type Service struct {
	XMLName    xml.Name   `xml:"service"`
	Spec       string     `xml:"spec,attr"`
	Version    string     `xml:"version,attr"`
	Name       string     `xml:"name,attr"`
	Operate    []Operate  `xml:"operate"`
	Protocol   []Protocol `xml:"protocol"`
	Dependency []string   `xml:"dependency>service"`
}

type Protocol struct {
	XMLName xml.Name `xml:"protocol"`
	Name    string   `xml:"name,attr"`
	URI     string   `xml:",chardata"`
}

type Operate struct {
	XMLName  xml.Name `xml:"operate"`
	Name     string   `xml:"name,attr"`
	Protocol string   `xml:"protocol,attr"`
	Argument string   `xml:",chardata"`
}

type SpecHost struct {
	XMLName xml.Name  `xml:"host"`
	Family  string    `xml:"family,attr"`
	Operate []Operate `xml:"operate"`
}
type Spec struct {
	XMLName xml.Name   `xml:"spec"`
	Host    []SpecHost `xml:"host"`
}

type Instance struct {
	XMLName xml.Name `xml:"instance"`
	Name    string   `xml:"name,attr"`
	Version string   `xml:"version,attr"`
	Host    []string `xml:"host"`
}
