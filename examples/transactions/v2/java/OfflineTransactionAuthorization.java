package com.algorand.javatest;

import java.nio.file.Files;
import java.nio.file.Paths;

import com.algorand.algosdk.account.Account;
import com.algorand.algosdk.crypto.Signature;
import com.algorand.algosdk.transaction.SignedTransaction;
import com.algorand.algosdk.transaction.Transaction;
import com.algorand.algosdk.util.Encoder;
import com.algorand.algosdk.v2.client.common.AlgodClient;
import com.algorand.algosdk.v2.client.common.Response;
import com.algorand.algosdk.v2.client.model.PendingTransactionResponse;
import com.algorand.algosdk.v2.client.model.TransactionParametersResponse;


public class OfflineTransactionAuthorization {

    // utility function to return and account from a mnemonic
    public static Account getAccount( String mnemonic ) throws Exception{
        com.algorand.algosdk.account.Account myAccount = new Account(mnemonic);
        System.out.println("My Address: " + myAccount.getAddress());

        return myAccount;
    }

    // utility function to wait on a transaction to be confirmed    
    public static void waitForConfirmation( AlgodClient algodClient, String txID ) throws Exception{
        Long lastRound = algodClient.GetStatus().execute().body().lastRound;
        while(true) {
            try {
                //Check the pending tranactions
                Response<PendingTransactionResponse> pendingInfo = algodClient.PendingTransactionInformation(txID).execute();
                if (pendingInfo.body().confirmedRound != null && pendingInfo.body().confirmedRound > 0) {
                    //Got the completed Transaction
                    System.out.println("Transaction " + txID + " confirmed in round " + pendingInfo.body().confirmedRound);
                    break;
                } 
                System.out.println("Waiting for confirmation...");
                lastRound++;
                algodClient.WaitForBlock(lastRound).execute();
            } catch (Exception e) {
                throw( e );
            }
        }
    }
    
    public static Transaction createTransaction(AlgodClient algodClient, Account myAccount) throws Exception{
        final String RECEIVER = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A";
        System.out.println("Creating transaction...");

        // Get suggested parameters from the node
        TransactionParametersResponse params = algodClient.TransactionParams().execute().body();

        // create transaction

        Transaction txnObj = Transaction.PaymentTransactionBuilder()
        .sender(myAccount.getAddress().toString())
        .receiver(RECEIVER)
        .amount(100000)
        .suggestedParams(params)
        .build();

        System.out.println("...tx1: from " + txnObj.sender + " to " + txnObj.receiver + " for " + txnObj.amount + " microAlgos");

                    
        return txnObj;
    }

    public static void saveUnsignedTransactionToFile(Transaction txnObj) throws Exception{
        System.out.println("Saving transaction to file...");
        // create new (null) Signature and SignedTransaction object to wrap the transaction object prior to writing
        Signature sig = new Signature();
        SignedTransaction signedTxn = new SignedTransaction(txnObj, sig, "");

        // Save transaction to a file 
        Files.write(Paths.get("./unsigned.txn"), Encoder.encodeToMsgPack(signedTxn));
        System.out.println("...transaction written to file");
    }

    public static Transaction readUnsigedTransactionFromFile() throws Exception{
        System.out.println("Reading transaction from file...");
        SignedTransaction decodedTxn = Encoder.decodeFromMsgPack(
            Files.readAllBytes(Paths.get("./unsigned.txn")), SignedTransaction.class);            
        Transaction txnObj = decodedTxn.tx;   
        
        // TODO: display file contents
        //System.out.println(txnObj);

        return txnObj;
    }

    public static SignedTransaction signTransaction(Transaction txnObj, Account myAccount) throws Exception{
        System.out.println("Signing transation...");
        SignedTransaction signedTxn = myAccount.signTransaction(txnObj);
        return signedTxn;
    }

    public static void saveSignedTransactionToFile(SignedTransaction signedTxn) throws Exception{
        System.out.println("Writing transation to file...");

        // Save transaction to a file 
        Files.write(Paths.get("./signed.txn"), Encoder.encodeToMsgPack(signedTxn));
        System.out.println("...transaction written to file");
    }

    public static byte[] readSignedTransactionFromFile() throws Exception{
        System.out.println("Reading signed transaction from file...");
        byte[] signedBytes = Files.readAllBytes(Paths.get("./signed.txn"));            

        return signedBytes;
    }

    public static void sendSignedTransaction(AlgodClient algodClient, byte[] signedBytes) throws Exception{
        System.out.println("Sending signed transaction...");

        String id = algodClient.RawTransaction().rawtxn(signedBytes).execute().body().txId;
        System.out.println("Successfully sent tx with ID: " + id);

        // Wait for transaction confirmation
        waitForConfirmation(algodClient, id);

        //Read the transaction
        PendingTransactionResponse pTrx = algodClient.PendingTransactionInformation(id).execute().body();
        System.out.println("Transaction information: " + pTrx.toString());
        
}
    public static void main(String args[]) throws Exception {

        // user declared mnemonic for myAccount
        final String MY_MNEMONIC = "boy kidney fall hamster ecology mercy inquiry vast deal normal vibrant labor couch economy embody glory possible color burger addict soap almost margin about negative"; // TODO:"Your 25-word mnemonic goes here";"

        // Initialize an algod client
        final String ALGOD_API_ADDR = "http://localhost"; //TODO:"http://localhost:4001"
        final Integer ALGOD_PORT = 49392; //TODO: 4001;
        final String ALGOD_API_TOKEN = "a31f09a18dbf7ad68c9e0ff22355774fb89c67ed2c4642d6c6822f9360cd7697"; //TODO:"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
        AlgodClient algodClient = (AlgodClient) new AlgodClient(ALGOD_API_ADDR, ALGOD_PORT, ALGOD_API_TOKEN);

        // Load account from Mymnemonic
        Account myAccount = getAccount(MY_MNEMONIC);

        // Create transaction object from account
        Transaction txnObj = createTransaction(algodClient, myAccount);

        // Save unsigned transaction to file
        saveUnsignedTransactionToFile(txnObj);

        // Read the unsigned transaction from the file
        txnObj = readUnsigedTransactionFromFile();

        // Sign the transaction using the mnemonic
        SignedTransaction signedTxn = signTransaction(txnObj, myAccount);

        // Save the signed transaction to file
        saveSignedTransactionToFile(signedTxn);

        // Read the signed transaction from file
        byte[] signedBytes = readSignedTransactionFromFile();

        // Send the transaction to the network
        sendSignedTransaction(algodClient, signedBytes);
    }

}