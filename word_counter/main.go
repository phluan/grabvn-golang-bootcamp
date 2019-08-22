package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/scanner"
)

var MaxWorker = 100
var finalCount = make(chan map[string]uint, 1)

func breakUpFileToLines(file *os.File, stringsChannel chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stringsChannel <- scanner.Text()
	}
}

func CountWord(str string) map[string]uint {
	var result = map[string]uint{}
	var s scanner.Scanner
	s.Init(strings.NewReader(str))

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if tok == scanner.Ident {
			result[strings.ToLower(s.TokenText())] += 1
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

func cummulate(resultChannel chan map[string]uint, wg *sync.WaitGroup) {
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
	var breakUpFileGroup = sync.WaitGroup{}
	var countGroup = sync.WaitGroup{}
	var cummulateGroup = sync.WaitGroup{}

	// Initize pool of workers (count word in each file)
	countGroup.Add(MaxWorker)
	for i := 0; i < MaxWorker; i++ {
		go handle(stringsChannel, resultChannel, &countGroup)
	}

	// worker for cummulate count results from each file
	cummulateGroup.Add(1)
	go cummulate(resultChannel, &cummulateGroup)

	// Feed data to workers
	for _, file := range files {
		file, err := os.Open(file)
		if err != nil {
			panic(err)
		}

		defer file.Close()
		breakUpFileGroup.Add(1)
		go breakUpFileToLines(file, stringsChannel, &breakUpFileGroup)
	}
	// Wait for worker to finish off their work
	breakUpFileGroup.Wait()
	close(stringsChannel)

	countGroup.Wait()
	close(resultChannel)

	cummulateGroup.Wait()

	fmt.Println("-------------------------------------------------------")
	fmt.Println("CUMMULATIVE COUNT: : ", <-finalCount)
}

// 1 workers
// real	0m0.303s
// user	0m0.313s
// sys	0m0.204s

// 2 workers
// real	0m0.284s
// user	0m0.283s
// sys	0m0.189s

// 5 workers
// real	0m0.292s
// user	0m0.321s
// sys	0m0.219s

// 100 workers
// real	0m0.279s
// user	0m0.311s
// sys	0m0.197s
