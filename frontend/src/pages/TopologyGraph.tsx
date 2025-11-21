import { useEffect, useRef, useState } from 'react'
import cytoscape from 'cytoscape'
// @ts-ignore - no type definitions available
import dagre from 'cytoscape-dagre'
import { servicesAPI, namespacesAPI } from '../services/api'

cytoscape.use(dagre)

const TopologyGraph = () => {
  const containerRef = useRef<HTMLDivElement>(null)
  const cyRef = useRef<cytoscape.Core | null>(null)
  const [namespace, setNamespace] = useState('default')
  const [namespaces, setNamespaces] = useState<string[]>([])

  useEffect(() => {
    loadNamespaces()
  }, [])

  useEffect(() => {
    if (containerRef.current) {
      initGraph()
    }
    
    return () => {
      if (cyRef.current) {
        cyRef.current.destroy()
      }
    }
  }, [namespace])

  const loadNamespaces = async () => {
    try {
      const res = await namespacesAPI.list()
      setNamespaces(res.data?.namespaces || ['default'])
    } catch (error) {
      console.error('Failed to load namespaces:', error)
    }
  }

  const initGraph = async () => {
    if (!containerRef.current) return

    try {
      const servicesRes = await servicesAPI.listInNamespace(namespace)
      const services = servicesRes.data || []

      const elements = services.map((svc: any) => ({
        data: {
          id: svc.name,
          label: svc.name,
        },
      }))

      // Add some example edges (in real implementation, derive from VirtualServices)
      const edges = []
      for (let i = 0; i < services.length - 1; i++) {
        if (Math.random() > 0.5) {
          edges.push({
            data: {
              id: `${services[i].name}-${services[i + 1].name}`,
              source: services[i].name,
              target: services[i + 1].name,
            },
          })
        }
      }

      cyRef.current = cytoscape({
        container: containerRef.current,
        elements: [...elements, ...edges],
        style: [
          {
            selector: 'node',
            style: {
              'background-color': '#3b82f6',
              'label': 'data(label)',
              'color': '#fff',
              'text-halign': 'center',
              'text-valign': 'center',
              'font-size': '12px',
              'width': '60px',
              'height': '60px',
            },
          },
          {
            selector: 'edge',
            style: {
              'width': 2,
              'line-color': '#94a3b8',
              'target-arrow-color': '#94a3b8',
              'target-arrow-shape': 'triangle',
              'curve-style': 'bezier',
            },
          },
        ],
        layout: {
          name: 'dagre',
          padding: 20,
        } as any,
      })

      // Add click handler
      cyRef.current.on('tap', 'node', (evt) => {
        const node = evt.target
        console.log('Clicked node:', node.data())
      })
    } catch (error) {
      console.error('Failed to load topology:', error)
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground">Topology Graph</h2>
        <p className="text-muted-foreground mt-1">
          Visualize your service mesh
        </p>
      </div>

      <div className="rounded-lg border bg-card p-4">
        <label className="block text-sm font-medium mb-2">Namespace</label>
        <select
          value={namespace}
          onChange={(e) => setNamespace(e.target.value)}
          className="w-full max-w-xs rounded-md border bg-background px-3 py-2"
        >
          {namespaces.map((ns) => (
            <option key={ns} value={ns}>{ns}</option>
          ))}
        </select>
      </div>

      <div className="rounded-lg border bg-card">
        <div
          ref={containerRef}
          className="w-full h-[600px] bg-background"
        />
      </div>
    </div>
  )
}

export default TopologyGraph
