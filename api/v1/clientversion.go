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

// ClientVersion uniquely identifies a client implementation and its version, mirroring the
// Engine API client version structure.
type ClientVersion struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

type clientVersionJSON struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

// MarshalJSON implements json.Marshaler.
func (c *ClientVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(&clientVersionJSON{
		Code:    c.Code,
		Name:    c.Name,
		Version: c.Version,
		Commit:  c.Commit,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *ClientVersion) UnmarshalJSON(input []byte) error {
	var data clientVersionJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.Code == "" {
		return errors.New("code missing")
	}
	if data.Name == "" {
		return errors.New("name missing")
	}
	if data.Version == "" {
		return errors.New("version missing")
	}
	if data.Commit == "" {
		return errors.New("commit missing")
	}

	c.Code = data.Code
	c.Name = data.Name
	c.Version = data.Version
	c.Commit = data.Commit

	return nil
}

func (c *ClientVersion) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
