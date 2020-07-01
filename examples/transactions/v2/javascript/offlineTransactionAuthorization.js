const algosdk = require('algosdk');
const fs = require('fs');

// user declared mnemonic for myAccount
const myMnemonic = "boy kidney fall hamster ecology mercy inquiry vast deal normal vibrant labor couch economy embody glory possible color burger addict soap almost margin about negative" // TODO:"Your 25-word mnemonic goes here";"
const receiver = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"

// user declared algod connection parameters
const algodAddress = "http://localhost"
const algodPort = 49392 //TODO:4001;
const algodToken = "a31f09a18dbf7ad68c9e0ff22355774fb89c67ed2c4642d6c6822f9360cd7697" //TODO:"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

// function used to wait for a tx confirmation
const waitForConfirmation = async function (algodclient, txId) {
    let status = (await algodclient.status().do());
    let lastRound = status["last-round"];
    while (true) {
        const pendingInfo = await algodclient.pendingTransactionInformation(txId).do();
        if (pendingInfo["confirmed-round"] !== null && pendingInfo["confirmed-round"] > 0) {
            //Got the completed Transaction
            console.log("Transaction " + txId + " confirmed in round " + pendingInfo["confirmed-round"]);
            break;
        }
        console.log("...waiting for confirmation");
        lastRound++;
        await algodclient.statusAfterBlock(lastRound).do();
    }
}; 

// recover account from mnemonic passphrase
function getAccount(passphrase) {
    console.log("Loading signing account...")
    try {
        let myAccount = algosdk.mnemonicToSecretKey(passphrase);
        console.log("...found address: ", myAccount.addr)
        return myAccount;
    }
    catch( e ){
        console.log( e );
    }
}

async function createTransaction(algodClient, myAccount) {
    try{
        // get network suggested parameters
        let params = await algodClient.getTransactionParams().do();
        
        // make transaction
        console.log("Creating transaction...");
        amount = 1000000;
        fromAddr = myAccount
        toAddr = receiver
        let txnObj = algosdk.makePaymentTxnWithSuggestedParams(fromAddr, toAddr, amount, undefined, undefined, params);    
        console.log("...tx1: from %s to %s for %s microAlgos\n", fromAddr, toAddr, amount)
        
        return txnObj
    }
    catch( e ){
        console.log( e );
    }
}

function saveUnsignedTransactionToFile(txnObj) {
    try{
        // assign Transaction object to SignedTxn struct
        let unsignedTxn = {
            txn: txnObj.get_obj_for_encoding(),
        }

        // Save unsigned transaction object to file
        console.log("Writing unsigned transaction to './unsigned.txn'...")
        let bytesToWrite = algosdk.encodeObj(unsignedTxn);
        fs.writeFileSync('./unsigned.txn', bytesToWrite);
    }
    catch( e ){
        console.log( e );
    }
}

function readUnsigedTransactionFromFile() {
    try {
        // read unsigned transaction from file
        console.log("Reading transaction from file...");
        let bytesRead = fs.readFileSync('./unsigned.txn');  

        console.log("Decoding file bytes...");
        let unsignedTxn = algosdk.decodeObj(bytesRead);

        // get the txnObj from unsignedTxn
        let txnObj = unsignedTxn.txn;
   
        //BLOCKER: txnObj is not properly decoded. gh, rcv, snd, etc. are still byte[]. See SDK Issuee https://github.com/algorand/js-algorand-sdk/issues/114

        return txnObj;
    }
    catch( e ){
        console.log( e );
    }
}

function signTransaction(txnObj, sk) {
    try {
        console.log("Signing transactions...");
        //TODO: display transaction object prior to signing

        // sign transaction
        //BLOCKER: algosdk.signTransaction(txnObj, sk) ERROR: genesis hash must be specified and in a base64 string. See SDK Issuee https://github.com/algorand/js-algorand-sdk/issues/114
        let signedTxn = algosdk.signTransaction(txnObj, sk);
        let txId = signedTxn.txID().toString();
        console.log("...signed transaction with txID: %s", txId);

        return signedTxn;
    }
    catch( e ){
        console.log( e );
    }    
}

function saveSignedTransactionToFile(signedTxn) {
    try{
        // Save signed transaction object to file
        console.log("Writing signed transaction to './signed.txn'...")
        fs.writeFileSync('./signed.txn', signedTxn);
    }
    catch( e ){
        console.log( e );
    }
} 

function readSignedTransactionFromFile() {
    try {
        // read unsigned transaction from file
        console.log("Reading transaction from file...");
        let signedBytes = fs.readFileSync('./signed.txn');  

        return signedBytes;
    }
    catch( e ){
        console.log( e );
    }
}

function sendSignedTransaction(algodClient, signedBytes) {
    try {
        let txn = (await algodClient.sendRawTransaction(signedBytes).do());
        console.log("Sent transaction : " + tx.txId);

        // Wait for transaction to be confirmed
        await waitForConfirmation(algodClient, txn.txId)
    } 
    catch (err) {
        console.log("err", err);  
    }
}

async function offlineTransctionAuthorization() {
    // Initialize an algodClient
    let algodClient = new algosdk.Algodv2(algodToken, algodAddress, algodPort);

    // Load account from Mymnemonic
    myAccount = getAccount(myMnemonic);

	// Create transaction object from account
    txnObj = await createTransaction(algodClient, myAccount.addr);

    // Save unsigned transaction to file
    saveUnsignedTransactionToFile(txnObj);

    // Read the unsigned transaction from the file
    txnObj = readUnsigedTransactionFromFile();

    // Sign the transaction using the secret key
    signedTxn = signTransaction(txnObj, myAccount.sk);

    // Save the signed transaction to file
    saveSignedTransactionToFile(signedTxn);

    // Read the signed transaction from file
    signedBytes = readSignedTransactionFromFile();

    // Send the transaction to the network
    sendSignedTransaction(algodClient, signedBytes);
};
offlineTransctionAuthorization();