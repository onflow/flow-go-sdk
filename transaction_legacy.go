/*
 * Flow Go SDK
 *
 * Copyright Flow Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package flow

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/onflow/go-ethereum/rlp"
)

type transactionLegacyCanonicalForm struct {
	Payload            payloadCanonicalForm
	PayloadSignatures  []transactionSignatureLegacyCanonicalForm
	EnvelopeSignatures []transactionSignatureLegacyCanonicalForm
}

type transactionSignatureLegacyCanonicalForm struct {
	SignerIndex uint
	KeyIndex    uint32
	Signature   []byte
}

var _ transactionSignatureCommonForm = (*transactionSignatureLegacyCanonicalForm)(nil)

func (s transactionSignatureLegacyCanonicalForm) dummy() {}

func (s *transactionLegacyCanonicalForm) convertToCanonicalForm() *transactionCanonicalForm {
	canonicalPayloadSigs := make([]transactionSignatureCanonicalForm, 0, len(s.PayloadSignatures))
	for _, sig := range s.PayloadSignatures {
		canonicalPayloadSigs = append(canonicalPayloadSigs, transactionSignatureCanonicalForm{
			SignerIndex: sig.SignerIndex,
			KeyIndex:    sig.KeyIndex,
			Signature:   sig.Signature,
		})
	}

	canonicalEnvelopSigs := make([]transactionSignatureCanonicalForm, 0, len(s.EnvelopeSignatures))
	for _, sig := range s.EnvelopeSignatures {
		canonicalEnvelopSigs = append(canonicalEnvelopSigs, transactionSignatureCanonicalForm{
			SignerIndex: sig.SignerIndex,
			KeyIndex:    sig.KeyIndex,
			Signature:   sig.Signature,
		})
	}

	return &transactionCanonicalForm{
		Payload:            s.Payload,
		PayloadSignatures:  canonicalPayloadSigs,
		EnvelopeSignatures: canonicalEnvelopSigs,
	}
}

func decodeTransactionLegacy(transactionMessage []byte) (*transactionCanonicalForm, error) {
	s := rlp.NewStream(bytes.NewReader(transactionMessage), 0)
	temp := &transactionLegacyCanonicalForm{}

	kind, _, err := s.Kind()
	if err != nil {
		return nil, err
	}

	// First kind should always be a list
	if kind != rlp.List {
		return nil, errors.New("unexpected rlp decoding type")
	}

	_, err = s.List()
	if err != nil {
		return nil, err
	}

	// Need to look at the type of the first element to determine if how we're going to be decoding
	kind, _, err = s.Kind()
	if err != nil {
		return nil, err
	}
	// If first kind is not list, safe to assume this is actually just encoded payload, and decrypt as such
	if kind != rlp.List {
		s.Reset(bytes.NewReader(transactionMessage), 0)
		txPayload := payloadCanonicalForm{}
		err := s.Decode(&txPayload)
		if err != nil {
			return nil, err
		}
		temp.Payload = txPayload
		return temp.convertToCanonicalForm(), nil
	}

	// If we're here, we will assume that we're decoding either a envelopeCanonicalForm
	// or a full transactionCanonicalForm

	// Decode the payload
	txPayload := payloadCanonicalForm{}
	err = s.Decode(&txPayload)
	if err != nil {
		return nil, err
	}
	temp.Payload = txPayload

	// Decode the payload sigs
	payloadSigs := []transactionSignatureLegacyCanonicalForm{}
	fmt.Println(s.Kind())
	err = s.Decode(&payloadSigs)
	if err != nil {
		return nil, err
	}
	temp.PayloadSignatures = payloadSigs

	// It's possible for the envelope signature to not exist (e.g. envelopeCanonicalForm).
	kind, _, err = s.Kind()
	if errors.Is(err, rlp.EOL) {
		return temp.convertToCanonicalForm(), nil
	} else if err != nil {
		return nil, err
	}
	// If we're not at EOL, and no error, finish decoding
	envelopeSigs := []transactionSignatureLegacyCanonicalForm{}
	err = s.Decode(&envelopeSigs)
	if err != nil {
		return nil, err
	}
	temp.EnvelopeSignatures = envelopeSigs

	return temp.convertToCanonicalForm(), nil
}
