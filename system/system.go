package system

import (
	"io"
	"strconv"
)

var SystemCache = make(map[string]*System)

type System struct {
	Version     string
	Name        string
	Service     []*Service
	TopService  []*Service
	LowService  []*Service
	LevelMatrix [][]*Service
	build       bool
}

func Get(name string) *System {
	return SystemCache[name]
}

func (sys *System) Cache() {
	sys.Build()
	SystemCache[sys.Name] = sys
}

// Build 计算实例化的和所有依赖的service，得出最上层最下层的服务列表和服务依赖矩阵
func (system *System) Build() {
	if system.build {
		return
	}
	serviceMap := make(map[string]*Service)
	for _, s := range system.Service {
		serviceMap[s.Name+" "+s.Version] = s
		if s.Dependency == nil {
			system.LowService = append(system.LowService, s)
		} else {
			for _, d := range s.Dependency {
				d.Reference = append(d.Reference, s)
			}
		}
	}
	recursion(&system.Service, &system.TopService, 1)
	for _, s := range system.Service {
		setLength(&system.LevelMatrix, s.Depth+1)
		system.LevelMatrix[s.Depth] = append(system.LevelMatrix[s.Depth], s)
	}
	system.build = true
}

func setLength(ss *[][]*Service, newLen int) {
	if len(*ss) < newLen {
		for l := len(*ss); l < newLen; l++ {
			*ss = append(*ss, nil)
		}
	}
}

// recursion 递归依赖，计算出每个service的最大深度
func recursion(ss *[]*Service, topService *[]*Service, depth int) {
	for _, s := range *ss {
		if s.Reference == nil {
			*topService = append(*topService, s)
		}
		if s.Depth < depth {
			s.Depth = depth
		}
		recursion(&s.Dependency, topService, depth+1)
	}
}

func (sys *System) Start(input io.Reader, output io.Writer) {
	for l := len(sys.LevelMatrix) - 1; l > 0; l-- {
		for j := len(sys.LevelMatrix[l]) - 1; j >= 0; j-- {
			s := sys.LevelMatrix[l][j]
			output.Write([]byte(strconv.Itoa(l) + " start " + s.Name + "(" + s.Version + "): "))
			s.Invoke(Operate_Start, input, output)
		}
	}
}

func (sys *System) Stop(input io.Reader, output io.Writer) {
	for l1, i := len(sys.LevelMatrix), 1; i < l1; i++ {
		for l2, j := len(sys.LevelMatrix[i]), 0; j < l2; j++ {
			s := sys.LevelMatrix[i][j]
			output.Write([]byte(strconv.Itoa(i) + " stop " + s.Name + "(" + s.Version + "): "))
			s.Invoke(Operate_Stop, input, output)
		}
	}
}

func (sys *System) Status(input io.Reader, output io.Writer) {
	for l1, i := len(sys.LevelMatrix), 1; i < l1; i++ {
		output.Write([]byte(strconv.Itoa(i)))
		for l2, j := len(sys.LevelMatrix[i]), 0; j < l2; j++ {
			s := sys.LevelMatrix[i][j]
			output.Write([]byte(" status " + s.Name + "(" + s.Version + "): "))
			s.Status(input, output)
		}
		output.Write([]byte{'\n'})
	}
}
