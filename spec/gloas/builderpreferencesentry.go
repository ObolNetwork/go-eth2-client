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

package gloas

import "github.com/attestantio/go-eth2-client/spec/phase0"

// BuilderPreferencesEntry is one builder preferences submission: a proposer's
// authorization and payment cap for one builder, submitted ahead of the bid request.
type BuilderPreferencesEntry struct {
	ProposerPubkey      phase0.BLSPubKey          `ssz-index:"0" ssz-size:"48"`
	URL                 []byte                    `ssz-index:"1" ssz-max:"2048"`
	Auth                *SignedBuilderRequestAuth `ssz-index:"2"`
	MaxExecutionPayment phase0.Gwei               `ssz-index:"3"`
}
