package starter

import (
	"time"

	enums "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

type Option func(*client.StartWorkflowOptions)

// ID sets the workflow ID.
func ID(id string) Option {
	return func(o *client.StartWorkflowOptions) {
		o.ID = id
	}
}

// TaskQueue sets the task queue for the workflow.
func TaskQueue(taskQueue string) Option {
	return func(o *client.StartWorkflowOptions) {
		o.TaskQueue = taskQueue
	}
}

// WorkflowExecutionTimeout sets the timeout for the entire workflow execution.
func WorkflowExecutionTimeout(timeout time.Duration) Option {
	return func(o *client.StartWorkflowOptions) {
		o.WorkflowExecutionTimeout = timeout
	}
}

// WorkflowRunTimeout sets the timeout for a single workflow run.
func WorkflowRunTimeout(timeout time.Duration) Option {
	return func(o *client.StartWorkflowOptions) {
		o.WorkflowRunTimeout = timeout
	}
}

// WorkflowTaskTimeout sets the timeout for a workflow task.
func WorkflowTaskTimeout(timeout time.Duration) Option {
	return func(o *client.StartWorkflowOptions) {
		o.WorkflowTaskTimeout = timeout
	}
}

// WorkflowIDReusePolicy sets the workflow ID reuse policy.
func WorkflowIDReusePolicy(policy enums.WorkflowIdReusePolicy) Option {
	return func(o *client.StartWorkflowOptions) {
		o.WorkflowIDReusePolicy = policy
	}
}

// WorkflowIDConflictPolicy sets the workflow ID conflict policy.
func WorkflowIDConflictPolicy(policy enums.WorkflowIdConflictPolicy) Option {
	return func(o *client.StartWorkflowOptions) {
		o.WorkflowIDConflictPolicy = policy
	}
}

// WorkflowExecutionErrorWhenAlreadyStarted sets whether to return an error when the workflow is already running.
func WorkflowExecutionErrorWhenAlreadyStarted(b bool) Option {
	return func(o *client.StartWorkflowOptions) {
		o.WorkflowExecutionErrorWhenAlreadyStarted = b
	}
}

// RetryPolicy sets the retry policy for the workflow.
func RetryPolicy(rp *temporal.RetryPolicy) Option {
	return func(o *client.StartWorkflowOptions) {
		o.RetryPolicy = rp
	}
}

// CronSchedule sets the cron schedule for the workflow.
func CronSchedule(cronSchedule string) Option {
	return func(o *client.StartWorkflowOptions) {
		o.CronSchedule = cronSchedule
	}
}

// Memo sets the memo for the workflow.
func Memo(memo map[string]any) Option {
	return func(o *client.StartWorkflowOptions) {
		o.Memo = memo
	}
}

// SearchAttributes sets the search attributes for the workflow.
// Deprecated: use TypedSearchAttributes instead.
func SearchAttributes(searchAttributes map[string]any) Option {
	return func(o *client.StartWorkflowOptions) {
		o.SearchAttributes = searchAttributes
	}
}

// TypedSearchAttributes sets the typed search attributes for the workflow.
func TypedSearchAttributes(searchAttributes temporal.SearchAttributes) Option {
	return func(o *client.StartWorkflowOptions) {
		o.TypedSearchAttributes = searchAttributes
	}
}

// EnableEagerStart requests eager execution for this workflow, if a local worker is available.
func EnableEagerStart(b bool) Option {
	return func(o *client.StartWorkflowOptions) {
		o.EnableEagerStart = b
	}
}

// StartDelay sets the delay before dispatching the first workflow task.
func StartDelay(delay time.Duration) Option {
	return func(o *client.StartWorkflowOptions) {
		o.StartDelay = delay
	}
}

// StaticSummary sets the static summary for the workflow.
func StaticSummary(summary string) Option {
	return func(o *client.StartWorkflowOptions) {
		o.StaticSummary = summary
	}
}

// StaticDetails sets the static details for the workflow.
func StaticDetails(details string) Option {
	return func(o *client.StartWorkflowOptions) {
		o.StaticDetails = details
	}
}

// VersioningOverride sets the versioning override for the workflow.
func VersioningOverride(override client.VersioningOverride) Option {
	return func(o *client.StartWorkflowOptions) {
		o.VersioningOverride = override
	}
}

// Priority sets the priority for the workflow.
func Priority(priority temporal.Priority) Option {
	return func(o *client.StartWorkflowOptions) {
		o.Priority = priority
	}
}

// NewOptions creates a new StartWorkflowOptions with the given options.
func NewOptions(taskQueue string, opts ...Option) client.StartWorkflowOptions {
	options := client.StartWorkflowOptions{
		TaskQueue: taskQueue,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
