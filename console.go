package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"oops/meta"
	"oops/system"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var logger = log.New(os.Stdout, "[console] ", log.Lshortfile|log.LstdFlags|log.Lmicroseconds)

// ProtocolListen listen a port
func ProtocolListen(addr string) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		logger.Fatalf("resolve tcpaddr(%s) fail: %s", addr, err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Fatalf("listen %s fail: %s", addr, err)
		return
	} else {
		logger.Println("successfully ", addr)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Fatalf("accept request error: %s", err)
			continue
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer func() {
		p := recover()
		if p != nil {
			logger.Println(p)
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			logger.Println(string(buf[:n]))
		}
	}()
	defer conn.Close()
	handler := NewHandler(conn)
	handler.Done()
}

type Handler struct {
	sys     *system.System
	svr     *system.Service
	conn    net.Conn
	bufR    *bufio.Reader
	prompt  []byte
	systems []string
}

func NewHandler(conn net.Conn) *Handler {
	h := &Handler{conn: conn, bufR: bufio.NewReader(conn)}
	h.resetPrompt()
	return h
}

func (h *Handler) Prompt() {
	h.conn.Write(h.prompt)
}
func (h *Handler) PrintlnPrompt() {
	h.conn.Write([]byte("\r\n"))
	h.conn.Write(h.prompt)
}

func (h *Handler) Println(s []byte) {
	h.conn.Write(s)
	h.conn.Write([]byte("\n"))
}

func (h *Handler) PrintlnEnd(s []byte) {
	h.conn.Write(s)
	h.conn.Write([]byte("\n"))
	h.conn.Write(h.prompt)
}

func (h *Handler) Done() {
	h.HandShake()
	for {
	sls:
		if h.SelectSystem() == 0 {
			h.Close()
			return
		}
	sli:
		switch h.SystemPhase() {
		case 2:
			h.Prompt()
		case 1:
			h.Prompt()
			continue
		case 0:
			h.Close()
			return
		default:
			continue
		}
		switch h.InstancePhase() {
		case 2:
			h.Prompt()
			goto sls
		case 1:
			h.Prompt()
			goto sli
		case 0:
			h.Close()
			return
		}
	}
}

func (h *Handler) Close() {
	h.Println([]byte("bye"))
	logger.Println(h.conn, "closed")
	h.conn.Close()
}

func (h *Handler) ReadCommand() []string {
	line, err := h.bufR.ReadString('\n')
	if err == nil {
		if strings.HasSuffix(line, "\r\n") {
			line = line[:len(line)-2]
		} else if strings.HasSuffix(line, "\n") {
			line = line[:len(line)-1]
		}
		line = strings.TrimSpace(line)

		n := strings.Index(line, " ")
		var cmd []string
		if n == -1 {
			cmd = []string{line}
		} else {
			svr := strings.TrimSpace(line[n:])
			m := strings.Index(svr, " ")
			if m == -1 {
				cmd = []string{line[:n], svr}
			} else {
				cmd = []string{line[:n], svr[:m], strings.TrimSpace(svr[m:])}
			}
		}
		return cmd
	} else {
		return nil
	}
}

func (h *Handler) HandShake() {
	h.Println([]byte("Oops... you are welcome!!!"))
	if h.systems = meta.GetSystemNames(); h.systems != nil {
		// [41;37m 红底白字 \033[0m
		buf := bytes.Buffer{}
		buf.Write([]byte("Please select system: "))
		for i, n := range h.systems {
			buf.Write([]byte{27})
			buf.Write([]byte("[41;37m" + strconv.Itoa(i+1) + ":" + n + "\033[0m "))
		}
		h.PrintlnEnd(buf.Bytes())
	}
}

// SelectSystem selected system return 1, exit console return 0
func (h *Handler) SelectSystem() int {
	for {
		cmd := h.ReadCommand()
		if cmd == nil {
			return 0
		}
		name := ""
		i, err := strconv.Atoi(cmd[0])
		if err == nil && len(h.systems) >= i && i >= 1 {
			name = h.systems[i-1]
		} else {
			name = cmd[0]
		}
		for _, s := range h.systems {
			if name == s {
				sys := system.Get(name)
				if sys == nil {
					sys = meta.GetSystem(name)
					sys.Cache()
				}
				if sys == nil {
					h.PrintlnEnd([]byte("build system(\"" + s + "\") failed!"))
				} else {
					h.sys = sys
					h.svr = nil
					h.resetPrompt()
					h.Prompt()
					return 1
				}
			}
		}
		h.PrintlnEnd([]byte("select(\"" + cmd[0] + "\") not found in system list!"))
	}
}

// SystemPhase select instance return 2, break return 1, quit console return 0
func (h *Handler) SystemPhase() int {
	for {
		cmd := h.ReadCommand()
		if cmd == nil {
			return 0
		}
		switch cmd[0] {
		case "start":
			h.sys.Start(h.conn, h.conn)
			h.PrintlnPrompt()
		case "stop":
			h.sys.Stop(h.conn, h.conn)
			h.PrintlnPrompt()
		case "status":
			h.sys.Status(h.conn, h.conn)
			h.PrintlnPrompt()
		case "l", "ls", "ll":
			buf := bytes.Buffer{}
			for i, s := range h.sys.Service {
				buf.Reset()
				buf.Write([]byte(strconv.Itoa(i + 1)))
				buf.Write([]byte{' '})
				buf.Write([]byte(s.Name))
				buf.Write([]byte{'('})
				buf.Write([]byte(s.Version))
				buf.Write([]byte{')'})
				h.Println(buf.Bytes())
			}
			h.Prompt()
		case "d", "ld":
			buf := bytes.Buffer{}
			for i := range h.sys.LevelMatrix {
				buf.Reset()
				buf.Write([]byte(strconv.Itoa(i + 1)))
				buf.Write([]byte{' '})
				it := h.sys.LevelMatrix[i]
				for _, ss := range it {
					buf.Write([]byte(ss.Name))
					buf.Write([]byte{'('})
					buf.Write([]byte(ss.Version))
					buf.Write([]byte{')', ' '})
				}
				h.Println(buf.Bytes())
			}
			h.Prompt()
		case "select", "s":
			if len(cmd) <= 1 {
				h.PrintlnEnd([]byte("select command like: select <instance-name>[ <instance-version>]"))
				continue
			}
			if len(cmd) == 2 {
				i, err := strconv.Atoi(cmd[1])
				if err == nil && i >= 1 && len(h.sys.Service) >= i {
					h.svr = h.sys.Service[i-1]
					h.resetPrompt()
					return 2
				}
			}
			for i := range h.sys.Service {
				if len(cmd) == 2 && h.sys.Service[i].Name == cmd[1] {
					h.svr = h.sys.Service[i]
					h.resetPrompt()
					return 2
				} else if len(cmd) == 3 && h.sys.Service[i].Name == cmd[1] && h.sys.Service[i].Version == cmd[2] {
					h.svr = h.sys.Service[i]
					h.resetPrompt()
					return 2
				}
			}

			errMsg := "not found instance(\"" + cmd[1]
			if len(cmd) == 3 {
				errMsg = errMsg + "(" + cmd[2] + ")"
			}
			errMsg = errMsg + "\")"
			h.PrintlnEnd([]byte(errMsg))
			continue
		case "b", "leave", "byte":
			h.sys, h.svr = nil, nil
			h.resetPrompt()
			return 1
		case "q", "quit", "exit":
			return 0
		case "h", "help":
			help := `l, ls, ll
	列出所有服务
d, ld
	列出服务，并按依赖排序
b, leave, byte
	返回系统选择
q, quit, exit
	退出命令行
s, select
	选择服务
h, help
	帮助信息`
			h.PrintlnEnd([]byte(help))
		default:
			h.PrintlnEnd([]byte("unknown command"))
		}
	}
}

// InstancePhase go back select system return 2, go back select instance return 1, exit console return 0
func (h *Handler) InstancePhase() int {
	for {
		cmd := h.ReadCommand()
		if cmd == nil {
			return 0
		}
		// service operate
		if op := h.svr.Operate[cmd[0]]; op != nil {
			h.svr.Invoke(cmd[0], h.conn, h.conn)
			h.PrintlnPrompt()
			continue
		}

		switch cmd[0] {
		case "!", "o", "op", "operate":
			buf := bytes.Buffer{}
			for i := range h.svr.Operate {
				buf.Write([]byte(h.svr.Operate[i].Name))
				buf.Write([]byte("(via (" + h.svr.Operate[i].Protocol.Name() + ")\n"))
			}
			if buf.Len() >= 1 {
				buf.Truncate(buf.Len() - 1)
			}
			h.PrintlnEnd(buf.Bytes())
		case "#", "p", "prot", "protocol":
			buf := bytes.Buffer{}
			for i := range h.svr.Protocol {
				buf.Write([]byte(i))
				buf.Write([]byte{' '})
				buf.Write([]byte(h.svr.Protocol[i].Name()))
				buf.Write([]byte{'\n'})
			}
			if buf.Len() >= 1 {
				buf.Truncate(buf.Len() - 1)
			}
			h.PrintlnEnd(buf.Bytes())
		case "@", "s", "start":
			h.svr.Invoke(system.Operate_Start, h.conn, h.conn)
			h.Prompt()
		case "$", "k", "stop":
			h.svr.Invoke(system.Operate_Stop, h.conn, h.conn)
			h.Prompt()
		case "%", "status":
			h.svr.Invoke(system.Operate_Status, h.conn, h.conn)
			h.Prompt()
		case "*", "info":
			h.PrintlnEnd([]byte(h.svr.Name))
		case "=", "q", "quit", "leave", "bye":
			h.svr = nil
			h.resetPrompt()
			return 1
		case "=s", "system":
			h.svr, h.sys = nil, nil
			h.resetPrompt()
			return 2
		case ".", "exit":
			h.svr, h.sys = nil, nil
			h.resetPrompt()
			return 0
		case "h", "help", "?":
			str := `!, op, operate
	列出该服务支持的操作
#, prot, protocol
	列出该服务支持使用的协议
@, start
	启动该服务
$, stop
	停止该服务
%, status
	查看服务状态
*, info
	查看服务信息
=, q, quit, leave, bye
	退出该服务控制台
=s, system
	跳回选择系统控制台
., exit
	退出控制台
?, h, help
	打印帮助信息`
			h.PrintlnEnd([]byte(str))
		default:
			if len(cmd[0]) != 0 {
				h.PrintlnEnd([]byte("unknown command"))
			}
		}
	}
}

func (h *Handler) resetPrompt() {
	if h.sys == nil {
		h.prompt = []byte{'o', 'o', 'p', 's', '>', '$', ' '}
		return
	}
	buf := bytes.Buffer{}
	buf.Write([]byte{27})
	buf.Write([]byte("[41;37m" + h.sys.Name + "\033[0m"))
	if h.svr != nil {
		buf.Write([]byte{'#', 27})
		buf.Write([]byte("[42;37m" + h.svr.Name + "\033[0m"))
	}
	buf.Write([]byte{'>', '$', ' '})
	h.prompt = buf.Bytes()
}
