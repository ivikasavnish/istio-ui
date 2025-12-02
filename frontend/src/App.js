import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import VirtualServices from './pages/VirtualServices';
import DestinationRules from './pages/DestinationRules';
import Gateways from './pages/Gateways';
import AuthorizationPolicies from './pages/AuthorizationPolicies';
import PeerAuthentications from './pages/PeerAuthentications';
import Topology from './pages/Topology';
import ScheduledActions from './pages/ScheduledActions';
import HelmManager from './pages/HelmManager';

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Layout>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/virtualservices" element={<VirtualServices />} />
            <Route path="/destinationrules" element={<DestinationRules />} />
            <Route path="/gateways" element={<Gateways />} />
            <Route path="/authorizationpolicies" element={<AuthorizationPolicies />} />
            <Route path="/peerauthentications" element={<PeerAuthentications />} />
            <Route path="/topology" element={<Topology />} />
            <Route path="/schedules" element={<ScheduledActions />} />
            <Route path="/helm" element={<HelmManager />} />
          </Routes>
        </Layout>
      </Router>
    </ThemeProvider>
  );
}

export default App;
