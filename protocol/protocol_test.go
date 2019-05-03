package protocol

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"testing"
)

import pssh "oops/ssh"

func TestSSH_Open(t *testing.T) {
	session, _ := pssh.OpenSSHSession("hujz", "123", "127.0.0.1:22")
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Println(err)
	}

	session.Run("sudo ls")
}
