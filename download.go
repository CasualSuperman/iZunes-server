package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/*
func main() {
	Download("https://www.youtube.com/watch?v=e-ORhEE9VVg", 3 * time.Second, time.Duration(3 * time.Minute + 54 * time.Second), 0, "Shake it Off", "1989", "Taylor Swift")
}
*/

// Download a youtube video using offLiberty.
func Download(downloadURL string, start, end time.Duration, trackNum int, trackName, albumName, artistName string) {
	resp, err := http.PostForm("http://offliberty.com/off02.php", url.Values{"track": {downloadURL}})
	if err != nil {
		println("offliberty failed")
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("readAll failed")
		return
	}

	body := string(data)
	i := strings.Index(body, "<A HREF=\"")
	if i == -1 {
		println("link start failed")
		return
	}

	body = body[i+len("<A HREF=\""):]

	i = strings.Index(body, "\"")
	if i == -1 {
		println("link end failed")
		return
	}

	downloadURL = body[:i]
	println(downloadURL)

	resp, err = http.Get(downloadURL)
	if err != nil {
		println("download failed")
		return
	}
	defer resp.Body.Close()

	f, err := ioutil.TempFile("", "iZunes-download-")
	if err != nil {
		println("temp file failed")
		return
	}
	defer os.Remove(f.Name())

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		println("copy failed")
		return
	}

	songFile := f.Name()
	f.Close()

	startSecs := strconv.FormatFloat(start.Seconds(), 'f', -1, 64)
	durationSecs := strconv.FormatFloat(end.Seconds()-start.Seconds(), 'f', -1, 64)
	cmd := exec.Command("ffmpeg", "-i", songFile, "-f", "mp3", "-acodec", "copy", "-ss", startSecs, "-t", durationSecs, songFile+".out")

	err = cmd.Run()
	if err != nil {
		println("convert failed")
		println(err.Error())
		return
	}

	println(songFile + ".out")
}
