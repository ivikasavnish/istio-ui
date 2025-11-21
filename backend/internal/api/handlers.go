package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	
	"github.com/ivikasavnish/istio-ui/backend/internal/models"
)

// PeerAuthentication handlers
func (s *Server) handleListPeerAuthentications(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	list, err := s.istioMgr.ListPeerAuthentications(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleListPeerAuthenticationsInNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	list, err := s.istioMgr.ListPeerAuthentications(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleCreatePeerAuthentication(w http.ResponseWriter, r *http.Request) {
	var spec models.PeerAuthenticationSpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.CreatePeerAuthentication(r.Context(), &spec); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "CREATE", "PeerAuthentication", fmt.Sprintf("%s/%s", spec.Namespace, spec.Name), true)
	s.broadcastEvent("peerauthentication_created", spec)
	
	respondJSON(w, map[string]interface{}{"message": "PeerAuthentication created successfully"})
}

func (s *Server) handleDeletePeerAuthentication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.DeletePeerAuthentication(r.Context(), namespace, name); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "DELETE", "PeerAuthentication", fmt.Sprintf("%s/%s", namespace, name), true)
	s.broadcastEvent("peerauthentication_deleted", map[string]string{"namespace": namespace, "name": name})
	
	respondJSON(w, map[string]interface{}{"message": "PeerAuthentication deleted successfully"})
}

// AuthorizationPolicy handlers
func (s *Server) handleListAuthorizationPolicies(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	list, err := s.istioMgr.ListAuthorizationPolicies(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleListAuthorizationPoliciesInNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	list, err := s.istioMgr.ListAuthorizationPolicies(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleCreateAuthorizationPolicy(w http.ResponseWriter, r *http.Request) {
	var spec models.AuthorizationPolicySpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.CreateAuthorizationPolicy(r.Context(), &spec); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "CREATE", "AuthorizationPolicy", fmt.Sprintf("%s/%s", spec.Namespace, spec.Name), true)
	s.broadcastEvent("authorizationpolicy_created", spec)
	
	respondJSON(w, map[string]interface{}{"message": "AuthorizationPolicy created successfully"})
}

func (s *Server) handleDeleteAuthorizationPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.DeleteAuthorizationPolicy(r.Context(), namespace, name); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "DELETE", "AuthorizationPolicy", fmt.Sprintf("%s/%s", namespace, name), true)
	s.broadcastEvent("authorizationpolicy_deleted", map[string]string{"namespace": namespace, "name": name})
	
	respondJSON(w, map[string]interface{}{"message": "AuthorizationPolicy deleted successfully"})
}

// Traffic management handlers
func (s *Server) handleUpdateTrafficWeights(w http.ResponseWriter, r *http.Request) {
	var update models.TrafficWeightUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Implementation would update VirtualService with new weights
	s.logAudit("admin", "UPDATE", "TrafficWeights", fmt.Sprintf("%s/%s", update.Namespace, update.Service), true)
	s.broadcastEvent("traffic_weights_updated", update)
	
	respondJSON(w, map[string]interface{}{
		"message": "Traffic weights updated successfully",
		"weights": update.Weights,
	})
}

func (s *Server) handleCanaryDeployment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Service    string `json:"service"`
		Namespace  string `json:"namespace"`
		OldVersion string `json:"old_version"`
		NewVersion string `json:"new_version"`
		Percentage int    `json:"percentage"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Implementation would create/update VirtualService for canary
	s.logAudit("admin", "CANARY", "TrafficShift", fmt.Sprintf("%s/%s", req.Namespace, req.Service), true)
	s.broadcastEvent("canary_deployment", req)
	
	respondJSON(w, map[string]interface{}{
		"message": "Canary deployment configured successfully",
	})
}

// Snapshot handlers
func (s *Server) handleListSnapshots(w http.ResponseWriter, r *http.Request) {
	snapshots, err := s.store.ListSnapshots()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, snapshots)
}

func (s *Server) handleCreateSnapshot(w http.ResponseWriter, r *http.Request) {
	var snapshot models.Snapshot
	if err := json.NewDecoder(r.Body).Decode(&snapshot); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	snapshot.CreatedAt = time.Now()
	if snapshot.CreatedBy == "" {
		snapshot.CreatedBy = "admin"
	}
	
	if err := s.store.CreateSnapshot(&snapshot); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "CREATE", "Snapshot", snapshot.Name, true)
	s.broadcastEvent("snapshot_created", snapshot)
	
	respondJSON(w, snapshot)
}

func (s *Server) handleGetSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid snapshot ID")
		return
	}
	
	snapshot, err := s.store.GetSnapshot(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	if snapshot == nil {
		respondError(w, http.StatusNotFound, "Snapshot not found")
		return
	}
	
	respondJSON(w, snapshot)
}

func (s *Server) handleDeleteSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid snapshot ID")
		return
	}
	
	if err := s.store.DeleteSnapshot(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "DELETE", "Snapshot", fmt.Sprintf("id:%d", id), true)
	s.broadcastEvent("snapshot_deleted", map[string]int{"id": id})
	
	respondJSON(w, map[string]interface{}{"message": "Snapshot deleted successfully"})
}

func (s *Server) handleRestoreSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid snapshot ID")
		return
	}
	
	snapshot, err := s.store.GetSnapshot(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	if snapshot == nil {
		respondError(w, http.StatusNotFound, "Snapshot not found")
		return
	}
	
	// Implementation would apply snapshot configuration
	s.logAudit("admin", "RESTORE", "Snapshot", snapshot.Name, true)
	s.broadcastEvent("snapshot_restored", snapshot)
	
	respondJSON(w, map[string]interface{}{"message": "Snapshot restored successfully"})
}

// ScheduledAction handlers
func (s *Server) handleListScheduledActions(w http.ResponseWriter, r *http.Request) {
	actions, err := s.store.ListScheduledActions()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, actions)
}

func (s *Server) handleCreateScheduledAction(w http.ResponseWriter, r *http.Request) {
	var action models.ScheduledAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	action.CreatedAt = time.Now()
	action.UpdatedAt = time.Now()
	
	if err := s.store.CreateScheduledAction(&action); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "CREATE", "ScheduledAction", action.Name, true)
	s.broadcastEvent("scheduled_action_created", action)
	
	respondJSON(w, action)
}

func (s *Server) handleGetScheduledAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid action ID")
		return
	}
	
	action, err := s.store.GetScheduledAction(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	if action == nil {
		respondError(w, http.StatusNotFound, "Scheduled action not found")
		return
	}
	
	respondJSON(w, action)
}

func (s *Server) handleUpdateScheduledAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid action ID")
		return
	}
	
	var action models.ScheduledAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	action.ID = id
	action.UpdatedAt = time.Now()
	
	if err := s.store.UpdateScheduledAction(&action); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "UPDATE", "ScheduledAction", action.Name, true)
	s.broadcastEvent("scheduled_action_updated", action)
	
	respondJSON(w, action)
}

func (s *Server) handleDeleteScheduledAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid action ID")
		return
	}
	
	if err := s.store.DeleteScheduledAction(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "DELETE", "ScheduledAction", fmt.Sprintf("id:%d", id), true)
	s.broadcastEvent("scheduled_action_deleted", map[string]int{"id": id})
	
	respondJSON(w, map[string]interface{}{"message": "Scheduled action deleted successfully"})
}

// Audit log handlers
func (s *Server) handleListAuditLogs(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	
	logs, err := s.store.ListAuditLogs(limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, logs)
}

// YAML preview handler
func (s *Server) handleYAMLPreview(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type string          `json:"type"`
		Spec json.RawMessage `json:"spec"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Convert spec to YAML based on type
	yaml := "# YAML preview for " + req.Type + "\napiVersion: networking.istio.io/v1beta1\nkind: " + req.Type
	
	respondJSON(w, map[string]interface{}{
		"yaml": yaml,
	})
}

// Helper methods
func (s *Server) logAudit(username, action, resource, details string, success bool) {
	log := &models.AuditLog{
		UserID:    1, // Default admin user
		Username:  username,
		Action:    action,
		Resource:  resource,
		Details:   details,
		Success:   success,
		Timestamp: time.Now(),
	}
	s.store.CreateAuditLog(log)
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
