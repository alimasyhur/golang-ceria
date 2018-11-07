package main

import (
	"flag"
	"log"

	"github.com/pkg/errors"
)

func main() {
	connect := flag.String("connect", "", "IP address of process to join. If empty, go into listen mode.")
	flag.Parse()
	if *connect != "" {
		err := client(*connect)
		if err != nil {
			log.Println("Error:", errors.WithStack(err))
		}
		log.Println("Client done.")
		return
	}
	err := server()
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}

	log.Println("Server done.")
}

func init() {
	log.SetFlags(log.Lshortfile)
}
