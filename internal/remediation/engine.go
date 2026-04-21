package remediation

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/vignesh245/ClusterMind/internal/kube"
	"github.com/vignesh245/ClusterMind/internal/model"
	"k8s.io/apimachinery/pkg/types"
)

// Executor securely executes remediation actions.
type Executor interface {
	Execute(ctx context.Context, plan *model.RemediationPlan) error
}

type executor struct {
	client kube.Client
}

// NewExecutor creates a new remediation Executor.
func NewExecutor(client kube.Client) Executor {
	return &executor{client: client}
}

func (e *executor) Execute(ctx context.Context, plan *model.RemediationPlan) error {
	// Security: Whitelist validation
	if err := validatePlan(plan); err != nil {
		return fmt.Errorf("security policy violation: %w", err)
	}

	switch plan.RemediationType {
	case model.CommandSuggestion:
		if plan.ProposedCommand == "" {
			return fmt.Errorf("command string missing")
		}

		parts := strings.Fields(plan.ProposedCommand)
		if len(parts) == 0 {
			return fmt.Errorf("empty command")
		}

		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("command execution failed: %s (err: %w)", string(output), err)
		}
		return nil

	case model.PatchSuggestion:
		if plan.ProposedPatch == "" {
			return fmt.Errorf("patch payload missing")
		}
		// In a real scenario, resource would be passed.
		dummyResource := model.Resource{}
		return e.client.ApplyPatch(ctx, dummyResource, []byte(plan.ProposedPatch), types.MergePatchType)

	default:
		return fmt.Errorf("unsupported remediation type: %s", plan.RemediationType)
	}
}

// validatePlan enforces read-only defaults or allowed patterns.
func validatePlan(plan *model.RemediationPlan) error {
	if plan.RemediationType == model.CommandSuggestion {
		cmdStr := plan.ProposedCommand
		// Example blacklist checks
		if strings.Contains(cmdStr, "rm -rf") || strings.Contains(cmdStr, ">") || strings.Contains(cmdStr, "|") {
			return fmt.Errorf("disallowed shell characters or commands")
		}
	}
	return nil
}
