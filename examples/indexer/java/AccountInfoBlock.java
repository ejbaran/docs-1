// AccountInfoBlock.java
// requires java-algorand-sdk 1.4.0 or higher (see pom.xml)
package com.algorand.AccountInfoBlock;
import com.algorand.algosdk.v2.client.common.IndexerClient;
import com.algorand.algosdk.v2.client.common.Client;
import com.algorand.algosdk.crypto.Address;
import org.json.JSONObject;

public class AccountInfoBlock {
    public Client indexerInstance = null;
    // utility function to connect to a node
    private Client connectToNetwork(){
        final String INDEXER_API_ADDR = "localhost";
        final int INDEXER_API_PORT = 8980;       
        IndexerClient indexerClient = new IndexerClient(INDEXER_API_ADDR, INDEXER_API_PORT); 
        return indexerClient;
    }
    public static void main(String args[]) throws Exception {
        AccountInfoBlock ex = new AccountInfoBlock();
        IndexerClient indexerClientInstance = (IndexerClient)ex.connectToNetwork();
        Address account = new Address("7WENHRCKEAZHD37QMB5T7I2KWU7IZGMCC3EVAO7TQADV7V5APXOKUBILCI");
        Long round = Long.valueOf(50);
        String response = indexerClientInstance.lookupAccountByID(account).round(round).execute().toString();
        JSONObject jsonObj = new JSONObject(response.toString());
        System.out.println("Account Info for block: " + jsonObj.toString(2)); // pretty print json
    }
 }