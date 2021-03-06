package telebot

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

// YtDownload handle video string and send path from saved file
func YtDownload(command string) (string, error) {
	splited := strings.Split(command, " ")
	videoID := splited[1]
	client := youtube.Client{}
	var path string

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
		return "", err
	}

	resp, err := client.GetStream(video, &video.Formats[2])
	if err != nil {
		panic(err)
		return "", err
	}
	defer resp.Body.Close()

	file, err := os.Create(fmt.Sprintf("tmp/%s.mp4", video.ID))
	path = fmt.Sprintf("tmp/%s.mp4", video.ID)
	if err != nil {
		panic(err)
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
		return "", err
	}

	return path, nil
}
