package com.algorand.javatest;

import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Paths;

import com.algorand.algosdk.account.Account;
import com.algorand.algosdk.algod.client.AlgodClient;
import com.algorand.algosdk.algod.client.api.AlgodApi;
import com.algorand.algosdk.algod.client.auth.ApiKeyAuth;
import com.algorand.algosdk.algod.client.model.TransactionID;
import com.algorand.algosdk.algod.client.model.TransactionParams;
import com.algorand.algosdk.crypto.Address;
import com.algorand.algosdk.crypto.Digest;
import com.algorand.algosdk.transaction.SignedTransaction;
import com.algorand.algosdk.transaction.Transaction;
import com.algorand.algosdk.util.Encoder;

public class OfflineTransactionAuthorization {
    public AlgodApi algodApiInstance = null;

    // utility function to connect to a node
    private AlgodApi connectToNetwork(){

        // Initialize an algod client
        final String ALGOD_API_ADDR = <algod-address>;
        final String ALGOD_API_TOKEN = <algod-token>;

        AlgodClient client = (AlgodClient) new AlgodClient().setBasePath(ALGOD_API_ADDR);
        ApiKeyAuth api_key = (ApiKeyAuth) client.getAuthentication("api_key");
        api_key.setApiKey(ALGOD_API_TOKEN);
        algodApiInstance = new AlgodApi(client);   
        return algodApiInstance;
    }
    // utility function to wait on a transaction to be confirmed    
    public void waitForConfirmation( String txID ) throws Exception{
        if( algodApiInstance == null ) connectToNetwork();
        while(true) {
            try {
                //Check the pending tranactions
                com.algorand.algosdk.algod.client.model.Transaction pendingInfo = algodApiInstance.pendingTransactionInformation(txID);
                if (pendingInfo.getRound() != null && pendingInfo.getRound().longValue() > 0) {
                    //Got the completed Transaction
                    System.out.println("Transaction " + pendingInfo.getTx() + " confirmed in round " + pendingInfo.getRound().longValue());
                    break;
                } 
                algodApiInstance.waitForBlock(BigInteger.valueOf( algodApiInstance.getStatus().getLastRound().longValue() +1 ) );
            } catch (Exception e) {
                throw( e );
            }
        }

    }
    public void writeUnsignedTransaction(){

        // connect to node
        if( algodApiInstance == null ) connectToNetwork();

        final String DEST_ADDR = <transaction-reciever>;
        final String SRC_ADDR = <transaction-sender>;

        try { 
            // Get suggested parameters from the node
            TransactionParams params = algodApiInstance.transactionParams();                     
            BigInteger firstRound = params.getLastRound();
            String genId = params.getGenesisID();
            Digest genesisHash = new Digest(params.getGenesishashb64());

            // create transaction
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
            System.out.println("Transaction written to a file");
        } catch (Exception e) { 
            System.out.println("Save Exception: " + e); 
        }

    }
    public void readUnsignedTransaction(){

        try {
            // connect to node
            if( algodApiInstance == null ) connectToNetwork();

            // read transaction from file
            SignedTransaction decodedTransaction = Encoder.decodeFromMsgPack(
                Files.readAllBytes(Paths.get("./unsigned.txn")), SignedTransaction.class);            
            Transaction tx = decodedTransaction.tx;           

            // recover account    
            String SRC_ACCOUNT = <25-word-passphrase>;
            Account src = new Account(SRC_ACCOUNT);

            // sign transaction
            SignedTransaction signedTx = src.signTransaction(tx);
            byte[] encodedTxBytes = Encoder.encodeToMsgPack(signedTx);

            // submit the encoded transaction to the network
            TransactionID id = algodApiInstance.rawTransaction(encodedTxBytes);
            System.out.println("Successfully sent tx with id: " + id);
            waitForConfirmation(id.getTxId());

        } catch (Exception e) {
            System.out.println("Submit Exception: " + e); 
        }


    }
    public void writeSignedTransaction(){

        // connect to node
        if( algodApiInstance == null ) connectToNetwork();

        final String DEST_ADDR = <transaction-reciever>;
        final String SRC_ADDR = <transaction-sender>;;

        try { 

            // Get suggested parameters from the node
            TransactionParams params = algodApiInstance.transactionParams();
            BigInteger firstRound = params.getLastRound();
            String genId = params.getGenesisID();
            Digest genesisHash = new Digest(params.getGenesishashb64());

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
        } catch (Exception e) { 
            System.out.println("Save Exception: " + e); 
        }

    }

    public void readSignedTransaction(){

        try {
            // connect to a node
            if( algodApiInstance == null ) connectToNetwork();

            //Read the transaction from a file 
            SignedTransaction decodedSignedTransaction = Encoder.decodeFromMsgPack(
                Files.readAllBytes(Paths.get("./signed.txn")), SignedTransaction.class);   
            System.out.println("Signed transaction with txid: " + decodedSignedTransaction.transactionID);           

            // Msgpack encode the signed transaction
            byte[] encodedTxBytes = Encoder.encodeToMsgPack(decodedSignedTransaction);

            //submit the encoded transaction to the network
            TransactionID id = algodApiInstance.rawTransaction(encodedTxBytes);
            System.out.println("Successfully sent tx with id: " + id); 
            waitForConfirmation(id.getTxId());

        } catch (Exception e) {
            System.out.println("Submit Exception: " + e); 
        }


    }
    public static void main(String args[]) throws Exception {
        OfflineTransactionAuthorization mn = new OfflineTransactionAuthorization();
        mn.writeUnsignedTransaction();
        mn.readUnsignedTransaction();

        //mn.writeSignedTransaction();
        //mn.readSignedTransaction();

    }

}