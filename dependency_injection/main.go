package main

import (
	configurer "github.com/phluan/GrabGoTrainingWeek5Assignment/configurer"
	serializer "github.com/phluan/GrabGoTrainingWeek5Assignment/serializer"
	services "github.com/phluan/GrabGoTrainingWeek5Assignment/services"
	"log"
	"net/http"
)

//TODO: how to separate API logic, business logic and response format logic
func main() {
	httpClient := http.DefaultClient
	serializer := serializer.JsonSerializer{}
	configurer, configurerErr := configurer.New(configurer.WithHttpClient(httpClient), configurer.WithSerializer(&serializer))
	if configurerErr != nil {
		log.Println("Cannot setup configurer: %s", configurerErr)
		return
	}

	http.HandleFunc("/postWithComments", func(writer http.ResponseWriter, request *http.Request) {
		postWithCommentsSerialized, err := services.GetPostWithComments(configurer)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(500)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(postWithCommentsSerialized)
	})

	log.Println("httpServer starts ListenAndServe at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
