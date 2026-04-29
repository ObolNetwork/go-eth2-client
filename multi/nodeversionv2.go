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

package multi

import (
	"context"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// NodeVersionV2 provides structured version information for the beacon node and, when available,
// its attached execution client.
func (s *Service) NodeVersionV2(ctx context.Context,
	opts *api.NodeVersionV2Opts,
) (
	*api.Response[*apiv1.NodeVersionV2],
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		aggregate, err := client.(consensusclient.NodeVersionV2Provider).NodeVersionV2(ctx, opts)
		if err != nil {
			return nil, err
		}

		return aggregate, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	response, isResponse := res.(*api.Response[*apiv1.NodeVersionV2])
	if !isResponse {
		return nil, ErrIncorrectType
	}

	return response, nil
}
