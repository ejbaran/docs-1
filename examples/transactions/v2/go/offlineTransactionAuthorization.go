package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// user declared account mnemonics for account1
const mnemonic1 = "boy kidney fall hamster ecology mercy inquiry vast deal normal vibrant labor couch economy embody glory possible color burger addict soap almost margin about negative" // TODO:"Your 25-word mnemonic goes here";"
const receiver = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"

// user declared algod connection parameters
const algodAddress = "http://localhost:49392"                                         //TODO:"http://localhost:4001"
const algodToken = "a31f09a18dbf7ad68c9e0ff22355774fb89c67ed2c4642d6c6822f9360cd7697" //TODO:"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

// Function that waits for a given txId to be confirmed by the network
func waitForConfirmation(txID string, client *algod.Client) {
	status, err := client.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting algod status: %s\n", err)
		return
	}
	lastRound := status.LastRound
	for {
		pt, _, err := client.PendingTransactionInformation(txID).Do(context.Background())
		if err != nil {
			fmt.Printf("error getting pending transaction: %s\n", err)
			return
		}
		if pt.ConfirmedRound > 0 {
			fmt.Printf("Transaction "+txID+" confirmed in round %d\n", pt.ConfirmedRound)
			break
		}
		fmt.Printf("...waiting for confirmation\n")
		lastRound++
		status, err = client.StatusAfterBlock(lastRound).Do(context.Background())
	}
}

// utility function to get address string
func getAddress(mn string) string {
	sk, err := mnemonic.ToPrivateKey(mn)
	if err != nil {
		fmt.Printf("error recovering account: %s\n", err)
		return ""
	}
	pk := sk.Public()
	var a types.Address
	cpk := pk.(ed25519.PublicKey)
	copy(a[:], cpk[:])
	fmt.Printf("...address: %s\n", a.String())
	address := a.String()
	return address
}

func saveUnsignedTransaction(algodClient *algod.Client) {
	// recover account for example
	account1 := getAddress(mnemonic1)

	// get network suggested parameters
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting suggested tx params: %s\n", err)
		return
	}
	txParams.FlatFee = true
	var minFee uint64 = 1000
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := uint64(txParams.FirstRoundValid)
	lastValidRound := uint64(txParams.LastRoundValid)

	// make transactions
	fmt.Println("Creating transactions...")
	// from account1 to receiver
	fromAddr := account1
	toAddr := receiver
	var amount uint64 = 1000000

	tx1, err := transaction.MakePaymentTxnWithFlatFee(fromAddr, toAddr, minFee, amount, firstValidRound, lastValidRound, nil, "", genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}
	fmt.Printf("...tx1: from %s to %s for %v microAlgos\n", fromAddr, toAddr, amount)

	// save unsigned Transaction to file
	err = ioutil.WriteFile("./unsigned.txn", msgpack.Encode(tx1), 0644)
	if err == nil {
		fmt.Printf("Saved unsigned transaction to file.\n")
		return
	}
	fmt.Printf("Failed in saving trx to file, error %s\n", err)
}

func readUnsignedTransaction(algodClient *algod.Client) []byte {
	// read unsigned transaction from file
	fmt.Println("Reading transaction from file...")
	dat, err := ioutil.ReadFile("./unsigned.txn")
	if err != nil {
		fmt.Printf("...Error reading transaction from file: %s\n", err)
	}
	fmt.Println("Decoding bytes...")
	var unsignedTxRaw types.SignedTxn
	msgpack.Decode(dat, &unsignedTxRaw)
	var unsignedTxn types.Transaction = unsignedTxRaw.Txn

	// get account from mnemonic
	fmt.Println("Loading account...")
	getAddress(mnemonic1)
	sk1, err := mnemonic.ToPrivateKey(mnemonic1)

	// sign the transaction
	fmt.Println("Signing transactions...")
	sTxID1, stx1, err := crypto.SignTransaction(sk1, unsignedTxn)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
	}
	fmt.Println("...account1 signed tx1: ", sTxID1)

	return stx1
}

func saveSignedTransaction(algodClient *algod.Client, signedTx []byte) {

	//Save the signed transaction to file
	err := ioutil.WriteFile("./signed.stxn", signedTx, 0644)
	if err == nil {
		fmt.Printf("Saved signed transaction to file\n")
		return
	}
	fmt.Printf("Failed in saving trx to file, error %s\n", err)
}

/*
func readSignedTransaction(){

    // setup connection
    algodClient := setupConnection()

    // read unsigned transaction from file
    dat, err := ioutil.ReadFile("./signed.stxn")
    if err != nil {
        fmt.Printf("Error reading signed transaction from file: %s\n", err)
        return
    }

    // send the transaction to the network
    sendResponse, err := algodClient.SendRawTransaction(dat)
    if err != nil {
        fmt.Printf("failed to send transaction: %s\n", err)
        return
    }

    fmt.Printf("Transaction ID: %s\n", sendResponse.TxID)
    waitForConfirmation(algodClient, sendResponse.TxID)
}
*/
func main() {
	// Initialize an algodClient
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %v\n", err)
		return
	}

	saveUnsignedTransaction(algodClient)
	var signedTx = readUnsignedTransaction(algodClient)
	saveSignedTransaction(algodClient, signedTx)
	//readSignedTransaction()
}
