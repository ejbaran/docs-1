def retrieve_from_file(path):
    """
    Retrieve signed or unsigned transactions from a file.

    Args:
        path (str): file to read from

    Returns:
        Transaction[], SignedTransaction[], or MultisigTransaction[]:\
            can be a mix of the three
    """

    f = open(path, "rb")
    txns = []
    unp = msgpack.Unpacker(f, raw=False)
    for tx in unp:
        txn=tx["txn"]
        if "msig" in txn:
            txns.append(MultisigTransaction.undictify(txn))
        elif "sig" in txn:
            txns.append(SignedTransaction.undictify(txn))
        elif "lsig" in txn:
            txns.append(LogicSigTransaction.undictify(txn))
        elif "type" in txn:
            txns.append(Transaction.undictify(txn))
    f.close()
    return txns