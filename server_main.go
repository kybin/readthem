package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func serverMain() {
	var (
		addr       string
		cert       string
		key        string
		configFile string
	)
	serverFlag := flag.NewFlagSet("server", flag.ExitOnError)
	serverFlag.StringVar(&addr, "addr", ":443", "https bind address")
	serverFlag.StringVar(&cert, "cert", "cert/cert.pem", "https certification file.\nreadthem ships with sample self-signed-certificate file,\nbut use your own certificate file for security.")
	serverFlag.StringVar(&key, "key", "cert/key.pem", "https key file.\nreadthem ships with sample self-signed-certificate file,\nbut use your own certificate file for security.")
	serverFlag.StringVar(&configFile, "config", "config.json", "json config file. where node name and address are filled in.")
	serverFlag.Parse(os.Args[2:])

	f, err := os.Open(configFile)
	if err != nil {
		die(err)
	}
	r := bufio.NewReader(f)
	dec := json.NewDecoder(r)
	nodes := make([]Node, 0)
	err = dec.Decode(&nodes)
	if err != nil {
		die(err)
	}

	http.HandleFunc("/node/", makeNodeHandler(nodes))
	http.HandleFunc("/", makeRootHandler(nodes))
	log.Fatal(http.ListenAndServeTLS(addr, cert, key, nil))
}

var serverTemplate = template.Must(template.New("server.html").ParseFiles("tmpl/server.html"))

type Node struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

func makeRootHandler(nodes []Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipt := struct {
			Nodes []Node
		}{
			Nodes: nodes,
		}
		err := serverTemplate.Execute(w, recipt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func makeNodeHandler(nodes []Node) http.HandlerFunc {
	address := make(map[string]string)
	for _, n := range nodes {
		address[n.Name] = n.Addr
	}
	return func(w http.ResponseWriter, r *http.Request) {
		node := strings.Split(r.URL.Path, "/")[2]
		if node == "" {
			http.Error(w, "node name not specified", http.StatusBadRequest)
			return
		}
		addr := address[node]
		if addr == "" {
			http.Error(w, "not found the node", http.StatusBadRequest)
			return
		}
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		bw := bufio.NewWriter(conn)
		_, err = bw.WriteString(node + "\n")
		if err != nil {
			log.Println(err)
			return
		}
		bw.Flush()
		_, err = bufio.NewReader(conn).WriteTo(w)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
