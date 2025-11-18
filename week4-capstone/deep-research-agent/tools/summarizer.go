package tools

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"
)

// Summarizer provides text summarization capabilities
type Summarizer struct{}

// NewSummarizer creates a new summarizer tool
func NewSummarizer() *Summarizer {
	return &Summarizer{}
}

// Name returns the tool name
func (s *Summarizer) Name() string {
	return "summarizer"
}

// Description returns the tool description
func (s *Summarizer) Description() string {
	return "Summarize long text using extractive summarization. Generates concise summaries and key points."
}

// Parameters returns the tool parameters
func (s *Summarizer) Parameters() []Parameter {
	return []Parameter{
		{
			Name:        "text",
			Type:        "string",
			Required:    true,
			Description: "Text to summarize",
		},
		{
			Name:        "action",
			Type:        "string",
			Required:    false,
			Description: "Action: summarize, key_points (default: summarize)",
			Default:     "summarize",
		},
		{
			Name:        "max_sentences",
			Type:        "int",
			Required:    false,
			Description: "Maximum number of sentences in summary (default: 5)",
			Default:     5,
		},
		{
			Name:        "compression_ratio",
			Type:        "float",
			Required:    false,
			Description: "Compression ratio for summary (0.1-0.9, default: 0.3)",
			Default:     0.3,
		},
	}
}

// Execute runs the summarizer
func (s *Summarizer) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
	// Extract parameters
	text, ok := params["text"].(string)
	if !ok || text == "" {
		return nil, fmt.Errorf("text parameter is required")
	}

	action := "summarize"
	if a, ok := params["action"].(string); ok {
		action = a
	}

	maxSentences := 5
	if ms, ok := params["max_sentences"].(float64); ok {
		maxSentences = int(ms)
	}

	compressionRatio := 0.3
	if cr, ok := params["compression_ratio"].(float64); ok {
		compressionRatio = cr
	}

	// Execute action
	switch action {
	case "summarize":
		return s.summarize(text, maxSentences, compressionRatio)
	case "key_points":
		return s.extractKeyPoints(text, maxSentences)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// summarize creates a summary using extractive summarization
func (s *Summarizer) summarize(text string, maxSentences int, compressionRatio float64) (*ToolResult, error) {
	// Split into sentences
	sentences := s.splitIntoSentences(text)

	if len(sentences) == 0 {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("no sentences found in text"),
		}, nil
	}

	// If text is already short, return as is
	if len(sentences) <= maxSentences {
		summary := &Summary{
			Original:         text,
			Summary:          text,
			KeyPoints:        sentences,
			CompressionRatio: 1.0,
			OriginalLength:   len(text),
			SummaryLength:    len(text),
		}

		return &ToolResult{
			Success: true,
			Data:    summary,
			Metadata: map[string]interface{}{
				"sentence_count": len(sentences),
				"compression":    1.0,
			},
		}, nil
	}

	// Calculate desired number of sentences
	desiredCount := int(float64(len(sentences)) * compressionRatio)
	if desiredCount < 1 {
		desiredCount = 1
	}
	if desiredCount > maxSentences {
		desiredCount = maxSentences
	}

	// Score sentences
	scores := s.scoreSentences(sentences, text)

	// Get top sentences
	topSentences := s.getTopSentences(sentences, scores, desiredCount)

	// Build summary
	summaryText := strings.Join(topSentences, " ")

	summary := &Summary{
		Original:         text,
		Summary:          summaryText,
		KeyPoints:        topSentences,
		CompressionRatio: float64(len(summaryText)) / float64(len(text)),
		OriginalLength:   len(text),
		SummaryLength:    len(summaryText),
	}

	return &ToolResult{
		Success: true,
		Data:    summary,
		Metadata: map[string]interface{}{
			"original_sentences": len(sentences),
			"summary_sentences":  len(topSentences),
			"compression_ratio":  summary.CompressionRatio,
		},
	}, nil
}

// extractKeyPoints extracts key points from text
func (s *Summarizer) extractKeyPoints(text string, maxPoints int) (*ToolResult, error) {
	sentences := s.splitIntoSentences(text)

	if len(sentences) == 0 {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("no sentences found in text"),
		}, nil
	}

	// Score sentences
	scores := s.scoreSentences(sentences, text)

	// Get top sentences as key points
	count := maxPoints
	if len(sentences) < maxPoints {
		count = len(sentences)
	}

	keyPoints := s.getTopSentences(sentences, scores, count)

	summary := &Summary{
		Original:         text,
		Summary:          strings.Join(keyPoints, "\n• "),
		KeyPoints:        keyPoints,
		CompressionRatio: float64(len(keyPoints)) / float64(len(sentences)),
		OriginalLength:   len(text),
		SummaryLength:    len(strings.Join(keyPoints, " ")),
	}

	return &ToolResult{
		Success: true,
		Data:    summary,
		Metadata: map[string]interface{}{
			"key_point_count": len(keyPoints),
			"total_sentences": len(sentences),
		},
	}, nil
}

// splitIntoSentences splits text into sentences
func (s *Summarizer) splitIntoSentences(text string) []string {
	// Simple sentence splitter
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	// Replace multiple spaces with single space
	text = strings.Join(strings.Fields(text), " ")

	// Split on sentence terminators
	sentences := make([]string, 0)
	current := ""

	for i, char := range text {
		current += string(char)

		// Check for sentence end
		if char == '.' || char == '!' || char == '?' {
			// Make sure it's not an abbreviation
			if i+1 < len(text) {
				nextChar := rune(text[i+1])
				if unicode.IsSpace(nextChar) && (i+2 >= len(text) || unicode.IsUpper(rune(text[i+2]))) {
					sentence := strings.TrimSpace(current)
					if len(sentence) > 10 { // Minimum sentence length
						sentences = append(sentences, sentence)
					}
					current = ""
				}
			} else {
				sentence := strings.TrimSpace(current)
				if len(sentence) > 10 {
					sentences = append(sentences, sentence)
				}
				current = ""
			}
		}
	}

	// Add remaining text
	if current != "" {
		sentence := strings.TrimSpace(current)
		if len(sentence) > 10 {
			sentences = append(sentences, sentence)
		}
	}

	return sentences
}

// scoreSentences scores sentences based on importance
func (s *Summarizer) scoreSentences(sentences []string, fullText string) []float64 {
	scores := make([]float64, len(sentences))

	// Calculate word frequencies
	wordFreq := s.calculateWordFrequency(fullText)

	// Score each sentence
	for i, sentence := range sentences {
		score := 0.0
		words := s.tokenize(sentence)

		// Sum word frequencies
		for _, word := range words {
			score += wordFreq[word]
		}

		// Normalize by sentence length
		if len(words) > 0 {
			score /= float64(len(words))
		}

		// Bonus for position (earlier sentences often more important)
		positionBonus := 1.0 - (float64(i) / float64(len(sentences)) * 0.3)
		score *= positionBonus

		// Bonus for sentence length (prefer medium-length sentences)
		lengthBonus := 1.0
		if len(words) < 5 {
			lengthBonus = 0.5 // Penalize very short sentences
		} else if len(words) > 30 {
			lengthBonus = 0.8 // Slightly penalize very long sentences
		}
		score *= lengthBonus

		scores[i] = score
	}

	return scores
}

// calculateWordFrequency calculates word frequency in text
func (s *Summarizer) calculateWordFrequency(text string) map[string]float64 {
	words := s.tokenize(text)
	freq := make(map[string]int)

	for _, word := range words {
		freq[word]++
	}

	// Normalize frequencies
	maxFreq := 0
	for _, count := range freq {
		if count > maxFreq {
			maxFreq = count
		}
	}

	normalized := make(map[string]float64)
	for word, count := range freq {
		normalized[word] = float64(count) / float64(maxFreq)
	}

	return normalized
}

// tokenize splits text into words
func (s *Summarizer) tokenize(text string) []string {
	// Convert to lowercase and split
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// Filter out stop words and short words
	stopWords := map[string]bool{
		"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
		"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
		"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
		"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
	}

	filtered := make([]string, 0)
	for _, word := range words {
		if len(word) > 2 && !stopWords[word] {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

// getTopSentences returns top N sentences by score, maintaining original order
func (s *Summarizer) getTopSentences(sentences []string, scores []float64, n int) []string {
	// Create sentence-score pairs
	type sentenceScore struct {
		sentence string
		score    float64
		index    int
	}

	pairs := make([]sentenceScore, len(sentences))
	for i := range sentences {
		pairs[i] = sentenceScore{
			sentence: sentences[i],
			score:    scores[i],
			index:    i,
		}
	}

	// Sort by score (descending)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].score > pairs[j].score
	})

	// Get top N
	topN := pairs[:int(math.Min(float64(n), float64(len(pairs))))]

	// Sort by original order
	sort.Slice(topN, func(i, j int) bool {
		return topN[i].index < topN[j].index
	})

	// Extract sentences
	result := make([]string, len(topN))
	for i, pair := range topN {
		result[i] = pair.sentence
	}

	return result
}

// Summary represents a text summary
type Summary struct {
	Original         string   `json:"original"`
	Summary          string   `json:"summary"`
	KeyPoints        []string `json:"key_points"`
	CompressionRatio float64  `json:"compression_ratio"`
	OriginalLength   int      `json:"original_length"`
	SummaryLength    int      `json:"summary_length"`
}
