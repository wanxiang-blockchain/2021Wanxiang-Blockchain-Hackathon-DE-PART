package main

type IPatientClient interface {
	GenReEncryptionKey (bPubKeyInput string)  (string, string, error)
}
