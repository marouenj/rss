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
	Title string `xml:"title"       json:"title"`
	Link  string `xml:"link"        json:"link"`
	Desc  string `xml:"description" json:"desc"`
	Date  string `xml:"pubDate"     json:"-"`
}

type Items []*Item

// Channel models a 'channel' in an RSS feed
type Channel struct {
	Owner string `xml:"-"           json:"-"`
	Title string `xml:"title"       json:"title"`
	Desc  string `xml:"description" json:"desc"`
	Items *Items `xml:"item"        json:"items"`
}

type Channels []*Channel

// Rss represents an RSS document
type Rss struct {
	XMLName  xml.Name `xml:"rss"`
	Channels Channels `xml:"channel"`
}

type Crawler struct {
	Rss Rss
}

func NewCrawler() (*Crawler, error) {
	return &Crawler{}, nil
}

func (c *Crawler) Crawl(loader *Loader) error {
	if loader == nil {
		return errors.New(fmt.Sprintf("[ERR] Loader is null"))
	}

	if loader.ChannelGroups == nil {
		return errors.New(fmt.Sprintf("[ERR] Loader not contain links"))
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
				continue
			}

			for _, channel := range rss.Channels {
				channel.Owner = group.Owner
			}

			c.merge(rss.Channels)
		}
	}

	c.clean()

	return nil
}

func (c *Crawler) merge(channels []*Channel) error {
	if channels == nil {
		return errors.New(fmt.Sprintf("[ERR] Unvalid arg 'rss>Channels', %v", channels))
	}

	c.Rss.Channels = append(c.Rss.Channels, channels...)
	return nil
}

func (c *Crawler) clean() error {
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

	return nil
}

func (c *Crawler) print() string {
	s := ""
	for _, channel := range c.Rss.Channels {
		s = strings.Join([]string{s, fmt.Sprintf("Title: @%s@\n", channel.Title)}, "")
		s = strings.Join([]string{s, fmt.Sprintf("Desc:  @%s@\n", channel.Desc)}, "")
		for idx, item := range *channel.Items {
			s = strings.Join([]string{s, fmt.Sprintf("\t%2d Title: @%s@\n", idx, item.Title)}, "")
			s = strings.Join([]string{s, fmt.Sprintf("\t%2d Link:  @%s@\n", idx, item.Link)}, "")
			s = strings.Join([]string{s, fmt.Sprintf("\t%2d Desc:  @%s@\n", idx, item.Desc)}, "")
		}
	}

	return s
}
