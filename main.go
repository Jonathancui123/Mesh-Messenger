package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/campoy/whispering-gophers/util"
)

var (
	self     string
	nickname = flag.String("name", "unknown", "name shown in messages")
	peerAddr = flag.String("dial", "", "peer host:port to dial")
)

// Message is a structure for sending messages in this mesh network
type Message struct {
	Nickname string
	ID       string
	Addr     string
	Body     string
}

func main() {
	flag.Parse()

	l, err := util.Listen() // listen for TCP connections on given address
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on", l.Addr())
	self = l.Addr().String()

	if len(*peerAddr) > 0 {
		go dial(*peerAddr) // start a go-routine for connection
	}

	go read() // start a go-routine for reading from stdIn

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(c)
	}

}

var peers = &Peers{m: make(map[string]chan<- Message)}

// Peers is a structure for concurrent-safe Peers registry
type Peers struct {
	m  map[string]chan<- Message
	mu sync.RWMutex
}

// Add creates and returns a new channel for a given peer address.
// If the address already exists in peer registry, return nil
func (p *Peers) Add(addr string) <-chan Message {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.m[addr]; ok {
		return nil
	}
	ch := make(chan Message)
	p.m[addr] = ch
	return ch
}

// Remove deletes the specified peer from the registry
func (p *Peers) Remove(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.m[addr]; ok {
		delete(p.m, addr)
	}
}

// List returns a slice of all peers in the registry
func (p *Peers) List() []chan<- Message {
	p.mu.RLock()
	defer p.mu.RUnlock()

	l := make([]chan<- Message, 0, len(p.m))
	for _, ch := range p.m {
		l = append(l, ch)
	}
	return l
}

func broadcast(m Message) {
	for _, ch := range peers.List() {
		select {
		case ch <- m:
		default: // In a mesh network, its okay to drop messages.
		}
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
		if !Seen(m.ID) {
			go dial(m.Addr)
			broadcast(m)
			fmt.Printf("%v @ %v: %v\n ", m.Nickname, m.Addr, m.Body) // Print in go-syntax representation
		}
	}
}

func read() {
	lines := bufio.NewScanner(os.Stdin)
	for lines.Scan() {
		message := Message{
			Nickname: *nickname,
			ID:       util.RandomID(),
			Addr:     self,
			Body:     lines.Text(),
		}
		Seen(message.ID)
		broadcast(message)
	}
	if err := lines.Err(); err != nil {
		log.Fatal(err)
	}
}

func dial(addr string) {
	if addr == self {
		return // Don't attempt to dial yourself
	}

	ch := peers.Add(addr)
	if ch == nil {
		return // Peer already connected.
	}
	defer peers.Remove(addr)

	conn, err := net.Dial("tcp", addr)
	fmt.Printf("* Connected to %v *\n", addr)
	if err != nil {
		log.Println(addr, err)
		return
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)
	for m := range ch {
		err := enc.Encode(m)

		if err != nil {
			log.Println(addr, err)
			return
		}
	}
}

var seenIDs = struct {
	m map[string]bool
	sync.Mutex
}{m: make(map[string]bool)}

// Seen returns true if the specified id has been seen before
// If not, it returns false and marks the id as seen
func Seen(id string) bool {
	seenIDs.Lock()
	defer seenIDs.Unlock()
	ok := seenIDs.m[id]
	seenIDs.m[id] = true
	return ok
}
