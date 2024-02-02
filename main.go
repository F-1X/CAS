package main

import (
	"fmt"
	"log"

	"github.com/F-1X/CAS/p2p"
)

func OnPeer(p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside TCP transport")
	return nil
}

func main() {

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":4000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		OnPeer: OnPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <- tr.Consume()
			fmt.Printf("msg: %+v\n",msg)
		}
	}()


	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select{}
}