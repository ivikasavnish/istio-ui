package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	
	"github.com/ivikasavnish/istio-ui/backend/internal/istio"
	"github.com/ivikasavnish/istio-ui/backend/internal/kube"
	"github.com/ivikasavnish/istio-ui/backend/internal/models"
	"github.com/ivikasavnish/istio-ui/backend/internal/storage"
)

type Server struct {
	store       *storage.Store
	kubeClient  *kube.Client
	istioMgr    *istio.Manager
	wsClients   map[*websocket.Conn]bool
	wsUpgrader  websocket.Upgrader
	wsMutex     sync.RWMutex
}

func NewServer(store *storage.Store, kubeClient *kube.Client) *Server {
	var istioMgr *istio.Manager
	if kubeClient != nil {
		istioMgr, _ = istio.NewManager(kubeClient)
	}
	
	return &Server{
		store:      store,
		kubeClient: kubeClient,
		istioMgr:   istioMgr,
		wsClients:  make(map[*websocket.Conn]bool),
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (s *Server) RegisterRoutes(r *mux.Router) {
	// Service discovery
	r.HandleFunc("/services", s.handleListServices).Methods("GET")
	r.HandleFunc("/services/{namespace}", s.handleListServicesInNamespace).Methods("GET")
	r.HandleFunc("/namespaces", s.handleListNamespaces).Methods("GET")
	
	// VirtualService routes
	r.HandleFunc("/virtualservices", s.handleListVirtualServices).Methods("GET")
	r.HandleFunc("/virtualservices/{namespace}", s.handleListVirtualServicesInNamespace).Methods("GET")
	r.HandleFunc("/virtualservices/{namespace}/{name}", s.handleGetVirtualService).Methods("GET")
	r.HandleFunc("/virtualservices", s.handleCreateVirtualService).Methods("POST")
	r.HandleFunc("/virtualservices/{namespace}/{name}", s.handleUpdateVirtualService).Methods("PUT")
	r.HandleFunc("/virtualservices/{namespace}/{name}", s.handleDeleteVirtualService).Methods("DELETE")
	
	// DestinationRule routes
	r.HandleFunc("/destinationrules", s.handleListDestinationRules).Methods("GET")
	r.HandleFunc("/destinationrules/{namespace}", s.handleListDestinationRulesInNamespace).Methods("GET")
	r.HandleFunc("/destinationrules/{namespace}/{name}", s.handleGetDestinationRule).Methods("GET")
	r.HandleFunc("/destinationrules", s.handleCreateDestinationRule).Methods("POST")
	r.HandleFunc("/destinationrules/{namespace}/{name}", s.handleUpdateDestinationRule).Methods("PUT")
	r.HandleFunc("/destinationrules/{namespace}/{name}", s.handleDeleteDestinationRule).Methods("DELETE")
	
	// Gateway routes
	r.HandleFunc("/gateways", s.handleListGateways).Methods("GET")
	r.HandleFunc("/gateways/{namespace}", s.handleListGatewaysInNamespace).Methods("GET")
	r.HandleFunc("/gateways", s.handleCreateGateway).Methods("POST")
	r.HandleFunc("/gateways/{namespace}/{name}", s.handleDeleteGateway).Methods("DELETE")
	
	// PeerAuthentication routes
	r.HandleFunc("/peerauthentications", s.handleListPeerAuthentications).Methods("GET")
	r.HandleFunc("/peerauthentications/{namespace}", s.handleListPeerAuthenticationsInNamespace).Methods("GET")
	r.HandleFunc("/peerauthentications", s.handleCreatePeerAuthentication).Methods("POST")
	r.HandleFunc("/peerauthentications/{namespace}/{name}", s.handleDeletePeerAuthentication).Methods("DELETE")
	
	// AuthorizationPolicy routes
	r.HandleFunc("/authorizationpolicies", s.handleListAuthorizationPolicies).Methods("GET")
	r.HandleFunc("/authorizationpolicies/{namespace}", s.handleListAuthorizationPoliciesInNamespace).Methods("GET")
	r.HandleFunc("/authorizationpolicies", s.handleCreateAuthorizationPolicy).Methods("POST")
	r.HandleFunc("/authorizationpolicies/{namespace}/{name}", s.handleDeleteAuthorizationPolicy).Methods("DELETE")
	
	// Traffic management
	r.HandleFunc("/traffic/weights", s.handleUpdateTrafficWeights).Methods("POST")
	r.HandleFunc("/traffic/canary", s.handleCanaryDeployment).Methods("POST")
	
	// Snapshots
	r.HandleFunc("/snapshots", s.handleListSnapshots).Methods("GET")
	r.HandleFunc("/snapshots", s.handleCreateSnapshot).Methods("POST")
	r.HandleFunc("/snapshots/{id}", s.handleGetSnapshot).Methods("GET")
	r.HandleFunc("/snapshots/{id}", s.handleDeleteSnapshot).Methods("DELETE")
	r.HandleFunc("/snapshots/{id}/restore", s.handleRestoreSnapshot).Methods("POST")
	
	// Scheduled actions
	r.HandleFunc("/scheduled-actions", s.handleListScheduledActions).Methods("GET")
	r.HandleFunc("/scheduled-actions", s.handleCreateScheduledAction).Methods("POST")
	r.HandleFunc("/scheduled-actions/{id}", s.handleGetScheduledAction).Methods("GET")
	r.HandleFunc("/scheduled-actions/{id}", s.handleUpdateScheduledAction).Methods("PUT")
	r.HandleFunc("/scheduled-actions/{id}", s.handleDeleteScheduledAction).Methods("DELETE")
	
	// Audit logs
	r.HandleFunc("/audit-logs", s.handleListAuditLogs).Methods("GET")
	
	// YAML preview
	r.HandleFunc("/yaml/preview", s.handleYAMLPreview).Methods("POST")
}

// Service handlers
func (s *Server) handleListServices(w http.ResponseWriter, r *http.Request) {
	if s.kubeClient == nil {
		respondError(w, http.StatusServiceUnavailable, "Kubernetes client not available")
		return
	}
	
	services, err := s.store.ListServiceCache()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, services)
}

func (s *Server) handleListServicesInNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	
	if s.kubeClient == nil {
		respondError(w, http.StatusServiceUnavailable, "Kubernetes client not available")
		return
	}
	
	services, err := s.kubeClient.DiscoverServices(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, services)
}

func (s *Server) handleListNamespaces(w http.ResponseWriter, r *http.Request) {
	if s.kubeClient == nil {
		respondError(w, http.StatusServiceUnavailable, "Kubernetes client not available")
		return
	}
	
	namespaces, err := s.kubeClient.ListNamespaces(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, map[string]interface{}{"namespaces": namespaces})
}

// VirtualService handlers
func (s *Server) handleListVirtualServices(w http.ResponseWriter, r *http.Request) {
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}
	
	list, err := s.istioMgr.ListVirtualServices(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleListVirtualServicesInNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	list, err := s.istioMgr.ListVirtualServices(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleGetVirtualService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	vs, err := s.istioMgr.GetVirtualService(r.Context(), namespace, name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, vs)
}

func (s *Server) handleCreateVirtualService(w http.ResponseWriter, r *http.Request) {
	var spec models.VirtualServiceSpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.CreateVirtualService(r.Context(), &spec); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "CREATE", "VirtualService", fmt.Sprintf("%s/%s", spec.Namespace, spec.Name), true)
	s.broadcastEvent("virtualservice_created", spec)
	
	respondJSON(w, map[string]interface{}{"message": "VirtualService created successfully"})
}

func (s *Server) handleUpdateVirtualService(w http.ResponseWriter, r *http.Request) {
	var spec models.VirtualServiceSpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	vars := mux.Vars(r)
	spec.Namespace = vars["namespace"]
	spec.Name = vars["name"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.UpdateVirtualService(r.Context(), &spec); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "UPDATE", "VirtualService", fmt.Sprintf("%s/%s", spec.Namespace, spec.Name), true)
	s.broadcastEvent("virtualservice_updated", spec)
	
	respondJSON(w, map[string]interface{}{"message": "VirtualService updated successfully"})
}

func (s *Server) handleDeleteVirtualService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.DeleteVirtualService(r.Context(), namespace, name); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "DELETE", "VirtualService", fmt.Sprintf("%s/%s", namespace, name), true)
	s.broadcastEvent("virtualservice_deleted", map[string]string{"namespace": namespace, "name": name})
	
	respondJSON(w, map[string]interface{}{"message": "VirtualService deleted successfully"})
}

// DestinationRule handlers
func (s *Server) handleListDestinationRules(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	list, err := s.istioMgr.ListDestinationRules(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleListDestinationRulesInNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	list, err := s.istioMgr.ListDestinationRules(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleGetDestinationRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	dr, err := s.istioMgr.GetDestinationRule(r.Context(), namespace, name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, dr)
}

func (s *Server) handleCreateDestinationRule(w http.ResponseWriter, r *http.Request) {
	var spec models.DestinationRuleSpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.CreateDestinationRule(r.Context(), &spec); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "CREATE", "DestinationRule", fmt.Sprintf("%s/%s", spec.Namespace, spec.Name), true)
	s.broadcastEvent("destinationrule_created", spec)
	
	respondJSON(w, map[string]interface{}{"message": "DestinationRule created successfully"})
}

func (s *Server) handleUpdateDestinationRule(w http.ResponseWriter, r *http.Request) {
	var spec models.DestinationRuleSpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	vars := mux.Vars(r)
	spec.Namespace = vars["namespace"]
	spec.Name = vars["name"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.UpdateDestinationRule(r.Context(), &spec); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "UPDATE", "DestinationRule", fmt.Sprintf("%s/%s", spec.Namespace, spec.Name), true)
	s.broadcastEvent("destinationrule_updated", spec)
	
	respondJSON(w, map[string]interface{}{"message": "DestinationRule updated successfully"})
}

func (s *Server) handleDeleteDestinationRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.DeleteDestinationRule(r.Context(), namespace, name); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "DELETE", "DestinationRule", fmt.Sprintf("%s/%s", namespace, name), true)
	s.broadcastEvent("destinationrule_deleted", map[string]string{"namespace": namespace, "name": name})
	
	respondJSON(w, map[string]interface{}{"message": "DestinationRule deleted successfully"})
}

// Gateway handlers  
func (s *Server) handleListGateways(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	list, err := s.istioMgr.ListGateways(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleListGatewaysInNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	list, err := s.istioMgr.ListGateways(r.Context(), namespace)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondJSON(w, list.Items)
}

func (s *Server) handleCreateGateway(w http.ResponseWriter, r *http.Request) {
	var spec models.GatewaySpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.CreateGateway(r.Context(), &spec); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "CREATE", "Gateway", fmt.Sprintf("%s/%s", spec.Namespace, spec.Name), true)
	s.broadcastEvent("gateway_created", spec)
	
	respondJSON(w, map[string]interface{}{"message": "Gateway created successfully"})
}

func (s *Server) handleDeleteGateway(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]
	
	if s.istioMgr == nil {
		respondError(w, http.StatusServiceUnavailable, "Istio manager not available")
		return
	}
	
	if err := s.istioMgr.DeleteGateway(r.Context(), namespace, name); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	s.logAudit("admin", "DELETE", "Gateway", fmt.Sprintf("%s/%s", namespace, name), true)
	s.broadcastEvent("gateway_deleted", map[string]string{"namespace": namespace, "name": name})
	
	respondJSON(w, map[string]interface{}{"message": "Gateway deleted successfully"})
}

// Continue in next part...
