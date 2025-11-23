package main

import (
	"context"
	"fmt"
	"memory/clients"
	"memory/types"
)

func exampleShortTermMemory() {
	ctx := context.Background()
	shortTermMemoryClient, err := clients.NewShortTermMemoryClient(clients.ShortTermMemoryClientConfig{})
	if err != nil {
		fmt.Println("error creating short term memory client: ", err)
		return
	}
	registerConversationOutput, err := shortTermMemoryClient.RegisterConversation(ctx, types.RegisterConversationInput{
		Agent: "agent",
		User:  "user",
	})
	if err != nil {
		fmt.Println("error registering conversation: ", err)
		return
	}
	conversationID := registerConversationOutput.ConversationID
	storeShortTermMemoryOutput, err := shortTermMemoryClient.Store(ctx, types.StoreShortTermMemoryInput{
		Query:          "how can i use the short memory client?",
		Response:       "this is how you can use the short memory client",
		ConversationID: conversationID,
	})
	if err != nil {
		fmt.Println("error storing short term memory: ", err)
		return
	}
	memoryID := storeShortTermMemoryOutput.MemoryID
	fmt.Println("memory id: ", memoryID)
	retrieveShortTermMemoryOutput, err := shortTermMemoryClient.Retrieve(ctx, types.RetrieveShortTermMemoryInput{
		TopK:           10,
		ConversationID: conversationID,
	})
	if err != nil {
		fmt.Println("error retrieving short term memory: ", err)
		return
	}
	retrievedMemories := retrieveShortTermMemoryOutput.Memories
	for _, memory := range retrievedMemories {
		fmt.Println("memory id: ", memory.ID)
		fmt.Println("memory query: ", memory.Query)
		fmt.Println("memory response: ", memory.Response)
		fmt.Println("memory created at: ", memory.CreatedAt)
		fmt.Println("--------------------------------")
	}
}
