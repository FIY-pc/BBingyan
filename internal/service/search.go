package service

import (
	"bytes"
	"encoding/json"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/es"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/FIY-pc/BBingyan/internal/model"
	"io"
)

func SearchPosts(dto dto.SearchPostDTO) ([]model.Post, error) {
	var buf bytes.Buffer
	// Build the query JSON
	queryJSON := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  dto.Query,
				"fields": []string{"title", "text"},
			},
		},
	}
	// Serialize the query JSON to the buffer
	if err := json.NewEncoder(&buf).Encode(queryJSON); err != nil {
		return nil, err
	}
	// Perform the search request
	from := (dto.Page - 1) * dto.PageSize
	res, err := es.ES.Search(
		es.ES.Search.WithIndex("posts"),
		es.ES.Search.WithBody(&buf),
		es.ES.Search.WithTrackTotalHits(true),
		es.ES.Search.WithPretty(),
		es.ES.Search.WithFrom(from),
		es.ES.Search.WithSize(dto.PageSize),
	)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Log.Error(nil, err.Error())
		}
	}(res.Body)
	// Parse the search result
	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source model.Post `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	// Decode the search result
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, err
	}
	// Extract the posts from the search result
	var posts []model.Post
	for _, hit := range searchResult.Hits.Hits {
		posts = append(posts, hit.Source)
	}
	return posts, nil
}
