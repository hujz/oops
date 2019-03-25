package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// ProtocolListen listen a port
func ProtocolListen(addr string) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", addr)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
		return
	} else {
		log.Println("listening", addr)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
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

	if system := handShake(conn); system != nil {
		doCommand(conn, system)
		return
	}

}

func handShake(conn net.Conn) *System {
	hello := "You are welcome, input `help<enter>` print usage."
	fmt.Fprintln(conn, hello)

	if systemNames := GetSystemNames(); systemNames != nil {
		nameListStr := ""
		for _, n := range systemNames {
			nameListStr += ", \"" + n + "\""
		}

		nameListStr = nameListStr[2:]
	loop:
		fmt.Fprintln(conn, "=========================")
		fmt.Fprintln(conn, "Please select system: "+nameListStr)
		fmt.Fprintln(conn, "=========================")
		printPrompt(conn)

		bufR := bufio.NewReader(conn)
		line, err := bufR.ReadString('\n')
		if err == nil {
			line = trimCommandLine(line)
			if !arrayContain(systemNames, line) {
				goto loop
			}
			printPrompt(conn)
			return BuildSystem(line)
		}

		return nil
	} else {
		fmt.Fprintln(conn, "Not found any system data!")
		return nil
	}
}

func arrayContain(arr []string, el string) bool {
	for _, e := range arr {
		if e == el {
			return true
		}
	}
	return false
}

func doCommand(conn net.Conn, system *System) {
	bufR := bufio.NewReader(conn)
	for {
		line, err := bufR.ReadString('\n')
		if err != nil && err == io.EOF {
			log.Println(err)
			break
		}

		switch protocolDispatch(line, conn, system) {
		case "exit":
			return
		}
	}
}
func help() string {
	return `ls, ll
	列出所有服务
q, quit, exit
	退出命令行
select
	选择服务
`
}

var prompt []byte = []byte{'>', '>', ':', ' '}

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

func protocolDispatch(line string, conn net.Conn, system *System) string {
	line = trimCommandLine(line)

	switch line {
	case "ls", "ll":
		for _, n := range system.Server {
			fmt.Fprintln(conn, n.Name)
		}
		printPrompt(conn)
	case "help":
		println(help(), conn)
	case "q", "exit", "quit":
		conn.Write([]byte("bye"))
		conn.Close()
		return "exit"
	case "start":
		system.Start()
		printPrompt(conn)
	default:
		printPrompt(conn)
	}

	return "ok"
}
