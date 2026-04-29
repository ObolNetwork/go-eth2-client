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

package v1

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// NodeVersionV2 contains structured version information for the beacon node and, when available,
// its attached execution client.
type NodeVersionV2 struct {
	BeaconNode      *ClientVersion `json:"beacon_node"`
	ExecutionClient *ClientVersion `json:"execution_client,omitempty"`
}

type nodeVersionV2JSON struct {
	BeaconNode      *ClientVersion `json:"beacon_node"`
	ExecutionClient *ClientVersion `json:"execution_client,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (n *NodeVersionV2) MarshalJSON() ([]byte, error) {
	return json.Marshal(&nodeVersionV2JSON{
		BeaconNode:      n.BeaconNode,
		ExecutionClient: n.ExecutionClient,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *NodeVersionV2) UnmarshalJSON(input []byte) error {
	var data nodeVersionV2JSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.BeaconNode == nil {
		return errors.New("beacon_node missing")
	}

	n.BeaconNode = data.BeaconNode
	n.ExecutionClient = data.ExecutionClient

	return nil
}

func (n *NodeVersionV2) String() string {
	data, err := json.Marshal(n)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
