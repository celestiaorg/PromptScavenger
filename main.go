package main

import (
  "context"
  "os"
  "log"
  nodeheader "github.com/celestiaorg/celestia-node/header"
  // state "github.com/celestiaorg/celestia-node/state"
  //"github.com/celestiaorg/nmt"
  namespace "github.com/celestiaorg/nmt/namespace"
  cosmosmath "cosmossdk.io/math"
  // please reference the go.mod file in this repository in order to correctly import
  // this package
  nodeclient "github.com/celestiaorg/celestia-node/api/rpc/client"
  "fmt"
  "github.com/joho/godotenv"
  openai "github.com/sashabaranov/go-openai"
)

func GPT3(msg string) {
  // set the authentication header
  OpenAIKey := os.Getenv("OPENAI_KEY")
  client := openai.NewClient(OpenAIKey)
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
 // send the request and process the response

 //fmt.Println(string(body))
}


func getData(client *nodeclient.Client, height uint64, namespaceID namespace.ID) {
	endHeight := height + 100
	fromParam := getHeader(client, height)
	responseRange, err := client.Header.GetVerifiedRangeByHeight(context.Background(), fromParam, endHeight)
	if err != nil {
	  panic(err)
	}	
	fmt.Printf("Got header: %v", responseRange)
}

func postData(client *nodeclient.Client, namespaceID namespace.ID, payLoad []byte, fee cosmosmath.Int, gasLimit uint64) {
	response, err := client.State.SubmitPayForBlob(context.Background(), namespaceID, payLoad, fee, gasLimit)
	if err != nil {
	  panic(err)
        }
	fmt.Printf("Got output: %v", response)
}

func getHeader(client *nodeclient.Client, height uint64) *nodeheader.ExtendedHeader {

   // call the GetByHeight method on the `HeaderModule` that returns a header to you
  // by the given height (20)
  header, err := client.Header.GetByHeight(context.Background(), height)
  if err != nil {
    panic(err)
  }
  return header
}

func loadEnv() {
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatal("Error loading .env file")
  }

}

func createClient(ctx context.Context) *nodeclient.Client {
  // Create a new client by dialing the celestia-node's RPC endpoint
  // By default, celestia-nodes run RPC on port 26658
  NodeRPCIP := os.Getenv("NODE_RPC_IP")
  JWTToken := os.Getenv("NODE_JWT_TOKEN")
  
  rpc, err := nodeclient.NewClient(ctx, NodeRPCIP, JWTToken)
  if err != nil {
    panic(err)
  }

  return rpc
}

func main() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()
  loadEnv()
  var namespaceID namespace.ID = []byte{5, 3, 0, 1, 4, 5, 9, 2}
  client := createClient(ctx)
  var gasLimit uint64 = 6000
  fee := cosmosmath.NewInt(2000)
  //header := getHeader(client, 20) 
  //fmt.Printf("Got header: %v", header)
  getData(client, 20, namespaceID)
  //var prompt []byte = "Tell me about modular blockchains"
  var gptprompt string = "Tell me about modular blockchains"
  prompt := []byte{0x00, 0x01, 0x02}
  prompt = append(prompt, []byte(gptprompt)...)
  fmt.Println(prompt)
  postData(client, namespaceID, prompt, fee, gasLimit)
  GPT3(gptprompt) 
  // close the client when you are finished
  client.Close()
}
