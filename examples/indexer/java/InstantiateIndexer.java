/**
  * requires java-algorand-sdk 1.4.2 or higher (see pom.xml)
  */
  package com.algorand.javatest.indexer;
  import com.algorand.algosdk.v2.client.common.IndexerClient;
  import com.algorand.algosdk.v2.client.common.Client;
  public class InstantiateIndexer {
      static final String host = "http://localhost";
      static final int port = 8980;
      public static void main(String args[]) throws Exception {
          IndexerClient indexerClient = new IndexerClient(host, port);
      }
   }