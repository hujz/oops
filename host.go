package main

import "encoding/xml"

type Host struct {
	XMLName  xml.Name `xml:"host"`
	VIP      []string `xml:"vip"`
	IP       []string `xml:"ip"`
	OS       string   `xml:"os"`
	Family   string   `xml:"via"`
	Protocol []Protocol
	Operate  []Operate
}
