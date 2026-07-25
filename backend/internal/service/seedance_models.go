package service

// Seedance / Binghuo video model catalog used for group model lists and billing defaults.
// Model IDs are case-sensitive and must match the Binghuo downstream handbook.

const DefaultSeedanceBaseURL = "https://api.7tai.cc/v1"

// SeedanceBillingModePerSecond charges unit_price * duration_seconds (* reference multiplier).
// SeedanceBillingModePerRequest charges a fixed price per successful generation.
const (
	SeedanceBillingModePerSecond  = "per_second"
	SeedanceBillingModePerRequest = "per_request"
)

// SeedanceModelSpec describes one publicly supported Binghuo Seedance-family model.
type SeedanceModelSpec struct {
	ID                    string
	BillingMode           string  // per_second | per_request
	UnitPriceCNY          float64 // 元/秒 或 元/次（上游对下游价，默认用作内部参考）
	ReferenceVideoFactor  float64 // 含 reference_videos 时的倍率；1 表示不加价；0 表示不支持参考视频
	DefaultResolution     string  // 480p / 720p / 1080p / 4k
	MinDurationSeconds    int
	MaxDurationSeconds    int
	FixedDurationSeconds  int // >0 表示固定时长
	MaxReferenceImages    int
	MaxReferenceVideos    int
	MaxReferenceAudios    int
	SupportsGenerateAudio bool
	RequiresImage         bool // 图生视频模型必须至少 1 张图
}

// seedanceModelCatalog is the authoritative list of Binghuo Seedance video models.
// Prices are upstream CNY list prices from the handbook; local billing still prefers
// group-level USD overrides when configured.
var seedanceModelCatalog = []SeedanceModelSpec{
	// bh2.0 series (per-second)
	{ID: "bh2.0-fast-480p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.38, ReferenceVideoFactor: 1.8, DefaultResolution: "480p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "bh2.0-fast-720p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.42, ReferenceVideoFactor: 1.8, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "bh2.0-480p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.48, ReferenceVideoFactor: 1.8, DefaultResolution: "480p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "bh2.0-720p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.58, ReferenceVideoFactor: 1.8, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "bh2.0-1080p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.79, ReferenceVideoFactor: 1.8, DefaultResolution: "1080p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "bh2.0-mini-480p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.25, ReferenceVideoFactor: 1.8, DefaultResolution: "480p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "bh2.0-mini-720p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.35, ReferenceVideoFactor: 1.8, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	// SD2.0 series
	{ID: "SD2.0-480P", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.40, ReferenceVideoFactor: 1.0, DefaultResolution: "480p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "SD2.0-720P-fast", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.39, ReferenceVideoFactor: 2.0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "SD2.0-720P", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.50, ReferenceVideoFactor: 2.0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "SD2.0-1080P", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 1.15, ReferenceVideoFactor: 2.0, DefaultResolution: "1080p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "sdvip4k", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 2.20, ReferenceVideoFactor: 2.0, DefaultResolution: "4k", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "sdvip720p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.33, ReferenceVideoFactor: 1.0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "sdvip1080p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.60, ReferenceVideoFactor: 1.0, DefaultResolution: "1080p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	// gz-sd series
	{ID: "gz-sd480p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.19, ReferenceVideoFactor: 2.0, DefaultResolution: "480p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "gz-sd720p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.35, ReferenceVideoFactor: 2.0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "gz-sd1080p", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.70, ReferenceVideoFactor: 2.0, DefaultResolution: "1080p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "gz-sd4k", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 1.60, ReferenceVideoFactor: 2.0, DefaultResolution: "4k", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	// other per-second
	{ID: "sdquan-2-miao", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.275, ReferenceVideoFactor: 0, DefaultResolution: "720p", MinDurationSeconds: 5, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 0, MaxReferenceAudios: 0, SupportsGenerateAudio: true},
	{ID: "wanneng1.1", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.18, ReferenceVideoFactor: 0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 0, MaxReferenceAudios: 0, SupportsGenerateAudio: true},
	{ID: "doubaofast", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.258, ReferenceVideoFactor: 1.0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	// per-request models
	{ID: "sd2-fast福利", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 2.36, ReferenceVideoFactor: 1.0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 4, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "sd2-福利", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 2.98, ReferenceVideoFactor: 1.0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 4, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "B-quannengship2.0", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 3.45, ReferenceVideoFactor: 0, DefaultResolution: "720p", MinDurationSeconds: 5, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 0, MaxReferenceAudios: 0, SupportsGenerateAudio: true},
	{ID: "quanneng2.0", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 5.25, ReferenceVideoFactor: 0, DefaultResolution: "720p", FixedDurationSeconds: 15, MinDurationSeconds: 15, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 0, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "quanneng2.0-9tu", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 1.58, ReferenceVideoFactor: 0, DefaultResolution: "720p", FixedDurationSeconds: 15, MinDurationSeconds: 15, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 0, MaxReferenceAudios: 0, SupportsGenerateAudio: true},
	{ID: "video2.0", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 4.85, ReferenceVideoFactor: 1.0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "sd2-vip720p", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 3.55, ReferenceVideoFactor: 0, DefaultResolution: "720p", FixedDurationSeconds: 15, MinDurationSeconds: 15, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 0, MaxReferenceAudios: 0, SupportsGenerateAudio: true},
	{ID: "sd2-vip720p-fast", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 3.75, ReferenceVideoFactor: 1.0, DefaultResolution: "720p", MinDurationSeconds: 4, MaxDurationSeconds: 15, MaxReferenceImages: 9, MaxReferenceVideos: 3, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "keling-3", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 0.90, ReferenceVideoFactor: 0, DefaultResolution: "720p", FixedDurationSeconds: 15, MinDurationSeconds: 15, MaxDurationSeconds: 15, MaxReferenceImages: 2, MaxReferenceVideos: 0, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "xb-sora2", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 0.78, ReferenceVideoFactor: 0, DefaultResolution: "720p", MinDurationSeconds: 8, MaxDurationSeconds: 12, MaxReferenceImages: 1, MaxReferenceVideos: 0, MaxReferenceAudios: 3, SupportsGenerateAudio: true, RequiresImage: true},
	{ID: "me-kuaile1.0", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 1.85, ReferenceVideoFactor: 0, DefaultResolution: "720p", MinDurationSeconds: 5, MaxDurationSeconds: 15, MaxReferenceImages: 5, MaxReferenceVideos: 0, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "sora-2-z", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 0.88, ReferenceVideoFactor: 0, DefaultResolution: "720p", FixedDurationSeconds: 12, MinDurationSeconds: 12, MaxDurationSeconds: 12, MaxReferenceImages: 1, MaxReferenceVideos: 0, MaxReferenceAudios: 3, SupportsGenerateAudio: true},
	{ID: "veo-omni-flash", BillingMode: SeedanceBillingModePerRequest, UnitPriceCNY: 0.88, ReferenceVideoFactor: 0, DefaultResolution: "720p", FixedDurationSeconds: 10, MinDurationSeconds: 10, MaxDurationSeconds: 10, MaxReferenceImages: 5, MaxReferenceVideos: 0, MaxReferenceAudios: 3, SupportsGenerateAudio: true, RequiresImage: true},
	// subtitle removal utility
	{ID: "去字幕", BillingMode: SeedanceBillingModePerSecond, UnitPriceCNY: 0.009, ReferenceVideoFactor: 0, DefaultResolution: "720p", MinDurationSeconds: 1, MaxDurationSeconds: 600, MaxReferenceImages: 0, MaxReferenceVideos: 0, MaxReferenceAudios: 0, SupportsGenerateAudio: false},
}

var seedanceModelByID map[string]*SeedanceModelSpec

func init() {
	seedanceModelByID = make(map[string]*SeedanceModelSpec, len(seedanceModelCatalog))
	for i := range seedanceModelCatalog {
		spec := &seedanceModelCatalog[i]
		seedanceModelByID[spec.ID] = spec
	}
}

func seedanceDefaultModelIDs() []string {
	return SeedanceDefaultModelIDs()
}

// SeedanceDefaultModelIDs returns the public catalog model IDs (case-sensitive).
func SeedanceDefaultModelIDs() []string {
	ids := make([]string, 0, len(seedanceModelCatalog))
	for _, spec := range seedanceModelCatalog {
		ids = append(ids, spec.ID)
	}
	return ids
}

// LookupSeedanceModel returns the catalog entry for model ID (case-sensitive).
func LookupSeedanceModel(modelID string) (*SeedanceModelSpec, bool) {
	spec, ok := seedanceModelByID[modelID]
	return spec, ok
}
