package main

import (
	"concept/pkg/env"
	"concept/pkg/git"
	"concept/pkg/loader"
	"concept/pkg/prompt"
	"concept/pkg/providers"
	"concept/pkg/rules"
	"context"
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// worker processes files using the provided provider
func worker(id int, files <-chan string, client *providers.ProviderClient, rules *rules.Rules, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		log.Info().
			Int("worker_id", id).
			Str("file", file).
			Msg("Processing file")

		messages := []providers.ProviderMessage{}

		// get the file's current git commit
		commit, err := git.GetFileCommit(file)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get file commit")
		}

		log.Trace().Str("commit", commit).Msg("File commit")

		messages = append(messages, providers.ProviderMessage{
			Content: "Current commit: " + commit,
			Role:    providers.ProviderMessageRoleUser,
		})

		// get the file's current git stage
		stage, err := git.GetFileStage(file)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get file stage")
		}

		log.Trace().Str("stage", stage).Msg("File stage")

		messages = append(messages, providers.ProviderMessage{
			Content: "Current stage: " + stage,
			Role:    providers.ProviderMessageRoleUser,
		})

		// create the initial prompt
		var initialPrompt []string = []string{
			"# File Review",
			"",
			"You are an AI agent with expert knowledge in programming.",
			"You'll be given some markdown component rules and a file to review.",
			"",
			"Expectations:",
			"",
			"- ALWAYS review the file against any rules provided.",
			"- ALWAYS determine if there are any changes that need to be made.",
			"- ALWAYS rewrite the entire file, including the changes.",
			"- ALWAYS make changes that are required by the rules.",
			"- NEVER code fence the output.",
			"- NEVER make unnecessary changes.",
			"- ALWAYS reply literally with 'x10Barry__Skipped' if there are no changes.",
			"- ALWAYS reply literally with 'x10Barry__Error' if there is an error.",
			"",
			"Filename: " + file,
			"",
		}

		// Create a new prompt instance
		prompt := prompt.NewPrompt()
		prompt.AppendString(strings.Join(initialPrompt, "\n"))

		// Get matching rules for this file
		matchingRules := rules.GetMatchingRules(file)
		log.Debug().
			Int("worker_id", id).
			Str("file", file).
			Int("matching_rules", len(matchingRules)).
			Msg("Found matching rules")

		prompt.AppendString("Rules:")

		for _, rule := range matchingRules {
			// append the rule content to the prompt
			prompt.AppendString("- " + rule.Path + " (" + rule.Description + ")")

			messages = append(messages, providers.ProviderMessage{
				Content: "Rule: " + rule.Path + " (" + rule.Description + ")\n\n" + "```md\n" + rule.Content + "\n```",
				Role:    providers.ProviderMessageRoleUser,
			})

			log.Debug().
				Int("worker_id", id).
				Str("file", file).
				Str("rule", rule.Description).
				Msg("Processing rule")

			log.Trace().Str("rule_content", rule.Content).Msg("Rule content")
		}

		// Get the file's content
		content, err := os.ReadFile(file)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read file")
			continue
		}

		log.Debug().Str("file", file).Str("content_length", strconv.Itoa(len(content))).Msg("File content")
		log.Trace().Str("content", string(content)).Msg("File content")

		fileToReview := providers.ProviderMessage{
			Content: "File: " + file + "\n\n```\n" + string(content) + "\n```",
			Role:    providers.ProviderMessageRoleUser,
		}
		messages = append(messages, fileToReview)

		systemMessage := providers.ProviderMessage{
			Content: prompt.GetAllAsString(),
			Role:    providers.ProviderMessageRoleSystem,
		}
		messages = append([]providers.ProviderMessage{systemMessage}, messages...)

		log.Debug().Int("num_messages", len(messages)).Msg("Messages")
		log.Trace().Interface("messages", messages).Msg("Messages")

		response, err := client.ChatCompletion(context.Background(), messages)
		if err != nil {
			log.Fatal().Str("file", file).Err(err).Msg("Failed to process file")
		}

		result, err := providers.UnmapProviderMessage(client.ProviderName, response)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmap response")
			continue
		}

		log.Trace().Interface("result", result).Msg("Result")

		if result.Content == "x10Barry__Skipped" {
			log.Info().
				Int("worker_id", id).
				Str("file", file).
				Msg("Skipping file")
			continue
		}

		if result.Content == "x10Barry__Error" {
			log.Error().
				Int("worker_id", id).
				Str("file", file).
				Msg("Error processing file")
			continue
		}

		// create the patch directory if it doesn't exist
		if _, err := os.Stat(".patches"); os.IsNotExist(err) {
			os.Mkdir(".patches", 0755)
			log.Info().Msg("Created .patches directory")
		}

		// write .patches/ to gitignore if it doesn't exist
		if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
			// read the gitignore file
			gitignore, err := os.ReadFile(".gitignore")
			if err != nil {
				log.Error().Err(err).Msg("Failed to read gitignore")
			}

			// check if .patches/ is already in the gitignore file
			if !strings.Contains(string(gitignore), ".patches/") {
				os.WriteFile(".gitignore", []byte(".patches/\n"), 0644)
				log.Info().Msg("Added .patches/ to gitignore")
			}
		}

		// write the result to a file
		os.WriteFile(".patches/"+file+".patch", []byte(result.Content), 0644)
		log.Info().Str("file", file).Str("patch_file", ".patches/"+file+".patch").Msg("Wrote patch to file")

		log.Debug().
			Int("worker_id", id).
			Str("file", file).
			Msg("Completed processing file")
	}
}

func main() {
	var (
		l string
		T string
		p string
		r string
		w int
	)

	flag.StringVar(&p, "p", "openai", "set the provider")
	flag.StringVar(&l, "l", "info", "set log level")
	flag.StringVar(&T, "T", "", "set the title of the project")
	flag.StringVar(&r, "r", ".", "set the root directory")
	flag.IntVar(&w, "w", 10, "number of workers")
	flag.Parse()

	level, err := zerolog.ParseLevel(l)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().
		Str("title", T).
		Str("provider", p).
		Str("log_level", l).
		Int("workers", w).
		Msg("Starting the project")

	// load environment variables
	env.Load(".env")

	// create a provider
	provider, err := providers.NewClient(p)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create provider")
	}

	log.Info().Str("provider", provider.ProviderName).Msg("Provider created")

	// load all files in the working directory
	files, err := loader.Load(r)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load files")
	}

	log.Info().Int("files", len(files)).Msg("Files loaded")

	// load all rules
	ruleFiles, err := loader.LoadRules()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load rules")
	}

	// Create Rules instance
	rulesInstance, err := rules.New(ruleFiles)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse rules")
	}

	// Create a buffered channel to hold the files
	filesChan := make(chan string, len(files))
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 1; i <= w; i++ {
		wg.Add(1)
		go worker(i, filesChan, provider, rulesInstance, &wg)
	}

	// Send files to the workers
	for _, file := range files {
		filesChan <- file
	}
	close(filesChan) // Close channel to signal no more files

	// Wait for all workers to finish
	wg.Wait()
	log.Info().Msg("All files processed")
}
