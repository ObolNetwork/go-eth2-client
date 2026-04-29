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

package mock

import (
	"context"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// NodeVersionV2 returns structured version information for the beacon node and, when available,
// its attached execution client.
func (s *Service) NodeVersionV2(ctx context.Context,
	opts *api.NodeVersionV2Opts,
) (
	*api.Response[*apiv1.NodeVersionV2],
	error,
) {
	if s.NodeVersionV2Func != nil {
		return s.NodeVersionV2Func(ctx, opts)
	}

	return &api.Response[*apiv1.NodeVersionV2]{
		Data: &apiv1.NodeVersionV2{
			BeaconNode: &apiv1.ClientVersion{
				Code:    "MK",
				Name:    "mock",
				Version: s.nodeVersion,
				Commit:  "0x00000000",
			},
		},
		Metadata: make(map[string]any),
	}, nil
}
