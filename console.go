package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"oops/meta"
	"oops/system"
	"strconv"
	"strings"
)

var logger = log.Logger{}

// ProtocolListen listen a port
func ProtocolListen(addr string) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		log.Fatalf("[console]: resolve tcpaddr(%s) fail: %s", addr, err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("[console]: listen %s fail: %s", addr, err)
		return
	} else {
		log.Println("[console]: successfully ", addr)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("[console]: accept request error: %s", err)
			continue
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer func() {
		p := recover()
		if p != nil {
			log.Println(p)
		}
	}()
	defer conn.Close()
loop:
	if system := handShake(conn); system != nil {
		switch doCommand(conn, system) {
		case "leave":
			goto loop
		}
		return
	}

}

func handShake(conn net.Conn) *system.System {
	hello := "You are welcome, input `help<enter>` print usage."
	fmt.Fprintln(conn, hello)

	if systemNames := meta.GetSystemNames(); systemNames != nil {
		// [41;37m 红底白字 \033[0m
		buf := bytes.Buffer{}
		buf.Write([]byte("Please select system: "))
		for _, n := range systemNames {
			buf.Write([]byte{27})
			buf.Write([]byte("[41;37m" + n + "\033[0m "))
		}

	loop:
		fmt.Fprintln(conn, "=========================")
		conn.Write(buf.Bytes())
		fmt.Fprintln(conn, "\n=========================")
		printPrompt(conn)

		bufR := bufio.NewReader(conn)
		line := readCommand(bufR)
		if line != "" {
			line = trimCommandLine(line)
			if !arrayContain(systemNames, line) {
				goto loop
			}
			sys := system.Get(line)
			if sys == nil {
				sys = meta.GetSystemMeta(line)
				if sys != nil {
					sys.Cache()
				} else {
					fmt.Fprintln(conn, "build system meta failed!!!")
					goto loop
				}
			}
			return sys
		}

		return nil
	} else {
		fmt.Fprintln(conn, "Not found any system data!")
		return nil
	}
}

func readCommand(bufR *bufio.Reader) string {
	line, err := bufR.ReadString('\n')
	if err == nil {
		return trimCommandLine(line)
	}
	return ""
}

func arrayContain(arr []string, el string) bool {
	for _, e := range arr {
		if e == el {
			return true
		}
	}
	return false
}
func help() string {
	return `ls, ll
	列出所有服务
ld
	列出服务，并按依赖排序
q, quit, exit
	退出命令行
select
	选择服务
`
}

var prompt = []byte{'>', '>', ':', ' '}

func println(data string, writer io.Writer) {
	writer.Write([]byte(data))
	writer.Write(prompt)
}

func printPrompt(writer io.Writer) {
	writer.Write(prompt)
}

func trimCommandLine(line string) string {
	if strings.HasSuffix(line, "\r\n") {
		line = line[:len(line)-2]
	} else if strings.HasSuffix(line, "\n") {
		line = line[:len(line)-1]
	}
	return line
}

func printPrompt2(writer io.Writer, system *system.System, service *system.Service) {
	writer.Write([]byte{27})
	writer.Write([]byte("[41;37m" + system.Name + "\033[0m"))
	writer.Write([]byte("/"))
	writer.Write([]byte{27})
	writer.Write([]byte("[42;37m" + system.Name + "\033[0m"))
	if service != nil {
		writer.Write([]byte(" -> "))
		writer.Write([]byte(service.Name))
	}
	writer.Write([]byte(">:"))
}

func doCommand(conn net.Conn, system *system.System) string {
	printPrompt2(conn, system, nil)
	bufR := bufio.NewReader(conn)
	for {
		line, err := bufR.ReadString('\n')
		if err != nil && err == io.EOF {
			log.Println(err)
			break
		}

		switch protocolDispatch(line, conn, system) {
		case "exit":
			return "leave"
		}
	}
	return "exit"
}

func protocolDispatch(line string, conn net.Conn, system *system.System) string {
	line = trimCommandLine(line)

	switch line {
	case "ls", "ll":
		for _, n := range system.Service {
			fmt.Fprintln(conn, n.Name)
		}
		printPrompt2(conn, system, nil)
	case "help":
		conn.Write([]byte(help()))
		conn.Write([]byte("\n"))
		printPrompt2(conn, system, nil)
	case "q", "exit", "quit":
		conn.Write([]byte("leave system(" + system.Name + ")!!!\n"))
		return "exit"
	case "start":
		system.Start()
		printPrompt2(conn, system, nil)
	case "ld":
		for i, n := range system.LevelMatrix {
			str := strconv.Itoa(i+1) + " - "
			if n != nil {
				for _, nn := range n {
					str = str + nn.Name + "  "
				}
			}
			conn.Write([]byte(str + "\n"))
		}
		printPrompt2(conn, system, nil)
	default:
		printPrompt2(conn, system, nil)
	}

	return "ok"
}
