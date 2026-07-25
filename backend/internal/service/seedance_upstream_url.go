package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/util/urlvalidator"
)

// SeedanceMediaEndpoint identifies Binghuo-compatible video media routes.
type SeedanceMediaEndpoint string

const (
	SeedanceMediaEndpointVideoGenerations SeedanceMediaEndpoint = "video_generations"
	SeedanceMediaEndpointVideoStatus      SeedanceMediaEndpoint = "video_status"
	SeedanceMediaEndpointAssetsUploads    SeedanceMediaEndpoint = "assets_uploads"
)

func (e SeedanceMediaEndpoint) RequiresRequestBody() bool {
	return e != SeedanceMediaEndpointVideoStatus
}

func (e SeedanceMediaEndpoint) IsGenerationRequest() bool {
	return e == SeedanceMediaEndpointVideoGenerations
}

func (e SeedanceMediaEndpoint) IsLookupRequest() bool {
	return e == SeedanceMediaEndpointVideoStatus
}

func (e SeedanceMediaEndpoint) IsUploadRequest() bool {
	return e == SeedanceMediaEndpointAssetsUploads
}

func seedanceBaseURLValidator(cfg *config.Config) func(string) (string, error) {
	if cfg == nil {
		return func(raw string) (string, error) {
			return urlvalidator.ValidateURLFormat(raw, true)
		}
	}
	if !cfg.Security.URLAllowlist.Enabled {
		return func(raw string) (string, error) {
			return urlvalidator.ValidateURLFormat(raw, cfg.Security.URLAllowlist.AllowInsecureHTTP)
		}
	}
	return func(raw string) (string, error) {
		return urlvalidator.ValidateHTTPSURL(raw, urlvalidator.ValidationOptions{
			AllowedHosts:     cfg.Security.URLAllowlist.UpstreamHosts,
			RequireAllowlist: true,
			AllowPrivate:     cfg.Security.URLAllowlist.AllowPrivateHosts,
		})
	}
}

func validateSeedanceBaseURL(account *Account, cfg *config.Config) (string, error) {
	if account == nil || !account.IsSeedance() {
		return "", fmt.Errorf("seedance account is required")
	}
	if account.Type != AccountTypeAPIKey {
		return "", fmt.Errorf("seedance accounts must use type=apikey")
	}
	raw := account.GetSeedanceBaseURL()
	if raw == "" {
		return "", fmt.Errorf("seedance base_url is empty")
	}
	validated, err := seedanceBaseURLValidator(cfg)(raw)
	if err != nil {
		return "", errors.New("base URL rejected by URL security policy")
	}
	return strings.TrimRight(validated, "/"), nil
}

func buildSeedanceMediaURL(account *Account, cfg *config.Config, endpoint SeedanceMediaEndpoint, taskID string) (string, error) {
	baseURL, err := validateSeedanceBaseURL(account, cfg)
	if err != nil {
		return "", err
	}
	switch endpoint {
	case SeedanceMediaEndpointVideoGenerations:
		return baseURL + "/video/generations", nil
	case SeedanceMediaEndpointVideoStatus:
		taskID = strings.TrimSpace(taskID)
		if taskID == "" {
			return "", fmt.Errorf("task_id is required")
		}
		return baseURL + "/video/generations/" + url.PathEscape(taskID), nil
	case SeedanceMediaEndpointAssetsUploads:
		return baseURL + "/assets/uploads", nil
	default:
		return "", fmt.Errorf("unsupported seedance media endpoint: %s", endpoint)
	}
}
