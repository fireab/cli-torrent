# CLI Torrent Downloader

A simple command-line tool to download torrents using a magnet link or torrent file, with a progress bar UI.

## Features

- Download torrent files from a URL
- Fetch and download torrent contents
- Shows download progress in the terminal

## Requirements

- Go 1.24+
- Internet connection

## Usage

1. Build the project:

   ```sh
   go build -o cli-torrent
   ```

2. Run the program:

   ```sh
   ./cli-torrent
   ```

3. Paste your magnet link when prompted.

## Project Structure

- [`main.go`](main.go): Main application logic
- [`progress.go`](progress.go): Progress bar UI
- [`files/`](files): Downloaded torrent files

## License

MIT
