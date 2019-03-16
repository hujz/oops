package main

import "encoding/xml"

type Operate struct {
	XMLName  xml.Name `xml:"operate"`
	Name     string   `xml:"name,attr"`
	Protocol string   `xml:"protocol,attr"`
	Argument string   `xml:",chardata"`
}
