package main

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"io/ioutil"
	"net/http"
)

func main() {
	client, err := statsd.New("127.0.0.1:8125",
		statsd.WithNamespace("http_echo_client."), // prefix every metric with the app name
		statsd.WithTags([]string{"localhost"}),    // send the EC2 availability zone as a tag with every metric
		// add more options here...
	)
	if err != nil {
		fmt.Print("fail to connect to Datadog agent")
		return
	}

	failCount := 0
	successCount := 0

	for i := 0; i < 100; i++ {
		res, err := http.Get("http://localhost:8080/")
		if err != nil {
			fmt.Printf("Request fails: %s", err)
			return
		}

		defer res.Body.Close()
		if res.StatusCode != 200 {
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Printf("Request fails with error: %d : %s\n", res.StatusCode, body)
			client.Count("echo_failure", 1, []string{"env:localhost"}, 0.5)
			failCount += 1
		} else {
			fmt.Println("Request success")
			client.Count("echo_success", 1, []string{"env:localhost"}, 0.5)
			successCount += 1
		}
	}

	fmt.Printf("Fail/success ratio: %v", float64(failCount)/float64(successCount))
}
