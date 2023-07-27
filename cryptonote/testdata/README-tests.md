# address_tests.json

In these 500 entries, `monero-wallet-rpc`'s `create_wallet` method was used to
create the seed sequences with a randomized choice of language. We then added
randomized password to roughly 50% of the entries, after which the seeds +
password were feed back to a wallet using `restore_deterministic_wallet`. The
generated spend/view keys were read back using `query_key` and the first 3
addresses of the first 3 accounts were read using `get_address`.
