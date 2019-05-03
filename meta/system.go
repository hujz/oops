package meta


type System struct {
	Version    string
	Name       string
	Service    []*Service
	TopService []*Service
	LowService []*Service
}

func (system *System) Build() {

}

func (system *System) Start() {
}

func (system *System) Stop() {
}

func (system *System) Status() {
}
