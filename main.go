package main

import (
	"flag"
)

var (
	h             bool
	v             string
	http, console string
)

func init() {
	flag.BoolVar(&h, "h", false, "print usage")
	flag.StringVar(&v, "v", "1.0.0", "print version info")
	flag.StringVar(&http, "http", ":9527", "config port for http-server mode, enable manage/control via http api")
	flag.StringVar(&console, "console", ":9528", "start oops-protocol console, listen this port")
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}

	stop := make(chan int)
	go func() {
		HttpServer(http)
		stop <- 1
	}()
	go func() {
		ProtocolListen(console)
		stop <- 1
	}()

	for n := <-stop; n == 2; {
		n += <-stop
	}

}
