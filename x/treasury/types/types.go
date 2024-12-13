package types

const (
	EventNewKeyRequest                  = "new_key_request"
	EventKeyRequestFulfilled            = "key_request_fulfilled"
	EventKeyRequestRejected             = "key_request_rejected"
	EventNewSignRequest                 = "new_sign_request"
	EventSignRequestFulfilled           = "sign_request_fulfilled"
	EventSignRequestRejected            = "sign_request_rejected"
	EventNewICATransactionRequest       = "new_ica_transaction_request"
	EventICATransactionRequestFulfilled = "ica_transaction_request_fulfilled"
)

const (
	AttributeRequestId = "request_id"
)

type VerificationStatus int

const (
	Verification_NotVerified VerificationStatus = iota
	Verification_Failed
	Verification_Suceeded
)

var ValidKeyTypes = []string{"ecdsa", "ed25519", "eddsa", "bitcoin", "btc"}

const KeyringCollectorName = "keyring_collector"
