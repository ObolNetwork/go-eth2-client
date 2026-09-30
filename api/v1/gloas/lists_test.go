// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gloas_test

import (
	"testing"

	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func testPayloadAttestationMessage(i byte) *gloas.PayloadAttestationMessage {
	return &gloas.PayloadAttestationMessage{
		ValidatorIndex: phase0.ValidatorIndex(i),
		Data: &gloas.PayloadAttestationData{
			BeaconBlockRoot:   phase0.Root{i},
			Slot:              phase0.Slot(i),
			PayloadPresent:    i%2 == 0,
			BlobDataAvailable: true,
		},
		Signature: phase0.BLSSignature{i},
	}
}

func testSignedProposerPreferences(i byte) *gloas.SignedProposerPreferences {
	return &gloas.SignedProposerPreferences{
		Message: &gloas.ProposerPreferences{
			DependentRoot:  phase0.Root{i},
			ProposalSlot:   phase0.Slot(i),
			ValidatorIndex: phase0.ValidatorIndex(i),
			FeeRecipient:   [20]byte{i},
			TargetGasLimit: 60_000_000,
		},
		Signature: phase0.BLSSignature{i},
	}
}

// TestPayloadAttestationMessagesSSZ asserts the List[PayloadAttestationMessage, PTC_SIZE] encoding
// used by POST /eth/v1/beacon/pool/payload_attestations.
func TestPayloadAttestationMessagesSSZ(t *testing.T) {
	list := apiv1gloas.PayloadAttestationMessages{testPayloadAttestationMessage(1), testPayloadAttestationMessage(2)}

	data, err := list.MarshalSSZ()
	require.NoError(t, err)

	// The list of fixed size containers is the concatenation of the elements.
	first, err := list[0].MarshalSSZ()
	require.NoError(t, err)
	require.Equal(t, first, data[:len(first)])
	require.Len(t, data, 2*len(first))

	var decoded apiv1gloas.PayloadAttestationMessages
	require.NoError(t, decoded.UnmarshalSSZ(data))
	require.Equal(t, list, decoded)

	var empty apiv1gloas.PayloadAttestationMessages
	require.NoError(t, empty.UnmarshalSSZ(nil))
	require.Empty(t, empty)

	require.Error(t, decoded.UnmarshalSSZ(data[:len(data)-1]), "misaligned body")

	tooMany := make([]byte, (512+1)*len(first))
	require.Error(t, decoded.UnmarshalSSZ(tooMany), "over PTC_SIZE")
}

// TestSignedProposerPreferencesListSSZ asserts the
// List[SignedProposerPreferences, (MIN_SEED_LOOKAHEAD + 1) * SLOTS_PER_EPOCH] encoding
// used by POST /eth/v1/validator/proposer_preferences.
func TestSignedProposerPreferencesListSSZ(t *testing.T) {
	list := apiv1gloas.SignedProposerPreferencesList{testSignedProposerPreferences(1), testSignedProposerPreferences(2)}

	data, err := list.MarshalSSZ()
	require.NoError(t, err)

	first, err := list[0].MarshalSSZ()
	require.NoError(t, err)
	require.Equal(t, first, data[:len(first)])
	require.Len(t, data, 2*len(first))

	var decoded apiv1gloas.SignedProposerPreferencesList
	require.NoError(t, decoded.UnmarshalSSZ(data))
	require.Equal(t, list, decoded)

	require.Error(t, decoded.UnmarshalSSZ(data[:len(data)-1]), "misaligned body")

	tooMany := make([]byte, (64+1)*len(first))
	require.Error(t, decoded.UnmarshalSSZ(tooMany), "over (MIN_SEED_LOOKAHEAD+1)*SLOTS_PER_EPOCH")
}
