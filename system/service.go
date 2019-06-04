package system

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "[service] ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)

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
	Protocol      map[string]IProtocol
	Operate       map[string]*Operate
	Depth         int
	Dependency    []*Service
	Reference     []*Service
}

type Protocol struct {
	Name, URI string
}

type Operate struct {
	Name, Argument string
	Protocol       IProtocol
	Env            map[string]string
}

func (o *Operate) Invoke(input io.Reader, output io.Writer) (string, error) {
	err := o.Protocol.Open()
	if err != nil {
		return "", errors.WithMessage(err, "open (\""+o.Protocol.Name()+"\") error!")
	}
	defer o.Protocol.Close()
	return o.Protocol.Invoke(o, input, output), nil
}

func (s *Service) Invoke(operate string, input io.Reader, output io.Writer) (Result, error) {
	logger.Println(s.Name, "->", operate)
	r, e := s.Operate[operate].Invoke(input, output)
	if e == nil {
		return Result_Ok, nil
	} else {
		logger.Println(r)
		return Result_Failed, e
	}
}

func (s *Service) Start(input io.Reader, output io.Writer) (Result, error) {
	if s.Status(input, output) {
		return Result_Service_Ready_Running, nil
	} else {
		return s.Invoke(Operate_Start, input, output)
	}
}

func (s *Service) Stop(input io.Reader, output io.Writer) (Result, error) {
	return s.Invoke("stop", input, output)
}

func (s *Service) Status(input io.Reader, output io.Writer) bool {
	buf := bytes.Buffer{}
	_, err := s.Invoke(Operate_Status, input, &buf)
	if err == nil {
		str := buf.String()
		str = strings.TrimSpace(str)
		output.Write([]byte{'o', 'k'})
		return strings.HasSuffix(str, "ok")
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
