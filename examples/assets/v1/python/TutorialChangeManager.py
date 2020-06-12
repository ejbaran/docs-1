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

# Update manager address.

# The current manager(Account 2) issues an asset configuration transaction that assigns Account 1 as the new manager.
# Keep reserve, freeze, and clawback address same as before, i.e. account 2

# paste in your asset_id
asset_id = 329044

data = {
    "sender": accounts[2]['pk'],
    "fee": min_fee,
    "first": first,
    "last": last,
    "gh": gh,
    "index": asset_id,
    "manager": accounts[1]['pk'],
    "reserve": accounts[2]['pk'],
    "freeze": accounts[2]['pk'],
    "clawback": accounts[2]['pk'],
    "flat_fee": True
}
txn = transaction.AssetConfigTxn(**data)
# sign by the current manager
stxn = txn.sign(accounts[2]['sk'])
txid = algod_client.send_transaction(stxn)
print("Transaction ID: ",txid)
# Wait for the transaction to be confirmed
wait_for_tx_confirmation(txid)

# Check asset info to view change in management.
asset_info = algod_client.asset_info(asset_id)
print(json.dumps(asset_info, indent=4))

# terminal output should be similar to...
# CPZZG7ZTMYKXUYDFHLPGKL4E6FF5BONF5ZZCD5U6PAJPNUMSSGTQ
# {
#     "creator": "THQHGD4HEESOPSJJYYF34MWKOI57HXBX4XR63EPBKCWPOJG5KUPDJ7QJCM",
#     "total": 1000,
#     "decimals": 0,
#     "defaultfrozen": false,
#     "unitname": "LATINUM",
#     "assetname": "latinum",
#     "url": "https://path/to/my/asset/details",
#     "managerkey": "THQHGD4HEESOPSJJYYF34MWKOI57HXBX4XR63EPBKCWPOJG5KUPDJ7QJCM",
#     "reserveaddr": "AJNNFQN7DSR7QEY766V7JDG35OPM53ZSNF7CU264AWOOUGSZBMLMSKCRIU",
#     "freezeaddr": "AJNNFQN7DSR7QEY766V7JDG35OPM53ZSNF7CU264AWOOUGSZBMLMSKCRIU",
#     "clawbackaddr": "AJNNFQN7DSR7QEY766V7JDG35OPM53ZSNF7CU264AWOOUGSZBMLMSKCRIU"
# }

