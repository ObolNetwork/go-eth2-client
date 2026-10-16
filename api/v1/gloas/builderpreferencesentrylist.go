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
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/pk910/dynamic-ssz/sszutils"
)

// BuilderPreferencesEntryList is the
// List[BuilderPreferencesEntry, MAX_BUILDER_ENTRIES * (MIN_SEED_LOOKAHEAD + 1) * SLOTS_PER_EPOCH]
// submitted to POST /eth/v1/validator/builder_preferences.
type BuilderPreferencesEntryList []*gloas.BuilderPreferencesEntry

var _ = sszutils.Annotate[BuilderPreferencesEntryList](`ssz-max:"4096" dynssz-max:"64*(MIN_SEED_LOOKAHEAD+1)*SLOTS_PER_EPOCH"`)

// The list is not generated: dynamic-ssz v1.3.2 generates code that does not
// compile for a list of variable-size containers, so these encode via reflection
// until dynamic-ssz is bumped to v1.4.0, which generates it correctly.

// MarshalSSZ marshals the list to SSZ-encoded bytes.
func (b *BuilderPreferencesEntryList) MarshalSSZ() ([]byte, error) {
	return dynssz.GetGlobalDynSsz().MarshalSSZ(b)
}

// UnmarshalSSZ unmarshals the list from SSZ-encoded bytes.
func (b *BuilderPreferencesEntryList) UnmarshalSSZ(buf []byte) error {
	return dynssz.GetGlobalDynSsz().UnmarshalSSZ(b, buf)
}

// String returns a string version of the structure.
func (b *BuilderPreferencesEntryList) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
