import React, { useState, useEffect } from 'react';
import {
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Alert,
  Chip,
} from '@mui/material';
import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

export default function ContextSelector({ onContextChange }) {
  const [contexts, setContexts] = useState([]);
  const [currentContext, setCurrentContext] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadContexts();
  }, []);

  const loadContexts = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/contexts`);
      setContexts(response.data || []);
      
      // Find the current context
      const current = response.data.find(ctx => ctx.is_current);
      if (current) {
        setCurrentContext(current.name);
      }
    } catch (error) {
      console.error('Error loading contexts:', error);
      setError('Failed to load contexts');
    } finally {
      setLoading(false);
    }
  };

  const handleContextChange = async (event) => {
    const newContext = event.target.value;
    setLoading(true);
    setError('');

    try {
      await axios.post(`${API_BASE_URL}/contexts/switch`, {
        context: newContext,
      });
      
      setCurrentContext(newContext);
      if (onContextChange) {
        onContextChange(newContext);
      }
      
      // Reload the page to refresh all data with new context
      window.location.reload();
    } catch (error) {
      console.error('Error switching context:', error);
      setError('Failed to switch context: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading && contexts.length === 0) {
    return (
      <Box display="flex" alignItems="center" gap={1}>
        <CircularProgress size={20} />
      </Box>
    );
  }

  return (
    <Box sx={{ minWidth: 200 }}>
      {error && (
        <Alert severity="error" onClose={() => setError('')} sx={{ mb: 1 }}>
          {error}
        </Alert>
      )}
      <FormControl fullWidth size="small">
        <InputLabel>K8s Context</InputLabel>
        <Select
          value={currentContext}
          label="K8s Context"
          onChange={handleContextChange}
          disabled={loading}
        >
          {contexts.map((ctx) => (
            <MenuItem key={ctx.name} value={ctx.name}>
              <Box display="flex" alignItems="center" gap={1}>
                {ctx.name}
                {ctx.is_current && (
                  <Chip label="Current" size="small" color="primary" />
                )}
              </Box>
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );
}
