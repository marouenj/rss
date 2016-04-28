package agent

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func Test_AddItem(t *testing.T) {
	item := Item{
		Title: "Apple iPhone SE owners bemoan audio bug - CNET",
		Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
		Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
	}

	item2 := Item{
		Title: "9 settings every new iPhone owner should change - CNET",
		Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
		Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
	}

	testCases := []struct {
		item         Item // in
		date         string
		owner        string
		channelTitle string
		channelDesc  string
		before       *Days
		after        Days // out
	}{
		{ // test case 0, date not exists (empty)
			item,
			"2016-04-22",
			"cnet",
			"CNET iPhone Update",
			"Tips, news, how tos, and troubleshooting help for the iPhone.",
			&Days{},
			Days{
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{&item},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 1, date not exists
			item,
			"2016-04-22",
			"cnet",
			"CNET iPhone Update",
			"Tips, news, how tos, and troubleshooting help for the iPhone.",
			&Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "wsj",
							Channels: &Channels{
								&Channel{
									Title: "WSJ.com: World News",
									Desc:  "World News",
									Items: &Items{},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "wsj",
							Channels: &Channels{
								&Channel{
									Title: "WSJ.com: World News",
									Desc:  "World News",
									Items: &Items{},
								},
							},
						},
					},
				},
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{&item},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 2, date exists, owner not exists
			item,
			"2016-04-22",
			"cnet",
			"CNET iPhone Update",
			"Tips, news, how tos, and troubleshooting help for the iPhone.",
			&Days{
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "wsj",
							Channels: &Channels{
								&Channel{
									Title: "WSJ.com: World News",
									Desc:  "World News",
									Items: &Items{},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "wsj",
							Channels: &Channels{
								&Channel{
									Title: "WSJ.com: World News",
									Desc:  "World News",
									Items: &Items{},
								},
							},
						},
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{&item},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 3, date exists, owner exists, channel not exists
			item,
			"2016-04-22",
			"cnet",
			"CNET iPhone Update",
			"Tips, news, how tos, and troubleshooting help for the iPhone.",
			&Days{
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET Gaming",
									Desc:  "Game on! Get the latest in gaming news, video game reviews, computer games & video game consoles.",
									Items: &Items{},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET Gaming",
									Desc:  "Game on! Get the latest in gaming news, video game reviews, computer games & video game consoles.",
									Items: &Items{},
								},
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{&item},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 4, date exists, owner exists, channel exists, item not exists
			item,
			"2016-04-22",
			"cnet",
			"CNET iPhone Update",
			"Tips, news, how tos, and troubleshooting help for the iPhone.",
			&Days{
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{&item2},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{&item2, &item},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 5, date exists, owner exists, channel exists, item exists
			item,
			"2016-04-22",
			"cnet",
			"CNET iPhone Update",
			"Tips, news, how tos, and troubleshooting help for the iPhone.",
			&Days{
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{&item},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-22",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{&item},
								},
							},
						},
					},
				},
			},
		},
	}

	for idx, testCase := range testCases {
		days := testCase.before
		err := days.AddItem(testCase.item, testCase.date, testCase.owner, testCase.channelTitle, testCase.channelDesc)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(*days, testCase.after) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.after, *days)
		}
	}
}

func Test_ReArrange(t *testing.T) {
	testCases := []struct {
		channels Channels // in
		days     Days
	}{
		{ // test case 0, one date, one owner, one channel, one item
			Channels{
				&Channel{
					Owner: "cnet",
					Title: "CNET iPhone Update",
					Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
					Items: &Items{
						&Item{
							Title: "Apple iPhone SE owners bemoan audio bug - CNET",
							Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
							Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
							Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-19",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "Apple iPhone SE owners bemoan audio bug - CNET",
											Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
											Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
											Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 1, one date, one owner, one channel, two items
			Channels{
				&Channel{
					Owner: "cnet",
					Title: "CNET iPhone Update",
					Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
					Items: &Items{
						&Item{
							Title: "Apple iPhone SE owners bemoan audio bug - CNET",
							Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
							Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
							Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
						},
						&Item{
							Title: "9 settings every new iPhone owner should change - CNET",
							Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
							Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
							Date:  "Tue, 19 Apr 2016 20:25:18 +0000",
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-19",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "Apple iPhone SE owners bemoan audio bug - CNET",
											Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
											Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
											Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
										},
										&Item{
											Title: "9 settings every new iPhone owner should change - CNET",
											Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
											Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
											Date:  "Tue, 19 Apr 2016 20:25:18 +0000",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 2, one date, one owner, two channels
			Channels{
				&Channel{
					Owner: "cnet",
					Title: "CNET iPhone Update",
					Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
					Items: &Items{
						&Item{
							Title: "Apple iPhone SE owners bemoan audio bug - CNET",
							Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
							Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
							Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
						},
					},
				},
				&Channel{
					Owner: "cnet",
					Title: "CNET Android Update",
					Desc:  "News, analysis and tips on the Google Android operating system, and devices and apps that use it.",
					Items: &Items{
						&Item{
							Title: "Google Play Music adds podcasts to the mix - CNET",
							Link:  "http://www.cnet.com/news/google-play-music-now-does-podcasts-too/#ftag=CADe34d7bf",
							Desc:  "The streaming service will offer up podcasts based on what users are doing or interested in, similar to its contextual playlists for music.",
							Date:  "Tue, 19 Apr 2016 13:25:18 +0000",
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-19",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "Apple iPhone SE owners bemoan audio bug - CNET",
											Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
											Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
											Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
										},
									},
								},
								&Channel{
									Title: "CNET Android Update",
									Desc:  "News, analysis and tips on the Google Android operating system, and devices and apps that use it.",
									Items: &Items{
										&Item{
											Title: "Google Play Music adds podcasts to the mix - CNET",
											Link:  "http://www.cnet.com/news/google-play-music-now-does-podcasts-too/#ftag=CADe34d7bf",
											Desc:  "The streaming service will offer up podcasts based on what users are doing or interested in, similar to its contextual playlists for music.",
											Date:  "Tue, 19 Apr 2016 13:25:18 +0000",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 3, one date, two owners
			Channels{
				&Channel{
					Owner: "cnet",
					Title: "CNET iPhone Update",
					Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
					Items: &Items{
						&Item{
							Title: "Apple iPhone SE owners bemoan audio bug - CNET",
							Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
							Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
							Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
						},
					},
				},
				&Channel{
					Owner: "wsj",
					Title: "WSJ.com: World News",
					Desc:  "World News",
					Items: &Items{
						&Item{
							Title: "Death Toll Rises Following Ecuador Earthquake",
							Link:  "http://www.wsj.com/articles/death-toll-in-ecuador-earthquake-climbs-as-correa-tours-ravaged-areas-1460993084?mod=fox_australian",
							Desc:  "The death toll in the magnitude-7.8 earthquake that struck this small country’s coast rose to 413, officials said.",
							Date:  "Tue, 19 Apr 2016 20:25:18 +0000",
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-19",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "Apple iPhone SE owners bemoan audio bug - CNET",
											Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
											Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
											Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
										},
									},
								},
							},
						},
						&Owner{
							Id: "wsj",
							Channels: &Channels{
								&Channel{
									Title: "WSJ.com: World News",
									Desc:  "World News",
									Items: &Items{
										&Item{
											Title: "Death Toll Rises Following Ecuador Earthquake",
											Link:  "http://www.wsj.com/articles/death-toll-in-ecuador-earthquake-climbs-as-correa-tours-ravaged-areas-1460993084?mod=fox_australian",
											Desc:  "The death toll in the magnitude-7.8 earthquake that struck this small country’s coast rose to 413, officials said.",
											Date:  "Tue, 19 Apr 2016 20:25:18 +0000",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 4, two dates
			Channels{
				&Channel{
					Owner: "cnet",
					Title: "CNET iPhone Update",
					Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
					Items: &Items{
						&Item{
							Title: "Apple iPhone SE owners bemoan audio bug - CNET",
							Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
							Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
							Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
						},
						&Item{
							Title: "9 settings every new iPhone owner should change - CNET",
							Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
							Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
							Date:  "Tue, 05 Apr 2016 20:25:18 +0000",
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-19",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "Apple iPhone SE owners bemoan audio bug - CNET",
											Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
											Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
											Date:  "Tue, 19 Apr 2016 17:25:18 +0000",
										},
									},
								},
							},
						},
					},
				},
				&Day{
					Date: "2016-04-05",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "9 settings every new iPhone owner should change - CNET",
											Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
											Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
											Date:  "Tue, 05 Apr 2016 20:25:18 +0000",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for idx, testCase := range testCases {
		marshaller := &Marshaller{
			Days: &Days{},
		}

		err := marshaller.ReArrange(testCase.channels)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(*marshaller.Days, testCase.days) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.days, *marshaller.Days)
		}
	}
}

func Test_load(t *testing.T) {
	testCases := []struct {
		date string // in
		load string
		day  Day
	}{
		{ // test case 0, date not exists
			"2016-04-25",
			"", // to say 'the file does not exist'
			Day{
				Date:   "2016-04-25",
				Owners: &Owners{},
			},
		},
		{ // test case 1, date exists
			"2016-04-25",
			`
		    {
		        "date": "2016-04-25",
		        "owners": [
		            {
		                "id": "cnet",
		                "channels": [
		                    {
		                        "title": "CNET iPhone Update",
		                        "desc": "Tips, news, how tos, and troubleshooting help for the iPhone.",
		                        "items": [
		                            {
		                                "title": "Apple iPhone SE owners bemoan audio bug - CNET",
		                                "link": "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
		                                "desc": "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners."
		                            }
		                        ]
		                    }
		                ]
		            }
		        ]
		    }
			`,
			Day{
				Date: "2016-04-25",
				Owners: &Owners{
					&Owner{
						Id: "cnet",
						Channels: &Channels{
							&Channel{
								Title: "CNET iPhone Update",
								Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
								Items: &Items{
									&Item{
										Title: "Apple iPhone SE owners bemoan audio bug - CNET",
										Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
										Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
									},
								},
							},
						},
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

		if strings.Compare(testCase.load, "") != 0 {
			// write json load to temp file
			file := filepath.Join(dir, testCase.date)
			err = ioutil.WriteFile(file, []byte(testCase.load), 0666)
			if err != nil {
				t.Error(err)
			}
		}

		marshaller, _ := NewMarshaller(dir)

		// under test
		day, err := marshaller.load(testCase.date)
		if err != nil {
			t.Error(err)
		}
		// assert
		if !reflect.DeepEqual(*day, testCase.day) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.day, *day)
		}

		// clean
		os.RemoveAll(dir)
	}
}

func Test_mergeItems(t *testing.T) {
	testCases := []struct {
		src  *Items // in
		dest *Items
		out  *Items //out
	}{
		{
			&Items{
				&Item{
					Title: "9 settings every new iPhone owner should change - CNET",
					Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
					Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
				},
			},
			&Items{
				&Item{
					Title: "Apple iPhone SE owners bemoan audio bug - CNET",
					Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
					Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
				},
			},
			&Items{
				&Item{
					Title: "Apple iPhone SE owners bemoan audio bug - CNET",
					Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
					Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
				},
				&Item{
					Title: "9 settings every new iPhone owner should change - CNET",
					Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
					Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
				},
			},
		},
	}

	for idx, testCase := range testCases {
		mergeItems(testCase.src, testCase.dest)

		// assert
		if !reflect.DeepEqual(*testCase.dest, *testCase.out) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.out, testCase.dest)
		}
	}
}

func Test_mergeChannels(t *testing.T) {
	testCases := []struct {
		src  *Channels // in
		dest *Channels
		out  *Channels //out
	}{
		{ // test case 0, different channels
			&Channels{
				&Channel{
					Title: "CNET iPhone Update",
					Items: &Items{
						&Item{
							Title: "9 settings every new iPhone owner should change - CNET",
							Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
							Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
						},
					},
				},
			},
			&Channels{
				&Channel{
					Title: "CNET Gaming",
					Items: &Items{
						&Item{
							Title: "Apple iPhone SE owners bemoan audio bug - CNET",
							Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
							Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
						},
					},
				},
			},
			&Channels{
				&Channel{
					Title: "CNET Gaming",
					Items: &Items{
						&Item{
							Title: "Apple iPhone SE owners bemoan audio bug - CNET",
							Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
							Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
						},
					},
				},
				&Channel{
					Title: "CNET iPhone Update",
					Items: &Items{
						&Item{
							Title: "9 settings every new iPhone owner should change - CNET",
							Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
							Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
						},
					},
				},
			},
		},
		{ // test case 1, same channels
			&Channels{
				&Channel{
					Title: "CNET iPhone Update",
					Items: &Items{
						&Item{
							Title: "9 settings every new iPhone owner should change - CNET",
							Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
							Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
						},
					},
				},
			},
			&Channels{
				&Channel{
					Title: "CNET iPhone Update",
					Items: &Items{
						&Item{
							Title: "Apple iPhone SE owners bemoan audio bug - CNET",
							Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
							Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
						},
					},
				},
			},
			&Channels{
				&Channel{
					Title: "CNET iPhone Update",
					Items: &Items{
						&Item{
							Title: "Apple iPhone SE owners bemoan audio bug - CNET",
							Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
							Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
						},
						&Item{
							Title: "9 settings every new iPhone owner should change - CNET",
							Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
							Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
						},
					},
				},
			},
		},
	}

	for idx, testCase := range testCases {
		mergeChannels(testCase.src, testCase.dest)

		// assert
		if !reflect.DeepEqual(*testCase.dest, *testCase.out) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.out, testCase.dest)
		}
	}
}

func Test_mergeOwners(t *testing.T) {
	testCases := []struct {
		src  *Owners // in
		dest *Owners
		out  *Owners //out
	}{
		{ // test case 0, different owners
			&Owners{
				&Owner{
					Id: "cnet01",
					Channels: &Channels{
						&Channel{
							Title: "CNET iPhone Update",
							Items: &Items{
								&Item{
									Title: "9 settings every new iPhone owner should change - CNET",
									Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
									Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
								},
							},
						},
					},
				},
			},
			&Owners{
				&Owner{
					Id: "cnet02",
					Channels: &Channels{
						&Channel{
							Title: "CNET Gaming",
							Items: &Items{
								&Item{
									Title: "Apple iPhone SE owners bemoan audio bug - CNET",
									Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
									Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
								},
							},
						},
					},
				},
			},
			&Owners{
				&Owner{
					Id: "cnet02",
					Channels: &Channels{
						&Channel{
							Title: "CNET Gaming",
							Items: &Items{
								&Item{
									Title: "Apple iPhone SE owners bemoan audio bug - CNET",
									Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
									Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
								},
							},
						},
					},
				},
				&Owner{
					Id: "cnet01",
					Channels: &Channels{
						&Channel{
							Title: "CNET iPhone Update",
							Items: &Items{
								&Item{
									Title: "9 settings every new iPhone owner should change - CNET",
									Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
									Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
								},
							},
						},
					},
				},
			},
		},
		{ // test case 1, same owner, different channels
			&Owners{
				&Owner{
					Id: "cnet01",
					Channels: &Channels{
						&Channel{
							Title: "CNET iPhone Update",
							Items: &Items{
								&Item{
									Title: "9 settings every new iPhone owner should change - CNET",
									Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
									Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
								},
							},
						},
					},
				},
			},
			&Owners{
				&Owner{
					Id: "cnet01",
					Channels: &Channels{
						&Channel{
							Title: "CNET Gaming",
							Items: &Items{
								&Item{
									Title: "Apple iPhone SE owners bemoan audio bug - CNET",
									Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
									Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
								},
							},
						},
					},
				},
			},
			&Owners{
				&Owner{
					Id: "cnet01",
					Channels: &Channels{
						&Channel{
							Title: "CNET Gaming",
							Items: &Items{
								&Item{
									Title: "Apple iPhone SE owners bemoan audio bug - CNET",
									Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
									Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
								},
							},
						},
						&Channel{
							Title: "CNET iPhone Update",
							Items: &Items{
								&Item{
									Title: "9 settings every new iPhone owner should change - CNET",
									Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
									Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
								},
							},
						},
					},
				},
			},
		},
		{ // test case 2, same owner, same channel
			&Owners{
				&Owner{
					Id: "cnet01",
					Channels: &Channels{
						&Channel{
							Title: "CNET iPhone Update",
							Items: &Items{
								&Item{
									Title: "9 settings every new iPhone owner should change - CNET",
									Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
									Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
								},
							},
						},
					},
				},
			},
			&Owners{
				&Owner{
					Id: "cnet01",
					Channels: &Channels{
						&Channel{
							Title: "CNET iPhone Update",
							Items: &Items{
								&Item{
									Title: "Apple iPhone SE owners bemoan audio bug - CNET",
									Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
									Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
								},
							},
						},
					},
				},
			},
			&Owners{
				&Owner{
					Id: "cnet01",
					Channels: &Channels{
						&Channel{
							Title: "CNET iPhone Update",
							Items: &Items{
								&Item{
									Title: "Apple iPhone SE owners bemoan audio bug - CNET",
									Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
									Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
								},
								&Item{
									Title: "9 settings every new iPhone owner should change - CNET",
									Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
									Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
								},
							},
						},
					},
				},
			},
		},
	}

	for idx, testCase := range testCases {
		mergeOwners(testCase.src, testCase.dest)

		// assert
		if !reflect.DeepEqual(*testCase.dest, *testCase.out) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.out, testCase.dest)
		}
	}
}

func Test_clean(t *testing.T) {
	testCases := []struct {
		in  Day // in
		out Day //out
	}{
		{
			Day{
				Owners: &Owners{
					&Owner{
						Id:       "cnet02",
						Channels: &Channels{},
					},
					&Owner{
						Id: "cnet01",
						Channels: &Channels{
							&Channel{
								Title: "CNET iPhone Update",
								Items: &Items{
									&Item{
										Title: "Apple iPhone SE owners bemoan audio bug - CNET",
										Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
										Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
									},
									&Item{
										Title: "9 settings every new iPhone owner should change - CNET",
										Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
										Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
									},
								},
							},
							&Channel{
								Title: "CNET Gaming",
								Items: &Items{},
							},
						},
					},
				},
			},
			Day{
				Owners: &Owners{
					&Owner{
						Id: "cnet01",
						Channels: &Channels{
							&Channel{
								Title: "CNET Gaming",
								Items: &Items{},
							},
							&Channel{
								Title: "CNET iPhone Update",
								Items: &Items{
									&Item{
										Title: "9 settings every new iPhone owner should change - CNET",
										Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
										Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
									},
									&Item{
										Title: "Apple iPhone SE owners bemoan audio bug - CNET",
										Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
										Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
									},
								},
							},
						},
					},
					&Owner{
						Id:       "cnet02",
						Channels: &Channels{},
					},
				},
			},
		},
	}

	for idx, testCase := range testCases {
		clean(testCase.in)

		// assert
		if !reflect.DeepEqual(testCase.in, testCase.out) {
			t.Errorf("[Test case %d], expected %+v, got %+v", idx, testCase.out, testCase.in)
		}
	}
}

func Test_Save(t *testing.T) {
	testCases := []struct {
		titles []string // in
		bodies []string
		src    Days
		merged Days // out
	}{
		{ // test case 0, same date, same owner, same channel, same item
			[]string{"2016-04-25"},
			[]string{
				`
				{
				    "date": "2016-04-25",
				    "owners": [
				        {
				            "id": "cnet",
				            "channels": [
				                {
				                    "title": "CNET iPhone Update",
				                    "desc": "Tips, news, how tos, and troubleshooting help for the iPhone.",
				                    "items": [
				                        {
				                            "title": "Apple iPhone SE owners bemoan audio bug - CNET",
				                            "link": "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
				                            "desc": "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners."
				                        }
				                    ]
				                }
				            ]
				        }
				    ]
				}
				`,
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "Apple iPhone SE owners bemoan audio bug - CNET",
											Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
											Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
										},
									},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "Apple iPhone SE owners bemoan audio bug - CNET",
											Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
											Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 1, same date, same owner, same channel, different items
			[]string{"2016-04-25"},
			[]string{
				`
				{
				    "date": "2016-04-25",
				    "owners": [
				        {
				            "id": "cnet",
				            "channels": [
				                {
				                    "title": "CNET iPhone Update",
				                    "desc": "Tips, news, how tos, and troubleshooting help for the iPhone.",
				                    "items": [
				                        {
				                            "title": "Apple iPhone SE owners bemoan audio bug - CNET",
				                            "link": "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
				                            "desc": "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners."
				                        }
				                    ]
				                }
				            ]
				        }
				    ]
				}
				`,
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "iPhone Upgrade Program launches at online Apple Store - CNET",
											Link:  "http://www.cnet.com/news/iphone-upgrade-program-launches-at-online-apple-store/#ftag=CAD4aa2096",
											Desc:  "The program that lets you upgrade your iPhone every year had been available only through Apple's retail outlets.",
										},
									},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "Apple iPhone SE owners bemoan audio bug - CNET",
											Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
											Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
										},
										&Item{
											Title: "iPhone Upgrade Program launches at online Apple Store - CNET",
											Link:  "http://www.cnet.com/news/iphone-upgrade-program-launches-at-online-apple-store/#ftag=CAD4aa2096",
											Desc:  "The program that lets you upgrade your iPhone every year had been available only through Apple's retail outlets.",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 2, same date, same owner, different channels
			[]string{"2016-04-25"},
			[]string{
				`
				{
				    "date": "2016-04-25",
				    "owners": [
				        {
				            "id": "cnet",
				            "channels": [
				                {
				                    "title": "CNET iPhone Update",
				                    "desc": "Tips, news, how tos, and troubleshooting help for the iPhone.",
				                    "items": [
				                        {
				                            "title": "Apple iPhone SE owners bemoan audio bug - CNET",
				                            "link": "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
				                            "desc": "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners."
				                        }
				                    ]
				                }
				            ]
				        }
				    ]
				}
				`,
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET Gaming",
									Desc:  "Game on! Get the latest in gaming news, video game reviews, computer games & video game consoles.",
									Items: &Items{
										&Item{
											Title: "The Ikea VR experience, now with the meatball update you've been waiting for - CNET",
											Link:  "http://www.cnet.com/news/the-ikea-vr-game-now-with-the-meatball-update-youve-been-waiting-for/#ftag=CADa872701",
											Desc:  "The fans have spoken. You can now interact with virtual meatballs in the Ikea VR experience.",
										},
									},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET Gaming",
									Desc:  "Game on! Get the latest in gaming news, video game reviews, computer games & video game consoles.",
									Items: &Items{
										&Item{
											Title: "The Ikea VR experience, now with the meatball update you've been waiting for - CNET",
											Link:  "http://www.cnet.com/news/the-ikea-vr-game-now-with-the-meatball-update-youve-been-waiting-for/#ftag=CADa872701",
											Desc:  "The fans have spoken. You can now interact with virtual meatballs in the Ikea VR experience.",
										},
									},
								},
								&Channel{
									Title: "CNET iPhone Update",
									Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
									Items: &Items{
										&Item{
											Title: "Apple iPhone SE owners bemoan audio bug - CNET",
											Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
											Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 3, same date, different owners
			[]string{"2016-04-25"},
			[]string{
				`
				{
				    "date": "2016-04-25",
				    "owners": [
				        {
				            "id": "wsj",
				            "channels": [
				                {
				                    "title": "WSJ.com: World News",
				                    "desc": "World News",
				                    "items": [
				                        {
				                            "title": "U.S., Turkey Step Up Border Campaign Against Islamic State",
				                            "link": "http://www.wsj.com/articles/u-s-turkey-step-up-border-campaign-against-islamic-state-1461684454?mod=fox_australian",
				                            "desc": "Ankara and Washington plan to deploy advanced rocket launchers and more Turkish forces to the Turkish-Syrian border in an effort to choke off a crucial Islamic State supply route."
				                        }
				                    ]
				                }
				            ]
				        }
				    ]
				}
				`,
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET Gaming",
									Desc:  "Game on! Get the latest in gaming news, video game reviews, computer games & video game consoles.",
									Items: &Items{
										&Item{
											Title: "The Ikea VR experience, now with the meatball update you've been waiting for - CNET",
											Link:  "http://www.cnet.com/news/the-ikea-vr-game-now-with-the-meatball-update-youve-been-waiting-for/#ftag=CADa872701",
											Desc:  "The fans have spoken. You can now interact with virtual meatballs in the Ikea VR experience.",
										},
									},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET Gaming",
									Desc:  "Game on! Get the latest in gaming news, video game reviews, computer games & video game consoles.",
									Items: &Items{
										&Item{
											Title: "The Ikea VR experience, now with the meatball update you've been waiting for - CNET",
											Link:  "http://www.cnet.com/news/the-ikea-vr-game-now-with-the-meatball-update-youve-been-waiting-for/#ftag=CADa872701",
											Desc:  "The fans have spoken. You can now interact with virtual meatballs in the Ikea VR experience.",
										},
									},
								},
							},
						},
						&Owner{
							Id: "wsj",
							Channels: &Channels{
								&Channel{
									Title: "WSJ.com: World News",
									Desc:  "World News",
									Items: &Items{
										&Item{
											Title: "U.S., Turkey Step Up Border Campaign Against Islamic State",
											Link:  "http://www.wsj.com/articles/u-s-turkey-step-up-border-campaign-against-islamic-state-1461684454?mod=fox_australian",
											Desc:  "Ankara and Washington plan to deploy advanced rocket launchers and more Turkish forces to the Turkish-Syrian border in an effort to choke off a crucial Islamic State supply route.",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // test case 4, different dates
			[]string{"2016-04-26"},
			[]string{
				`
				{
				    "date": "2016-04-26",
				    "owners": [
				        {
				            "id": "wsj",
				            "channels": [
				                {
				                    "title": "WSJ.com: World News",
				                    "desc": "World News",
				                    "items": [
				                        {
				                            "title": "U.S., Turkey Step Up Border Campaign Against Islamic State",
				                            "link": "http://www.wsj.com/articles/u-s-turkey-step-up-border-campaign-against-islamic-state-1461684454?mod=fox_australian",
				                            "desc": "Ankara and Washington plan to deploy advanced rocket launchers and more Turkish forces to the Turkish-Syrian border in an effort to choke off a crucial Islamic State supply route."
				                        }
				                    ]
				                }
				            ]
				        }
				    ]
				}
				`,
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET Gaming",
									Desc:  "Game on! Get the latest in gaming news, video game reviews, computer games & video game consoles.",
									Items: &Items{
										&Item{
											Title: "The Ikea VR experience, now with the meatball update you've been waiting for - CNET",
											Link:  "http://www.cnet.com/news/the-ikea-vr-game-now-with-the-meatball-update-youve-been-waiting-for/#ftag=CADa872701",
											Desc:  "The fans have spoken. You can now interact with virtual meatballs in the Ikea VR experience.",
										},
									},
								},
							},
						},
					},
				},
			},
			Days{
				&Day{
					Date: "2016-04-25",
					Owners: &Owners{
						&Owner{
							Id: "cnet",
							Channels: &Channels{
								&Channel{
									Title: "CNET Gaming",
									Desc:  "Game on! Get the latest in gaming news, video game reviews, computer games & video game consoles.",
									Items: &Items{
										&Item{
											Title: "The Ikea VR experience, now with the meatball update you've been waiting for - CNET",
											Link:  "http://www.cnet.com/news/the-ikea-vr-game-now-with-the-meatball-update-youve-been-waiting-for/#ftag=CADa872701",
											Desc:  "The fans have spoken. You can now interact with virtual meatballs in the Ikea VR experience.",
										},
									},
								},
							},
						},
					},
				},
				&Day{
					Date: "2016-04-26",
					Owners: &Owners{
						&Owner{
							Id: "wsj",
							Channels: &Channels{
								&Channel{
									Title: "WSJ.com: World News",
									Desc:  "World News",
									Items: &Items{
										&Item{
											Title: "U.S., Turkey Step Up Border Campaign Against Islamic State",
											Link:  "http://www.wsj.com/articles/u-s-turkey-step-up-border-campaign-against-islamic-state-1461684454?mod=fox_australian",
											Desc:  "Ankara and Washington plan to deploy advanced rocket launchers and more Turkish forces to the Turkish-Syrian border in an effort to choke off a crucial Islamic State supply route.",
										},
									},
								},
							},
						},
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

		for idx, title := range testCase.titles {
			// write json load to temp file
			file := filepath.Join(dir, title)
			err = ioutil.WriteFile(file, []byte(testCase.bodies[idx]), 0666)
			if err != nil {
				t.Error(err)
			}
		}

		marshaller, _ := NewMarshaller(dir)
		marshaller.Days = &testCase.src

		// under test
		err = marshaller.Save()
		if err != nil {
			t.Error(err)
		}

		for _, expected := range testCase.merged {
			actual, err := marshaller.load(expected.Date)
			if err != nil {
				t.Error(err)
			}

			// assert
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("[Test case %d], expected %+v, got %+v", idx, expected, *actual)
			}
		}

		// clean
		os.RemoveAll(dir)
	}
}
