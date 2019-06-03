package meta

import (
	"fmt"
	"testing"
)

func TestParseParams(t *testing.T) {
	m := map[string]string{
		"m": "123",
		"k": "hujj",
	}
	s := "$m'd'$k kk$"
	fmt.Println(ParseParams(s, m))
}
