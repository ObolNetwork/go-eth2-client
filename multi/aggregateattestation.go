// Copyright © 2021 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// AggregateAttestation fetches the aggregate attestation given an attestation to v1 beacon node endpoint.
func (s *Service) AggregateAttestation(ctx context.Context,
	opts *api.AggregateAttestationOpts,
) (
	*api.Response[*phase0.Attestation],
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		aggregate, err := client.(consensusclient.AggregateAttestationProvider).AggregateAttestation(ctx, opts)
		if err != nil {
			return nil, err
		}

		return aggregate, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	response, isResponse := res.(*api.Response[*phase0.Attestation])
	if !isResponse {
		return nil, ErrIncorrectType
	}

	return response, nil
}

// AggregateAttestationV2 fetches the aggregate attestation for the given options to v2 beacon node endpoint.
func (s *Service) AggregateAttestationV2(ctx context.Context,
	opts *api.AggregateAttestationOpts,
) (
	*api.Response[*spec.VersionedAttestation],
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		aggregate, err := client.(consensusclient.AggregateAttestationProvider).AggregateAttestationV2(ctx, opts)
		if err != nil {
			return nil, err
		}

		return aggregate, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	response, isResponse := res.(*api.Response[*spec.VersionedAttestation])
	if !isResponse {
		return nil, ErrIncorrectType
	}

	return response, nil
}
