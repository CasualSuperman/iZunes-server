package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const musicDir = "/bulk/music/"

type songInfo struct {
	path string
	mod  time.Time
}

func main() {
	l := Library{}
	songList, err := getSongList()
	if err != nil {
		panic(err.Error())
	}
	for _, s := range songList {
		err := l.AddSongByPath(s)
		if err != nil {
			panic(err.Error())
		}
	}
	for _, p := range l.ListSongs() {
		fmt.Println(p)
	}
}

func getSongList() ([]string, error) {
	var songInfos []songInfo

	err := filepath.Walk(musicDir, func(folderPath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			songPath, err := filepath.Rel(musicDir, folderPath)
			if err != nil {
				return err
			}

			switch strings.ToLower(path.Ext(songPath)) {
			case ".jpg":
			case ".png":
			case ".gif":
			case ".bmp":
			default:
				songInfos = append(songInfos, songInfo{
					songPath,
					info.ModTime(),
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Sort(songList(songInfos))
	songs := make([]string, len(songInfos))
	for i, song := range songInfos {
		songs[i] = song.path
	}

	return songs, nil
}

type songList []songInfo

func (s songList) Len() int {
	return len(s)
}

func (s songList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s songList) Less(i, j int) bool {
	return s[i].mod.Before(s[j].mod)
}
