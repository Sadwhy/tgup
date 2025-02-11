# tgup  

A simple CLI script to upload media from URLs to a Telegram channel.  

## Prerequisites  

Ensure the following dependencies are installed:  

1. **ffmpeg** – Generates thumbnails, retrieves duration, and remuxes/merges m3u8 streams for `yt-dlp`.  
2. **aria2** – The primary downloader used for all downloads.  
3. **yt-dlp** – Supports m3u8 downloads and additional sites, using `aria2` as an external downloader.  
4. **Go Compiler** – Required to build and run the script.  

---

## Installation  

### 1. Install `ffmpeg`, `aria2`, and `yt-dlp`  

- **Ubuntu/Debian:**  
  ```bash  
  sudo apt update  
  sudo apt install ffmpeg aria2  
  sudo pip install yt-dlp  
  ```  

- **macOS (Homebrew):**  
  ```bash  
  brew install ffmpeg aria2 yt-dlp  
  ```  

- **Windows:**  
  Download and install them from their official websites:  
  - [ffmpeg](https://ffmpeg.org/download.html)  
  - [aria2](https://aria2.github.io/)  
  - [yt-dlp](https://github.com/yt-dlp/yt-dlp)  

### 2. Build the Telegram Bot API  

Follow [this guide](https://tdlib.github.io/telegram-bot-api/build.html) to build the Telegram API server.  
After building, save the generated `telegram-bot-api` binary for later use.  

### 3. Install the Go Compiler  

- **Ubuntu/Debian:**  
  ```bash  
  sudo apt install golang  
  ```  
- **macOS (Homebrew):**  
  ```bash  
  brew install go  
  ```  
- **Windows:**  
  Download from the [official Go website](https://golang.org/dl/).  

---

## Setup & Usage  

### 1. Clone the Repository  

```bash  
git clone https://github.com/Sadwhy/tgup.git  
cd tgup  
```  

### 2. Install Dependencies  

```bash  
go mod tidy  
```  

### 3. Build the Script  

```bash  
go build -o tgup ./src/main.go  
```  

This generates a `tgup` binary.  

### 4. Prepare the Environment  

- Place the `telegram-bot-api` binary in the same directory as `tgup`.  
- Create a `.env` file with the following:  

  ```env  
  BOT_TOKEN=12345678:qwertyuiopasdfghjkl  
  API_HASH=zxcvbnmasdfghjkl  
  API_ID=1234567890  
  ```  

- Get your API credentials:  
  - [API_ID & API_HASH](https://core.telegram.org/api/obtaining_api_id)  
  - [BOT_TOKEN](https://t.me/BotFather)  

⚠️ **Keep these credentials private** to prevent account compromise.  

---

## Using the Script  

### Single File Upload  

```bash  
./tgup  
```  

- Enter your Channel ID (and Topic ID for supergroups, if applicable).  
- Provide the media URL when prompted.  

### Multi-File Upload  

1. Create a `list.json` file with media links:  

   ```json  
   [  
     "https://example.com/file1.mp4",  
     "https://example.com/file2.jpg",  
     "https://example.com/file3.png"  
   ]  
   ```  

2. Run the script:  

   ```bash  
   ./tgup  
   ```  

- Enter your Channel/Group ID (and Topic ID if using supergroups).  
- The script will process each file in `list.json`.  

---

## Troubleshooting  

- **Missing Dependencies?** Verify `ffmpeg`, `aria2`, `yt-dlp`, and Go are installed.  
- **Invalid API Credentials?** Double-check `.env` values (`BOT_TOKEN`, `API_HASH`, `API_ID`).  
- **Execution Issues?** Ensure `tgup` and `telegram-bot-api` binaries have execute permissions.