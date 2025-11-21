package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ivikasavnish/istio-ui/backend/internal/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// InitSchema initializes the database schema
func (s *Store) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS snapshots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		config TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS scheduled_actions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		action_type TEXT NOT NULL,
		cron_expr TEXT NOT NULL,
		config TEXT NOT NULL,
		enabled BOOLEAN DEFAULT 1,
		next_run TIMESTAMP,
		last_run TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		username TEXT NOT NULL,
		action TEXT NOT NULL,
		resource TEXT NOT NULL,
		details TEXT,
		success BOOLEAN NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS service_cache (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		namespace TEXT NOT NULL,
		name TEXT NOT NULL,
		version TEXT,
		subsets TEXT,
		last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(namespace, name, version)
	);

	CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp DESC);
	CREATE INDEX IF NOT EXISTS idx_scheduled_actions_enabled ON scheduled_actions(enabled, next_run);
	CREATE INDEX IF NOT EXISTS idx_service_cache_namespace ON service_cache(namespace, name);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Create default admin user if not exists
	_, err = s.db.Exec(`
		INSERT OR IGNORE INTO users (username, email, role)
		VALUES ('admin', 'admin@meshcontrol.local', 'admin')
	`)
	
	return err
}

// User operations
func (s *Store) GetUser(username string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(`
		SELECT id, username, email, role, created_at
		FROM users WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (s *Store) ListUsers() ([]models.User, error) {
	rows, err := s.db.Query(`
		SELECT id, username, email, role, created_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// Snapshot operations
func (s *Store) CreateSnapshot(snapshot *models.Snapshot) error {
	result, err := s.db.Exec(`
		INSERT INTO snapshots (name, description, config, created_by)
		VALUES (?, ?, ?, ?)
	`, snapshot.Name, snapshot.Description, snapshot.Config, snapshot.CreatedBy)
	
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	snapshot.ID = int(id)
	return nil
}

func (s *Store) GetSnapshot(id int) (*models.Snapshot, error) {
	snapshot := &models.Snapshot{}
	err := s.db.QueryRow(`
		SELECT id, name, description, config, created_at, created_by
		FROM snapshots WHERE id = ?
	`, id).Scan(&snapshot.ID, &snapshot.Name, &snapshot.Description, 
		&snapshot.Config, &snapshot.CreatedAt, &snapshot.CreatedBy)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return snapshot, err
}

func (s *Store) ListSnapshots() ([]models.Snapshot, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, config, created_at, created_by
		FROM snapshots ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.Snapshot
	for rows.Next() {
		var snapshot models.Snapshot
		if err := rows.Scan(&snapshot.ID, &snapshot.Name, &snapshot.Description,
			&snapshot.Config, &snapshot.CreatedAt, &snapshot.CreatedBy); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, rows.Err()
}

func (s *Store) DeleteSnapshot(id int) error {
	_, err := s.db.Exec("DELETE FROM snapshots WHERE id = ?", id)
	return err
}

// ScheduledAction operations
func (s *Store) CreateScheduledAction(action *models.ScheduledAction) error {
	result, err := s.db.Exec(`
		INSERT INTO scheduled_actions (name, action_type, cron_expr, config, enabled, next_run)
		VALUES (?, ?, ?, ?, ?, ?)
	`, action.Name, action.ActionType, action.CronExpr, action.Config, action.Enabled, action.NextRun)
	
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	action.ID = int(id)
	return nil
}

func (s *Store) GetScheduledAction(id int) (*models.ScheduledAction, error) {
	action := &models.ScheduledAction{}
	var lastRun sql.NullTime
	
	err := s.db.QueryRow(`
		SELECT id, name, action_type, cron_expr, config, enabled, next_run, last_run, created_at, updated_at
		FROM scheduled_actions WHERE id = ?
	`, id).Scan(&action.ID, &action.Name, &action.ActionType, &action.CronExpr,
		&action.Config, &action.Enabled, &action.NextRun, &lastRun, &action.CreatedAt, &action.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if lastRun.Valid {
		action.LastRun = lastRun.Time
	}
	
	return action, err
}

func (s *Store) ListScheduledActions() ([]models.ScheduledAction, error) {
	rows, err := s.db.Query(`
		SELECT id, name, action_type, cron_expr, config, enabled, next_run, last_run, created_at, updated_at
		FROM scheduled_actions ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []models.ScheduledAction
	for rows.Next() {
		var action models.ScheduledAction
		var lastRun sql.NullTime
		
		if err := rows.Scan(&action.ID, &action.Name, &action.ActionType, &action.CronExpr,
			&action.Config, &action.Enabled, &action.NextRun, &lastRun, &action.CreatedAt, &action.UpdatedAt); err != nil {
			return nil, err
		}
		
		if lastRun.Valid {
			action.LastRun = lastRun.Time
		}
		
		actions = append(actions, action)
	}
	return actions, rows.Err()
}

func (s *Store) UpdateScheduledAction(action *models.ScheduledAction) error {
	_, err := s.db.Exec(`
		UPDATE scheduled_actions
		SET name = ?, action_type = ?, cron_expr = ?, config = ?, enabled = ?, 
		    next_run = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, action.Name, action.ActionType, action.CronExpr, action.Config, 
		action.Enabled, action.NextRun, action.ID)
	return err
}

func (s *Store) UpdateScheduledActionLastRun(id int, lastRun time.Time) error {
	_, err := s.db.Exec(`
		UPDATE scheduled_actions SET last_run = ? WHERE id = ?
	`, lastRun, id)
	return err
}

func (s *Store) DeleteScheduledAction(id int) error {
	_, err := s.db.Exec("DELETE FROM scheduled_actions WHERE id = ?", id)
	return err
}

func (s *Store) GetPendingScheduledActions() ([]models.ScheduledAction, error) {
	rows, err := s.db.Query(`
		SELECT id, name, action_type, cron_expr, config, enabled, next_run, last_run, created_at, updated_at
		FROM scheduled_actions 
		WHERE enabled = 1 AND next_run <= ?
		ORDER BY next_run ASC
	`, time.Now())
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []models.ScheduledAction
	for rows.Next() {
		var action models.ScheduledAction
		var lastRun sql.NullTime
		
		if err := rows.Scan(&action.ID, &action.Name, &action.ActionType, &action.CronExpr,
			&action.Config, &action.Enabled, &action.NextRun, &lastRun, &action.CreatedAt, &action.UpdatedAt); err != nil {
			return nil, err
		}
		
		if lastRun.Valid {
			action.LastRun = lastRun.Time
		}
		
		actions = append(actions, action)
	}
	return actions, rows.Err()
}

// AuditLog operations
func (s *Store) CreateAuditLog(log *models.AuditLog) error {
	result, err := s.db.Exec(`
		INSERT INTO audit_logs (user_id, username, action, resource, details, success)
		VALUES (?, ?, ?, ?, ?, ?)
	`, log.UserID, log.Username, log.Action, log.Resource, log.Details, log.Success)
	
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	log.ID = int(id)
	return nil
}

func (s *Store) ListAuditLogs(limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	
	rows, err := s.db.Query(`
		SELECT id, user_id, username, action, resource, details, success, timestamp
		FROM audit_logs 
		ORDER BY timestamp DESC
		LIMIT ?
	`, limit)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.Username, &log.Action,
			&log.Resource, &log.Details, &log.Success, &log.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

// ServiceCache operations
func (s *Store) UpsertServiceCache(service *models.ServiceCache) error {
	subsetsJSON, err := json.Marshal(service.Subsets)
	if err != nil {
		return err
	}
	
	_, err = s.db.Exec(`
		INSERT INTO service_cache (namespace, name, version, subsets, last_seen)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(namespace, name, version) 
		DO UPDATE SET 
			subsets = excluded.subsets,
			last_seen = excluded.last_seen,
			updated_at = CURRENT_TIMESTAMP
	`, service.Namespace, service.Name, service.Version, string(subsetsJSON), service.LastSeen)
	
	return err
}

func (s *Store) ListServiceCache() ([]models.ServiceCache, error) {
	rows, err := s.db.Query(`
		SELECT id, namespace, name, version, subsets, last_seen, updated_at
		FROM service_cache 
		ORDER BY namespace, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.ServiceCache
	for rows.Next() {
		var service models.ServiceCache
		if err := rows.Scan(&service.ID, &service.Namespace, &service.Name,
			&service.Version, &service.Subsets, &service.LastSeen, &service.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, rows.Err()
}

func (s *Store) DeleteOldServiceCache(olderThan time.Time) error {
	_, err := s.db.Exec("DELETE FROM service_cache WHERE last_seen < ?", olderThan)
	return err
}
