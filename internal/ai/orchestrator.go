package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vignesh245/ClusterMind/internal/model"
)

// Orchestrator routes requests to the active AI provider.
type Orchestrator interface {
	Explain(ctx context.Context, pkg *model.EvidencePackage) (*model.ExplainResult, error)
	SuggestRemediation(ctx context.Context, pkg *model.EvidencePackage) (*model.RemediationPlan, error)
	ResolveIntent(ctx context.Context, query string, candidates []model.Resource) (*model.IntentResult, error)
}

// Provider represents an AI model backend.
type Provider interface {
	Name() string
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
	HealthCheck(ctx context.Context) error
}

type CompletionRequest struct {
	SystemPrompt string
	UserPrompt   string
	MaxTokens    int
	Temperature  float32
	Params       map[string]any
}

type CompletionResponse struct {
	Content      string
	FinishReason string
	ProviderMeta map[string]any
}

type orchestrator struct {
	provider Provider
}

// NewOrchestrator creates a new AI orchestrator using the given provider.
func NewOrchestrator(p Provider) Orchestrator {
	return &orchestrator{
		provider: p,
	}
}

func (o *orchestrator) Explain(ctx context.Context, pkg *model.EvidencePackage) (*model.ExplainResult, error) {
	// 1. Construct Prompt
	systemPrompt := `You are a Kubernetes diagnostics assistant. You receive structured evidence
about a failing workload and produce a root cause analysis in JSON.

Rules:
- Base all conclusions on provided evidence only. Do not speculate.
- If evidence is insufficient, set confidence to "low" and state why.
- Output must match the exact schema. No markdown outside the JSON block.`

	pkgJSON, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal evidence: %w", err)
	}
	userPrompt := fmt.Sprintf("Evidence:\n%s\n\nProvide the JSON response:", string(pkgJSON))

	// 2. Call Provider
	resp, err := o.provider.Complete(ctx, CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  0.2,
		MaxTokens:    1024,
	})
	if err != nil {
		return nil, fmt.Errorf("provider error: %w", err)
	}

	// 3. Enforce Schema
	var result model.ExplainResult
	err = json.Unmarshal([]byte(resp.Content), &result)
	if err != nil {
		// Basic fallback/retry logic would go here.
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &result, nil
}

func (o *orchestrator) SuggestRemediation(ctx context.Context, pkg *model.EvidencePackage) (*model.RemediationPlan, error) {
	systemPrompt := `You are a Kubernetes SRE. Based on the evidence provided, 
suggest a single remediation action in JSON format.
Include RemediationType (CommandSuggestion/PatchSuggestion), Rationale, RiskLevel (Low/Medium/High), and ProposedCommand (if CommandSuggestion) or ProposedPatch.
Schema must exactly match model.RemediationPlan.`

	pkgJSON, _ := json.MarshalIndent(pkg, "", "  ")
	userPrompt := fmt.Sprintf("Evidence:\n%s\n\nProvide the JSON response:", string(pkgJSON))

	resp, err := o.provider.Complete(ctx, CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  0.1,
		MaxTokens:    500,
	})
	if err != nil {
		return nil, fmt.Errorf("provider error: %w", err)
	}

	var plan model.RemediationPlan
	if err := json.Unmarshal([]byte(resp.Content), &plan); err != nil {
		return nil, fmt.Errorf("failed to parse remediation plan: %w", err)
	}

	return &plan, nil
}

func (o *orchestrator) ResolveIntent(ctx context.Context, query string, candidates []model.Resource) (*model.IntentResult, error) {
	// Not implemented in Phase 2
	return nil, fmt.Errorf("not implemented")
}
