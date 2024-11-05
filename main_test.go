package main

import (
	"fmt"
	"testing"
)

// var q = "The Five Elements of Effective Thinking (Unabridged)"

// var q = "The Curse of the High IQ"

// var q = "Killing Crazy Horse: The Merciless Indian Wars in America"

// var q = "Yongey Mingyur Rinpoche, Helen Tworkov"

var q = "Robert Greene"

func TestOpenLibrary(t *testing.T) {
	// Open Library API Test
	fmt.Println("Open Library API Results:")
	olAudiobooks, err := searchOpenLibrary(q)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, book := range olAudiobooks {
			fmt.Printf("Title: %s\nAuthors: %v\nInfoLink: %s\n\n",
				book.Title, book.Authors, book.InfoLink)
		}
	}
}

func TestITunesSearch(t *testing.T) {
	// iTunes Search API Test
	fmt.Println("iTunes Search API Results:")
	itAudiobooks, err := searchITunes(q)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, book := range itAudiobooks {
			fmt.Printf("Title: %s\nAuthors: %v\nDescription: %s\nInfoLink: %s\nCoverImage: %s\n\n",
				book.Title, book.Authors, book.Description, book.InfoLink, book.CoverImage)
		}
	}
}

func TestLibriVoxSearch(t *testing.T) {
	// LibriVox API Test
	fmt.Println("LibriVox API Results:")
	lvAudiobooks, err := searchLibriVox(q)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, book := range lvAudiobooks {
			fmt.Printf("Title: %s\nAuthors: %v\nDescription: %s\nInfoLink: %s\nCoverImage: %s\n\n",
				book.Title, book.Authors, book.Description, book.InfoLink, book.CoverImage)
		}
	}
}

func TestArchiveOrgSearch(t *testing.T) {
	// Archive.org API Test
	fmt.Println("Archive.org API Results:")
	aoAudiobooks, err := searchArchiveOrg(q)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, book := range aoAudiobooks {
			fmt.Printf("Title: %s\nAuthors: %v\nDescription: %s\nInfoLink: %s\nCoverImage: %s\n\n",
				book.Title, book.Authors, book.Description, book.InfoLink, book.CoverImage)
		}
	}
}
