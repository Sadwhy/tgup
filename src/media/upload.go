package media

import (
  "fmt"
  "io"
  "os"
  "log"
  "time"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/Sadwhy/tgup/src/utils"
	"github.com/Sadwhy/tgup/src/dif"
)



func Upload(path string) error {

  mime := utils.GetType(path)

	file, err := os.Open(path)
	if err != nil {
	  return fmt.Errorf("u Failed to open file: %s", err)
	}
	defer file.Close()


	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("u Failed to get file info: %s", err)
	}

  if fileInfo.Size() > 10*1024*1024 && mime == "image" {
      mime = "document"
  }

  tracker := NewProgressTracker()

	progressReader := &ProgressReader{
		File:      file,
		TotalSize: fileInfo.Size(),
		ProgressFn: tracker.Callback,
	}

  var media tgbotapi.Chattable

  switch mime {
  case "video":
    video := tgbotapi.NewVideo(dif.ChatId, tgbotapi.FileReader{
	  	Name:   fileInfo.Name(),
		  Reader: progressReader,
	  })
	  video.SupportsStreaming = true
	  video.Duration = utils.Duration(path)
	  video.Thumb = utils.Thumbnail(path)
	  if dif.ForumId != 0 {
	    video.MessageThreadID = dif.ForumId
	  }
	  media = video
  case "image":
    photo := tgbotapi.NewPhoto(dif.ChatId, tgbotapi.FileReader{
      Name:   fileInfo.Name(),
      Reader: progressReader,
    })
	  if dif.ForumId != 0 {
	    photo.MessageThreadID = dif.ForumId
	  }
	  media = photo
  case "audio":
    audio := tgbotapi.NewAudio(dif.ChatId, tgbotapi.FileReader{
  		Name:   fileInfo.Name(),
  		Reader: progressReader,
  	})
	  audio.Duration = utils.Duration(path)
	  if dif.ForumId != 0 {
	    audio.MessageThreadID = dif.ForumId
	  }
	  media = audio
  case "document":
    doc := tgbotapi.NewDocument(dif.ChatId, tgbotapi.FileReader{
  		Name:   fileInfo.Name(),
		  Reader: progressReader,
  	})
	  if dif.ForumId != 0 {
	    doc.MessageThreadID = dif.ForumId
	  }
	  media = doc
	default:
    doc := tgbotapi.NewDocument(dif.ChatId, tgbotapi.FileReader{
  		Name:   fileInfo.Name(),
		  Reader: progressReader,
  	})
	  if dif.ForumId != 0 {
	    doc.MessageThreadID = dif.ForumId
	  }
	  media = doc
  }
 
  log.Println("\nFinishing upload...")
	_, err = dif.Bot.Send(media)
	if err != nil {
		return fmt.Errorf("u Failed to open file: %s", err)
	}

  log.Println("\nUpload complete.")
  os.RemoveAll("download")
  log.Println("\nCleared downloads")
  return nil
}


type ProgressReader struct {
	File       io.Reader
	TotalSize  int64
	Uploaded   int64
	ProgressFn func(current, total int64)
}

func (p *ProgressReader) Read(b []byte) (int, error) {
	n, err := p.File.Read(b)
	if n > 0 {
		p.Uploaded += int64(n)
		if p.ProgressFn != nil {
			p.ProgressFn(p.Uploaded, p.TotalSize)
		}
	}
	return n, err
}


type ProgressTracker struct {
    startTime    time.Time
    lastCurrent  int64
    lastTotal    int64
    lastTime     time.Time
    lastSpeed    float64
}


func (pt *ProgressTracker) Callback(current, total int64) {
    if total != pt.lastTotal {
        pt.startTime = time.Now()
        pt.lastCurrent = current
        pt.lastTotal = total
        pt.lastTime = time.Now()
        pt.lastSpeed = 0
    }

    now := time.Now()
    elapsed := now.Sub(pt.startTime)
    
    // Calculate speed
    var speed float64
    deltaTime := now.Sub(pt.lastTime).Seconds()
    if deltaTime >= 0.1 { // Minimum interval to avoid spikes
        deltaBytes := float64(current - pt.lastCurrent)
        instantSpeed := deltaBytes / deltaTime
        // Smooth the speed
        if pt.lastSpeed == 0 {
            speed = instantSpeed
        } else {
            speed = 0.7*pt.lastSpeed + 0.3*instantSpeed
        }
        pt.lastSpeed = speed
        pt.lastCurrent = current
        pt.lastTime = now
    } else {
        speed = pt.lastSpeed
    }

    // Calculate percentage
    progress := float64(current) / float64(total) * 100

    // output
    output := fmt.Sprintf("\rUploading: %s/%s | %.1f%% | %s/s | %s elapsed",
        utils.FormatBytes(float64(current)),
        utils.FormatBytes(float64(total)),
        progress,
        utils.FormatBytes(speed),
        utils.FormatDuration(elapsed))

    fmt.Print(output)
}


func NewProgressTracker() *ProgressTracker {
    return &ProgressTracker{
        startTime: time.Now(),
        lastTime:  time.Now(),
    }
}