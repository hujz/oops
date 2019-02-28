package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// OpenSSHSession open a ssh
func OpenSSHSession(user, password, addr string) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

// ScpPushFile push file via scp
func ScpPushFile(client *ssh.Session, src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	fileInof, _ := file.Stat()
	size := fileInof.Size()
	fileName := fileInof.Name()

	go func() {
		w, _ := client.StdinPipe()
		fmt.Fprintln(w, "C0644", size, fileName)
		io.CopyN(w, file, size)
		fmt.Fprint(w, "\x00")
		w.Close()
	}()

	if err := client.Run(fmt.Sprintf("scp -qrt %s", "~/test.text")); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("succefully")

	return nil
}
