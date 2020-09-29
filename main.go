package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Message struct {
	Body string
}

func main() {
	lines := bufio.NewScanner(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
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
