// Reference Update Transaction (RUT)
package transaction

import (
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/spx/gitchain/repository"
	"github.com/spx/gitchain/types"
)

func init() {
	gob.Register(&ReferenceUpdate{})
}

const (
	REFERENCE_UPDATE_VERSION = 1
)

type ReferenceUpdate struct {
	Version    uint32
	Repository string
	Ref        string
	Old        repository.Ref
	New        repository.Ref
}

func (tx *ReferenceUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Type":       "Reference Update Tranasction",
		"Version":    tx.Version,
		"Repository": tx.Repository,
		"Ref":        tx.Ref,
		"Old":        hex.EncodeToString(tx.Old),
		"New":        hex.EncodeToString(tx.New),
	})
}

func NewReferenceUpdate(repository, ref string, old, new repository.Ref) *ReferenceUpdate {
	return &ReferenceUpdate{
		Version:    REFERENCE_UPDATE_VERSION,
		Repository: repository,
		Ref:        ref,
		Old:        old,
		New:        new}
}

func (txn *ReferenceUpdate) Valid() bool {
	return (txn.Version == REFERENCE_UPDATE_VERSION && len(txn.Repository) > 0 &&
		len(txn.Ref) > 0)
}

func (txn *ReferenceUpdate) Encode() ([]byte, error) {
	return encode(txn)
}

func (txn *ReferenceUpdate) Hash() types.Hash {
	return hash(txn)
}

func (txn *ReferenceUpdate) String() string {
	return fmt.Sprintf("RUT %s %s %s:%s", txn.Repository, txn.Ref, txn.Old, txn.New)
}
