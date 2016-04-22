package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// represent a group of channels grouped by their common owner
type ChannelGroup struct {
	Owner    string   `json:"owner"`
	Channels []string `json:"channels"`
}

type ChannelGroups []ChannelGroup

// implement the sort interface for ChannelGroups
func (cg ChannelGroups) Len() int {
	return len(cg)
}
func (cg ChannelGroups) Less(i, j int) bool {
	return strings.Compare(cg[i].Owner, cg[j].Owner) < 0
}
func (cg ChannelGroups) Swap(i, j int) {
	cg[i], cg[j] = cg[j], cg[i]
}

// init a ChannelGroups from a json file
func NewChannelGroups(dir string, fname string) (*ChannelGroups, error) {
	path := filepath.Join(dir, fname)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[ERR] Unable to read '%s': %v", path, err)
	}

	var channelGroups ChannelGroups
	err = json.Unmarshal(file, &channelGroups)
	if err != nil {
		return nil, fmt.Errorf("[ERR] Unable to decode '%s': %v", path, err)
	}

	return &channelGroups, nil
}

// group by owner
func (cg *ChannelGroups) mergeOwners() error {
	curr := 0
	for idx, _ := range (*cg)[1:] {
		if strings.Compare((*cg)[curr].Owner, (*cg)[idx+1].Owner) == 0 { // merge
			(*cg)[curr].Channels = append((*cg)[curr].Channels, (*cg)[idx+1].Channels...)
		} else {
			curr++
			(*cg)[curr] = (*cg)[idx+1]
		}
	}

	// resize
	t := *cg
	*cg = make(ChannelGroups, curr+1)
	copy(*cg, t)

	return nil
}

// remove duplicate links (scope is within same owner)
func (cg *ChannelGroups) cleanLinks() error {
	for idxG, _ := range *cg {
		curr := 0
		for idx, _ := range (*cg)[idxG].Channels[1:] {
			if strings.Compare((*cg)[idxG].Channels[curr], (*cg)[idxG].Channels[idx+1]) != 0 {
				curr++
				(*cg)[idxG].Channels[curr] = (*cg)[idxG].Channels[idx+1]
			}
		}

		// resize
		t := (*cg)[idxG].Channels
		(*cg)[idxG].Channels = make([]string, curr+1)
		copy((*cg)[idxG].Channels, t)
	}

	return nil
}

// the agent responsible for loading and managing the links to the rss resources
type Loader struct {
	ChannelGroups ChannelGroups
}

func NewLoader() (*Loader, error) {
	return &Loader{}, nil
}

func (l *Loader) Load(file string) error {
	// open file
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("[ERR] Unable to open '%s': %v", file, err)
	}
	defer f.Close()

	// get file info
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("[ERR] Unable to read stats of '%s': %v", file, err)
	}

	if !fi.IsDir() { // is a file
		groups, err := NewChannelGroups("", file)
		if err != nil {
			return err // already formatted
		}
		l.ChannelGroups = *groups
	} else { // is a dir
		contents, err := f.Readdir(-1)
		if err != nil {
			return fmt.Errorf("[ERR] Unable to list dir entries of '%s': %v", file, err)
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

			groups, err := NewChannelGroups(file, fi.Name())
			if err != nil {
				return err // already formatted
			}
			l.ChannelGroups = append(l.ChannelGroups, *groups...)
		}
	}

	// sort owners
	sort.Sort(l.ChannelGroups)

	// similar entries (entries of the same owner) are merged into one entry
	if err := l.ChannelGroups.mergeOwners(); err != nil {
		return fmt.Errorf("[ERR] Unable to merge owners: %v", err)
	}

	// sort channels
	for _, ChannelGroup := range l.ChannelGroups {
		sort.Strings(ChannelGroup.Channels)
	}

	// clean links
	if err := l.ChannelGroups.cleanLinks(); err != nil {
		return fmt.Errorf("[ERR] Unable to clean links: %v", err)
	}

	return nil
}

type dirEntries []os.FileInfo

// Implement the sort interface for dirEntries
func (d dirEntries) Len() int {
	return len(d)
}
func (d dirEntries) Less(i, j int) bool {
	return d[i].Name() < d[j].Name()
}
func (d dirEntries) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
