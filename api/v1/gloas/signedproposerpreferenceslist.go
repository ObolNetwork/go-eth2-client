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

import (
	"encoding/json"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/pk910/dynamic-ssz/sszutils"
)

// SignedProposerPreferencesList is the List[SignedProposerPreferences, (MIN_SEED_LOOKAHEAD + 1) * SLOTS_PER_EPOCH]
// submitted to POST /eth/v1/validator/proposer_preferences.
type SignedProposerPreferencesList []*gloas.SignedProposerPreferences

var _ = sszutils.Annotate[SignedProposerPreferencesList](`ssz-max:"64" dynssz-max:"(MIN_SEED_LOOKAHEAD+1)*SLOTS_PER_EPOCH"`)

// String returns a string version of the structure.
func (s *SignedProposerPreferencesList) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
