package model

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime"
)

// Resource represents a generic Kubernetes resource reference
// used across ClusterMind modules (UI, Remediation, Intent, etc.).
type Resource struct {
	Kind      string
	Name      string
	Namespace string
	Object    runtime.Object
}

// ExplainResult is the final output of the AI orchestrator's RCA flow.
type ExplainResult struct {
	Summary            string          `json:"summary"`
	LikelyRootCause    string          `json:"likely_root_cause"`
	Evidence           []EvidenceRef   `json:"evidence"`
	Confidence         ConfidenceLevel `json:"confidence"`
	RecommendedActions []Action        `json:"recommended_actions"`
}

type EvidenceRef struct {
	Ref         string `json:"ref"`
	Description string `json:"description"`
}

type ConfidenceLevel string

const (
	ConfidenceLow    ConfidenceLevel = "low"
	ConfidenceMedium ConfidenceLevel = "medium"
	ConfidenceHigh   ConfidenceLevel = "high"
)

type Action struct {
	Description string  `json:"description"`
	Command     *string `json:"command"`
}

// IntentResult is returned by the orchestrator when resolving natural language to queries.
type IntentResult struct {
	MatchedResourceKind string
	Filters             map[string]string
}

// RemediationPlan represents an actionable set of fixes approved by AI.
type RemediationPlan struct {
	RemediationType RemediationType
	Rationale       string
	ProposedCommand string
	ProposedPatch   string
	ExpectedOutcome string
	RiskNotes       string
	RiskLevel       RiskLevel
}

type RemediationType string

const (
	CommandSuggestion RemediationType = "CommandSuggestion"
	PatchSuggestion   RemediationType = "PatchSuggestion"
	RunbookSuggestion RemediationType = "RunbookSuggestion"
)

type RiskLevel string

const (
	RiskLow         RiskLevel = "Low"
	RiskMedium      RiskLevel = "Medium"
	RiskHigh        RiskLevel = "High"
	RiskDestructive RiskLevel = "Destructive"
)

// EvidencePackage is the normalized context given to the AI.
type EvidencePackage struct {
	ResourceKind     string
	ResourceName     string
	Namespace        string
	StatusConditions []Condition
	RecentEvents     []EventSummary
	LogExcerpt       string
	RestartHistory   []ContainerRestart
	ProbeConfig      *ProbeConfig
	OwnerChain       []OwnerRef
	MetricsSummary   *MetricsSummary // Optional, nil if not available
	AnalyzerFindings []Finding
	RedactedFields   []string
}

type Condition struct {
	Type               string
	Status             string
	Reason             string
	Message            string
	LastTransitionTime time.Time
}

type EventSummary struct {
	Reason    string
	Message   string
	Count     int32
	FirstSeen time.Time
	LastSeen  time.Time
	Type      string
}

type ContainerRestart struct {
	Container string
	Count     int32
	Reason    string
	ExitCode  int32
	FinishedAt time.Time
}

type ProbeConfig map[string]string

type OwnerRef struct {
	Kind string
	Name string
}

type MetricsSummary struct {
	CPUUsage    string
	MemoryUsage string
}

type Finding struct {
	Severity       Severity
	Category       Category
	Title          string
	Detail         string
	SourceAnalyzer string
	Evidence       []string
}

type Severity string

const (
	SeverityInfo     Severity = "Info"
	SeverityWarning  Severity = "Warning"
	SeverityCritical Severity = "Critical"
)

type Category string

const (
	CategoryScheduling Category = "Scheduling"
	CategoryRuntime    Category = "Runtime"
	CategoryNetwork    Category = "Network"
	CategoryStorage    Category = "Storage"
	CategoryConfig     Category = "Config"
	CategoryResource   Category = "Resource"
)
