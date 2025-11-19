package jobs

import (
	"encoding/json"
	"time"

	"deep-research-agent/models"
)

// ResearchJobRequest represents a request to start a research job
type ResearchJobRequest struct {
	Query      string              `json:"query"`
	Documents  []string            `json:"documents"`
	UseWeb     bool                `json:"use_web"`
	UseWiki    bool                `json:"use_wiki"`
	Depth      models.ResearchDepth `json:"depth"`
	MaxSources int                 `json:"max_sources"`
	MaxSteps   int                 `json:"max_steps"`
}

// ToResearchQuery converts a request to a ResearchQuery
func (r *ResearchJobRequest) ToResearchQuery() models.ResearchQuery {
	return models.ResearchQuery{
		Query:      r.Query,
		Documents:  r.Documents,
		UseWeb:     r.UseWeb,
		UseWiki:    r.UseWiki,
		Depth:      r.Depth,
		MaxSources: r.MaxSources,
		MaxSteps:   r.MaxSteps,
	}
}

// ResearchJobData wraps the research query with additional job metadata
type ResearchJobData struct {
	Query     models.ResearchQuery
	CreatedAt time.Time
}

// MarshalJobData serializes job data to JSON
func MarshalJobData(query models.ResearchQuery) (string, error) {
	data := ResearchJobData{
		Query:     query,
		CreatedAt: time.Now(),
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// UnmarshalJobData deserializes job data from JSON
func UnmarshalJobData(jsonData string) (*models.ResearchQuery, error) {
	var data ResearchJobData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}
	return &data.Query, nil
}

// MarshalJobResult serializes research result to JSON
func MarshalJobResult(result *models.ResearchResult) (string, error) {
	bytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// UnmarshalJobResult deserializes research result from JSON
func UnmarshalJobResult(jsonData string) (*models.ResearchResult, error) {
	var result models.ResearchResult
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
