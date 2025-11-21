import { useState, useEffect } from 'react'
import { snapshotsAPI, auditLogsAPI } from '../services/api'
import { History as HistoryIcon, Camera, RotateCcw } from 'lucide-react'

const History = () => {
  const [snapshots, setSnapshots] = useState<any[]>([])
  const [auditLogs, setAuditLogs] = useState<any[]>([])
  const [activeTab, setActiveTab] = useState<'snapshots' | 'audit'>('snapshots')

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [snapshotsRes, logsRes] = await Promise.all([
        snapshotsAPI.list(),
        auditLogsAPI.list(50),
      ])
      setSnapshots(snapshotsRes.data || [])
      setAuditLogs(logsRes.data || [])
    } catch (error) {
      console.error('Failed to load history data:', error)
    }
  }

  const handleRestore = async (id: number) => {
    if (!confirm('Restore this snapshot? This will apply the saved configuration.')) return
    
    try {
      await snapshotsAPI.restore(id)
      alert('Snapshot restored successfully!')
    } catch (error) {
      console.error('Failed to restore snapshot:', error)
      alert('Failed to restore snapshot')
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground">History</h2>
        <p className="text-muted-foreground mt-1">
          View snapshots and audit logs
        </p>
      </div>

      <div className="rounded-lg border bg-card">
        <div className="border-b flex">
          <button
            onClick={() => setActiveTab('snapshots')}
            className={`px-6 py-3 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'snapshots'
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <div className="flex items-center space-x-2">
              <Camera className="h-4 w-4" />
              <span>Snapshots</span>
            </div>
          </button>
          <button
            onClick={() => setActiveTab('audit')}
            className={`px-6 py-3 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'audit'
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <div className="flex items-center space-x-2">
              <HistoryIcon className="h-4 w-4" />
              <span>Audit Logs</span>
            </div>
          </button>
        </div>

        <div className="p-4">
          {activeTab === 'snapshots' && (
            <div className="space-y-2">
              {snapshots.length === 0 ? (
                <p className="text-muted-foreground text-center py-8">
                  No snapshots yet
                </p>
              ) : (
                snapshots.map((snapshot: any) => (
                  <div key={snapshot.id} className="flex items-center justify-between rounded-md border p-4">
                    <div>
                      <p className="font-medium">{snapshot.name}</p>
                      <p className="text-sm text-muted-foreground">{snapshot.description}</p>
                      <p className="text-xs text-muted-foreground mt-1">
                        Created {new Date(snapshot.created_at).toLocaleString()} by {snapshot.created_by}
                      </p>
                    </div>
                    <button
                      onClick={() => handleRestore(snapshot.id)}
                      className="flex items-center space-x-2 rounded-md border px-3 py-1 hover:bg-accent"
                    >
                      <RotateCcw className="h-4 w-4" />
                      <span className="text-sm">Restore</span>
                    </button>
                  </div>
                ))
              )}
            </div>
          )}

          {activeTab === 'audit' && (
            <div className="overflow-x-auto">
              {auditLogs.length === 0 ? (
                <p className="text-muted-foreground text-center py-8">
                  No audit logs yet
                </p>
              ) : (
                <table className="w-full">
                  <thead className="border-b">
                    <tr>
                      <th className="text-left p-2 text-sm font-medium text-muted-foreground">Time</th>
                      <th className="text-left p-2 text-sm font-medium text-muted-foreground">User</th>
                      <th className="text-left p-2 text-sm font-medium text-muted-foreground">Action</th>
                      <th className="text-left p-2 text-sm font-medium text-muted-foreground">Resource</th>
                      <th className="text-left p-2 text-sm font-medium text-muted-foreground">Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {auditLogs.map((log: any) => (
                      <tr key={log.id} className="border-b last:border-0 hover:bg-accent">
                        <td className="p-2 text-sm">{new Date(log.timestamp).toLocaleString()}</td>
                        <td className="p-2 text-sm">{log.username}</td>
                        <td className="p-2 text-sm">{log.action}</td>
                        <td className="p-2 text-sm">{log.resource}</td>
                        <td className="p-2 text-sm">
                          <span className={`inline-flex items-center rounded-full px-2 py-1 text-xs ${
                            log.success ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                          }`}>
                            {log.success ? 'Success' : 'Failed'}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default History
