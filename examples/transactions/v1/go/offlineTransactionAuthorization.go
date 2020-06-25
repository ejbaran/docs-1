package main

import (
    "fmt"
    "io/ioutil"

    "golang.org/x/crypto/ed25519"

    "github.com/algorand/go-algorand-sdk/client/algod"
    "github.com/algorand/go-algorand-sdk/crypto"
    "github.com/algorand/go-algorand-sdk/mnemonic"
    "github.com/algorand/go-algorand-sdk/encoding/msgpack"
    "github.com/algorand/go-algorand-sdk/transaction"
    "github.com/algorand/go-algorand-sdk/types"
)

// Function that waits for a given txId to be confirmed by the network
func waitForConfirmation(algodClient algod.Client, txID string) {
    for {
        pt, err := algodClient.PendingTransactionInformation(txID)
        if err != nil {
            fmt.Printf("waiting for confirmation... (pool error, if any): %s\n", err)
            continue
        }
        if pt.ConfirmedRound > 0 {
            fmt.Printf("Transaction "+pt.TxID+" confirmed in round %d\n", pt.ConfirmedRound)
            break
        }
        nodeStatus, err := algodClient.Status()
        if err != nil {
            fmt.Printf("error getting algod status: %s\n", err)
            return
        }
        algodClient.StatusAfterBlock( nodeStatus.LastRound + 1)
    }
}
// utility function to recover account and return sk and address
func recoverAccount()(string, ed25519.PrivateKey) {
    const passphrase = <your-25-word-mnemonic>

    sk, err := mnemonic.ToPrivateKey(passphrase)
    if err != nil {
        fmt.Printf("error recovering account: %s\n", err)
        return "", nil
    }
    pk := sk.Public()
    var a types.Address
    cpk := pk.(ed25519.PublicKey)
    copy(a[:], cpk[:])
    fmt.Printf("Address: %s\n", a.String()) 
    address := a.String()
    return address, sk 
}
// utility funciton to setup connection to node
func setupConnection()( algod.Client ){
    const algodToken = <algod-token>
    const algodAddress = <algod-address>
    algodClient, err := algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        fmt.Printf("failed to make algod client: %s\n", err)
    }
    return algodClient
}

func saveUnsignedTransaction() {

    // setup connection
    algodClient := setupConnection()

    // recover account for example
    addr, _ := recoverAccount();

    // get network suggested parameters
    txParams, err := algodClient.SuggestedParams()
    if err != nil {
        fmt.Printf("error getting suggested tx params: %s\n", err)
        return
    }

    // create transaction
    toAddr := <transaction-receiver>
    genID := txParams.GenesisID
    tx, err := transaction.MakePaymentTxn(addr, toAddr, 1, 100000,
        txParams.LastRound, txParams.LastRound+100, nil, "", 
        genID, txParams.GenesisHash)
    if err != nil {
        fmt.Printf("Error creating transaction: %s\n", err)
        return
    }
    unsignedTx := types.SignedTxn{
        Txn:  tx,
    }        

    // save unsigned Transaction to file
    err = ioutil.WriteFile("./unsigned.txn", msgpack.Encode(unsignedTx), 0644)
    if err == nil {
        fmt.Printf("Saved unsigned transaction to file\n")
        return
    }
    fmt.Printf("Failed in saving trx to file, error %s\n", err)

}
func readUnsignedTransaction(){

    // setup connection
    algodClient := setupConnection()

    // read unsigned transaction from file
    dat, err := ioutil.ReadFile("./unsigned.txn")
    if err != nil {
        fmt.Printf("Error reading transaction from file: %s\n", err)
        return
    }
    var unsignedTxRaw types.SignedTxn 
    var unsignedTxn types.Transaction

    msgpack.Decode(dat, &unsignedTxRaw)

    unsignedTxn = unsignedTxRaw.Txn

    // recover account and sign transaction
    addr, sk := recoverAccount();
    fmt.Printf("Address is: %s\n", addr)
    txid, stx, err := crypto.SignTransaction(sk, unsignedTxn)
    if err != nil {
        fmt.Printf("Failed to sign transaction: %s\n", err)
        return
    }
    fmt.Printf("Transaction id: %s\n", txid)

    // send transaction to the network
    sendResponse, err := algodClient.SendRawTransaction(stx)
    if err != nil {
        fmt.Printf("failed to send transaction: %s\n", err)
        return
    }
    fmt.Printf("Transaction ID: %s\n", sendResponse.TxID)
    waitForConfirmation(algodClient, sendResponse.TxID)
}


func saveSignedTransaction() {

    // setup connection
    algodClient := setupConnection()

    // recover account
    addr, sk := recoverAccount();

    // get network suggested parameters
    txParams, err := algodClient.SuggestedParams()
    if err != nil {
        fmt.Printf("error getting suggested tx params: %s\n", err)
        return
    }

    // create transaction
    toAddr := <transaction-receiver>
    genID := txParams.GenesisID
    tx, err := transaction.MakePaymentTxn(addr, toAddr, 1, 100000,
        txParams.LastRound, txParams.LastRound+100, nil, "", 
        genID, txParams.GenesisHash)
    if err != nil {
        fmt.Printf("Error creating transaction: %s\n", err)
        return
    }

    // sign the Transaction, msgpack encoding happens in sign
    txid, stx, err := crypto.SignTransaction(sk, tx)
    if err != nil {
        fmt.Printf("Failed to sign transaction: %s\n", err)
        return
    }
    fmt.Printf("Made signed transaction with TxID %s: %x\n", txid, stx)

    //Save the signed transaction to file
    err = ioutil.WriteFile("./signed.stxn", stx, 0644)
    if err == nil {
        fmt.Printf("Saved signed transaction to file\n")
        return
    }
    fmt.Printf("Failed in saving trx to file, error %s\n", err)

}
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
func main() {
    //saveUnsignedTransaction()
    //readUnsignedTransaction()

    saveSignedTransaction()
    readSignedTransaction()

}    