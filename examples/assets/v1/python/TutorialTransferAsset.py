# Asset ID: 329044
import json
from algosdk import account, algod, mnemonic, transaction

# Shown for demonstration purposes. NEVER reveal secret mnemonics in practice.
# Change these values with your mnemonics
# mnemonic1 = "PASTE your phrase for account 1"
# mnemonic2 = "PASTE your phrase for account 2"
# mnemonic3 = "PASTE your phrase for account 3"

mnemonic1 = "portion never forward pill lunch organ biology weird catch curve isolate plug innocent skin grunt bounce clown mercy hole eagle soul chunk type absorb trim"
mnemonic2 = "place blouse sad pigeon wing warrior wild script problem team blouse camp soldier breeze twist mother vanish public glass code arrow execute convince ability there"
mnemonic3 = "image travel claw climb bottom spot path roast century also task cherry address curious save item clean theme amateur loyal apart hybrid steak about blanket"

# For ease of reference, add account public and private keys to
# an accounts dict.
accounts = {}
counter = 1
for m in [mnemonic1, mnemonic2, mnemonic3]:
    accounts[counter] = {}
    accounts[counter]['pk'] = mnemonic.to_public_key(m)
    accounts[counter]['sk'] = mnemonic.to_private_key(m)
    counter += 1

# Specify your node address and token. This must be updated.
# algod_address = ""  # ADD ADDRESS
# algod_token = ""  # ADD TOKEN

# algod_address = "http://hackathon.algodev.network:9100"
# algod_token = "ef920e2e7e002953f4b29a8af720efe8e4ecc75ff102b165e0472834b25832c1"

algod_address = "http://localhost:4001"
algod_token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

# Initialize an algod client
algod_client = algod.AlgodClient(algod_token, algod_address)

# Get network params for transactions.
params = algod_client.suggested_params()
first = params.get("lastRound")
last = first + 1000
gen = params.get("genesisID")
gh = params.get("genesishashb64")
min_fee = params.get("minFee")

# Utility function to wait for a transaction to be confirmed by network


def wait_for_tx_confirmation(txid):
   last_round = algod_client.status().get('lastRound')
   while True:
       txinfo = algod_client.pending_transaction_info(txid)
       if txinfo.get('round') and txinfo.get('round') > 0:
           print("Transaction {} confirmed in round {}.".format(
               txid, txinfo.get('round')))
           break
       else:
           print("Waiting for confirmation...")
           last_round += 1
           algod_client.status_after_block(last_round)


print("Account 1 address: {}".format(accounts[1]['pk']))
print("Account 2 address: {}".format(accounts[2]['pk']))
print("Account 3 address: {}".format(accounts[3]['pk']))


# your terminal output should look similar to the following
# Account 1 account: THQHGD4HEESOPSJJYYF34MWKOI57HXBX4XR63EPBKCWPOJG5KUPDJ7QJCM
# Account 1 account: AJNNFQN7DSR7QEY766V7JDG35OPM53ZSNF7CU264AWOOUGSZBMLMSKCRIU
# Account 1 account: 3ZQ3SHCYIKSGK7MTZ7PE7S6EDOFWLKDQ6RYYVMT7OHNQ4UJ774LE52AQCU


# copy in your assetID
asset_id = 329044
# transfer asset of 10 from account 1 to account 3

data = {
    "sender": accounts[1]['pk'],
    "fee": min_fee,
    "first": first,
    "last": last,
    "gh": gh,
    "receiver": accounts[3]["pk"],
    "amt": 10,
    "index": asset_id,
    "flat_fee": True
}

txn = transaction.AssetTransferTxn(**data)
stxn = txn.sign(accounts[1]['sk'])
txid = algod_client.send_transaction(stxn)
print(txid)
# Wait for the transaction to be confirmed
wait_for_tx_confirmation(txid)
# The balance should now be 10.
account_info = algod_client.account_info(accounts[3]['pk'])
print(json.dumps(account_info['assets'][str(asset_id)], indent=4))

# terminal output should look similar to this...

# BX72JE7T2NTOO33VZFK75MT4JNQXFVEUOMFKV2LRPZFH5VDL27AQ
# {
#     "creator": "THQHGD4HEESOPSJJYYF34MWKOI57HXBX4XR63EPBKCWPOJG5KUPDJ7QJCM",
#     "amount": 10,
#     "frozen": false
# }
