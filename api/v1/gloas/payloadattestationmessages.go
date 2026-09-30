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

// PayloadAttestationMessages is the List[PayloadAttestationMessage, PTC_SIZE] submitted to
// POST /eth/v1/beacon/pool/payload_attestations.
type PayloadAttestationMessages []*gloas.PayloadAttestationMessage

var _ = sszutils.Annotate[PayloadAttestationMessages](`ssz-max:"512" dynssz-max:"PTC_SIZE"`)

// String returns a string version of the structure.
func (p *PayloadAttestationMessages) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
