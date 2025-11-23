package clients

import (
	"context"
	"fmt"
	"log"
	"memory/internal/cache"
	"memory/internal/conversation"
	"memory/internal/memory"
	"memory/internal/persistence/rdbms"
	"memory/internal/summarizer"
	"memory/types"

	"github.com/google/uuid"
)

type shortTermMemoryClient struct {
	memoryService       memory.ServiceInterface
	conversationService conversation.ConversationServiceInterface
	config              ShortTermMemoryClientConfig
}

func NewShortTermMemoryClient(config ShortTermMemoryClientConfig) (ShortTermMemoryClient, error) {
	duckdbClient, err := rdbms.NewDuckDBClient()
	if err != nil {
		log.Printf("[ERROR] NewShortTermMemoryClient: Failed to connect to DuckDB - %v", err)
		return nil, err
	}
	memoryRepo := cache.NewInMemMemoryRepo()
	conversationRepo := rdbms.NewConversationRepo(duckdbClient.GetDB())
	conversationService := conversation.NewConversationService(conversationRepo)
	summarizerService := summarizer.NewNoOpService()
	memoryService := memory.NewCachedService(memoryRepo, summarizerService)
	return &shortTermMemoryClient{
		config:              config,
		memoryService:       memoryService,
		conversationService: conversationService,
	}, nil
}

func (r *shortTermMemoryClient) Store(ctx context.Context, input types.StoreShortTermMemoryInput) (types.StoreShortTermMemoryOutput, error) {
	if input.Query == "" || input.Response == "" {
		log.Printf("[ERROR] Store: Query and response are required but one or both were empty (query: %q, response: %q)", input.Query, input.Response)
		return types.StoreShortTermMemoryOutput{}, fmt.Errorf("query and response are required")
	}
	if input.ConversationID == "" {
		log.Printf("[ERROR] Store: Conversation ID is required but was empty")
		return types.StoreShortTermMemoryOutput{}, fmt.Errorf("conversation id is required")
	}
	conversationID, err := uuid.Parse(input.ConversationID)
	if err != nil {
		log.Printf("[ERROR] Store: Invalid conversation ID format - %q, error: %v", input.ConversationID, err)
		return types.StoreShortTermMemoryOutput{}, fmt.Errorf("invalid conversation id")
	}
	exists, err := r.conversationService.Exists(ctx, conversationID)
	if err != nil {
		log.Printf("[ERROR] Store: Failed to check if conversation exists (conversationID: %s) - %v", conversationID, err)
		return types.StoreShortTermMemoryOutput{}, fmt.Errorf("error checking if conversation exists")
	}
	if !exists {
		log.Printf("[ERROR] Store: Conversation does not exist (conversationID: %s)", conversationID)
		return types.StoreShortTermMemoryOutput{}, fmt.Errorf("conversation does not exist")
	}
	id, err := r.memoryService.Store(ctx, conversationID, input.Query, input.Response)
	if err != nil {
		log.Printf("[ERROR] Store: Failed to store memory (conversationID: %s) - %v", conversationID, err)
		return types.StoreShortTermMemoryOutput{}, fmt.Errorf("error storing memory")
	}
	return types.StoreShortTermMemoryOutput{
		MemoryID: id.String(),
	}, nil
}

func (r *shortTermMemoryClient) Retrieve(ctx context.Context, input types.RetrieveShortTermMemoryInput) (types.RetrieveShortTermMemoryOutput, error) {
	if input.ConversationID == "" {
		log.Printf("[ERROR] Retrieve: Conversation ID is required but was empty")
		return types.RetrieveShortTermMemoryOutput{}, fmt.Errorf("conversation id is required")
	}
	conversationID, err := uuid.Parse(input.ConversationID)
	if err != nil {
		log.Printf("[ERROR] Retrieve: Invalid conversation ID format - %q, error: %v", input.ConversationID, err)
		return types.RetrieveShortTermMemoryOutput{}, fmt.Errorf("invalid conversation id")
	}
	exists, err := r.conversationService.Exists(ctx, conversationID)
	if err != nil {
		log.Printf("[ERROR] Retrieve: Failed to check if conversation exists (conversationID: %s) - %v", conversationID, err)
		return types.RetrieveShortTermMemoryOutput{}, fmt.Errorf("error checking if conversation exists")
	}
	if !exists {
		log.Printf("[ERROR] Retrieve: Conversation does not exist (conversationID: %s)", conversationID)
		return types.RetrieveShortTermMemoryOutput{}, fmt.Errorf("conversation does not exist")
	}
	var topK int
	if input.TopK == 0 {
		topK = 10
	}
	retrievedMemories, err := r.memoryService.Retrieve(ctx, conversationID, topK)
	if err != nil {
		log.Printf("[ERROR] Retrieve: Failed to retrieve memories (conversationID: %s, topK: %d) - %v", conversationID, topK, err)
		return types.RetrieveShortTermMemoryOutput{}, fmt.Errorf("error retrieving memories")
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
	return types.RetrieveShortTermMemoryOutput{
		Memories: memories,
	}, nil
}

func (r *shortTermMemoryClient) RegisterConversation(ctx context.Context, input types.RegisterConversationInput) (types.RegisterConversationOutput, error) {
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
