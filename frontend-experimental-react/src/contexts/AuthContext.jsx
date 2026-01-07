import { createContext, useContext, useState, useEffect, useCallback } from 'react'
import axios from 'axios'

const AuthContext = createContext(null)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null)
  const [token, setToken] = useState(() => {
    // Initialize token from localStorage on mount
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      // Set axios header immediately
      axios.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
    }
    return storedToken
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (token) {
      // Ensure header is set
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
      fetchUserProfile()
    } else {
      setLoading(false)
    }
  }, [token])

  const fetchUserProfile = async () => {
    // Ensure token is set in headers before making request
    const currentToken = localStorage.getItem('token')
    if (currentToken) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${currentToken}`
    }
    
    try {
      const response = await axios.get('/user/profile')
      setUser(response.data)
      return response.data
    } catch (error) {
      console.error('Failed to fetch user profile:', error)
      // Only logout on 401, not on other errors
      if (error.response?.status === 401) {
        logout()
      }
      throw error
    } finally {
      setLoading(false)
    }
  }

  // Public method to refresh user profile
  const refreshUserProfile = async () => {
    if (token) {
      try {
        const userData = await fetchUserProfile()
        return userData
      } catch (error) {
        console.error('Failed to refresh user profile:', error)
        return null
      }
    }
    return null
  }

  const login = async (username, password) => {
    try {
      const response = await axios.post('/user/login', { username, password })
      const newToken = response.data.token
      setToken(newToken)
      localStorage.setItem('token', newToken)
      axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      await fetchUserProfile()
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Login failed',
      }
    }
  }

  const register = async (username, email, password) => {
    try {
      await axios.post('/user/register', { username, email, password })
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Registration failed',
      }
    }
  }

  const logout = useCallback(() => {
    setToken(null)
    setUser(null)
    localStorage.removeItem('token')
    delete axios.defaults.headers.common['Authorization']
  }, [])

  // Setup axios interceptors
  useEffect(() => {
    // Request interceptor: ensure token is always set from localStorage
    const requestInterceptor = axios.interceptors.request.use(
      (config) => {
        const currentToken = localStorage.getItem('token')
        if (currentToken) {
          config.headers.Authorization = `Bearer ${currentToken}`
        }
        return config
      },
      (error) => Promise.reject(error)
    )

    // Response interceptor: handle token errors
    const responseInterceptor = axios.interceptors.response.use(
      (response) => response,
      (error) => {
        // Only handle 401 if we have a token (means token is invalid/expired)
        // If no token, it's a normal unauthenticated request
        if (error.response?.status === 401) {
          const currentToken = localStorage.getItem('token')
          if (currentToken) {
            // Token exists but is invalid/expired
            console.warn('Token is invalid or expired, logging out')
            logout()
            // Redirect to login if not already there
            if (window.location.pathname !== '/login' && window.location.pathname !== '/register') {
              window.location.href = '/login'
            }
          }
          // If no token, just reject the error normally (user needs to login)
        }
        return Promise.reject(error)
      }
    )

    return () => {
      axios.interceptors.request.eject(requestInterceptor)
      axios.interceptors.response.eject(responseInterceptor)
    }
  }, [logout])

  return (
    <AuthContext.Provider value={{ user, token, loading, login, register, logout, refreshUserProfile }}>
      {children}
    </AuthContext.Provider>
  )
}
