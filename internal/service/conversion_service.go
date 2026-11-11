package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Sidi1901/mediamp3/internal/config"
	"github.com/Sidi1901/mediamp3/internal/model"
	"github.com/Sidi1901/mediamp3/internal/validation"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ConversionResult struct {
	Response model.ConversionResponse
	FilePath string
	Format   string
}

var (
	ErrInvalidConversionRequest = errors.New("invalid conversion request")
	ErrConversionNotImplemented = errors.New("conversion functionality is not implemented yet")
)

type ConversionService interface {
	ConvertFormat(ctx context.Context, request model.ConversionRequest) (ConversionResult, error)
}

type conversionService struct {
	cfg config.Config
}

func NewConversionService(cfg config.Config) ConversionService {
	return &conversionService{
		cfg: cfg,
	}
}

func (s *conversionService) ConvertFormat(ctx context.Context, request model.ConversionRequest) (ConversionResult, error) {

	request.URL = strings.TrimSpace(request.URL)
	request.Format = strings.ToLower(strings.TrimSpace(request.Format))
	request.Quality = strings.ToLower(strings.TrimSpace(request.Quality))

	if err := validateConversionRequest(request); err != nil {
		return ConversionResult{}, err
	}

	id, err := generateID()

	if err != nil {
		return ConversionResult{}, fmt.Errorf("generate conversion ID: %w", err)
	}

	conversionDir := filepath.Join(s.cfg.OutputDir, id)

	if err := os.MkdirAll(conversionDir, 0750); err != nil {
		return ConversionResult{}, fmt.Errorf("create output directory: %w", err)
	}

	conversionSucceeded := false

	defer func() {
		if !conversionSucceeded {
			_ = os.RemoveAll(conversionDir)
		}
	}()

	outputTemplate := filepath.Join(conversionDir, "output.%(ext)s")

	if err := s.runYtDlp(ctx, request, outputTemplate); err != nil {
		return ConversionResult{}, fmt.Errorf(
			"convert media: %w", err,
		)
	}

	outputPath, err := findOutputFile(conversionDir, request.Format)
	if err != nil {
		return ConversionResult{}, err
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		return ConversionResult{}, fmt.Errorf(
			"inspect converted file: %w", err,
		)
	}

	if !info.Mode().IsRegular() || info.Size() == 0 {
		return ConversionResult{}, fmt.Errorf(
			"converted output is empty or invalid",
		)
	}

	conversionSucceeded = true

	return ConversionResult{
		Response: model.ConversionResponse{
			ID:     id,
			Status: model.ConversionStatusCompleted,
		},
		FilePath: outputPath,
		Format:   request.Format,
	}, nil

}

func (s *conversionService) runYtDlp(ctx context.Context, request model.ConversionRequest, outputTemplate string) error {

	args := []string{"--no-playlist", "--ffmpeg-location", s.cfg.FFmpegPath, "-o", outputTemplate}

	switch request.Format {
	case "mp3":
		audioQuality, err := parseAudioQuality(request.Quality)
		if err != nil {
			return err
		}

		args = append(
			args,
			"-x",
			"--audio-format", "mp3",
			"--audio-quality", audioQuality,
		)

	case "mp4":
		height, err := parseVideoQuality(request.Quality)

		if err != nil {
			return err
		}

		format := fmt.Sprintf("bv*[height<=%d][ext=mp4]+ba[ext=m4a]/b[height<=%d][ext=mp4]", height, height)

		args = append(args, "-f", format, "--merge-output-format", "mp4")

	default:
    	return ErrInvalidConversionRequest
	}

	args = append(args, request.URL)

	cmd := exec.CommandContext(ctx, s.cfg.YtDlpPath, args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("yt-dlp execution failed: %w: %s",
			err, string(output))
	}

	return nil
}

func generateID() (string, error) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func validateConversionRequest(request model.ConversionRequest) error {

	// valid host/domain
	if validation.IsValidYoutubeURL(request.URL){

	}else{
		return fmt.Errorf("%w: unsupported URL host", ErrInvalidConversionRequest)
	}


	switch request.Format {
	case "mp3":
		switch request.Quality {
		case "128kbps", "192kbps", "320kbps":
			return nil
		default:
			return fmt.Errorf("%w: unsupported MP3 quality", ErrInvalidConversionRequest)
		}

	case "mp4":
		switch request.Quality {
		case "360p", "480p", "720p", "1080p", "1440p", "2160p":
			return nil
		default:
			return fmt.Errorf("%w: unsupported MP4 quality", ErrInvalidConversionRequest)
		}

	default:
		return fmt.Errorf("%w: unsupported output format", ErrInvalidConversionRequest)
	}

}


func parseAudioQuality(quality string) (string, error) {
    switch quality {
    case "128kbps":
        return "128K", nil
    case "192kbps":
        return "192K", nil
    case "320kbps":
        return "320K", nil
    default:
        return "", ErrInvalidConversionRequest
    }
}


func parseVideoQuality(quality string) (int, error) {
	switch quality {
	case "360p":
		return 360, nil
	case "480p":
		return 480, nil
	case "720p":
		return 720, nil
	case "1080p":
		return 1080, nil
	case "1440p":
		return 1440, nil
	case "2160p":
		return 2160, nil
	default:
    return 0, ErrInvalidConversionRequest
	}
}

func findOutputFile(dir, format string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read conversion directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) == "."+format &&
			name == "output."+format {
			return filepath.Join(dir, name), nil
		}
	}

	return "", fmt.Errorf(
		"expected %s output was not created", format,
	)
}
