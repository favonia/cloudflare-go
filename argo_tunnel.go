package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// ArgoTunnel is the struct definition of a tunnel.
type ArgoTunnel struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Secret      string                 `json:"tunnel_secret,omitempty"`
	CreatedAt   *time.Time             `json:"created_at,omitempty"`
	DeletedAt   *time.Time             `json:"deleted_at,omitempty"`
	Connections []ArgoTunnelConnection `json:"connections,omitempty"`
}

// ArgoTunnelConnection represents the connections associated with a tunnel.
type ArgoTunnelConnection struct {
	ColoName           string `json:"colo_name"`
	UUID               string `json:"uuid"`
	IsPendingReconnect bool   `json:"is_pending_reconnect"`
}

// ArgoTunnelsDetailResponse is used for representing the API response payload for
// multiple tunnels.
type ArgoTunnelsDetailResponse struct {
	Result []ArgoTunnel `json:"result"`
	Response
}

// ArgoTunnelDetailResponse is used for representing the API response payload for
// a single tunnel.
type ArgoTunnelDetailResponse struct {
	Result ArgoTunnel `json:"result"`
	Response
}

// ArgoTunnels lists all tunnels.
//
// API reference: https://api.cloudflare.com/#argo-tunnel-list-argo-tunnels
func (api *API) ArgoTunnels(ctx context.Context, accountID string) ([]ArgoTunnel, error) {
	uri := "/accounts/" + accountID + "/tunnels"

	res, err := api.makeRequestContext(ctx, "GET", uri, nil)
	if err != nil {
		return []ArgoTunnel{}, errors.Wrap(err, errMakeRequestError)
	}

	var argoDetailsResponse ArgoTunnelsDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return []ArgoTunnel{}, errors.Wrap(err, errUnmarshalError)
	}
	return argoDetailsResponse.Result, nil
}

// ArgoTunnel returns a single Argo tunnel.
//
// API reference: https://api.cloudflare.com/#argo-tunnel-get-argo-tunnel
func (api *API) ArgoTunnel(ctx context.Context, accountID, tunnelUUID string) (ArgoTunnel, error) {
	uri := fmt.Sprintf("/accounts/%s/tunnels/%s", accountID, tunnelUUID)

	res, err := api.makeRequestContext(ctx, "GET", uri, nil)
	if err != nil {
		return ArgoTunnel{}, errors.Wrap(err, errMakeRequestError)
	}

	var argoDetailsResponse ArgoTunnelDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return ArgoTunnel{}, errors.Wrap(err, errUnmarshalError)
	}
	return argoDetailsResponse.Result, nil
}

// CreateArgoTunnel creates a new tunnel for the account.
//
// API reference: https://api.cloudflare.com/#argo-tunnel-create-argo-tunnel
func (api *API) CreateArgoTunnel(ctx context.Context, accountID, name, secret string) (ArgoTunnel, error) {
	uri := "/accounts/" + accountID + "/tunnels"

	tunnel := ArgoTunnel{Name: name, Secret: secret}

	res, err := api.makeRequestContext(ctx, "POST", uri, tunnel)
	if err != nil {
		return ArgoTunnel{}, errors.Wrap(err, errMakeRequestError)
	}

	var argoDetailsResponse ArgoTunnelDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return ArgoTunnel{}, errors.Wrap(err, errUnmarshalError)
	}
	return argoDetailsResponse.Result, nil
}

// DeleteArgoTunnel removes a single Argo tunnel.
//
// API reference: https://api.cloudflare.com/#argo-tunnel-delete-argo-tunnel
func (api *API) DeleteArgoTunnel(ctx context.Context, accountID, tunnelUUID string) error {
	uri := fmt.Sprintf("/accounts/%s/tunnels/%s", accountID, tunnelUUID)

	res, err := api.makeRequestContext(ctx, "DELETE", uri, nil)
	if err != nil {
		return errors.Wrap(err, errMakeRequestError)
	}

	var argoDetailsResponse ArgoTunnelDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return errors.Wrap(err, errUnmarshalError)
	}

	return nil
}

// CleanupArgoTunnelConnections deletes any inactive connections on a tunnel.
//
// API reference: https://api.cloudflare.com/#argo-tunnel-clean-up-argo-tunnel-connections
func (api *API) CleanupArgoTunnelConnections(ctx context.Context, accountID, tunnelUUID string) error {
	uri := fmt.Sprintf("/accounts/%s/tunnels/%s/connections", accountID, tunnelUUID)

	res, err := api.makeRequestContext(ctx, "DELETE", uri, nil)
	if err != nil {
		return errors.Wrap(err, errMakeRequestError)
	}

	var argoDetailsResponse ArgoTunnelDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return errors.Wrap(err, errUnmarshalError)
	}

	return nil
}