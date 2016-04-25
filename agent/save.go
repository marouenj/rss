package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/marouenj/rss/util"
)

type Owner struct {
	Id       string    `json:"id"`
	Channels *Channels `json:"channels"`
}

type Owners []*Owner

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

// organizes the crawler channel-centric data into date-centric data
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

func (m *Marshaller) Save() error {
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

func (m *Marshaller) merge() error {
	return nil
}

func (m *Marshaller) clean() error {
	return nil
}
