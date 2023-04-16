// レスポンスを全表示（初期）

package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const APIKey = "AIzaSyDAKFWvZl0y9DQI47F2tEaJPncMvnDPdkE"


func search_youtube_list(APIKey string) string{
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
 
	// fmt.Println(string(body))
	return string(body)
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