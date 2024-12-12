package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"golang.org/x/net/context"
	"io"
)

func IndexPost(post model.Post, content model.Content) error {
	// Create a document
	doc := map[string]interface{}{
		"title":      post.Title,
		"user_id":    post.UserID,
		"text":       content.Text,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
	}
	// Serialize the document to JSON
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	// Create an index request
	req := esapi.IndexRequest{
		Index:      "posts",
		DocumentID: fmt.Sprintf("%d", post.ID),
		Body:       bytes.NewReader(jsonData),
		Refresh:    "true",
	}
	// Perform the request with the client
	res, err := req.Do(context.Background(), ES)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Log.Error(nil, err.Error())
		}
	}(res.Body)
	return nil
}
