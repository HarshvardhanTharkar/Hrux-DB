package main

import (
	"fmt"
	"log"

	"kv-distributed/client"
)

func main() {
	// Create client connection
	kvClient, err := client.NewKVClient("localhost:8080")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer kvClient.Close()

	// Example usage - Basic KV operations
	fmt.Println("Testing KV operations...")

	// Put operation
	err = kvClient.Put("users", "user1", []byte("Harsh"), 0)
	if err != nil {
		log.Fatal("Put failed:", err)
	}
	fmt.Println("✓ Put user1 successful")

	// Get operation
	value, err := kvClient.Get("users", "user1")
	if err != nil {
		log.Fatal("Get failed:", err)
	}
	fmt.Printf("✓ Get user1: %s\n", string(value))

	// Put with TTL
	err = kvClient.Put("sessions", "session1", []byte("active"), 60)
	if err != nil {
		log.Fatal("Put with TTL failed:", err)
	}
	fmt.Println("✓ Put session1 with TTL successful")

	// List keys by prefix
	err = kvClient.Put("users", "user2", []byte("Jane Smith"), 0)
	if err != nil {
		log.Fatal("Put failed:", err)
	}

	keys, err := kvClient.ListKeysByPrefix("users", "user")
	if err != nil {
		log.Fatal("ListKeysByPrefix failed:", err)
	}
	fmt.Printf("✓ Keys with prefix 'user': %v\n", keys)

	// Set operations
	err = kvClient.SetAdd("admins", "user1")
	if err != nil {
		log.Fatal("SetAdd failed:", err)
	}
	fmt.Println("✓ SetAdd successful")

	setList, err := kvClient.SetList("admins")
	if err != nil {
		log.Fatal("SetList failed:", err)
	}
	fmt.Printf("✓ Set members: %v\n", setList)

	// List all in bucket
	data, err := kvClient.List("users")
	if err != nil {
		log.Fatal("List failed:", err)
	}
	fmt.Printf("✓ All users: %+v\n", data)

	// Queue operations
	err = kvClient.QueuePush("tasks", "task1")
	if err != nil {
		log.Fatal("QueuePush failed:", err)
	}
	fmt.Println("✓ QueuePush successful")

	// Stack operations
	err = kvClient.StackPush("history", "action1")
	if err != nil {
		log.Fatal("StackPush failed:", err)
	}
	fmt.Println("✓ StackPush successful")

	fmt.Println("All operations completed successfully!")
}
