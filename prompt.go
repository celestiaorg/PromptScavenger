package prompt

import (
	"fmt"
	"context"
	"log"
	"os"
	"github.com/celestiaorg/celestia-node/api/rpc/client"
)

func getCurrentBlockHeight() uint64 {
	// TODO
}

func submitPrompt(client *client.Client, namespaceID namespace.ID, payload []byte, fee cosmosmath.Int, gasLimit uint64) uint64 {
        response, err := client.State.SubmitPayForBlob(context.Background()
, namespaceID, payLoad, fee, gasLimit)
        if err != nil {
                log.Fatalf("Error submitting pay for blob: %v", err)
        }
        fmt.Printf("Got output: %v", response)
        height := uint64(response.Height)
        fmt.Printf("Height that data was submitted at: %v", height)
        return height
}

func main(){
	// TODO
}

