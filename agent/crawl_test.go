package agent

import (
	"container/list"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// check we're correcly parsing xml
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

func Test_Crawl(t *testing.T) {
	testCases := []struct {
		bodies []string
		rss    Rss
	}{
		{ // test case 0, one channel
			[]string{
				`
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
    xmlns:atom="http://www.w3.org/2005/Atom" >
    <channel>
        <title>WSJ.com: World News</title>
        <link>http://online.wsj.com/page/2_0006.html</link>
        <description>World News</description>
        <language>en-us</language>
        <copyright>copyright  &#169; 2016 Dow Jones &amp; Company, Inc.</copyright>
        <lastBuildDate>Tue, 19 Apr 2016 21:45:53 EDT</lastBuildDate>
        <image>
            <title>WSJ.com: World News</title>
            <url>http://online.wsj.com/img/wsj_sm_logo.gif</url>
            <link>http://online.wsj.com/page/2_0006.html</link>
        </image>
        <atom:link href="http://online.wsj.com/page/2_0006.html" rel="self" type="application/rss+xml" />
        <!--Item 1 of World_Section_Front_2 -->
        <item>
            <title>Obama's Mideast Mission: Get Saudis, Iran to Make Nice</title>
            <guid isPermaLink="false">SB10225542119583864159404582016363723781068</guid>
            <link>http://www.wsj.com/articles/obamas-mideast-mission-get-saudis-iran-to-make-nice-1461111595?mod=fox_australian</link>
            <description>President Obama, visiting Saudi Arabia, will encourage Mideast stability through better relations between Saudis and Iran, but America is seen as part of the problem.</description>
            <media:content
                xmlns:media="http://search.yahoo.com/mrss"
               url="http://s.wsj.net/public/resources/images/BN-NQ267_SAUDIR_G_20160419200245.jpg" 
               type="image/jpeg"
               medium="image"
               height="369"
               width="553">
                <media:description>image</media:description>
            </media:content>
            <category>PAID</category>
            <pubDate>Tue, 19 Apr 2016 20:20:01 EDT</pubDate>
        </item>
        <!--Item 2 of World_Section_Front_2 -->
        <item>
            <title>Taliban Coordinated Attack Kills at Least 28 in Kabul</title>
            <guid isPermaLink="false">SB10834865168797973818204582015193317713700</guid>
            <link>http://www.wsj.com/articles/kabul-rocked-by-suicide-attack-and-gunfire-afghan-official-says-1461044819?mod=fox_australian</link>
            <description>The deadliest attack in the Afghan capital since August was carried out on a compound housing the agency charged with protecting top officials and visiting dignitaries.</description>
            <media:content
                xmlns:media="http://search.yahoo.com/mrss"
               url="http://s.wsj.net/public/resources/images/P1-BX153_CATDOO_G_20160419213530.jpg" 
               type="image/jpeg"
               medium="image"
               height="369"
               width="553">
                <media:description>image</media:description>
            </media:content>
            <category>FREE</category>
            <pubDate>Tue, 19 Apr 2016 21:38:51 EDT</pubDate>
        </item>
    </channel>
</rss>
<!-- fastdynapage - sbkj2kwebappp10 - Tue 04/19/16 - 21:45:53 EDT -->
				`,
			},
			Rss{
				XMLName: xml.Name{
					Space: "",
					Local: "",
				},
				Channels: []*Channel{
					&Channel{
						Title: "WSJ.com: World News",
						Desc:  "World News",
						Items: []*Item{
							&Item{
								Title: "Obama's Mideast Mission: Get Saudis, Iran to Make Nice",
								Link:  "http://www.wsj.com/articles/obamas-mideast-mission-get-saudis-iran-to-make-nice-1461111595?mod=fox_australian",
								Desc:  "President Obama, visiting Saudi Arabia, will encourage Mideast stability through better relations between Saudis and Iran, but America is seen as part of the problem.",
							},
							&Item{
								Title: "Taliban Coordinated Attack Kills at Least 28 in Kabul",
								Link:  "http://www.wsj.com/articles/kabul-rocked-by-suicide-attack-and-gunfire-afghan-official-says-1461044819?mod=fox_australian",
								Desc:  "The deadliest attack in the Afghan capital since August was carried out on a compound housing the agency charged with protecting top officials and visiting dignitaries.",
							},
						},
					},
				},
			},
		},
		{ // test case 1, multiple channels
			[]string{
				`
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
    xmlns:atom="http://www.w3.org/2005/Atom" >
    <channel>
        <title>WSJ.com: World News</title>
        <link>http://online.wsj.com/page/2_0006.html</link>
        <description>World News</description>
        <language>en-us</language>
        <copyright>copyright  &#169; 2016 Dow Jones &amp; Company, Inc.</copyright>
        <lastBuildDate>Tue, 19 Apr 2016 21:45:53 EDT</lastBuildDate>
        <image>
            <title>WSJ.com: World News</title>
            <url>http://online.wsj.com/img/wsj_sm_logo.gif</url>
            <link>http://online.wsj.com/page/2_0006.html</link>
        </image>
        <atom:link href="http://online.wsj.com/page/2_0006.html" rel="self" type="application/rss+xml" />
        <!--Item 1 of World_Section_Front_2 -->
        <item>
            <title>Obama's Mideast Mission: Get Saudis, Iran to Make Nice</title>
            <guid isPermaLink="false">SB10225542119583864159404582016363723781068</guid>
            <link>http://www.wsj.com/articles/obamas-mideast-mission-get-saudis-iran-to-make-nice-1461111595?mod=fox_australian</link>
            <description>President Obama, visiting Saudi Arabia, will encourage Mideast stability through better relations between Saudis and Iran, but America is seen as part of the problem.</description>
            <media:content
                xmlns:media="http://search.yahoo.com/mrss"
               url="http://s.wsj.net/public/resources/images/BN-NQ267_SAUDIR_G_20160419200245.jpg" 
               type="image/jpeg"
               medium="image"
               height="369"
               width="553">
                <media:description>image</media:description>
            </media:content>
            <category>PAID</category>
            <pubDate>Tue, 19 Apr 2016 20:20:01 EDT</pubDate>
        </item>
        <!--Item 2 of World_Section_Front_2 -->
        <item>
            <title>Taliban Coordinated Attack Kills at Least 28 in Kabul</title>
            <guid isPermaLink="false">SB10834865168797973818204582015193317713700</guid>
            <link>http://www.wsj.com/articles/kabul-rocked-by-suicide-attack-and-gunfire-afghan-official-says-1461044819?mod=fox_australian</link>
            <description>The deadliest attack in the Afghan capital since August was carried out on a compound housing the agency charged with protecting top officials and visiting dignitaries.</description>
            <media:content
                xmlns:media="http://search.yahoo.com/mrss"
               url="http://s.wsj.net/public/resources/images/P1-BX153_CATDOO_G_20160419213530.jpg" 
               type="image/jpeg"
               medium="image"
               height="369"
               width="553">
                <media:description>image</media:description>
            </media:content>
            <category>FREE</category>
            <pubDate>Tue, 19 Apr 2016 21:38:51 EDT</pubDate>
        </item>
    </channel>
</rss>
<!-- fastdynapage - sbkj2kwebappp10 - Tue 04/19/16 - 21:45:53 EDT -->
				`,
				`
<rss
    xmlns:media="http://search.yahoo.com/mrss/"
    xmlns:dc="http://purl.org/dc/elements/1.1/"
    xmlns:atom="http://www.w3.org/2005/Atom"
    xmlns:content="http://purl.org/rss/1.0/modules/content/" version="2.0">
    <channel>
        <atom:link href="http://feed.cnet.com/feed/collection" rel="self" type="application/rss+xml" />
        <title>CNET iPhone Update</title>
        <description>Tips, news, how tos, and troubleshooting help for the iPhone.</description>
        <image>
            <url>http://i.i.cbsi.com/cnwk.1d/i/ne/gr/prtnr/CNET_Logo_150.gif</url>
            <title>CNET iPhone Update</title>
            <link>http://www.cnet.com/#ftag=CAD4aa2096</link>
        </image>
        <link>http://www.cnet.com/#ftag=CAD4aa2096</link>
        <item>
            <title>
                                        9 settings every new iPhone owner should change - CNET                                    </title>
            <link>
                                        http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096                                    </link>
            <guid isPermaLink="false">
                                        51ce1a3e-cff1-4703-ac79-88be56124acc                                    </guid>
            <pubDate>
                                        Tue, 19 Apr 2016 17:25:18 +0000                                    </pubDate>
            <description>
                                        Whether you&#039;re a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.                                    </description>
            <media:thumbnail url="https://cnet4.cbsistatic.com/hub/i/r/2016/04/18/fa2b5a53-f02e-4a48-86a7-0f2709b17d14/thumbnail/300x230/158c0d9a8bfc13ac81bd2bf514d9971f/ios-brightness-slider.jpg"/>
            <dc:creator>
                                        Rick
                                                                            Broida                                    </dc:creator>
        </item>
        <item>
            <title>
                                        2017 iPhone will replace aluminum body with glass, says analyst - CNET                                    </title>
            <link>
                                        http://www.cnet.com/news/2017-iphone-glass-body-aluminum-ming-chi-kuo-samsung/#ftag=CAD4aa2096                                    </link>
            <guid isPermaLink="false">
                                        926169d8-dcac-45d3-a00e-d846a994e7b0                                    </guid>
            <pubDate>
                                        Mon, 18 Apr 2016 16:06:29 +0000                                    </pubDate>
            <description>
                                        If true, the goal would be to make the iPhone more distinctive from smartphones that use metal bodies.                                    </description>
            <media:thumbnail url="https://cnet2.cbsistatic.com/hub/i/r/2014/09/09/c311a2da-e4c3-4d27-9956-7931b9bc2235/thumbnail/300x230/f757c1a5b610753462620e73ba7bfeb1/medium-carousel-apple-iphone-6-plus-7-new-isight-camera-770.jpg"/>
            <dc:creator>
                                        Lance
                                                                            Whitney                                    </dc:creator>
        </item>
    </channel>
</rss>
				`,
			},
			Rss{
				XMLName: xml.Name{
					Space: "",
					Local: "",
				},
				Channels: []*Channel{
					&Channel{
						Title: "WSJ.com: World News",
						Desc:  "World News",
						Items: []*Item{
							&Item{
								Title: "Obama's Mideast Mission: Get Saudis, Iran to Make Nice",
								Link:  "http://www.wsj.com/articles/obamas-mideast-mission-get-saudis-iran-to-make-nice-1461111595?mod=fox_australian",
								Desc:  "President Obama, visiting Saudi Arabia, will encourage Mideast stability through better relations between Saudis and Iran, but America is seen as part of the problem.",
							},
							&Item{
								Title: "Taliban Coordinated Attack Kills at Least 28 in Kabul",
								Link:  "http://www.wsj.com/articles/kabul-rocked-by-suicide-attack-and-gunfire-afghan-official-says-1461044819?mod=fox_australian",
								Desc:  "The deadliest attack in the Afghan capital since August was carried out on a compound housing the agency charged with protecting top officials and visiting dignitaries.",
							},
						},
					},
					&Channel{
						Title: "CNET iPhone Update",
						Desc:  "Tips, news, how tos, and troubleshooting help for the iPhone.",
						Items: []*Item{
							&Item{
								Title: "9 settings every new iPhone owner should change - CNET",
								Link:  "http://www.cnet.com/how-to/9-settings-you-should-change-on-your-new-iphone/#ftag=CAD4aa2096",
								Desc:  "Whether you're a newcomer to iOS or just upgrading to a newer model, consider tweaking these settings to improve performance and battery life.",
							},
							&Item{
								Title: "2017 iPhone will replace aluminum body with glass, says analyst - CNET",
								Link:  "http://www.cnet.com/news/2017-iphone-glass-body-aluminum-ming-chi-kuo-samsung/#ftag=CAD4aa2096",
								Desc:  "If true, the goal would be to make the iPhone more distinctive from smartphones that use metal bodies.",
							},
						},
					},
				},
			},
		},
	}

	for idx, testCase := range testCases {
		// create a list of resp to pass to the http server
		resp := list.New()
		for _, body := range testCase.bodies {
			resp.PushBack(body)
		}

		// bootstrap a test http server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/xml")
			// feed the response from the list
			fmt.Fprintln(w, resp.Front().Value.(string))
			// remove to recently used element to simulate the behaviour of a stack
			resp.Remove(resp.Front())
		}))
		defer ts.Close()

		loader, err := NewLoader()
		if err != nil {
			t.Error(err)
		}
		loader.Urls = make([]string, len(testCase.bodies))
		for idx, _ := range loader.Urls {
			loader.Urls[idx] = ts.URL
		}

		crawler, err := NewCrawler()
		if err != nil {
			t.Error(err)
		}

		crawler.Crawl(loader)

		if !reflect.DeepEqual(crawler.Rss, testCase.rss) {
			t.Errorf("[Test case %d] expecting %v, got %v", idx, testCase.rss, crawler.Rss)
		}
	}
}

// Smoke test
func Test_Crawl_(t *testing.T) {
	t.Skip()
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
