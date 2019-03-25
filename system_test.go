package main

import (
	"fmt"
	"testing"
)

func Test_BuildSystem(t *testing.T) {
	system := BuildSystem("link")
	if system != nil {
		for _, d := range system.LaunchNode.Dependency {
			fmt.Println(d.Server.Name)
		}
	}
}

func TestSystem_Start(t *testing.T) {
	system := BuildSystem("link")
	system.Start()
}

func TestServer_Start(t *testing.T) {
	server := Service{Name: "test"}
	server.Start()
	fmt.Println(server.Status)
}

