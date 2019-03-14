package oops

import (
	"fmt"
	"testing"
)

func Test_BuildSystem(t *testing.T) {
	system := BuildSystem("link")
	if system != nil {
		fmt.Println(system.Name)
	}
}
