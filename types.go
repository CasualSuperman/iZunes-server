package main

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
)

type Library struct {
	Artists []Artist
	Albums  []Album
	Songs   []Song
}

type Artist struct {
	name string

	singles []*Song
	albums  []*Album
}

type Album struct {
	name    string
	songs   []*Song
	numSize uint8

	artist *Artist
}

type Song struct {
	trackNum uint8
	name     string
	ext      string

	album  *Album
	artist *Artist
}

func NewArtist(name string) Artist {
	return Artist{name, nil, nil}
}

func NewAlbum(name string, artist *Artist) Album {
	return Album{name, nil, 2, artist}
}

func NewNumberedSong(song string) (Song, error) {
	baseSong, err := NewSong(song)
	if err != nil {
		return Song{}, err
	}
	dashIndx := strings.Index(baseSong.name, "-")
	if dashIndx == -1 {
		return Song{}, errors.New("no dash separator")
	}
	songNum := baseSong.name[:dashIndx]
	baseSong.name = baseSong.name[dashIndx+1:]
	trackNum, err := strconv.Atoi(songNum)
	baseSong.trackNum = uint8(trackNum)
	return baseSong, err
}

func NewSong(song string) (Song, error) {
	ext := path.Ext(song)
	if len(ext) == 0 {
		return Song{}, errors.New("song has no extension")
	}
	song = song[:len(song)-len(ext)]
	return Song{0, song, ext[1:], nil, nil}, nil
}

func (l *Library) AddSongByPath(path string) error {
	parts := unsanitize(path)
	var artistName, albumName, songName string
	if len(parts) == 2 {
		artistName = parts[0]
		albumName = ""
		songName = parts[1]
	} else if len(parts) == 3 {
		artistName = parts[0]
		albumName = parts[1]
		songName = parts[2]
	} else {
		return errors.New("can't handle " + strconv.Itoa(len(parts)) + " path segments")
	}

	var artistInstance *Artist
	for i, artist := range l.Artists {
		if artist.name == artistName {
			artistInstance = &l.Artists[i]
			break
		}
	}
	if artistInstance == nil {
		artist := NewArtist(artistName)
		l.Artists = append(l.Artists, artist)
		artistInstance = &l.Artists[len(l.Artists)-1]
	}

	var albumInstance *Album
	if albumName == "" {
		albumInstance = nil
	} else {
		for _, album := range artistInstance.albums {
			if album.name == albumName {
				albumInstance = album
				break
			}
		}
		if albumInstance == nil {
			album := NewAlbum(albumName, artistInstance)
			l.Albums = append(l.Albums, album)
			albumInstance = &l.Albums[len(l.Albums)-1]
			artistInstance.albums = append(artistInstance.albums, albumInstance)
		}
	}

	var songInstance *Song
	var song Song
	var err error

	if albumInstance == nil {
		song, err = NewSong(songName)
	} else {
		song, err = NewNumberedSong(songName)
	}

	if err != nil {
		return err
	}

	song.album = albumInstance
	song.artist = artistInstance

	l.Songs = append(l.Songs, song)
	songInstance = &l.Songs[len(l.Songs)-1]

	if albumInstance == nil {
		artistInstance.singles = append(artistInstance.singles, songInstance)
	} else {
		albumInstance.songs = append(albumInstance.songs, songInstance)
		digits := uint8(len(strconv.Itoa(int(song.trackNum))))

		if digits > albumInstance.numSize {
			albumInstance.numSize = digits
		}
	}
	return nil
}

func (l Library) ListSongs() []string {
	songs := make([]string, len(l.Songs))
	for _, song := range l.Songs {
		songs = append(songs, song.Path())
	}
	return songs
}

func unsanitize(path string) []string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		parts[i] = strings.Replace(strings.Replace(part, "%2F", "/", -1), "%%", "%", -1)
	}
	return parts
}

func sanitize(path string) string {
	return strings.Replace(strings.Replace(path, "%", "%%", -1), "/", "%2F", -1)
}

func (s Song) Path() string {
	path := sanitize(s.artist.name) + "/"
	if s.album != nil {
		path += sanitize(s.album.name) + "/"
		path += fmt.Sprintf("%0"+strconv.Itoa(int(s.album.numSize))+"d", s.trackNum)
		path += "-"
	}
	path += sanitize(s.name)
	path += "." + s.ext
	return path
}
