/**
 * Main Application Module
 * Handles UI initialization, routing, and user interactions
 */

// Global state
let currentPage = 'new-research';
let currentSessionId = null;
let activeJobs = new Map(); // jobId -> connection object
let sessionsPage = 1;
const sessionsLimit = 20;

/**
 * Initialize application on page load
 */
document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
});

/**
 * Initialize application
 */
function initializeApp() {
    // Check authentication status
    if (isAuthenticated()) {
        showDashboard();
        loadUserInfo();
        setupEventListeners();
        loadInitialData();
    } else {
        showAuthPage();
        setupAuthListeners();
    }
}

/**
 * Show authentication page
 */
function showAuthPage() {
    document.getElementById('auth-page').classList.add('active');
    document.getElementById('dashboard-page').classList.remove('active');
}

/**
 * Show dashboard page
 */
function showDashboard() {
    document.getElementById('auth-page').classList.remove('active');
    document.getElementById('dashboard-page').classList.add('active');
}

/**
 * Setup authentication form listeners
 */
function setupAuthListeners() {
    // Login form
    const loginForm = document.getElementById('login-form');
    loginForm.addEventListener('submit', handleLogin);

    // Register form
    const registerForm = document.getElementById('register-form');
    registerForm.addEventListener('submit', handleRegister);
}

/**
 * Setup dashboard event listeners
 */
function setupEventListeners() {
    // Tab navigation
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            switchTab(btn.dataset.tab);
        });
    });

    // Research form
    const researchForm = document.getElementById('research-form');
    researchForm.addEventListener('submit', handleResearchSubmit);

    // File input preview
    const fileInput = document.getElementById('research-files');
    fileInput.addEventListener('change', handleFilePreview);

    // Sessions search
    const searchInput = document.getElementById('sessions-search');
    searchInput.addEventListener('input', debounce(handleSessionsSearch, 500));

    // Upload modal file input
    const uploadFileInput = document.getElementById('upload-files-input');
    if (uploadFileInput) {
        uploadFileInput.addEventListener('change', handleUploadFilePreview);
    }
}

/**
 * Load user information
 */
async function loadUserInfo() {
    try {
        const user = getCurrentUser();
        if (user) {
            document.getElementById('user-name-display').textContent = user.name || user.email;
        } else {
            // Fetch from API if not in localStorage
            const profile = await fetchUserProfile();
            document.getElementById('user-name-display').textContent = profile.name || profile.email;
        }
    } catch (error) {
        console.error('Failed to load user info:', error);
    }
}

/**
 * Load initial data
 */
function loadInitialData() {
    loadSessions();
}

/**
 * Handle login form submission
 */
async function handleLogin(e) {
    e.preventDefault();

    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        showToast('Logging in...', 'info');
        await login(email, password);
        showToast('Login successful!', 'success');

        // Reload page to show dashboard
        window.location.reload();
    } catch (error) {
        showToast(error.message || 'Login failed', 'error');
    }
}

/**
 * Handle registration form submission
 */
async function handleRegister(e) {
    e.preventDefault();

    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    const confirmPassword = document.getElementById('register-confirm-password').value;

    // Validate passwords match
    if (password !== confirmPassword) {
        showToast('Passwords do not match', 'error');
        return;
    }

    try {
        showToast('Creating account...', 'info');
        await register(email, password, name);
        showToast('Registration successful!', 'success');

        // Reload page to show dashboard
        window.location.reload();
    } catch (error) {
        showToast(error.message || 'Registration failed', 'error');
    }
}

/**
 * Handle logout
 */
function handleLogout() {
    if (confirm('Are you sure you want to logout?')) {
        // Disconnect all active WebSockets
        disconnectAllWebSockets();
        activeJobs.clear();

        logout();
        showToast('Logged out successfully', 'success');
        window.location.reload();
    }
}

/**
 * Switch between login and register forms
 */
function switchToRegister() {
    document.getElementById('login-form-container').classList.remove('active');
    document.getElementById('register-form-container').classList.add('active');
}

function switchToLogin() {
    document.getElementById('register-form-container').classList.remove('active');
    document.getElementById('login-form-container').classList.add('active');
}

/**
 * Switch tabs
 */
function switchTab(tabName) {
    // Update tab buttons
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.toggle('active', btn.dataset.tab === tabName);
    });

    // Update tab content
    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.toggle('active', content.id === tabName);
    });

    currentPage = tabName;

    // Load data for specific tabs
    if (tabName === 'sessions') {
        loadSessions();
    } else if (tabName === 'documents') {
        loadDocuments();
    } else if (tabName === 'active-research') {
        refreshActiveJobs();
    }
}

/**
 * Handle research form submission
 */
async function handleResearchSubmit(e) {
    e.preventDefault();

    const query = document.getElementById('research-query').value;
    const depth = document.getElementById('research-depth').value;
    const filesInput = document.getElementById('research-files');
    const files = filesInput.files ? Array.from(filesInput.files) : [];

    // Get options
    const options = {
        use_web: document.querySelector('input[name="use_web"]').checked,
        use_wikipedia: document.querySelector('input[name="use_wikipedia"]').checked,
        fact_check: document.querySelector('input[name="fact_check"]').checked
    };

    // Show loading state
    const submitBtn = e.target.querySelector('button[type="submit"]');
    const submitText = document.getElementById('submit-text');
    const submitLoading = document.getElementById('submit-loading');
    submitBtn.disabled = true;
    submitText.style.display = 'none';
    submitLoading.style.display = 'inline-block';

    try {
        const result = await startResearch(query, depth, files, options);

        showToast('Research started successfully!', 'success');

        // Reset form
        e.target.reset();

        // Switch to active research tab
        switchTab('active-research');

        // Start tracking the job
        trackResearchJob(result.job_id, query);

    } catch (error) {
        showToast(error.message || 'Failed to start research', 'error');
    } finally {
        // Reset button state
        submitBtn.disabled = false;
        submitText.style.display = 'inline';
        submitLoading.style.display = 'none';
    }
}

/**
 * Track research job with WebSocket
 */
function trackResearchJob(jobId, query) {
    const token = getToken();

    // Create job card
    const jobCard = createJobCard(jobId, query);
    addJobCardToUI(jobCard);

    // Connect WebSocket or start polling
    const connection = createResearchConnection(jobId, token, {
        onProgress: (data) => updateJobProgress(jobId, data),
        onStepStarted: (data) => addJobLog(jobId, `⟳ ${data.step}...`, 'pending'),
        onStepCompleted: (data) => addJobLog(jobId, `✓ ${data.step} (${data.duration})`, 'success'),
        onComplete: (data) => handleJobComplete(jobId, data),
        onError: (data) => handleJobError(jobId, data),
        onMessage: (data) => {
            if (data.type === 'tool') {
                addJobLog(jobId, `🔧 ${data.tool}: ${data.status}`, 'info');
            }
        }
    });

    // Store connection
    activeJobs.set(jobId, {
        connection,
        query,
        startTime: Date.now()
    });

    // Update elapsed time
    startElapsedTimer(jobId);
}

/**
 * Create job card element
 */
function createJobCard(jobId, query) {
    const card = document.createElement('div');
    card.className = 'research-card';
    card.dataset.jobId = jobId;
    card.innerHTML = `
        <h3 class="query">${escapeHtml(query)}</h3>
        <div class="progress-bar">
            <div class="progress" style="width: 0%"></div>
        </div>
        <p class="status">Initializing research...</p>
        <div class="logs"></div>
        <small class="elapsed-time">Elapsed: 00:00</small>
    `;
    return card;
}

/**
 * Add job card to UI
 */
function addJobCardToUI(card) {
    const container = document.getElementById('active-jobs-container');

    // Remove empty state if present
    const emptyState = container.querySelector('.empty-state');
    if (emptyState) {
        emptyState.remove();
    }

    container.insertBefore(card, container.firstChild);
}

/**
 * Update job progress
 */
function updateJobProgress(jobId, data) {
    const card = document.querySelector(`.research-card[data-job-id="${jobId}"]`);
    if (!card) return;

    const progressBar = card.querySelector('.progress');
    const statusText = card.querySelector('.status');

    if (progressBar) {
        progressBar.style.width = `${data.percentage}%`;
    }

    if (statusText) {
        statusText.textContent = data.message || `Progress: ${data.percentage}%`;
    }
}

/**
 * Add log entry to job card
 */
function addJobLog(jobId, message, type = 'info') {
    const card = document.querySelector(`.research-card[data-job-id="${jobId}"]`);
    if (!card) return;

    const logsContainer = card.querySelector('.logs');
    if (!logsContainer) return;

    const logEntry = document.createElement('div');
    logEntry.className = `log log-${type}`;
    logEntry.textContent = message;
    logsContainer.appendChild(logEntry);

    // Auto-scroll to bottom
    logsContainer.scrollTop = logsContainer.scrollHeight;

    // Limit log entries
    const logs = logsContainer.querySelectorAll('.log');
    if (logs.length > 20) {
        logs[0].remove();
    }
}

/**
 * Handle job completion
 */
function handleJobComplete(jobId, data) {
    const card = document.querySelector(`.research-card[data-job-id="${jobId}"]`);
    if (card) {
        card.style.borderLeftColor = 'var(--success)';
        const statusText = card.querySelector('.status');
        if (statusText) {
            statusText.textContent = '✓ Research completed!';
            statusText.style.color = 'var(--success)';
        }

        // Add view results button
        const viewBtn = document.createElement('button');
        viewBtn.className = 'btn btn-primary btn-sm';
        viewBtn.textContent = 'View Results';
        viewBtn.onclick = () => viewSessionResults(data.sessionId);
        card.appendChild(viewBtn);
    }

    // Stop elapsed timer
    stopElapsedTimer(jobId);

    // Disconnect and remove from active jobs
    const job = activeJobs.get(jobId);
    if (job && job.connection) {
        job.connection.disconnect();
    }
    activeJobs.delete(jobId);

    // Reload sessions
    loadSessions();

    showToast('Research completed successfully!', 'success');
}

/**
 * Handle job error
 */
function handleJobError(jobId, data) {
    const card = document.querySelector(`.research-card[data-job-id="${jobId}"]`);
    if (card) {
        card.style.borderLeftColor = 'var(--danger)';
        const statusText = card.querySelector('.status');
        if (statusText) {
            statusText.textContent = `✗ Error: ${data.error}`;
            statusText.style.color = 'var(--danger)';
        }
    }

    // Stop elapsed timer
    stopElapsedTimer(jobId);

    // Disconnect
    const job = activeJobs.get(jobId);
    if (job && job.connection) {
        job.connection.disconnect();
    }
    activeJobs.delete(jobId);

    showToast(`Research failed: ${data.error}`, 'error');
}

/**
 * Start elapsed time timer
 */
function startElapsedTimer(jobId) {
    const timerId = setInterval(() => {
        const job = activeJobs.get(jobId);
        if (!job) {
            clearInterval(timerId);
            return;
        }

        const elapsed = Math.floor((Date.now() - job.startTime) / 1000);
        const minutes = Math.floor(elapsed / 60);
        const seconds = elapsed % 60;
        const timeStr = `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;

        const card = document.querySelector(`.research-card[data-job-id="${jobId}"]`);
        if (card) {
            const timeEl = card.querySelector('.elapsed-time');
            if (timeEl) {
                timeEl.textContent = `Elapsed: ${timeStr}`;
            }
        }
    }, 1000);

    // Store timer ID
    const job = activeJobs.get(jobId);
    if (job) {
        job.timerId = timerId;
    }
}

/**
 * Stop elapsed timer
 */
function stopElapsedTimer(jobId) {
    const job = activeJobs.get(jobId);
    if (job && job.timerId) {
        clearInterval(job.timerId);
    }
}

/**
 * Refresh active jobs
 */
function refreshActiveJobs() {
    // This function can poll for any jobs that might have been missed
    // For now, it just updates the UI
    if (activeJobs.size === 0) {
        const container = document.getElementById('active-jobs-container');
        container.innerHTML = `
            <div class="empty-state">
                <p>No active research jobs</p>
                <small>Start a new research to see progress here</small>
            </div>
        `;
    }
}

/**
 * Load sessions
 */
async function loadSessions(page = 1) {
    try {
        const data = await listSessions(page, sessionsLimit);
        displaySessions(data.sessions || [], data.pagination);
        sessionsPage = page;
    } catch (error) {
        console.error('Failed to load sessions:', error);
        showToast('Failed to load sessions', 'error');
    }
}

/**
 * Display sessions in table
 */
function displaySessions(sessions, pagination) {
    const tbody = document.querySelector('#sessions-table tbody');

    if (sessions.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="5" style="text-align: center;">No sessions found</td>
            </tr>
        `;
        return;
    }

    tbody.innerHTML = sessions.map(session => `
        <tr>
            <td><strong>${escapeHtml(session.query || 'N/A')}</strong></td>
            <td>${formatDate(session.created_at)}</td>
            <td>${formatDuration(session.duration)}</td>
            <td><span class="status-badge ${session.status}">${session.status}</span></td>
            <td>
                <div class="action-buttons">
                    <button class="action-btn" onclick="viewSessionResults('${session.id}')">View</button>
                    <button class="action-btn" onclick="handleExportSession('${session.id}', 'markdown')">Export</button>
                    <button class="action-btn danger" onclick="handleDeleteSession('${session.id}')">Delete</button>
                </div>
            </td>
        </tr>
    `).join('');

    // Update pagination
    displayPagination(pagination);
}

/**
 * Display pagination controls
 */
function displayPagination(pagination) {
    if (!pagination) return;

    const paginationEl = document.getElementById('pagination');
    const { current_page, total_pages, has_next, has_prev } = pagination;

    let html = '';

    // Previous button
    html += `<button ${!has_prev ? 'disabled' : ''} onclick="loadSessions(${current_page - 1})">Previous</button>`;

    // Page numbers
    for (let i = 1; i <= total_pages; i++) {
        if (i === 1 || i === total_pages || (i >= current_page - 2 && i <= current_page + 2)) {
            html += `<button class="${i === current_page ? 'active' : ''}" onclick="loadSessions(${i})">${i}</button>`;
        } else if (i === current_page - 3 || i === current_page + 3) {
            html += '<span>...</span>';
        }
    }

    // Next button
    html += `<button ${!has_next ? 'disabled' : ''} onclick="loadSessions(${current_page + 1})">Next</button>`;

    paginationEl.innerHTML = html;
}

/**
 * View session results
 */
async function viewSessionResults(sessionId) {
    try {
        showToast('Loading session details...', 'info');
        const session = await getSession(sessionId);
        currentSessionId = sessionId;

        // Display in modal
        const detailsEl = document.getElementById('session-details');
        detailsEl.innerHTML = `
            <div class="session-info">
                <h4>${escapeHtml(session.query)}</h4>
                <div class="session-meta">
                    <p><strong>Status:</strong> <span class="status-badge ${session.status}">${session.status}</span></p>
                    <p><strong>Date:</strong> ${formatDate(session.created_at)}</p>
                    <p><strong>Duration:</strong> ${formatDuration(session.duration)}</p>
                </div>
            </div>
            <div class="session-results">
                <h4>Results</h4>
                <div class="results-content">
                    ${session.results ? formatResults(session.results) : '<p>No results available</p>'}
                </div>
            </div>
        `;

        openModal('session-modal');
    } catch (error) {
        console.error('Failed to load session:', error);
        showToast('Failed to load session details', 'error');
    }
}

/**
 * Format results for display
 */
function formatResults(results) {
    if (typeof results === 'string') {
        // Markdown or plain text
        return `<pre style="white-space: pre-wrap;">${escapeHtml(results)}</pre>`;
    } else if (typeof results === 'object') {
        // JSON results
        return `<pre>${escapeHtml(JSON.stringify(results, null, 2))}</pre>`;
    }
    return '<p>Unable to display results</p>';
}

/**
 * Handle export session
 */
async function handleExportSession(sessionId, format) {
    try {
        currentSessionId = sessionId;
        await exportSessionFile(format);
    } catch (error) {
        console.error('Export failed:', error);
        showToast(`Failed to export: ${error.message}`, 'error');
    }
}

/**
 * Export session (called from modal)
 */
async function exportSessionFile(format) {
    if (!currentSessionId) return;

    try {
        showToast(`Exporting as ${format}...`, 'info');
        const { blob, filename } = await exportSession(currentSessionId, format);
        downloadFile(blob, filename);
        showToast('Export successful!', 'success');
    } catch (error) {
        console.error('Export failed:', error);
        showToast('Export failed', 'error');
    }
}

/**
 * Handle delete session
 */
async function handleDeleteSession(sessionId) {
    if (!confirm('Are you sure you want to delete this session?')) {
        return;
    }

    try {
        await deleteSession(sessionId);
        showToast('Session deleted successfully', 'success');
        loadSessions(sessionsPage);
    } catch (error) {
        console.error('Delete failed:', error);
        showToast('Failed to delete session', 'error');
    }
}

/**
 * Handle sessions search
 */
async function handleSessionsSearch(e) {
    const searchTerm = e.target.value.trim();

    if (!searchTerm) {
        loadSessions();
        return;
    }

    try {
        const sessions = await searchSessions(searchTerm);
        displaySessions(sessions, null);
    } catch (error) {
        console.error('Search failed:', error);
        showToast('Search failed', 'error');
    }
}

/**
 * Load documents
 */
async function loadDocuments() {
    try {
        const documents = await listDocuments();
        displayDocuments(documents);
    } catch (error) {
        console.error('Failed to load documents:', error);
        showToast('Failed to load documents', 'error');
    }
}

/**
 * Display documents grid
 */
function displayDocuments(documents) {
    const grid = document.getElementById('documents-grid');

    if (documents.length === 0) {
        grid.innerHTML = `
            <div class="empty-state">
                <p>No documents uploaded</p>
                <small>Upload PDF or DOCX files to use in your research</small>
            </div>
        `;
        return;
    }

    grid.innerHTML = documents.map(doc => `
        <div class="document-card">
            <h4>${escapeHtml(doc.filename)}</h4>
            <div class="document-meta">
                <p>Size: ${formatFileSize(doc.size)}</p>
                <p>Uploaded: ${formatDate(doc.uploaded_at)}</p>
            </div>
            <div class="document-actions">
                <button class="btn btn-secondary btn-sm" onclick="viewDocumentInfo('${doc.id}')">Info</button>
                <button class="btn btn-danger btn-sm" onclick="handleDeleteDocument('${doc.id}')">Delete</button>
            </div>
        </div>
    `).join('');
}

/**
 * Show upload modal
 */
function showUploadModal() {
    openModal('upload-modal');
}

/**
 * Upload documents
 */
async function uploadDocuments() {
    const fileInput = document.getElementById('upload-files-input');
    const files = Array.from(fileInput.files);

    if (files.length === 0) {
        showToast('Please select files to upload', 'warning');
        return;
    }

    try {
        showToast('Uploading files...', 'info');
        await uploadFiles(files);
        showToast('Files uploaded successfully!', 'success');
        closeModal('upload-modal');
        loadDocuments();
        fileInput.value = '';
        document.getElementById('file-preview').innerHTML = '';
    } catch (error) {
        console.error('Upload failed:', error);
        showToast('Upload failed', 'error');
    }
}

/**
 * Handle file preview for research form
 */
function handleFilePreview(e) {
    const files = Array.from(e.target.files);
    // Could show file names/sizes here if desired
}

/**
 * Handle file preview for upload modal
 */
function handleUploadFilePreview(e) {
    const files = Array.from(e.target.files);
    const preview = document.getElementById('file-preview');

    if (files.length === 0) {
        preview.innerHTML = '';
        return;
    }

    preview.innerHTML = `
        <h4>Selected Files:</h4>
        ${files.map(file => `
            <div class="file-preview-item">
                <span>${escapeHtml(file.name)}</span>
                <small>${formatFileSize(file.size)}</small>
            </div>
        `).join('')}
    `;
}

/**
 * View document info
 */
async function viewDocumentInfo(documentId) {
    try {
        const doc = await getDocument(documentId);
        alert(`Document: ${doc.filename}\nSize: ${formatFileSize(doc.size)}\nType: ${doc.type}\nUploaded: ${formatDate(doc.uploaded_at)}`);
    } catch (error) {
        showToast('Failed to load document info', 'error');
    }
}

/**
 * Handle delete document
 */
async function handleDeleteDocument(documentId) {
    if (!confirm('Are you sure you want to delete this document?')) {
        return;
    }

    try {
        await deleteDocument(documentId);
        showToast('Document deleted successfully', 'success');
        loadDocuments();
    } catch (error) {
        console.error('Delete failed:', error);
        showToast('Failed to delete document', 'error');
    }
}

/**
 * Modal functions
 */
function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.add('active');
    }
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.remove('active');
    }
}

// Close modal when clicking outside
window.addEventListener('click', (e) => {
    if (e.target.classList.contains('modal')) {
        e.target.classList.remove('active');
    }
});

/**
 * Show toast notification
 */
function showToast(message, type = 'info') {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.innerHTML = `
        <span class="toast-message">${escapeHtml(message)}</span>
        <button class="toast-close" onclick="this.parentElement.remove()">&times;</button>
    `;

    container.appendChild(toast);

    // Auto-remove after 5 seconds
    setTimeout(() => {
        if (toast.parentElement) {
            toast.remove();
        }
    }, 5000);
}

/**
 * Utility functions
 */

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateStr) {
    if (!dateStr) return 'N/A';
    const date = new Date(dateStr);
    return date.toLocaleString();
}

function formatDuration(seconds) {
    if (!seconds) return 'N/A';
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}m ${secs}s`;
}

function formatFileSize(bytes) {
    if (!bytes) return 'N/A';
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(2)} ${sizes[i]}`;
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
    disconnectAllWebSockets();
});
