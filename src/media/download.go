package media

import (
  "fmt"
  "os"
  "io/fs"
  "path/filepath"
  "os/exec"
  "time"
  "strings"
   "github.com/Sadwhy/tgup/src/utils"
  )

const header = "User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"

var speedParams = []string{
  "--max-concurrent-downloads=10",
  "--split=16",
  "--max-connection-per-server=12",
  "--min-split-size=8M",
  "--max-tries=10",
  "--retry-wait=10",
  "--connect-timeout=30",
  "--timeout=60",
  "--auto-file-renaming=true",
  "--allow-overwrite=true",
  "--disk-cache=64M",
  "--piece-length=1M",
  "--optimize-concurrent-downloads=true",
  }

var torrentParams = []string{
  "--seed-time=0",
  "--bt-enable-lpd=true",
  "--enable-dht=true",
  "--enable-peer-exchange=true",
  "--bt-max-peers=500",
  "--bt-request-peer-speed-limit=50M",
  "--max-overall-upload-limit=50M",
  "--seed-ratio=0.0",
  }

var ytdlpParams = []string{
  "--max-concurrent-downloads=4",
  "--split=4",
  "--max-connection-per-server=4",
  "--min-split-size=8M",
  "--max-tries=10",
  "--retry-wait=10",
  "--connect-timeout=30",
  "--timeout=60",
  "--auto-file-renaming=true",
  "--allow-overwrite=true",
  "--disk-cache=64M",
  "--piece-length=1M",
  "--optimize-concurrent-downloads=true",
}


func Download(link string) error {
  var downloadPath string
  os.MkdirAll("download", os.ModePerm)
  isTorrent := strings.HasPrefix(link, "magnet:") || strings.HasSuffix(link, ".torrent")
  
  useYtdlp := false
  
  if !isTorrent {
    useYtdlp = utils.ForYtdlp(link)
    torrentParams = []string{}
  }
  
  var cmd *exec.Cmd

	if useYtdlp {
	  
	  ytdlpArgs := strings.Join(ytdlpParams, " ")
	  
		args := []string{
			link,
			"--output",
			"download/%(title)s.%(ext)s",
			"--remux-video", "aac>mp4/mov>mp4/m4a>mp4/mkv>mp4",
			"--downloader", "aria2c",
			"--downloader-args", "aria2c:" + ytdlpArgs,
		}
		cmd = exec.Command("yt-dlp", args...)
	} else {
		args := append([]string{
			"-d", "download",
			"--console-log-level=warn",
			"--summary-interval=1",
			"--show-console-readout=true",
			"--header", header,
			link,
		}, append(speedParams, torrentParams...)...)
		cmd = exec.Command("aria2c", args...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
	  return fmt.Errorf("d Download failed: %s", err)
	}

	files, err := os.ReadDir("download")
	if err != nil {
		return fmt.Errorf("d Failed to read download directory: %s", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("d No files found in download directory")
	}

	var recentFile fs.DirEntry
	var recentTime time.Time

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return fmt.Errorf("d Failed to get file info in download.go: %s", err)
		}

		if info.ModTime().After(recentTime) {
			recentTime = info.ModTime()
			recentFile = file
		}
	}

	if recentFile == nil {
		return fmt.Errorf("d Failed to determine the most recent file")
	}

	downloadPath = filepath.Join("download", recentFile.Name())
	fmt.Printf("Download completed: %s\n", recentFile.Name())
	fmt.Println(downloadPath)

	err = Upload(downloadPath)
	if err != nil {
	  return fmt.Errorf("d Failed to upload file: %s", err)
	}
	
	return nil
}