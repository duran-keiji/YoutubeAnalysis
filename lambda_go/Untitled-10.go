// チャンネルIDと動画IDと視聴回数取得（初期）

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
	// Kind     string `json:"kind"`
	// Etag       string    `json:"etag"`
	NextPageToken       string    `json:"nextPageToken"`
	// RegionCode       string    `json:"regionCode"`
	// PageInfo       Info    `json:"pageInfo"`
	Items []*SearchItem `json:"items"`
}

// [2] JSONと同様のネスト構造で構造体を定義
type SearchItem struct {
	// Kind  string    `json:"kind"`
	// Etag string `json:"etag"`
	ID SearchID `json:"id"`
	Snippet SearchSnippet `json:"snippet"`
	ChannelTitle string `json:"channelTitle"`
	// LiveBroadcastContent string `json:"liveBroadcastContent"`
	// PublishTime string `json:"publishTime"`
}

type SearchID struct {
	// Kind string `json:"kind"`
	VideoId string `json:"videoId"`
}

type SearchSnippet struct {
	// PublishedAt string `json:"publishedAt"`
	ChannelId string `json:"channelId"`
	Title string `json:"title"`
	// Description string `json:"description"`
	// Thumbnails Thumbnails `json:"thumbnails"`
}

// type Thumbnails struct {
// 	Default Default `json:"default"`
// 	Medium Medium `json:"medium"`
// 	High High `json:"high"`
// }

// type Default struct {
// 	Url string `json:"url"`
// 	Width int `json:"width"`
// 	Height int `json:"height"`
// }

// type Medium struct {
// 	Url string `json:"url"`
// 	Width int `json:"width"`
// 	Height int `json:"height"`
// }

// type High struct {
// 	Url string `json:"url"`
// 	Width int `json:"width"`
// 	Height int `json:"height"`
// }

// type Info struct {
// 	TotalResults string `json:"totalResults"`
// 	ResultsPerPage string `json:"resultsPerPage"`
// }


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
	params.Add("maxResults", "1")

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
	// Kind     string `json:"kind"`
	// Etag       string    `json:"etag"`
	Items []*VideoItem `json:"items"`
	// PageInfo     Info   `json:"pageInfo"`
}

// [2] JSONと同様のネスト構造で構造体を定義
type VideoItem struct {
	// Kind  string    `json:"kind"`
	// Etag string `json:"etag"`
	ID string `json:"id"`
	Statistics VideoStatistic `json:"statistics"`
}

type VideoStatistic struct {
	ViewCount string `json:"viewCount"`
	// LikeCount string `json:"likeCount"`
	// FavoriteCount string `json:"favoriteCount"`
	// CommentCount string `json:"commentCount"`
}

// type Info struct {
// 	TotalResults int `json:"totalResults"`
// 	ResultsPerPage int `json:"resultsPerPage"`
// }



func get_youtube_viewCount(APIKey string, videoId string) string {

	// なぜかvideoIdにダブルクォートが含まれるためダブルクォートをカットする。これをしないとparams.Add("id", videoId)のvideoIdでダブルクォート込みで検索してしまい値が取れない。
	if strings.Contains(videoId, "\"") {
		videoId = strings.Replace(videoId, "\"", "", -1)
	}

	url := "https://www.googleapis.com/youtube/v3/videos"
 
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Fatal(err)
	}
	
	//クエリパラメータ
	params := request.URL.Query()
	params.Add("key", APIKey)
	params.Add("id", videoId)
	params.Add("part", "statistics") 
	params.Add("maxResults", "1")

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



func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	channelId, videoId := get_youtube_id(APIKey)
	viewCount := get_youtube_viewCount(APIKey, videoId)

	return events.APIGatewayProxyResponse{
 
		 Body:       fmt.Sprintf(channelId, videoId, viewCount),
		 StatusCode: 200,
	 }, nil
 }

func main() {

    lambda.Start(handler)
}