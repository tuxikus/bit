package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
)

type Bookmark struct {
	Name string   `json:"name"`
	Link string   `json:"link"`
	Tags []string `json:"tags"`
}

func printHelp() {
	helpText := `usage: 
	$ bit <command> <args>

commands:
	add <name> <link> <tags> - add a bookmark
	list - list all bookmarks
	delete <index> - delete a bookmark
	open <index> - open a bookmark

	help - print this message
	`

	fmt.Println(helpText)
}

func getDefaultBookmarksFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".bookmarks.json"), nil
}

func createDefaultBookmarksFile() error {
	path, err := getDefaultBookmarksFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err = os.Create(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func store(bookmarks []Bookmark) error {
	bmJSON, err := json.Marshal(bookmarks)
	if err != nil {
		return err
	}

	path, err := getDefaultBookmarksFilePath()
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(bmJSON)
	return err
}

func getOpenCmd() string {
	cmd := ""

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	return cmd
}

func openBookmark(bookmarks []Bookmark, idx int) error {
	cmd := getOpenCmd()
	link := bookmarks[idx].Link

	err := exec.Command(cmd, link).Start()
	if err != nil {
		return err
	}

	return nil
}

func listBookmarks(bookmarks []Bookmark) {
	for i, bm := range bookmarks {
		fmt.Println(i, bm.Name, bm.Link, bm.Tags)
	}
}

func load() ([]Bookmark, error) {
	path, err := getDefaultBookmarksFilePath()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Bookmark{}, nil
		}
		return nil, err
	}

	if len(content) == 0 || string(content) == "" {
		return []Bookmark{}, nil
	}

	var bookmarks []Bookmark
	err = json.Unmarshal(content, &bookmarks)
	if err != nil {
		return nil, err
	}

	return bookmarks, nil
}

func addBookmark(bookmarks []Bookmark, bookmark Bookmark) []Bookmark {
	bookmarks = append(bookmarks, bookmark)
	return bookmarks
}

func deleteBookmark(bookmarks []Bookmark, idx int) []Bookmark {
	bookmarks = append(bookmarks[:idx], bookmarks[idx+1:]...)
	return bookmarks
}

func buildBookmark() Bookmark {
	name := os.Args[2]
	link := os.Args[3]
	tags := make([]string, 0)

	// got tags
	if len(os.Args) > 3 {
		tags = append(tags, os.Args[4:]...)
	}

	return Bookmark{
		Name: name,
		Link: link,
		Tags: tags,
	}
}

func main() {
	err := createDefaultBookmarksFile()
	if err != nil {
		fmt.Println(
			errors.Join(errors.New("[ERROR] unable to create bookmarks file"), err),
		)
	}

	bookmarks, err := load()
	if err != nil {
		fmt.Println(
			errors.Join(errors.New("[ERROR] unable to load bookmarks from file"), err),
		)
		bookmarks = []Bookmark{}
	}

	if len(os.Args) < 2 {
		fmt.Println("incorrect argument count")
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		bookmarks = addBookmark(bookmarks, buildBookmark())
	case "list":
		listBookmarks(bookmarks)
	case "delete":
		idx, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(
				errors.Join(errors.New("[ERROR] unable to parse index to int"), err),
			)
			fmt.Println(err)
			os.Exit(-1)
		}
		bookmarks = deleteBookmark(bookmarks, idx)
	case "open":
		idx, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(
				errors.Join(errors.New("[ERROR] unable to parse index to int"), err),
			)
			os.Exit(-1)
		}

		err = openBookmark(bookmarks, idx)
		if err != nil {
			fmt.Println(
				errors.Join(errors.New("[ERROR] unable to open bookmark"), err),
			)
			os.Exit(-1)
		}
	default:
		printHelp()
	}

	err = store(bookmarks)
	if err != nil {
		fmt.Println(
			errors.Join(errors.New("[ERROR] unable to store bookmarks"), err),
		)
	}
}
