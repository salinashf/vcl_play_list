package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var extList = []string{".mp4", ".mkv", ".avi", ".flv", ".mov", ".wmv", ".vob", ".mpg", ".3gp", ".m4v"}
var checkSubdirectories = true

type FileFound struct {
	Name   string
	Path   string
	Tracks []string
}

type Playlist struct {
	XMLName   xml.Name `xml:"playlist"`
	Xmlns     string   `xml:"xmlns,attr"`
	Vlc       string   `xml:"xmlns:vlc,attr"`
	Version   string   `xml:"version,attr"`
	Title     string   `xml:"title"`
	TrackList struct {
		Tracks []Track `xml:"track"`
	} `xml:"trackList"`
}

type Track struct {
	Location  string `xml:"location"`
	Title     string `xml:"title"`
	Album     string `xml:"album"`
	TrackNum  string `xml:"trackNum"`
	Extension struct {
		Application string `xml:"application,attr"`
		VlcID       string `xml:"vlc:id"`
	} `xml:"extension"`
}

func NewPlaylist(nameList string) *Playlist {
	return &Playlist{
		Xmlns:   "http://xspf.org/ns/0/",
		Vlc:     "http://www.videolan.org/vlc/playlist/ns/0/",
		Version: "1",
		Title:   nameList,
	}
}

func (p *Playlist) AddTrack(pathFolder *FileFound, seqno int) int {
	nameMaster := pathFolder.Name
	for _, tk := range pathFolder.Tracks {
		seqno++
		track := Track{
			Location: fmt.Sprintf("file:///%s", strings.ReplaceAll(tk, "\\", "/")),
			Title:    filepath.Base(tk),
			Album:    nameMaster,
			TrackNum: fmt.Sprintf("%d", seqno),
			Extension: struct {
				Application string `xml:"application,attr"`
				VlcID       string `xml:"vlc:id"`
			}{
				Application: "http://www.videolan.org/vlc/playlist/0",
				VlcID:       fmt.Sprintf("%d", seqno),
			},
		}
		p.TrackList.Tracks = append(p.TrackList.Tracks, track)
	}
	return seqno
}

func (p *Playlist) GetPlaylist() *Playlist {
	return p
}

type Videos struct{}

func (v *Videos) RemoveNonVideoFiles(fileList []string) []string {
	var filteredList []string
	for _, fileName := range fileList {
		for _, ext := range extList {
			if strings.HasSuffix(fileName, ext) || strings.HasSuffix(fileName, strings.ToUpper(ext)) {
				filteredList = append(filteredList, fileName)
				break
			}
		}
	}
	return filteredList
}

func (v *Videos) GetVideos(pathFinder string) []*FileFound {
	var listFolder []*FileFound
	if checkSubdirectories {
		listFolder = append(listFolder, &FileFound{
			Name: filepath.Base(pathFinder),
			Path: pathFinder,
		})
		filepath.Walk(pathFinder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && !strings.HasPrefix(filepath.Base(path), ".") {
				listFolder = append(listFolder, &FileFound{
					Name: filepath.Base(path),
					Path: path,
				})
			}
			return nil
		})
		for _, folder := range listFolder {
			files, _ := ioutil.ReadDir(folder.Path)
			for _, f := range files {
				if !f.IsDir() {
					folder.Tracks = append(folder.Tracks, filepath.Join(folder.Path, f.Name()))
				}
			}
			folder.Tracks = v.RemoveNonVideoFiles(folder.Tracks)
		}
	} else {
		listFolder = append(listFolder, &FileFound{
			Name: filepath.Base(pathFinder),
			Path: pathFinder,
		})
		files, _ := ioutil.ReadDir(pathFinder)
		for _, f := range files {
			if !f.IsDir() {
				listFolder[0].Tracks = append(listFolder[0].Tracks, filepath.Join(pathFinder, f.Name()))
			}
		}
		listFolder[0].Tracks = v.RemoveNonVideoFiles(listFolder[0].Tracks)
	}
	return listFolder
}

func main() {
	args := os.Args[1:]
	pathFinder := ""
	if len(args) == 2 && args[0] == "-path" {
		pathFinder = args[1]
	} else if len(args) == 1 {
		pathFinder = args[0]
	} else if len(args) == 0 {
		pathFinder = "."
	} else {
		fmt.Println("Invalid arguments.")
		return
	}

	if pathFinder == "." {
		pathFinder, _ = os.Getwd()
	}

	namePlayList := "playList"
	nameList := filepath.Base(pathFinder)
	playlist := NewPlaylist(nameList)
	videos := &Videos{}
	videoFiles := videos.GetVideos(pathFinder)
	seqno := 0
	for _, path := range videoFiles {
		seqno = playlist.AddTrack(path, seqno)
	}

	playlistXML, err := xml.MarshalIndent(playlist.GetPlaylist(), "", "  ")
	if err != nil {
		fmt.Println("Error marshalling XML:", err)
		return
	}

	pathPlayList := filepath.Join(pathFinder, fmt.Sprintf("%s.xspf", namePlayList))
	err = ioutil.WriteFile(pathPlayList, playlistXML, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}
