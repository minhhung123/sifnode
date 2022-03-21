import siftool_path
from siftool import eth, test_utils, sifchain
from siftool.common import *


fund_amount_eth = 10 * eth.ETH
fund_amount_sif = 10 * test_utils.sifnode_funds_for_transfer_peggy1  # TODO How much rowan do we need? (this is 10**18)


def test_eth_to_ceth_and_back_grpc(ctx):
    # Create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])

    # Verify initial balance
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    # Number of burn transactions that we want to do
    count = 1  # TODO Multiply the total amount by count

    # Send from ethereum to sifchain by locking
    amount_to_send_in_tx = 123456 * eth.GWEI
    total_amount_to_send = amount_to_send_in_tx * count
    assert total_amount_to_send < fund_amount_eth

    ctx.bridge_bank_lock_eth(test_eth_account, test_sif_account, total_amount_to_send)
    ctx.advance_blocks()

    # Verify final balance
    test_sif_account_final_balance = ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance)
    balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)
    assert exactly_one(list(balance_diff.keys())) == ctx.ceth_symbol
    assert balance_diff[ctx.ceth_symbol] == total_amount_to_send

    # Send from sifchain to ethereum by burning on sifchain side,
    # > sifnoded tx ethbridge burn
    # Reduce amount for cross-chain fee. The same formula is used inside this function.
    eth_balance_before = ctx.eth.get_eth_balance(test_eth_account)
    amount_to_send = amount_to_send_in_tx - ctx.eth.cross_chain_fee_base * ctx.eth.cross_chain_burn_fee

    log.debug("Generating {} transactions...".format(count))
    signed_encoded_txs = []
    for i in range(1):
        # "generate_only" tells sifnode to print a transaction as JSON instead of signing and broadcasting it
        tx = ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, amount_to_send,
            ctx.ceth_symbol, generate_only=True)
        signed_tx = ctx.sifnode_client.sign_transaction(tx, test_sif_account)
        encoded_tx = ctx.sifnode_client.encode_transaction(signed_tx)
        signed_encoded_txs.append(encoded_tx)

    log.debug("Broadcasting {} transactions...".format(count))
    for tx in signed_encoded_txs:
        result = ctx.sifnode_client.broadcast_tx(tx)

    # Verify final balance
    new_eth_balance = ctx.wait_for_eth_balance_change(test_eth_account, eth_balance_before)
