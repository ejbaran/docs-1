Title: Rekeying

## Overview

Rekeying is a powerful protocol feature which enables an Algorand account holder to maintain a static receiving public address while rotating the spending key(s) by defining an authorized address which may be a single address, MultiSig or LogicSig account. Key management is an important concept to understand and Algorand provide tools to accomplish relevant tasks securely. 

!!! Warning
    This documentation is based on pre-release information and considered neither authoritative nor final. Critiques are welcome to the PR or @ryanRfox on Slack. Testing is _only_ available (to my knowledge) on a private network built from the [maxj/applications](https://github.com/justicz/go-algorand/tree/maxj/applications) repo. Unfortunately, the code samples provided will not (yet) work on BetaNet or TestNet. There are many "TODO:" items reminding me to change them prior to publication.

### Quick Review

The [account overview](https://staging.new-dev-site.algorand.org/docs/features/accounts/#keys-and-addresses0) page introduced _keys_, _addresses_ and _accounts_. Initially, all Algorand accounts are comprised of a public key (pk) and the spending key (sk), which derive the address. The address is commonly displayed within wallet software and remains static for each account. When you receive Algos or other assets, they will be sent to your address. When you send from your account, the transaction must be authorized using the appropriate spending key (sk).  

### Introducing Authorized Address

The account object includes a field _auth-addr_ which, when populated, defines the address authorized (aa) to sign transactions from this account. Initially, the _auth-addr_ field is implicitly set to the account's address. The only valid (aa) is the (sk) created during account generation. To conserve resources, the _auth-addr_ field is only stored and displayed when an authorized `rekey-to` transaction is confirmed. 

#### Initial Account

Use the following code samples to view a brand new account on TestNet:

```bash tab="goal"
# TODO: Change ACCOUNT_ID to match BetaNet, then TestNet later.
# Set the appropriate values for your environment. These are the defaults for Sandbox on TestNet
$ HEADER="x-algo-api-token:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
$ curl "localhost:4001/v2/accounts/2OOTR3IDG57QCHTCFWHMN6YEEQUMGHCT3HIBMPH5A3VQCDZG5GK7DIKDXI" -H $HEADER
```

```bash tab="JavaScript"

``` 

```bash tab="Python"

``` 

```bash tab="Java"

``` 

```bash tab="Go"

``` 

Response:
Notice the account object lacks the _auth-addr_ field. The (sk) for _address_ is the implicit (aa) for an initial account.

```json 
{
    "address": "2OOTR3IDG57QCHTCFWHMN6YEEQUMGHCT3HIBMPH5A3VQCDZG5GK7DIKDXI",
    "amount": 5000020000000000,
    "amount-without-pending-rewards": 5000020000000000,
    "apps-local-state": [],
    "apps-total-schema": {
        "num-byte-slice": 0,
        "num-uint": 0
    },
    "assets": [],
    "created-apps": [],
    "created-assets": [],
    "pending-rewards": 0,
    "reward-base": 0,
    "rewards": 0,
    "round": 1602,
    "status": "Offline"
}
TODO: highlight="2"
```

#### Authorized Account

Next, modify your code slightly to display results for this account `YQP5SHSZOGOUUXIXMBM445HM5N67SCV3XZQ7TXNFOTLC7ZY5QW24DHBDOY`.

Response:
Notice the _auth-addr_ field is populated, which means any transactions from `YQP5SHSZOGOUUXIXMBM445HM5N67SCV3XZQ7TXNFOTLC7ZY5QW24DHBDOY` must now be signed by `2OOTR3IDG57QCHTCFWHMN6YEEQUMGHCT3HIBMPH5A3VQCDZG5GK7DIKDXI` to become authorized. 

```json 
{
    "address": "YQP5SHSZOGOUUXIXMBM445HM5N67SCV3XZQ7TXNFOTLC7ZY5QW24DHBDOY",
    "amount": 5000020000000000,
    "amount-without-pending-rewards": 5000020000000000,
    "apps-local-state": [],
    "apps-total-schema": {
        "num-byte-slice": 0,
        "num-uint": 0
    },
    "assets": [],
    "auth-addr": "2OOTR3IDG57QCHTCFWHMN6YEEQUMGHCT3HIBMPH5A3VQCDZG5GK7DIKDXI",
    "created-apps": [],
    "created-assets": [],
    "pending-rewards": 0,
    "reward-base": 0,
    "rewards": 0,
    "round": 1602,
    "status": "Offline"
}
TODO: highlight="2,11"
```

#### rekey-to Event

The only way to change the (aa) is to authorize a "rekey-to" transaction. Initially, the (aa) is implicitly the (sk), as shown in the first example above. The second example account above completed at least one successful "rekey-to" event, perhaps many, but we only observe the most recent and authoritative result. The "rekey-to" event is comprised of: 

- Constructing a payment transaction defining the _rekey-to_ field
- Signing the transaction with the then current (aa) 
- Confirming the transaction on the network

The result will be the _auth-addr_ field of the account object is defined, modified or removed. Defining or modifying means only the (sk) of the corresponding _auth-addr_ may authorize future transactions for this _address_. Removing the _auth-addr_ is really an explicit redefining of the (aa) back to the (sk) of this _address_ (observed implicitly). 

!!! Info
    Scenarios with code samples will demonstrate key rotation below.

### Potential Use Cases

## Scenario 1

### Send from Initial Account


### Rekey to Authorized Account

#### TEST: Send from Initial Account

### Generate Unsigned Transaction

### Sign Using Authorized Account

### Broadcast Signed Transaction

## Scenario 2

### Change Address with Secret Key to MultiSig

## Scenario 3

### Change Key for MultiSig Address
