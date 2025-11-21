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
  MenuItem,
} from '@mui/material';
import { DataGrid } from '@mui/x-data-grid';
import DeleteIcon from '@mui/icons-material/Delete';
import AddIcon from '@mui/icons-material/Add';
import { scheduleApi } from '../services/api';
import { formatDate } from '../utils/helpers';

export default function ScheduledActions() {
  const [schedules, setSchedules] = useState([]);
  const [loading, setLoading] = useState(true);
  const [openDialog, setOpenDialog] = useState(false);
  const [error, setError] = useState('');
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    schedule: '',
    action_type: 'update_virtualservice',
    namespace: 'default',
    resource_name: '',
  });

  useEffect(() => {
    loadSchedules();
  }, []);

  const loadSchedules = async () => {
    try {
      const response = await scheduleApi.list();
      setSchedules(response.data || []);
    } catch (error) {
      setError('Failed to load schedules: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setFormData({
      name: '',
      description: '',
      schedule: '0 0 * * *',
      action_type: 'update_virtualservice',
      namespace: 'default',
      resource_name: '',
    });
    setOpenDialog(true);
  };

  const handleSave = async () => {
    try {
      await scheduleApi.create({
        ...formData,
        payload: {},
      });
      setOpenDialog(false);
      loadSchedules();
    } catch (error) {
      setError('Failed to save schedule: ' + error.message);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this scheduled action?')) {
      try {
        await scheduleApi.delete(id);
        loadSchedules();
      } catch (error) {
        setError('Failed to delete schedule: ' + error.message);
      }
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', width: 200 },
    { field: 'description', headerName: 'Description', width: 250 },
    { field: 'schedule', headerName: 'Schedule', width: 150 },
    { field: 'action_type', headerName: 'Action Type', width: 200 },
    { field: 'next_run', headerName: 'Next Run', width: 200, valueGetter: (params) => formatDate(params.row.next_run) },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 100,
      renderCell: (params) => (
        <IconButton size="small" onClick={() => handleDelete(params.row.id)}>
          <DeleteIcon />
        </IconButton>
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
        <Typography variant="h4">Scheduled Actions</Typography>
        <Button variant="contained" startIcon={<AddIcon />} onClick={handleCreate}>
          Create Schedule
        </Button>
      </Box>

      {error && (
        <Alert severity="error" onClose={() => setError('')} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Paper>
        <DataGrid
          rows={schedules}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 25, 50]}
          autoHeight
          getRowId={(row) => row.id}
          disableSelectionOnClick
        />
      </Paper>

      {/* Create Dialog */}
      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Create Scheduled Action</DialogTitle>
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
            label="Description"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          />
          <TextField
            fullWidth
            margin="normal"
            label="Cron Schedule"
            value={formData.schedule}
            onChange={(e) => setFormData({ ...formData, schedule: e.target.value })}
            helperText="e.g., '0 0 * * *' for daily at midnight"
          />
          <TextField
            fullWidth
            select
            margin="normal"
            label="Action Type"
            value={formData.action_type}
            onChange={(e) => setFormData({ ...formData, action_type: e.target.value })}
          >
            <MenuItem value="update_virtualservice">Update VirtualService</MenuItem>
            <MenuItem value="update_destinationrule">Update DestinationRule</MenuItem>
            <MenuItem value="delete_virtualservice">Delete VirtualService</MenuItem>
            <MenuItem value="delete_destinationrule">Delete DestinationRule</MenuItem>
          </TextField>
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
            label="Resource Name"
            value={formData.resource_name}
            onChange={(e) => setFormData({ ...formData, resource_name: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button onClick={handleSave} variant="contained">
            Save
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
