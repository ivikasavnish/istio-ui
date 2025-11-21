import { useState } from 'react'
import { AlertCircle } from 'lucide-react'

const FaultInjection = () => {
  const [service, setService] = useState('')
  const [namespace, setNamespace] = useState('default')
  const [faultType, setFaultType] = useState<'delay' | 'abort'>('delay')
  
  // Delay settings
  const [delayMs, setDelayMs] = useState(5000)
  const [delayPercentage, setDelayPercentage] = useState(100)
  
  // Abort settings
  const [httpStatus, setHttpStatus] = useState(503)
  const [abortPercentage, setAbortPercentage] = useState(100)

  const handleApplyFault = async () => {
    try {
      const fault: any = {}
      
      if (faultType === 'delay') {
        fault.delay = {
          fixedDelay: `${delayMs}ms`,
          percentage: { value: delayPercentage },
        }
      } else {
        fault.abort = {
          httpStatus,
          percentage: { value: abortPercentage },
        }
      }

      alert(`Fault injection configured for ${service} in ${namespace}`)
      console.log('Fault configuration:', fault)
    } catch (error) {
      console.error('Failed to apply fault injection:', error)
      alert('Failed to apply fault injection')
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground">Fault Injection</h2>
        <p className="text-muted-foreground mt-1">
          Test resilience with delay and abort faults
        </p>
      </div>

      <div className="rounded-lg border bg-card p-6">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">Namespace</label>
            <input
              type="text"
              value={namespace}
              onChange={(e) => setNamespace(e.target.value)}
              className="w-full rounded-md border bg-background px-3 py-2"
            />
          </div>

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
            <label className="block text-sm font-medium mb-2">Fault Type</label>
            <div className="flex space-x-4">
              <label className="flex items-center space-x-2">
                <input
                  type="radio"
                  checked={faultType === 'delay'}
                  onChange={() => setFaultType('delay')}
                />
                <span>Delay</span>
              </label>
              <label className="flex items-center space-x-2">
                <input
                  type="radio"
                  checked={faultType === 'abort'}
                  onChange={() => setFaultType('abort')}
                />
                <span>Abort</span>
              </label>
            </div>
          </div>

          {faultType === 'delay' && (
            <>
              <div>
                <label className="block text-sm font-medium mb-2">
                  Delay (ms): {delayMs}
                </label>
                <input
                  type="range"
                  min="100"
                  max="10000"
                  step="100"
                  value={delayMs}
                  onChange={(e) => setDelayMs(parseInt(e.target.value))}
                  className="w-full"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">
                  Percentage: {delayPercentage}%
                </label>
                <input
                  type="range"
                  min="0"
                  max="100"
                  value={delayPercentage}
                  onChange={(e) => setDelayPercentage(parseInt(e.target.value))}
                  className="w-full"
                />
              </div>
            </>
          )}

          {faultType === 'abort' && (
            <>
              <div>
                <label className="block text-sm font-medium mb-2">
                  HTTP Status Code
                </label>
                <select
                  value={httpStatus}
                  onChange={(e) => setHttpStatus(parseInt(e.target.value))}
                  className="w-full rounded-md border bg-background px-3 py-2"
                >
                  <option value={400}>400 Bad Request</option>
                  <option value={401}>401 Unauthorized</option>
                  <option value={403}>403 Forbidden</option>
                  <option value={404}>404 Not Found</option>
                  <option value={500}>500 Internal Server Error</option>
                  <option value={503}>503 Service Unavailable</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">
                  Percentage: {abortPercentage}%
                </label>
                <input
                  type="range"
                  min="0"
                  max="100"
                  value={abortPercentage}
                  onChange={(e) => setAbortPercentage(parseInt(e.target.value))}
                  className="w-full"
                />
              </div>
            </>
          )}

          <div className="flex items-start space-x-2 rounded-md border border-yellow-500 bg-yellow-500/10 p-4">
            <AlertCircle className="h-5 w-5 text-yellow-500 flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-sm font-medium">Warning</p>
              <p className="text-sm text-muted-foreground">
                Fault injection will affect traffic to this service. Use carefully in production environments.
              </p>
            </div>
          </div>

          <button
            onClick={handleApplyFault}
            disabled={!service}
            className="rounded-md bg-primary px-4 py-2 text-primary-foreground hover:opacity-90 disabled:opacity-50"
          >
            Apply Fault Injection
          </button>
        </div>
      </div>
    </div>
  )
}

export default FaultInjection
