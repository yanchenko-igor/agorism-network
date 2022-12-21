package p2p

import (
	"agorism-network/blockchain"
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
)

// // Message is the structure of the messages that are exchanged between nodes.
// type Message struct {
// 	Nonce int
// 	Data  string
// }

// Peer represents a peer in the p2p network.
type Peer struct {
	conn  net.Conn
	enc   *gob.Encoder
	dec   *gob.Decoder
	Chain *blockchain.Blockchain
}

// FetchChain fetches the chain from the peer.
func (p *Peer) FetchChain() error {
	// Send a message to the peer to request the chain.
	err := p.SendMessage("fetch_chain")
	if err != nil {
		return err
	}
	// Receive the chain from the peer.
	var blocks []*blockchain.Block
	err = p.dec.Decode(&blocks)
	if err != nil {
		return err
	}
	// Set the peer's chain to the received blocks.
	p.Chain = &blockchain.Blockchain{Blocks: blocks}
	return nil
}

// SendChain sends the chain to the peer.
func (p *Peer) SendChain() error {
	// Send a message to the peer indicating that the chain is being sent.
	err := p.SendMessage("send_chain")
	if err != nil {
		return err
	}
	// Send the chain to the peer.
	err = p.enc.Encode(p.Chain.Blocks)
	if err != nil {
		return err
	}
	return nil
}

// SendMessage sends a message to the node.
func (p *Peer) SendMessage(m string) error {
// TODO implement messages
	return nil
}

// NewPeer creates a new Peer instance.
func NewPeer(conn net.Conn) *Peer {
	peer := &Peer{conn: conn}
	peer.enc = gob.NewEncoder(conn)
	peer.dec = gob.NewDecoder(conn)
	return peer
}

// Listen listens for incoming connections and handles them.
func Listen(port int, bc *blockchain.Blockchain) error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleConnection(conn, bc)
	}
}

//
// func handleConnection(conn net.Conn) {
//     defer conn.Close()
//     peer := NewPeer(conn)
//     for {
//         block, err := peer.ReceiveBlock()
//         if err != nil {
//             break
//         }
//         fmt.Println("Received block:", block)
// 	bc.AddBlock(block)
//         err = peer.SendBlock(block)
//         if err != nil {
//             break
//         }
//     }
// }
//

func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	defer conn.Close()
	peer := NewPeer(conn)
	for {
		block, err := peer.ReceiveBlock()
		if err != nil {
			break
		}
		if !bc.IsValidBlock(block) {
			fmt.Println("Received invalid block")
			continue
		}
		bc.AddBlock(block)
		fmt.Println("Received and added block:", block)
		err = peer.SendBlock(block)
		if err != nil {
			break
		}
		err = forwardBlockToPeers(block)
		if err != nil {
			break
		}
	}
}

// SendBlock sends a block to the connected peer.
func (p *Peer) SendBlock(block *blockchain.Block) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(block)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(p.conn, buf.String())
	return err
}

// ReceiveBlock receives a block from the connected peer.
func (p *Peer) ReceiveBlock() (*blockchain.Block, error) {
	message, err := bufio.NewReader(p.conn).ReadString('\n')
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.WriteString(message)
	dec := gob.NewDecoder(&buf)
	var block blockchain.Block
	err = dec.Decode(&block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// Peers is a slice of Peer instances.
var Peers []*Peer

// AddPeer adds a Peer to the Peers slice.
func AddPeer(peer *Peer) {
	Peers = append(Peers, peer)
}

func forwardBlockToPeers(block *blockchain.Block) error {
	for _, peer := range Peers {
		err := peer.SendBlock(block)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendNonce sends the nonce to all connected peers.
func SendBlock(b *blockchain.Block) {
	for _, peer := range Peers {
		peer.SendBlock(b)
	}
}

// ConnectToPeers connects to a list of predefined peers.
func ConnectToPeers(peers []string, bc *blockchain.Blockchain) error {
	for _, peer := range peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			return err
		}
		go handleConnection(conn, bc)
	}
	return nil
}
