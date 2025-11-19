# Deep Research Agent - Web UI

A modern, clean web interface for the Deep Research Agent API. Built with vanilla JavaScript, CSS3, and WebSocket for real-time updates.

## Features

### Core Functionality

- **User Authentication**
  - User registration with email and password
  - Secure login with JWT token-based authentication
  - Session management with automatic token expiration handling

- **Research Operations**
  - Start new research with customizable parameters
  - Real-time progress tracking via WebSocket
  - Upload PDF/DOCX documents for research context
  - Configure research depth (shallow, medium, deep)
  - Enable/disable web search and Wikipedia
  - Optional fact-checking

- **Active Research Monitoring**
  - Live progress updates with percentage completion
  - Real-time step-by-step logging
  - Tool execution tracking
  - Elapsed time counter
  - WebSocket connection with automatic reconnection

- **Session Management**
  - View all past research sessions
  - Search and filter sessions
  - Pagination for large result sets
  - View detailed results in modal
  - Export results in Markdown or JSON format
  - Delete unwanted sessions

- **Document Management**
  - Upload multiple documents (PDF, DOCX)
  - View uploaded documents
  - Delete documents
  - Track document usage in research

## Architecture

```
web/
├── static/
│   ├── css/
│   │   └── style.css           # All styles with CSS variables
│   ├── js/
│   │   ├── auth.js             # Authentication module
│   │   ├── research.js         # Research operations module
│   │   ├── websocket.js        # WebSocket handling module
│   │   └── app.js              # Main application logic
│   └── index.html              # Single-page application
├── templates/                   # Optional server-side templates
├── assets/                      # Images, icons, etc.
├── server.go                    # Go static file server
└── README.md                    # This file
```

## Technology Stack

- **Frontend**: Vanilla JavaScript (ES6+)
- **Styling**: CSS3 with CSS Variables, Flexbox, Grid
- **Real-time**: WebSocket API with polling fallback
- **HTTP**: Fetch API for RESTful calls
- **Storage**: LocalStorage for token and user data
- **Backend Integration**: Gin framework (Go)

## Getting Started

### Prerequisites

The API server must be running with the following endpoints:

**Authentication**:
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/profile` - Get user profile

**Research**:
- `POST /api/v1/research/start` - Start research job
- `GET /api/v1/research/{jobId}/status` - Get job status
- `WS /api/v1/research/{jobId}/stream` - WebSocket for real-time updates
- `GET /api/v1/research/sessions` - List all sessions
- `GET /api/v1/research/sessions/{sessionId}` - Get session details
- `DELETE /api/v1/research/sessions/{sessionId}` - Delete session
- `GET /api/v1/research/sessions/{sessionId}/export` - Export session

**Documents**:
- `POST /api/v1/documents/upload` - Upload documents
- `GET /api/v1/documents` - List documents
- `GET /api/v1/documents/{documentId}` - Get document info
- `DELETE /api/v1/documents/{documentId}` - Delete document

### Integration with Go API Server

Add the following to your main API server file:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "deep-research-agent/web"
    // ... other imports
)

func main() {
    router := gin.Default()

    // Optional: Enable CORS for development
    router.Use(web.CORSMiddleware())

    // Setup API routes first
    api := router.Group("/api/v1")
    {
        // Authentication routes
        auth := api.Group("/auth")
        {
            auth.POST("/register", handleRegister)
            auth.POST("/login", handleLogin)
            auth.GET("/profile", authMiddleware, handleProfile)
        }

        // Research routes
        research := api.Group("/research")
        research.Use(authMiddleware)
        {
            research.POST("/start", handleStartResearch)
            research.GET("/:jobId/status", handleJobStatus)
            research.GET("/:jobId/stream", handleWebSocket)
            research.GET("/sessions", handleListSessions)
            research.GET("/sessions/:sessionId", handleGetSession)
            research.DELETE("/sessions/:sessionId", handleDeleteSession)
            research.GET("/sessions/:sessionId/export", handleExportSession)
        }

        // Document routes
        documents := api.Group("/documents")
        documents.Use(authMiddleware)
        {
            documents.POST("/upload", handleUploadDocuments)
            documents.GET("", handleListDocuments)
            documents.GET("/:documentId", handleGetDocument)
            documents.DELETE("/:documentId", handleDeleteDocument)
        }
    }

    // Setup web UI routes (MUST be last)
    web.ServeStaticFiles(router)

    // Start server
    router.Run(":8080")
}
```

### Accessing the UI

Once the server is running:

1. Open your browser and navigate to `http://localhost:8080`
2. You'll see the login/register page
3. Register a new account or login with existing credentials
4. Start using the research agent!

## User Interface

### Login/Register Page

The initial page shows a clean authentication form with:
- Toggle between login and register modes
- Email and password fields
- Form validation
- Password strength requirements (8+ characters)

### Dashboard

After login, you'll see the main dashboard with four tabs:

#### 1. New Research Tab

- **Research Query**: Text area for your research question
- **Research Depth**: Dropdown to select shallow (2-3 min), medium (5-10 min), or deep (15-20 min)
- **File Upload**: Upload PDF or DOCX files (optional)
- **Options**: Checkboxes for web search, Wikipedia, and fact-checking
- **Submit Button**: Start the research

#### 2. Active Research Tab

Shows all currently running research jobs with:
- Query display
- Progress bar (0-100%)
- Current step description
- Real-time log entries
- Elapsed time counter
- View Results button (when complete)

#### 3. Sessions History Tab

Table view of all past research sessions:
- Query column
- Date and time
- Duration
- Status badge (completed/failed)
- Actions: View, Export, Delete
- Search bar to filter sessions
- Pagination controls

#### 4. Documents Tab

Grid view of uploaded documents:
- Document name
- File size
- Upload date
- Info and Delete actions
- Upload button to add new documents

## UI Components

### Navigation Bar

Sticky header with:
- Research Agent logo
- User name display
- Logout button

### Progress Card

Real-time research progress display:
```
┌─────────────────────────────────────┐
│ What are the latest AI trends?     │
│ ▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░ 60%          │
│ Step 3 of 6: Analyzing Wikipedia   │
│                                     │
│ ✓ Web search completed (2.3s)      │
│ ⟳ Wikipedia search in progress...  │
│                                     │
│ Elapsed: 00:45                      │
└─────────────────────────────────────┘
```

### Modals

- **Session Details Modal**: View full research results
- **Upload Modal**: Upload multiple documents with preview

### Toast Notifications

Non-intrusive notifications for:
- Success messages (green)
- Error messages (red)
- Info messages (blue)
- Auto-dismiss after 5 seconds

## JavaScript Modules

### auth.js

Handles all authentication operations:

```javascript
// Login user
await login(email, password);

// Register new user
await register(email, password, name);

// Logout
logout();

// Check if authenticated
const isAuth = isAuthenticated();

// Get current user
const user = getCurrentUser();

// Make authenticated API request
const response = await authenticatedFetch(url, options);
```

### research.js

Manages research operations:

```javascript
// Start research
const result = await startResearch(query, depth, files, options);

// Get job status
const status = await getResearchStatus(jobId);

// List sessions
const sessions = await listSessions(page, limit);

// Get session details
const session = await getSession(sessionId);

// Export session
const { blob, filename } = await exportSession(sessionId, 'markdown');

// Upload files
const result = await uploadFiles(files);
```

### websocket.js

Real-time WebSocket communication:

```javascript
// Create WebSocket connection
const ws = connectResearchWebSocket(jobId, token);

// Register callbacks
ws.on('onProgress', (data) => {
    console.log(`Progress: ${data.percentage}%`);
});

ws.on('onStepCompleted', (data) => {
    console.log(`Completed: ${data.step}`);
});

ws.on('onComplete', (data) => {
    console.log('Research completed!');
});

// Disconnect
ws.disconnect();
```

### app.js

Main application logic and UI management:
- Page initialization
- Tab switching
- Form handling
- Event listeners
- UI updates
- Modal management
- Toast notifications

## Styling

### Design System

**Color Palette**:
- Primary: `#2563eb` (Blue)
- Success: `#10b981` (Green)
- Danger: `#ef4444` (Red)
- Warning: `#f59e0b` (Orange)

**Typography**:
- System font stack for native look
- Font sizes: 0.85rem - 2rem
- Line height: 1.6

**Spacing**:
- Uses consistent spacing scale (0.25rem - 2rem)
- CSS variables for easy customization

**Components**:
- Rounded corners (0.25rem - 0.75rem)
- Subtle shadows for depth
- Smooth transitions (0.2s ease)

### Responsive Design

Mobile-friendly with breakpoints:
- Desktop: Full layout
- Tablet: Adjusted spacing and column layout
- Mobile (< 768px): Single column, stacked elements

## Browser Compatibility

**Tested and working on**:
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

**Required features**:
- ES6+ JavaScript support
- CSS Grid and Flexbox
- WebSocket API
- Fetch API
- LocalStorage API

## WebSocket Protocol

### Connection

```
ws://localhost:8080/api/v1/research/{jobId}/stream?token={jwt_token}
```

### Message Types

```json
{
  "type": "progress",
  "data": {
    "percentage": 60,
    "current": 3,
    "total": 5,
    "message": "Analyzing data..."
  }
}
```

```json
{
  "type": "step_started",
  "data": {
    "step": "Web Search",
    "step_number": 2,
    "total_steps": 6
  }
}
```

```json
{
  "type": "step_completed",
  "data": {
    "step": "Web Search",
    "duration": "2.3s",
    "result": "Found 10 relevant sources"
  }
}
```

```json
{
  "type": "research_completed",
  "data": {
    "session_id": "session123",
    "results": "...",
    "total_duration": "5m 32s"
  }
}
```

```json
{
  "type": "error",
  "data": {
    "message": "Failed to fetch data",
    "details": "Network timeout"
  }
}
```

## Error Handling

### Network Errors

- Displays user-friendly error messages
- Automatically retries WebSocket connections
- Falls back to polling if WebSocket fails

### Authentication Errors

- Detects expired tokens (401 responses)
- Automatically redirects to login
- Clears local storage on logout

### Validation Errors

- Client-side form validation
- Server-side error display
- Field-specific error messages

## Security Considerations

### Token Storage

- JWT tokens stored in `localStorage`
- Tokens sent in `Authorization` header
- Automatic token expiration checking

### XSS Prevention

- All user input is escaped before display
- Uses `textContent` instead of `innerHTML` where possible
- Sanitizes HTML in results display

### CORS

- CORS middleware included for development
- Configure properly for production deployment

## Customization

### Changing Colors

Edit CSS variables in `style.css`:

```css
:root {
    --primary: #2563eb;        /* Your primary color */
    --success: #10b981;        /* Success color */
    --danger: #ef4444;         /* Error color */
    /* ... */
}
```

### Modifying API Endpoints

Edit constants in JavaScript files:

```javascript
// In auth.js
const AUTH_API_BASE = '/api/v1/auth';

// In research.js
const RESEARCH_API_BASE = '/api/v1/research';
```

### Adding New Features

1. Add HTML structure to `index.html`
2. Style components in `style.css`
3. Implement logic in appropriate JS module
4. Connect to API endpoints

## Performance

### Optimization Techniques

- **Code Splitting**: Separate modules for maintainability
- **Debouncing**: Search input debounced to 500ms
- **Lazy Loading**: Data loaded only when needed
- **WebSocket**: Real-time updates without polling overhead
- **Caching**: User data cached in localStorage

### Best Practices

- Minimal DOM manipulation
- Event delegation where possible
- Efficient CSS selectors
- Compressed assets in production
- HTTP/2 for multiplexing

## Troubleshooting

### WebSocket not connecting

1. Check if API server is running
2. Verify WebSocket endpoint is correct
3. Check browser console for errors
4. Falls back to polling automatically

### Login not working

1. Verify API server is running on correct port
2. Check network tab for API response
3. Ensure CORS is configured correctly
4. Check JWT token format

### Styles not loading

1. Verify static file server is configured
2. Check file paths in `index.html`
3. Clear browser cache
4. Check browser console for 404 errors

## Development

### Running locally

1. Ensure Go API server is running
2. Access `http://localhost:8080`
3. Open browser DevTools for debugging
4. Check console for any errors

### Making changes

1. Edit files in `web/static/`
2. Refresh browser (may need hard refresh: Ctrl+F5)
3. Use browser DevTools for debugging
4. Test in multiple browsers

### Adding dependencies

This project uses vanilla JavaScript with no build step. To add libraries:

1. Download library files to `static/js/`
2. Add `<script>` tag to `index.html`
3. Or use CDN links

## Production Deployment

### Build checklist

- [ ] Minify CSS and JavaScript
- [ ] Optimize images
- [ ] Configure CORS properly
- [ ] Enable HTTPS
- [ ] Set secure cookie flags
- [ ] Configure CSP headers
- [ ] Enable gzip compression
- [ ] Set cache headers
- [ ] Remove console.log statements
- [ ] Test in all target browsers

### Environment configuration

- Use environment variables for API URLs
- Configure different settings for dev/staging/prod
- Set appropriate token expiration times
- Configure secure WebSocket (wss://)

## License

Part of the Deep Research Agent project.

## Support

For issues or questions:
1. Check the troubleshooting section
2. Review browser console for errors
3. Check API server logs
4. Verify all endpoints are working

## Future Enhancements

Potential features to add:
- Dark mode toggle
- Keyboard shortcuts
- Advanced search filters
- Research templates
- Collaborative features
- Export to more formats (PDF, Word)
- Statistics dashboard
- Research history graphs
- Syntax highlighting for code results
- Markdown preview in results
