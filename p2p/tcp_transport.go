package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransportOpts struct {
	ListenAddr string
	HandshakeFunc HandshakeFunc
	Decoder Decoder
	OnPeer func(Peer) error
}

type TCPPeer struct {
	conn net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn: conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	TCPTransportOpts
	listener   net.Listener
	rpcch chan RPC

	mu sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch: make(chan RPC),
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	fmt.Println("start listening")
	ln, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	t.listener = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		fmt.Printf("new incoming connection: %+v\n", conn)
		go t.handleConn(conn)

	}
}


func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	
	var err error
	defer func() {
		fmt.Printf("dropping peer connection: %+v\n",err)
		conn.Close()
	}()

	if err = t.HandshakeFunc(peer); err != nil {
		//conn.Close()
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			fmt.Printf("here: %+v\n",err)
			return
		}
	}
	
	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			if err == net.ErrClosed {
				return
			}
			fmt.Printf("tcp read error: %+v\n",err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <-rpc
		fmt.Println(rpc)
	}

	

}