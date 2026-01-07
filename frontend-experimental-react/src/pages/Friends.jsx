import { useState, useEffect } from 'react'
import {
  Box,
  Container,
  Typography,
  Card,
  CardContent,
  TextField,
  Button,
  Alert,
  Avatar,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  InputAdornment,
  CircularProgress,
} from '@mui/material'
import {
  PersonAdd,
  Search,
  Person,
} from '@mui/icons-material'
import axios from 'axios'
import Layout from '../components/Layout'

const Friends = () => {
  const [friends, setFriends] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [addDialogOpen, setAddDialogOpen] = useState(false)
  const [friendName, setFriendName] = useState('')
  const [addingFriend, setAddingFriend] = useState(false)

  useEffect(() => {
    fetchFriends()
  }, [])

  const fetchFriends = async () => {
    setLoading(true)
    setError('')
    try {
      const response = await axios.get('/friends')
      setFriends(response.data || [])
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to load friends')
    } finally {
      setLoading(false)
    }
  }

  const handleAddFriend = async () => {
    if (!friendName.trim()) {
      setError('Please enter a friend name')
      return
    }

    setAddingFriend(true)
    setError('')
    try {
      await axios.post(`/friends/by-name/${friendName}`)
      setSuccess('Friend added successfully!')
      setFriendName('')
      setAddDialogOpen(false)
      fetchFriends()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to add friend')
    } finally {
      setAddingFriend(false)
    }
  }

  return (
    <Layout>
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
          <Typography variant="h2" sx={{ fontWeight: 700 }}>
            Friends
          </Typography>
          <Button
            variant="contained"
            startIcon={<PersonAdd />}
            onClick={() => setAddDialogOpen(true)}
            sx={{
              bgcolor: 'primary.main',
              '&:hover': {
                bgcolor: 'primary.dark',
              },
            }}
          >
            Add Friend
          </Button>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }}>
            {error}
          </Alert>
        )}

        {success && (
          <Alert severity="success" sx={{ mb: 3, borderRadius: 2 }}>
            {success}
          </Alert>
        )}

        <Card>
          <CardContent>
            {loading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
                <CircularProgress />
              </Box>
            ) : friends.length === 0 ? (
              <Box sx={{ textAlign: 'center', py: 6 }}>
                <Person sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  No friends yet
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                  Add friends to share quests and compete together
                </Typography>
                <Button
                  variant="contained"
                  startIcon={<PersonAdd />}
                  onClick={() => setAddDialogOpen(true)}
                >
                  Add Your First Friend
                </Button>
              </Box>
            ) : (
              <List>
                {friends.map((friend, index) => (
                  <ListItem
                    key={friend.id}
                    sx={{
                      borderRadius: 2,
                      mb: 1,
                      '&:hover': {
                        bgcolor: 'action.hover',
                      },
                    }}
                  >
                    <ListItemAvatar>
                      <Avatar sx={{ bgcolor: 'primary.main' }}>
                        {friend.username?.[0]?.toUpperCase()}
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={
                        <Typography variant="body1" fontWeight={600}>
                          {friend.username}
                        </Typography>
                      }
                      secondary={
                        <Typography variant="caption" color="text.secondary">
                          Added {new Date(friend.created_at).toLocaleDateString()}
                        </Typography>
                      }
                    />
                  </ListItem>
                ))}
              </List>
            )}
          </CardContent>
        </Card>

        <Dialog open={addDialogOpen} onClose={() => setAddDialogOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Add Friend</DialogTitle>
          <DialogContent>
            <TextField
              fullWidth
              label="Friend Username"
              value={friendName}
              onChange={(e) => {
                setFriendName(e.target.value)
                setError('')
              }}
              margin="normal"
              autoFocus
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search />
                  </InputAdornment>
                ),
              }}
              helperText="Enter your friend's username"
            />
            {error && (
              <Alert severity="error" sx={{ mt: 2, borderRadius: 2 }}>
                {error}
              </Alert>
            )}
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setAddDialogOpen(false)}>Cancel</Button>
            <Button
              onClick={handleAddFriend}
              variant="contained"
              disabled={addingFriend || !friendName.trim()}
            >
              {addingFriend ? 'Adding...' : 'Add Friend'}
            </Button>
          </DialogActions>
        </Dialog>
      </Container>
    </Layout>
  )
}

export default Friends
