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

// utility function to get Account object
func getAccount(mn string) (string, ed25519.PrivateKey) {
	sk, err := mnemonic.ToPrivateKey(mn)
	if err != nil {
		fmt.Printf("...error recovering account: %s\n", err)
	}
	pk := sk.Public()
	var a types.Address
	cpk := pk.(ed25519.PublicKey)
	copy(a[:], cpk[:])
	fmt.Printf("...found address: %s\n", a.String())
	address := a.String()
	return address, sk
}

func createTransaction(algodClient *algod.Client, myAccount string) types.Transaction {
	// get network suggested parameters
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("...error getting suggested tx params: %s\n", err)
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

	txnObj, err := transaction.MakePaymentTxnWithFlatFee(fromAddr, toAddr, minFee, amount, firstValidRound, lastValidRound, nil, "", genID, genHash)
	if err != nil {
		fmt.Printf("...error creating transaction: %s\n", err)
	}
	fmt.Printf("...tx1: from %s to %s for %v microAlgos\n", fromAddr, toAddr, amount)

	return txnObj
}
func saveUnsignedTransactionToFile(txnObj types.Transaction) {
	// assign Transaction object to SignedTxn struct
	unsignedTxn := types.SignedTxn{
		Txn: txnObj,
	}

	// save unsigned Transaction to file
	err := ioutil.WriteFile("./unsigned.txn", msgpack.Encode(unsignedTxn), 0644)
	if err == nil {
		fmt.Printf("Saving unsigned transaction to file...\n")
		return
	}
	fmt.Printf("...failed to save transaction to file, error %s\n", err)
}

func readUnsigedTransactionFromFile() types.Transaction {
	// read unsigned transaction from file
	fmt.Println("Reading transaction from file...")
	bytesRead, err := ioutil.ReadFile("./unsigned.txn")
	if err != nil {
		fmt.Printf("...error reading transaction from file: %s\n", err)
	}
	fmt.Println("Decoding file bytes...")
	var unsignedTxn types.SignedTxn
	msgpack.Decode(bytesRead, &unsignedTxn)
	txnObj := unsignedTxn.Txn

	// TODO: display unsigned transaction

	return txnObj
}

func signTransaction(txnObj types.Transaction, sk ed25519.PrivateKey) []byte {
	// sign the transaction
	fmt.Println("Signing transactions...")
	signedTxnID, signedBytes, err := crypto.SignTransaction(sk, txnObj)
	if err != nil {
		fmt.Printf("...Failed to sign transaction: %s\n", err)
	}
	fmt.Println("...Signed transaction: ", signedTxnID)

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
	var signedTxn types.SignedTxn
	var signedTxnObj types.Transaction
	msgpack.Decode(signedBytes, &signedTxn)
	signedTxnObj = signedTxn.Txn

	txnJSON, err := json.MarshalIndent(signedTxnObj, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall txn data: %s\n", err)
	}
	fmt.Printf("Transaction information: %s\n", txnJSON)

	return signedBytes
}

func sendSignedTransaction(algodClient *algod.Client, signedBytes []byte) {
	// send the transaction to the network
	fmt.Println("Sending signed transction to network...")
	txID, err := algodClient.SendRawTransaction(signedBytes).Do(context.Background())
	if err != nil {
		fmt.Printf("...Failed to send transaction: %s\n", err)
		return
	}

	// wait for response
	waitForConfirmation(txID, algodClient)
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
	address, sk := getAccount(myMnemonic)

	// Create transaction object from account
	txnObj := createTransaction(algodClient, address)

	// Save unsigned transaction to file
	saveUnsignedTransactionToFile(txnObj)

	// Read the unsigned transaction from the file
	unsignedTxn := readUnsigedTransactionFromFile()

	// Sign the transaction using the mnemonic
	signedBytes := signTransaction(unsignedTxn, sk)

	// Save the signed transaction to file
	saveSignedTransactionToFile(signedBytes)

	// Read the signed transaction from file
	signedBytes = readSignedTransactionFromFile()

	// Send the transaction to the network
	sendSignedTransaction(algodClient, signedBytes)
}
