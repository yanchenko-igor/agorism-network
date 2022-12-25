// main.go
package main

import (
	"agorism-network/blockchain"
	"fmt"
	//"log"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(2)

	newBlockChannel := make(chan *blockchain.Block)

	bc := blockchain.NewBlockchain()
	// Launch two goroutines
	go func() {
		defer wg.Done()
		blockchain.HandleNewBlocks(bc, newBlockChannel)
	}()
	// Start mining for new blocks in a goroutine
	go func() {
		defer wg.Done()
		blockchain.MineBlock(bc, newBlockChannel)
	}()

	// Wait for the goroutines to finish
	wg.Wait()
	fmt.Println("All goroutines completed")
}
