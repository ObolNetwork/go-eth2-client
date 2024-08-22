package deneb

import (
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/holiman/uint256"
)

// GnosisBeaconBlock represents a beacon block.
type GnosisBeaconBlock struct {
	Slot          phase0.Slot
	ProposerIndex phase0.ValidatorIndex
	ParentRoot    phase0.Root `ssz-size:"32"`
	StateRoot     phase0.Root `ssz-size:"32"`
	Body          *GnosisBeaconBlockBody
}

// String returns a string version of the structure.
func (b *GnosisBeaconBlock) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

// GnosisBeaconBlockBody represents the body of a beacon block.
type GnosisBeaconBlockBody struct {
	RANDAOReveal          phase0.BLSSignature `ssz-size:"96"`
	ETH1Data              *phase0.ETH1Data
	Graffiti              [32]byte                      `ssz-size:"32"`
	ProposerSlashings     []*phase0.ProposerSlashing    `ssz-max:"16"`
	AttesterSlashings     []*phase0.AttesterSlashing    `ssz-max:"2"`
	Attestations          []*phase0.Attestation         `ssz-max:"128"`
	Deposits              []*phase0.Deposit             `ssz-max:"16"`
	VoluntaryExits        []*phase0.SignedVoluntaryExit `ssz-max:"16"`
	SyncAggregate         *altair.SyncAggregate
	ExecutionPayload      *GnosisExecutionPayload
	BLSToExecutionChanges []*capella.SignedBLSToExecutionChange `ssz-max:"16"`
	BlobKZGCommitments    []KZGCommitment                       `ssz-max:"4096" ssz-size:"?,48"`
}

// String returns a string version of the structure.
func (b *GnosisBeaconBlockBody) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

// GnosisExecutionPayload represents an execution layer payload.
type GnosisExecutionPayload struct {
	ParentHash    phase0.Hash32              `ssz-size:"32"`
	FeeRecipient  bellatrix.ExecutionAddress `ssz-size:"20"`
	StateRoot     phase0.Root                `ssz-size:"32"`
	ReceiptsRoot  phase0.Root                `ssz-size:"32"`
	LogsBloom     [256]byte                  `ssz-size:"256"`
	PrevRandao    [32]byte                   `ssz-size:"32"`
	BlockNumber   uint64
	GasLimit      uint64
	GasUsed       uint64
	Timestamp     uint64
	ExtraData     []byte                  `ssz-max:"32"`
	BaseFeePerGas *uint256.Int            `ssz-size:"32"`
	BlockHash     phase0.Hash32           `ssz-size:"32"`
	Transactions  []bellatrix.Transaction `ssz-max:"1048576,1073741824" ssz-size:"?,?"`
	Withdrawals   []*capella.Withdrawal   `ssz-max:"8"`
	BlobGasUsed   uint64
	ExcessBlobGas uint64
}

// String returns a string version of the structure.
func (e *GnosisExecutionPayload) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
