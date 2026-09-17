package index

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

type Document struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Category    string `json:"category"`
	Requester   string `json:"requester"`
	AssigneeID  string `json:"assignee_id"`
	UpdatedAt   string `json:"updated_at"`
}

type Store struct {
	client meilisearch.ServiceManager
	index  meilisearch.IndexManager
	name   string
}

func New(host, apiKey, indexName string) (*Store, error) {
	client := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))
	s := &Store{client: client, name: indexName}
	return s, nil
}

func (s *Store) Ensure(ctx context.Context) error {
	_, err := s.client.CreateIndexWithContext(ctx, &meilisearch.IndexConfig{
		Uid:        s.name,
		PrimaryKey: "id",
	})
	// index may already exist
	if err != nil {
		// continue; GetIndex will confirm
	}
	s.index = s.client.Index(s.name)

	_, err = s.index.UpdateSearchableAttributesWithContext(ctx, &[]string{
		"title", "description", "category", "requester", "assignee_id", "status", "priority", "id",
	})
	if err != nil {
		return fmt.Errorf("searchable attrs: %w", err)
	}
	_, err = s.index.UpdateFilterableAttributesWithContext(ctx, &[]string{
		"status", "priority", "assignee_id", "category",
	})
	if err != nil {
		return fmt.Errorf("filterable attrs: %w", err)
	}
	return nil
}

func (s *Store) Upsert(ctx context.Context, docs ...Document) error {
	if len(docs) == 0 {
		return nil
	}
	task, err := s.index.AddDocumentsWithContext(ctx, docs, "id")
	if err != nil {
		return err
	}
	_, err = s.client.WaitForTaskWithContext(ctx, task.TaskUID, 100*time.Millisecond)
	return err
}

type Hit struct {
	Document
}

func (s *Store) Search(ctx context.Context, query string, limit int64) ([]Hit, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	res, err := s.index.SearchWithContext(ctx, query, &meilisearch.SearchRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, 0, err
	}
	hits := make([]Hit, 0, len(res.Hits))
	for _, raw := range res.Hits {
		b, err := json.Marshal(raw)
		if err != nil {
			continue
		}
		var d Document
		if err := json.Unmarshal(b, &d); err != nil {
			continue
		}
		hits = append(hits, Hit{Document: d})
	}
	return hits, res.EstimatedTotalHits, nil
}
