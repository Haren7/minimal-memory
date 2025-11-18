package embedding

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	openAiClient *openai.Client
}

func NewOpenAIService(openAiClient *openai.Client) Service {
	return &OpenAIService{
		openAiClient: openAiClient,
	}
}

func (r *OpenAIService) EmbedOne(ctx context.Context, text string) (Embedding, error) {
	input := openai.EmbeddingRequest{
		Input: text,
		Model: openai.SmallEmbedding3,
	}
	response, err := r.openAiClient.CreateEmbeddings(ctx, input)
	if err != nil {
		return Embedding{}, err
	}
	embedding := response.Data[0].Embedding
	return Embedding{
		Dim:    len(embedding),
		Vector: embedding,
	}, nil
}

func (r *OpenAIService) EmbedMany(ctx context.Context, texts []string) ([]Embedding, error) {
	input := openai.EmbeddingRequest{
		Input: texts,
		Model: openai.SmallEmbedding3,
	}
	resp, err := r.openAiClient.CreateEmbeddings(ctx, input)
	if err != nil {
		return []Embedding{}, err
	}
	var embeddings []Embedding
	for _, emb := range resp.Data {
		embeddings = append(embeddings, Embedding{
			Dim:    len(emb.Embedding),
			Vector: emb.Embedding,
		})
	}
	return embeddings, nil
}
