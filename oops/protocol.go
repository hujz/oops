package oops

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

// ProtocolLinsten listen a port
func ProtocolLinsten(addr string) {

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

	hello := "You are welcome, input `help<enter>` print usage.\n"
	println(hello, conn)

	bufR := bufio.NewReader(conn)
	for {
		line, err := bufR.ReadString('\n')
		if err != nil && err == io.EOF {
			log.Println(err)
			break
		}

		switch protocolDispatch(line, conn) {
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

func protocolDispatch(line string, conn net.Conn) string {
	if strings.HasSuffix(line, "\r\n") {
		line = line[:len(line)-2]
	} else if strings.HasSuffix(line, "\n") {
		line = line[:len(line)-1]
	}

	//switch line {
	//case "ls", "ll":
	//	for _, n := range systemCache.Server {
	//		fmt.Fprintln(conn, n.Name)
	//	}
	//	printPrompt(conn)
	//case "help":
	//	println(help(), conn)
	//case "q", "exit", "quit":
	//	conn.Write([]byte("bye"))
	//	conn.Close()
	//	return "exit"
	//default:
	//	printPrompt(conn)
	//}

	return "ok"
}
