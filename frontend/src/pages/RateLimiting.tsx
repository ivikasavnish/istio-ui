import { Gauge } from 'lucide-react'

const RateLimiting = () => {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground">Rate Limiting</h2>
        <p className="text-muted-foreground mt-1">
          Configure rate limits for your services
        </p>
      </div>

      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center space-x-2 mb-4">
          <Gauge className="h-5 w-5 text-primary" />
          <h3 className="text-lg font-semibold">Rate Limit Configuration</h3>
        </div>
        <p className="text-muted-foreground">
          Rate limiting configuration coming soon. This feature will allow you to:
        </p>
        <ul className="list-disc list-inside mt-2 space-y-1 text-muted-foreground">
          <li>Set request rate limits per service</li>
          <li>Configure quotas and burst limits</li>
          <li>Apply local and global rate limiting</li>
          <li>Monitor rate limit metrics</li>
        </ul>
      </div>
    </div>
  )
}

export default RateLimiting
