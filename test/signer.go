package test

type MockSigner []byte

func (s MockSigner) Sign(message []byte) ([]byte, error) {
	return s, nil
}
