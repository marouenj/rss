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

func Test_NewChannelGroups(t *testing.T) {
	testCases := []struct {
		json   string
		groups ChannelGroups
	}{
		{ // test case 0, single channel
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
		{ // test case 1, multiple channels
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
		{ // test case 2, multiple channels, duplicated
			`
			[
				{
					"owner": "wsj",
					"channels": ["http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7014.xml"]
				},
				{
					"owner": "cnet",
					"channels": ["http://www.cnet.com/rss/iphone-update/"]
				},
				{
					"owner": "cnet",
					"channels": ["http://www.cnet.com/rss/android-update/"]
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
					},
				},
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
						"http://www.cnet.com/rss/android-update/",
					},
				},
			},
		},
		{ // test case 3, multiple channels, duplicated
			`
			[
				{
					"owner": "wsj",
					"channels": ["http://www.wsj.com/xml/rss/3_7085.xml"]
				},
				{
					"owner": "cnet",
					"channels": ["http://www.cnet.com/rss/iphone-update/"]
				},
				{
					"owner": "wsj",
					"channels": ["http://www.wsj.com/xml/rss/3_7014.xml"]
				},
				{
					"owner": "cnet",
					"channels": ["http://www.cnet.com/rss/android-update/"]
				}
			]
			`,
			ChannelGroups{
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7085.xml",
					},
				},
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
						"http://www.cnet.com/rss/iphone-update/",
					},
				},
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7014.xml",
					},
				},
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
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

		// under test
		groups, err := NewChannelGroups(dir, "file")

		// assert
		if !reflect.DeepEqual(*groups, testCase.groups) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.groups, groups)
		}

		// clean
		os.RemoveAll(dir)
	}
}

func Test_Load(t *testing.T) {
	testCases := []struct {
		json   []string
		groups ChannelGroups
	}{
		{ // test case 0, one file, one owner
			[]string{
				`
				[
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7014.xml"]
					}
				]
				`,
			},
			ChannelGroups{
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7014.xml",
						"http://www.wsj.com/xml/rss/3_7085.xml",
					},
				},
			},
		},
		{ // test case 1, one file, multiple owners
			[]string{
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
			},
			ChannelGroups{
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
						"http://www.cnet.com/rss/android-update/",
						"http://www.cnet.com/rss/iphone-update/",
					},
				},
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7014.xml",
						"http://www.wsj.com/xml/rss/3_7085.xml",
					},
				},
			},
		},
		{ // test case 2, one file, multiple owners appearing many times
			[]string{
				`
				[
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7014.xml"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/iphone-update/"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/android-update/"]
					}
				]
				`,
			},
			ChannelGroups{
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
						"http://www.cnet.com/rss/android-update/",
						"http://www.cnet.com/rss/iphone-update/",
					},
				},
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7014.xml",
						"http://www.wsj.com/xml/rss/3_7085.xml",
					},
				},
			},
		},
		{ // test case 3, one file, multiple owners appearing many times
			[]string{
				`
				[
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7085.xml"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/iphone-update/"]
					},
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7014.xml"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/android-update/"]
					}
				]
				`,
			},
			ChannelGroups{
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
						"http://www.cnet.com/rss/android-update/",
						"http://www.cnet.com/rss/iphone-update/",
					},
				},
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7014.xml",
						"http://www.wsj.com/xml/rss/3_7085.xml",
					},
				},
			},
		},
		{ // test case 4, multiple files, multiple owners appearing many times
			[]string{
				`
				[
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7085.xml"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/iphone-update/"]
					}
				]
				`,
				`
				[
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7014.xml"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/android-update/"]
					}
				]
				`,
			},
			ChannelGroups{
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
						"http://www.cnet.com/rss/android-update/",
						"http://www.cnet.com/rss/iphone-update/",
					},
				},
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7014.xml",
						"http://www.wsj.com/xml/rss/3_7085.xml",
					},
				},
			},
		},
		{ // test case 5, one file, one owner, duplicate links
			[]string{
				`
				[
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7014.xml", "http://www.wsj.com/xml/rss/3_7085.xml"]
					}
				]
				`,
			},
			ChannelGroups{
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7014.xml",
						"http://www.wsj.com/xml/rss/3_7085.xml",
					},
				},
			},
		},
		{ // test case 6, one file, multiple owners appearing many times, duplicate links
			[]string{
				`
				[
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7014.xml", "http://www.wsj.com/xml/rss/3_7085.xml"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/iphone-update/"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/iphone-update/", "http://www.cnet.com/rss/android-update/"]
					}
				]
				`,
			},
			ChannelGroups{
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
						"http://www.cnet.com/rss/android-update/",
						"http://www.cnet.com/rss/iphone-update/",
					},
				},
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7014.xml",
						"http://www.wsj.com/xml/rss/3_7085.xml",
					},
				},
			},
		},
		{ // test case 7, multiple files, multiple owners appearing many times, duplicate links
			[]string{
				`
				[
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7085.xml", "http://www.wsj.com/xml/rss/3_7014.xml"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/iphone-update/"]
					}
				]
				`,
				`
				[
					{
						"owner": "wsj",
						"channels": ["http://www.wsj.com/xml/rss/3_7014.xml"]
					},
					{
						"owner": "cnet",
						"channels": ["http://www.cnet.com/rss/android-update/", "http://www.cnet.com/rss/iphone-update/"]
					}
				]
				`,
			},
			ChannelGroups{
				ChannelGroup{
					Owner: "cnet",
					Channels: []string{
						"http://www.cnet.com/rss/android-update/",
						"http://www.cnet.com/rss/iphone-update/",
					},
				},
				ChannelGroup{
					Owner: "wsj",
					Channels: []string{
						"http://www.wsj.com/xml/rss/3_7014.xml",
						"http://www.wsj.com/xml/rss/3_7085.xml",
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
		for idx, i := range testCase.json {
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

		// assert
		if !reflect.DeepEqual(loader.ChannelGroups, testCase.groups) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.groups, loader.ChannelGroups)
		}

		os.RemoveAll(dir)
	}
}
