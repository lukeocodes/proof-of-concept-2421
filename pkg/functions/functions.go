package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/openai/openai-go"
)

type FunctionProcessor struct{}

func NewFunctionProcessor() *FunctionProcessor {
	return &FunctionProcessor{}
}

func (fp *FunctionProcessor) ProcessFunctionCall(f openai.ChatCompletionMessageToolCall) (string, error) {
	switch f.Function.Name {
	case "load_file":
		return fp.loadFile(f.Function.Arguments)
	case "download_file":
		return fp.downloadFile(f.Function.Arguments)
	default:
		return "", fmt.Errorf("unknown function: %s", f.Function.Name)
	}
}

type LoadFileArgs struct {
	Filepath string `json:"filepath"`
}

func (fp *FunctionProcessor) loadFile(args string) (string, error) {
	var params LoadFileArgs
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	content, err := os.ReadFile(params.Filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

type DownloadFileArgs struct {
	URL string `json:"url"`
}

func (fp *FunctionProcessor) downloadFile(args string) (string, error) {
	var params DownloadFileArgs
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	resp, err := http.Get(params.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(content), nil
}
