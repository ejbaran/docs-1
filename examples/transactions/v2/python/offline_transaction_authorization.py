import json
import time
import base64
import os
from algosdk.v2client import algod
from algosdk import mnemonic
from algosdk.future import transaction
from algosdk import encoding
from algosdk import account

# user declared mnemonic for myAccount
my_mnemonic = "boy kidney fall hamster ecology mercy inquiry vast deal normal vibrant labor couch economy embody glory possible color burger addict soap almost margin about negative" # TODO:"Your 25-word mnemonic goes here";"
receiver = "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"

# user declared algod connection parameters
algod_address = "http://localhost:49392"                                          #TODO:"http:#localhost:4001"
algod_token = "a31f09a18dbf7ad68c9e0ff22355774fb89c67ed2c4642d6c6822f9360cd7697" #TODO:"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

# def that waits for a given txId to be confirmed by the network
def wait_for_confirmation(client, txid) :
    last_round = client.status().get('last-round')
    txinfo = client.pending_transaction_info(txid)
    while not (txinfo.get('confirmed-round') and txinfo.get('confirmed-round') > 0):
        print("Waiting for confirmation...")
        last_round += 1
        client.status_after_block(last_round)
        txinfo = client.pending_transaction_info(txid)
    print("Transaction {} confirmed in round {}.".format(txid, txinfo.get('confirmed-round')))
    return txinfo

# utility def to get address string
def get_account(mn) :
    print("Loading signing account...")
    sk = mnemonic.to_private_key(mn)
    address = account.address_from_private_key(sk)
    print("...found address : ", address)
    return address, sk

def create_transaction(algod_client, address) :
    print("Creating transaction...")
	# get node suggested parameters
    params = algod_client.suggested_params()
    # comment out the next two (2) lines to use suggested fees
    params.flat_fee = True
    params.fee = 1000

	# from account 1 to account 3
    sender = address
    #receiver = receiver
    amount = 1000000
    txn_obj = transaction.PaymentTxn(sender, params, receiver, amount)
    print("...txn: from {} to {} for {} microAlgos".format(sender, receiver, amount))
    print("...with txid:", txn_obj.get_txid())

def saveUnsignedTransactionToFile(txnObj) :
    print("Saving signed transction to file...")
    dir_path = os.path.dirname(os.path.realpath(__file__))
    transaction.write_to_file([txnObj], dir_path + "/unsigned.txn")


def main() :
	# Initialize an algod_client
    algod_client = algod.AlgodClient(algod_token, algod_address)

	# Load account from Mymnemonic
    address, sk = get_account(my_mnemonic)

	# Create transaction object from account
    txnObj = create_transaction(algod_client, address)

	# Save unsigned transaction to file
    saveUnsignedTransactionToFile(txnObj)

	# # Read the unsigned transaction from the file
	# unsignedTxn := readUnsigedTransactionFromFile()

	# # Sign the transaction using the mnemonic
	# signedBytes := signTransaction(unsignedTxn, sk)

	# # Save the signed transaction to file
	# saveSignedTransactionToFile(signedBytes)

	# # Read the signed transaction from file
	# signedBytes = readSignedTransactionFromFile()

	# # Send the transaction to the network
	# sendSignedTransaction(algod_client, signedBytes)

main()