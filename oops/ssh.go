/**
 * https://github.com/andesli/gossh/blob/master/scp/scp.go
 */

package oops

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

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
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
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

	go func() {
		writer, err := client.StdinPipe()
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		defer writer.Close()
		upstream(writer, src)
	}()

	if err := client.Run(fmt.Sprintf("scp -rt %s", dst)); err != nil {
		log.Fatalln(err.Error())
		return err
	}

	return nil
}

// ScpPullFile pull file via scp
func ScpPullFile(client *ssh.Session, src, dst string) error {
	go func() {
		relDir, whiteBit := dst, make([]byte, 1)
		writer, err := client.StdinPipe()
		if err != nil {
			log.Fatalln("[pullfile] get stdin pipe failed, " + err.Error())
			return
		}
		reader, err := client.StdoutPipe()
		if err != nil {
			log.Fatalln("[pullfile] get stdout pipe failed, " + err.Error())
			return
		}
		buf := bufio.NewReader(reader)
		for true {
			fmt.Fprint(writer, "\x00")
			headLine, err := buf.ReadString('\n')
			if err != nil {
				fmt.Fprint(writer, "\x02")
				continue
			}
			log.Println("[scp] headline:" + headLine[:len(headLine)-1])
			if strings.HasPrefix(headLine, "D") {
				_, fileMode, _, name := parseHead(headLine)
				dstDir := path.Join(relDir, name)
				err := os.MkdirAll(dstDir, fileMode)
				if err != nil {
					log.Fatalln("[scp] mkdir failed, `" + dstDir + "`! " + err.Error())
					return
				}
				log.Println("[scp] mkdir `" + dstDir + "`.")
				relDir = dstDir
			} else if strings.HasPrefix(headLine, "C") {
				_, fileMode, size, name := parseHead(headLine)
				dstFile := path.Join(relDir, name)
				_, err := os.Open(relDir)

				if err != nil {
					fmt.Println(err)
					os.MkdirAll(relDir, fileMode)
				}
				file, err := os.OpenFile(dstFile, os.O_CREATE|os.O_RDWR, fileMode)
				if err != nil {
					log.Fatalln("[scp] create file failed, `" + dstFile + "`, " + err.Error())
					return
				}
				log.Println("[scp] save file `" + dstFile + "`.")
				defer file.Close()
				fmt.Fprint(writer, "\x00")
				n, err := io.CopyN(file, buf, size)
				if n != size {
					return
				}
				buf.Read(whiteBit)
				if err != nil {
					log.Fatalln(err.Error())
					return
				}
			} else if strings.HasPrefix(headLine, "E") {
				relDir = path.Join(relDir, "../")
			} else {
				break
			}
		}
	}()

	client.Stderr = os.Stdout

	if err := client.Run(fmt.Sprintf("scp -rf %s", src)); err != nil {
		log.Fatalln("[pullfile] run scp failed, " + err.Error())
		return err
	}
	return nil
}

func parseHead(hl string) (string, os.FileMode, int64, string) {
	lineItem := strings.Split(hl, " ")
	mod, name := lineItem[0][2:], lineItem[2][:len(lineItem[2])-1]
	fileMode := os.FileMode(Oct2Dec(mod))
	size, _ := strconv.Atoi(lineItem[1])
	return string(hl[1:]), fileMode, int64(size), name
}

// ctrl+c from https://github.com/andesli/gossh/blob/master/scp/scp.go
const (
	ScpPushBeginFile      = "C"
	ScpPushBeginFolder    = "D"
	ScpPushBeginEndFolder = "0"
	ScpPushEndFolder      = "E"
	ScpPushEnd            = "\x00"
)

// ctrl+c from https://github.com/andesli/gossh/blob/master/scp/scp.go
func upstream(writer io.WriteCloser, src string) {
	file, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileInof, _ := file.Stat()
	size, fileName := fileInof.Size(), fileInof.Name()
	mode, isDir := getFileMode(*file)
	if isDir {
		fmt.Fprintln(writer, ScpPushBeginFolder+mode, ScpPushBeginEndFolder, fileName)
		subDirList, err := file.Readdir(0)
		if err == nil {
			for _, info := range subDirList {
				if err == nil {
					upstream(writer, path.Join(src, info.Name()))
				} else {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println(err)
		}
		fmt.Fprintln(writer, ScpPushEndFolder)
	} else {
		fmt.Fprintln(writer, ScpPushBeginFile+mode, size, fileName)
		io.CopyN(writer, file, size)
		fmt.Fprint(writer, ScpPushEnd)
	}
}

func getFileMode(file os.File) (string, bool) {
	fileStat, err := file.Stat()
	if err == nil {
		fileMode := fileStat.Mode()
		perm := fileMode.Perm()
		return "0" + Dec2Oct(int(perm)), fileStat.IsDir()
	}
	return "", false
}

// Oct2Dec Oct2Dec
func Oct2Dec(oct string) int {
	var sum, temp, temp2 int
	for l, n := len(oct), 1; n <= l; n++ {
		temp = 1
		for m := 1; m < n; m++ {
			temp *= 8
		}
		temp2, _ = strconv.Atoi(string(oct[l-n]))
		sum += temp * temp2
	}
	return sum
}

// Dec2Oct Dec2Oct
func Dec2Oct(dec int) string {
	var mod, remain int
	var ret string
	for true {
		remain = dec / 8
		mod = dec % 8
		ret = strconv.Itoa(mod) + ret
		if remain == 0 {
			break
		}
		dec = remain
	}
	return ret
}
