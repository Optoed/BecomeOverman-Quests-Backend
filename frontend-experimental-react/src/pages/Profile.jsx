import { useState, useEffect } from 'react'
import {
  Box,
  Container,
  Typography,
  Card,
  CardContent,
  Grid,
  Avatar,
  LinearProgress,
  Chip,
} from '@mui/material'
import {
  Person,
  TrendingUp,
  LocalAtm,
  Star,
  CalendarToday,
} from '@mui/icons-material'
import axios from 'axios'
import { useAuth } from '../contexts/AuthContext'
import Layout from '../components/Layout'

const Profile = () => {
  const { user: authUser, refreshUserProfile } = useAuth()
  const [user, setUser] = useState(authUser)
  const [loading, setLoading] = useState(!authUser)

  // Sync with AuthContext user
  useEffect(() => {
    if (authUser) {
      setUser(authUser)
    } else if (!loading) {
      fetchProfile()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [authUser])

  const fetchProfile = async () => {
    setLoading(true)
    try {
      const response = await axios.get('/user/profile')
      setUser(response.data)
    } catch (error) {
      console.error('Failed to fetch profile:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <Layout>
        <Container maxWidth="md" sx={{ py: 4 }}>
          <Card>
            <CardContent>
              <Typography>Loading...</Typography>
            </CardContent>
          </Card>
        </Container>
      </Layout>
    )
  }

  if (!user) {
    return (
      <Layout>
        <Container maxWidth="md" sx={{ py: 4 }}>
          <Card>
            <CardContent>
              <Typography>Failed to load profile</Typography>
            </CardContent>
          </Card>
        </Container>
      </Layout>
    )
  }

  // Calculate XP progress to next level using quadratic formula
  // Formula: level = floor(sqrt(XP / 100)) + 1
  // To find XP needed for next level: (level^2 * 100) - currentXP
  const calculateLevel = (xp) => {
    if (xp < 0) xp = 0
    const baseXP = 100.0
    return Math.floor(Math.sqrt(xp / baseXP)) + 1
  }

  const currentLevel = user.level
  const currentXP = user.xp_points
  
  // XP needed for current level: (level - 1)^2 * 100
  const currentLevelXP = Math.pow(currentLevel - 1, 2) * 100
  
  // XP needed for next level: level^2 * 100
  const nextLevelXP = Math.pow(currentLevel, 2) * 100
  
  const xpInCurrentLevel = currentXP - currentLevelXP
  const xpNeededForNextLevel = nextLevelXP - currentLevelXP
  const xpProgress = xpNeededForNextLevel > 0 ? (xpInCurrentLevel / xpNeededForNextLevel) * 100 : 0

  return (
    <Layout>
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Typography variant="h2" sx={{ mb: 4, fontWeight: 700 }}>
          Profile
        </Typography>

        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 3, mb: 3 }}>
                  <Avatar
                    sx={{
                      width: 80,
                      height: 80,
                      bgcolor: 'primary.main',
                      fontSize: '2rem',
                    }}
                  >
                    {user.username?.[0]?.toUpperCase()}
                  </Avatar>
                  <Box>
                    <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
                      {user.username}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {user.email}
                    </Typography>
                  </Box>
                </Box>

                <Box sx={{ mb: 3 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                    <Typography variant="body2" fontWeight={600}>
                      Level {user.level}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {user.xp_points} / {nextLevelXP} XP
                    </Typography>
                  </Box>
                  <LinearProgress
                    variant="determinate"
                    value={xpProgress}
                    sx={{
                      height: 10,
                      borderRadius: 5,
                      bgcolor: 'action.hover',
                    }}
                  />
                  <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: 'block' }}>
                    {xpNeededForNextLevel - xpInCurrentLevel} XP to next level
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} sm={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                  <Box
                    sx={{
                      p: 1.5,
                      borderRadius: 2,
                      bgcolor: 'primary.main',
                      color: 'white',
                    }}
                  >
                    <TrendingUp />
                  </Box>
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Experience Points
                    </Typography>
                    <Typography variant="h5" fontWeight={700}>
                      {user.xp_points}
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} sm={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                  <Box
                    sx={{
                      p: 1.5,
                      borderRadius: 2,
                      bgcolor: 'secondary.main',
                      color: 'white',
                    }}
                  >
                    <LocalAtm />
                  </Box>
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Coins
                    </Typography>
                    <Typography variant="h5" fontWeight={700}>
                      {user.coin_balance}
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} sm={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                  <Box
                    sx={{
                      p: 1.5,
                      borderRadius: 2,
                      bgcolor: 'warning.main',
                      color: 'white',
                    }}
                  >
                    <Star />
                  </Box>
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Level
                    </Typography>
                    <Typography variant="h5" fontWeight={700}>
                      {user.level}
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} sm={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                  <Box
                    sx={{
                      p: 1.5,
                      borderRadius: 2,
                      bgcolor: 'info.main',
                      color: 'white',
                    }}
                  >
                    <CalendarToday />
                  </Box>
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Member Since
                    </Typography>
                    <Typography variant="body1" fontWeight={600}>
                      {new Date(user.created_at).toLocaleDateString('en-US', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                      })}
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Container>
    </Layout>
  )
}

export default Profile
