# Quick Start Guide - Deep Research Agent Web UI

Get up and running with the web UI in 5 minutes!

## Prerequisites

- Go 1.19 or higher
- Gin framework installed
- API server with required endpoints (see INTEGRATION_GUIDE.md)

## Step 1: Verify File Structure

Your `web/` directory should look like this:

```
web/
├── static/
│   ├── css/
│   │   └── style.css
│   ├── js/
│   │   ├── app.js
│   │   ├── auth.js
│   │   ├── research.js
│   │   └── websocket.js
│   └── index.html
├── server.go
└── README.md
```

All files are already in place at:
`/home/user/golang-learning-projects/week4-capstone/deep-research-agent/web/`

## Step 2: Add Web Routes to Your Server

Open your `main.go` and add these lines:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "deep-research-agent/web"
)

func main() {
    router := gin.Default()

    // Your API routes here
    setupAPIRoutes(router)

    // Add web UI (must be last!)
    web.ServeStaticFiles(router)

    router.Run(":8080")
}
```

## Step 3: Start the Server

```bash
cd /home/user/golang-learning-projects/week4-capstone/deep-research-agent
go run main.go
```

## Step 4: Access the UI

Open your browser and go to:

```
http://localhost:8080
```

You should see the login/register page!

## Step 5: Create an Account

1. Click "Register"
2. Enter your name, email, and password (min 8 characters)
3. Click "Register"
4. You'll be automatically logged in

## Step 6: Start Your First Research

1. Go to the "New Research" tab (default)
2. Enter a research query (e.g., "What are the latest AI trends?")
3. Select research depth (Medium is default)
4. Check your desired options (Web Search, Wikipedia)
5. Click "Start Research"

## Step 7: Watch Real-Time Progress

1. You'll be automatically taken to "Active Research" tab
2. Watch the progress bar move
3. See real-time logs as research progresses
4. Wait for completion

## Step 8: View Results

1. When complete, click "View Results"
2. See the full research report
3. Export as Markdown or JSON
4. Share or save your research

## Troubleshooting

### Can't access http://localhost:8080

- Check if server is running: `ps aux | grep research-agent`
- Check for port conflicts: `lsof -i :8080`
- Try a different port in main.go

### WebSocket not connecting

- Check browser console for errors (F12)
- Verify API endpoint returns proper WebSocket upgrade
- Falls back to polling automatically

### Login not working

- Check API server logs for errors
- Verify `/api/v1/auth/login` endpoint exists
- Check network tab in browser DevTools

## Next Steps

- Read the full [README.md](README.md) for all features
- Check [INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md) for API integration
- View [UI_MOCKUP.md](UI_MOCKUP.md) for visual reference
- See [SUMMARY.md](SUMMARY.md) for complete overview

## API Endpoints Needed

For the UI to work, implement these endpoints:

**Authentication**:
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/profile`

**Research**:
- `POST /api/v1/research/start`
- `GET /api/v1/research/:jobId/status`
- `WS /api/v1/research/:jobId/stream`

**Sessions**:
- `GET /api/v1/research/sessions`
- `GET /api/v1/research/sessions/:sessionId`

**Documents**:
- `POST /api/v1/documents/upload`
- `GET /api/v1/documents`

See INTEGRATION_GUIDE.md for complete implementation examples.

## Features Available

✅ User registration and login
✅ JWT token authentication
✅ Start new research with files
✅ Real-time progress via WebSocket
✅ View research results
✅ Session history with search
✅ Export results (Markdown/JSON)
✅ Document management
✅ Responsive design (mobile-friendly)
✅ Toast notifications
✅ Modal dialogs

## Support

For detailed documentation, see:
- [README.md](README.md) - Full documentation
- [INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md) - API integration
- [UI_MOCKUP.md](UI_MOCKUP.md) - Visual mockups
- [SUMMARY.md](SUMMARY.md) - Project summary

Happy researching! 🔬
