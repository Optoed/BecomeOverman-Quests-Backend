import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  Avatar,
  Menu,
  MenuItem,
  IconButton,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Divider,
} from '@mui/material'
import {
  Logout,
  AccountCircle,
  Menu as MenuIcon,
  Dashboard,
  ShoppingBag,
  People,
  Person,
} from '@mui/icons-material'
import { useAuth } from '../contexts/AuthContext'

const Layout = ({ children }) => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuth()
  const [anchorEl, setAnchorEl] = useState(null)
  const [drawerOpen, setDrawerOpen] = useState(false)

  const menuItems = [
    { text: 'Quests', icon: <Dashboard />, path: '/quests' },
    { text: 'Shop', icon: <ShoppingBag />, path: '/shop' },
    { text: 'Friends', icon: <People />, path: '/friends' },
    { text: 'Profile', icon: <Person />, path: '/profile' },
  ]

  const handleMenuOpen = (event) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
    handleMenuClose()
  }

  const handleNavigate = (path) => {
    navigate(path)
    setDrawerOpen(false)
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <AppBar
        position="sticky"
        elevation={0}
        sx={{
          bgcolor: 'background.paper',
          borderBottom: '1px solid',
          borderColor: 'divider',
        }}
      >
        <Toolbar>
          <IconButton
            edge="start"
            onClick={() => setDrawerOpen(true)}
            sx={{ mr: 2, display: { md: 'none' } }}
          >
            <MenuIcon />
          </IconButton>

          <Typography
            variant="h6"
            sx={{
              flexGrow: 1,
              fontWeight: 700,
              cursor: 'pointer',
              background: 'linear-gradient(135deg, #000 0%, #6366f1 100%)',
              backgroundClip: 'text',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
            }}
            onClick={() => navigate('/quests')}
          >
            Become Overman
          </Typography>

          <Box sx={{ display: { xs: 'none', md: 'flex' }, gap: 1, mr: 2 }}>
            {menuItems.map((item) => (
              <Box
                key={item.path}
                onClick={() => handleNavigate(item.path)}
                sx={{
                  px: 2,
                  py: 1,
                  borderRadius: 2,
                  cursor: 'pointer',
                  bgcolor: location.pathname === item.path ? 'action.selected' : 'transparent',
                  '&:hover': {
                    bgcolor: 'action.hover',
                  },
                  transition: 'background-color 0.2s',
                }}
              >
                <Typography
                  variant="body2"
                  sx={{
                    fontWeight: location.pathname === item.path ? 600 : 500,
                    color: location.pathname === item.path ? 'primary.main' : 'text.primary',
                  }}
                >
                  {item.text}
                </Typography>
              </Box>
            ))}
          </Box>

          {user && (
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
              <Box sx={{ textAlign: 'right', display: { xs: 'none', sm: 'block' } }}>
                <Typography variant="body2" fontWeight={600}>
                  {user.username}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  Level {user.level} • {user.xp_points} XP
                </Typography>
              </Box>
              <IconButton onClick={handleMenuOpen} size="small">
                <Avatar sx={{ width: 32, height: 32, bgcolor: 'primary.main' }}>
                  {user.username?.[0]?.toUpperCase()}
                </Avatar>
              </IconButton>
            </Box>
          )}

          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'right',
            }}
            transformOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
          >
            <MenuItem onClick={() => { handleNavigate('/profile'); handleMenuClose(); }}>
              <AccountCircle sx={{ mr: 1 }} />
              Profile
            </MenuItem>
            <Divider />
            <MenuItem onClick={handleLogout}>
              <Logout sx={{ mr: 1 }} />
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      <Drawer
        anchor="left"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
      >
        <Box sx={{ width: 250, pt: 2 }}>
          <List>
            {menuItems.map((item) => (
              <ListItem key={item.path} disablePadding>
                <ListItemButton
                  selected={location.pathname === item.path}
                  onClick={() => handleNavigate(item.path)}
                >
                  <ListItemIcon>{item.icon}</ListItemIcon>
                  <ListItemText primary={item.text} />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
        </Box>
      </Drawer>

      <Box component="main" sx={{ flexGrow: 1, bgcolor: 'background.default' }}>
        {children}
      </Box>
    </Box>
  )
}

export default Layout
