package main

import (
	"log"

	"innoseti_task/claim"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	asset, err := contractapi.NewChaincode(&claim.SmartContract{})
	if err != nil {
		log.Panicf("Error creating claim smartcontract: %v", err)
	}

	if err := asset.Start(); err != nil {
		log.Panicf("Error starting claim smartcontract: %v", err)
	}
}
