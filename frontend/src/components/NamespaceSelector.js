import React, { useState, useEffect } from 'react';
import {
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
} from '@mui/material';
import { kubernetesApi } from '../services/api';

export default function NamespaceSelector({ value, onChange, showAllOption = true }) {
  const [namespaces, setNamespaces] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadNamespaces();
  }, []);

  const loadNamespaces = async () => {
    try {
      const response = await kubernetesApi.listNamespaces();
      setNamespaces(response.data || []);
    } catch (error) {
      console.error('Error loading namespaces:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" alignItems="center" gap={1}>
        <CircularProgress size={20} />
      </Box>
    );
  }

  return (
    <Box sx={{ minWidth: 150 }}>
      <FormControl fullWidth size="small">
        <InputLabel>Namespace</InputLabel>
        <Select
          value={value || ''}
          label="Namespace"
          onChange={(e) => onChange(e.target.value)}
        >
          {showAllOption && (
            <MenuItem value="">
              <em>All Namespaces</em>
            </MenuItem>
          )}
          {namespaces.map((ns) => (
            <MenuItem key={ns} value={ns}>
              {ns}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );
}
