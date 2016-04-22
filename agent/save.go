package agent

import (
	"strings"
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

	// check if the channel exist, create a new one otherwise
	if idxItem == -1 {
		(*items) = append(*items, &item)
	}

	return nil
}
