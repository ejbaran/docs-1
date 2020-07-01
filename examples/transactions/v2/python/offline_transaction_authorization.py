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

# Function that waits for a given txid to be confirmed by the network
def wait_for_confirmation(client, txid):
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
    sender = address
    amount = 1000000
    txn_obj = transaction.PaymentTxn(sender, params, receiver, amount)
    print("...txn: from {} to {} for {} microAlgos".format(sender, receiver, amount))
    print("...with txid:", txn_obj.get_txid())

    return txn_obj

def save_unsigned_transaction_to_file(txn_obj) :
    print("Saving unsigned transction to file...")
    dir_path = os.path.dirname(os.path.realpath(__file__))
    transaction.write_to_file([txn_obj], dir_path + "/unsigned.txn")

def read_unsiged_transaction_from_file() :
    print("Reading unsigned transction from file...")
    dir_path = os.path.dirname(os.path.realpath(__file__))
    txns = transaction.retrieve_from_file(dir_path + "/unsigned.txn")
    unsigned_txn = txns[0]

    return unsigned_txn

def sign_transaction(unsigned_txn, sk) :
    print("Signing transaction...")
    signed_txn = unsigned_txn.sign(sk)
    txid = signed_txn.transaction.get_txid()
    print("Signed transaction with txID: {}".format(txid))

    return signed_txn

def save_signed_transaction_to_file(signed_txn) :
    print("Saving signed transction to file...")
    dir_path = os.path.dirname(os.path.realpath(__file__))
    transaction.write_to_file([signed_txn], dir_path + "/signed.txn")

def read_signed_transaction_from_file() :
    print("Reading signed transction from file...")
    dir_path = os.path.dirname(os.path.realpath(__file__))
    txns = transaction.retrieve_from_file(dir_path + "/signed.txn")
    signed_txn = txns[0]

    #TODO: Ensure signed_txn contains a signature. See issue: https://github.com/algorand/py-algorand-sdk/issues/124

    return signed_txn

def send_signed_transaction(algod_client, signed_txn) :
    # send transactions
    print("Sending transactions...")
    tx_id = algod_client.send_transactions([signed_txn])

    # wait for confirmation
    wait_for_confirmation(algod_client, tx_id) 

    confirmed_txn = algod_client.pending_transaction_info(tx_id.get_txid())
    print("Transaction information: {}".format(json.dumps(confirmed_txn, indent=4)))

def main() :
	# Initialize an algod_client
    algod_client = algod.AlgodClient(algod_token, algod_address)

	# Load account from Mymnemonic
    address, sk = get_account(my_mnemonic)

	# Create transaction object from account
    txn_obj = create_transaction(algod_client, address)

	# Save unsigned transaction to file
    save_unsigned_transaction_to_file(txn_obj)

	# Read the unsigned transaction from the file
    unsigned_txn = read_unsiged_transaction_from_file()

	# Sign the transaction using the mnemonic
    signed_txn = sign_transaction(unsigned_txn, sk)

	# Save the signed transaction to file
    save_signed_transaction_to_file(signed_txn)

	# Read the signed transaction from file
    signed_txn = read_signed_transaction_from_file()

	# Send the transaction to the network
    send_signed_transaction(algod_client, signed_txn)

main()