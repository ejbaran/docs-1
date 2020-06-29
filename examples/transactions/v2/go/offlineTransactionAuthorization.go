package main

import (
	"context"
	"encoding/json"
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

// user declared mnemonic for myAccount
const myMnemonic = "boy kidney fall hamster ecology mercy inquiry vast deal normal vibrant labor couch economy embody glory possible color burger addict soap almost margin about negative" // TODO:"Your 25-word mnemonic goes here";"
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

func createTransaction(algodClient *algod.Client, myAccount string) types.Transaction {
	// get network suggested parameters
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting suggested tx params: %s\n", err)
	}
	txParams.FlatFee = true
	var minFee uint64 = 1000
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := uint64(txParams.FirstRoundValid)
	lastValidRound := uint64(txParams.LastRoundValid)

	// make transaction
	fmt.Println("Creating transaction...")
	// from account to receiver
	fromAddr := myAccount
	toAddr := receiver
	var amount uint64 = 1000000

	tx1, err := transaction.MakePaymentTxnWithFlatFee(fromAddr, toAddr, minFee, amount, firstValidRound, lastValidRound, nil, "", genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
	}
	fmt.Printf("...tx1: from %s to %s for %v microAlgos\n", fromAddr, toAddr, amount)

	return tx1
}
func saveUnsignedTransactionToFile(algodClient *algod.Client, txn types.Transaction) {
	// assign Transaction data to SignedTxn struct
	unsignedTx := types.SignedTxn{
		Txn: txn,
	}

	// save unsigned Transaction to file
	err := ioutil.WriteFile("./unsigned.txn", msgpack.Encode(unsignedTx), 0644)
	if err == nil {
		fmt.Printf("Saving unsigned transaction to file...\n")
		return
	}
	fmt.Printf("...Failed to save transaction to file, error %s\n", err)
}

func readUnsigedTransactionFromFile() types.Transaction {
	// read unsigned transaction from file
	fmt.Println("Reading transaction from file...")
	dat, err := ioutil.ReadFile("./unsigned.txn")
	if err != nil {
		fmt.Printf("...Error reading transaction from file: %s\n", err)
	}
	fmt.Println("Decoding file bytes...")
	var unsignedTxRaw types.SignedTxn
	msgpack.Decode(dat, &unsignedTxRaw)
	unsignedTxn := unsignedTxRaw.Txn

	// TODO: display unsigned transaction

	return unsignedTxn
}

func signTransaction(unsignedTxn types.Transaction, myMnemonic string) []byte {
	// Load private key from mnemonic
	sk1, err := mnemonic.ToPrivateKey(myMnemonic)
	if err != nil {
		fmt.Printf("...Failed to convert mnemonic to private key: %v\n", err)
	}

	// sign the transaction
	fmt.Println("Signing transactions...")
	signedTxnId, signedBytes, err := crypto.SignTransaction(sk1, unsignedTxn)
	if err != nil {
		fmt.Printf("...Failed to sign transaction: %s\n", err)
	}
	fmt.Println("...Signed transaction: ", signedTxnId)

	return signedBytes
}

func saveSignedTransactionToFile(signedBytes []byte) {
	// save the signed transaction to file
	fmt.Println("Saving signed transction to file...")
	err := ioutil.WriteFile("./signed.stxn", signedBytes, 0644)
	if err == nil {
		fmt.Printf("...Saved signed transaction to file\n")
		return
	}
	fmt.Printf("...Failed to save transaction to file, error %s\n", err)
}

func readSignedTransactionFromFile() []byte {
	// read unsigned transaction from file
	fmt.Println("Reading transaction from file...")
	signedBytes, err := ioutil.ReadFile("./signed.stxn")
	if err != nil {
		fmt.Printf("...Error reading signed transaction from file: %s\n", err)
	}
	fmt.Println("Decoding bytes...")

	// display signed transaction
	var signedTxRaw types.SignedTxn
	var signedTxn types.Transaction
	msgpack.Decode(signedBytes, &signedTxRaw)
	signedTxn = signedTxRaw.Txn

	txnJSON, err := json.MarshalIndent(signedTxn, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall txn data: %s\n", err)
	}
	fmt.Printf("Transaction information: %s\n", txnJSON)

	return signedBytes
}

func sendSignedTransaction(algodClient *algod.Client, signedBytes []byte) {
	// send the transaction to the network
	fmt.Println("Sending signed transction to network...")
	txId, err := algodClient.SendRawTransaction(signedBytes).Do(context.Background())
	if err != nil {
		fmt.Printf("...Failed to send transaction: %s\n", err)
		return
	}

	// wait for response
	waitForConfirmation(txId, algodClient)
}

func main() {
	// Initialize an algodClient
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %v\n", err)
		return
	}

	// Load account from Mymnemonic
	fmt.Println("Loading signing account...")
	myAccount := getAddress(myMnemonic)

	// Create transaction from myAccount
	unsignedTxn := createTransaction(algodClient, myAccount)

	// Save unsigned transaction to file
	saveUnsignedTransactionToFile(algodClient, unsignedTxn)

	// Read the unsigned transaction from the file
	unsignedTxn = readUnsigedTransactionFromFile()

	// Sign the transaction using the mnemonic
	signedBytes := signTransaction(unsignedTxn, myMnemonic)

	// Save the signed transaction to file
	saveSignedTransactionToFile(signedBytes)

	// Read the signed transaction from file
	signedBytes = readSignedTransactionFromFile()

	// Send the transaction to the network
	sendSignedTransaction(algodClient, signedBytes)
}
