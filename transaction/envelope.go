package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/spx/gitchain/keys"
	"github.com/spx/gitchain/types"
	"github.com/spx/gitchain/util"
)

type Envelope struct {
	PreviousEnvelopeHash types.Hash
	SignatureR           []byte
	SignatureS           []byte
	PublicKey            []byte
	NextPublicKey        []byte
	Transaction          T
}

func NewEnvelope(prev types.Hash, txn T, args ...[]byte) *Envelope {
	e := &Envelope{
		PreviousEnvelopeHash: prev,
		Transaction:          txn}
	if len(args) == 1 {
		e.PublicKey = args[0]
		e.NextPublicKey = args[0]
	}
	return e
}

func (e *Envelope) Hash() types.Hash {
	return util.SHA256(append(append(e.Transaction.Hash(), e.PreviousEnvelopeHash...), e.NextPublicKey...))
}

func (e *Envelope) Sign(privateKey *ecdsa.PrivateKey) error {
	pubkey, err := keys.EncodeECDSAPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	if len(e.NextPublicKey) == 0 || bytes.Compare(e.NextPublicKey, e.PublicKey) == 0 {
		e.NextPublicKey = pubkey
	}
	e.PublicKey = pubkey

	sigR, sigS, err := ecdsa.Sign(rand.Reader, privateKey, e.Hash())
	if err != nil {
		return err
	}
	e.SignatureR = sigR.Bytes()
	e.SignatureS = sigS.Bytes()

	return nil
}

func (e *Envelope) Verify() (bool, error) {
	publicKey, err := keys.DecodeECDSAPublicKey(e.PublicKey)
	if err != nil {
		return false, err
	}
	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(e.SignatureR)
	s.SetBytes(e.SignatureS)
	return ecdsa.Verify(publicKey, e.Hash(), r, s), nil
}

func (e *Envelope) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&e)
	return buf.Bytes(), err
}

func DecodeEnvelope(b []byte) (*Envelope, error) {
	var t *Envelope
	buf := bytes.NewBuffer(b)
	enc := gob.NewDecoder(buf)
	err := enc.Decode(&t)
	return t, err
}

func (e *Envelope) String() string {
	return fmt.Sprintf("[%s %s]", e.Hash(), e.Transaction)
}
