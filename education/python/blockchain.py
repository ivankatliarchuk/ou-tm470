class Blockchain(object):
    def __init__(self):
        self.chain = []
        self.current_transactions = []

    def new_block(self):
        # Creates a new block and adds it to the Blockchain
        pass

    def new_transaction(self, sender, recipient, amount):
        """
        Adds a new transaction to the list of new current_transactions to go into the next Blockchain

        :param sender: <str> who sent the amount
        :param recipient: <str> who is recieving
        :param amount: howmuch is sent
        :return: <int> the block id this transaction will enter
        """

        self.current_transactions.append({
        'sender' : sender
        'recipient' : recipient
        'amount' : amount
        })

        return self.last_block['index'] + 1

    @staticmethod
    def hash(block):
        # Hash a Blockchain
        pass

    @property
    def last_block(self):
        # Return the last block on the Blockchain
        pass
