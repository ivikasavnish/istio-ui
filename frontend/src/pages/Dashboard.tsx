import { useEffect, useState } from 'react'
import { servicesAPI, namespacesAPI } from '../services/api'
import { useWebSocket } from '../hooks/useWebSocket'
import { Activity, Database, Shield, GitBranch } from 'lucide-react'

const Dashboard = () => {
  const [services, setServices] = useState<any[]>([])
  const [namespaces, setNamespaces] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const { connected, messages } = useWebSocket()

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [servicesRes, namespacesRes] = await Promise.all([
        servicesAPI.list(),
        namespacesAPI.list(),
      ])
      setServices(servicesRes.data || [])
      setNamespaces(namespacesRes.data?.namespaces || [])
    } catch (error) {
      console.error('Failed to load dashboard data:', error)
    } finally {
      setLoading(false)
    }
  }

  const stats = [
    { label: 'Services', value: services.length, icon: Database, color: 'text-blue-500' },
    { label: 'Namespaces', value: namespaces.length, icon: GitBranch, color: 'text-green-500' },
    { label: 'WebSocket', value: connected ? 'Connected' : 'Disconnected', icon: Activity, color: connected ? 'text-green-500' : 'text-red-500' },
    { label: 'Events', value: messages.length, icon: Shield, color: 'text-purple-500' },
  ]

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-muted-foreground">Loading...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground">Dashboard</h2>
        <p className="text-muted-foreground mt-1">
          Monitor your Istio service mesh
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => {
          const Icon = stat.icon
          return (
            <div key={stat.label} className="rounded-lg border bg-card p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">{stat.label}</p>
                  <p className="text-2xl font-bold text-foreground">{stat.value}</p>
                </div>
                <Icon className={`h-8 w-8 ${stat.color}`} />
              </div>
            </div>
          )
        })}
      </div>

      {/* Services Table */}
      <div className="rounded-lg border bg-card">
        <div className="border-b p-4">
          <h3 className="text-lg font-semibold text-foreground">Services</h3>
        </div>
        <div className="p-4">
          {services.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">
              No services found. Connect to a Kubernetes cluster to see services.
            </p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="border-b">
                  <tr>
                    <th className="text-left p-2 text-sm font-medium text-muted-foreground">Namespace</th>
                    <th className="text-left p-2 text-sm font-medium text-muted-foreground">Name</th>
                    <th className="text-left p-2 text-sm font-medium text-muted-foreground">Version</th>
                    <th className="text-left p-2 text-sm font-medium text-muted-foreground">Last Seen</th>
                  </tr>
                </thead>
                <tbody>
                  {services.map((service, idx) => (
                    <tr key={idx} className="border-b last:border-0 hover:bg-accent">
                      <td className="p-2 text-sm">{service.namespace}</td>
                      <td className="p-2 text-sm font-medium">{service.name}</td>
                      <td className="p-2 text-sm">{service.version || '-'}</td>
                      <td className="p-2 text-sm">{new Date(service.last_seen).toLocaleString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>

      {/* Recent Events */}
      <div className="rounded-lg border bg-card">
        <div className="border-b p-4">
          <h3 className="text-lg font-semibold text-foreground">Recent Events</h3>
        </div>
        <div className="p-4">
          {messages.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">No events yet</p>
          ) : (
            <div className="space-y-2">
              {messages.slice(-10).reverse().map((msg, idx) => (
                <div key={idx} className="rounded-md border p-3">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">{msg.type}</span>
                    <span className="text-xs text-muted-foreground">Just now</span>
                  </div>
                  <pre className="text-xs text-muted-foreground mt-1 overflow-auto">
                    {JSON.stringify(msg.payload, null, 2)}
                  </pre>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Dashboard
