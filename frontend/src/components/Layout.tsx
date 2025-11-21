import { ReactNode } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { Moon, Sun, Home, GitBranch, AlertCircle, Shield, Gauge, Network, FileCode, Clock, History } from 'lucide-react'

interface LayoutProps {
  children: ReactNode
  darkMode: boolean
  setDarkMode: (value: boolean) => void
}

const Layout = ({ children, darkMode, setDarkMode }: LayoutProps) => {
  const location = useLocation()

  const navItems = [
    { path: '/', label: 'Dashboard', icon: Home },
    { path: '/traffic', label: 'Traffic Management', icon: GitBranch },
    { path: '/fault-injection', label: 'Fault Injection', icon: AlertCircle },
    { path: '/security', label: 'Security', icon: Shield },
    { path: '/rate-limiting', label: 'Rate Limiting', icon: Gauge },
    { path: '/topology', label: 'Topology', icon: Network },
    { path: '/yaml-editor', label: 'YAML Editor', icon: FileCode },
    { path: '/scheduler', label: 'Scheduler', icon: Clock },
    { path: '/history', label: 'History', icon: History },
  ]

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="sticky top-0 z-50 border-b bg-card">
        <div className="flex h-16 items-center px-4">
          <div className="flex items-center space-x-2">
            <Network className="h-8 w-8 text-primary" />
            <h1 className="text-2xl font-bold text-foreground">MeshControl Center</h1>
          </div>
          <div className="ml-auto flex items-center space-x-4">
            <button
              onClick={() => setDarkMode(!darkMode)}
              className="rounded-md p-2 hover:bg-accent"
              aria-label="Toggle dark mode"
            >
              {darkMode ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
            </button>
          </div>
        </div>
      </header>

      <div className="flex">
        {/* Sidebar */}
        <aside className="sticky top-16 h-[calc(100vh-4rem)] w-64 border-r bg-card overflow-y-auto">
          <nav className="space-y-1 p-4">
            {navItems.map((item) => {
              const Icon = item.icon
              const isActive = location.pathname === item.path
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  className={`flex items-center space-x-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
                    isActive
                      ? 'bg-primary text-primary-foreground'
                      : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                  }`}
                >
                  <Icon className="h-5 w-5" />
                  <span>{item.label}</span>
                </Link>
              )
            })}
          </nav>
        </aside>

        {/* Main Content */}
        <main className="flex-1 p-6 overflow-y-auto">
          {children}
        </main>
      </div>
    </div>
  )
}

export default Layout
