package workflows

import (
	"context"

	"agent-commission/workflows/activities"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

const (
	// TaskQueueAgentOnboarding is the task queue for agent onboarding workflows
	TaskQueueAgentOnboarding = "agent-onboarding"
)

// WorkerParams defines the dependencies for Temporal worker
type WorkerParams struct {
	fx.In
	Client         client.Client
	AgentActivities *activities.AgentActivities
}

// NewTemporalWorker creates and starts a Temporal worker
func NewTemporalWorker(params WorkerParams) worker.Worker {
	// Create worker for agent onboarding task queue
	w := worker.New(params.Client, TaskQueueAgentOnboarding, worker.Options{})

	// Register workflows
	w.RegisterWorkflow(AgentOnboardingWorkflow)

	// Register activities
	w.RegisterActivity(params.AgentActivities.ValidateAgentInputActivity)
	w.RegisterActivity(params.AgentActivities.CheckPANExistsActivity)
	w.RegisterActivity(params.AgentActivities.FetchEmployeeFromHRMSActivity)
	w.RegisterActivity(params.AgentActivities.ValidateCoordinatorActivity)
	w.RegisterActivity(params.AgentActivities.GenerateAgentCodeActivity)
	w.RegisterActivity(params.AgentActivities.CreateAgentProfileActivity)
	w.RegisterActivity(params.AgentActivities.CreateHierarchyActivity)
	w.RegisterActivity(params.AgentActivities.SendOnboardingNotificationActivity)

	return w
}

// StartWorker starts the Temporal worker
func StartWorker(lc fx.Lifecycle, w worker.Worker) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start worker in a goroutine
			go func() {
				if err := w.Run(worker.InterruptCh()); err != nil {
					// Log error but don't stop the application
					// In production, implement proper error handling and restart logic
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			w.Stop()
			return nil
		},
	})
}
