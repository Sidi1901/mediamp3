package service

import(
	"context"
	"os/exec"
	"encoding/json"
	"strings"
	"sort"
	"github.com/Sidi1901/mediamp3/internal/model"

)

type FormatService interface {
	GetFormats(ctx context.Context, videoURL string) (model.FormatResponse,error)
}

type formatService struct {
	cfg config.Config
}

func NewFormatService(config config.Config) FormatService {
	return &formatService{
		cfg : config,
	}
}

type ytDlpFormat struct {
	FormatID string `json:"format_id"`
	Ext      string `json:"ext"`
	Height   int    `json:"height"`
	ACodec   string `json:"acodec"`
	VCodec   string `json:"vcodec"`
}

type ytDlpInfo struct {
	Formats []ytDlpFormat `json:"formats"`
}

func (s *formatService) GetFormats(ctx context.Context, videoURL string) (model.FormatResponse, error) {
	
	cmd := exec.CommandContext(
		ctx,
		s.cfg.YtDlpPath,
		"--no-warnings",
		"--skip-download",
		"--dump-single-json",
		"--no-playlist",
		"--",
		videoURL,
	)

	output, err := cmd.Output()

	if err != nil {
		return model.FormatResponse{}, fmt.Errorf("Inspect video format %w", err)
	}

	var info ytDlpInfo

	if err := json.Unmarshal(output, &info); err != nil {
		return model.FormatResponse, fmt.Errorf("parse video formats: %w", err)
	} 

	mp3Qualities := []string{"128kbps","192kbps", "320kbps"}
	mp4Qualities := map[string]struct{}

	for _, format := range info.Formats {
		if format.Height <= 0 || !strings.EqualFold(format.Ext, "mp4") || format.VCodec == "" || format.VCodec == "none"{
			continue;
		}

		mp4QualitySet[quality] = struct{}{}
	}

	var heights []int

	for h := range mp4QualitySet {
		heights = append(heights, h)
	}

	sort.Ints(heights)

	var mp4Qualities []string

	for _, h := range heights {
		mp4Qualities = append(mp4Qualities, fmt.Sprintf("%dp", h))
	}

	modelResponse := model.FormatResponse{
		Formats:  []model.Format{
			{Type:"mp3", Qualities: mp3Qualities},
			{Type:"mp4", Qualities: mp4Qualities}, 
		},
	}

	return modelResponse, nil
}