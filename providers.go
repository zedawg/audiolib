package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Open Library API Function
func searchOpenLibrary(query string) ([]AudiobookResult, error) {
	baseURL := "https://openlibrary.org/search.json"
	params := url.Values{}
	params.Add("q", query)
	params.Add("subject", "audiobooks")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Docs []struct {
			Title   string   `json:"title"`
			Authors []string `json:"author_name"`
			Key     string   `json:"key"`
			Subject []string `json:"subject"`
		} `json:"docs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var audiobooks []AudiobookResult
	for _, doc := range result.Docs {
		// Check if 'Audiobooks' is in the subject list
		isAudiobook := false
		for _, subject := range doc.Subject {
			if subject == "Audiobooks" {
				isAudiobook = true
				break
			}
		}
		if isAudiobook {
			audiobooks = append(audiobooks, AudiobookResult{
				Title:    doc.Title,
				Authors:  doc.Authors,
				InfoLink: "https://openlibrary.org" + doc.Key,
			})
		}
	}

	return audiobooks, nil
}

// iTunes Search API Function
func searchITunes(query string) ([]AudiobookResult, error) {
	baseURL := "https://itunes.apple.com/search"
	params := url.Values{}
	params.Add("term", query)
	params.Add("media", "audiobook")
	params.Add("limit", "25")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			ArtistName        string `json:"artistName"`
			CollectionName    string `json:"collectionName"`
			TrackName         string `json:"trackName"`
			Description       string `json:"description"`
			ArtworkURL100     string `json:"artworkUrl100"`
			CollectionViewURL string `json:"collectionViewUrl"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var audiobooks []AudiobookResult
	for _, item := range result.Results {
		audiobooks = append(audiobooks, AudiobookResult{
			Title:       item.CollectionName,
			Authors:     []string{item.ArtistName},
			Description: item.Description,
			InfoLink:    item.CollectionViewURL,
			CoverImage:  item.ArtworkURL100,
		})
	}

	return audiobooks, nil
}

// LibriVox API Function
func searchLibriVox(query string) ([]AudiobookResult, error) {
	baseURL := "https://librivox.org/api/feed/audiobooks"
	params := url.Values{}
	params.Add("format", "json")
	params.Add("title", query)
	params.Add("limit", "10")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Books []struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			URL         string `json:"url_librivox"`
			Description string `json:"description"`
			Authors     []struct {
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"authors"`
			CoverImage string `json:"url_zip_file"`
		} `json:"books"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var audiobooks []AudiobookResult
	for _, book := range result.Books {
		var authors []string
		for _, author := range book.Authors {
			fullName := fmt.Sprintf("%s %s", author.FirstName, author.LastName)
			authors = append(authors, fullName)
		}
		audiobooks = append(audiobooks, AudiobookResult{
			Title:       book.Title,
			Authors:     authors,
			Description: book.Description,
			InfoLink:    book.URL,
			CoverImage:  book.CoverImage,
		})
	}

	return audiobooks, nil
}

// Archive.org API Function
func searchArchiveOrg(query string) ([]AudiobookResult, error) {
	baseURL := "https://archive.org/advancedsearch.php"
	params := url.Values{}
	params.Add("q", fmt.Sprintf("%s AND mediatype:(audio)", query))
	params.Add("fl[]", "identifier,title,creator,description")
	params.Add("rows", "10")
	params.Add("page", "1")
	params.Add("output", "json")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Response struct {
			Docs []struct {
				Identifier  string `json:"identifier"`
				Title       string `json:"title"`
				Creator     string `json:"creator"`
				Description string `json:"description"`
			} `json:"docs"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var audiobooks []AudiobookResult
	for _, doc := range result.Response.Docs {
		audiobooks = append(audiobooks, AudiobookResult{
			Title:       doc.Title,
			Authors:     []string{doc.Creator},
			Description: doc.Description,
			InfoLink:    fmt.Sprintf("https://archive.org/details/%s", doc.Identifier),
			CoverImage:  fmt.Sprintf("https://archive.org/services/img/%s", doc.Identifier),
		})
	}

	return audiobooks, nil
}
