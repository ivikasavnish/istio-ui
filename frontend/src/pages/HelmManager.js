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
  Chip,
} from '@mui/material';
import { DataGrid } from '@mui/x-data-grid';
import DeleteIcon from '@mui/icons-material/Delete';
import HistoryIcon from '@mui/icons-material/History';
import AddIcon from '@mui/icons-material/Add';
import UpgradeIcon from '@mui/icons-material/Upgrade';
import axios from 'axios';
import NamespaceSelector from '../components/NamespaceSelector';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

export default function HelmManager() {
  const [releases, setReleases] = useState([]);
  const [loading, setLoading] = useState(true);
  const [openInstallDialog, setOpenInstallDialog] = useState(false);
  const [openUpgradeDialog, setOpenUpgradeDialog] = useState(false);
  const [selectedRelease, setSelectedRelease] = useState(null);
  const [error, setError] = useState('');
  const [namespace, setNamespace] = useState('default');
  const [formData, setFormData] = useState({
    name: '',
    namespace: 'default',
    chartPath: '',
    values: '',
  });

  useEffect(() => {
    loadReleases();
  }, [namespace]);

  const loadReleases = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/helm/releases?namespace=${namespace}`);
      setReleases(response.data || []);
    } catch (error) {
      setError('Failed to load Helm releases: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleInstall = () => {
    setFormData({
      name: '',
      namespace: namespace || 'default',
      chartPath: '',
      values: '',
    });
    setOpenInstallDialog(true);
  };

  const handleUpgrade = (release) => {
    setSelectedRelease(release);
    setFormData({
      name: release.name,
      namespace: release.namespace,
      chartPath: '',
      values: '',
    });
    setOpenUpgradeDialog(true);
  };

  const handleSaveInstall = async () => {
    try {
      let values = {};
      if (formData.values) {
        values = JSON.parse(formData.values);
      }

      await axios.post(`${API_BASE_URL}/helm/releases`, {
        name: formData.name,
        namespace: formData.namespace,
        chart_path: formData.chartPath,
        values: values,
      });

      setOpenInstallDialog(false);
      loadReleases();
    } catch (error) {
      setError('Failed to install chart: ' + error.message);
    }
  };

  const handleSaveUpgrade = async () => {
    try {
      let values = {};
      if (formData.values) {
        values = JSON.parse(formData.values);
      }

      await axios.put(`${API_BASE_URL}/helm/releases/${selectedRelease.name}`, {
        namespace: formData.namespace,
        chart_path: formData.chartPath,
        values: values,
      });

      setOpenUpgradeDialog(false);
      loadReleases();
    } catch (error) {
      setError('Failed to upgrade release: ' + error.message);
    }
  };

  const handleUninstall = async (name, namespace) => {
    if (window.confirm(`Are you sure you want to uninstall Helm release ${name}?`)) {
      try {
        await axios.delete(`${API_BASE_URL}/helm/releases/${name}?namespace=${namespace}`);
        loadReleases();
      } catch (error) {
        setError('Failed to uninstall release: ' + error.message);
      }
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', width: 200 },
    { field: 'namespace', headerName: 'Namespace', width: 150 },
    { field: 'chart', headerName: 'Chart', width: 200 },
    { field: 'version', headerName: 'Version', width: 100 },
    {
      field: 'status',
      headerName: 'Status',
      width: 150,
      renderCell: (params) => (
        <Chip
          label={params.value}
          color={params.value === 'deployed' ? 'success' : 'default'}
          size="small"
        />
      ),
    },
    { field: 'app_version', headerName: 'App Version', width: 150 },
    { field: 'updated', headerName: 'Updated', width: 200 },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 200,
      renderCell: (params) => (
        <>
          <IconButton size="small" onClick={() => handleUpgrade(params.row)}>
            <UpgradeIcon />
          </IconButton>
          <IconButton
            size="small"
            onClick={() => handleUninstall(params.row.name, params.row.namespace)}
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
        <Typography variant="h4">Helm Releases</Typography>
        <Box display="flex" gap={2} alignItems="center">
          <NamespaceSelector
            value={namespace}
            onChange={setNamespace}
            showAllOption={false}
          />
          <Button variant="contained" startIcon={<AddIcon />} onClick={handleInstall}>
            Install Chart
          </Button>
        </Box>
      </Box>

      {error && (
        <Alert severity="error" onClose={() => setError('')} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Paper>
        <DataGrid
          rows={releases}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 25, 50]}
          autoHeight
          getRowId={(row) => `${row.namespace}/${row.name}`}
          disableSelectionOnClick
        />
      </Paper>

      {/* Install Dialog */}
      <Dialog open={openInstallDialog} onClose={() => setOpenInstallDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Install Helm Chart</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            margin="normal"
            label="Release Name"
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
            label="Chart Path or URL"
            value={formData.chartPath}
            onChange={(e) => setFormData({ ...formData, chartPath: e.target.value })}
            helperText="e.g., stable/nginx-ingress or /path/to/chart"
          />
          <TextField
            fullWidth
            margin="normal"
            label="Values (JSON)"
            multiline
            rows={6}
            value={formData.values}
            onChange={(e) => setFormData({ ...formData, values: e.target.value })}
            helperText='Optional: {"key": "value"}'
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenInstallDialog(false)}>Cancel</Button>
          <Button onClick={handleSaveInstall} variant="contained">
            Install
          </Button>
        </DialogActions>
      </Dialog>

      {/* Upgrade Dialog */}
      <Dialog open={openUpgradeDialog} onClose={() => setOpenUpgradeDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Upgrade Helm Release</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            margin="normal"
            label="Release Name"
            value={formData.name}
            disabled
          />
          <TextField
            fullWidth
            margin="normal"
            label="Namespace"
            value={formData.namespace}
            disabled
          />
          <TextField
            fullWidth
            margin="normal"
            label="Chart Path or URL"
            value={formData.chartPath}
            onChange={(e) => setFormData({ ...formData, chartPath: e.target.value })}
            helperText="e.g., stable/nginx-ingress or /path/to/chart"
          />
          <TextField
            fullWidth
            margin="normal"
            label="Values (JSON)"
            multiline
            rows={6}
            value={formData.values}
            onChange={(e) => setFormData({ ...formData, values: e.target.value })}
            helperText='Optional: {"key": "value"}'
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenUpgradeDialog(false)}>Cancel</Button>
          <Button onClick={handleSaveUpgrade} variant="contained">
            Upgrade
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
