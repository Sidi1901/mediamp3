

## Tech stack

### Frotend
Html Css Javascript

### Backend

GoLang

### External Services / APIs

yt-dlp - It acts as the fetcher: it parses URLs, bypasses restrictions, selects streams, and pulls raw data.   FFmpeg - It acts as the processor: it merges separate high-definition video and audio tracks, handles format conversions, embeds subtitles/metadata, and extracts audio.   

### Build Tools
Vite for frontend
Go for Backend



## yt-dlp | FFmpeg

yt-dlp works through a few systematic steps:

Parsing & Extraction: It takes the URL you pass it, analyzes the page structure, and extracts the raw streaming formats and metadata directly from the source site's internal APIs (improving heavily on the legacy youtube-dl extractor logic).

Format Selection: By default, it looks for the highest quality video and audio streams available (e.g., separate DASH video and audio tracks on YouTube).

Downloading & Muxing: It downloads the streams via multi-threaded chunking, and if FFmpeg is detected on your system, it automatically runs a post-processing pass to mux (merge) the video and audio into a single playback-ready .mkv or .mp4 container.


#### How Standalone Tools & Backends Handle YouTube Extraction
To build a reliable software service or backend API that extracts YouTube videos for end-users, you should combine two architectural components:

┌───────────────────────────────┐      ┌───────────────────────────────┐
│     JavaScript Runtime        │  +   │      PO Token Provider        │
│   (Node.js / Deno / Bun)      │      │    (Proof of Origin Tokens)   │
└───────────────┬───────────────┘      └───────────────┬───────────────┘
                │                                      │
                └───────────────────┬──────────────────┘
                                    │
                         ┌──────────▼──────────┐
                         │   yt-dlp Execution  │
                         └─────────────────────┘

1. Bundle a Local JavaScript Engine
YouTube’s sign-in bot detection and player signature challenges (n challenge, EJS) depend on executing obfuscated client-side JavaScript.

When running yt-dlp in a server/production environment, ensure Node.js, Deno, or Bun is installed on the host system. yt-dlp automatically detects the JS runtime to evaluate signature challenges dynamically without requiring human authentication.

2. Implement a PO Token (Proof of Origin) Provider
YouTube requires clients (especially non-web clients) to present a valid Proof-of-Origin (PO) Token generated via Botguard JS execution.

Instead of cookies, production systems run a background daemon or container (like bgutil-ytdlp-pot-provider or similar microservices) that generates fresh, valid PO tokens on demand.

#### Deno
Deno is acting as a background JavaScript execution engine. Deno handles the JavaScript code execution required to crack YouTube's security layer

When yt-dlp runs, it passes those obfuscated YouTube challenge scripts directly to deno.exe. Deno safely executes the JavaScript code in a sandboxed environment, computes the required cryptographic signatures, and returns them to yt-dlp so it can successfully unlock and download the video streams.

**YouTube's JS Challenges (EJS):** YouTube protects its video streams by forcing clients to execute heavily obfuscated JavaScript code snippets (called EJS or signature/n challenges) to decrypt streaming URLs.

**No JS runtime Error**
No supported JavaScript runtime could be found. Only deno is enabled by default; to use another runtime add --js-runtimes RUNTIME[:PATH]  to your command/config. YouTube extraction without a JS runtime has been deprecated, and some formats may be missing. See  https://github.com/yt-dlp/yt-dlp/wiki/EJS  for details on installing one


## Alternate without Deno



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