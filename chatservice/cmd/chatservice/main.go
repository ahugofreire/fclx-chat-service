package main

import (
	"database/sql"
	"fmt"
	"github.com/ahugofreire/chatservice/configs"
	"github.com/ahugofreire/chatservice/internal/infra/grpc/server"
	"github.com/ahugofreire/chatservice/internal/infra/repository"
	"github.com/ahugofreire/chatservice/internal/infra/web"
	"github.com/ahugofreire/chatservice/internal/infra/web/webserver"
	"github.com/ahugofreire/chatservice/internal/usecase/chatcompletion"
	"github.com/ahugofreire/chatservice/internal/usecase/chatcompletionstream"
	"github.com/sashabaranov/go-openai"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/sashabaranov/go-openai"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	conn, err := sql.Open(config.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName))

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	repo := repository.NewChatRepositoryMySQL(conn)
	client := openai.NewClient(config.OpenAIApiKey)

	chatConfig := chatcompletion.ChatCompletionConfigInputDTO{
		Model:                config.Model,
		ModelMaxTokens:       config.ModelMaxTokens,
		Temperature:          float32(config.Temperature),
		TopP:                 float32(config.TopP),
		N:                    config.N,
		Stop:                 config.Stop,
		MaxTokens:            config.MaxTokens,
		InitialSystemMessage: config.InitialChatMessage,
	}

	chatConfigStream := chatcompletionstream.ChatCompletionConfigDTO{
		Model:                config.Model,
		ModelMaxTokens:       config.ModelMaxTokens,
		Temperature:          float32(config.Temperature),
		TopP:                 float32(config.TopP),
		N:                    config.N,
		Stop:                 config.Stop,
		MaxTokens:            config.MaxTokens,
		InitialSystemMessage: config.InitialChatMessage,
	}

	useCase := chatcompletion.NewChatCompletionUseCase(repo, client)

	streamChannel := make(chan chatcompletionstream.ChatCompletionOutputDTO)
	useCaseStream := chatcompletionstream.NewChatCompletionUseCase(repo, client, streamChannel)

	fmt.Println("Server running on port " + config.GRPCServerPort)
	grpcServer := server.NewGRPCServer(*useCaseStream, chatConfigStream, config.GRPCServerPort, config.AuthToken, streamChannel)
	go grpcServer.Start()

	webServer := webserver.NewWebServer(":" + config.WebServerPort)
	webServerChatHandler := web.NewWebChatGPTHandler(*useCase, chatConfig, config.AuthToken)
	webServer.AddHandler("/chat", webServerChatHandler.Handle)

	fmt.Println("Server running on port " + config.WebServerPort)
	webServer.Start()

}
