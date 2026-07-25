package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// SeedanceMediaRequestInfo captures billing-relevant fields from a Binghuo video request.
type SeedanceMediaRequestInfo struct {
	Model              string
	Prompt             string
	DurationSeconds    int
	Resolution         string
	Ratio              string
	ImageURLs          []string
	StartFrameURLs     []string
	EndFrameURLs       []string
	ReferenceVideos    []string
	ReferenceAudios    []string
	VideoURL           string // 去字幕
	GenerateAudio      *bool
	HasReferenceVideo  bool
	HasReferenceAudio  bool
	HasReferenceImage  bool
}

func ParseSeedanceMediaRequest(body []byte) SeedanceMediaRequestInfo {
	info := SeedanceMediaRequestInfo{}
	if !gjson.ValidBytes(body) {
		return info
	}
	info.Model = firstNonEmptySeedanceString(
		strings.TrimSpace(gjson.GetBytes(body, "model").String()),
		strings.TrimSpace(gjson.GetBytes(body, "model_id").String()),
	)
	info.Prompt = firstNonEmptySeedanceString(
		strings.TrimSpace(gjson.GetBytes(body, "prompt").String()),
		strings.TrimSpace(gjson.GetBytes(body, "text").String()),
	)
	if duration := gjson.GetBytes(body, "duration"); duration.Exists() && duration.Type == gjson.Number {
		info.DurationSeconds = int(duration.Int())
	} else if seconds := gjson.GetBytes(body, "seconds"); seconds.Exists() && seconds.Type == gjson.Number {
		info.DurationSeconds = int(seconds.Int())
	}
	info.Resolution = strings.TrimSpace(gjson.GetBytes(body, "resolution").String())
	info.Ratio = firstNonEmptySeedanceString(
		strings.TrimSpace(gjson.GetBytes(body, "ratio").String()),
		strings.TrimSpace(gjson.GetBytes(body, "aspect_ratio").String()),
		strings.TrimSpace(gjson.GetBytes(body, "size").String()),
	)
	info.ImageURLs = collectSeedanceURLArray(body, "images", "reference_images")
	info.StartFrameURLs = collectSeedanceURLArray(body, "start_frame", "first_frame")
	info.EndFrameURLs = collectSeedanceURLArray(body, "end_frame", "last_frame")
	info.ReferenceVideos = collectSeedanceURLArray(body, "reference_videos")
	info.ReferenceAudios = collectSeedanceURLArray(body, "reference_audios")
	info.VideoURL = firstNonEmptySeedanceString(
		strings.TrimSpace(gjson.GetBytes(body, "video_url").String()),
		strings.TrimSpace(gjson.GetBytes(body, "url").String()),
	)
	if videos := gjson.GetBytes(body, "videos"); videos.IsArray() && len(videos.Array()) > 0 {
		if info.VideoURL == "" {
			info.VideoURL = strings.TrimSpace(videos.Array()[0].String())
		}
	}
	if ga := gjson.GetBytes(body, "generate_audio"); ga.Exists() {
		v := ga.Bool()
		info.GenerateAudio = &v
	}

	info.HasReferenceImage = len(info.ImageURLs) > 0 || len(info.StartFrameURLs) > 0 || len(info.EndFrameURLs) > 0
	info.HasReferenceVideo = len(info.ReferenceVideos) > 0
	info.HasReferenceAudio = len(info.ReferenceAudios) > 0

	// Normalize duration/resolution against catalog defaults when possible.
	if spec, ok := LookupSeedanceModel(info.Model); ok {
		if info.DurationSeconds <= 0 {
			if spec.FixedDurationSeconds > 0 {
				info.DurationSeconds = spec.FixedDurationSeconds
			} else {
				info.DurationSeconds = 5
			}
		}
		if info.Resolution == "" {
			info.Resolution = spec.DefaultResolution
		}
	} else if info.DurationSeconds <= 0 {
		info.DurationSeconds = 5
	}
	info.Resolution = NormalizeSeedanceBillingResolutionOrDefault(info.Resolution)
	return info
}

func collectSeedanceURLArray(body []byte, fields ...string) []string {
	out := make([]string, 0)
	seen := make(map[string]struct{})
	for _, field := range fields {
		value := gjson.GetBytes(body, field)
		if !value.Exists() {
			continue
		}
		appendOne := func(raw string) {
			raw = strings.TrimSpace(raw)
			if raw == "" {
				return
			}
			if _, ok := seen[raw]; ok {
				return
			}
			seen[raw] = struct{}{}
			out = append(out, raw)
		}
		switch {
		case value.IsArray():
			for _, item := range value.Array() {
				if item.Type == gjson.String {
					appendOne(item.String())
					continue
				}
				if u := strings.TrimSpace(item.Get("url").String()); u != "" {
					appendOne(u)
				}
			}
		case value.Type == gjson.String:
			appendOne(value.String())
		default:
			if u := strings.TrimSpace(value.Get("url").String()); u != "" {
				appendOne(u)
			}
		}
	}
	return out
}

func firstNonEmptySeedanceString(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// NormalizeSeedanceBillingResolutionOrDefault normalizes resolution labels for billing.
func NormalizeSeedanceBillingResolutionOrDefault(resolution string) string {
	switch strings.ToLower(strings.TrimSpace(resolution)) {
	case "480", "480p", "sd":
		return VideoBillingResolution480P
	case "540", "540p":
		return "540p"
	case "720", "720p", "hd":
		return VideoBillingResolution720P
	case "1080", "1080p", "full_hd", "full-hd", "fhd":
		return VideoBillingResolution1080P
	case "4k", "2160", "2160p", "uhd":
		return "4k"
	default:
		if resolution == "" {
			return VideoBillingResolution720P
		}
		return strings.ToLower(strings.TrimSpace(resolution))
	}
}

// SeedanceVideoRequestSessionHash binds a task_id to the submitting user/api key.
func SeedanceVideoRequestSessionHash(taskID string, userID, apiKeyID int64) string {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" || userID <= 0 || apiKeyID <= 0 {
		return ""
	}
	ownerSeed := fmt.Sprintf("%d:%d:%s", userID, apiKeyID, taskID)
	return "seedance-video:" + DeriveSessionHashFromSeed(ownerSeed)
}

func (s *OpenAIGatewayService) BindSeedanceVideoRequestAccount(
	ctx context.Context,
	groupID *int64,
	taskID string,
	userID, apiKeyID, accountID int64,
) error {
	if s == nil || s.cache == nil {
		return fmt.Errorf("seedance video request binding cache is unavailable")
	}
	sessionHash := SeedanceVideoRequestSessionHash(taskID, userID, apiKeyID)
	cacheKey := s.openAISessionCacheKey(sessionHash)
	if cacheKey == "" || accountID <= 0 {
		return fmt.Errorf("seedance video request binding is invalid")
	}
	ttl := openaiStickySessionTTL
	if s.cfg != nil && s.cfg.Gateway.OpenAIWS.StickySessionTTLSeconds > 0 {
		ttl = time.Duration(s.cfg.Gateway.OpenAIWS.StickySessionTTLSeconds) * time.Second
	}
	return s.cache.SetSessionAccountID(ctx, derefGroupID(groupID), cacheKey, accountID, ttl)
}

func (s *OpenAIGatewayService) ResolveSeedanceVideoRequestAccount(
	ctx context.Context,
	groupID *int64,
	taskID string,
	userID, apiKeyID int64,
) (int64, error) {
	if s == nil || s.cache == nil {
		return 0, fmt.Errorf("seedance video request binding cache is unavailable")
	}
	cacheKey := s.openAISessionCacheKey(SeedanceVideoRequestSessionHash(taskID, userID, apiKeyID))
	if cacheKey == "" {
		return 0, fmt.Errorf("seedance video request binding is invalid")
	}
	return s.cache.GetSessionAccountID(ctx, derefGroupID(groupID), cacheKey)
}

// ForwardSeedanceMedia proxies Binghuo Seedance video endpoints with passthrough semantics.
func (s *OpenAIGatewayService) ForwardSeedanceMedia(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	endpoint SeedanceMediaEndpoint,
	taskID string,
	body []byte,
	contentType string,
) (*OpenAIForwardResult, error) {
	startTime := time.Now()
	if account == nil || !account.IsSeedance() {
		return nil, fmt.Errorf("seedance account is required")
	}
	apiKey := account.GetSeedanceAPIKey()
	if apiKey == "" {
		return nil, fmt.Errorf("seedance api_key is required")
	}

	targetURL, err := buildSeedanceMediaURL(account, s.cfg, endpoint, taskID)
	if err != nil {
		return nil, err
	}

	requestInfo := SeedanceMediaRequestInfo{}
	upstreamModel := ""
	if endpoint.IsGenerationRequest() && gjson.ValidBytes(body) {
		requestInfo = ParseSeedanceMediaRequest(body)
		upstreamModel = requestInfo.Model
		if mapped := strings.TrimSpace(account.GetMappedModel(requestInfo.Model)); mapped != "" && mapped != requestInfo.Model {
			upstreamModel = mapped
			body, err = sjson.SetBytes(body, "model", upstreamModel)
			if err != nil {
				return nil, fmt.Errorf("rewrite seedance account mapped model: %w", err)
			}
			requestInfo.Model = upstreamModel
		}
	}

	method := http.MethodPost
	if endpoint.IsLookupRequest() {
		method = http.MethodGet
	}

	var bodyReader io.Reader
	if endpoint.RequiresRequestBody() {
		bodyReader = bytes.NewReader(body)
	}
	upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
	defer releaseUpstreamCtx()
	upstreamReq, err := http.NewRequestWithContext(upstreamCtx, method, targetURL, bodyReader)
	if err != nil {
		return nil, err
	}
	upstreamReq.Header.Set("Authorization", "Bearer "+apiKey)
	upstreamReq.Header.Set("Accept", "application/json")
	if endpoint.RequiresRequestBody() {
		contentType = strings.TrimSpace(contentType)
		if contentType == "" {
			contentType = "application/json"
		}
		upstreamReq.Header.Set("Content-Type", contentType)
	}
	// Preserve multipart boundary for asset uploads.
	if endpoint.IsUploadRequest() && c != nil {
		if ct := strings.TrimSpace(c.GetHeader("Content-Type")); ct != "" {
			upstreamReq.Header.Set("Content-Type", ct)
		}
	}
	account.ApplyHeaderOverrides(upstreamReq.Header)

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}
	upstreamStart := time.Now()
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	SetOpsLatencyMs(c, OpsUpstreamLatencyMsKey, time.Since(upstreamStart).Milliseconds())
	if err != nil {
		return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, false)
	}
	defer func() { _ = resp.Body.Close() }()

	requestIDHeader := firstNonEmpty(resp.Header.Get("x-request-id"), resp.Header.Get("request-id"))
	if resp.StatusCode >= 400 {
		return s.handleSeedanceMediaErrorResponse(ctx, resp, c, account, requestIDHeader, requestInfo.Model)
	}

	respBody, err := ReadUpstreamResponseBody(resp.Body, s.cfg, c, openAITooLargeError)
	if err != nil {
		return nil, err
	}

	writeSeedanceMediaResponse(c, resp, respBody, s.responseHeaderFilter)

	responseTaskID := extractSeedanceTaskID(respBody)
	usage := seedanceMediaUsageFromResponse(endpoint, requestInfo, respBody)
	return &OpenAIForwardResult{
		RequestID:            requestIDHeader,
		ResponseID:           responseTaskID,
		Usage:                usage.Usage,
		Model:                requestInfo.Model,
		BillingModel:         requestInfo.Model,
		UpstreamModel:        firstNonEmptySeedanceString(upstreamModel, requestInfo.Model),
		ResponseHeaders:      resp.Header.Clone(),
		Duration:             time.Since(startTime),
		VideoCount:           usage.VideoCount,
		VideoResolution:      usage.VideoResolution,
		VideoDurationSeconds: usage.VideoDurationSeconds,
	}, nil
}

type seedanceMediaUsageMetadata struct {
	Usage                OpenAIUsage
	VideoCount           int
	VideoResolution      string
	VideoDurationSeconds int
}

func seedanceMediaUsageFromResponse(endpoint SeedanceMediaEndpoint, requestInfo SeedanceMediaRequestInfo, responseBody []byte) seedanceMediaUsageMetadata {
	meta := seedanceMediaUsageMetadata{}
	// Charge on successful task submission only (not on status poll / upload).
	if !endpoint.IsGenerationRequest() {
		return meta
	}
	taskID := extractSeedanceTaskID(responseBody)
	if taskID == "" {
		return meta
	}
	// Binghuo accepts with {task_id, status:processing}; bill on acceptance.
	meta.VideoCount = 1
	meta.VideoResolution = requestInfo.Resolution
	meta.VideoDurationSeconds = requestInfo.DurationSeconds
	if meta.VideoDurationSeconds <= 0 {
		meta.VideoDurationSeconds = 5
	}
	return meta
}

func extractSeedanceTaskID(body []byte) string {
	if !gjson.ValidBytes(body) {
		return ""
	}
	return firstNonEmptySeedanceString(
		strings.TrimSpace(gjson.GetBytes(body, "task_id").String()),
		strings.TrimSpace(gjson.GetBytes(body, "id").String()),
	)
}

func (s *OpenAIGatewayService) handleSeedanceMediaErrorResponse(
	ctx context.Context,
	resp *http.Response,
	c *gin.Context,
	account *Account,
	requestIDHeader, requestModel string,
) (*OpenAIForwardResult, error) {
	respBody, _ := ReadUpstreamResponseBody(resp.Body, s.cfg, c, openAITooLargeError)
	message := strings.TrimSpace(gjson.GetBytes(respBody, "error").String())
	if message == "" {
		message = strings.TrimSpace(gjson.GetBytes(respBody, "error.message").String())
	}
	if message == "" {
		message = strings.TrimSpace(gjson.GetBytes(respBody, "message").String())
	}
	if message == "" {
		message = truncateString(string(respBody), 512)
	}
	setOpsUpstreamError(c, resp.StatusCode, message, truncateString(string(respBody), 512))

	// 4xx validation errors should not failover; 5xx/429 may.
	failover := &UpstreamFailoverError{
		StatusCode:      resp.StatusCode,
		ResponseBody:    respBody,
		ResponseHeaders: resp.Header.Clone(),
	}
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		// allow retry next account
	} else if resp.StatusCode >= 400 {
		// write error to client and stop
		writeSeedanceMediaErrorResponse(c, resp.StatusCode, "invalid_request_error", message)
		return nil, failover
	}
	return nil, failover
}

func writeSeedanceMediaErrorResponse(c *gin.Context, statusCode int, errType, message string) {
	if c == nil || serviceResponseCommitted(c) {
		return
	}
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"type":    errType,
			"message": message,
		},
	})
}

func writeSeedanceMediaResponse(c *gin.Context, resp *http.Response, body []byte, filter *responseheaders.CompiledHeaderFilter) {
	if c == nil || resp == nil {
		return
	}
	writeOpenAIPassthroughResponseHeaders(c.Writer.Header(), resp.Header, filter)
	contentType := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/json"
	}
	c.Data(resp.StatusCode, contentType, body)
}

func serviceResponseCommitted(c *gin.Context) bool {
	if c == nil {
		return false
	}
	return IsResponseCommitted(c)
}
