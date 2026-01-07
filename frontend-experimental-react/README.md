# Become Overman - React Frontend

Modern React frontend for Become Overman quest management system with a clean, minimal design inspired by Continue AI.

## Features

- 🎨 Modern UI with Material-UI (MUI)
- 🔐 Authentication (Login/Register)
- 📋 Quest Management (Available, Active, Completed)
- 🛒 Quest Shop with search and recommendations
- 👥 Friends Management
- 👤 User Profile with XP and level tracking
- 📝 Quest Details with task completion
- 🎯 Clean, minimal design inspired by Continue AI

## Tech Stack

- React 18
- Material-UI (MUI) 5
- React Router 6
- Vite
- Axios

## Getting Started

### Installation

```bash
cd frontend-experimental-react
npm install
```

### Development

```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Build

```bash
npm run build
```

## Project Structure

```
src/
  ├── components/     # Reusable components
  │   └── Layout.jsx  # Main layout with navigation
  ├── contexts/       # React contexts
  │   └── AuthContext.jsx
  ├── pages/          # Page components
  │   ├── Login.jsx
  │   ├── Register.jsx
  │   ├── Quests.jsx
  │   ├── QuestDetails.jsx
  │   ├── Shop.jsx
  │   ├── Friends.jsx
  │   └── Profile.jsx
  ├── App.jsx         # Main app component with routing
  ├── main.jsx        # Entry point
  └── theme.js        # MUI theme configuration
```

## Pages

- **Login** - User authentication
- **Register** - User registration
- **Quests** - View available, active, and completed quests
- **Quest Details** - Detailed view of a quest with tasks
- **Shop** - Browse and purchase quests, search, recommendations
- **Friends** - Manage friends list
- **Profile** - User profile with stats and progress

## API Integration

The frontend expects the backend API to be running on `http://localhost:8080`.

Proxy configuration is set up in `vite.config.js` to forward requests to the backend.

## Design

- **Color Scheme**: Black, white, and accent color (#6366f1)
- **Typography**: System fonts with clean, modern styling
- **Components**: Minimal shadows, rounded corners, smooth transitions
- **Layout**: Responsive design with mobile-friendly navigation
