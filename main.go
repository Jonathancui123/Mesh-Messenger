package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/campoy/whispering-gophers/util"
)

var (
	self     string
	dialAddr = flag.String("dial", "localhost:8000", "host:port to dial")
	dialTrue = flag.Bool("dialTrue", true, "determine whether or not to dial a connection")
)

type Message struct {
	Addr string
	Body string
}

func main() {
	flag.Parse()

	messageCh := make(chan Message)

	if *dialTrue {
		go dial(*dialAddr, messageCh) // start a go-routine for starting connection
	}

	l, err := util.Listen() // listen for TCP connections on given address
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on", l.Addr())
	self = l.Addr().String()
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

func read(messageCh chan Message) {
	lines := bufio.NewScanner(os.Stdin)

	for lines.Scan() {
		message := Message{
			Addr: self,
			Body: lines.Text()}

		messageCh <- message

	}
	if err := lines.Err(); err != nil {
		log.Fatal(err)
	}
}

func dial(addr string, messageCh chan Message) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(conn)
	for {
		err := enc.Encode(<-messageCh)

		if err != nil {
			log.Fatal(err)
		}
	}

}
