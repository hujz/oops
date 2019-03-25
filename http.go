package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/netutil"
)

// HttpServer start a HttpServer
func HttpServer(host string) {
	http.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir("dist"))))

	l, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	limitL := netutil.LimitListener(l, 4000)

	// err = http.Serve(limitL, nil)
	srv := &http.Server{Handler: nil, ReadTimeout: time.Second * 5, WriteTimeout: time.Second * 10}
	srv.SetKeepAlivesEnabled(false)
	err = srv.Serve(limitL)

	//	err = http.ListenAndServe(listenPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
