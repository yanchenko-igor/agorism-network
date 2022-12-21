package main

import (
	"agorism-network/blockchain"
	"agorism-network/p2p"
	"fmt"
	//"time"
)

func main() {
	bc := blockchain.Blockchain{}
	bc.NewGenesisBlock()

	go p2p.Listen(8080, &bc)
	peers := []string{"127.0.0.1:8081"}
	p2p.ConnectToPeers(peers, &bc)
	done := make(chan struct{})

	go func() {
		for {
			block := bc.FindBlock("Send 1 BTC to Alice")
			bc.AddBlock(block)
			p2p.SendBlock(block)
			//fmt.Println("Alice found a valid block: %+v", block)
			//fmt.Println("Blockchain: %+v", bc)
			for _, block := range bc.Blocks() {
				//fmt.Printf("%v\n", block)
				fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
				fmt.Printf("Data: %s\n", block.Data)
				fmt.Printf("Hash: %x\n", block.Hash)
				fmt.Printf("Nonce: %x\n", block.Nonce)
				fmt.Println()
			}
			fmt.Print("===================================\n")
		}
	}()

	// Wait for the application to be interrupted.
	<-done
}
