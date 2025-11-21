package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	
	"github.com/ivikasavnish/istio-ui/backend/internal/kube"
	"github.com/ivikasavnish/istio-ui/backend/internal/storage"
)

type Scheduler struct {
	store      *storage.Store
	kubeClient *kube.Client
	cron       *cron.Cron
}

func NewScheduler(store *storage.Store, kubeClient *kube.Client) *Scheduler {
	return &Scheduler{
		store:      store,
		kubeClient: kubeClient,
		cron:       cron.New(),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	// Add periodic job to check for pending actions
	s.cron.AddFunc("* * * * *", s.checkPendingActions) // Every minute
	s.cron.Start()
	log.Println("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Scheduler stopped")
}

// checkPendingActions checks for pending scheduled actions
func (s *Scheduler) checkPendingActions() {
	actions, err := s.store.GetPendingScheduledActions()
	if err != nil {
		log.Printf("Failed to get pending actions: %v", err)
		return
	}

	for _, action := range actions {
		if err := s.executeAction(&action); err != nil {
			log.Printf("Failed to execute action %s: %v", action.Name, err)
		} else {
			log.Printf("Successfully executed action: %s", action.Name)
			
			// Update last run time
			s.store.UpdateScheduledActionLastRun(action.ID, time.Now())
			
			// Calculate next run time
			if action.CronExpr != "" {
				schedule, err := cron.ParseStandard(action.CronExpr)
				if err == nil {
					nextRun := schedule.Next(time.Now())
					action.NextRun = nextRun
					s.store.UpdateScheduledAction(&action)
				}
			}
		}
	}
}

// executeAction executes a scheduled action
func (s *Scheduler) executeAction(action *storage.ScheduledAction) error {
	ctx := context.Background()
	
	log.Printf("Executing action: %s (type: %s)", action.Name, action.ActionType)
	
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(action.Config), &config); err != nil {
		return err
	}
	
	switch action.ActionType {
	case "traffic_shift":
		return s.executeTrafficShift(ctx, config)
	case "apply_mtls":
		return s.applyMTLS(ctx, config)
	case "remove_fault":
		return s.removeFaultInjection(ctx, config)
	case "snapshot_capture":
		return s.captureSnapshot(config)
	default:
		log.Printf("Unknown action type: %s", action.ActionType)
	}
	
	return nil
}

// executeTrafficShift executes a traffic shift action
func (s *Scheduler) executeTrafficShift(ctx context.Context, config map[string]interface{}) error {
	log.Printf("Executing traffic shift with config: %v", config)
	// Implementation would update VirtualService weights
	return nil
}

// applyMTLS applies mTLS configuration
func (s *Scheduler) applyMTLS(ctx context.Context, config map[string]interface{}) error {
	log.Printf("Applying mTLS with config: %v", config)
	// Implementation would create/update PeerAuthentication
	return nil
}

// removeFaultInjection removes fault injection
func (s *Scheduler) removeFaultInjection(ctx context.Context, config map[string]interface{}) error {
	log.Printf("Removing fault injection with config: %v", config)
	// Implementation would update VirtualService to remove fault
	return nil
}

// captureSnapshot captures a configuration snapshot
func (s *Scheduler) captureSnapshot(config map[string]interface{}) error {
	name := config["name"].(string)
	description := config["description"].(string)
	
	log.Printf("Capturing snapshot: %s", name)
	
	// Implementation would capture current Istio configuration
	snapshotConfig := map[string]interface{}{
		"timestamp": time.Now(),
		"resources": []string{}, // Would contain actual resources
	}
	
	configJSON, err := json.Marshal(snapshotConfig)
	if err != nil {
		return err
	}
	
	snapshot := &storage.Snapshot{
		Name:        name,
		Description: description,
		Config:      string(configJSON),
		CreatedBy:   "scheduler",
	}
	
	return s.store.CreateSnapshot(snapshot)
}
