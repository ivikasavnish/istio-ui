import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Paper,
  Typography,
  CircularProgress,
  IconButton,
  Alert,
} from '@mui/material';
import { DataGrid } from '@mui/x-data-grid';
import DeleteIcon from '@mui/icons-material/Delete';
import CodeIcon from '@mui/icons-material/Code';
import { authorizationPolicyApi } from '../services/api';

export default function AuthorizationPolicies() {
  const [policies, setPolicies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadPolicies();
  }, []);

  const loadPolicies = async () => {
    try {
      const response = await authorizationPolicyApi.list();
      setPolicies(response.data || []);
    } catch (error) {
      setError('Failed to load authorization policies: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (namespace, name) => {
    if (window.confirm(`Are you sure you want to delete AuthorizationPolicy ${name}?`)) {
      try {
        await authorizationPolicyApi.delete(namespace, name);
        loadPolicies();
      } catch (error) {
        setError('Failed to delete authorization policy: ' + error.message);
      }
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', width: 200, valueGetter: (params) => params.row.metadata?.name },
    { field: 'namespace', headerName: 'Namespace', width: 150, valueGetter: (params) => params.row.metadata?.namespace },
    { field: 'action', headerName: 'Action', width: 150, valueGetter: (params) => params.row.spec?.action || 'ALLOW' },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 200,
      renderCell: (params) => (
        <>
          <IconButton
            size="small"
            onClick={() => handleDelete(params.row.metadata?.namespace, params.row.metadata?.name)}
          >
            <DeleteIcon />
          </IconButton>
        </>
      ),
    },
  ];

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Authorization Policies</Typography>
      </Box>

      {error && (
        <Alert severity="error" onClose={() => setError('')} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Paper>
        <DataGrid
          rows={policies}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 25, 50]}
          autoHeight
          getRowId={(row) => `${row.metadata?.namespace}/${row.metadata?.name}`}
          disableSelectionOnClick
        />
      </Paper>
    </Box>
  );
}
