package agent

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Test_unmarshal(t *testing.T) {
	testCases := []struct {
		json   string
		groups ChannelGroups
	}{
		{
			`
			[
				{
					"owner": "wsj",
					"channels": ["http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7014.xml"]
				}
			]
			`,
			ChannelGroups{
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7085.xml",
						"http://www.wsj.com/xml/rss/3_7014.xml",
					},
				},
			},
		},
		{
			`
			[
				{
					"owner": "wsj",
					"channels": ["http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7014.xml"]
				},
				{
					"owner": "cnet",
					"channels": ["http://www.cnet.com/rss/iphone-update/", "http://www.cnet.com/rss/android-update/"]
				}
			]
			`,
			ChannelGroups{
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7085.xml",
						"http://www.wsj.com/xml/rss/3_7014.xml",
					},
				},
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
						"http://www.cnet.com/rss/iphone-update/",
						"http://www.cnet.com/rss/android-update/",
					},
				},
			},
		},
	}

	for idx, testCase := range testCases {
		// create temp dir
		dir, err := ioutil.TempDir("", "dir")
		if err != nil {
			t.Error(err)
		}

		// write json load to temp file
		file := filepath.Join(dir, "file")
		err = ioutil.WriteFile(file, []byte(testCase.json), 0666)
		if err != nil {
			t.Error(err)
		}

		groups, err := unmarshal(dir, "file")

		if !reflect.DeepEqual(groups, testCase.groups) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.groups, groups)
		}

		os.RemoveAll(dir)
	}
}

func Test_forEachFile(t *testing.T) {
	cases := []struct {
		in     string
		length int
	}{
		{`[]`, 0},
		{`["a", "b", "c"]`, 3},
	}

	for _, c := range cases {
		// create temp dir
		dir, err := ioutil.TempDir("", "dir")
		if err != nil {
			t.Error(err)
		}

		// write json load to temp file
		file := filepath.Join(dir, "file")
		err = ioutil.WriteFile(file, []byte(c.in), 0666)
		if err != nil {
			t.Error(err)
		}

		urls, err := forEachFile(dir, "file")

		// check length
		if urls.Len() != c.length {
			t.Errorf("Length of %v == %d, want %d", c.in, urls.Len(), c.length)
		}

		os.RemoveAll(dir)
	}
}

func Test_Load(t *testing.T) {
	cases := []struct {
		in     []string
		length int
		load   []string
	}{
		{[]string{`[]`}, 0, []string{}},
		{[]string{`["a", "b", "c"]`}, 3, []string{"a", "b", "c"}},
		{[]string{`["c", "b", "a"]`}, 3, []string{"a", "b", "c"}},
		{[]string{`["c", "b", "a"]`, `["f", "e", "d"]`}, 6, []string{"a", "b", "c", "d", "e", "f"}},
		{[]string{`["c", "e", "d"]`, `["a", "e", "b"]`}, 6, []string{"a", "b", "c", "d", "e", "e"}},
	}

	for _, c := range cases {
		// create temp dir
		dir, err := ioutil.TempDir("", "dir")
		if err != nil {
			t.Error(err)
		}

		// write json load to temp file
		for idx, i := range c.in {
			file := filepath.Join(dir, strings.Join([]string{"file", strconv.Itoa(idx), ".json"}, ""))
			err = ioutil.WriteFile(file, []byte(i), 0666)
			if err != nil {
				t.Error(err)
			}
		}

		// init loader
		loader, err := NewLoader()
		if err != nil {
			t.Error(err)
		}

		// under test
		err = loader.Load(dir)
		if err != nil {
			t.Error(err)
		}
		urls := loader.Urls

		// check length
		if len(urls) != c.length {
			t.Errorf("[input: %v] Length is %d, want %d", c.in, len(urls), c.length)
		}

		// check content, order
		for idx, url := range urls {
			if strings.Compare(url, c.load[idx]) != 0 {
				t.Errorf("[input: %v] %dth element is %v, want %v", c.in, idx, url, c.load[idx])
			}
		}

		os.RemoveAll(dir)
	}
}
