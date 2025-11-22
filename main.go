package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type Bookmark struct {
	Name string   `json:"name"`
	Link string   `json:"link"`
	Tags []string `json:"tags"`
}

func getDefaultBookmarskFilePath() (string, error) {
	// home, err := os.UserHomeDir()
	// if err != nil {
	// 	return "", err
	// }

	// return path.Join(home, ".bookmarks.json"), nil
	return path.Join(".", "test.json"), nil
}

func store(bookmarks []Bookmark) error {
	bmJSON, err := json.Marshal(bookmarks)
	if err != nil {
		return err
	}

	path := filepath.Join(".", "test.json")
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	// dont need to know how much byte written
	_, err = f.Write(bmJSON)
	return err
}

func listBookmarks(bookmarks []Bookmark) {
	for i, bm := range bookmarks {
		fmt.Println(i, bm.Name, bm.Link)
	}
}

func load(bookmarks *[]Bookmark) error {
	path, err := getDefaultBookmarskFilePath()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &bookmarks)
	if err != nil {
		return err
	}

	return nil
}

func addBookmark(bookmarks *[]Bookmark, bookmark Bookmark) {
	*bookmarks = append(*bookmarks, bookmark)
}

func deleteBookmark(bookmarks *[]Bookmark, idx int) {
	*bookmarks = append((*bookmarks)[:idx], (*bookmarks)[idx+1:]...)
}

func buildBookmark() Bookmark {
	name := os.Args[2]
	link := os.Args[3]
	tags := make([]string, 0)

	// got tags
	if len(os.Args) > 3 {
		tags = append(tags, os.Args[3:]...)
	}

	return Bookmark{
		Name: name,
		Link: link,
		Tags: tags,
	}
}

func main() {
	bookmarks := make([]Bookmark, 0)
	err := load(&bookmarks)
	if err != nil {
		fmt.Println(err)
	}

	switch os.Args[1] {
	case "add":
		addBookmark(&bookmarks, buildBookmark())
	case "list":
		listBookmarks(bookmarks)
	case "delete":
		idx, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		deleteBookmark(bookmarks, idx)
	}

	err = store(bookmarks)
	if err != nil {
		fmt.Println(err)
	}
}
