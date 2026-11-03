package config

import(
	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	YtDlpPath string `env:"YT_DLP_PATH" envDefault:"yt-dlp"`
	FFmpegPath string `env:"FFMPEG_PATH" envDefault:"ffmpeg"`
	OutputDir string `env:"OUTPUT_DIR" envDefault:"./tmp/mediamp3"`
}


func LoadEnv() (Config,error) {

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}