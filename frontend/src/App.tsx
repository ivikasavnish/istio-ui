import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { useState, useEffect } from 'react'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import TrafficManagement from './pages/TrafficManagement'
import FaultInjection from './pages/FaultInjection'
import Security from './pages/Security'
import RateLimiting from './pages/RateLimiting'
import TopologyGraph from './pages/TopologyGraph'
import YAMLEditor from './pages/YAMLEditor'
import Scheduler from './pages/Scheduler'
import History from './pages/History'

function App() {
  const [darkMode, setDarkMode] = useState(false)

  useEffect(() => {
    if (darkMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [darkMode])

  return (
    <Router>
      <Layout darkMode={darkMode} setDarkMode={setDarkMode}>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/traffic" element={<TrafficManagement />} />
          <Route path="/fault-injection" element={<FaultInjection />} />
          <Route path="/security" element={<Security />} />
          <Route path="/rate-limiting" element={<RateLimiting />} />
          <Route path="/topology" element={<TopologyGraph />} />
          <Route path="/yaml-editor" element={<YAMLEditor />} />
          <Route path="/scheduler" element={<Scheduler />} />
          <Route path="/history" element={<History />} />
        </Routes>
      </Layout>
    </Router>
  )
}

export default App
