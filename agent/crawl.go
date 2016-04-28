package agent

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Item models an 'item' of an RSS feed
type Item struct {
	Title string `xml:"title"       json:"title"`
	Link  string `xml:"link"        json:"link"`
	Desc  string `xml:"description" json:"desc"`
	Date  string `xml:"pubDate"     json:"-"`
}

type Items []*Item

// implement the sort interface for Items
func (it Items) Len() int {
	return len(it)
}
func (it Items) Less(i, j int) bool {
	return strings.Compare(it[i].Title, it[j].Title) < 0
}
func (it Items) Swap(i, j int) {
	it[i], it[j] = it[j], it[i]
}

// Channel models a 'channel' in an RSS feed
type Channel struct {
	Owner string `xml:"-"           json:"-"`
	Title string `xml:"title"       json:"title"`
	Desc  string `xml:"description" json:"desc"`
	Items *Items `xml:"item"        json:"items"`
}

type Channels []*Channel

// implement the sort interface for Channels
func (ch Channels) Len() int {
	return len(ch)
}
func (ch Channels) Less(i, j int) bool {
	return strings.Compare(ch[i].Title, ch[j].Title) < 0
}
func (ch Channels) Swap(i, j int) {
	ch[i], ch[j] = ch[j], ch[i]
}

// Rss represents an RSS document
type Rss struct {
	XMLName  xml.Name `xml:"rss"`
	Channels Channels `xml:"channel"`
}

type Crawler struct {
	Rss Rss
}

func NewCrawler() (*Crawler, error) {
	return &Crawler{
		Rss: Rss{
			Channels: Channels{},
		},
	}, nil
}

func (c *Crawler) Crawl(loader *Loader) error {
	if loader == nil {
		return fmt.Errorf("[ERR] 'loader' is nil")
	}

	if loader.ChannelGroups == nil {
		return fmt.Errorf("[ERR] 'loader->ChannelGroups' is nil")
	}

	for _, group := range loader.ChannelGroups {
		for _, url := range group.Channels {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("[ERR] Unable to GET '%s': %v", url, err)
				continue
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)

			var rss Rss
			err = xml.Unmarshal(body, &rss)
			if err != nil {
				fmt.Printf("[ERR] Unable to unmarshal '%s': %v", string(body[:]), err)
				continue
			}

			for _, channel := range rss.Channels {
				channel.Owner = group.Owner
			}

			err = c.merge(rss.Channels)
			if err != nil {
				fmt.Errorf("[ERR] Unable to merge channels '%v'", rss.Channels)
			}
		}
	}

	c.clean()

	return nil
}

func (c *Crawler) merge(channels []*Channel) error {
	if channels == nil {
		return fmt.Errorf("[ERR] Unvalid arg 'rss>Channels', %v", channels)
	}

	c.Rss.Channels = append(c.Rss.Channels, channels...)
	return nil
}

func (c *Crawler) clean() {
	for _, channel := range c.Rss.Channels {
		channel.Title = strings.TrimSpace(channel.Title)
		channel.Desc = strings.TrimSpace(channel.Desc)
		for _, item := range *channel.Items {
			item.Title = strings.TrimSpace(item.Title)
			item.Link = strings.TrimSpace(item.Link)
			item.Desc = strings.TrimSpace(item.Desc)
			item.Date = strings.TrimSpace(item.Date)
		}
	}
}
