package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Fetcher interface {
	FetchUrls() []string
}

type UrlFetcher struct {
	ImageNumber int
	QueryString string
	ApiKey      string
	SerpApiUrl  string
}

func NewUrlFetcher(imageNumber int, queryString, apiKey, SerpApiUrl string) Fetcher {
	return &UrlFetcher{
		ImageNumber: imageNumber,
		QueryString: queryString,
		ApiKey:      apiKey,
		SerpApiUrl:  SerpApiUrl,
	}
}

func (f *UrlFetcher) FetchUrls() []string {
	ijn := 0
	imageUrls := make([]string, 0, f.ImageNumber)
	for f.ImageNumber > 0 {
		// The API only returns 100 images at a time, so we need to make multiple requests
		// to get the desired number of images
		// if google does not have enough images, it will return less than 100 or even 0
		url := fmt.Sprintf("%s?engine=google_images&q=%s&api_key=%s&ijn=%d", f.SerpApiUrl, f.QueryString, f.ApiKey, ijn)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}

		urls := extractUrlFromJson(resp)
		if f.ImageNumber > 100 {
			imageUrls = append(imageUrls, urls...)
		} else {
			if len(urls) < f.ImageNumber {
				imageUrls = append(imageUrls, urls...)
			} else {
				imageUrls = append(imageUrls, urls[:f.ImageNumber]...)
			}
		}
		f.ImageNumber -= 100
		ijn++
		resp.Body.Close()
	}
	return imageUrls
}

func extractUrlFromJson(resp *http.Response) []string {
	urls := make([]string, 0, 100)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}
	imagesResults, ok := data["images_results"].([]interface{})
	if !ok {
		// if google does not have enough images, it will return less than 100 or even 0 and the status code will be 200
		if resp.StatusCode == 200 {
			fmt.Println("not_enough_results")
			return nil
		}
		fmt.Println("Error extracting images_results")
		return nil
	}
	for _, image := range imagesResults {
		img, ok := image.(map[string]interface{})
		if !ok {
			fmt.Println("Error parsing image data")
			continue
		}
		originalURL, ok := img["original"].(string)
		if !ok {
			fmt.Println("Error extracting original URL")
			continue
		}
		urls = append(urls, originalURL)
	}
	return urls
}
