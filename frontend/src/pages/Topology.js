import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  CircularProgress,
  Alert,
} from '@mui/material';
import { topologyApi } from '../services/api';

export default function Topology() {
  const [topology, setTopology] = useState({ nodes: [], edges: [] });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadTopology();
  }, []);

  const loadTopology = async () => {
    try {
      const response = await topologyApi.get();
      setTopology(response.data);
    } catch (error) {
      setError('Failed to load topology: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Service Mesh Topology
      </Typography>

      {error && (
        <Alert severity="error" onClose={() => setError('')} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Paper sx={{ p: 3, minHeight: '500px' }}>
        <Typography variant="h6" gutterBottom>
          Nodes: {topology.nodes.length}
        </Typography>
        <Typography variant="body2" color="text.secondary" paragraph>
          Services: {topology.nodes.filter(n => n.type === 'service').length} |
          Virtual Services: {topology.nodes.filter(n => n.type === 'virtualservice').length} |
          Gateways: {topology.nodes.filter(n => n.type === 'gateway').length} |
          Destination Rules: {topology.nodes.filter(n => n.type === 'destinationrule').length}
        </Typography>

        <Box sx={{ mt: 3 }}>
          <Typography variant="subtitle2" gutterBottom>
            Visualization coming soon...
          </Typography>
          <Typography variant="body2" color="text.secondary">
            This page will display an interactive graph of your service mesh topology,
            showing relationships between services, virtual services, gateways, and destination rules.
          </Typography>
        </Box>

        {/* In a real implementation, you would use react-flow-renderer or similar library here */}
        <Box sx={{ mt: 3 }}>
          {topology.nodes.slice(0, 10).map((node, idx) => (
            <Box key={idx} sx={{ p: 1, borderBottom: '1px solid #eee' }}>
              <Typography variant="body2">
                <strong>{node.name}</strong> ({node.type}) - {node.namespace}
              </Typography>
            </Box>
          ))}
        </Box>
      </Paper>
    </Box>
  );
}
