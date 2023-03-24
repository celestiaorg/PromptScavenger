package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
	nodeheader "github.com/celestiaorg/celestia-node/header"
	"github.com/celestiaorg/nmt/namespace"
	"github.com/joho/godotenv"
	cosmosmath "cosmossdk.io/math"
	openai "github.com/sashabaranov/go-openai"
	"encoding/base64"
	"encoding/hex"
)

// gpt3 processes a given message using GPT-3 and prints the response.
func gpt3(msg string) {
	// Set the authentication header
	openAIKey := os.Getenv("OPENAI_KEY")
	client := openai.NewClient(openAIKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}

// getData fetches data from the Celestia node based on the provided height and namespace ID.
func getData(client *client.Client, height uint64, namespaceID namespace.ID) {
	endHeight := height + 100
	fromParam := getHeader(client, height)
	responseRange, err := client.Header.GetVerifiedRangeByHeight(context.Background(), fromParam, endHeight)
	if err != nil {
		log.Fatalf("Error getting verified range by height: %v", err)
	}
	fmt.Printf("Got header: %v", responseRange)
}

func getDataAsPrompt(client *client.Client, height uint64, namespaceID namespace.ID) string {
	headerParam := getHeader(client, height)
	fmt.Println(headerParam.DAH)
	response, err := client.Share.GetSharesByNamespace(context.Background(), headerParam.DAH, namespaceID)
	if err != nil {
		log.Fatalf("Error getting shares by namespace data for block height: %v. Error is %v", height, err)
	}
	fmt.Println(response)
	fmt.Println(len(response))
	var dataString string
	for _, shares := range response {
		fmt.Println(shares)
		for _, share := range shares.Shares {
			fmt.Println(string(share[8:]))
			dataString = string(share[8:])
		}
	}
	return dataString
}

// postData submits a new transaction with the provided data to the Celestia node.
func postDataAndGetHeight(client *client.Client, namespaceID namespace.ID, payLoad []byte, fee cosmosmath.Int, gasLimit uint64) uint64 {
	response, err := client.State.SubmitPayForBlob(context.Background(), namespaceID, payLoad, fee, gasLimit)
	if err != nil {
		log.Fatalf("Error submitting pay for blob: %v", err)
	}
	fmt.Printf("Got output: %v", response)
	height := uint64(response.Height)
	fmt.Printf("Height that data was submitted at: %v", height)
	return height
}

// getHeader fetches a header from the Celestia node based on the provided height.
func getHeader(client *client.Client, height uint64) *nodeheader.ExtendedHeader {
	header, err := client.Header.GetByHeight(context.Background(), height)
	if err != nil {
		log.Fatalf("Error getting header by height: %v", err)
	}
	return header
}

// loadEnv loads environment variables from the .env file.
func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// createClient initializes a new Celestia node client.
func createClient(ctx context.Context) *client.Client {
	nodeRPCIP := os.Getenv("NODE_RPC_IP")
	jwtToken := os.Getenv("NODE_JWT_TOKEN")

	rpc, err := client.NewClient(ctx, nodeRPCIP, jwtToken)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	return rpc
}

func createNamespaceID() []byte {
	nIDString := os.Getenv("NAMESPACE_ID")
	fmt.Println(nIDString)
	data, err := hex.DecodeString(nIDString)
	if err != nil {
		log.Fatalf("Error decoding hex string:", err)
	}
	fmt.Println(data)
	// Encode the byte array in Base64
	base64Str := base64.StdEncoding.EncodeToString(data)
	fmt.Println(base64Str)
	namespaceID, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		log.Fatalf("Error decoding Base64 string:", err)
	}
	fmt.Println(namespaceID)
	return namespaceID
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	loadEnv()
	var namespaceID namespace.ID = createNamespaceID() 
	client := createClient(ctx)
	var gasLimit uint64 = 6000000
	fee := cosmosmath.NewInt(10000)
	getData(client, 20, namespaceID)
	var gptPrompt string = "Tell me about modular blockchains"
	prompt := []byte{0x00, 0x01, 0x02}
	prompt = append(prompt, []byte(gptPrompt)...)
	fmt.Println(prompt)
	height := postDataAndGetHeight(client, namespaceID, prompt, fee, gasLimit)
	promptString := getDataAsPrompt(client, height, namespaceID)
	gpt3(promptString)
	// Close the client when you are finished
	client.Close()
}
