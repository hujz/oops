package protocol

import (
	"encoding/hex"
	"golang.org/x/crypto/ssh"
	"io"
	"oops/util"
	"os"
	"strings"
)

import ossh "oops/ssh"

func BuildProtocol(uri string) IProtocol {
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
	Invoke(string) string
	SetInOut(io.Writer, io.Reader)
}

type SSH struct {
	Url, Identity, Addr string
	session             *ssh.Session
}

func (s *SSH) Open() error {
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

func (s *SSH) Invoke(cmd string) string {
	s.session.Run(cmd)
	return ""
}

func (s *SSH) SetInOut(writer io.Writer, reader io.Reader) {
	s.session.Stderr = writer
	s.session.Stdin = reader
	s.session.Stdout = writer
}
