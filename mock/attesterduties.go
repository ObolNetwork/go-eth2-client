// Copyright © 2020 Attestant Limited.
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

// AttesterDuties obtains attester duties.
// If validatorIndices is nil it will return all duties for the given epoch.
func (s *Service) AttesterDuties(ctx context.Context,
	opts *api.AttesterDutiesOpts,
) (
	*api.Response[[]*apiv1.AttesterDuty],
	error,
) {
	if s.AttesterDutiesFunc != nil {
		return s.AttesterDutiesFunc(ctx, opts)
	}

	data := make([]*apiv1.AttesterDuty, len(opts.Indices))
	for i := range opts.Indices {
		data[i] = &apiv1.AttesterDuty{
			ValidatorIndex: opts.Indices[i],
		}
	}

	return &api.Response[[]*apiv1.AttesterDuty]{
		Data:     data,
		Metadata: make(map[string]any),
	}, nil
}
