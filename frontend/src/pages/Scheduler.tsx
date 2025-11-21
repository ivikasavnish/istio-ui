import { useState, useEffect } from 'react'
import { scheduledActionsAPI } from '../services/api'
import { Plus, Trash2, Clock } from 'lucide-react'

const Scheduler = () => {
  const [actions, setActions] = useState<any[]>([])
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [name, setName] = useState('')
  const [actionType, setActionType] = useState('traffic_shift')
  const [cronExpr, setCronExpr] = useState('0 0 * * *')

  useEffect(() => {
    loadActions()
  }, [])

  const loadActions = async () => {
    try {
      const res = await scheduledActionsAPI.list()
      setActions(res.data || [])
    } catch (error) {
      console.error('Failed to load scheduled actions:', error)
    }
  }

  const handleCreate = async () => {
    try {
      const action = {
        name,
        action_type: actionType,
        cron_expr: cronExpr,
        config: JSON.stringify({}),
        enabled: true,
        next_run: new Date().toISOString(),
      }
      
      await scheduledActionsAPI.create(action)
      await loadActions()
      setShowCreateForm(false)
      setName('')
      alert('Scheduled action created successfully!')
    } catch (error) {
      console.error('Failed to create scheduled action:', error)
      alert('Failed to create scheduled action')
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this scheduled action?')) return
    
    try {
      await scheduledActionsAPI.delete(id)
      await loadActions()
      alert('Scheduled action deleted successfully!')
    } catch (error) {
      console.error('Failed to delete scheduled action:', error)
      alert('Failed to delete scheduled action')
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground">Scheduler</h2>
        <p className="text-muted-foreground mt-1">
          Automate Istio configuration changes
        </p>
      </div>

      {showCreateForm && (
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Create Scheduled Action</h3>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">Name</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full rounded-md border bg-background px-3 py-2"
                placeholder="Daily traffic shift"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">Action Type</label>
              <select
                value={actionType}
                onChange={(e) => setActionType(e.target.value)}
                className="w-full rounded-md border bg-background px-3 py-2"
              >
                <option value="traffic_shift">Traffic Shift</option>
                <option value="apply_mtls">Apply mTLS</option>
                <option value="remove_fault">Remove Fault Injection</option>
                <option value="snapshot_capture">Capture Snapshot</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">
                Cron Expression
              </label>
              <input
                type="text"
                value={cronExpr}
                onChange={(e) => setCronExpr(e.target.value)}
                className="w-full rounded-md border bg-background px-3 py-2"
                placeholder="0 0 * * *"
              />
              <p className="text-xs text-muted-foreground mt-1">
                Format: minute hour day month weekday (e.g., "0 0 * * *" = daily at midnight)
              </p>
            </div>

            <div className="flex space-x-2">
              <button
                onClick={handleCreate}
                disabled={!name || !cronExpr}
                className="rounded-md bg-primary px-4 py-2 text-primary-foreground hover:opacity-90 disabled:opacity-50"
              >
                Create Action
              </button>
              <button
                onClick={() => setShowCreateForm(false)}
                className="rounded-md border px-4 py-2 hover:bg-accent"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {!showCreateForm && (
        <button
          onClick={() => setShowCreateForm(true)}
          className="flex items-center space-x-2 rounded-md bg-primary px-4 py-2 text-primary-foreground hover:opacity-90"
        >
          <Plus className="h-4 w-4" />
          <span>Create Scheduled Action</span>
        </button>
      )}

      <div className="rounded-lg border bg-card">
        <div className="border-b p-4 flex items-center space-x-2">
          <Clock className="h-5 w-5" />
          <h3 className="text-lg font-semibold">Scheduled Actions</h3>
        </div>
        <div className="p-4">
          {actions.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">
              No scheduled actions yet
            </p>
          ) : (
            <div className="space-y-2">
              {actions.map((action: any) => (
                <div key={action.id} className="flex items-center justify-between rounded-md border p-4">
                  <div>
                    <p className="font-medium">{action.name}</p>
                    <p className="text-sm text-muted-foreground">
                      Type: {action.action_type} | Cron: {action.cron_expr}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      Next run: {new Date(action.next_run).toLocaleString()}
                    </p>
                  </div>
                  <button
                    onClick={() => handleDelete(action.id)}
                    className="rounded-md p-2 text-destructive hover:bg-destructive/10"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Scheduler
