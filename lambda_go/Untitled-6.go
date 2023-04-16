// 動画の視聴回数取得
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
// 	"kind": "youtube#videoListResponse",
// 	"etag": "KJNVbMj-AJp_G3TXR4kX4VAAVKI",
// 	"items": [
// 	  {
// 		"kind": "youtube#video",
// 		"etag": "VVC10I-E1sMTGnzJkCRSFFdWELQ",
// 		"id": "TqZ91Hw4rjs",
// 		"statistics": {
// 		  "viewCount": "276",
// 		  "likeCount": "14",
// 		  "favoriteCount": "0",
// 		  "commentCount": "2"
// 		}
// 	  }
// 	],
// 	"pageInfo": {
// 	  "totalResults": 1,
// 	  "resultsPerPage": 1
// 	}
//   }

// [1] 構造体のフィールドとJSONのキーを紐付け
type Video struct {
	Kind     string `json:"kind"`
	Etag       string    `json:"etag"`
	Items []*Item `json:"items"`
	PageInfo     Info   `json:"pageInfo"`
}

// [2] JSONと同様のネスト構造で構造体を定義
type Item struct {
	Kind  string    `json:"kind"`
	Etag string `json:"etag"`
	ID string `json:"id"`
	Statistics Statistic `json:"statistics"`
}

type Statistic struct {
	ViewCount string `json:"viewCount"`
	LikeCount string `json:"likeCount"`
	FavoriteCount string `json:"favoriteCount"`
	CommentCount string `json:"commentCount"`
}

type Info struct {
	TotalResults int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}



func search_youtube_list(APIKey string) string{
	url := "https://www.googleapis.com/youtube/v3/videos"
 
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Fatal(err)
	}
	
	//クエリパラメータ
	params := request.URL.Query()
	params.Add("key", APIKey)
	params.Add("id", "TqZ91Hw4rjs")
	// params.Add("part", "snippet") 
	params.Add("part", "statistics") 
	// params.Add("part", "contentDetails") 
	// params.Add("part", "fileDetails") 
	// params.Add("part", "liveStreamingDetails") 
	// params.Add("part", "processingDetails") 
	// params.Add("part", "recordingDetails") 
	// params.Add("part", "statistics") 
	// params.Add("part", "status") 
	// params.Add("part", "suggestions") 
	// params.Add("part", "topicDetails") 
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


	var video Video
	json.Unmarshal(body, &video)
	fmt.Println(video.Items[0].Statistics.ViewCount)

	count, err := json.Marshal(video.Items[0].Statistics.ViewCount)
	if err != nil {
		log.Fatal(err)
	}



	

	// // [3] 配列型のJSONデータを読み込む
	// videoData := make([]*Video, 0)
	// err = json.Unmarshal(body, &videoData)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	
	// for _, video := range videoData {
	// 	// [2] ネストした構造体のフィールドを呼び出す
	// 	fmt.Printf("%s (@%s)\n", video.Items.Statistics.ViewCount)
	// }

	// body, _ = json.Marshal(videoData)
	// fmt.Println(string(body)) 

	



	
 
	// fmt.Println(string(body))
	// return string(count)
	return string(count)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	youtube_search_list := search_youtube_list(APIKey)

	return events.APIGatewayProxyResponse{
 
		 Body:       fmt.Sprintf(youtube_search_list),
		 StatusCode: 200,
	 }, nil
 }

func main() {

    lambda.Start(handler)
}

//やりたいこと：チャンネル登録者数が少なくて再生回数の多い動画を洗い出す。
// ①search関数で検索ワードを基に動画IDとチャンネルIDを取得
// ②①の動画IDを基にvideoのlist関数のstatics使って再生回数を取得
// ③②のチャンネルIDをもとにchannnel関数のstatisticsで登録者数を取得
// ④再生回数とチャンネル登録者数を比較して効果の高い動画を見つける