package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/scanner"
)

var MaxWorker = 5
var finalCount = make(chan map[string]uint, 1)

func CountWord(str string) map[string]uint {
	var result = map[string]uint{}
	var s scanner.Scanner
	s.Init(strings.NewReader(str))

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if tok == scanner.Ident {
			result[s.TokenText()] += 1
		}
	}
	return result
}

func handle(strings chan string, resultChannel chan map[string]uint, wg *sync.WaitGroup) {
	defer wg.Done()

	for str := range strings {
		result := CountWord(str)
		fmt.Println(result)
		resultChannel <- result
	}
}

func cummulateCount(resultChannel chan map[string]uint, wg *sync.WaitGroup) {
	defer wg.Done()

	finalResult := map[string]uint{}
	for result := range resultChannel {
		fmt.Println("Pick up another result")
		for word, count := range result {
			finalResult[word] += count
		}
	}

	finalCount <- finalResult
}

func main() {
	// Extract files name
	var files []string
	root := "./data"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".txt" {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	var stringsChannel = make(chan string)
	var resultChannel = make(chan map[string]uint)
	var countGroup = sync.WaitGroup{}
	var cummulateGroup = sync.WaitGroup{}

	// worker for cummulate count results from each file
	cummulateGroup.Add(1)
	go cummulateCount(resultChannel, &cummulateGroup)

	// Initize pool of workers (count word in each file)
	countGroup.Add(MaxWorker)
	for i := 0; i < MaxWorker; i++ {
		go handle(stringsChannel, resultChannel, &countGroup)
	}

	// Feed data to workers
	for _, file := range files {
		dat, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		stringsChannel <- string(dat)
	}
	close(stringsChannel)

	// Wait for worker to finish off their work
	countGroup.Wait()
	close(resultChannel)
	cummulateGroup.Wait()

	fmt.Println("-------------------------------------------------------")
	fmt.Println("CUMMULATIVE COUNT: : ", <-finalCount)
}
