import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Box,
  Container,
  Typography,
  Card,
  CardContent,
  CardActions,
  Button,
  Grid,
  Chip,
  CircularProgress,
  Alert,
  Tabs,
  Tab,
} from '@mui/material'
import {
  LocalAtm,
  TrendingUp,
  PlayArrow,
  Visibility,
} from '@mui/icons-material'
import axios from 'axios'
import { useAuth } from '../contexts/AuthContext'
import Layout from '../components/Layout'

const Quests = () => {
  const navigate = useNavigate()
  const { refreshUserProfile } = useAuth()
  const [tabValue, setTabValue] = useState(0)
  const [quests, setQuests] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchQuests()
  }, [tabValue])

  const fetchQuests = async () => {
    setLoading(true)
    setError('')
    try {
      let endpoint = '/quests/available'
      if (tabValue === 1) endpoint = '/quests/active'
      if (tabValue === 2) endpoint = '/quests/completed'

      const response = await axios.get(endpoint)
      setQuests(response.data || [])
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to load quests')
    } finally {
      setLoading(false)
    }
  }

  const handleMenuOpen = (event) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const handleStartQuest = async (questId) => {
    try {
      await axios.post(`/quests/${questId}/start`)
      // Refresh quests and user profile (in case quest was purchased)
      await Promise.all([fetchQuests(), refreshUserProfile()])
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to start quest')
    }
  }

  const handleViewDetails = (questId) => {
    navigate(`/quests/${questId}/details`)
  }

  const getRarityColor = (rarity) => {
    const colors = {
      free: '#6b7280',
      common: '#3b82f6',
      rare: '#8b5cf6',
      epic: '#ec4899',
      legendary: '#f59e0b',
    }
    return colors[rarity] || colors.common
  }

  return (
    <Layout>
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h2" sx={{ mb: 1, fontWeight: 700 }}>
            Quests
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Complete quests to level up and earn rewards
          </Typography>
        </Box>

        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 4 }}>
          <Tabs
            value={tabValue}
            onChange={(e, newValue) => setTabValue(newValue)}
            sx={{
              '& .MuiTab-root': {
                textTransform: 'none',
                fontWeight: 600,
                fontSize: '0.95rem',
              },
            }}
          >
            <Tab label="Available" />
            <Tab label="Active" />
            <Tab label="Completed" />
          </Tabs>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }}>
            {error}
          </Alert>
        )}

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
            <CircularProgress />
          </Box>
        ) : quests.length === 0 ? (
          <Card sx={{ textAlign: 'center', py: 6 }}>
            <CardContent>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                No quests found
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {tabValue === 0
                  ? 'Check back later for new quests'
                  : tabValue === 1
                  ? 'You have no active quests'
                  : 'You have not completed any quests yet'}
              </Typography>
            </CardContent>
          </Card>
        ) : (
          <Grid container spacing={3}>
            {quests.map((quest) => (
              <Grid item xs={12} sm={6} md={4} key={quest.id}>
                <Card
                  sx={{
                    height: '100%',
                    display: 'flex',
                    flexDirection: 'column',
                    transition: 'transform 0.2s, box-shadow 0.2s',
                    '&:hover': {
                      transform: 'translateY(-4px)',
                      boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
                    },
                  }}
                >
                  <CardContent sx={{ flexGrow: 1 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                      <Chip
                        label={quest.rarity}
                        size="small"
                        sx={{
                          bgcolor: getRarityColor(quest.rarity),
                          color: 'white',
                          fontWeight: 600,
                          textTransform: 'capitalize',
                        }}
                      />
                      <Chip
                        label={`Difficulty ${quest.difficulty}`}
                        size="small"
                        variant="outlined"
                      />
                    </Box>

                    <Typography variant="h5" sx={{ mb: 1, fontWeight: 600 }}>
                      {quest.title}
                    </Typography>

                    <Typography
                      variant="body2"
                      color="text.secondary"
                      sx={{ mb: 2, minHeight: 40 }}
                    >
                      {quest.description?.substring(0, 100)}
                      {quest.description?.length > 100 ? '...' : ''}
                    </Typography>

                    <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                        <TrendingUp fontSize="small" color="primary" />
                        <Typography variant="body2" fontWeight={600}>
                          {quest.reward_xp} XP
                        </Typography>
                      </Box>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                        <LocalAtm fontSize="small" color="primary" />
                        <Typography variant="body2" fontWeight={600}>
                          {quest.reward_coin} coins
                        </Typography>
                      </Box>
                    </Box>

                    <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                      <Chip
                        label={quest.category}
                        size="small"
                        variant="outlined"
                        sx={{ textTransform: 'capitalize' }}
                      />
                      {quest.tasks_count && (
                        <Chip
                          label={`${quest.tasks_count} tasks`}
                          size="small"
                          variant="outlined"
                        />
                      )}
                    </Box>
                  </CardContent>

                  <CardActions sx={{ p: 2, pt: 0 }}>
                    {tabValue === 0 && (
                      <Button
                        fullWidth
                        variant="contained"
                        startIcon={<PlayArrow />}
                        onClick={() => handleStartQuest(quest.id)}
                        sx={{
                          bgcolor: 'primary.main',
                          '&:hover': {
                            bgcolor: 'primary.dark',
                          },
                        }}
                      >
                        Start Quest
                      </Button>
                    )}
                    {tabValue === 1 && (
                      <Button
                        fullWidth
                        variant="outlined"
                        startIcon={<Visibility />}
                        onClick={() => handleViewDetails(quest.id)}
                      >
                        View Details
                      </Button>
                    )}
                    {tabValue === 2 && (
                      <Button
                        fullWidth
                        variant="outlined"
                        disabled
                        sx={{ color: 'text.secondary' }}
                      >
                        Completed
                      </Button>
                    )}
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </Grid>
        )}
      </Container>
    </Layout>
  )
}

export default Quests
