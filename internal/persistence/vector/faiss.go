package vector

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/haren7/minimal-memory/internal/embedding"

	"github.com/DataIntelligenceCrew/go-faiss"
)

var ErrindexDoesNotExist = errors.New("faiss index does not exist")

type FaissSearchResponse struct {
	Distances []float32
	Ids       []int64
}

type FaissClient struct {
	dir                   string
	conversationIDVsIndex map[string]*faiss.IndexImpl
}

func NewFaissClient() *FaissClient {
	return &FaissClient{
		conversationIDVsIndex: make(map[string]*faiss.IndexImpl),
	}
}

func (r *FaissClient) Index(ctx context.Context, conversationID string, id int, embedding embedding.Embedding) error {
	index, exists := r.conversationIDVsIndex[conversationID]
	if !exists {
		newIndex, err := faiss.IndexFactory(embedding.Dim, "IDMap,Flat", 1)
		if err != nil {
			return fmt.Errorf("error creating idmap + flat index with dim %d - %w", embedding.Dim, err)
		}
		r.conversationIDVsIndex[conversationID] = newIndex
		index = newIndex
	}
	err := index.AddWithIDs(embedding.Vector, []int64{int64(id)})
	if err != nil {
		return fmt.Errorf("error adding vector to index: %d - %w", id, err)
	}
	return nil
}

func (r *FaissClient) Search(ctx context.Context, conversationID string, query embedding.Embedding, topK int) (FaissSearchResponse, error) {
	index, exists := r.conversationIDVsIndex[conversationID]
	if !exists {
		return FaissSearchResponse{}, ErrindexDoesNotExist
	}
	distances, labels, err := index.Search(query.Vector, int64(topK))
	if err != nil {
		return FaissSearchResponse{}, fmt.Errorf("error searching index: %w", err)
	}
	var validIds []int64
	for _, label := range labels {
		if label != -1 {
			validIds = append(validIds, label)
		}
	}
	return FaissSearchResponse{
		Distances: distances,
		Ids:       validIds,
	}, nil
}

func (r *FaissClient) Mount(dir string, files map[string]io.Reader) error {
	conversationIDVsIndex := make(map[string]*faiss.IndexImpl)
	for fileName, reader := range files {
		conversationID := strings.TrimSuffix(fileName, ".index")
		writePath := r.getFilePath(dir, conversationID)
		bytes, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", fileName, err)
		}
		err = os.WriteFile(writePath, bytes, 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s: %w", fileName, err)
		}
		index, err := faiss.ReadIndex(writePath, 0)
		if err != nil {
			return fmt.Errorf("error reading index from file %s: %w", fileName, err)
		}
		conversationIDVsIndex[conversationID] = index
	}
	r.conversationIDVsIndex = conversationIDVsIndex
	return nil
}

func (r *FaissClient) Export(dir string) ([]os.File, error) {
	var files []os.File
	for conversationID, index := range r.conversationIDVsIndex {
		filePath := r.getFilePath(dir, conversationID)
		err := faiss.WriteIndex(index, filePath)
		if err != nil {
			return nil, fmt.Errorf("error exporting index for conversation id %s: %w", conversationID, err)
		}
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
		}
		files = append(files, *file)
	}
	return files, nil
}

func (r *FaissClient) getFilePath(dir string, conversationID string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.index", conversationID))
}
