package client

import (
	llmv1 "github.com/GitEval/GitEval-Backend/client/gen"
	"github.com/GitEval/GitEval-Backend/conf"
	"google.golang.org/grpc"
	"log"
)

func NewLLMClient(config *conf.LLMConfig) llmv1.LLMServiceClient {
	// 建立 gRPC 连接
	conn, err := grpc.Dial(config.Addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}

	// 创建并返回 gRPC 客户端
	client := llmv1.NewLLMServiceClient(conn)
	return client
}
