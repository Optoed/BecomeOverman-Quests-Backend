import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Box,
  Container,
  Typography,
  Card,
  CardContent,
  Button,
  Chip,
  LinearProgress,
  Alert,
  CircularProgress,
  List,
  ListItem,
  ListItemText,
  ListItemButton,
  Divider,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material'
import {
  ArrowBack,
  CheckCircle,
  RadioButtonUnchecked,
  PlayArrow,
  LocalAtm,
  TrendingUp,
  Timer,
} from '@mui/icons-material'
import axios from 'axios'
import { useAuth } from '../contexts/AuthContext'
import Layout from '../components/Layout'

const QuestDetails = () => {
  const { questID } = useParams()
  const navigate = useNavigate()
  const { refreshUserProfile, token, loading: authLoading } = useAuth()
  const [quest, setQuest] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [completingTask, setCompletingTask] = useState(null)
  const [completeDialogOpen, setCompleteDialogOpen] = useState(false)

  useEffect(() => {
    // Only fetch if we have a token and auth is not loading
    if (!authLoading && token) {
      fetchQuestDetails()
    }
  }, [questID, token, authLoading])

  const fetchQuestDetails = async () => {
    setLoading(true)
    setError('')
    try {
      const response = await axios.get(`/quests/${questID}/details`)
      console.log('Quest details received:', response.data)
      console.log('Tasks statuses:', response.data.tasks?.map(t => ({ id: t.id, status: t.status })))
      setQuest(response.data)
    } catch (err) {
      console.error('Error fetching quest details:', err)
      const errorMessage = err.response?.data?.error || err.message || 'Failed to load quest details'
      setError(errorMessage)
      // If it's a 401 error, the interceptor will handle logout
      if (err.response?.status === 401) {
        console.warn('Unauthorized - token may be invalid')
      }
    } finally {
      setLoading(false)
    }
  }

  const handleCompleteTask = async (taskID) => {
    setCompletingTask(taskID)
    try {
      await axios.post(`/quests/${questID}/${taskID}/complete`)
      // Refresh quest details and user profile to update XP/coins
      await Promise.all([fetchQuestDetails(), refreshUserProfile()])
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to complete task')
    } finally {
      setCompletingTask(null)
    }
  }

  const handleCompleteQuest = async () => {
    try {
      await axios.post(`/quests/${questID}/complete`)
      // Refresh user profile to update XP/coins/level
      await refreshUserProfile()
      setCompleteDialogOpen(false)
      navigate('/quests')
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to complete quest')
    }
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

  const completedTasks = quest?.tasks?.filter((t) => t.status === 'completed')?.length || 0
  const totalTasks = quest?.tasks?.length || 0
  const progress = totalTasks > 0 ? (completedTasks / totalTasks) * 100 : 0

  if (authLoading || loading) {
    return (
      <Layout>
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
          <CircularProgress />
        </Box>
      </Layout>
    )
  }

  if (error && !quest) {
    return (
      <Layout>
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
          <Button startIcon={<ArrowBack />} onClick={() => navigate('/quests')}>
            Back to Quests
          </Button>
        </Container>
      </Layout>
    )
  }

  return (
    <Layout>
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Button
          startIcon={<ArrowBack />}
          onClick={() => navigate('/quests')}
          sx={{ mb: 3 }}
        >
          Back to Quests
        </Button>

        {error && (
          <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }}>
            {error}
          </Alert>
        )}

        {quest && (
          <>
            <Card sx={{ mb: 4 }}>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2, flexWrap: 'wrap', gap: 2 }}>
                  <Box>
                    <Chip
                      label={quest.rarity}
                      size="small"
                      sx={{
                        bgcolor: getRarityColor(quest.rarity),
                        color: 'white',
                        fontWeight: 600,
                        textTransform: 'capitalize',
                        mb: 1,
                      }}
                    />
                    <Typography variant="h3" sx={{ mb: 1, fontWeight: 700 }}>
                      {quest.title}
                    </Typography>
                    <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
                      {quest.description}
                    </Typography>
                  </Box>
                </Box>

                <Box sx={{ display: 'flex', gap: 3, mb: 3, flexWrap: 'wrap' }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <TrendingUp fontSize="small" color="primary" />
                    <Typography variant="body2" fontWeight={600}>
                      {quest.reward_xp} XP
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <LocalAtm fontSize="small" color="primary" />
                    <Typography variant="body2" fontWeight={600}>
                      {quest.reward_coin} coins
                    </Typography>
                  </Box>
                  <Chip
                    label={quest.category}
                    size="small"
                    variant="outlined"
                    sx={{ textTransform: 'capitalize' }}
                  />
                  <Chip
                    label={`Difficulty ${quest.difficulty}`}
                    size="small"
                    variant="outlined"
                  />
                </Box>

                <Box sx={{ mb: 2 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                    <Typography variant="body2" fontWeight={600}>
                      Progress
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {completedTasks} / {totalTasks} tasks completed
                    </Typography>
                  </Box>
                  <LinearProgress
                    variant="determinate"
                    value={progress}
                    sx={{
                      height: 8,
                      borderRadius: 4,
                      bgcolor: 'action.hover',
                    }}
                  />
                </Box>
              </CardContent>
            </Card>

            <Card>
              <CardContent>
                <Typography variant="h5" sx={{ mb: 3, fontWeight: 600 }}>
                  Tasks
                </Typography>

                <List>
                  {quest.tasks?.map((task, index) => (
                    <Box key={task.id}>
                      <ListItem
                        sx={{
                          bgcolor: task.status === 'completed' ? 'action.selected' : 'transparent',
                          borderRadius: 2,
                          mb: 1,
                        }}
                      >
                        <ListItemButton
                          disabled={task.status === 'completed' || completingTask === task.id}
                          onClick={() => {
                            // Allow completing if status is 'active' or null (not started yet but quest is active)
                            if (task.status !== 'completed') {
                              handleCompleteTask(task.id)
                            }
                          }}
                        >
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%' }}>
                            {task.status === 'completed' ? (
                              <CheckCircle color="success" />
                            ) : (
                              <RadioButtonUnchecked color="action" />
                            )}
                            <ListItemText
                              primary={
                                <Typography variant="body1" fontWeight={600}>
                                  {task.title}
                                </Typography>
                              }
                              secondary={
                                <Typography variant="body2" color="text.secondary">
                                  {task.description}
                                </Typography>
                              }
                            />
                            {/* Show Complete button if task is not completed (status is 'active' or null) */}
                            {task.status !== 'completed' && (
                              <Button
                                variant="contained"
                                size="small"
                                startIcon={<PlayArrow />}
                                disabled={completingTask === task.id}
                                onClick={(e) => {
                                  e.stopPropagation()
                                  handleCompleteTask(task.id)
                                }}
                              >
                                {completingTask === task.id ? 'Completing...' : 'Complete'}
                              </Button>
                            )}
                            {task.status === 'completed' && task.xp_gained && (
                              <Box sx={{ textAlign: 'right' }}>
                                <Typography variant="caption" color="success.main" fontWeight={600}>
                                  +{task.xp_gained} XP
                                </Typography>
                                {task.coin_gained && (
                                  <Typography variant="caption" color="success.main" display="block">
                                    +{task.coin_gained} coins
                                  </Typography>
                                )}
                              </Box>
                            )}
                          </Box>
                        </ListItemButton>
                      </ListItem>
                      {index < quest.tasks.length - 1 && <Divider sx={{ my: 1 }} />}
                    </Box>
                  ))}
                </List>

                {completedTasks === totalTasks && totalTasks > 0 && (
                  <Box sx={{ mt: 4, textAlign: 'center' }}>
                    <Button
                      variant="contained"
                      size="large"
                      onClick={() => setCompleteDialogOpen(true)}
                      sx={{
                        bgcolor: 'success.main',
                        '&:hover': {
                          bgcolor: 'success.dark',
                        },
                      }}
                    >
                      Complete Quest
                    </Button>
                  </Box>
                )}
              </CardContent>
            </Card>
          </>
        )}

        <Dialog open={completeDialogOpen} onClose={() => setCompleteDialogOpen(false)}>
          <DialogTitle>Complete Quest</DialogTitle>
          <DialogContent>
            <Typography>
              Are you sure you want to complete this quest? You will receive{' '}
              {quest?.reward_xp} XP and {quest?.reward_coin} coins.
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setCompleteDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleCompleteQuest} variant="contained" color="success">
              Complete
            </Button>
          </DialogActions>
        </Dialog>
      </Container>
    </Layout>
  )
}

export default QuestDetails
