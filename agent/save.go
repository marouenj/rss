package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/marouenj/rss/util"
)

type Owner struct {
	Id       string    `json:"id"`
	Channels *Channels `json:"channels"`
}

type Owners []*Owner

// implement the sort interface for Owners
func (ow Owners) Len() int {
	return len(ow)
}
func (ow Owners) Less(i, j int) bool {
	return strings.Compare(ow[i].Id, ow[j].Id) < 0
}
func (ow Owners) Swap(i, j int) {
	ow[i], ow[j] = ow[j], ow[i]
}

type Day struct {
	Date   string  `json:"date"`
	Owners *Owners `json:"owners"`
}

type Days []*Day

func (d *Days) AddItem(item Item, date string, ownerId string, channelTitle string, channelDesc string) error {
	// optimistic search for the date
	idxDate := -1
	for idx, day := range *d {
		if strings.Compare(day.Date, date) == 0 {
			idxDate = idx
			break
		}
	}

	// check if the date exist, create a new one otherwise
	selectedDay := &Day{
		Date:   date,
		Owners: &Owners{},
	}
	if idxDate == -1 {
		(*d) = append(*d, selectedDay)
	} else {
		selectedDay = (*d)[idxDate]
	}

	// optimistic search for the owner
	owners := selectedDay.Owners
	idxOwner := -1
	for idx, owner := range *owners {
		if strings.Compare(owner.Id, ownerId) == 0 {
			idxOwner = idx
			break
		}
	}

	// check if the owner exist, create a new one otherwise
	selectedOwner := &Owner{
		Id:       ownerId,
		Channels: &Channels{},
	}
	if idxOwner == -1 {
		(*owners) = append(*owners, selectedOwner)
	} else {
		selectedOwner = (*owners)[idxOwner]
	}

	// optimistic search for the channel
	channels := selectedOwner.Channels
	idxChannel := -1
	for idx, channel := range *channels {
		if strings.Compare(channel.Title, channelTitle) == 0 {
			idxChannel = idx
			break
		}
	}

	// check if the channel exist, create a new one otherwise
	selectedChannel := &Channel{
		Title: channelTitle,
		Desc:  channelDesc,
		Items: &Items{},
	}
	if idxChannel == -1 {
		(*channels) = append(*channels, selectedChannel)
	} else {
		selectedChannel = (*channels)[idxChannel]
	}

	// optimistic search for the item
	items := selectedChannel.Items
	idxItem := -1
	for idx, i := range *items {
		if strings.Compare(i.Title, item.Title) == 0 {
			idxItem = idx
			break
		}
	}

	// check if the channel exist, add it to the list otherwise
	if idxItem == -1 {
		(*items) = append(*items, &item)
	}

	return nil
}

// agent that's responsible for merging new feeds with existing ones
// then persisting them back to disk
type Marshaller struct {
	Days *Days
	dir  string // dir to load from/save to
}

// init a new agent
func NewMarshaller(dir string) (*Marshaller, error) {
	return &Marshaller{
		Days: &Days{},
		dir:  dir,
	}, nil
}

// organizes the crawler channel-centric data into the marshaller date-centric data
func (m *Marshaller) ReArrange(channels Channels) error {
	if channels == nil {
		return errors.New(fmt.Sprintf("[ERR] Argument is nil"))
	}

	for _, channel := range channels {
		for _, item := range *channel.Items {
			date, err := util.ParsePubDate(item.Date)
			if err != nil {
				continue
			}

			m.Days.AddItem(*item, util.DateInUtc(date), channel.Owner, channel.Title, channel.Desc)
		}
	}

	return nil
}

// for each day...
// previous data is loaded from disk, if exists
// current data is then merged with previous data
// the whole is persisted back to disk
// merging operation insures no duplicates in 'owner', 'channel' and 'item' levels
// cleaning operation insures entries are sorted by 'owner', 'channel' and 'item'
func (m *Marshaller) Save() error {
	for _, src := range *m.Days {
		dest, err := m.load(src.Date)
		if err != nil {
			return err // already formatted
		}

		merge(*src, *dest)
		clean(*dest)

		// persist back to disk
		bytes, err := json.Marshal(*dest)
		if err != nil {
			return errors.New(fmt.Sprintf("[ERR] Unable to marshal: %v", err))
		}

		path := filepath.Join(m.dir, src.Date)
		err = ioutil.WriteFile(path, bytes, 0666)
		if err != nil {
			return errors.New(fmt.Sprintf("[ERR] Unable to write to '%s': %v", path, err))
		}
	}

	return nil
}

func (m *Marshaller) load(date string) (*Day, error) {
	path := filepath.Join(m.dir, date)

	// file for this date hasn't been initialized yet
	if _, err := os.Stat(path); err != nil {
		return &Day{
			Date:   date,
			Owners: &Owners{},
		}, nil
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[ERR] Unable to read '%s': %v", path, err)
	}

	var day Day
	err = json.Unmarshal(file, &day)
	if err != nil {
		return nil, fmt.Errorf("[ERR] Unable to decode '%s': %v", path, err)
	}

	return &day, nil
}

func merge(src, dest Day) error {
	mergeOwners(src.Owners, dest.Owners)

	return nil
}

func mergeOwners(src, dest *Owners) error {
	for _, ownerSrc := range *src {
		// optimistic search for the owner
		idxOwner := -1
		for idx, ownerDest := range *dest {
			if strings.Compare(ownerSrc.Id, ownerDest.Id) == 0 {
				idxOwner = idx
				break
			}
		}

		// check if the owner exist, append it otherwise
		if idxOwner == -1 {
			*dest = append(*dest, ownerSrc)
		} else {
			mergeChannels(ownerSrc.Channels, (*dest)[idxOwner].Channels)
		}
	}

	return nil
}

func mergeChannels(src, dest *Channels) error {
	for _, channelSrc := range *src {
		// optimistic search for the channel
		idxChannel := -1
		for idx, channelDest := range *dest {
			if strings.Compare(channelSrc.Title, channelDest.Title) == 0 {
				idxChannel = idx
				break
			}
		}

		// check if the channel exist, append otherwise
		if idxChannel == -1 {
			*dest = append(*dest, channelSrc)
		} else {
			mergeItems(channelSrc.Items, (*dest)[idxChannel].Items)
		}
	}

	return nil
}

func mergeItems(src, dest *Items) error {
	for _, itemSrc := range *src {
		// optimistic search for the item
		idxItem := -1
		for idx, itemDest := range *dest {
			if strings.Compare(itemSrc.Title, itemDest.Title) == 0 {
				idxItem = idx
				break
			}
		}

		// check if the item exist, append otherwise
		if idxItem == -1 {
			*dest = append(*dest, itemSrc)
		}
	}

	// *dest = append(*dest, *src...)
	return nil
}

func clean(day Day) error {
	// sort owners
	sort.Sort(*day.Owners)

	// sort channels for each owner
	for _, owner := range *day.Owners {
		sort.Sort(*owner.Channels)
		for _, channel := range *owner.Channels {
			sort.Sort(*channel.Items)
		}
	}

	return nil
}
