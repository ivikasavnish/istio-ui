import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Paper,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  CircularProgress,
  IconButton,
  Alert,
} from '@mui/material';
import { DataGrid } from '@mui/x-data-grid';
import DeleteIcon from '@mui/icons-material/Delete';
import AddIcon from '@mui/icons-material/Add';
import CodeIcon from '@mui/icons-material/Code';
import { gatewayApi } from '../services/api';
import { objectToYaml } from '../utils/helpers';

export default function Gateways() {
  const [gateways, setGateways] = useState([]);
  const [loading, setLoading] = useState(true);
  const [openDialog, setOpenDialog] = useState(false);
  const [openYamlDialog, setOpenYamlDialog] = useState(false);
  const [yamlContent, setYamlContent] = useState('');
  const [error, setError] = useState('');
  const [formData, setFormData] = useState({
    name: '',
    namespace: 'default',
    selector: 'istio: ingressgateway',
    host: '*',
    port: '80',
  });

  useEffect(() => {
    loadGateways();
  }, []);

  const loadGateways = async () => {
    try {
      const response = await gatewayApi.list();
      setGateways(response.data || []);
    } catch (error) {
      setError('Failed to load gateways: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setFormData({
      name: '',
      namespace: 'default',
      selector: 'istio: ingressgateway',
      host: '*',
      port: '80',
    });
    setOpenDialog(true);
  };

  const handleViewYaml = (gw) => {
    setYamlContent(objectToYaml(gw));
    setOpenYamlDialog(true);
  };

  const handleSave = async () => {
    try {
      const gwObject = {
        apiVersion: 'networking.istio.io/v1beta1',
        kind: 'Gateway',
        metadata: {
          name: formData.name,
          namespace: formData.namespace,
        },
        spec: {
          selector: {
            istio: 'ingressgateway',
          },
          servers: [
            {
              port: {
                number: parseInt(formData.port),
                name: 'http',
                protocol: 'HTTP',
              },
              hosts: [formData.host],
            },
          ],
        },
      };

      await gatewayApi.create(gwObject);
      setOpenDialog(false);
      loadGateways();
    } catch (error) {
      setError('Failed to save gateway: ' + error.message);
    }
  };

  const handleDelete = async (namespace, name) => {
    if (window.confirm(`Are you sure you want to delete Gateway ${name}?`)) {
      try {
        await gatewayApi.delete(namespace, name);
        loadGateways();
      } catch (error) {
        setError('Failed to delete gateway: ' + error.message);
      }
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', width: 200, valueGetter: (params) => params.row.metadata?.name },
    { field: 'namespace', headerName: 'Namespace', width: 150, valueGetter: (params) => params.row.metadata?.namespace },
    { field: 'selector', headerName: 'Selector', width: 200, valueGetter: (params) => JSON.stringify(params.row.spec?.selector || {}) },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 200,
      renderCell: (params) => (
        <>
          <IconButton size="small" onClick={() => handleViewYaml(params.row)}>
            <CodeIcon />
          </IconButton>
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
        <Typography variant="h4">Gateways</Typography>
        <Button variant="contained" startIcon={<AddIcon />} onClick={handleCreate}>
          Create Gateway
        </Button>
      </Box>

      {error && (
        <Alert severity="error" onClose={() => setError('')} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Paper>
        <DataGrid
          rows={gateways}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 25, 50]}
          autoHeight
          getRowId={(row) => `${row.metadata?.namespace}/${row.metadata?.name}`}
          disableSelectionOnClick
        />
      </Paper>

      {/* Create Dialog */}
      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Create Gateway</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            margin="normal"
            label="Name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          />
          <TextField
            fullWidth
            margin="normal"
            label="Namespace"
            value={formData.namespace}
            onChange={(e) => setFormData({ ...formData, namespace: e.target.value })}
          />
          <TextField
            fullWidth
            margin="normal"
            label="Host"
            value={formData.host}
            onChange={(e) => setFormData({ ...formData, host: e.target.value })}
            helperText="e.g., * or myapp.example.com"
          />
          <TextField
            fullWidth
            margin="normal"
            label="Port"
            type="number"
            value={formData.port}
            onChange={(e) => setFormData({ ...formData, port: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button onClick={handleSave} variant="contained">
            Save
          </Button>
        </DialogActions>
      </Dialog>

      {/* YAML Dialog */}
      <Dialog open={openYamlDialog} onClose={() => setOpenYamlDialog(false)} maxWidth="md" fullWidth>
        <DialogTitle>YAML Preview</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            multiline
            rows={20}
            value={yamlContent}
            InputProps={{
              readOnly: true,
              style: { fontFamily: 'monospace', fontSize: '12px' },
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenYamlDialog(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
