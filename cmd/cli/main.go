package main

import (
	"context"
	"fmt"

	"github.com/haren7/minimal-memory/clients"
	"github.com/haren7/minimal-memory/types"
)

func main() {
	ctx := context.Background()
	semanticMemoryClient, err := clients.NewSemanticMemoryClient(clients.SemanticMemoryClientConfig{
		ContextWindowSize: 10,
		OpenAIApiKey:      "sk-proj-",
	})
	if err != nil {
		fmt.Println("error creating semantic memory client: ", err)
		return
	}
	registerConversationOutput, err := semanticMemoryClient.RegisterConversation(ctx, types.RegisterConversationInput{
		Agent: "agent",
		User:  "user",
	})
	if err != nil {
		fmt.Println("error registering conversation: ", err)
		return
	}
	conversationID := registerConversationOutput.ConversationID
	storeSemanticMemoryOutput, err := semanticMemoryClient.Store(ctx, types.StoreSemanticMemoryInput{
		ConversationID: conversationID,
		Query:          "how can i use the semantic memory client?",
		Response:       "this is how you can use the semantic memory client",
	})
	if err != nil {
		fmt.Println("error storing semantic memory: ", err)
		return
	}
	memoryID := storeSemanticMemoryOutput.MemoryID
	fmt.Println("memory id: ", memoryID)
	retrieveSemanticMemoryOutput, err := semanticMemoryClient.Retrieve(ctx, types.RetrieveSemanticMemoryInput{
		ConversationID: conversationID,
		Query:          "help me use the semantic memory client",
	})
	if err != nil {
		fmt.Println("error retrieving semantic memory: ", err)
		return
	}
	memories := retrieveSemanticMemoryOutput.Memories
	similarMemories := retrieveSemanticMemoryOutput.SimilarMemories
	for _, memory := range memories {
		fmt.Println("memories:")
		fmt.Println("memory id: ", memory.ID)
		fmt.Println("memory query: ", memory.Query)
		fmt.Println("memory response: ", memory.Response)
		fmt.Println("memory created at: ", memory.CreatedAt)
		fmt.Println("--------------------------------")
	}
	for _, memory := range similarMemories {
		fmt.Println("similar memories:")
		fmt.Println("memory id: ", memory.ID)
		fmt.Println("memory query: ", memory.Query)
		fmt.Println("memory response: ", memory.Response)
		fmt.Println("memory created at: ", memory.CreatedAt)
		fmt.Println("--------------------------------")
	}
}
