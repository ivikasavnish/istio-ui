# Example: Progressive Canary Deployment with Scheduled Actions

This example demonstrates how to use MeshControl Center to perform a progressive canary deployment
with automated traffic shifting using scheduled actions.

## Scenario

You have a service called `product-service` with two versions:
- v1 (current stable version)
- v2 (new version to be deployed)

Goal: Gradually shift traffic from v1 to v2 over 24 hours, with automatic rollback if needed.

## Step 1: Create DestinationRule

Define subsets for both versions:

```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: product-service
  namespace: default
spec:
  host: product-service.default.svc.cluster.local
  trafficPolicy:
    loadBalancer:
      simple: ROUND_ROBIN
  subsets:
    - name: v1
      labels:
        version: v1
    - name: v2
      labels:
        version: v2
```

**Via UI:**
1. Navigate to "Destination Rules"
2. Click "Create Destination Rule"
3. Fill in:
   - Name: `product-service`
   - Namespace: `default`
   - Host: `product-service.default.svc.cluster.local`
   - Subset Name: `v1`, Version: `v1`
4. Save and repeat for v2

## Step 2: Initial VirtualService (100% to v1)

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: product-service
  namespace: default
spec:
  hosts:
    - product-service.default.svc.cluster.local
  http:
    - route:
        - destination:
            host: product-service.default.svc.cluster.local
            subset: v1
          weight: 100
        - destination:
            host: product-service.default.svc.cluster.local
            subset: v2
          weight: 0
```

**Via UI:**
1. Navigate to "Virtual Services"
2. Click "Create Virtual Service"
3. Configure initial routing to v1

## Step 3: Create Scheduled Actions

### Action 1: Shift 25% traffic to v2 (after 6 hours)

```json
{
  "name": "Canary Phase 1 - 25%",
  "description": "Shift 25% of traffic to v2",
  "schedule": "0 6 * * *",
  "action_type": "update_virtualservice",
  "namespace": "default",
  "resource_name": "product-service",
  "payload": {
    "spec": {
      "hosts": ["product-service.default.svc.cluster.local"],
      "http": [{
        "route": [
          {
            "destination": {
              "host": "product-service.default.svc.cluster.local",
              "subset": "v1"
            },
            "weight": 75
          },
          {
            "destination": {
              "host": "product-service.default.svc.cluster.local",
              "subset": "v2"
            },
            "weight": 25
          }
        ]
      }]
    }
  }
}
```

### Action 2: Shift 50% traffic to v2 (after 12 hours)

```json
{
  "name": "Canary Phase 2 - 50%",
  "description": "Shift 50% of traffic to v2",
  "schedule": "0 12 * * *",
  "action_type": "update_virtualservice",
  "namespace": "default",
  "resource_name": "product-service",
  "payload": {
    "spec": {
      "hosts": ["product-service.default.svc.cluster.local"],
      "http": [{
        "route": [
          {
            "destination": {
              "host": "product-service.default.svc.cluster.local",
              "subset": "v1"
            },
            "weight": 50
          },
          {
            "destination": {
              "host": "product-service.default.svc.cluster.local",
              "subset": "v2"
            },
            "weight": 50
          }
        ]
      }]
    }
  }
}
```

### Action 3: Shift 100% traffic to v2 (after 24 hours)

```json
{
  "name": "Canary Phase 3 - 100%",
  "description": "Complete migration to v2",
  "schedule": "0 0 * * *",
  "action_type": "update_virtualservice",
  "namespace": "default",
  "resource_name": "product-service",
  "payload": {
    "spec": {
      "hosts": ["product-service.default.svc.cluster.local"],
      "http": [{
        "route": [
          {
            "destination": {
              "host": "product-service.default.svc.cluster.local",
              "subset": "v2"
            },
            "weight": 100
          }
        ]
      }]
    }
  }
}
```

**Via UI:**
1. Navigate to "Scheduled Actions"
2. Click "Create Schedule"
3. Fill in the details from each action above
4. Adjust cron schedules based on your deployment window

## Step 4: Monitor the Deployment

1. Check the "Topology" page to visualize traffic flow
2. Monitor application metrics (error rates, latency)
3. Check logs for any issues
4. Use the "Scheduled Actions" page to see upcoming executions

## Rollback Strategy

If issues are detected with v2, create an immediate rollback schedule:

```json
{
  "name": "Emergency Rollback",
  "description": "Rollback all traffic to v1",
  "schedule": "* * * * *",
  "action_type": "update_virtualservice",
  "namespace": "default",
  "resource_name": "product-service",
  "payload": {
    "spec": {
      "hosts": ["product-service.default.svc.cluster.local"],
      "http": [{
        "route": [
          {
            "destination": {
              "host": "product-service.default.svc.cluster.local",
              "subset": "v1"
            },
            "weight": 100
          }
        ]
      }]
    }
  }
}
```

Or manually update the VirtualService through the UI to shift traffic back to v1.

## Advanced: Header-based Routing

Add header-based routing for beta testers before the canary:

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: product-service
  namespace: default
spec:
  hosts:
    - product-service.default.svc.cluster.local
  http:
    - match:
        - headers:
            user-group:
              exact: "beta-testers"
      route:
        - destination:
            host: product-service.default.svc.cluster.local
            subset: v2
          weight: 100
    - route:
        - destination:
            host: product-service.default.svc.cluster.local
            subset: v1
          weight: 100
```

This allows beta testers to access v2 while all other traffic goes to v1.

## Cron Schedule Reference

```
# Every 6 hours
0 */6 * * *

# Daily at midnight
0 0 * * *

# Every hour during business hours (9 AM - 5 PM)
0 9-17 * * *

# Weekdays at 9 AM
0 9 * * 1-5

# Every 15 minutes
*/15 * * * *
```

## Best Practices

1. **Start Small**: Begin with a small percentage (10-25%) to v2
2. **Monitor Closely**: Watch metrics closely during the initial phase
3. **Have Rollback Ready**: Always have a rollback plan
4. **Test Before Scheduling**: Manually test traffic shifting before automating
5. **Use Appropriate Intervals**: Don't rush; give time to detect issues
6. **Document Changes**: Keep track of what changes were made and when
7. **Communicate**: Inform your team about scheduled deployments

## Cleanup

After successful migration:
1. Delete the scheduled actions (they're one-time or recurring)
2. Update the DestinationRule to remove v1 subset (optional)
3. Scale down v1 deployments
4. Archive old configurations
