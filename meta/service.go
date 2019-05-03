package meta

import (
	"github.com/pkg/errors"
	"log"
	"oops/protocol"
	"os"
)

var logger = log.New(os.Stdout, "[service]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)

type Status string
type Result string

const (
	Status_Starting              Status = "starting"
	Status_Running               Status = "running"
	Status_Downing               Status = "downing"
	Status_Down                  Status = "down"
	Result_Ok                    Result = "ok"
	Result_Failed                Result = "filed"
	Result_Service_Ready_Running Result = "running"
	Result_Unsupport             Result = "unsupport"
	Operate_Start                       = "start"
	Operate_Status                      = "status"
	Operate_Stop                        = "stop"
)

type Service struct {
	Version, Name string
	Env           map[string]string
	Protocol      map[string]*Protocol
	Operate       map[string]*Operate
	Dependency    []*Service
	Reference     []*Service
}

type Protocol struct {
	Name, URI string
}

type Operate struct {
	Name, Argument string
	Protocol       protocol.IProtocol
}

func (o *Operate) Invoke() (string, error) {
	err := o.Protocol.Open()
	return "", errors.WithMessage(err, "[protocol] open (\""+o.Protocol.Name()+"\") error!")
	defer o.Protocol.Close()
	return o.Protocol.Invoke(o.Argument), nil
}

func (s *Service) Invoke(operate string) (Result, error) {
	logger.Println(s.Name, ": --> ", operate)
	s.Operate[operate].Invoke()
	return Result_Ok, nil
}

func (s *Service) Start() (Result, error) {
	if s.Status() {
		return Result_Service_Ready_Running, nil
	} else {
		return s.Invoke(Operate_Start)
	}
}

func (s *Service) Stop() (Result, error) {
	return s.Invoke("stop")
}

func (s *Service) Status() bool {
	res, err := s.Invoke(Operate_Status)
	if err == nil {
		return res == "ok"
	} else {
		logger.Println("check status failed!", err)
		return false
	}
}

type IApplication interface {
	Start() (Result, error)
	Stop() (Result, error)
	Status() (Result, error)
}
