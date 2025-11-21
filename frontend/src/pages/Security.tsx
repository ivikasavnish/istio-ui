import { useState, useEffect } from 'react'
import { peerAuthAPI, authzPolicyAPI, namespacesAPI } from '../services/api'
import { Shield, Lock } from 'lucide-react'

const Security = () => {
  const [namespace, setNamespace] = useState('default')
  const [namespaces, setNamespaces] = useState<string[]>([])
  const [mtlsMode, setMtlsMode] = useState<'STRICT' | 'PERMISSIVE' | 'DISABLE'>('STRICT')
  const [peerAuths, setPeerAuths] = useState<any[]>([])

  useEffect(() => {
    loadNamespaces()
  }, [])

  useEffect(() => {
    loadPeerAuthentications()
  }, [namespace])

  const loadNamespaces = async () => {
    try {
      const res = await namespacesAPI.list()
      setNamespaces(res.data?.namespaces || ['default'])
    } catch (error) {
      console.error('Failed to load namespaces:', error)
    }
  }

  const loadPeerAuthentications = async () => {
    try {
      const res = await peerAuthAPI.list(namespace)
      setPeerAuths(res.data || [])
    } catch (error) {
      console.error('Failed to load peer authentications:', error)
    }
  }

  const handleApplyMTLS = async () => {
    try {
      const spec = {
        name: 'default',
        namespace,
        mtls_mode: mtlsMode,
      }
      
      await peerAuthAPI.create(spec)
      await loadPeerAuthentications()
      alert('mTLS policy applied successfully!')
    } catch (error) {
      console.error('Failed to apply mTLS:', error)
      alert('Failed to apply mTLS policy')
    }
  }

  const handleDeletePeerAuth = async (name: string) => {
    if (!confirm(`Delete PeerAuthentication ${name}?`)) return
    
    try {
      await peerAuthAPI.delete(namespace, name)
      await loadPeerAuthentications()
      alert('PeerAuthentication deleted successfully!')
    } catch (error) {
      console.error('Failed to delete peer authentication:', error)
      alert('Failed to delete PeerAuthentication')
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground">Security</h2>
        <p className="text-muted-foreground mt-1">
          Configure mTLS, authentication, and authorization
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

      {/* mTLS Configuration */}
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center space-x-2 mb-4">
          <Lock className="h-5 w-5 text-primary" />
          <h3 className="text-lg font-semibold">mTLS Configuration</h3>
        </div>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">mTLS Mode</label>
            <div className="space-y-2">
              <label className="flex items-center space-x-2">
                <input
                  type="radio"
                  checked={mtlsMode === 'STRICT'}
                  onChange={() => setMtlsMode('STRICT')}
                />
                <div>
                  <span className="font-medium">Strict</span>
                  <p className="text-sm text-muted-foreground">
                    All connections must use mTLS
                  </p>
                </div>
              </label>

              <label className="flex items-center space-x-2">
                <input
                  type="radio"
                  checked={mtlsMode === 'PERMISSIVE'}
                  onChange={() => setMtlsMode('PERMISSIVE')}
                />
                <div>
                  <span className="font-medium">Permissive</span>
                  <p className="text-sm text-muted-foreground">
                    Accept both mTLS and plain text
                  </p>
                </div>
              </label>

              <label className="flex items-center space-x-2">
                <input
                  type="radio"
                  checked={mtlsMode === 'DISABLE'}
                  onChange={() => setMtlsMode('DISABLE')}
                />
                <div>
                  <span className="font-medium">Disable</span>
                  <p className="text-sm text-muted-foreground">
                    Disable mTLS
                  </p>
                </div>
              </label>
            </div>
          </div>

          <button
            onClick={handleApplyMTLS}
            className="rounded-md bg-primary px-4 py-2 text-primary-foreground hover:opacity-90"
          >
            Apply mTLS Policy
          </button>
        </div>
      </div>

      {/* PeerAuthentications List */}
      <div className="rounded-lg border bg-card">
        <div className="border-b p-4 flex items-center space-x-2">
          <Shield className="h-5 w-5" />
          <h3 className="text-lg font-semibold">PeerAuthentications</h3>
        </div>
        <div className="p-4">
          {peerAuths.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">
              No PeerAuthentications found in this namespace
            </p>
          ) : (
            <div className="space-y-2">
              {peerAuths.map((pa: any, idx: number) => (
                <div key={idx} className="flex items-center justify-between rounded-md border p-4">
                  <div>
                    <p className="font-medium">{pa.metadata?.name || 'Unknown'}</p>
                    <p className="text-sm text-muted-foreground">
                      Mode: {pa.spec?.mtls?.mode || '-'}
                    </p>
                  </div>
                  <button
                    onClick={() => handleDeletePeerAuth(pa.metadata?.name)}
                    className="rounded-md px-3 py-1 text-sm text-destructive hover:bg-destructive/10"
                  >
                    Delete
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

export default Security
