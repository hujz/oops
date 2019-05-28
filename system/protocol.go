package system

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"oops/util"
	"os"
	"strings"
)

import ossh "oops/ssh"

var sshLogger = log.New(os.Stdout, "[ssh] ", log.Llongfile|log.LstdFlags|log.Lmicroseconds)

func BuildProtocol(protocol Protocol) IProtocol {
	uri := protocol.URI
	switch {
	case strings.HasPrefix(uri, "ssh:"):
		ssh_ := &SSH{Url: uri[4:], Identity: uri[4:strings.Index(uri, "@")], Addr: uri[strings.Index(uri, "@")+1:]}
		return ssh_
	}
	return nil
}

type IProtocol interface {
	Open() error
	Close() error
	Name() string
	Invoke(string, io.Reader, io.Writer) string
}

type SSH struct {
	Url, Identity, Addr string
	session             *ssh.Session
}

func (s *SSH) Open() error {
	sshLogger.Println("open")
	desPw := []byte("123456781234567812345678")
	id, _ := hex.DecodeString(s.Identity)
	plain, err := util.DesDecrypt(id, desPw)
	if plain != nil {
		plainStr := string(plain)
		user, passwd := plainStr[:strings.Index(plainStr, ":")], plainStr[strings.Index(plainStr, ":")+1:]
		session, err := ossh.OpenSSHSession(user, passwd, s.Addr)
		if err == nil {
			s.session = session
			session.Stdin = os.Stdin
			session.Stdout = os.Stdout
			session.Stderr = os.Stderr
			return nil
		}
		return err
	}
	return err
}

func (s *SSH) Close() error {
	s.session.Close()
	return nil
}

func (s *SSH) Name() string {
	return "ssh:@" + s.Addr
}

func (s *SSH) Invoke(cmd string, reader io.Reader, writer io.Writer) string {
	if strings.HasPrefix(cmd, "!") {
		s.session.Stderr = writer
		s.session.Stdout = writer
		s.session.Stdin = reader
		s.session.Run(cmd[1:])
	} else {
		fmt.Println(cmd)
	}
	return ""
}
