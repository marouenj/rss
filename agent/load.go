package agent

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ChannelGroup struct {
	Owner    string   `json:"owner"`
	Channels []string `json:"channels"`
}

type ChannelGroups []ChannelGroup

type Loader struct {
	Urls []string
}

func NewLoader() (*Loader, error) {
	return &Loader{}, nil
}

func (l *Loader) Load(file string) error {
	// open file
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening '%s': %s\n", file, err)
		os.Exit(1)
	}
	defer f.Close()

	// get file info
	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("Error reading stats for '%s': %s\n", file, err)
		os.Exit(1)
	}

	urls := list.New()

	if !fi.IsDir() { // is a file
		l, err := forEachFile("", file)
		if err != nil {
			fmt.Printf("Error reading '%s': %s\n", file, err)
			os.Exit(1)
		}
		urls.PushBackList(l)
	} else { // is a dir
		contents, err := f.Readdir(-1)
		if err != nil {
			fmt.Printf("Error reading '%s': %s\n", file, err)
			os.Exit(1)
		}

		// sort the contents, ensures lexical order
		sort.Sort(dirEntries(contents))

		for _, fi := range contents {
			// don't recursively read contents
			if fi.IsDir() {
				continue
			}

			// if it's not a json file, ignore it
			if !strings.HasSuffix(fi.Name(), ".json") {
				continue
			}

			l, err := forEachFile(file, fi.Name())
			if err != nil {
				fmt.Printf("Error reading '%s': %s\n", fi.Name(), err)
				os.Exit(1)
			}
			urls.PushBackList(l)
		}
	}

	l.Urls = make([]string, urls.Len())
	idx := -1
	for e := urls.Front(); e != nil; e = e.Next() {
		url, _ := e.Value.(string)
		idx++
		l.Urls[idx] = url
	}

	sort.Strings(l.Urls)

	return nil
}

func unmarshal(prefix string, name string) (ChannelGroups, error) {
	in := filepath.Join(prefix, name)
	file, err := ioutil.ReadFile(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading '%s': %s", in, err)
	}

	var channelGroups ChannelGroups
	err = json.Unmarshal(file, &channelGroups)
	if err != nil {
		return nil, fmt.Errorf("Error decoding '%s': %s", in, err)
	}

	return channelGroups, nil
}

func forEachFile(prefix string, name string) (*list.List, error) {
	urls := list.New()

	in := filepath.Join(prefix, name)
	file, err := ioutil.ReadFile(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading '%s': %s", in, err)
	}

	var entries []string
	err = json.Unmarshal(file, &entries)
	if err != nil {
		fmt.Errorf("Error decoding '%s': %s", in, err)
	}

	for _, url := range entries {
		urls.PushBack(url)
	}

	return urls, nil
}

type dirEntries []os.FileInfo

// Implement the sort interface for dirEnts
func (d dirEntries) Len() int {
	return len(d)
}

func (d dirEntries) Less(i, j int) bool {
	return d[i].Name() < d[j].Name()
}

func (d dirEntries) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
