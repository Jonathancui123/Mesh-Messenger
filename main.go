package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	listenAddr = flag.String("listen", "", "host:port to listen on")
	dialAddr   = flag.String("dial", "localhost:8000", "host:port to dial")
)

type Message struct {
	Body string
}

func main() {
	flag.Parse()

	go dial(*dialAddr) // start a go-routine for starting connection

	l, err := net.Listen("tcp", *listenAddr) // listen for TCP connections on given address

	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(c)
	}

}

func serve(c net.Conn) {
	defer c.Close()
	d := json.NewDecoder(c)
	for {
		var m Message
		err := d.Decode(&m)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%#v\n", m) // Print in go-syntax representation
	}
}

func dial(addr string) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	lines := bufio.NewScanner(os.Stdin)
	enc := json.NewEncoder(conn)
	for lines.Scan() {
		message := Message{Body: lines.Text()}
		err := enc.Encode(message)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := lines.Err(); err != nil {
		log.Fatal(err)
	}
}
