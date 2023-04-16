// チャンネルIDと動画ID取得

package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const APIKey = "AIzaSyDAKFWvZl0y9DQI47F2tEaJPncMvnDPdkE"

// {
// 	"kind": "youtube#searchListResponse",
// 	"etag": "ZrjEiOaLy1kD94AwBPLsflr3N_M",
// 	"nextPageToken": "CAEQAA",
// 	"regionCode": "JP",
// 	"pageInfo": {
// 	  "totalResults": 1000000,
// 	  "resultsPerPage": 1
// 	},
// 	"items": [
// 	  {
// 		"kind": "youtube#searchResult",
// 		"etag": "blz4aX-nVla7yAGTe1zJziUZ2d8",
// 		"id": {
// 		  "kind": "youtube#video",
// 		  "videoId": "mMx_xlrbuhg"
// 		},
// 		"snippet": {
// 		  "publishedAt": "2020-04-15T09:17:19Z",
// 		  "channelId": "UC0LwRQO6woCv1lvX4X0QGlg",
// 		  "title": "［和訳］Girl in the Mirror - Bebe Rexha",
// 		  "description": "https://youtu.be/4sw_3Ru5_Sg.",
// 		  "thumbnails": {
// 			"default": {
// 			  "url": "https://i.ytimg.com/vi/mMx_xlrbuhg/default.jpg",
// 			  "width": 120,
// 			  "height": 90
// 			},
// 			"medium": {
// 			  "url": "https://i.ytimg.com/vi/mMx_xlrbuhg/mqdefault.jpg",
// 			  "width": 320,
// 			  "height": 180
// 			},
// 			"high": {
// 			  "url": "https://i.ytimg.com/vi/mMx_xlrbuhg/hqdefault.jpg",
// 			  "width": 480,
// 			  "height": 360
// 			}
// 		  },
// 		  "channelTitle": "blue",
// 		  "liveBroadcastContent": "none",
// 		  "publishTime": "2020-04-15T09:17:19Z"
// 		}
// 	  }
// 	]
//   }

  // [1] 構造体のフィールドとJSONのキーを紐付け
type Search struct {
	Kind     string `json:"kind"`
	Etag       string    `json:"etag"`
	NextPageToken       string    `json:"nextPageToken"`
	RegionCode       string    `json:"regionCode"`
	PageInfo       Info    `json:"pageInfo"`
	Items []*Item `json:"items"`
}

// [2] JSONと同様のネスト構造で構造体を定義
type Item struct {
	Kind  string    `json:"kind"`
	Etag string `json:"etag"`
	ID ID `json:"id"`
	Snippet Snippet `json:"snippet"`
	ChannelTitle string `json:"channelTitle"`
	LiveBroadcastContent string `json:"liveBroadcastContent"`
	PublishTime string `json:"publishTime"`
}

type ID struct {
	Kind string `json:"kind"`
	VideoId string `json:"videoId"`
}

type Snippet struct {
	PublishedAt string `json:"publishedAt"`
	ChannelId string `json:"channelId"`
	Title string `json:"title"`
	Description string `json:"description"`
	Thumbnails Thumbnails `json:"thumbnails"`
}

type Thumbnails struct {
	Default Default `json:"default"`
	Medium Medium `json:"medium"`
	High High `json:"high"`
}

type Default struct {
	Url string `json:"url"`
	Width int `json:"width"`
	Height int `json:"height"`
}

type Medium struct {
	Url string `json:"url"`
	Width int `json:"width"`
	Height int `json:"height"`
}

type High struct {
	Url string `json:"url"`
	Width int `json:"width"`
	Height int `json:"height"`
}

type Info struct {
	TotalResults string `json:"totalResults"`
	ResultsPerPage string `json:"resultsPerPage"`
}


func search_youtube_list(APIKey string) (string, string){
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

	var search Search
	json.Unmarshal(body, &search)
	fmt.Println(search.Items[0].Snippet.ChannelId)
	fmt.Println(search.Items[0].ID.VideoId)

	ChannelId, err := json.Marshal(search.Items[0].Snippet.ChannelId)
	if err != nil {
		log.Fatal(err)
	}

	VideoId, err := json.Marshal(search.Items[0].ID.VideoId)
	if err != nil {
		log.Fatal(err)
	}

	// ChannelId := search.Items[0].Snippet.ChannelId
	// VideoId := search.Items[0].ID.VideoId
 
	// fmt.Println(string(body))
	return string(ChannelId), string(VideoId)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	channelId, videoId := search_youtube_list(APIKey)

	return events.APIGatewayProxyResponse{
 
		 Body:       fmt.Sprintf(channelId, videoId),
		 StatusCode: 200,
	 }, nil
 }

func main() {

    lambda.Start(handler)
}