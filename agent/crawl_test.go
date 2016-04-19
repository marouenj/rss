package agent

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func Test_RssUnmarshal(t *testing.T) {
	testCases := []struct {
		in  string
		out Rss
	}{
		{
			`
            <rss>
                <channel>
                    <title>CNET iPhone Update</title>
                    <description>Tips, news, how tos, and troubleshooting help for the iPhone.</description>
                    <item>
                        <title>Apple iPhone SE owners bemoan audio bug - CNET</title>
                        <link>http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096</link>
                        <description>Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.</description>
                    </item>
                </channel>
            </rss>`,
			Rss{
				XMLName: xml.Name{
					Space: "",
					Local: "rss",
				},
				Channels: []*Channel{
					&Channel{
						Title: "CNET iPhone Update",
						Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
						Items: []*Item{
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
	}

	for _, testCase := range testCases {
		var root Rss
		err := xml.Unmarshal([]byte(testCase.in), &root)
		if err != nil {
			t.Error(err)
		}

		//check
		if len(root.Channels) != len(testCase.out.Channels) {
			t.Errorf("[Input: %v], expected %v, got %v", testCase.in, len(testCase.out.Channels), len(root.Channels))
		}

		if len(root.Channels[0].Items) != len(testCase.out.Channels[0].Items) {
			t.Errorf("[Input: %v], expected %v, got %v", testCase.in, len(testCase.out.Channels), len(root.Channels))
		}

		if !reflect.DeepEqual(root, testCase.out) {
			t.Errorf("[Input: %v], expected %v, got %v", testCase.in, testCase.out, root)
		}
	}
}

func Test_merge(t *testing.T) {
	testCases := []struct {
		src    Rss
		dest   Rss
		merged Rss
	}{
		{ // test case 0
			Rss{
				XMLName: xml.Name{
					Space: "",
					Local: "rss",
				},
				Channels: []*Channel{
					&Channel{
						Title: "CNET iPhone Update",
						Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
						Items: []*Item{
							&Item{
								Title: "Apple iPhone SE owners bemoan audio bug - CNET",
								Link:  "http://www.cnet.com/news/apple-iphone-se-owners-complain-of-phone-call-audio-bug/#ftag=CAD4aa2096",
								Desc:  "Introduced with the latest update to iOS, the glitch distorts the quality of phone calls made via Bluetooth, according to some owners.",
							},
						},
					},
				},
			},
			Rss{
				XMLName: xml.Name{
					Space: "",
					Local: "rss",
				},
				Channels: []*Channel{
					&Channel{
						Title: "WSJ.com: World News",
						Desc:  "World News",
						Items: []*Item{
							&Item{
								Title: "Death Toll Rises Following Ecuador Earthquake",
								Link:  "http://www.wsj.com/articles/death-toll-in-ecuador-earthquake-climbs-as-correa-tours-ravaged-areas-1460993084?mod=fox_australian",
								Desc:  "The death toll in the magnitude-7.8 earthquake that struck this small country’s coast rose to 413, officials said.",
							},
						},
					},
				},
			},
			Rss{
				XMLName: xml.Name{
					Space: "",
					Local: "rss",
				},
				Channels: []*Channel{
					&Channel{
						Title: "WSJ.com: World News",
						Desc:  "World News",
						Items: []*Item{
							&Item{
								Title: "Death Toll Rises Following Ecuador Earthquake",
								Link:  "http://www.wsj.com/articles/death-toll-in-ecuador-earthquake-climbs-as-correa-tours-ravaged-areas-1460993084?mod=fox_australian",
								Desc:  "The death toll in the magnitude-7.8 earthquake that struck this small country’s coast rose to 413, officials said.",
							},
						},
					},
					&Channel{
						Title: "CNET iPhone Update",
						Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
						Items: []*Item{
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
		{ // test case 1, the src is empty
			Rss{
				XMLName: xml.Name{
					Space: "",
					Local: "rss",
				},
				Channels: []*Channel{ // the src is empty
				},
			},
			Rss{
				XMLName: xml.Name{
					Space: "",
					Local: "rss",
				},
				Channels: []*Channel{
					&Channel{
						Title: "WSJ.com: World News",
						Desc:  "World News",
						Items: []*Item{
							&Item{
								Title: "Death Toll Rises Following Ecuador Earthquake",
								Link:  "http://www.wsj.com/articles/death-toll-in-ecuador-earthquake-climbs-as-correa-tours-ravaged-areas-1460993084?mod=fox_australian",
								Desc:  "The death toll in the magnitude-7.8 earthquake that struck this small country’s coast rose to 413, officials said.",
							},
						},
					},
				},
			},
			Rss{
				XMLName: xml.Name{
					Space: "",
					Local: "rss",
				},
				Channels: []*Channel{
					&Channel{
						Title: "WSJ.com: World News",
						Desc:  "World News",
						Items: []*Item{
							&Item{
								Title: "Death Toll Rises Following Ecuador Earthquake",
								Link:  "http://www.wsj.com/articles/death-toll-in-ecuador-earthquake-climbs-as-correa-tours-ravaged-areas-1460993084?mod=fox_australian",
								Desc:  "The death toll in the magnitude-7.8 earthquake that struck this small country’s coast rose to 413, officials said.",
							},
						},
					},
				},
			},
		},
	}
	for idx, testCase := range testCases {
		crawler, err := NewCrawler()
		if err != nil {
			t.Error(err)
		}

		crawler.Rss = testCase.dest
		crawler.merge(testCase.src.Channels)

		if !reflect.DeepEqual(crawler.Rss, testCase.merged) {
			t.Errorf("[Test case %d] expecting %v, got %v", idx, testCase.merged, crawler.Rss)
		}
	}
}

// Smoke test
func Test_Crawl(t *testing.T) {
	testCases := []struct {
		urls []string
	}{
		{ //test case 0
			[]string{"http://www.wsj.com/xml/rss/3_7085.xml", "http://www.cnet.com/rss/iphone-update/"},
		},
	}

	for _, testCase := range testCases {
		loader, err := NewLoader()
		if err != nil {
			t.Error(err)
		}
		loader.Urls = testCase.urls

		crawler, err := NewCrawler()
		if err != nil {
			t.Error(err)
		}

		crawler.Crawl(loader)
	}
}
