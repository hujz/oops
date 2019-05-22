package system

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

func (system *System) Cache() {
	system.Build()
	SystemCache[system.Name] = system
}

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

func (system *System) Start() {

}

func (system *System) Stop() {
}

func (system *System) Status() {
}
