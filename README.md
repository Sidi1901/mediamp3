## Architecture

                    ┌─────────────────────┐
                    │   Browser / UI      │
                    │  html + css + js    │
                    │                     │
                    │ YouTube URL         │
                    │ Format: MP3 / MP4   │
                    │ Quality: 128/320... │
                    └──────────┬──────────┘
                               │ HTTP
                               ▼
                    ┌─────────────────────┐
                    │     Go API          │
                    │                     │
                    │ /info               │
                    │ /download           │
                    │ /progress           │
                    └──────────┬──────────┘

                               │
                    ┌──────────┴──────────┐
                    ▼                     ▼
              ┌───────────┐         ┌───────────┐
              │  yt-dlp   │         │  FFmpeg   │
              │ download  │───────▶│ convert    │
              │ streams   │         │ /merge    │
              └───────────┘         └───────────┘
                                          │
                                          ▼
                                  MP3 / MP4 file


User pastes URL
       │
       ▼
Frontend
       │
       │ POST /api/formats
       ▼
Go BFF
       │
       ▼
YouTube service
       │
       ▼
Get available streams/formats
       │
       ▼
BFF returns formats
       │
       ▼
Frontend displays quality options
       │
       ▼
User selects MP3/MP4 + quality
       │
       ▼
POST /api/download

## Project structure

mediamp3/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── handler/
│   │   ├── format_handler.go
│   │   └── convert_handler.go
│   │
│   ├── route/
│   │   └── router.go
│   │
│   └── service/
│       └── converter_service.go
│
├── web/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   │
│   ├── src/
│   │   ├── js/
│   │   │   ├── app.js
│   │   │   ├── api.js
│   │   │   └── ui.js
│   │   │
│   │   └── css/
│   │       └── style.css
│   │
│   └── public/
│       └── favicon.ico
│
├── go.mod
├── go.sum
└── README.md