import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  CircularProgress,
  Alert,
} from '@mui/material';
import { DataGrid } from '@mui/x-data-grid';
import { peerAuthenticationApi } from '../services/api';

export default function PeerAuthentications() {
  const [peerAuths, setPeerAuths] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadPeerAuths();
  }, []);

  const loadPeerAuths = async () => {
    try {
      const response = await peerAuthenticationApi.list();
      setPeerAuths(response.data || []);
    } catch (error) {
      setError('Failed to load peer authentications: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', width: 200, valueGetter: (params) => params.row.metadata?.name },
    { field: 'namespace', headerName: 'Namespace', width: 150, valueGetter: (params) => params.row.metadata?.namespace },
    { field: 'mtls', headerName: 'mTLS Mode', width: 150, valueGetter: (params) => params.row.spec?.mtls?.mode || 'PERMISSIVE' },
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
        <Typography variant="h4">Peer Authentications</Typography>
      </Box>

      {error && (
        <Alert severity="error" onClose={() => setError('')} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Paper>
        <DataGrid
          rows={peerAuths}
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
