// AccountsAssetIDMinBalance.java
// requires java-algorand-sdk 1.4.0 or higher (see pom.xml)
package com.algorand.AccountsAssetIDMinBalance;

import com.algorand.algosdk.v2.client.common.IndexerClient;
import com.algorand.algosdk.v2.client.common.Client;
import org.json.JSONObject;

public class AccountsAssetIDMinBalance {
    public Client indexerInstance = null;
    // utility function to connect to a node
    private Client connectToNetwork(){
        final String INDEXER_API_ADDR = "localhost";
        final int INDEXER_API_PORT = 8980;       
        IndexerClient indexerClient = new IndexerClient(INDEXER_API_ADDR, INDEXER_API_PORT); 
        return indexerClient;
    }
    public static void main(String args[]) throws Exception {
        AccountsAssetIDMinBalance ex = new AccountsAssetIDMinBalance();
        IndexerClient indexerClientInstance = (IndexerClient)ex.connectToNetwork();
        Long asset_id = Long.valueOf(312769);
        Long currencyGreaterThan = Long.valueOf(100);
        // searches for asset greater than currencyGreaterThan
        String response = indexerClientInstance.searchForAccounts()
                .assetId(asset_id)
                .currencyGreaterThan(currencyGreaterThan).execute().toString();
        JSONObject jsonObj = new JSONObject(response.toString());
        System.out.println("Account Info for Asset Min Balance: " + jsonObj.toString(2)); // pretty print json
    }
 }