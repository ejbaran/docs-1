title: Offline Authorization

This section explains how to authorize transactions with private keys that are kept **offline**. In particular, this guide shows how to create and save transactions to a file that can then be transferred to an offline device for signing. To learn about the structure of transactions and how to authorize them in general visit the [Transactions Structure](./index.md) and [Authorizing Transactions](./signatures.md) sections, respectively.

The same methodology described here can also be used to work with [LogicSignatures](../asc1/modes.md#logic-signatures) and [Multisignatures](./signatures.md#multisignatures). All objects in the following examples use msgpack to store the transaction object ensuring interoperability with the SDKs and `goal`.

!!! info
    Storing keys _offline_ is also referred to as placing them in **cold storage**. An _online_ device that stores private keys is often referred to as a **hot wallet**.  

Algorand SDK's and `goal` support writing and reading both signed and unsigned transactions to a file. Examples of these scenarios are shown in the following code snippets.

There are three basic steps when working with transactions: create, sign and send. The offline authorization scenario assumes the signing step is performed by a disconnected _offline device_ and the create and send steps are performed by an _online device_ connected to the network. 

# Online Device (Create)

The first step is to create an unsigned transaction. This requires online connectivity to gather relevant network parameters. Other information required includes the sender, receiver and amount. The result will be an unsigned transaction file able to be exported to an offline device for signing.

## Declarations and Instantiations

First, declare the required transaction details and instantiate an algod client:

``` javascript tab="JavaScript"
const algosdk = require('algosdk');
const fs = require('fs');

// User defined constants. The receiver is the TestNet faucet.
const sender = "YOUR_ACCOUNT_TO_SEND_FROM";
const receiver = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A";

// User defined constants. The receiver is the TestNet faucet.
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

// User defined constants. The receiver is the TestNet faucet.
sender = "YOUR_ACCOUNT_TO_SEND_FROM"
receiver = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"

# User declared algod settings. These are the defaults for Sandbox.
algod_address = "http://localhost:4001"
algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

# Initialize an algod client.
algod_client = algod.AlgodClient(algod_token, algod_address)
```

``` java tab="Java"
final String SENDER = "YOUR_ACCOUNT_TO_SEND_FROM";
final String RECEIVER = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A";

// User defined constants. The receiver is the TestNet faucet.
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

// User defined constants. The receiver is the TestNet faucet.
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
AMOUNT=1000000
```

## Create Unsigned Transaction

Next, create the unsigned transaction. The result will be a transaction object `txnObj` which will be used in the next step.

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
// get network suggested parameters
txParams, err := algodClient.SuggestedParams().Do(context.Background())
txParams.FlatFee = true
var minFee uint64 = 1000
genID := txParams.GenesisID
genHash := txParams.GenesisHash
firstValidRound := uint64(txParams.FirstRoundValid)
lastValidRound := uint64(txParams.LastRoundValid)

// create transaction
fromAddr := sender
toAddr := receiver
var amount uint64 = 1000000

txnObj, err := transaction.MakePaymentTxnWithFlatFee(fromAddr, toAddr, minFee, amount, firstValidRound, lastValidRound, nil, "", genID, genHash)

fmt.Printf("...txn: from %s to %s for %v microAlgos\n", fromAddr, toAddr, amount)
```

``` bash tab="goal"
$ goal clerk send --from $SENDER --to $RECEIVER --amount $AMOUNT --out "unsigned.txn"

# This command also writes the file to "unsigned.txn"
```

## Save Unsigned Txn to File

The last step in the create process is saving the transaction to an "unsigned.txn" file. The `txnObj` is encoded using msgpack. The resulting file is interoperable between `goal` and the SDKs for ease with offline signing.

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
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

```

``` bash tab="goal"
# The "unsigned.txn" file was written by the previous `goal clerk send` command
```

# Offline Device (Sign)

Next, are the three signing steps, which are performed on a separate _offline device_. This requires the "unsigned.txn" file from the create steps and will result in "signed.txn" for use in the sending steps below.

## Read Unsigned Txn from File

First, the "unsigned.txn" file must be transported to this offline device. The file is read in, decoded from msgpack and stored as a `txnObj` for use in the next step.

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
// read unsigned transaction from file
bytesRead, err := ioutil.ReadFile("./unsigned.txn")
if err != nil {
    fmt.Printf("...error reading transaction from file: %s\n", err)
}
fmt.Println("Decoding file bytes...")
var unsignedTxn types.SignedTxn
msgpack.Decode(bytesRead, &unsignedTxn)
txnObj := unsignedTxn.Txn
```

``` bash tab="goal"
# Ensure goal has read access to the "unsigned.txn" file. Signing takes place in the next step.
```

## Sign Transaction

Now it's time to sign the `txnObj` with the sender's signing key. The SDK code snippets below import the signing key from a mnemonic passphrase, while `goal` assumes the connected wallet holds the key. This signing step converts the `txnObj` into `signedBytes` which will be saved to a file in the next step.

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
// load signing key from mnemonic passphrase
const mnemonic = "Your 25-word secret mnemonic passphrase here"
sk, err := mnemonic.ToPrivateKey(mnemonic)
if err != nil {
    fmt.Printf("...error recovering account: %s\n", err)
}
pk := sk.Public()
var a types.Address
cpk := pk.(ed25519.PublicKey)
copy(a[:], cpk[:])
fmt.Printf("...found address: %s\n", a.String())
address := a.String()

// sign the transaction
signedTxnID, signedBytes, err := crypto.SignTransaction(sk, txnObj)
if err != nil {
    fmt.Printf("...Failed to sign transaction: %s\n", err)
}
fmt.Println("...Signed transaction: ", signedTxnID)
```

``` bash tab="goal"
$ goal clerk sign --in "unsigned.txn" --out "signed.txn"
```

## Save Signed Txn to File

The final step on the _offline device_ is saving the `signedBytes` to a file "signed.txn" for later broadcast to the network.

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
// save the signed transaction to file
err := ioutil.WriteFile("./signed.stxn", signedBytes, 0644)
if err == nil {
    fmt.Printf("...Saved signed transaction to file\n")
    return
}
fmt.Printf("...Failed to save transaction to file, error %s\n", err)

```

``` bash tab="goal"
# The "signed.txn" file was created in the previous `goal clerk sign` step.
```

# Online Device (Send)

The final two steps sending the transaction are performed from the _online devive_. The "signed.txn" created in the previous step must be transported to this device.

## Read Signed Txn from File

Read in the file bytes and decode the from msgpack to yield a `signedTnxObj` which contains both a _Transaction_ `txn` and a _Signature_ `sig`. 

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
// read unsigned transaction from file
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
```

``` bash tab="goal"
# Ensure goal has read access to the "signed.txn" file. It will be sent in the final step below.
```

## Send Signed Transaction

Finally, the `signedTxnObj` is broadcast to the network. Wait for a couple seconds for confirmation the transaction was committed. 

``` javascript tab="JavaScript"
```

``` python tab="Python"
```

``` java tab="Java"
```

``` go tab="Go"
// send the transaction to the network
txID, err := algodClient.SendRawTransaction(signedBytes).Do(context.Background())
if err != nil {
	fmt.Printf("...Failed to send transaction: %s\n", err)
	return
}

// wait for response
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

```

``` bash tab="goal"
$ goal clerk rawsend --file signed.txn
```