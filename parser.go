package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Transaction represents a basic structure for a transaction
type Transaction struct {
	From        common.Address `json:"from"`
	To          common.Address `json:"to"`
	Value       *big.Int       `json:"value"`
	Hash        common.Hash    `json:"hash"`
	BlockNumber *big.Int       `json:"blockNumber"`
}

// Parser implements the interface for parsing Ethereum transactions
type Parser struct {
	client       *ethclient.Client
	currentBlock *big.Int
	observers    map[string]bool
}

// NewParser creates a new Ethereum blockchain parser
func NewParser(url string) (*Parser, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Parser{
		client:       client,
		currentBlock: nil,
		observers:    make(map[string]bool),
	}, nil
}

// GetCurrentBlock returns the last parsed block number
func (p *Parser) GetCurrentBlock() int {
	if p.currentBlock == nil {
		return -1
	}
	return int(p.currentBlock.Int64())
}

// Subscribe adds an address to the list of observers
func (p *Parser) Subscribe(address string) bool {
	_, ok := p.observers[address]
	if !ok {
		p.observers[address] = true
		return true
	}
	return false
}

// GetTransactions retrieves a list of inbound or outbound transactions for an address
func (p *Parser) GetTransactions(address string) []Transaction {
	var transactions []Transaction

	// Get latest block number
	blockNumber, err := p.client.BlockNumber(context.Background())
	if err != nil {
		fmt.Println("Error getting latest block number:", err)
		return transactions
	}

	bigIntBlockNumber := big.NewInt(int64(blockNumber))

	// Update currentBlock if needed
	if p.currentBlock == nil || p.currentBlock.Cmp(bigIntBlockNumber) < 0 {
		p.currentBlock = bigIntBlockNumber
	}

	// Check if address is subscribed
	if !p.observers[address] {
		return transactions
	}

	// Iterate through blocks starting from the last parsed block
	for blockI := p.currentBlock.Int64(); blockI >= p.currentBlock.Int64(); blockI-- {

		block, err := p.client.BlockByNumber(context.Background(), big.NewInt(blockI))
		if err != nil {
			log.Println("Error getting block:", err)
			continue
		}

		for _, tx := range block.Transactions() {
			signer := types.LatestSignerForChainID(tx.ChainId())
			addr, err := signer.Sender(tx)
			if err != nil {
				log.Printf("Warning: skip transaction - unable to find sender for transaction %+v at block %v", tx, blockI)
				continue
			}
			if tx.To() == nil {
				log.Printf("Warning: skip transaction with empty receiver %+v at block %v", tx, blockI)
				continue
			}

			if tx.To().Hex() == address || addr.Hex() == address {
				transaction := Transaction{
					To:          *tx.To(),
					From:        addr,
					Value:       tx.Value(),
					Hash:        tx.Hash(),
					BlockNumber: big.NewInt(blockI),
				}
				transactions = append(transactions, transaction)
			}
		}
	}

	return transactions
}

func main() {
	// Define a flag for the Ethereum node URL
	nodeURL := flag.String("url", "https://cloudflare-eth.com", "Ethereum node URL (e.g. https://cloudflare-eth.com)")
	address := flag.String("addr", "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD", "an address to subscribe")
	flag.Parse()

	// Check if URL is provided
	if *nodeURL == "" {
		fmt.Println("Error: Please provide an Ethereum node URL using the -url flag.")
		return
	}

	// Create a new parser instance
	p, err := NewParser(*nodeURL)
	if err != nil {
		fmt.Println("Error creating parser:", err)
		return
	}

	// Example usage: Subscribe to an address and get transactions (modify as needed)
	subscribed := p.Subscribe(*address)
	if subscribed {
		fmt.Println("Subscribed to address:", address)
		transactions := p.GetTransactions(*address)
		if len(transactions) > 0 {
			fmt.Println("Transactions for", address, ":")
			for i, tx := range transactions {
				fmt.Printf("Transaction %v: %+v\n", i, tx)
			}
		} else {
			fmt.Println("No transactions found for", address)
		}
	} else {
		fmt.Println("Error subscribing to address:", address)
	}
}
