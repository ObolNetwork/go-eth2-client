package multi

import (
	"context"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// NodeIdentity provides the identity information of the node.
func (s *Service) NodeIdentity(ctx context.Context, opts *api.NodeIdentityOpts) (*api.Response[*apiv1.NodeIdentity], error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		nodeIdentity, err := client.(consensusclient.NodeIdentityProvider).NodeIdentity(ctx, opts)
		if err != nil {
			return nil, err
		}

		return nodeIdentity, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	response, isResponse := res.(*api.Response[*apiv1.NodeIdentity])
	if !isResponse {
		return nil, ErrIncorrectType
	}

	return response, nil
}
