package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func DownloadTorrentFile(url, filename string) (string, error) {
	fmt.Println("Downloading torrent file from:", url)
	// Make sure the folder exists
	dir := "files"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Build full save path relative to current working dir
	savePath := filepath.Join(dir, filename)

	// Download file
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create file
	out, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Stream data directly into file
	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	return savePath, nil
}
func main() {
	var link string
	// var file string
	huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("enter magnet link").
				Description("Paste the magnet link here").Value(&link).Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("name cannot be empty")
				}
				return nil
			}),
		).WithShowHelp(true),
	).Run()

	fmt.Println("You entered:", link)

	// magnetURI := os.Args[1]
	clientConfig := torrent.NewDefaultClientConfig()

	// change the directory
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// setup path
	downloadPath := path.Join(user.HomeDir, "Downloads", "fire torrent")

	// check if the folder exist
	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		err = os.MkdirAll(downloadPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	clientConfig.DataDir = downloadPath

	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		fmt.Println("Error creating torrent client:", err)
		panic(err)
	}
	fmt.Println("Torrent client created successfully.")
	defer client.Close()

	// download the torrent file first
	path, err := DownloadTorrentFile(link, "downloaded2.torrent")

	fmt.Println("Torrent file downloaded to:---=-=-=", path)

	if err != nil {
		fmt.Println("Error downloading torrent file:", err)
		panic(err)
	}

	t, err := client.AddTorrentFromFile(path)
	if err != nil {
		panic(err)
	}

	fmt.Println("Fetching torrent metadata...")
	// await info until download is complete
	<-t.GotInfo()

	t.DownloadAll()
	fmt.Println("Downloading:", t.Name())

	for t.BytesMissing() > 0 {
		fmt.Printf("\rProgress: %.2f%%", 100-(float64(t.BytesMissing())/float64(t.Info().TotalLength())*100))
		time.Sleep(2 * time.Second)
	}
	//setup progress bar
	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}

	fmt.Println("\nDownload complete!")
}
