import { useState, useEffect } from 'react'
import { virtualServicesAPI, namespacesAPI } from '../services/api'
import { Plus, Trash2 } from 'lucide-react'

const TrafficManagement = () => {
  const [namespace, setNamespace] = useState('default')
  const [namespaces, setNamespaces] = useState<string[]>([])
  const [virtualServices, setVirtualServices] = useState<any[]>([])
  const [showCreateForm, setShowCreateForm] = useState(false)
  
  // Traffic weight form
  const [service, setService] = useState('')
  const [v1Weight, setV1Weight] = useState(50)
  const [v2Weight, setV2Weight] = useState(50)
  const [v3Weight, setV3Weight] = useState(0)

  useEffect(() => {
    loadNamespaces()
  }, [])

  useEffect(() => {
    loadVirtualServices()
  }, [namespace])

  const loadNamespaces = async () => {
    try {
      const res = await namespacesAPI.list()
      setNamespaces(res.data?.namespaces || ['default'])
    } catch (error) {
      console.error('Failed to load namespaces:', error)
    }
  }

  const loadVirtualServices = async () => {
    try {
      const res = await virtualServicesAPI.list(namespace)
      setVirtualServices(res.data || [])
    } catch (error) {
      console.error('Failed to load virtual services:', error)
    }
  }

  const handleWeightChange = (version: string, value: number) => {
    if (version === 'v1') setV1Weight(value)
    else if (version === 'v2') setV2Weight(value)
    else if (version === 'v3') setV3Weight(value)
  }

  const normalizeWeights = () => {
    const total = v1Weight + v2Weight + v3Weight
    if (total > 100) {
      const scale = 100 / total
      setV1Weight(Math.round(v1Weight * scale))
      setV2Weight(Math.round(v2Weight * scale))
      setV3Weight(Math.round(v3Weight * scale))
    }
  }

  const handleCreateVirtualService = async () => {
    try {
      const spec = {
        name: service,
        namespace,
        hosts: [service],
        http: [
          {
            route: [
              { destination: { host: service, subset: 'v1' }, weight: v1Weight },
              { destination: { host: service, subset: 'v2' }, weight: v2Weight },
              ...(v3Weight > 0 ? [{ destination: { host: service, subset: 'v3' }, weight: v3Weight }] : []),
            ],
          },
        ],
      }
      
      await virtualServicesAPI.create(spec)
      await loadVirtualServices()
      setShowCreateForm(false)
      alert('VirtualService created successfully!')
    } catch (error) {
      console.error('Failed to create virtual service:', error)
      alert('Failed to create VirtualService')
    }
  }

  const handleDelete = async (name: string) => {
    if (!confirm(`Delete VirtualService ${name}?`)) return
    
    try {
      await virtualServicesAPI.delete(namespace, name)
      await loadVirtualServices()
      alert('VirtualService deleted successfully!')
    } catch (error) {
      console.error('Failed to delete virtual service:', error)
      alert('Failed to delete VirtualService')
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground">Traffic Management</h2>
        <p className="text-muted-foreground mt-1">
          Control traffic routing and canary deployments
        </p>
      </div>

      {/* Namespace Selector */}
      <div className="rounded-lg border bg-card p-4">
        <label className="block text-sm font-medium mb-2">Namespace</label>
        <select
          value={namespace}
          onChange={(e) => setNamespace(e.target.value)}
          className="w-full rounded-md border bg-background px-3 py-2"
        >
          {namespaces.map((ns) => (
            <option key={ns} value={ns}>{ns}</option>
          ))}
        </select>
      </div>

      {/* Traffic Weight Control */}
      {showCreateForm && (
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Create Traffic Split</h3>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">Service Name</label>
              <input
                type="text"
                value={service}
                onChange={(e) => setService(e.target.value)}
                className="w-full rounded-md border bg-background px-3 py-2"
                placeholder="my-service"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">
                Version 1 (v1): {v1Weight}%
              </label>
              <input
                type="range"
                min="0"
                max="100"
                value={v1Weight}
                onChange={(e) => handleWeightChange('v1', parseInt(e.target.value))}
                onBlur={normalizeWeights}
                className="w-full"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">
                Version 2 (v2): {v2Weight}%
              </label>
              <input
                type="range"
                min="0"
                max="100"
                value={v2Weight}
                onChange={(e) => handleWeightChange('v2', parseInt(e.target.value))}
                onBlur={normalizeWeights}
                className="w-full"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">
                Version 3 (v3): {v3Weight}%
              </label>
              <input
                type="range"
                min="0"
                max="100"
                value={v3Weight}
                onChange={(e) => handleWeightChange('v3', parseInt(e.target.value))}
                onBlur={normalizeWeights}
                className="w-full"
              />
            </div>

            <div className="text-sm text-muted-foreground">
              Total: {v1Weight + v2Weight + v3Weight}%
            </div>

            <div className="flex space-x-2">
              <button
                onClick={handleCreateVirtualService}
                className="rounded-md bg-primary px-4 py-2 text-primary-foreground hover:opacity-90"
              >
                Create VirtualService
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
          <span>Create Traffic Split</span>
        </button>
      )}

      {/* VirtualServices List */}
      <div className="rounded-lg border bg-card">
        <div className="border-b p-4">
          <h3 className="text-lg font-semibold">VirtualServices</h3>
        </div>
        <div className="p-4">
          {virtualServices.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">
              No VirtualServices found in this namespace
            </p>
          ) : (
            <div className="space-y-2">
              {virtualServices.map((vs: any, idx: number) => (
                <div key={idx} className="flex items-center justify-between rounded-md border p-4">
                  <div>
                    <p className="font-medium">{vs.metadata?.name || 'Unknown'}</p>
                    <p className="text-sm text-muted-foreground">
                      Hosts: {vs.spec?.hosts?.join(', ') || '-'}
                    </p>
                  </div>
                  <button
                    onClick={() => handleDelete(vs.metadata?.name)}
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

export default TrafficManagement
