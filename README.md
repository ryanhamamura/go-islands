# React Islands Architecture with Go

This project demonstrates the React Islands Architecture pattern using:
- Go backend with Chi router
- React frontend with Vite build system
- Islands of interactivity in mostly server-rendered pages

## What is Islands Architecture?

Islands Architecture is a pattern where:
- Most of the page is static HTML rendered by the server
- Interactive "islands" of React components are hydrated on the client
- JavaScript is only loaded for the interactive components that need it

This approach provides:
- Better performance (reduced JavaScript bundle size)
- Better SEO (server-rendered content)
- Progressive enhancement

## Project Structure

```
go-react-islands/
├── backend/            # Go server with Chi router
│   ├── main.go         # Server implementation
│   ├── templates/      # HTML templates
│   └── static/         # Built assets (after frontend build)
└── frontend/           # React/Vite project
    ├── src/
    │   ├── islands/    # React component islands
    │   │   ├── Counter.jsx
    │   │   └── UserProfile.jsx
    │   ├── main.jsx    # Main entry point
    │   └── islands-client.js # Island hydration logic
    ├── index.html      # Development index
    ├── package.json
    └── vite.config.js  # Vite configuration
```

## Getting Started

### Prerequisites
- Go 1.16+
- Node.js 14+ and npm

### Development

1. Clone the repository
2. Install frontend dependencies:
   ```
   cd go-react-islands/frontend
   npm install
   ```
3. Install Air for hot reloading (optional):
   ```
   go install github.com/air-verse/air@latest
   ```
4. Configure environment (optional):
   ```
   cd go-react-islands/backend
   cp .env.example .env
   # Edit .env file with your configuration
   ```
5. Run development servers (uses Make):
   ```
   make dev
   ```
   This starts:
   - Vite dev server on port 5173 (hot reloading for frontend)
   - Go backend server on port 8080
   
   Alternatively, run with Air for backend hot reloading:
   ```
   make dev-air
   ```
   This uses Air to automatically rebuild and restart the Go server when files change.

### Production Build

1. Build frontend assets:
   ```
   make build
   ```
   This bundles React components into the backend/static directory

2. Run production server:
   ```
   make run-backend
   ```

## How It Works

1. **Server-Side Rendering**: Go templates render the main HTML structure
2. **Island Components**: Specific div elements are marked with data attributes
3. **Selective Hydration**: JavaScript only hydrates the interactive islands
4. **Asset Optimization**: Vite bundles JavaScript per island component

## Key Implementation Details

- **Frontend**: Vite is configured to build separate chunks for each island
- **Backend**: Chi router serves both the HTML templates and the API
- **Integration**: Templates include the necessary JavaScript only in production
- **Development**: Proxy setup allows frontend dev server to work with backend API

## Key Integration Features

### 1. Vite-Go Integration

The project uses [vite-go](https://github.com/torenware/vite-go) for seamless integration between the Go backend and Vite frontend:

- Automatically parses the Vite manifest file
- Handles script and link tags for both development and production
- Provides template functions to include the correct assets

### 2. Air Hot Reloading

For development, the project supports [Air](https://github.com/air-verse/air) for hot reloading:

- Automatically rebuilds and restarts the Go server when code changes
- Works alongside Vite's hot module replacement
- Configured via `.air.toml` in the backend directory

## Next Steps

- Add a more comprehensive backend API
- Implement proper CSS bundling and optimization
- Add server-side rendering of the initial component state
- Implement progressive enhancement for core functionality