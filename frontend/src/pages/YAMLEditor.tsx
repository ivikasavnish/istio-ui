import { useState } from 'react'
import { FileCode, Eye, CheckCircle } from 'lucide-react'

const YAMLEditor = () => {
  const [resourceType, setResourceType] = useState('VirtualService')
  const [yaml, setYaml] = useState(`apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: example
  namespace: default
spec:
  hosts:
  - example.com
  http:
  - route:
    - destination:
        host: example
        subset: v1
      weight: 100
`)

  const handleApply = () => {
    alert('YAML would be applied to cluster')
    console.log('Applying YAML:', yaml)
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground">YAML Editor</h2>
        <p className="text-muted-foreground mt-1">
          Create and edit Istio resources
        </p>
      </div>

      <div className="rounded-lg border bg-card p-4">
        <label className="block text-sm font-medium mb-2">Resource Type</label>
        <select
          value={resourceType}
          onChange={(e) => setResourceType(e.target.value)}
          className="w-full max-w-xs rounded-md border bg-background px-3 py-2"
        >
          <option>VirtualService</option>
          <option>DestinationRule</option>
          <option>Gateway</option>
          <option>PeerAuthentication</option>
          <option>AuthorizationPolicy</option>
        </select>
      </div>

      <div className="rounded-lg border bg-card">
        <div className="border-b p-4 flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <FileCode className="h-5 w-5" />
            <h3 className="text-lg font-semibold">YAML Configuration</h3>
          </div>
          <button className="flex items-center space-x-2 rounded-md border px-3 py-1 hover:bg-accent">
            <Eye className="h-4 w-4" />
            <span className="text-sm">Preview</span>
          </button>
        </div>
        <div className="p-4">
          <textarea
            value={yaml}
            onChange={(e) => setYaml(e.target.value)}
            className="w-full h-96 rounded-md border bg-background p-4 font-mono text-sm"
            spellCheck={false}
          />
        </div>
        <div className="border-t p-4 flex space-x-2">
          <button
            onClick={handleApply}
            className="flex items-center space-x-2 rounded-md bg-primary px-4 py-2 text-primary-foreground hover:opacity-90"
          >
            <CheckCircle className="h-4 w-4" />
            <span>Apply to Cluster</span>
          </button>
        </div>
      </div>
    </div>
  )
}

export default YAMLEditor
