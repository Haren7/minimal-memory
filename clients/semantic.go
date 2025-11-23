package clients

import (
	"context"
	"fmt"
	"log"
	"memory/internal/conversation"
	"memory/internal/embedding"
	"memory/internal/memory"
	"memory/internal/persistence/rdbms"
	"memory/internal/persistence/vector"
	"memory/internal/summarizer"
	"memory/types"

	"github.com/google/uuid"
)

type semanticMemoryClient struct {
	memoryService       memory.SemanticServiceInterface
	conversationService conversation.ConversationServiceInterface
	config              SemanticMemoryClientConfig
}

func NewSemanticMemoryClient(config SemanticMemoryClientConfig) (SemanticMemoryClient, error) {
	if config.OpenAIApiKey == "" {
		log.Printf("[ERROR] NewSemanticMemoryClient: OpenAI API key is required but was not provided")
		return nil, fmt.Errorf("error openai api key is required")
	}
	embeddingService := embedding.NewOpenAIService(config.OpenAIApiKey)
	summarizerService := summarizer.NewNoOpService()
	duckdbClient, err := rdbms.NewDuckDBClient()
	if err != nil {
		log.Printf("[ERROR] NewSemanticMemoryClient: Failed to connect to DuckDB - %v", err)
		return nil, fmt.Errorf("error connecting to duckdb")
	}
	conversationRepo := rdbms.NewConversationRepo(duckdbClient.GetDB())
	conversationService := conversation.NewConversationService(conversationRepo)
	memoryRepo := rdbms.NewMemoryRepo(duckdbClient.GetDB())
	faissMemoryRepo := rdbms.NewFaissMemoryRepo(duckdbClient.GetDB())
	// chromemDB := vector.NewChromem()
	// vectorMemoryRepo := vector.NewChromemMemoryRepo(chromemDB, embeddingService)
	faiss := vector.NewFaissClient()
	vectorMemoryRepo := vector.NewFaissMemoryRepo(faiss, embeddingService, faissMemoryRepo)
	memoryService := memory.NewSemanticService(vectorMemoryRepo, memoryRepo, conversationRepo, summarizerService)
	return &semanticMemoryClient{
		config:              config,
		memoryService:       memoryService,
		conversationService: conversationService,
	}, nil
}

func (r *semanticMemoryClient) Store(ctx context.Context, input types.StoreSemanticMemoryInput) (types.StoreSemanticMemoryOutput, error) {
	if input.Query == "" || input.Response == "" {
		log.Printf("[ERROR] Store: Query and response are required but one or both were empty (query: %q, response: %q)", input.Query, input.Response)
		return types.StoreSemanticMemoryOutput{}, fmt.Errorf("query and response are required")
	}
	if input.ConversationID == "" {
		log.Printf("[ERROR] Store: Conversation ID is required but was empty")
		return types.StoreSemanticMemoryOutput{}, fmt.Errorf("conversation id is required")
	}
	conversationID, err := uuid.Parse(input.ConversationID)
	if err != nil {
		log.Printf("[ERROR] Store: Invalid conversation ID format - %q, error: %v", input.ConversationID, err)
		return types.StoreSemanticMemoryOutput{}, fmt.Errorf("invalid conversation id")
	}
	exists, err := r.conversationService.Exists(ctx, conversationID)
	if err != nil {
		log.Printf("[ERROR] Store: Failed to check if conversation exists (conversationID: %s) - %v", conversationID, err)
		return types.StoreSemanticMemoryOutput{}, fmt.Errorf("error checking if conversation exists")
	}
	if !exists {
		log.Printf("[ERROR] Store: Conversation does not exist (conversationID: %s)", conversationID)
		return types.StoreSemanticMemoryOutput{}, fmt.Errorf("conversation does not exist")
	}
	id, err := r.memoryService.Store(ctx, conversationID, input.Query, input.Response)
	if err != nil {
		log.Printf("[ERROR] Store: Failed to store memory (conversationID: %s) - %v", conversationID, err)
		return types.StoreSemanticMemoryOutput{}, fmt.Errorf("error storing memory")
	}
	return types.StoreSemanticMemoryOutput{
		MemoryID: id.String(),
	}, nil
}

func (r *semanticMemoryClient) Retrieve(ctx context.Context, input types.RetrieveSemanticMemoryInput) (types.RetrieveSemanticMemoryOutput, error) {
	if input.ConversationID == "" {
		log.Printf("[ERROR] Retrieve: Conversation ID is required but was empty")
		return types.RetrieveSemanticMemoryOutput{}, fmt.Errorf("conversation id is required")
	}
	conversationID, err := uuid.Parse(input.ConversationID)
	if err != nil {
		log.Printf("[ERROR] Retrieve: Invalid conversation ID format - %q, error: %v", input.ConversationID, err)
		return types.RetrieveSemanticMemoryOutput{}, fmt.Errorf("invalid conversation id")
	}
	exists, err := r.conversationService.Exists(ctx, conversationID)
	if err != nil {
		log.Printf("[ERROR] Retrieve: Failed to check if conversation exists (conversationID: %s) - %v", conversationID, err)
		return types.RetrieveSemanticMemoryOutput{}, fmt.Errorf("error checking if conversation exists")
	}
	if !exists {
		log.Printf("[ERROR] Retrieve: Conversation does not exist (conversationID: %s)", conversationID)
		return types.RetrieveSemanticMemoryOutput{}, fmt.Errorf("conversation does not exist")
	}
	var topK int
	if input.TopK == 0 {
		topK = 10
	}
	retrievedMemories, err := r.memoryService.Retrieve(ctx, conversationID, r.config.ContextWindowSize)
	if err != nil {
		log.Printf("[ERROR] Retrieve: Failed to retrieve memories (conversationID: %s, contextWindowSize: %d) - %v", conversationID, r.config.ContextWindowSize, err)
		return types.RetrieveSemanticMemoryOutput{}, fmt.Errorf("error retrieving memories")
	}
	retrievedSimilarMemories, err := r.memoryService.RetrieveSimilar(ctx, conversationID, input.Query, topK)
	if err != nil {
		log.Printf("[ERROR] Retrieve: Failed to retrieve similar memories (conversationID: %s, query: %q, topK: %d) - %v", conversationID, input.Query, topK, err)
		return types.RetrieveSemanticMemoryOutput{}, fmt.Errorf("error retrieving similar memories")
	}
	var memories []types.Memory
	for _, memory := range retrievedMemories {
		memories = append(memories, types.Memory{
			ID:        memory.ID.String(),
			Query:     memory.Query,
			Response:  memory.Response,
			CreatedAt: memory.CreatedAt,
		})
	}
	var similarMemories []types.SemanticMemory
	for _, memory := range retrievedSimilarMemories {
		similarMemories = append(similarMemories, types.SemanticMemory{
			ID:        memory.ID.String(),
			Query:     memory.Query,
			Response:  memory.Response,
			CreatedAt: memory.CreatedAt,
		})
	}
	return types.RetrieveSemanticMemoryOutput{
		Memories:        memories,
		SimilarMemories: similarMemories,
	}, nil
}

func (r *semanticMemoryClient) RegisterConversation(ctx context.Context, input types.RegisterConversationInput) (types.RegisterConversationOutput, error) {
	if input.Agent == "" || input.User == "" {
		log.Printf("[ERROR] RegisterConversation: Agent and user are required but one or both were empty (agent: %q, user: %q)", input.Agent, input.User)
		return types.RegisterConversationOutput{}, fmt.Errorf("agent and user are required")
	}
	id, err := r.conversationService.Create(ctx, input.Agent, input.User)
	if err != nil {
		log.Printf("[ERROR] RegisterConversation: Failed to create conversation (agent: %q, user: %q) - %v", input.Agent, input.User, err)
		return types.RegisterConversationOutput{}, fmt.Errorf("error creating conversation")
	}
	return types.RegisterConversationOutput{
		ConversationID: id.String(),
	}, nil
}
