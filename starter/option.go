package starter

import (
	"time"

	"github.com/gogo/protobuf/proto"
	enums "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

type Options struct {
	client.StartWorkflowOptions
	Result proto.Message
}

type Option func(*Options)

func GetResult(res proto.Message) Option {
	return func(o *Options) {
		o.Result = res
	}
}

// ID sets the workflow ID.
func ID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

// TaskQueue sets the task queue for the workflow.
func TaskQueue(taskQueue string) Option {
	return func(o *Options) {
		o.TaskQueue = taskQueue
	}
}

// WorkflowExecutionTimeout sets the timeout for the entire workflow execution.
func WorkflowExecutionTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.WorkflowExecutionTimeout = timeout
	}
}

// WorkflowRunTimeout sets the timeout for a single workflow run.
func WorkflowRunTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.WorkflowRunTimeout = timeout
	}
}

// WorkflowTaskTimeout sets the timeout for a workflow task.
func WorkflowTaskTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.WorkflowTaskTimeout = timeout
	}
}

// WorkflowIDReusePolicy sets the workflow ID reuse policy.
func WorkflowIDReusePolicy(policy enums.WorkflowIdReusePolicy) Option {
	return func(o *Options) {
		o.WorkflowIDReusePolicy = policy
	}
}

// WorkflowIDConflictPolicy sets the workflow ID conflict policy.
func WorkflowIDConflictPolicy(policy enums.WorkflowIdConflictPolicy) Option {
	return func(o *Options) {
		o.WorkflowIDConflictPolicy = policy
	}
}

// WorkflowExecutionErrorWhenAlreadyStarted sets whether to return an error when the workflow is already running.
func WorkflowExecutionErrorWhenAlreadyStarted(b bool) Option {
	return func(o *Options) {
		o.WorkflowExecutionErrorWhenAlreadyStarted = b
	}
}

// RetryPolicy sets the retry policy for the workflow.
func RetryPolicy(rp *temporal.RetryPolicy) Option {
	return func(o *Options) {
		o.RetryPolicy = rp
	}
}

// CronSchedule sets the cron schedule for the workflow.
func CronSchedule(cronSchedule string) Option {
	return func(o *Options) {
		o.CronSchedule = cronSchedule
	}
}

// Memo sets the memo for the workflow.
func Memo(memo map[string]any) Option {
	return func(o *Options) {
		o.Memo = memo
	}
}

// TypedSearchAttributes sets the typed search attributes for the workflow.
func TypedSearchAttributes(searchAttributes temporal.SearchAttributes) Option {
	return func(o *Options) {
		o.TypedSearchAttributes = searchAttributes
	}
}

// EnableEagerStart requests eager execution for this workflow, if a local worker is available.
func EnableEagerStart(b bool) Option {
	return func(o *Options) {
		o.EnableEagerStart = b
	}
}

// StartDelay sets the delay before dispatching the first workflow task.
func StartDelay(delay time.Duration) Option {
	return func(o *Options) {
		o.StartDelay = delay
	}
}

// StaticSummary sets the static summary for the workflow.
func StaticSummary(summary string) Option {
	return func(o *Options) {
		o.StaticSummary = summary
	}
}

// StaticDetails sets the static details for the workflow.
func StaticDetails(details string) Option {
	return func(o *Options) {
		o.StaticDetails = details
	}
}

// VersioningOverride sets the versioning override for the workflow.
func VersioningOverride(override client.VersioningOverride) Option {
	return func(o *Options) {
		o.VersioningOverride = override
	}
}

// Priority sets the priority for the workflow.
func Priority(priority temporal.Priority) Option {
	return func(o *Options) {
		o.Priority = priority
	}
}

// NewOptions creates a new StartWorkflowOptions with the given options.
func NewOptions(taskQueue string, opts ...Option) Options {
	options := Options{
		StartWorkflowOptions: client.StartWorkflowOptions{
			TaskQueue: taskQueue,
		},
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
