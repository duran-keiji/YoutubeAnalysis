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

func get_youtube_id(APIKey string) (string, string) {
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

	var search Search
	json.Unmarshal(body, &search)
	fmt.Println("ChannelId=", search.Items[0].Snippet.ChannelId)
	fmt.Println("VideoId=", search.Items[0].ID.VideoId)

	ChannelId, err := json.Marshal(search.Items[0].Snippet.ChannelId)
	if err != nil {
		log.Fatal(err)
	}

	VideoId, err := json.Marshal(search.Items[0].ID.VideoId)
	if err != nil {
		log.Fatal(err)
	}

	return string(ChannelId), string(VideoId)
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

func get_youtube_viewCount(APIKey string, videoId string) string {

	url := "https://www.googleapis.com/youtube/v3/videos"
 
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Fatal(err)
	}

	// なぜかvideoIdにダブルクォートが含まれるためダブルクォートをカットする。これをしないとparams.Add("id", videoId)のvideoIdでダブルクォート込みで検索してしまい値が取れない。
	if strings.Contains(videoId, "\"") {
		videoId = strings.Replace(videoId, "\"", "", -1)
	}
	
	//クエリパラメータ
	params := request.URL.Query()
	params.Add("key", APIKey)
	params.Add("id", videoId)
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
	fmt.Println("body=", string(body))

	var video Video
	json.Unmarshal(body, &video)
	fmt.Println("ViewCount=", video.Items[0].Statistics.ViewCount)


	ViewCount, err := json.Marshal(video.Items[0].Statistics.ViewCount)
	if err != nil {
		log.Fatal(err)
	}

	return string(ViewCount)
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

func get_youtube_subscriberCount(APIKey string, channelId string) string {

	url := "https://www.googleapis.com/youtube/v3/channels"
 
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Fatal(err)
	}

	// なぜかchannelIdにダブルクォートが含まれるためダブルクォートをカットする。これをしないとparams.Add("id", channelId)のchannelIdでダブルクォート込みで検索してしまい値が取れない。
	if strings.Contains(channelId, "\"") {
		channelId = strings.Replace(channelId, "\"", "", -1)
	}
	
	//クエリパラメータ
	params := request.URL.Query()
	params.Add("key", APIKey)
	params.Add("id", channelId)
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
	fmt.Println("body=", string(body))

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

	channelId, videoId := get_youtube_id(APIKey)
	viewCount := get_youtube_viewCount(APIKey, videoId)
	subscriberCount := get_youtube_subscriberCount(APIKey, channelId)

	return events.APIGatewayProxyResponse{
 
		 Body:       fmt.Sprintf(channelId, videoId, viewCount, subscriberCount),
		 StatusCode: 200,
	 }, nil
 }

func main() {

    lambda.Start(handler)
}