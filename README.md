# bit - Bookmarks in terminal
`bit` is a simple terminal based bookmarks manager. Simple because `bit` is storing the bookmarks in a plain JSON file, located at `~/.bookmarks.json`.

## Usage
```shell
# add a bookmark
$ bit add "Example website" https://example.com tag1 tag2
$ bit add "Go website" https://go.dev go programming

# list bookmarks
$ bit list
-> 0 Example website https://example.com [tag1 tag2]
-> 1 Go website https://go.dev [go programming]

# delete a bookmark
$ bit delete 1

# open a bookmark
$ bit open 0
```
