package electra

import "github.com/attestantio/go-eth2-client/spec/deneb"

// GnosisBeaconBlock

func (b *GnosisBeaconBlock) UnmarshalJSON(bytes []byte) error {
	var ep BeaconBlock
	if err := ep.UnmarshalJSON(bytes); err != nil {
		return err
	}

	*b = BeaconBlockToGnosis(ep)

	return nil
}

func (b *GnosisBeaconBlock) MarshalJSON() ([]byte, error) {
	ep := gnosisBBToStd(*b)

	return ep.MarshalJSON()
}

// MarshalYAML implements yaml.Marshaler.
func (b *GnosisBeaconBlock) MarshalYAML() ([]byte, error) {
	ep := gnosisBBToStd(*b)

	return ep.MarshalYAML()
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *GnosisBeaconBlock) UnmarshalYAML(input []byte) error {
	var ep BeaconBlock
	if err := ep.UnmarshalYAML(input); err != nil {
		return err
	}

	*b = BeaconBlockToGnosis(ep)

	return nil
}

// GnosisBeaconBlockBody

func (b *GnosisBeaconBlockBody) UnmarshalJSON(bytes []byte) error {
	var ep BeaconBlockBody
	if err := ep.UnmarshalJSON(bytes); err != nil {
		return err
	}

	*b = stdBBBToGnosis(ep)

	return nil
}

func (b *GnosisBeaconBlockBody) MarshalJSON() ([]byte, error) {
	ep := gnosisBBBToStd(*b)

	return ep.MarshalJSON()
}

// MarshalYAML implements yaml.Marshaler.
func (b *GnosisBeaconBlockBody) MarshalYAML() ([]byte, error) {
	ep := gnosisBBBToStd(*b)

	return ep.MarshalYAML()
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *GnosisBeaconBlockBody) UnmarshalYAML(input []byte) error {
	var ep BeaconBlockBody
	if err := ep.UnmarshalYAML(input); err != nil {
		return err
	}

	*b = stdBBBToGnosis(ep)

	return nil
}

// GnosisExecutionPayload

func (b *GnosisExecutionPayload) UnmarshalJSON(bytes []byte) error {
	var ep deneb.ExecutionPayload
	if err := ep.UnmarshalJSON(bytes); err != nil {
		return err
	}

	*b = stdEPToGnosis(ep)

	return nil
}

func (b *GnosisExecutionPayload) MarshalJSON() ([]byte, error) {
	ep := gnosisEPToStd(*b)

	return ep.MarshalJSON()
}

// MarshalYAML implements yaml.Marshaler.
func (b *GnosisExecutionPayload) MarshalYAML() ([]byte, error) {
	ep := gnosisEPToStd(*b)

	return ep.MarshalYAML()
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *GnosisExecutionPayload) UnmarshalYAML(input []byte) error {
	var ep deneb.ExecutionPayload
	if err := ep.UnmarshalYAML(input); err != nil {
		return err
	}

	*b = stdEPToGnosis(ep)

	return nil
}

func gnosisEPToStd(ep GnosisExecutionPayload) deneb.ExecutionPayload {
	var b deneb.ExecutionPayload

	b.ParentHash = ep.ParentHash
	b.FeeRecipient = ep.FeeRecipient
	b.StateRoot = ep.StateRoot
	b.ReceiptsRoot = ep.ReceiptsRoot
	b.LogsBloom = ep.LogsBloom
	b.PrevRandao = ep.PrevRandao
	b.BlockNumber = ep.BlockNumber
	b.GasLimit = ep.GasLimit
	b.GasUsed = ep.GasUsed
	b.Timestamp = ep.Timestamp
	b.ExtraData = ep.ExtraData
	b.BaseFeePerGas = ep.BaseFeePerGas
	b.BlockHash = ep.BlockHash
	b.Transactions = ep.Transactions
	b.Withdrawals = ep.Withdrawals
	b.BlobGasUsed = ep.BlobGasUsed
	b.ExcessBlobGas = ep.ExcessBlobGas

	return b
}

func stdEPToGnosis(b deneb.ExecutionPayload) GnosisExecutionPayload {
	var ep GnosisExecutionPayload

	ep.ParentHash = b.ParentHash
	ep.FeeRecipient = b.FeeRecipient
	ep.StateRoot = b.StateRoot
	ep.ReceiptsRoot = b.ReceiptsRoot
	ep.LogsBloom = b.LogsBloom
	ep.PrevRandao = b.PrevRandao
	ep.BlockNumber = b.BlockNumber
	ep.GasLimit = b.GasLimit
	ep.GasUsed = b.GasUsed
	ep.Timestamp = b.Timestamp
	ep.ExtraData = b.ExtraData
	ep.BaseFeePerGas = b.BaseFeePerGas
	ep.BlockHash = b.BlockHash
	ep.Transactions = b.Transactions
	ep.Withdrawals = b.Withdrawals
	ep.BlobGasUsed = b.BlobGasUsed
	ep.ExcessBlobGas = b.ExcessBlobGas

	return ep
}

func gnosisBBBToStd(eb GnosisBeaconBlockBody) BeaconBlockBody {
	ep := gnosisEPToStd(*eb.ExecutionPayload)

	b := BeaconBlockBody{
		RANDAOReveal:          eb.RANDAOReveal,
		ETH1Data:              eb.ETH1Data,
		Graffiti:              eb.Graffiti,
		ProposerSlashings:     eb.ProposerSlashings,
		AttesterSlashings:     eb.AttesterSlashings,
		Attestations:          eb.Attestations,
		Deposits:              eb.Deposits,
		VoluntaryExits:        eb.VoluntaryExits,
		SyncAggregate:         eb.SyncAggregate,
		ExecutionPayload:      &ep,
		BLSToExecutionChanges: eb.BLSToExecutionChanges,
		BlobKZGCommitments:    eb.BlobKZGCommitments,
	}

	return b
}

func stdBBBToGnosis(b BeaconBlockBody) GnosisBeaconBlockBody {
	ep := stdEPToGnosis(*b.ExecutionPayload)

	gb := GnosisBeaconBlockBody{
		RANDAOReveal:          b.RANDAOReveal,
		ETH1Data:              b.ETH1Data,
		Graffiti:              b.Graffiti,
		ProposerSlashings:     b.ProposerSlashings,
		AttesterSlashings:     b.AttesterSlashings,
		Attestations:          b.Attestations,
		Deposits:              b.Deposits,
		VoluntaryExits:        b.VoluntaryExits,
		SyncAggregate:         b.SyncAggregate,
		ExecutionPayload:      &ep,
		BLSToExecutionChanges: b.BLSToExecutionChanges,
		BlobKZGCommitments:    b.BlobKZGCommitments,
	}

	return gb
}

func gnosisBBToStd(gb GnosisBeaconBlock) BeaconBlock {
	body := gnosisBBBToStd(*gb.Body)

	return BeaconBlock{
		Slot:          gb.Slot,
		ProposerIndex: gb.ProposerIndex,
		ParentRoot:    gb.ParentRoot,
		StateRoot:     gb.StateRoot,
		Body:          &body,
	}
}

// BeaconBlockToGnosis returns b's data under a GnosisBeaconBlock.
func BeaconBlockToGnosis(b BeaconBlock) GnosisBeaconBlock {
	body := stdBBBToGnosis(*b.Body)

	return GnosisBeaconBlock{
		Slot:          b.Slot,
		ProposerIndex: b.ProposerIndex,
		ParentRoot:    b.ParentRoot,
		StateRoot:     b.StateRoot,
		Body:          &body,
	}
}
