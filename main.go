package main

import (
	"flag"
	"github.com/matrix-go/bitcoin/server"
	"github.com/matrix-go/bitcoin/wallet_server"
	"log"
	"os"
	"os/signal"
)

func main() {
	port := flag.Int("port", 5000, "server port")
	flag.Parse()

	// wallet server
	wserv := wallet_server.NewServer(8001)
	go func() {
		if err := wserv.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// chain server
	serv := server.NewServer(*port)
	if err := serv.Start(); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
