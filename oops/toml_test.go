package oops

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml"
)

func TestToml(t *testing.T) {
	tree, err := toml.Load(
		`name="mysql-5.7"
		[start]
		cmd="service mysql start"
		[status]
		cmd="mysql -uroot -pLink$2013 -e select now()"
		[stop]
		cmd="service mysql stop"
		[check]
		cmd=""
		`)
	if err != nil {
		fmt.Println(err)
	} else {
		config := tree.ToMap()
		for k, v := range config {
			fmt.Println(k, "=", v)
		}
	}
}
