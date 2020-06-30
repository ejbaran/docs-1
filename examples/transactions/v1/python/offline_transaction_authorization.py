
#/usr/bin/python3
import json
import time
import base64
import os
from algosdk import algod
from algosdk import mnemonic
from algosdk import transaction

def connect_to_network():
    algod_address = <algod-address>
    algod_token = <algod-token>
    algod_client = algod.AlgodClient(algod_token, algod_address)
    return algod_client

def wait_for_confirmation( algod_client, txid ):
    while True:
        txinfo = algod_client.pending_transaction_info(txid)
        if txinfo.get('round') and txinfo.get('round') > 0:
            print("Transaction {} confirmed in round {}.".format(txid, txinfo.get('round')))
            break
        else:
            print("Waiting for confirmation...")
            algod_client.status_after_block(algod_client.status().get('lastRound') +1)

def write_unsigned():
    # setup none connection
    algod_client = connect_to_network()

    # recover account
    passphrase = <25-word-passphrase>
    private_key = mnemonic.to_private_key(passphrase)
    my_address = mnemonic.to_public_key(passphrase)
    print("My address: {}".format(my_address))

    # get suggested parameters
    params = algod_client.suggested_params()

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


def read_unsigned():

    # setup node connection
    algod_client = connect_to_network()

    # recover account
    passphrase = "25-word-passphrase"
    private_key = mnemonic.to_private_key(passphrase)
    my_address = mnemonic.to_public_key(passphrase)
    print("My address: {}".format(my_address))

    # read from file
    txns = transaction.retrieve_from_file("./unsigned.txn")

    # sign and submit transaction
    txn = txns[0]
    signed_txn = txn.sign(private_key)
    txid = signed_txn.transaction.get_txid()
    print("Signed transaction with txID: {}".format(txid))
    algod_client.send_transaction(signed_txn)

    # wait for confirmation
    wait_for_confirmation( algod_client, txid)

def write_signed():

    # setup connection to node
    algod_client = connect_to_network()

    # recovere account
    passphrase = <25-word-passphrase>
    private_key = mnemonic.to_private_key(passphrase)
    my_address = mnemonic.to_public_key(passphrase)
    print("My address: {}".format(my_address))

    # get node suggested parameters
    params = algod_client.suggested_params()

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


def read_signed():

    # set up connection to node
    algod_client = connect_to_network()

    # read signed transaction from file
    txns = transaction.retrieve_from_file("./signed.txn")
    signed_txn = txns[0]
    txid = signed_txn.transaction.get_txid()
    print("Signed transaction with txID: {}".format(txid))

    # send transaction to network
    algod_client.send_transaction(signed_txn)

    # wait for confirmation
    wait_for_confirmation( algod_client, txid)

# Test Runs     
#write_unsigned()
#read_unsigned()
write_signed()
read_signed()    