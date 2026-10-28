
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strconv"

	"github.com/Sidi1901/mediamp3/internal/config"
	"github.com/Sidi1901/mediamp3/internal/model"
)

type FormatService interface {
	GetFormats(ctx context.Context, videoURL string) (model.FormatResponse, error)
}

type formatService struct {
	cfg config.Config
}

func NewFormatService(cfg config.Config) FormatService {
	return &formatService{
		cfg: cfg,
	}
}

type ytDlpResponse struct {
	Formats []ytDlpFormat `json:"formats"`
}

type ytDlpFormat struct {
	Ext       string `json:"ext"`
	Height    int    `json:"height"`
	VCodec    string `json:"vcodec"`
	ACodec    string `json:"acodec"`
}

func (s *formatService) GetFormats(
	ctx context.Context,
	videoURL string,
) (model.FormatResponse, error) {
	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		"--no-warnings",
		"--skip-download",
		"--dump-single-json",
		"--no-playlist",
		"--",
		videoURL,
	)

	output, err := cmd.Output()
	if err != nil {
		return model.FormatResponse{}, fmt.Errorf(
			"inspect video formats: %w",
			err,
		)
	}

	var result ytDlpResponse
	if err := json.Unmarshal(output, &result); err != nil {
		return model.FormatResponse{}, fmt.Errorf(
			"parse yt-dlp response: %w",
			err,
		)
	}

	heights := make(map[int]struct{})

	for _, format := range result.Formats {
		if format.Ext != "mp4" {
			continue
		}

		if format.Height <= 0 || format.VCodec == "none" {
			continue
		}

		heights[format.Height] = struct{}{}
	}

	var resolutions []int
	for height := range heights {
		resolutions = append(resolutions, height)
	}

	sort.Ints(resolutions)

	qualities := make([]string, 0, len(resolutions))
	for _, height := range resolutions {
		qualities = append(
			qualities,
			strconv.Itoa(height)+"p",
		)
	}

	return model.FormatResponse{
		Formats: []model.Format{
			{
				Type:      "mp3",
				Qualities: []string{"128kbps", "192kbps", "320kbps"},
			},
			{
				Type:      "mp4",
				Qualities: qualities,
			},
		},
	}, nil
}
