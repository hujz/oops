package meta

type System struct {
	Version, Name string
	Service       map[string]*Service
	Instance      map[string]*Instance
	Host          map[string]*Host
}

type Host struct {
	Name, OS, Family, Virt string
	VIP, IP                []string
	Protocol               map[string]*Protocol
	Operate                map[string]*Operate
}

type Service struct {
	Spec, Version, Name string
	Protocol            map[string]*Protocol
	Operate             map[string]*Operate
	Dependency          []*Service
	Reference           []*Service
	Instance            []*Instance
}

type Instance struct {
	Name, Version string
	Host          Host
}

type Protocol struct {
	Name, URI string
}

type Operate struct {
	Name, Protocol, Argument string
}

func (sys *System) Build(system *XMLSystem, hostList []*XMLHost) {
	sys.Name = system.Name
	sys.Version = system.Version
	hostMap := make(map[string]Host)
	for _, h := range hostList {
		host := Host{Name: h.Name, OS: h.OS, Family: h.Family, Virt: h.Virt, VIP: h.VIP, IP: h.IP}
		protocolMap := make(map[string]Protocol)
		for _, p := range h.Protocol {
			protocolMap[p.Name] = Protocol{Name: p.Name, URI: p.URI}
		}
		operateMap := make(map[string]Operate)
		for _, o := range h.Operate {
			operateMap[o.Name] = Operate{Name: o.Name, Protocol: o.Protocol, Argument: o.Argument}
		}
		host.Protocol = protocolMap
		hostMap[h.Name] = host
	}

}

type IApplication interface {
	Start() (string, error)
	Stop() (string, error)
	Status() (string, error)
}
