package main

import "encoding/xml"

type Service struct {
	XMLName    xml.Name     `xml:"server"`
	Spec       string       `xml:"spec,attr"`
	Version    string       `xml:"version,attr"`
	Name       string       `xml:"name,attr"`
	Operate    []Operate    `xml:"operate"`
	Protocol   []Protocol   `xml:"protocol"`
	Host       Host         `xml:"host"`
	Dependency []Dependency `xml:"dependency>server"`
	Status     ServerStatus
}
