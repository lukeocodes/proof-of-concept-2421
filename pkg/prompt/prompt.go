package prompt

import "bytes"

// Prompt holds a collection of byte slices
type Prompt struct {
	data [][]byte
}

// NewPrompt creates a new instance of Prompt
func NewPrompt() *Prompt {
	return &Prompt{
		data: make([][]byte, 0),
	}
}

// Append adds a new byte slice to the store
func (p *Prompt) Append(bytes []byte) {
	p.data = append(p.data, bytes)
}

// Appends a string to the byte store as a new byte slice
func (p *Prompt) AppendString(s string) {
	p.data = append(p.data, []byte(s))
}

// GetAll returns all stored byte slices
func (p *Prompt) GetAll() [][]byte {
	return p.data
}

// GetAllAsString returns all stored byte slices as a single string
func (p *Prompt) GetAllAsString() string {
	return string(bytes.Join(p.GetAll(), []byte("\n")))
}
