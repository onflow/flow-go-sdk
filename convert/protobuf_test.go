package convert

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dapperlabs/flow-go-sdk/utils/unittest"
)

func TestBlockHeader(t *testing.T) {
	RegisterTestingT(t)

	blockHeaderA := unittest.BlockHeaderFixture()

	message := BlockHeaderToMessage(blockHeaderA)
	blockHeaderB := MessageToBlockHeader(message)

	Expect(blockHeaderA).To(Equal(blockHeaderB))
}

func TestAccountSignature(t *testing.T) {
	RegisterTestingT(t)

	sigA := unittest.AccountSignatureFixture()

	message := AccountSignatureToMessage(sigA)
	sigB := MessageToAccountSignature(message)

	Expect(sigA).To(Equal(sigB))
}

func TestTransaction(t *testing.T) {
	RegisterTestingT(t)

	txA := unittest.TransactionFixture()

	message := TransactionToMessage(txA)

	txB, err := MessageToTransaction(message)
	Expect(err).ToNot(HaveOccurred())

	Expect(txA).To(Equal(txB))
}

func TestAccount(t *testing.T) {
	RegisterTestingT(t)

	accA := unittest.AccountFixture()

	message, err := AccountToMessage(accA)
	Expect(err).ToNot(HaveOccurred())

	accB, err := MessageToAccount(message)
	Expect(err).ToNot(HaveOccurred())

	Expect(accA).To(Equal(accB))
}

func TestAccountKey(t *testing.T) {
	RegisterTestingT(t)

	keyA := unittest.AccountPublicKeyFixture()

	message, err := AccountPublicKeyToMessage(keyA)
	Expect(err).ToNot(HaveOccurred())

	keyB, err := MessageToAccountPublicKey(message)
	Expect(err).ToNot(HaveOccurred())

	Expect(keyA).To(Equal(keyB))
}
