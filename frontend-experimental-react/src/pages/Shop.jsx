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
  TextField,
  InputAdornment,
  Tabs,
  Tab,
} from '@mui/material'
import {
  Search,
  ShoppingCart,
  LocalAtm,
  TrendingUp,
  PlayArrow,
} from '@mui/icons-material'
import axios from 'axios'
import { useAuth } from '../contexts/AuthContext'
import Layout from '../components/Layout'

const Shop = () => {
  const navigate = useNavigate()
  const { refreshUserProfile } = useAuth()
  const [quests, setQuests] = useState([])
  const [recommendedQuests, setRecommendedQuests] = useState([])
  const [loading, setLoading] = useState(true)
  const [searchLoading, setSearchLoading] = useState(false)
  const [error, setError] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [tabValue, setTabValue] = useState(0)

  useEffect(() => {
    if (tabValue === 0) {
      fetchShopQuests()
    } else if (tabValue === 1) {
      fetchRecommendedQuests()
    }
  }, [tabValue])

  const fetchShopQuests = async () => {
    setLoading(true)
    setError('')
    try {
      const response = await axios.get('/quests/shop')
      setQuests(response.data || [])
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to load shop quests')
    } finally {
      setLoading(false)
    }
  }

  const fetchRecommendedQuests = async () => {
    setLoading(true)
    setError('')
    try {
      const response = await axios.post('/quests/recommend')
      // Recommendations structure: { recommendations: [...], user_profile_info: {...} }
      setRecommendedQuests(response.data?.recommendations || [])
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to load recommendations')
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      fetchShopQuests()
      return
    }

    setSearchLoading(true)
    setError('')
    try {
      const response = await axios.post('/quests/search', {
        query: searchQuery,
        status: 'all',
      })
      // Search returns QuestWithSimilarityScore[] - extract quest objects
      const searchResults = response.data || []
      const questsFromSearch = searchResults.map((item) => ({
        ...item.quest,
        similarity_score: item.similarity_score,
      }))
      setQuests(questsFromSearch)
    } catch (err) {
      setError(err.response?.data?.error || 'Search failed')
    } finally {
      setSearchLoading(false)
    }
  }

  const handlePurchase = async (questId) => {
    try {
      await axios.post(`/quests/${questId}/purchase`)
      // Refresh shop quests and user profile to update coin balance
      await Promise.all([fetchShopQuests(), refreshUserProfile()])
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to purchase quest')
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

  const displayQuests = tabValue === 0 ? quests : recommendedQuests

  return (
    <Layout>
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h2" sx={{ mb: 1, fontWeight: 700 }}>
            Quest Shop
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Discover and purchase new quests
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
            <Tab label="All Quests" />
            <Tab label="Recommended" />
          </Tabs>
        </Box>

        {tabValue === 0 && (
          <Box sx={{ mb: 4 }}>
            <TextField
              fullWidth
              placeholder="Search quests..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyPress={(e) => {
                if (e.key === 'Enter') {
                  handleSearch()
                }
              }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search />
                  </InputAdornment>
                ),
              }}
              sx={{ mb: 2 }}
            />
            <Button
              variant="outlined"
              onClick={handleSearch}
              disabled={searchLoading}
              sx={{ mb: 3 }}
            >
              {searchLoading ? 'Searching...' : 'Search'}
            </Button>
          </Box>
        )}

        {error && (
          <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }}>
            {error}
          </Alert>
        )}

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
            <CircularProgress />
          </Box>
        ) : displayQuests.length === 0 ? (
          <Card sx={{ textAlign: 'center', py: 6 }}>
            <CardContent>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                No quests found
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {tabValue === 0
                  ? 'Try searching for something else'
                  : 'Complete more quests to get recommendations'}
              </Typography>
            </CardContent>
          </Card>
        ) : (
          <Grid container spacing={3}>
            {displayQuests.map((quest) => (
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
                      {quest.rarity && (
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
                      )}
                      {quest.price !== undefined && (
                        <Chip
                          label={`${quest.price} coins`}
                          size="small"
                          variant="outlined"
                          icon={<LocalAtm fontSize="small" />}
                        />
                      )}
                      {(quest.similarity_score !== undefined || quest.similarityScore !== undefined) && (
                        <Chip
                          label={`${((quest.similarity_score || quest.similarityScore) * 100).toFixed(0)}% match`}
                          size="small"
                          variant="outlined"
                          color="secondary"
                        />
                      )}
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

                    {(quest.reward_xp !== undefined || quest.reward_coin !== undefined) && (
                      <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
                        {quest.reward_xp !== undefined && (
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <TrendingUp fontSize="small" color="primary" />
                            <Typography variant="body2" fontWeight={600}>
                              {quest.reward_xp} XP
                            </Typography>
                          </Box>
                        )}
                        {quest.reward_coin !== undefined && (
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <LocalAtm fontSize="small" color="primary" />
                            <Typography variant="body2" fontWeight={600}>
                              {quest.reward_coin} coins
                            </Typography>
                          </Box>
                        )}
                      </Box>
                    )}

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
                    {quest.price !== undefined ? (
                      <Button
                        fullWidth
                        variant="contained"
                        startIcon={<ShoppingCart />}
                        onClick={() => handlePurchase(quest.id)}
                        sx={{
                          bgcolor: 'primary.main',
                          '&:hover': {
                            bgcolor: 'primary.dark',
                          },
                        }}
                      >
                        Purchase
                      </Button>
                    ) : (
                      <Button
                        fullWidth
                        variant="outlined"
                        onClick={() => navigate(`/quests/${quest.id}/details`)}
                      >
                        View Details
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

export default Shop
