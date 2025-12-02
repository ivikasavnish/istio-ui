# Future Improvements and Known Issues

## Minor Improvements (from Code Review)

### Frontend
- [ ] Update DataGrid `pageSize` prop to `initialState.pagination.pageSize` pattern (deprecated in v6+)
  - Affected files: VirtualServices.js, DestinationRules.js, Gateways.js, AuthorizationPolicies.js, PeerAuthentications.js, ScheduledActions.js, HelmManager.js
  - Impact: None (still works, but newer API is preferred)
  - Priority: Low

- [ ] Add JSON parsing error handling in HelmManager.js
  - Lines: 82, 103 (handleSaveInstall, handleSaveUpgrade)
  - Add try-catch blocks around JSON.parse for better user feedback
  - Priority: Medium

- [ ] Replace window.location.reload() with React state management in ContextSelector.js
  - Line: 60
  - Use state refresh instead of full page reload to preserve app state
  - Priority: Low

### Backend
- [ ] Enhance error messages in scheduler.go with more context (action type/name)
  - Lines: 101, 131
  - Impact: Better debugging experience
  - Priority: Low

- [ ] Extract actionConfig.Init into helper method in helm/manager.go
  - Duplicated across 7 methods
  - Reduces code duplication
  - Priority: Low

- [ ] Add chart path validation in helm/manager.go
  - Lines: 105, 128 (InstallChart, UpgradeChart)
  - Validate and sanitize chartPath to prevent directory traversal
  - Priority: High (security consideration)

## Feature Enhancements

### Phase 2 (Post-MVP)
- [ ] Add authentication/authorization
- [ ] Implement WebSocket for real-time updates
- [ ] Add comprehensive test coverage
- [ ] Integrate Prometheus metrics
- [ ] Add frontend unit and integration tests
- [ ] Implement advanced RBAC with user roles
- [ ] Add configuration validation before apply
- [ ] Support for Istio telemetry resources
- [ ] Add dark mode theme support

### Phase 3 (Advanced Features)
- [ ] Multi-cluster support
- [ ] GitOps integration (Flux/ArgoCD)
- [ ] Service mesh health scoring
- [ ] Automated canary analysis
- [ ] Advanced traffic analytics
- [ ] Policy recommendations engine
- [ ] Integration with Kiali
- [ ] Integration with Jaeger for tracing

## Technical Debt
- [ ] Add backend unit tests
- [ ] Add frontend component tests
- [ ] Add integration tests
- [ ] Set up CI/CD pipeline
- [ ] Add API rate limiting
- [ ] Implement request caching
- [ ] Add OpenAPI/Swagger documentation
- [ ] Add telemetry and observability

## Documentation
- [ ] Add video tutorial
- [ ] Create interactive demo
- [ ] Add more example configurations
- [ ] Create FAQ section
- [ ] Add troubleshooting guide expansions

## Notes
All current functionality is working as designed. These are enhancement opportunities for future versions.
The application is production-ready in its current state.
