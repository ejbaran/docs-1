title: Offline Authorization

This section explains how to authorize transactions with private keys that are kept **offline**. In particular, this guide shows how to create and save transactions to a file that can then be transferred to an offline device for signing. To learn about the structure of transactions and how to authorize them in general visit the [Transactions Structure](./index.md) and [Authorizing Transactions](./signatures.md) sections, respectively.

The same methodology described here can also be used to work with [LogicSignatures](../asc1/modes.md#logic-signatures) and [Multisignatures](./signatures.md#multisignatures). All objects in the following examples use msgpack to store the transaction object ensuring interoperability with the SDKs and `goal`.

!!! info
    Storing keys _offline_ is also referred to as placing them in **cold storage**. An _online_ device that stores private keys is often referred to as a **hot wallet**.  

Algorand SDK's and `goal` support writing and reading both signed and unsigned transactions to a file. Examples of these scenarios are shown in the following code snippets.

There are three basic steps when working with transactions: create, sign and send. The offline authorization scenario assumes the signing step is performed by a disconnected _offline device_ and the create and send steps are performed by an _online device_ connected to the network. 

# Online Device (Create)

The first step is to create an unsigned transaction. This requires online connectivity to gather relevant network parameters. Other information required includes the sender, reciever and amount. The result will be an unsigned transaction file able to be exported to an offline device for signing.

## Initializations

First, gather the required transaction details:

``` javascript tab="JavaScript"
const algosdk = require('algosdk');
const fs = require('fs');

// User declared accounts. The receiver is the TestNet faucet.
const sender = "YOUR_ACCOUNT_TO_SEND_FROM";
const receiver = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A";

// User declared algod settings. These are the defaults for Sandbox.
const algodAddress = "http://localhost";
const algodPort = 4001;
const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

// Initialize and algod client.
let algodClient = new algosdk.Algodv2(algodToken, algodAddress, algodPort);

```

``` python tab="Python"
import json
import time
import base64
import os
from algosdk.v2client import algod
from algosdk import mnemonic
from algosdk.future import transaction
from algosdk import encoding
from algosdk import account

# User declared accounts. The receiver is the TestNet faucet.
sender = "YOUR_ACCOUNT_TO_SEND_FROM"
receiver = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"

# User declared algod settings. These are the defaults for Sandbox.
algod_address = "http://localhost:4001"
algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

# Initialize an algod client.
algod_client = algod.AlgodClient(algod_token, algod_address)
```

``` java tab="Java"
final String SENDER = "YOUR_ACCOUNT_TO_SEND_FROM"
final String RECEIVER = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"

// User declared algod settings. These are the defaults for Sandbox.
final String ALGOD_API_ADDR = "http://localhost";
final Integer ALGOD_PORT = 4001;
final String ALGOD_API_TOKEN = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

// Initialize an algod client.
AlgodClient algodClient = (AlgodClient) new AlgodClient(ALGOD_API_ADDR, ALGOD_PORT, ALGOD_API_TOKEN);
```

``` go tab="Go"
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

// User declared accounts. The receiver is the TestNet faucet.
const sender = "YOUR_ACCOUNT_TO_SEND_FROM"
const receiver = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"

// User declared algod settings. These are the defaults for Sandbox.
const algodAddress = "http://localhost"
const algodPort = 4001
const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

// Initialize and algod client.
algodClient, err := algod.MakeClient(algodAddress, algodToken)
```

``` bash tab="goal"
# User declared accounts. The receiver is the TestNet faucet.
SENDER="YOUR_ACCOUNT_TO_SEND_FROM"
RECEIVER="GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"
```

## Create Unsigned Transaction

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
```

``` bash tab="goal"
```

## Save Unsigned Txn to File

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
```

``` bash tab="goal"
```

# Offline Device (Sign)

## Read Unsigned Txn from File

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
```

``` bash tab="goal"
```

## Sign Transaction

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
```

``` bash tab="goal"
```

## Save Signed Txn to File

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
```

``` bash tab="goal"
```

# Online Device (Send)

## Read Signed Txn from File

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
```

``` bash tab="goal"
```

## Send Signed Transaction

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
```

``` bash tab="goal"
```


# Unsigned Transaction File Operations
Algorand SDK's and `goal` support writing and reading both signed and unsigned transactions to a file. Examples of these scenarios are shown in the following code snippets. There are three basic steps *create, sign* and *send*. 

Unsigned transactions require the transaction object to be created before writing to a file.


``` javascript tab="JavaScript"
	let txn = {
		"from": myAccount.addr,
		"to": receiver,
		"fee": params.minFee,
		"flatFee": true,
		"amount": 1000000,
		"firstRound": params.lastRound,
		"lastRound": params.lastRound + 1000,
		"genesisID": params.genesisID,
		"genesisHash": params.genesishashb64
	};
	// Save transaction to file
    fs.writeFileSync('./unsigned.txn', algosdk.encodeObj( txn ));
    
	// read transaction from file and sign it
    let txn = algosdk.decodeObj(fs.readFileSync('./unsigned.txn')); 
	let signedTxn = algosdk.signTransaction(txn, myAccount.sk);
	let txId = signedTxn.txID;
	
	// send signed transaction to node
	await algodClient.sendRawTransaction(signedTxn.blob);         
```

``` python tab="Python"
	# create transaction
	receiver = <transaction-receiver>
	data = {
		"sender": my_address,
		"receiver": receiver,
		"fee": params.get('minFee'),
		"flat_fee": True,
		"amt": <amount>,
		"first": params.get('lastRound'),
		"last": params.get('lastRound') + 1000,
		"gen": params.get('genesisID'),
		"gh": params.get('genesishashb64')
	}
	txn = transaction.PaymentTxn(**data)

	# write to file
	dir_path = os.path.dirname(os.path.realpath(__file__))
	transaction.write_to_file([txn], dir_path + "/unsigned.txn")

	# read from file
	txns = transaction.retrieve_from_file("./unsigned.txn")

	# sign and submit transaction
	txn = txns[0]
	signed_txn = txn.sign(private_key)
	txid = signed_txn.transaction.get_txid()
	algod_client.send_transaction(signed_txn)    
```

``` java tab="Java"
    BigInteger amount = BigInteger.valueOf(200000);
    BigInteger lastRound = firstRound.add(BigInteger.valueOf(1000));  
    Transaction tx = new Transaction(new Address(SRC_ADDR),  
            BigInteger.valueOf(1000), firstRound, lastRound, 
            null, amount, new Address(DEST_ADDR), genId, genesisHash);
    // save as signed even though it has not been
    SignedTransaction stx = new SignedTransaction();
    stx.tx = tx;  
    // Save transaction to a file 
    Files.write(Paths.get("./unsigned.txn"), Encoder.encodeToMsgPack(stx));

    // read transaction from file
    SignedTransaction decodedTransaction = Encoder.decodeFromMsgPack(
        Files.readAllBytes(Paths.get("./unsigned.txn")), 
        SignedTransaction.class);            
    Transaction tx = decodedTransaction.tx;          

    // recover account    
    String SRC_ACCOUNT = <25-word-passphrase>;
    Account src = new Account(SRC_ACCOUNT);

    // sign transaction
    SignedTransaction signedTx = src.signTransaction(tx);
    byte[] encodedTxBytes = Encoder.encodeToMsgPack(signedTx);
            
    // submit the encoded transaction to the network
    TransactionID id = algodApiInstance.rawTransaction(encodedTxBytes);
```

``` go tab="Go"
	tx, err := transaction.MakePaymentTxn(addr, toAddr, 1, 100000,
		 txParams.LastRound, txParams.LastRound+100, nil, "", 
		 genID, txParams.GenesisHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
    }
    // save as signed tx object without sig
    unsignedTx := types.SignedTxn{
		Txn:  tx,
	 }

	// save unsigned Transaction to file
	err = ioutil.WriteFile("./unsigned.txn", msgpack.Encode(unsignedTx), 0644)
	if err == nil {
		fmt.Printf("Saved unsigned transaction to file\n")
		return
    }

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
```


``` goal tab="goal"
$ goal clerk send --from=<my-account> --to=GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A --fee=1000 --amount=1000000 --out="unsigned.txn"

$ goal clerk sign --infile unsigned.txn --outfile signed.txn

$ goal clerk rawsend --filename signed.txn

```
# Signed Transaction File Operations 
Signed Transactions are similar, but require an account to sign the transaction before writing it to a file.

``` javascript tab="JavaScript"
	let txn = {
		"from": myAccount.addr,
		"to": receiver,
		"fee": params.minFee,
		"flatFee": true,
		"amount": 1000000,
		"firstRound": params.lastRound,
		"lastRound": params.lastRound + 1000,
		"genesisID": params.genesisID,
		"genesisHash": params.genesishashb64
	};

	// sign transaction and write to file
	let signedTxn = algosdk.signTransaction(txn, myAccount.sk);
    fs.writeFileSync('./signed.stxn', algosdk.encodeObj( signedTxn ));
    
	// read signed transaction from file
	let stx = algosdk.decodeObj(fs.readFileSync("./signed.stxn"));
		
	// send signed transaction to node
	let tx = await algodClient.sendRawTransaction(stx.blob);    
```

``` python tab="Python"
	# create transaction
    receiver = <transaction-receiver>
    data = {
        "sender": my_address,
        "receiver": receiver,
        "fee": params.get('minFee'),
        "flat_fee": True,
        "amt": <amount>,
        "first": params.get('lastRound'),
        "last": params.get('lastRound') + 1000,
        "gen": params.get('genesisID'),
        "gh": params.get('genesishashb64')
    }
    txn = transaction.PaymentTxn(**data)

    # sign transaction
    signed_txn = txn.sign(private_key)

    # write to file
    dir_path = os.path.dirname(os.path.realpath(__file__))
    transaction.write_to_file([signed_txn], dir_path + "/signed.txn")

	# read signed transaction from file
	txns = transaction.retrieve_from_file("./signed.txn")
	signed_txn = txns[0]
	txid = signed_txn.transaction.get_txid()
	print("Signed transaction with txID: {}".format(txid))
	
	# send transaction to network
	algod_client.send_transaction(signed_txn)    
```

``` java tab="Java"
    // create transaction 
    BigInteger amount = BigInteger.valueOf(200000);
    BigInteger lastRound = firstRound.add(BigInteger.valueOf(1000));  
    Transaction tx = new Transaction(new Address(SRC_ADDR),  
            BigInteger.valueOf(1000), firstRound, lastRound, 
            null, amount, new Address(DEST_ADDR), genId, genesisHash);

    // recover account    
    String SRC_ACCOUNT = <25-word-passphrase>;                    
    Account src = new Account(SRC_ACCOUNT);

    // sign transaction
    SignedTransaction signedTx = src.signTransaction(tx);                    

    // save signed transaction to  a file 
    Files.write(Paths.get("./signed.txn"), Encoder.encodeToMsgPack(signedTx));

    //Read the transaction from a file 
    SignedTransaction decodedSignedTransaction = Encoder.decodeFromMsgPack(
        Files.readAllBytes(Paths.get("./signed.txn")), SignedTransaction.class);   
    System.out.println("Signed transaction with txid: " + decodedSignedTransaction.transactionID);           

    // Msgpack encode the signed transaction
    byte[] encodedTxBytes = Encoder.encodeToMsgPack(decodedSignedTransaction);

    //submit the encoded transaction to the network
    TransactionID id = algodApiInstance.rawTransaction(encodedTxBytes);    
```

``` go tab="Go"
	tx, err := transaction.MakePaymentTxn(addr, toAddr, 1, 100000,
		 txParams.LastRound, txParams.LastRound+100, nil, "", 
		 genID, txParams.GenesisHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}

	//Sign the Transaction
	txid, stx, err := crypto.SignTransaction(sk, tx)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
		return
	}
	fmt.Printf("Made signed transaction with TxID %s: %x\n", txid, stx)

	//Save the signed transaction to file
	err = ioutil.WriteFile("./signed.stxn", msgpack.Encode(stx), 0644)
	if err == nil {
		fmt.Printf("Saved signed transaction to file\n")
		return
    }
    
	// read unsigned transaction from file
	dat, err := ioutil.ReadFile("./signed.stxn")
	if err != nil {
		fmt.Printf("Error reading signed transaction from file: %s\n", err)
		return
	}
	var signedTx []byte 
	msgpack.Decode(dat, &signedTx)
	
	// send the transaction to the network
	sendResponse, err := algodClient.SendRawTransaction(signedTx)
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}

```

``` goal tab="goal"
$ goal clerk rawsend --filename signed.txn
```
