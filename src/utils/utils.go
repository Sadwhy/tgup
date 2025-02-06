package utils

import (
  "context"
  "fmt"
  "bufio"
  "os"
  "encoding/json"
  "io"
  "net"
  "net/http"
  "os/exec"
  "path/filepath"
  "strings"
  "strconv"
  "time"
  "crypto/tls"
  "github.com/gabriel-vasile/mimetype"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
  )

var supportedExtensions = []string{
    "mp4", "mkv", "avi", "mov", "flv", "wmv", "webm",   // Video formats
    "mp3", "aac", "m4a", "flac", "opus", "ogg", "wav",  // Audio formats
    "srt", "vtt", "ass", "ttml",  // Subtitle formats
}  

func GetInput(prompt string) string {
  reader := bufio.NewReader(os.Stdin)
  fmt.Print(prompt)
  input, _ := reader.ReadString('\n')
  return strings.TrimSpace(input)
}

func ForYtdlp(link string) bool {
  cmd := exec.Command(
    "yt-dlp",
    "--dump-json",
    link,
    )

  out, _ := cmd.Output()

	var data map[string]interface{}
	json.Unmarshal(out, &data)

	if ext, ok := data["ext"].(string); ok {
		for _, supportedExt := range supportedExtensions {
			if strings.Contains(strings.ToLower(ext), supportedExt) {
				fmt.Println("Using yt-dlp")
				return true
			}
		}
	}

	if ext, ok := data["format"].(string); ok {
		for _, supportedExt := range supportedExtensions {
			if strings.Contains(strings.ToLower(ext), supportedExt) {
				fmt.Println("Using yt-dlp")
				return true
			}
		}
	}
	fmt.Println("Using Aria2")
	return false
}

type ThumbStruct struct {
	FilePath string
}

func (lf ThumbStruct) NeedsUpload() bool {
	return true
}

func (lf ThumbStruct) UploadData() (string, io.Reader, error) {
	file, err := os.Open(lf.FilePath)
	if err != nil {
		return "", nil, err
	}
	return filepath.Base(lf.FilePath), file, nil
}

func (lf ThumbStruct) SendData() string {
	return ""
}

func Thumbnail(path string) tgbotapi.RequestFileData {
  thumbPath := filepath.Join(filepath.Dir(path), filepath.Base(path[:len(path)-len(filepath.Ext(path))]) + "_thumb.jpg")
  
  cmd := exec.Command(
    "ffmpeg",
    "-i", path,
    "-ss", "00:00:05",
    "-frames:v", "1",
    "-y",
    thumbPath,
    )

  _, err := cmd.CombinedOutput()

  if err != nil{
    fmt.Println("Failed to create Thumbnail:\n", err)
  }
  return ThumbStruct{FilePath: thumbPath}
}

func Duration(path string) int {
  cmd := exec.Command(
    "ffprobe", "-v",
    "error", "-show_entries",
    "format=duration", "-of",
    "default=noprint_wrappers=1:nokey=1",
    path,
    )

  dura, err := cmd.CombinedOutput()
  if err != nil {
    return 0
  }

  duraStr := strings.TrimSpace(string(dura))
  
  duration, err := strconv.ParseFloat(duraStr, 64)
  if err != nil {
    return 0
  }
  return int(duration)
}

func GetType(path string) string {
  if mType, _ := mimetype.DetectFile(path); mType != nil {
    mime := strings.Split(mType.String(), ";")[0]
    switch {
    case strings.HasPrefix(mime, "video/"):
      return "video"
    case strings.HasPrefix(mime, "image/"):
      return "image"
    case strings.HasPrefix(mime, "audio/"):
      return "audio"
    default:
      return "document"
    }
  }
  return "document"
}

func FormatBytes(bytes float64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%.2f B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.2f %ciB", bytes/float64(div), "KMGTPE"[exp])
}

func FormatDuration(d time.Duration) string {
    hours := int(d.Hours())
    minutes := int(d.Minutes()) % 60
    seconds := int(d.Seconds()) % 60
    return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func GetRequest(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

  client := HTTPClient()

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func HTTPClient() *http.Client {
	dnsIP := "8.8.8.8:53"

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 5 * time.Second,
				}
				return d.DialContext(ctx, "udp", dnsIP)
			},
		},
	}

	transport := &http.Transport{
		DialContext: dialer.DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &http.Client{
		Transport: transport,
	}
}
