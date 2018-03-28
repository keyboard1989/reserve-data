package signer

type FileSigner struct {
	KeystoreD   string `json:"keystore_deposit_path"`
	PassphraseD string `json:"passphrase_deposit"`
	KeystoreI   string `json:"keystore_intermediator_path"`
	PassphraseI string `json:"passphrase_intermediate_account"`
}
