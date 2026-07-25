package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSeedanceMediaRequest_TextToVideo(t *testing.T) {
	body := []byte(`{
		"model": "bh2.0-fast-720p",
		"prompt": "电影感城市夜景",
		"duration": 8,
		"ratio": "16:9",
		"generate_audio": true
	}`)
	info := ParseSeedanceMediaRequest(body)
	require.Equal(t, "bh2.0-fast-720p", info.Model)
	require.Equal(t, "电影感城市夜景", info.Prompt)
	require.Equal(t, 8, info.DurationSeconds)
	require.Equal(t, "720p", info.Resolution)
	require.Equal(t, "16:9", info.Ratio)
	require.NotNil(t, info.GenerateAudio)
	require.True(t, *info.GenerateAudio)
	require.False(t, info.HasReferenceVideo)
}

func TestParseSeedanceMediaRequest_Multimodal(t *testing.T) {
	body := []byte(`{
		"model": "SD2.0-720P",
		"prompt": "保持角色外观",
		"duration": 10,
		"images": ["https://example.com/a.png"],
		"reference_videos": ["https://example.com/v.mp4"],
		"reference_audios": ["https://example.com/a.mp3"]
	}`)
	info := ParseSeedanceMediaRequest(body)
	require.Equal(t, "SD2.0-720P", info.Model)
	require.True(t, info.HasReferenceImage)
	require.True(t, info.HasReferenceVideo)
	require.True(t, info.HasReferenceAudio)
	require.Equal(t, []string{"https://example.com/a.png"}, info.ImageURLs)
	require.Equal(t, []string{"https://example.com/v.mp4"}, info.ReferenceVideos)
	require.Equal(t, 10, info.DurationSeconds)
	require.Equal(t, "720p", info.Resolution)
}

func TestParseSeedanceMediaRequest_Aliases(t *testing.T) {
	body := []byte(`{
		"model_id": "quanneng2.0",
		"text": "hello",
		"seconds": 15,
		"aspect_ratio": "9:16",
		"reference_images": ["https://example.com/b.png"],
		"first_frame": ["https://example.com/start.png"],
		"last_frame": ["https://example.com/end.png"]
	}`)
	info := ParseSeedanceMediaRequest(body)
	require.Equal(t, "quanneng2.0", info.Model)
	require.Equal(t, "hello", info.Prompt)
	require.Equal(t, 15, info.DurationSeconds)
	require.Equal(t, "9:16", info.Ratio)
	require.True(t, info.HasReferenceImage)
	require.Equal(t, []string{"https://example.com/start.png"}, info.StartFrameURLs)
	require.Equal(t, []string{"https://example.com/end.png"}, info.EndFrameURLs)
}

func TestLookupSeedanceModel(t *testing.T) {
	spec, ok := LookupSeedanceModel("bh2.0-720p")
	require.True(t, ok)
	require.Equal(t, SeedanceBillingModePerSecond, spec.BillingMode)
	require.Equal(t, 1.8, spec.ReferenceVideoFactor)

	spec, ok = LookupSeedanceModel("sd2-vip720p")
	require.True(t, ok)
	require.Equal(t, SeedanceBillingModePerRequest, spec.BillingMode)
	require.Equal(t, 15, spec.FixedDurationSeconds)

	_, ok = LookupSeedanceModel("does-not-exist")
	require.False(t, ok)
}

func TestBuildSeedanceMediaURL(t *testing.T) {
	account := &Account{
		Platform: PlatformSeedance,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"base_url": "https://api.7tai.cc/v1",
			"api_key":  "test-key",
		},
	}
	url, err := buildSeedanceMediaURL(account, nil, SeedanceMediaEndpointVideoGenerations, "")
	require.NoError(t, err)
	require.Equal(t, "https://api.7tai.cc/v1/video/generations", url)

	url, err = buildSeedanceMediaURL(account, nil, SeedanceMediaEndpointVideoStatus, "task_abc 1")
	require.NoError(t, err)
	require.Equal(t, "https://api.7tai.cc/v1/video/generations/task_abc%201", url)

	url, err = buildSeedanceMediaURL(account, nil, SeedanceMediaEndpointAssetsUploads, "")
	require.NoError(t, err)
	require.Equal(t, "https://api.7tai.cc/v1/assets/uploads", url)
}

func TestSeedanceMediaUsageFromResponse_OnlyOnAcceptedTask(t *testing.T) {
	info := SeedanceMediaRequestInfo{
		Model:           "bh2.0-fast-720p",
		DurationSeconds: 5,
		Resolution:      "720p",
	}
	meta := seedanceMediaUsageFromResponse(SeedanceMediaEndpointVideoGenerations, info, []byte(`{"task_id":"task_1","status":"processing"}`))
	require.Equal(t, 1, meta.VideoCount)
	require.Equal(t, 5, meta.VideoDurationSeconds)
	require.Equal(t, "720p", meta.VideoResolution)

	meta = seedanceMediaUsageFromResponse(SeedanceMediaEndpointVideoStatus, info, []byte(`{"task_id":"task_1","status":"succeeded"}`))
	require.Equal(t, 0, meta.VideoCount)

	meta = seedanceMediaUsageFromResponse(SeedanceMediaEndpointVideoGenerations, info, []byte(`{"error":"bad"}`))
	require.Equal(t, 0, meta.VideoCount)
}

func TestExtractSeedanceTaskID(t *testing.T) {
	require.Equal(t, "task_x", extractSeedanceTaskID([]byte(`{"task_id":"task_x"}`)))
	require.Equal(t, "id_y", extractSeedanceTaskID([]byte(`{"id":"id_y"}`)))
	require.Equal(t, "", extractSeedanceTaskID([]byte(`{}`)))
}
