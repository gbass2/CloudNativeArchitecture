// Find the top K most common words in a text document.
// Input path: location of the document, K top words
// Output: Slice of top K words
// For this excercise, word is defined as characters separated by a whitespace

// Note: You should use `checkError` to handle potential errors.

package textproc

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
)

func topWords(path string, K int) []WordCount {
	if K <= 0 {
		checkError(errors.New("Value for K is less than 0"))
	}

	wordMap := map[string]int{} // Creating the map to store the words and the word count.

	// Reading the contents of the file.
	file, err := os.Open("passage")
	checkError(err)

	// Converting the contents of the text file to a Slice with each element a different word
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	defer file.Close()

	// Going word by word and if the word is in the map then increment the integer, if not then add the word to the map with a value of 1.
	for _, word := range words {
		_, ok := wordMap[word]
		if ok {
			wordMap[word]++
		} else {
			wordMap[word] = 1
		}
	}

	// Covnert each map element to WordCount objects and place in a slice.
	var wordSlice []WordCount

	for k, v := range wordMap {
		wordSlice = append(wordSlice, WordCount{Word: k, Count: v})
	}

	// Sort the Slice
	sortWordCounts(wordSlice)

	// Returning the top Kth elements of the slice. If K > length of the slice then return all of the slice.
	if K <= len(wordSlice) {
		wordSlice = wordSlice[:K]
	}

	fmt.Println(wordSlice)
	return wordSlice
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
