package system

import (
	"fmt"
	"testing"
)

func Test_Append(t *testing.T) {
	var v []string
	v = append(v, "123")
	fmt.Println(v)
}

func Test_Append2(t *testing.T) {
	ss := make([]*Service, 0)
	append2(&ss, 3, nil)
	fmt.Println(ss)
}

func append2(ss *[]*Service, i int, v *Service) {
	if len(*ss) < i+1 {
		for l := len(*ss); l < i; l++ {
			*ss = append(*ss, nil)
		}
	}
	*ss = append(*ss, v)
	fmt.Println(*ss)
}
