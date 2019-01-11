package main

import (
	"bufio"
	"flag"
	"html/template"
	"log"
	"net"
	"os"
	"os/exec"
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

var nodeTemplate = template.Must(template.New("node.html").Funcs(template.FuncMap{
	"TextOutput": func(cmd ...string) string {
		if len(cmd) == 0 {
			return ""
		}
		c := exec.Command(cmd[0], cmd[1:]...)
		out, _ := c.CombinedOutput()
		return string(out)
	},
}).ParseFiles("tmpl/node.html"))

func sendReadMe(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(time.Second))
	defer conn.Close()
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	conn.SetWriteDeadline(time.Now().Add(time.Second))
	recipt := struct {
		Name string
	}{
		Name: name,
	}
	err = nodeTemplate.Execute(conn, recipt)
	if err != nil {
		log.Print(err)
		return
	}
}
