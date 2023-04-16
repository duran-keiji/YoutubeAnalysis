// チャンネルIDと動画IDと視聴回数と登録者数取得（ver3）

package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const APIKey = "AIzaSyDAKFWvZl0y9DQI47F2tEaJPncMvnDPdkE"

  // [1] 構造体のフィールドとJSONのキーを紐付け
type Search struct {
	NextPageToken       string    `json:"nextPageToken"`
	Items []*SearchItem `json:"items"`
}

type SearchItem struct {
	ID SearchID `json:"id"`
	Snippet SearchSnippet `json:"snippet"`
	ChannelTitle string `json:"channelTitle"`
}

type SearchID struct {
	VideoId string `json:"videoId"`
}


type SearchSnippet struct {
	ChannelId string `json:"channelId"`
	Title string `json:"title"`
}

type Data struct {
	ChannelId string `json:"channelId"`
	VideoId  string    `json:"videoId"`
	ViewCount  string    `json:"viewCount"`
	SubscriberCount  string    `json:"subscriberCount"`
}
func get_youtube_data() string {
	url := "https://www.googleapis.com/youtube/v3/search"
 
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Fatal(err)
	}
	
	//クエリパラメータ
	params := request.URL.Query()
	params.Add("key", APIKey)
	params.Add("q", "洋楽")
	params.Add("part", "snippet, id")
	params.Add("maxResults", "2")

    request.URL.RawQuery = params.Encode()
 
	fmt.Println(request.URL.String()) //https://jsonplaceholder.typicode.com/todos?userId=1
	
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
        Timeout: timeout,
	}
 
	response, err := client.Do(request)
	if err != nil{
		log.Fatal(err)
	}
	
	defer response.Body.Close()
 
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("body=", string(body))

	// [3] 配列型のJSONデータを読み込む
	// FIXME: objが取得できていない
	// obj := make([]*Search, 2)
	var obj []*Search
	
	// obj := make([]*Search, 0)
	// obj := []*Search{}
	fmt.Printf("全 %d 件\n\n", obj)
	json.Unmarshal(body, &obj)
	fmt.Printf("全 %d 件\n\n", obj)

	// スライスの要素数を取得
	// fmt.Printf("全 %d 件\n\n", len(obj))

	var ChannelId string
	var VideoId string
	var ViewCount string
	var SubscriberCount string
	var data []*Data

	// スライスの要素を順に取り出す
	for _, search := range obj {
		if len(search.Items) != 0 {
			for i, items := range search.Items {
				ChannelId = items.Snippet.ChannelId
				VideoId = items.ID.VideoId
				ViewCount = get_youtube_viewCount(string(VideoId))
				SubscriberCount = get_youtube_subscriberCount(string(ChannelId))
				fmt.Printf(string(i),ChannelId,VideoId,ViewCount,SubscriberCount)
				// obj[i].ChannelId = ChannelId
				// obj[i]["VideoId"] = VideoId
				// obj[i]["ViewCount"] = ViewCount
				// obj[i]["SubscriberCount"] = SubscriberCount
				data = append(data, &Data{ChannelId: string(ChannelId), VideoId: string(VideoId), ViewCount: string(ViewCount), SubscriberCount: string(SubscriberCount)})
			}
		}
	} 
	fmt.Printf("respons=", data)

	res, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return string(res)


	// var search Search
	// json.Unmarshal(body, &search)
	// fmt.Println("ChannelId=", search.Items[0].Snippet.ChannelId)
	// fmt.Println("VideoId=", search.Items[0].ID.VideoId)

	

	// ChannelId, err := json.Marshal(search.Items[0].Snippet.ChannelId)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// VideoId, err := json.Marshal(search.Items[0].ID.VideoId)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ViewCount := get_youtube_viewCount(string(ChannelId), string(VideoId))
	// SubscriberCount := get_youtube_subscriberCount(string(ChannelId))

	// return string(ChannelId), string(VideoId)
}



// [1] 構造体のフィールドとJSONのキーを紐付け
type Video struct {
	Items []*VideoItem `json:"items"`
}

type VideoItem struct {
	ID string `json:"id"`
	Statistics VideoStatistic `json:"statistics"`
}

type VideoStatistic struct {
	ViewCount string `json:"viewCount"`
}

func get_youtube_viewCount(VideoId string) string {

	url := "https://www.googleapis.com/youtube/v3/videos"
 
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Fatal(err)
	}

	// なぜかVideoIdにダブルクォートが含まれるためダブルクォートをカットする。これをしないとparams.Add("id", VideoId)のvideoIdでダブルクォート込みで検索してしまい値が取れない。
	if strings.Contains(VideoId, "\"") {
		VideoId = strings.Replace(VideoId, "\"", "", -1)
	}
	
	//クエリパラメータ
	params := request.URL.Query()
	params.Add("key", APIKey)
	params.Add("id", VideoId)
	params.Add("part", "statistics") 
	params.Add("maxResults", "2")

    request.URL.RawQuery = params.Encode()
 
	fmt.Println(request.URL.String()) //https://jsonplaceholder.typicode.com/todos?userId=1

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
        Timeout: timeout,
	}
 
	response, err := client.Do(request)
	if err != nil{
		log.Fatal(err)
	}
	
	defer response.Body.Close()
 
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("body=", string(body))

	var video Video
	json.Unmarshal(body, &video)
	fmt.Println("ViewCount=", video.Items[0].Statistics.ViewCount)


	ViewCount, err := json.Marshal(video.Items[0].Statistics.ViewCount)
	if err != nil {
		log.Fatal(err)
	}

	return string(ViewCount)

	// return string(ViewCount)
}

// [1] 構造体のフィールドとJSONのキーを紐付け
type Channel struct {
	Items []*ChannelItem `json:"items"`
}

// [2] JSONと同様のネスト構造で構造体を定義
type ChannelItem struct {
	Statistics ChannelStatistic `json:"statistics"`
}

type ChannelStatistic struct {
	ViewCount string `json:"viewCount"`
	SubscriberCount string `json:"subscriberCount"`
}

func get_youtube_subscriberCount(ChannelId string) string {

	url := "https://www.googleapis.com/youtube/v3/channels"
 
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Fatal(err)
	}

	// なぜかchannelIdにダブルクォートが含まれるためダブルクォートをカットする。これをしないとparams.Add("id", channelId)のchannelIdでダブルクォート込みで検索してしまい値が取れない。
	if strings.Contains(ChannelId, "\"") {
		ChannelId = strings.Replace(ChannelId, "\"", "", -1)
	}
	
	//クエリパラメータ
	params := request.URL.Query()
	params.Add("key", APIKey)
	params.Add("id", ChannelId)
	params.Add("part", "statistics") 
	params.Add("maxResults", "2")

    request.URL.RawQuery = params.Encode()
 
	fmt.Println(request.URL.String()) //https://jsonplaceholder.typicode.com/todos?userId=1

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
        Timeout: timeout,
	}
 
	response, err := client.Do(request)
	if err != nil{
		log.Fatal(err)
	}
	
	defer response.Body.Close()
 
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("body=", string(body))

	var channel Channel
	json.Unmarshal(body, &channel)
	fmt.Println("SubscriberCount=", channel.Items[0].Statistics.SubscriberCount)


	SubscriberCount, err := json.Marshal(channel.Items[0].Statistics.SubscriberCount)
	if err != nil {
		log.Fatal(err)
	}

	return string(SubscriberCount)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	data := get_youtube_data()

	return events.APIGatewayProxyResponse{
 
		 Body:       fmt.Sprintf(data),
		 StatusCode: 200,
	 }, nil
 }

func main() {

    lambda.Start(handler)
}