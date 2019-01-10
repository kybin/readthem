package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"
)

func nodeMain() {
	var (
		addr string
	)
	nodeFlag := flag.NewFlagSet("node", flag.ExitOnError)
	nodeFlag.StringVar(&addr, "addr", ":52939", "tcp bind address")
	nodeFlag.Parse(os.Args[2:])

	l, err := net.Listen("tcp", addr)
	if err != nil {
		die(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go sendReadMe(conn)
	}
}

func sendReadMe(conn net.Conn) {
	conn.SetWriteDeadline(time.Now().Add(time.Second))
	defer conn.Close()
	err := templates.ExecuteTemplate(conn, "node.html", nil)
	if err != nil {
		log.Print(err)
	}
}
