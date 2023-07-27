# test_all_seeds_all_langs.json

This test set was created by using each 24-word increment from the 1626 seed
word lists for all languages. The final set starts at index 1626 instead of
1608, as 1626 is not evenly divisible by 24.

monero-wallet-rpc was used with restore_deterministic_wallet and passed the 24
seeds which returns the full 25-word seed list in the response that was saved
into the tests. The code creating the test set verified that monero-wallet-rpc
returned the same secret_key (using query_key method) for every language over
the same indices.

# test_seeds_with_passwords.json

These tests were created with monero-wallet-rpc's create_wallet method,
passing a randomly generated password. The mnemonic seeds and spend key
values were read back using query_key.
