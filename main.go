package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.Println("Hello world")
}
