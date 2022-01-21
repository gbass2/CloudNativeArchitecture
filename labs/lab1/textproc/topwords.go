// Find the top K most common words in a text document.
// Input path: location of the document, K top words
// Output: Slice of top K words
// For this excercise, word is defined as characters separated by a whitespace

// Note: You should use `checkError` to handle potential errors.

package textproc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func topWords(path string, K int) []WordCount {
	wordMap := map[string]int{} // Creating the map to store the words and the word count.

	// Reading the contents of the file.
	file, err := os.Open("passage")
	checkError(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	// Going line by line of the read in contents of the file.
	for _, each_ln := range text {
		input := strings.Fields(each_ln)

		// Going word by word and if the word is in the map then increment the integer, if not then add the word to the map with a value of 1.
		for _, word := range input {
			_, ok := wordMap[word]
			if ok {
				wordMap[word] = wordMap[word] + 1
			} else {
				wordMap[word] = 1
			}
		}
	}
	fmt.Println(wordMap)

}

//--------------- DO NOT MODIFY----------------!

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

// Method to convert struct to string format
func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.

func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
