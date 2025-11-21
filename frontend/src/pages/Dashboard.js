import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  Grid,
  Card,
  CardContent,
  CircularProgress,
} from '@mui/material';
import {
  virtualServiceApi,
  destinationRuleApi,
  gatewayApi,
  authorizationPolicyApi,
} from '../services/api';

export default function Dashboard() {
  const [stats, setStats] = useState({
    virtualServices: 0,
    destinationRules: 0,
    gateways: 0,
    authorizationPolicies: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const [vs, dr, gw, ap] = await Promise.all([
        virtualServiceApi.list(),
        destinationRuleApi.list(),
        gatewayApi.list(),
        authorizationPolicyApi.list(),
      ]);

      setStats({
        virtualServices: vs.data.length,
        destinationRules: dr.data.length,
        gateways: gw.data.length,
        authorizationPolicies: ap.data.length,
      });
    } catch (error) {
      console.error('Error loading stats:', error);
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
        Dashboard
      </Typography>
      <Typography variant="body1" color="text.secondary" paragraph>
        Overview of your Istio service mesh resources
      </Typography>

      <Grid container spacing={3}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Virtual Services
              </Typography>
              <Typography variant="h4">{stats.virtualServices}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Destination Rules
              </Typography>
              <Typography variant="h4">{stats.destinationRules}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Gateways
              </Typography>
              <Typography variant="h4">{stats.gateways}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Authorization Policies
              </Typography>
              <Typography variant="h4">{stats.authorizationPolicies}</Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Paper sx={{ mt: 3, p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Welcome to MeshControl Center
        </Typography>
        <Typography variant="body1" paragraph>
          MeshControl Center provides a comprehensive interface for managing your Istio service mesh.
          Use the navigation menu to access different resource types and features.
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Features:
        </Typography>
        <ul>
          <li>Visual configuration of traffic management rules</li>
          <li>Security policy management</li>
          <li>Service mesh topology visualization</li>
          <li>Scheduled actions for automated operations</li>
          <li>YAML preview and editing for all resources</li>
        </ul>
      </Paper>
    </Box>
  );
}
