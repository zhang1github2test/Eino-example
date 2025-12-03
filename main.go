package main

import (
	workflow_demo "Eino-example/workflow-demo"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"io"
	"log"
)

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	workflow_demo.FlowDemo()

}

func reportStream(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()
	i := 0
	for {
		message, err := sr.Recv()
		if err == io.EOF { // 流式输出结束
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		log.Printf("message[%d]: %+v\n", i, message)
		i++
	}
}
