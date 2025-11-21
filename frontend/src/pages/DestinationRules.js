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
import { destinationRuleApi } from '../services/api';
import { objectToYaml } from '../utils/helpers';

export default function DestinationRules() {
  const [destinationRules, setDestinationRules] = useState([]);
  const [loading, setLoading] = useState(true);
  const [openDialog, setOpenDialog] = useState(false);
  const [openYamlDialog, setOpenYamlDialog] = useState(false);
  const [selectedDR, setSelectedDR] = useState(null);
  const [yamlContent, setYamlContent] = useState('');
  const [error, setError] = useState('');
  const [formData, setFormData] = useState({
    name: '',
    namespace: 'default',
    host: '',
    subsetName: '',
    subsetVersion: '',
  });

  useEffect(() => {
    loadDestinationRules();
  }, []);

  const loadDestinationRules = async () => {
    try {
      const response = await destinationRuleApi.list();
      setDestinationRules(response.data || []);
    } catch (error) {
      setError('Failed to load destination rules: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setSelectedDR(null);
    setFormData({
      name: '',
      namespace: 'default',
      host: '',
      subsetName: '',
      subsetVersion: '',
    });
    setOpenDialog(true);
  };

  const handleEdit = (dr) => {
    setSelectedDR(dr);
    setFormData({
      name: dr.metadata?.name || '',
      namespace: dr.metadata?.namespace || 'default',
      host: dr.spec?.host || '',
      subsetName: dr.spec?.subsets?.[0]?.name || '',
      subsetVersion: dr.spec?.subsets?.[0]?.labels?.version || '',
    });
    setOpenDialog(true);
  };

  const handleViewYaml = (dr) => {
    setYamlContent(objectToYaml(dr));
    setOpenYamlDialog(true);
  };

  const handleSave = async () => {
    try {
      const drObject = {
        apiVersion: 'networking.istio.io/v1beta1',
        kind: 'DestinationRule',
        metadata: {
          name: formData.name,
          namespace: formData.namespace,
        },
        spec: {
          host: formData.host,
          subsets: [
            {
              name: formData.subsetName,
              labels: {
                version: formData.subsetVersion,
              },
            },
          ],
        },
      };

      if (selectedDR) {
        await destinationRuleApi.update(formData.namespace, formData.name, drObject);
      } else {
        await destinationRuleApi.create(drObject);
      }

      setOpenDialog(false);
      loadDestinationRules();
    } catch (error) {
      setError('Failed to save destination rule: ' + error.message);
    }
  };

  const handleDelete = async (namespace, name) => {
    if (window.confirm(`Are you sure you want to delete DestinationRule ${name}?`)) {
      try {
        await destinationRuleApi.delete(namespace, name);
        loadDestinationRules();
      } catch (error) {
        setError('Failed to delete destination rule: ' + error.message);
      }
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', width: 200, valueGetter: (params) => params.row.metadata?.name },
    { field: 'namespace', headerName: 'Namespace', width: 150, valueGetter: (params) => params.row.metadata?.namespace },
    { field: 'host', headerName: 'Host', width: 250, valueGetter: (params) => params.row.spec?.host },
    { field: 'subsets', headerName: 'Subsets', width: 150, valueGetter: (params) => params.row.spec?.subsets?.length || 0 },
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
        <Typography variant="h4">Destination Rules</Typography>
        <Button variant="contained" startIcon={<AddIcon />} onClick={handleCreate}>
          Create Destination Rule
        </Button>
      </Box>

      {error && (
        <Alert severity="error" onClose={() => setError('')} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Paper>
        <DataGrid
          rows={destinationRules}
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
        <DialogTitle>{selectedDR ? 'Edit' : 'Create'} Destination Rule</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            margin="normal"
            label="Name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            disabled={!!selectedDR}
          />
          <TextField
            fullWidth
            margin="normal"
            label="Namespace"
            value={formData.namespace}
            onChange={(e) => setFormData({ ...formData, namespace: e.target.value })}
            disabled={!!selectedDR}
          />
          <TextField
            fullWidth
            margin="normal"
            label="Host"
            value={formData.host}
            onChange={(e) => setFormData({ ...formData, host: e.target.value })}
            helperText="e.g., myapp.default.svc.cluster.local"
          />
          <TextField
            fullWidth
            margin="normal"
            label="Subset Name"
            value={formData.subsetName}
            onChange={(e) => setFormData({ ...formData, subsetName: e.target.value })}
            helperText="e.g., v1"
          />
          <TextField
            fullWidth
            margin="normal"
            label="Subset Version"
            value={formData.subsetVersion}
            onChange={(e) => setFormData({ ...formData, subsetVersion: e.target.value })}
            helperText="Label value for version"
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
