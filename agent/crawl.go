package agent

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Item models an 'item' of an RSS feed
type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Date  string `xml:"pubDate"`
}

func (i *Item) setTitle(newTitle string) {
	i.Title = newTitle
}

func (i *Item) setLink(newLink string) {
	i.Link = newLink
}

func (i *Item) setDesc(newDesc string) {
	i.Desc = newDesc
}

func (i *Item) setDate(newDate string) {
	i.Date = newDate
}

// Channel models a 'channel' in an RSS feed
type Channel struct {
	Title string  `xml:"title"`
	Desc  string  `xml:"description"`
	Items []*Item `xml:"item"`
}

func (c *Channel) setTitle(newTitle string) {
	c.Title = newTitle
}

func (c *Channel) setDesc(newDesc string) {
	c.Desc = newDesc
}

// Rss represents an RSS document
type Rss struct {
	XMLName  xml.Name   `xml:"rss"`
	Channels []*Channel `xml:"channel"`
}

type Crawler struct {
	Rss Rss
}

func NewCrawler() (*Crawler, error) {
	return &Crawler{}, nil
}

func (c *Crawler) Crawl(loader *Loader) error {
	if loader == nil {
		return nil
	}

	if loader.ChannelGroups == nil {
		return nil
	}

	for _, group := range loader.ChannelGroups {
		for _, url := range group.Channels {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("[ERR] Unable to GET %v: %v", url, err)
				continue
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			var rss Rss
			err = xml.Unmarshal(body, &rss)
			if err != nil {
				fmt.Printf("[ERR] Unable to unmarshal %v: %v", string(body[:]), err)
			}

			c.merge(rss.Channels)
		}
	}

	c.clean()

	return nil
}

func (c *Crawler) merge(channels []*Channel) error {
	if channels == nil {
		return errors.New(fmt.Sprintf("Unvalid arg 'rss>Channels', %v", channels))
	}

	c.Rss.Channels = append(c.Rss.Channels, channels...)
	return nil
}

func (c *Crawler) clean() error {
	for _, channel := range c.Rss.Channels {
		channel.setTitle(strings.TrimSpace(channel.Title))
		channel.setDesc(strings.TrimSpace(channel.Desc))
		for _, item := range channel.Items {
			item.setTitle(strings.TrimSpace(item.Title))
			item.setLink(strings.TrimSpace(item.Link))
			item.setDesc(strings.TrimSpace(item.Desc))
			item.setDate(strings.TrimSpace(item.Date))
		}
	}

	return nil
}

func (c *Crawler) print() string {
	s := ""
	for _, channel := range c.Rss.Channels {
		s = strings.Join([]string{s, fmt.Sprintf("Title: @%s@\n", channel.Title)}, "")
		s = strings.Join([]string{s, fmt.Sprintf("Desc:  @%s@\n", channel.Desc)}, "")
		for idx, item := range channel.Items {
			s = strings.Join([]string{s, fmt.Sprintf("\t%2d Title: @%s@\n", idx, item.Title)}, "")
			s = strings.Join([]string{s, fmt.Sprintf("\t%2d Link:  @%s@\n", idx, item.Link)}, "")
			s = strings.Join([]string{s, fmt.Sprintf("\t%2d Desc:  @%s@\n", idx, item.Desc)}, "")
		}
	}

	return s
}
