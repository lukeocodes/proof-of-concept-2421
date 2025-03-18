package main

import (
	"bytes"
	"concept/pkg/env"
	"concept/pkg/filetree"
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// var standardSources = []string{
// 	"https://raw.githubusercontent.com/deepgram/deepgram-api-specs/refs/heads/main/openapi.yml",
// 	"https://raw.githubusercontent.com/deepgram/deepgram-api-specs/refs/heads/main/asyncapi.yml",
// 	"https://raw.githubusercontent.com/deepgram/deepgram-js-sdk/refs/heads/main/README.md",
// 	"https://raw.githubusercontent.com/deepgram/deepgram-python-sdk/refs/heads/main/README.md",
// 	"https://raw.githubusercontent.com/deepgram/deepgram-go-sdk/refs/heads/main/README.md",
// 	"https://raw.githubusercontent.com/deepgram/deepgram-dotnet-sdk/refs/heads/main/README.md",
// 	"https://raw.githubusercontent.com/deepgram/deepgram-rust-sdk/refs/heads/main/README.md",
// }

func getMimeType(path string) string {
	ext := filepath.Ext(path)

	switch ext {
	case ".cpp":
		return "text/x-c++"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".doc":
		return "application/msword"
	case ".py":
		return "text/x-script.python"
	case ".html":
		return "text/html"
	case ".c":
		return "text/x-c"
	case ".cs":
		return "text/x-csharp"
	case ".txt":
		return "text/plain"
	case ".xml":
		return "application/xml"
	case ".jpeg":
		return "image/jpeg"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".php":
		return "text/x-php"
	case ".a":
		return "application/octet-stream"
	case ".bin":
		return "application/octet-stream"
	case ".bpk":
		return "application/octet-stream"
	case ".deploy":
		return "application/octet-stream"
	case ".dist":
		return "application/octet-stream"
	case ".distz":
		return "application/octet-stream"
	case ".dmg":
		return "application/octet-stream"
	case ".dms":
		return "application/octet-stream"
	case ".dump":
		return "application/octet-stream"
	case ".elc":
		return "application/octet-stream"
	case ".lha":
		return "application/octet-stream"
	case ".lrf":
		return "application/octet-stream"
	case ".lzh":
		return "application/octet-stream"
	case ".o":
		return "application/octet-stream"
	case ".obj":
		return "application/octet-stream"
	case ".pkg":
		return "application/octet-stream"
	case ".so":
		return "application/octet-stream"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".gif":
		return "image/gif"
	case ".tex":
		return "text/x-tex"
	case ".json":
		return "application/json"
	case ".sh":
		return "text/x-sh"
	case ".ts":
		return "application/typescript"
	case ".tsx":
		return "application/typescript"
	case ".js":
		return "text/javascript"
	case ".jsx":
		return "text/javascript"
	case ".md":
		return "text/markdown"
	case ".css":
		return "text/css"
	case ".csv":
		return "text/csv"
	case ".webp":
		return "image/webp"
	case ".tar":
		return "application/x-tar"
	case ".ruby":
		return "text/x-ruby"
	case ".png":
		return "text/x-java"
	default:
		return "text/plain"
	}
}

func main() {
	// Parse command line flags
	projectTitle := flag.String("project-title", "", "Title of the project")
	rootDir := flag.String("root", ".", "Path to the root directory of the project")
	logLevel := flag.String("log-level", "info", "Set the logging level (debug, info, warn, error, fatal, panic)")
	envFile := flag.String("env-file", ".env", "Path to .env file")
	flag.Parse()

	log.Debug().
		Str("root", *rootDir).
		Str("log-level", *logLevel).
		Str("project-title", *projectTitle).
		Str("env-file", *envFile).
		Msg("Starting application with flags")

	if *projectTitle == "" {
		log.Fatal().Msg("Project title is required")
	}

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
	client := openai.NewClient()
	log.Debug().Msg("Created OpenAI handler")

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

	// create a vector store
	store, err := client.Beta.VectorStores.New(ctx, openai.BetaVectorStoreNewParams{
		Name: openai.F("Repository Review " + *projectTitle + " " + time.Now().Format("2006-01-02")),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create vector store")
	}
	log.Info().Msg("Vector store created successfully")

	// Loop over every file in the rulesTree and upload them to OpenAI
	for _, file := range rulesTree.Files {
		// read the file
		content, err := os.ReadFile(".cursor/rules/" + file.Path)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read rule file")
			continue
		}

		// create a upload object
		upload, err := client.Uploads.New(ctx, openai.UploadNewParams{
			Purpose:  openai.F(openai.FilePurposeAssistants),
			Filename: openai.F(file.Path),
			Bytes:    openai.F(int64(len(content))),
			MimeType: openai.F("text/markdown"),
		})

		if err != nil {
			log.Error().Err(err).Msg("Failed to create upload object")
			continue
		}

		// loop over every 64mb chunk of the file
		chunkSize := 64 * 1024 * 1024
		numChunks := len(content) / chunkSize
		if len(content)%chunkSize != 0 {
			numChunks++
		}

		var uploadParts []openai.UploadPart

		for i := 0; i < numChunks; i++ {
			start := i * chunkSize
			end := start + chunkSize
			if end > len(content) {
				end = len(content)
			}
			chunk := content[start:end]

			// create a upload part
			uploadPart, err := client.Uploads.Parts.New(ctx, upload.ID, openai.UploadPartNewParams{
				Data: openai.F(io.Reader(bytes.NewReader(chunk))),
			})

			if err != nil {
				log.Error().Err(err).Msg("Failed to create upload part")
				continue
			}

			uploadParts = append(uploadParts, *uploadPart)
		}

		complete, err := client.Uploads.Complete(ctx, upload.ID, openai.UploadCompleteParams{
			PartIDs: openai.F(func() []string {
				partIds := make([]string, len(uploadParts))
				for i, part := range uploadParts {
					partIds[i] = part.ID
				}
				return partIds
			}()),
		})

		if err != nil {
			log.Error().Err(err).Msg("Failed to complete upload")
			continue
		}

		log.Info().Str("file_id", complete.File.ID).Int("file_size", int(complete.Bytes)).Str("file_name", complete.File.Filename).Msg("Rule uploaded successfully")
	}

	var uploads []openai.Upload
	var completedUploads []openai.Upload

	// Loop over every file in the projectTree and upload them to OpenAI
	for _, file := range projectTree.Files {
		// read the file
		content, err := os.ReadFile(file.Path)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read project file")
			continue
		}

		// create a upload object
		upload, err := client.Uploads.New(ctx, openai.UploadNewParams{
			Purpose:  openai.F(openai.FilePurposeAssistants),
			Filename: openai.F(file.Path + ".txt"),
			Bytes:    openai.F(int64(len(content))),
			MimeType: openai.F(getMimeType(file.Path)),
		})

		if err != nil {
			log.Error().Err(err).Msg("Failed to create upload object")
			continue
		}

		uploads = append(uploads, *upload)

		// loop over every 64mb chunk of the file
		chunkSize := 64 * 1024 * 1024
		numChunks := len(content) / chunkSize
		if len(content)%chunkSize != 0 {
			numChunks++
		}

		var uploadParts []openai.UploadPart

		for i := 0; i < numChunks; i++ {
			start := i * chunkSize
			end := start + chunkSize
			if end > len(content) {
				end = len(content)
			}
			chunk := content[start:end]

			// create a upload part
			uploadPart, err := client.Uploads.Parts.New(ctx, upload.ID, openai.UploadPartNewParams{
				Data: openai.F(io.Reader(bytes.NewReader(chunk))),
			})

			if err != nil {
				log.Error().Err(err).Msg("Failed to create upload part")
				continue
			}

			uploadParts = append(uploadParts, *uploadPart)
		}

		complete, err := client.Uploads.Complete(ctx, upload.ID, openai.UploadCompleteParams{
			PartIDs: openai.F(func() []string {
				partIds := make([]string, len(uploadParts))
				for i, part := range uploadParts {
					partIds[i] = part.ID
				}
				return partIds
			}()),
		})

		completedUploads = append(completedUploads, *complete)

		if err != nil {
			log.Error().Err(err).Msg("Failed to complete upload")
			continue
		}

		log.Info().Str("file_id", complete.File.ID).Int("file_size", int(complete.Bytes)).Str("file_name", complete.File.Filename).Msg("Upload completed successfully")
	}

	batch, err := client.Beta.VectorStores.FileBatches.New(ctx, store.ID, openai.BetaVectorStoreFileBatchNewParams{
		FileIDs: openai.F(func() []string {
			fileIds := make([]string, len(completedUploads))
			for i, upload := range completedUploads {
				fileIds[i] = upload.File.ID
			}
			return fileIds
		}()),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file batch")
	}

	log.Info().Str("file_batch_id", batch.ID).Int("file_count", len(uploads)).Msg("Vector store file batch created successfully")

	log.Debug().
		Int("project_tree_length", len(projectTree.Tree)).
		Int("rules_tree_length", len(rulesTree.Tree)).
		Msg("Generated tree strings")

	// Add initial message with file trees and request confirmation
	var developerMessage = []string{
		"Analyse all files in the project tree according to these rules:",
		"- Respond ONLY with the git diff patch for the file.",
		"- Rules files are always in `.cursor/rules/`.",
		"- Review each file one at a time.",
		"- Do not include any other text than the git diff patch.",
		"- Review each file against against any rules with a glob value that makes them relevant.",
		"",
		"No files have been processed yet.",
	}

	log.Debug().Int("message_length", len(strings.Join(developerMessage, "\n"))).Msg("Created developer message")

	// Send request to /v1/responses endpoint
	responseAPIRequestParameters := map[string]interface{}{
		"model": openai.ChatModelO3Mini,
		"reasoning": map[string]interface{}{
			"effort": openai.ChatCompletionReasoningEffortMedium,
		},
		"store": false,
		"input": []map[string]interface{}{
			{
				"role": openai.ChatCompletionMessageParamRoleDeveloper,
				"content": []map[string]interface{}{
					{
						"type": "input_text",
						"text": strings.Join(developerMessage, "\n"),
					},
				},
			},
			{
				"role": openai.ChatCompletionMessageParamRoleUser,
				"content": []map[string]interface{}{
					{
						"type": "input_text",
						"text": "Get ready to review the " + *projectTitle + " project! Let me know when you are ready.",
					},
				},
			},
		},
		"tools": []map[string]interface{}{
			{
				"type":             openai.AssistantToolTypeFileSearch,
				"vector_store_ids": []string{store.ID},
			},
		},
	}

	type responsesApiResult struct {
		ID                string `json:"id"`
		Object            string `json:"object"`
		CreatedAt         int    `json:"created_at"`
		Status            string `json:"status"`
		Error             string `json:"error"`
		IncompleteDetails string `json:"incomplete_details"`
		Instructions      string `json:"instructions"`
		MaxOutputTokens   int    `json:"max_output_tokens"`
		Model             string `json:"model"`
		Output            []struct {
			Content []struct {
				Type        string   `json:"type"`
				Text        string   `json:"text"`
				Annotations []string `json:"annotations"`
			} `json:"content"`
			ID     string `json:"id"`
			Role   string `json:"role"`
			Status string `json:"status"`
			Type   string `json:"type"`
		} `json:"output"`
		ParallelToolCalls  bool   `json:"parallel_tool_calls"`
		PreviousResponseID string `json:"previous_response_id"`
		Reasoning          struct {
			Effort  string `json:"effort"`
			Summary string `json:"summary"`
		} `json:"reasoning"`
		Store       bool    `json:"store"`
		Temperature float64 `json:"temperature"`
		Text        struct {
			Format struct {
				Type string `json:"type"`
			} `json:"format"`
		} `json:"text"`
		ToolChoice string `json:"tool_choice"`
		Tools      []struct {
			Filters        interface{} `json:"filters"`
			MaxNumResults  int         `json:"max_num_results"`
			RankingOptions struct {
				Ranker         string  `json:"ranker"`
				ScoreThreshold float64 `json:"score_threshold"`
			} `json:"ranking_options"`
			Type           string   `json:"type"`
			VectorStoreIDs []string `json:"vector_store_ids"`
		} `json:"tools"`
		TopP       float64 `json:"top_p"`
		Truncation string  `json:"truncation"`
		Usage      struct {
			InputTokens        int `json:"input_tokens"`
			InputTokensDetails struct {
				CachedTokens int `json:"cached_tokens"`
			} `json:"input_tokens_details"`
			OutputTokens        int `json:"output_tokens"`
			OutputTokensDetails struct {
				ReasoningTokens int `json:"reasoning_tokens"`
			} `json:"output_tokens_details"`
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
		User     interface{}            `json:"user"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	var responseAPIResult responsesApiResult

	err = client.Post(context.TODO(), "/responses", responseAPIRequestParameters, &responseAPIResult)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get response from /v1/responses")
	}

	log.Info().Str("response_id", responseAPIResult.ID).Msg("Response received")

	log.Info().
		Interface("responseResult", responseAPIResult).
		Msg("Full response result")

	var previousResponseID string = responseAPIResult.ID

	for _, upload := range uploads {
		log.Info().Str("file", upload.Filename).Msg("Starting analysis")

		// Send request to /v1/responses endpoint
		responseParams := map[string]interface{}{
			"model": openai.ChatModelO3Mini,
			"reasoning": map[string]interface{}{
				"effort": openai.ChatCompletionReasoningEffortMedium,
			},
			"previous_response_id": previousResponseID,
			"store":                false,
			"input": []map[string]interface{}{
				{
					"role": openai.ChatCompletionMessageParamRoleDeveloper,
					"content": []map[string]interface{}{
						{
							"type": "input_text",
							"text": strings.Join(developerMessage, "\n"),
						},
					},
				},
				{
					"role": openai.ChatCompletionMessageParamRoleUser,
					"content": []map[string]interface{}{
						{
							"type": "input_text",
							"text": "Get ready to review the " + *projectTitle + " project!",
						},
					},
				},
			},
			"tools": []map[string]interface{}{
				{
					"type":             openai.AssistantToolTypeFileSearch,
					"vector_store_ids": []string{store.ID},
				},
			},
		}

		err = client.Post(context.TODO(), "/responses", responseParams, &responseAPIResult)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to get response from /v1/responses")
			continue
		}

		log.Info().Str("response_id", responseAPIResult.ID).Msg("Response received")
		log.Info().Str("response_message", responseAPIResult.Output[0].Content[0].Text).Msg("Response message")

		// Update previousResponseID to the current response ID
		previousResponseID = responseAPIResult.ID
	}

	_, err = client.Beta.VectorStores.Delete(ctx, store.ID)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to delete vector store")
	}

	log.Info().Msg("Analysis process completed")
}
