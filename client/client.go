package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"grpc-redis-client/protos/todo/protos/todo"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	//serverAddress := os.Getenv("SERVER_ADDRESS")
	//address := fmt.Sprintf("%s:8888", serverAddress)

	conn, err := grpc.NewClient("grpc-redis-server:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := todo.NewTodoServiceClient(conn)

	fmt.Println("Connecting Client on Port 8888...")

	// Connect to Redis
	redisAddr := "redis:6379"
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		args := parseInput(input)
		if len(args) == 0 {
			continue
		}

		command := args[0]

		switch command {
		case "add":
			if len(args) != 3 {
				fmt.Println("Usage: add \"<title>\" \"<description>\"")
				continue
			}
			title := args[1]
			description := args[2]
			addTodo(client, title, description)
		case "get":
			if len(args) != 2 {
				fmt.Println("Usage: get <id>")
				continue
			}
			id := args[1]
			getTodo(client, id)
		case "update":
			if len(args) != 4 {
				fmt.Println("Usage: update <id> \"<title>\" \"<description>\"")
				continue
			}
			id := args[1]
			title := args[2]
			description := args[3]
			updateTodo(client, id, title, description)
		case "delete":
			if len(args) != 2 {
				fmt.Println("Usage: delete <id>")
				continue
			}
			id := args[1]
			deleteTodo(client, id)
		case "list":
			listTodos(client)
		default:
			fmt.Println("Unknown command:", command)
			fmt.Println("Available commands: add, get, update, delete, list")
		}
	}
}

func parseInput(input string) []string {
	re := regexp.MustCompile(`"[^"]*"|\S+`)
	matches := re.FindAllString(input, -1)
	for i, match := range matches {
		// Trim the double quotes from the arguments
		matches[i] = strings.Trim(match, `"`)
	}
	return matches
}

func addTodo(client todo.TodoServiceClient, title, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.AddTodo(ctx, &todo.AddTodoRequest{Title: title, Description: description})
	if err != nil {
		log.Fatalf("could not add todo: %v", err)
	}
	fmt.Printf("Todo added with ID: %s\n", resp.Id)
}

func getTodo(client todo.TodoServiceClient, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetTodo(ctx, &todo.GetTodoRequest{Id: id})
	if err != nil {
		log.Fatalf("could not get todo: %v", err)
	}
	fmt.Printf("Todo: ID: %s, Title: %s, Description: %s\n", resp.Id, resp.Title, resp.Description)
}

func updateTodo(client todo.TodoServiceClient, id, title, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.UpdateTodo(ctx, &todo.UpdateTodoRequest{Id: id, Title: title, Description: description})
	if err != nil {
		log.Fatalf("could not update todo: %v", err)
	}
	fmt.Printf("Todo updated: %s\n", resp.Success)
}

func deleteTodo(client todo.TodoServiceClient, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.DeleteTodo(ctx, &todo.DeleteTodoRequest{Id: id})
	if err != nil {
		log.Fatalf("could not delete todo: %v", err)
	}
	fmt.Printf("Todo deleted: %s\n", resp.Success)
}

func listTodos(client todo.TodoServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.ListTodos(ctx, &todo.ListTodosRequest{})
	if err != nil {
		log.Fatalf("could not list todos: %v", err)
	}
	fmt.Println("Todos:")
	for _, todoItem := range resp.Todos {
		fmt.Printf("ID: %s, Title: %s, Description: %s\n", todoItem.Id, todoItem.Title, todoItem.Description)
	}
}
