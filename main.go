/*Author: Bobby Williams
---------------------
| Word Counter v1.0 |
--------------------
This program counts the number of total words/non-whitespace in a file as well as # of instances per word/non-whitespace displaying the top 3 words/non-white space in the file as well as the total.

usage: wordcount.exe <filename>
*/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var mutex sync.Mutex
var wg sync.WaitGroup

func main() {
	switch {
	case len(os.Args) < 2:
		log.Fatal("Missing data file to process, exiting..")
	case len(os.Args) > 2:
		log.Fatal("Too many arguments, exiting..")
	}

	start := time.Now()
	jar := make(map[string]int)

	file, err := os.Open(os.Args[1])
	checkError(err)
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	checkError(err)

	fileData := string(fileBytes)

	fileLines := strings.Split(fileData, "\n")

	totalWords := 0
	for _, i := range fileLines {
		wg.Add(1)
		go func(s string, t int) {
			words := strings.Split(s, " ")
			mutex.Lock()
			totalWords += len(words)
			for _, w := range words {
				w = strings.ToLower(w)
				if w == "" {
					continue
				}
				jar[w] = jar[w] + 1
			}
			mutex.Unlock()
			wg.Done()
		}(i, totalWords)
	}
	wg.Wait()

	highest := make(map[string]map[string]int)

	topCounter, secondCounter, thirdCounter := 0, 0, 0

	for k, v := range jar {
		fmt.Printf("There are %v instances of '%v'\n", v, k)
		switch {
		case v > topCounter:
			topCounter = v
			highest["top"] = map[string]int{k: v}
		case v > secondCounter:
			secondCounter = v
			highest["second"] = map[string]int{k: v}
		case v > thirdCounter:
			thirdCounter = v
			highest["third"] = map[string]int{k: v}
		}
	}

	fmt.Println("--------------------------------")
	for k, v := range highest["top"] {
		fmt.Printf("The top word is: '%v' with %v total instances\n", k, v)
	}
	for k, v := range highest["second"] {
		fmt.Printf("The second highest word is: '%v' with %v total instances\n", k, v)
	}

	for k, v := range highest["third"] {
		fmt.Printf("The third highest word is: '%v' with %v total instances\n", k, v)
	}
	fmt.Printf("** Total words/non-whitespace = %v ***\n", totalWords)
	fmt.Println("--------------------------------")
	elapsed := time.Since(start)
	fmt.Printf("Execution time = %v\n", elapsed.String())
	fmt.Println("--------------------------------")
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
