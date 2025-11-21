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
import EditIcon from '@mui/icons-material/Edit';
import AddIcon from '@mui/icons-material/Add';
import CodeIcon from '@mui/icons-material/Code';
import { virtualServiceApi, kubernetesApi } from '../services/api';
import { objectToYaml, yamlToObject } from '../utils/helpers';

export default function VirtualServices() {
  const [virtualServices, setVirtualServices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [openDialog, setOpenDialog] = useState(false);
  const [openYamlDialog, setOpenYamlDialog] = useState(false);
  const [selectedVS, setSelectedVS] = useState(null);
  const [yamlContent, setYamlContent] = useState('');
  const [error, setError] = useState('');
  const [formData, setFormData] = useState({
    name: '',
    namespace: 'default',
    host: '',
    destinationHost: '',
    destinationSubset: '',
  });

  useEffect(() => {
    loadVirtualServices();
  }, []);

  const loadVirtualServices = async () => {
    try {
      const response = await virtualServiceApi.list();
      setVirtualServices(response.data || []);
    } catch (error) {
      setError('Failed to load virtual services: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setSelectedVS(null);
    setFormData({
      name: '',
      namespace: 'default',
      host: '',
      destinationHost: '',
      destinationSubset: '',
    });
    setOpenDialog(true);
  };

  const handleEdit = (vs) => {
    setSelectedVS(vs);
    setFormData({
      name: vs.metadata?.name || '',
      namespace: vs.metadata?.namespace || 'default',
      host: vs.spec?.hosts?.[0] || '',
      destinationHost: vs.spec?.http?.[0]?.route?.[0]?.destination?.host || '',
      destinationSubset: vs.spec?.http?.[0]?.route?.[0]?.destination?.subset || '',
    });
    setOpenDialog(true);
  };

  const handleViewYaml = (vs) => {
    setYamlContent(objectToYaml(vs));
    setOpenYamlDialog(true);
  };

  const handleSave = async () => {
    try {
      const vsObject = {
        apiVersion: 'networking.istio.io/v1beta1',
        kind: 'VirtualService',
        metadata: {
          name: formData.name,
          namespace: formData.namespace,
        },
        spec: {
          hosts: [formData.host],
          http: [
            {
              route: [
                {
                  destination: {
                    host: formData.destinationHost,
                    subset: formData.destinationSubset || undefined,
                  },
                },
              ],
            },
          ],
        },
      };

      if (selectedVS) {
        await virtualServiceApi.update(formData.namespace, formData.name, vsObject);
      } else {
        await virtualServiceApi.create(vsObject);
      }

      setOpenDialog(false);
      loadVirtualServices();
    } catch (error) {
      setError('Failed to save virtual service: ' + error.message);
    }
  };

  const handleDelete = async (namespace, name) => {
    if (window.confirm(`Are you sure you want to delete VirtualService ${name}?`)) {
      try {
        await virtualServiceApi.delete(namespace, name);
        loadVirtualServices();
      } catch (error) {
        setError('Failed to delete virtual service: ' + error.message);
      }
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', width: 200, valueGetter: (params) => params.row.metadata?.name },
    { field: 'namespace', headerName: 'Namespace', width: 150, valueGetter: (params) => params.row.metadata?.namespace },
    { field: 'hosts', headerName: 'Hosts', width: 250, valueGetter: (params) => params.row.spec?.hosts?.join(', ') },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 200,
      renderCell: (params) => (
        <>
          <IconButton size="small" onClick={() => handleEdit(params.row)}>
            <EditIcon />
          </IconButton>
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
        <Typography variant="h4">Virtual Services</Typography>
        <Button variant="contained" startIcon={<AddIcon />} onClick={handleCreate}>
          Create Virtual Service
        </Button>
      </Box>

      {error && (
        <Alert severity="error" onClose={() => setError('')} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Paper>
        <DataGrid
          rows={virtualServices}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 25, 50]}
          autoHeight
          getRowId={(row) => `${row.metadata?.namespace}/${row.metadata?.name}`}
          disableSelectionOnClick
        />
      </Paper>

      {/* Create/Edit Dialog */}
      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>{selectedVS ? 'Edit' : 'Create'} Virtual Service</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            margin="normal"
            label="Name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            disabled={!!selectedVS}
          />
          <TextField
            fullWidth
            margin="normal"
            label="Namespace"
            value={formData.namespace}
            onChange={(e) => setFormData({ ...formData, namespace: e.target.value })}
            disabled={!!selectedVS}
          />
          <TextField
            fullWidth
            margin="normal"
            label="Host"
            value={formData.host}
            onChange={(e) => setFormData({ ...formData, host: e.target.value })}
            helperText="e.g., myapp.example.com"
          />
          <TextField
            fullWidth
            margin="normal"
            label="Destination Host"
            value={formData.destinationHost}
            onChange={(e) => setFormData({ ...formData, destinationHost: e.target.value })}
            helperText="e.g., myapp.default.svc.cluster.local"
          />
          <TextField
            fullWidth
            margin="normal"
            label="Destination Subset (Optional)"
            value={formData.destinationSubset}
            onChange={(e) => setFormData({ ...formData, destinationSubset: e.target.value })}
            helperText="e.g., v1"
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
