package main

import (
	"bytes"
	"concept/pkg/env"
	"concept/pkg/filetree"
	"concept/pkg/functions"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Parse command line flags
	rootDir := flag.String("root", ".", "Path to the root directory of the project")
	logLevel := flag.String("log-level", "info", "Set the logging level (debug, info, warn, error, fatal, panic)")
	envFile := flag.String("env-file", ".env", "Path to .env file")
	flag.Parse()

	log.Debug().
		Str("root", *rootDir).
		Str("log-level", *logLevel).
		Str("env-file", *envFile).
		Msg("Starting application with flags")

	// Load environment variables
	if err := env.Load(*envFile); err != nil {
		log.Warn().Err(err).Str("file", *envFile).Msg("Failed to load environment variables")
	} else {
		log.Debug().Str("file", *envFile).Msg("Successfully loaded environment variables")
	}

	// Configure logging
	previousLevel := zerolog.GlobalLevel()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	switch strings.ToLower(*logLevel) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Debug().
		Str("previous", previousLevel.String()).
		Str("current", zerolog.GlobalLevel().String()).
		Msg("Log level changed")

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal().Msg("OPENAI_API_KEY environment variable is required")
	}
	log.Debug().Msg("API key found in environment")

	// Create OpenAI handler
	client := openai.NewClient(
		option.WithAPIKey("My API Key"), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
	log.Debug().Str("model", "gpt-4").Msg("Created OpenAI handler")

	// Create file trees
	log.Debug().Str("root", *rootDir).Msg("Creating project file tree")
	projectTree, err := filetree.NewFileTree(*rootDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create project file tree")
	}
	log.Info().Msg("Project file tree created successfully")

	log.Debug().Str("root", *rootDir).Msg("Creating rules file tree")
	rulesTree, err := filetree.NewRulesFileTree(*rootDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create rules file tree")
	}
	log.Info().Msg("Rules file tree created successfully")

	ctx := context.Background()

	uploadService := openai.NewUploadService()

	// Loop over every file in the projectTree and upload them to OpenAI
	for _, file := range projectTree.Files {
		// read the file
		content, err := os.ReadFile(file.Path)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read project file")
			continue
		}

		upload, err := uploadService.New(ctx, openai.UploadNewParams{
			Bytes:    openai.F(int64(len(content))),
			Filename: openai.F(file.Path),
			MimeType: openai.F("text/plain"),
			Purpose:  openai.F(openai.FilePurposeAssistants),
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to upload project file")
		}
		log.Debug().Str("file", file.Path).Msg("Uploaded project file")

		// add parts
		part := openai.UploadPartNewParams{
			Data: openai.F(io.Reader(bytes.NewReader(content))),
		}

		uploadPart, err := uploadService.Parts.New(ctx, upload.ID, part)
		if err != nil {
			log.Error().Err(err).Msg("Failed to add part to upload")
			continue
		}
		log.Debug().Str("part_id", uploadPart.ID).Msg("Added part to upload")
	}

	// Loop over every file in the rulesTree and upload them to OpenAI
	// for _, file := range rulesTree.Tree {
	// 	// Upload the file for search purpose
	// 	uploadReq := openai.FileUploadRequest{
	// 		File:    file,
	// 		Purpose: "search",
	// 	}
	// 	uploadResp, err := client.Files.Upload(ctx, uploadReq)
	// 	if err != nil {
	// 		fmt.Println("Error uploading rules file:", err)
	// 		return
	// 	}
	// 	log.Debug().Str("file", file.Name).Msg("Uploaded rules file for search")
	// }

	log.Debug().
		Int("project_tree_length", len(projectTree.Tree)).
		Int("rules_tree_length", len(rulesTree.Tree)).
		Msg("Generated tree strings")

	// Add initial message with file trees and request confirmation
	initialMessage := fmt.Sprintf(
		// "Here are the file trees for analysis:\n\nProject Files:\n%s\n\nRules Files:\n%s\n\n"+
		"Your task is to analyse all files in the project tree according to the rules. " +
			// "You have access to two tools that you should use directly:\n"+
			// "- 'load_file' tool: Use this to read file contents. You must provide the complete file path from the root, reconstructed from the file tree structure\n"+
			// "- 'download_file' tool: Use this to download external sources\n\n"+
			"Follow these steps:\n" +
			// "1. First, use the 'load_file' tool to load ALL rules files. Rules are ALWAYS located in the '.cursor/rules/' directory. For each rules file:\n"+
			// "   - Look at the rules file tree structure and use the complete path (e.g., '.cursor/rules/repository-structure.mdc')\n"+
			// "   - Use the 'load_file' tool with the complete path including '.cursor/rules/'\n"+
			// "   - You can load the rules files in parallel function calls\n"+
			// "2. Then, use the 'download_file' tool to download ALL necessary sources\n"+
			// "   - You can download the sources in parallel function calls\n"+
			// "   - The sources are always explicitly mentioned in the rules files\n"+
			// "   - Don't just download all URLs, only sources specifically mentioned\n"+
			"3. For each project file that needs processing:\n" +
			"   a. Look at the file tree structure and reconstruct its complete path (e.g., 'cmd/main/main.go')\n" +
			// "   b. Use the 'load_file' tool with the complete path to read its contents\n" +
			"   c. Process the file according to the rules\n" +
			"   d. Generate a git diff patch if changes are needed\n" +
			"   e. Log the results\n\n" +
			"IMPORTANT:\n" +
			"- Keep your responses tiny, they're only used for logging and are not shown to the user\n" +
			"- Do not describe what tools you want to use - instead, use the tool_calls feature to actually execute them\n" +
			"- Always reconstruct complete file paths by following the nesting structure shown in the file tree\n" +
			"- Rules files are ALWAYS in '.cursor/rules/' directory - never use 'rules/' or just the filename\n" +
			"When you have processed all significant and relevant files, simply respond with 'DGLuke__done: <reason>'.\n" +
			"If you encounter an error please respond with 'DGLuke__error: <reason>'.",
		// projectTreeStr, rulesTreeStr
	)

	log.Debug().Int("message_length", len(initialMessage)).Msg("Created initial system message")

	// Example messages and tools
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(initialMessage),
		openai.UserMessage("Let's begin the analysis of the project files according to the rules."),
	}
	log.Debug().Int("initial_messages", len(messages)).Msg("Created initial message array")

	tools := []openai.ChatCompletionToolParam{
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.F("load_file"),
				Description: openai.F("Load and read the contents of a file"),
				Parameters: openai.F(openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"filepath": map[string]interface{}{
							"type":        "string",
							"description": "The complete path to the file from the root directory",
						},
					},
					"required": []string{"filepath"},
				}),
			}),
		},
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.F("download_file"),
				Description: openai.F("Download a file from an external source"),
				Parameters: openai.F(openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"url": map[string]interface{}{
							"type":        "string",
							"description": "The URL of the file to download",
						},
					},
					"required": []string{"url"},
				}),
			}),
		},
	}
	log.Debug().Int("tools_count", len(tools)).Msg("Created tools array")

	log.Info().Msg("Starting analysis loop")
	for {
		log.Debug().Int("message_count", len(messages)).Msg("Sending completion request")

		// Send completion request
		completion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
			Messages: openai.F(messages),
			Model:    openai.F(openai.ChatModelGPT4o),
			Tools:    openai.F(tools),
		})

		if err != nil {
			log.Fatal().
				Err(err).
				Int("message_count", len(messages)).
				Msg("Failed to get completion")
		}

		if len(completion.Choices) == 0 {
			log.Fatal().
				Int("message_count", len(messages)).
				Msg("No completion choices returned")
		}

		choice := completion.Choices[0]
		messages = append(messages, choice.Message)
		log.Debug().Str("content", choice.Message.Content).Msg("Received assistant message")

		// Check for special responses
		response := strings.TrimSpace(choice.Message.Content)
		if strings.HasPrefix(response, "DGLuke__done:") {
			reason := strings.TrimPrefix(response, "DGLuke__done:")
			log.Info().Str("reason", reason).Msg("Analysis completed")
			break
		}
		if strings.HasPrefix(response, "DGLuke__error:") {
			reason := strings.TrimPrefix(response, "DGLuke__error:")
			log.Error().
				Str("reason", reason).
				Int("message_count", len(messages)).
				Msg("Analysis error")
			break
		}

		// Process tool calls if any
		if len(choice.Message.ToolCalls) > 0 {
			log.Debug().Int("tool_calls", len(choice.Message.ToolCalls)).Msg("Processing tool calls")

			for _, toolCall := range choice.Message.ToolCalls {
				log.Debug().
					Str("tool", toolCall.Function.Name).
					Str("args", fmt.Sprintf("%.100s", toolCall.Function.Arguments)).
					Msg("Processing tool call")

				functionProcessor := functions.NewFunctionProcessor()
				result, err := functionProcessor.ProcessFunctionCall(toolCall)
				if err != nil {
					log.Error().
						Err(err).
						Str("tool", toolCall.Function.Name).
						Str("args", toolCall.Function.Arguments).
						Int("message_count", len(messages)).
						Msg("Tool call failed")
					messages = append(messages, openai.ToolMessage(toolCall.ID, fmt.Sprintf("Error: %v", err)))
					continue
				}

				log.Debug().
					Str("tool", toolCall.Function.Name).
					Int("result_length", len(result)).
					Msg("Tool call successful")
				messages = append(messages, openai.ToolMessage(toolCall.ID, result))
			}
		}

		// Add a message asking what to do next
		messages = append(messages, openai.UserMessage("What would you like to do next?"))
		log.Debug().Msg("Added next step prompt")
	}

	log.Info().Msg("Analysis process completed")
}
