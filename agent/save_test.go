package agent

import (
	"reflect"
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
