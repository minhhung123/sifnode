# Design document: https://docs.google.com/document/d/1yxxQ3RtftCvCJp_vDlSR5MiXvtOEJxlGXg3GLvf9dls/edit?skip_itp2_check=true

import json
import sys

chains = {
    "akash": {"binary": "akash", "relayer": "ibc"},
    "iris": {"binary": "iris", "relayer": "hermes"},
    "sentinel": {"binary": "sentinelhub", "relayer": "ibc"},
    "persistence": {"binary": "persistenceCore", "relayer": "hermes"},
    "sifchain": {"binary": "sifnoded", "relayer": "ibc"},
}

# TODO Define format and usage
configs = {
    "local": {
        "akash": "http://127.0.0.1:26656",
        "sifchain": "http://127.0.0.1:26657",
        "chainId": "akash-testnet-6",
        "fees": 5000,
        "gas": 1000,
        "denom": "uakt",
    },
    "ci": {
    },
}

# Runs the command synchronously, checks that the exit code is 0 and returns standard output and error.
def run_command(args, stdin=None, cwd=None, env=None, ignore_errors=False):
    return None  # TODO

# Starts a process and returns a Popen object for it
def start_process(args):
    return None  # TODO

def get_binary_for_chain(chain_name):
    return chains[chain_name]["binary"]

def get_config(config_name):
    return None  # TODO

# Generates a sifnoded key and stores it into test keyring. Returns the mnemonic that can be used to
# recreate it.
def add_new_key_to_keyring(chain, key_name):
    binary = get_binary_for_chain(chain)
    res = run_command([binary, "keys", "add", key_name, "--keyring-backend", "test", "--output", "json"], stdin=["y"])
    return json.loads(res.stdout)["mnemonic"]

def add_existing_key_to_keyring(chain, key_name, mnemonic):
    binary = get_binary_for_chain(chain)
    run_command([binary, "keys", "delete", key_name, "--keyring-backend", "test", "-y"], ignore_errors=True)
    run_command([binary, "keys", "add", key_name, "-i", "--recover", "--keyring-backend", "test"],
        stdin=[mnemonic, ""])

def start_chain(chain):
    pass

# Can be initialized either manually or from genesis file
def init_chain(chain):
    pass

def start_relayer(chain_a, chain_b, channel_id, counterchannel_id):
    # TODO Determine which relayer to use
    relayer_binary = None
    relayer_args = []
    relayer_process = start_process([relayer_binary] + relayer_args)
    return relayer_process

def send_transaction(chain, channel, amount, denom, src_addr, dst_addr, sequence, chain_id, node, broadcast_mode,
    fees, gas, account_number, dry_run=False
):
    if not broadcast_mode in ["async", "block"]:
        raise ValueError("Invalid broadcast_mode '{}'".format(broadcast_mode))
    args = [get_binary_for_chain(chain), "tx", "ibc-transfer", "transfer", "transfer", f"channel-{channel}",
        dst_addr, f"{amount}{denom}", "--from", src_addr, "--keyring-backend", "test", "--chain-id", chain_id,
        "--node", node, "--sequence", str(sequence), "--account-number", account_number] + \
        (["--fees", f"{fees}{denom}"] if fees else []) + \
        (["--gas", gas] if gas else []) + \
        (["--broadcast-mode", broadcast_mode if broadcast_mode else "async"]) + \
        (["--offline"] if broadcast_mode == "async" else []) + \
        (["--dry-run"] if dry_run else [])
    run_command(args)

def query_bank_balance(chain, addr, denom):
    node = None  # TODO
    chain_id = None  # TODO
    result = json.loads(run_command([get_binary_for_chain(chain), "q", "bank", "balances", addr, "--node", node,
        "--chain-id", chain_id, "--output", "json"]).stdout)
    return result[denom]

def run_tests_for_one_chain_in_one_direction(config, other_chain, direction_flag, number_of_iterations):
    from_chain = "sifchain" if direction_flag else other_chain
    to_chain = other_chain if direction_flag else "sifchain"
    sequence = 0   # TODO
    broadcast_mode = "block"  # TODO
    account_number = None  # TODO
    chain_id = int(config["chain)id"])
    denom = config["denom"]
    channel_id = int(config["channel_id"])
    counterchannel_id = int(config["counterchannel_id"])
    from_account = config["from_account"]
    to_account = config["to_account"]
    amount = int(config["amount"])
    node = config[chain_id]["node"]
    fees = config["fees"]
    gas = config["gas"]
    sifchain_proc = start_chain("sifchain")
    other_chain_proc = start_chain(other_chain)
    init_chain(from_chain)
    init_chain(to_chain)
    relayer_proc = start_relayer(from_chain, to_chain, channel_id, counterchannel_id)
    mnemonic = add_new_key_to_keyring("sifchain", from_account)
    add_existing_key_to_keyring(other_chain, to_account, mnemonic)
    from_balance_before = query_bank_balance(from_chain, from_account, denom)
    to_balance_before = query_bank_balance(to_chain, to_account, denom)
    for i in range(number_of_iterations):
        send_transaction(from_chain, channel_id if direction_flag else counterchannel_id, amount, denom, from_account,
            to_account, sequence + i, chain_id, node, broadcast_mode, fees, gas, account_number)
    # TODO Wait for transaction to complete (if async)
    from_balance_after = query_bank_balance(from_chain, from_account, denom)
    to_balance_after = query_bank_balance(to_chain, to_account, denom)
    relayer_proc.stop()
    sifchain_proc.stop()
    other_chain_proc.stop()
    assert from_balance_after == from_balance_before - number_of_iterations * amount
    assert to_balance_after == to_balance_before + number_of_iterations * amount

def run_tests_for_all_chains_in_both_directions(config, number_of_iterations):
    for chain in chains:
        run_tests_for_one_chain_in_one_direction(config, chain, True, number_of_iterations)
        run_tests_for_one_chain_in_one_direction(config, chain, False, number_of_iterations)

# This is called from GitHub CI/CD (i.e. .github/workflows)
def run_from_ci(args):
    config = get_config("local")
    run_tests_for_all_chains_in_both_directions(config, 1000)

def run_locally(args):
    config = get_config("local")
    other_chain = args[0]
    direction_flag = args[1] == "receiver"
    number_of_iterations = int(args[2])
    run_tests_for_one_chain_in_one_direction(config, other_chain, direction_flag, number_of_iterations)

def main(argv):
    action = argv[0]
    action_args = argv[1:]
    if action == "ci":
        run_from_ci(action_args)
    elif action == "local":
        run_locally(action_args)

if __name__ == "__main__":
    main(sys.argv)
