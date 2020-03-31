package crypto

type Signable interface {
	Message() []byte
}

type Signer interface {
	Sign(obj Signable) ([]byte, error)
}
