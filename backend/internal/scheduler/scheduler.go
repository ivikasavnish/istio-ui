package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ivikasavnish/istio-ui/internal/istio"
	"github.com/robfig/cron/v3"
)

// ActionType defines the type of scheduled action
type ActionType string

const (
	ActionUpdateVirtualService  ActionType = "update_virtualservice"
	ActionUpdateDestinationRule ActionType = "update_destinationrule"
	ActionDeleteVirtualService  ActionType = "delete_virtualservice"
	ActionDeleteDestinationRule ActionType = "delete_destinationrule"
)

// ScheduledAction represents a scheduled action
type ScheduledAction struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schedule    string                 `json:"schedule"` // cron expression
	ActionType  ActionType             `json:"action_type"`
	Namespace   string                 `json:"namespace"`
	ResourceName string                `json:"resource_name"`
	Payload     map[string]interface{} `json:"payload"`
	CreatedAt   time.Time              `json:"created_at"`
	NextRun     time.Time              `json:"next_run"`
	LastRun     *time.Time             `json:"last_run,omitempty"`
	Enabled     bool                   `json:"enabled"`
}

// Scheduler manages scheduled actions
type Scheduler struct {
	cron     *cron.Cron
	client   *istio.Client
	actions  map[string]*ScheduledAction
	entries  map[string]cron.EntryID
	mu       sync.RWMutex
}

// NewScheduler creates a new scheduler
func NewScheduler(client *istio.Client) *Scheduler {
	return &Scheduler{
		cron:    cron.New(),
		client:  client,
		actions: make(map[string]*ScheduledAction),
		entries: make(map[string]cron.EntryID),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// AddAction adds a new scheduled action
func (s *Scheduler) AddAction(action *ScheduledAction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if action.ID == "" {
		action.ID = uuid.New().String()
	}
	action.CreatedAt = time.Now()
	action.Enabled = true

	entryID, err := s.cron.AddFunc(action.Schedule, func() {
		s.executeAction(action)
	})
	if err != nil {
		return fmt.Errorf("failed to schedule action: %w", err)
	}

	s.actions[action.ID] = action
	s.entries[action.ID] = entryID

	return nil
}

// RemoveAction removes a scheduled action
func (s *Scheduler) RemoveAction(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, exists := s.entries[id]
	if !exists {
		return fmt.Errorf("action not found: %s", id)
	}

	s.cron.Remove(entryID)
	delete(s.actions, id)
	delete(s.entries, id)

	return nil
}

// ListActions returns all scheduled actions
func (s *Scheduler) ListActions() []*ScheduledAction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	actions := make([]*ScheduledAction, 0, len(s.actions))
	for _, action := range s.actions {
		actions = append(actions, action)
	}

	return actions
}

// GetAction returns a specific action by ID
func (s *Scheduler) GetAction(id string) (*ScheduledAction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	action, exists := s.actions[id]
	if !exists {
		return nil, fmt.Errorf("action not found: %s", id)
	}

	return action, nil
}

// executeAction executes a scheduled action
func (s *Scheduler) executeAction(action *ScheduledAction) {
	ctx := context.Background()
	now := time.Now()

	s.mu.Lock()
	action.LastRun = &now
	s.mu.Unlock()

	switch action.ActionType {
	case ActionUpdateVirtualService:
		s.updateVirtualService(ctx, action)
	case ActionUpdateDestinationRule:
		s.updateDestinationRule(ctx, action)
	case ActionDeleteVirtualService:
		s.deleteVirtualService(ctx, action)
	case ActionDeleteDestinationRule:
		s.deleteDestinationRule(ctx, action)
	}
}

func (s *Scheduler) updateVirtualService(ctx context.Context, action *ScheduledAction) {
	vs, err := s.client.GetVirtualService(ctx, action.Namespace, action.ResourceName)
	if err != nil {
		fmt.Printf("Error getting VirtualService: %v\n", err)
		return
	}

	// Apply payload updates to the VirtualService
	payloadBytes, _ := json.Marshal(action.Payload)
	var updates map[string]interface{}
	json.Unmarshal(payloadBytes, &updates)

	// Update spec if provided in payload
	if specData, ok := updates["spec"]; ok {
		specBytes, _ := json.Marshal(specData)
		json.Unmarshal(specBytes, &vs.Spec)
	}

	_, err = s.client.UpdateVirtualService(ctx, vs)
	if err != nil {
		fmt.Printf("Error updating VirtualService: %v\n", err)
		return
	}

	fmt.Printf("Successfully executed scheduled action: %s\n", action.Name)
}

func (s *Scheduler) updateDestinationRule(ctx context.Context, action *ScheduledAction) {
	dr, err := s.client.GetDestinationRule(ctx, action.Namespace, action.ResourceName)
	if err != nil {
		fmt.Printf("Error getting DestinationRule: %v\n", err)
		return
	}

	// Apply payload updates
	payloadBytes, _ := json.Marshal(action.Payload)
	var updates map[string]interface{}
	json.Unmarshal(payloadBytes, &updates)

	if specData, ok := updates["spec"]; ok {
		specBytes, _ := json.Marshal(specData)
		json.Unmarshal(specBytes, &dr.Spec)
	}

	_, err = s.client.UpdateDestinationRule(ctx, dr)
	if err != nil {
		fmt.Printf("Error updating DestinationRule: %v\n", err)
		return
	}

	fmt.Printf("Successfully executed scheduled action: %s\n", action.Name)
}

func (s *Scheduler) deleteVirtualService(ctx context.Context, action *ScheduledAction) {
	err := s.client.DeleteVirtualService(ctx, action.Namespace, action.ResourceName)
	if err != nil {
		fmt.Printf("Error deleting VirtualService: %v\n", err)
		return
	}

	fmt.Printf("Successfully executed scheduled action: %s\n", action.Name)
	// Remove this action after execution since the resource is deleted
	s.RemoveAction(action.ID)
}

func (s *Scheduler) deleteDestinationRule(ctx context.Context, action *ScheduledAction) {
	err := s.client.DeleteDestinationRule(ctx, action.Namespace, action.ResourceName)
	if err != nil {
		fmt.Printf("Error deleting DestinationRule: %v\n", err)
		return
	}

	fmt.Printf("Successfully executed scheduled action: %s\n", action.Name)
	// Remove this action after execution
	s.RemoveAction(action.ID)
}

// Helper to validate cron expression
func ValidateCronExpression(expr string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expr)
	return err
}
