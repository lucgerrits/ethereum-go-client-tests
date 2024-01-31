package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// number of runs
	const N = 1
	// Set the amount to send (in wei)
	value := big.NewInt(100000000000000000) // 0.1 ETH

	// Start the profiling server
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	// Start CPU profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("Could not create CPU profile: ", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("Could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// Create a file to save the heap profile
	fHeap, err := os.Create("heap.prof")
	if err != nil {
		log.Fatal("Could not create heap profile: ", err)
	}
	defer fHeap.Close()

	// Start heap profiling
	if err := pprof.WriteHeapProfile(fHeap); err != nil {
		log.Fatal("Could not start heap profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// ---------------- Now the interesting part ---------------------------
	start := time.Now() // get the current time
	// Call sendEther function 100 times
	for i := 1; i <= N; i++ {
		fmt.Printf("#%d ", i)
		err := sendEther(value)
		if err != nil {
			log.Fatal(err)
		}
	}
	elapsed := time.Since(start) // calculate the elapsed time
	fmt.Printf("Elapsed time: %s\n", elapsed)
	// Delay to allow time for the heap profiling to run
	time.Sleep(30 * time.Second)
}

func sendEther(value *big.Int) error {

	// Connect to an Ethereum client
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	// Set the sender's private key
	privateKey, err := crypto.HexToECDSA("a47f9857331dee233af0100e9d18ae3f52cb2f21aa51ef530a60b8562a6b8868") // no 0x prefix
	// privateKey, err := crypto.HexToECDSA("99b3c12287537e38c90a9219d4cb074a89a16e9cdb20bf85728ebd97c343e342") // no 0x prefix
	if err != nil {
		log.Fatal(err)
	}

	// Set the recipient's address
	toAddress := common.HexToAddress("0xD92373D4BB00e93FaF91fD127492763d289e9487")

	// Set the gas price and gas limit
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(21000)

	// Set the sender's address
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Create the transaction
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// Sign the transaction
	// chainID, err := client.NetworkID(context.Background())
	chainID, err := new(big.Int).SetInt64(1337), nil
	// chainID, err := new(big.Int).SetInt64(42), nil
	if err != nil {
		return err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	// Print the transaction hash
	fmt.Printf("Transaction hash: 0x%x\n", signedTx.Hash())

	return nil
}
