const algosdk = require('/Users/ryanrfox/algorand/docs/testing/js/node_modules/algosdk') //TODO:'algosdk');
const fs = require('fs');
var client = null;
// make connection to node
async function setupClient() {
    if( client == null){
        const token = "a31f09a18dbf7ad68c9e0ff22355774fb89c67ed2c4642d6c6822f9360cd7697"//TODO:"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
        const server = "http://localhost";
        const port = 49392 //TODO:4001;
        let algodClient = new algosdk.Algodv2(token, server, port);
        client = algodClient;
    } else {
        return client;
    }
    return client;
}
// recover account from mnemonic passphrase
function recoverAccount(){
    const passphrase = "boy kidney fall hamster ecology mercy inquiry vast deal normal vibrant labor couch economy embody glory possible color burger addict soap almost margin about negative"; // TODO:"Your 25-word mnemonic goes here";"
    let myAccount = algosdk.mnemonicToSecretKey(passphrase);
    return myAccount;
}
// Function used to wait for a tx confirmation
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
        lastRound++;
        await algodclient.statusAfterBlock(lastRound).do();
      }
}; 
async function writeUnsignedTransactionToFile() {

    try{
        // setup accounts
        const receiver = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A";

        // recover account from mnemonic passphrase
        let myAccount = await recoverAccount();
        console.log("My address: %s", myAccount.addr)

        // connect to node
        let algodClient = await setupClient();

        // get network suggested parameters
        let params = await algodClient.getTransactionParams().do();
        let txn = algosdk.makePaymentTxnWithSuggestedParams(myAccount.addr, receiver, 1000000, undefined, undefined, params);         
        
        // Save unsigned transaction object to file
        console.log("Writing unsigned transaction to './unsigned.txn'...")
        let bytesToWrite = algosdk.encodeObj(txn.get_obj_for_encoding());
        fs.writeFileSync('./unsigned.txn', bytesToWrite);
    }
    catch( e ){
        console.log( e );
    }
}; 
async function readUnsignedTransactionFromFile() {

    try{
        // setup connection to node
        let algodClient = await setupClient();

        // recover account
        let myAccount = await recoverAccount(); 
        console.log("My address: %s", myAccount.addr)

        // read transaction from file and sign it
        let txn = algosdk.decodeObj(fs.readFileSync('./unsigned.txn'));  

        let signedTxn = txn.signTxn(myAccount.sk)
        let txId = txn.txID().toString();

        console.log("Signed transaction with txID: %s", txId);

        // send signed transaction to node
        await algodClient.sendRawTransaction(signedTxn).do();

        // Wait for transaction to be confirmed
        await waitForConfirmation(algodClient, txId);
    } catch ( e ){
        console.log( e );
    }   
}; 

/*
async function writeSignedTransactionToFile() {

    try{
        const receiver = <transaction-receiver>;

        // setup connection to node
        let algodClient = await setupClient();
        let myAccount = await recoverAccount();
        console.log("My address: %s", myAccount.addr)

        // get network suggested parameters
        let params = await algodClient.getTransactionParams();

        // setup a transaction
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
    } catch( e ) {
        console.log(e);
    }
}; 
async function readSignedTransactionFromFile() {

    try{
        // setup connection to node
        let algodClient = await setupClient();

        // read signed transaction from file
        let stx = algosdk.decodeObj(fs.readFileSync("./signed.stxn"));

        // send signed transaction to node
        let tx = await algodClient.sendRawTransaction(stx.blob);
        console.log("Signed transaction with txID: %s", tx.txId);

        // Wait for transaction to be confirmed
        await waitForConfirmation(algodClient, tx.txId);
    } catch( e ) {
        console.log(e);
    }   
}; 
*/
async function testUnsigned(){
    await writeUnsignedTransactionToFile();
   // await readUnsignedTransactionFromFile();
}
/*
async function testSigned(){
    await writeSignedTransactionToFile();
    await readSignedTransactionFromFile();
}*/
testUnsigned();
//testSigned();