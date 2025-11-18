/**
 * Research Operations Module
 * Handles research job creation, status tracking, session management
 */

const RESEARCH_API_BASE = '/api/v1/research';
const DOCUMENTS_API_BASE = '/api/v1/documents';

/**
 * Start a new research job
 * @param {string} query - Research query
 * @param {string} depth - Research depth (shallow, medium, deep)
 * @param {Array<File>} files - Optional files to upload
 * @param {Object} options - Research options (use_web, use_wikipedia, etc.)
 * @returns {Promise<Object>} - Job information with job_id
 */
async function startResearch(query, depth = 'medium', files = [], options = {}) {
    try {
        // First upload files if any
        let documentIds = [];
        if (files && files.length > 0) {
            const uploadResult = await uploadFiles(files);
            documentIds = uploadResult.document_ids || [];
        }

        // Start research with the query
        const response = await authenticatedFetch(`${RESEARCH_API_BASE}/start`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                query,
                depth,
                document_ids: documentIds,
                use_web: options.use_web !== false,
                use_wikipedia: options.use_wikipedia !== false,
                fact_check: options.fact_check || false,
                ...options
            })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to start research');
        }

        return data;
    } catch (error) {
        console.error('Start research error:', error);
        throw error;
    }
}

/**
 * Get research job status
 * @param {string} jobId - Job ID
 * @returns {Promise<Object>} - Job status information
 */
async function getResearchStatus(jobId) {
    try {
        const response = await authenticatedFetch(`${RESEARCH_API_BASE}/${jobId}/status`);
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to get research status');
        }

        return data;
    } catch (error) {
        console.error('Get research status error:', error);
        throw error;
    }
}

/**
 * List all research sessions
 * @param {number} page - Page number (default: 1)
 * @param {number} limit - Items per page (default: 20)
 * @param {string} status - Filter by status (optional)
 * @returns {Promise<Object>} - Sessions list with pagination
 */
async function listSessions(page = 1, limit = 20, status = null) {
    try {
        let url = `${RESEARCH_API_BASE}/sessions?page=${page}&limit=${limit}`;
        if (status) {
            url += `&status=${status}`;
        }

        const response = await authenticatedFetch(url);
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to list sessions');
        }

        return data;
    } catch (error) {
        console.error('List sessions error:', error);
        throw error;
    }
}

/**
 * Get detailed session information
 * @param {string} sessionId - Session ID
 * @returns {Promise<Object>} - Complete session details with results
 */
async function getSession(sessionId) {
    try {
        const response = await authenticatedFetch(`${RESEARCH_API_BASE}/sessions/${sessionId}`);
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to get session');
        }

        return data;
    } catch (error) {
        console.error('Get session error:', error);
        throw error;
    }
}

/**
 * Delete a research session
 * @param {string} sessionId - Session ID to delete
 * @returns {Promise<Object>} - Deletion confirmation
 */
async function deleteSession(sessionId) {
    try {
        const response = await authenticatedFetch(`${RESEARCH_API_BASE}/sessions/${sessionId}`, {
            method: 'DELETE'
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to delete session');
        }

        return data;
    } catch (error) {
        console.error('Delete session error:', error);
        throw error;
    }
}

/**
 * Export session results in specified format
 * @param {string} sessionId - Session ID
 * @param {string} format - Export format (markdown, json, pdf)
 * @returns {Promise<Blob>} - File blob for download
 */
async function exportSession(sessionId, format = 'markdown') {
    try {
        const response = await authenticatedFetch(
            `${RESEARCH_API_BASE}/sessions/${sessionId}/export?format=${format}`
        );

        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || 'Failed to export session');
        }

        // Get filename from Content-Disposition header if available
        const contentDisposition = response.headers.get('Content-Disposition');
        let filename = `research-${sessionId}.${format}`;
        if (contentDisposition) {
            const filenameMatch = contentDisposition.match(/filename="?(.+)"?/);
            if (filenameMatch) {
                filename = filenameMatch[1];
            }
        }

        const blob = await response.blob();
        return { blob, filename };
    } catch (error) {
        console.error('Export session error:', error);
        throw error;
    }
}

/**
 * Download exported file
 * @param {Blob} blob - File blob
 * @param {string} filename - Filename for download
 */
function downloadFile(blob, filename) {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
}

/**
 * Upload files/documents
 * @param {Array<File>} files - Files to upload
 * @returns {Promise<Object>} - Upload result with document IDs
 */
async function uploadFiles(files) {
    try {
        const formData = new FormData();
        for (const file of files) {
            formData.append('files', file);
        }

        const token = getToken();
        const response = await fetch(`${DOCUMENTS_API_BASE}/upload`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to upload files');
        }

        return data;
    } catch (error) {
        console.error('Upload files error:', error);
        throw error;
    }
}

/**
 * List all uploaded documents
 * @returns {Promise<Array>} - List of documents
 */
async function listDocuments() {
    try {
        const response = await authenticatedFetch(`${DOCUMENTS_API_BASE}`);
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to list documents');
        }

        return data.documents || [];
    } catch (error) {
        console.error('List documents error:', error);
        throw error;
    }
}

/**
 * Get document details
 * @param {string} documentId - Document ID
 * @returns {Promise<Object>} - Document information
 */
async function getDocument(documentId) {
    try {
        const response = await authenticatedFetch(`${DOCUMENTS_API_BASE}/${documentId}`);
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to get document');
        }

        return data;
    } catch (error) {
        console.error('Get document error:', error);
        throw error;
    }
}

/**
 * Delete a document
 * @param {string} documentId - Document ID to delete
 * @returns {Promise<Object>} - Deletion confirmation
 */
async function deleteDocument(documentId) {
    try {
        const response = await authenticatedFetch(`${DOCUMENTS_API_BASE}/${documentId}`, {
            method: 'DELETE'
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to delete document');
        }

        return data;
    } catch (error) {
        console.error('Delete document error:', error);
        throw error;
    }
}

/**
 * Search sessions
 * @param {string} searchTerm - Search query
 * @returns {Promise<Array>} - Matching sessions
 */
async function searchSessions(searchTerm) {
    try {
        const response = await authenticatedFetch(
            `${RESEARCH_API_BASE}/sessions/search?q=${encodeURIComponent(searchTerm)}`
        );
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to search sessions');
        }

        return data.sessions || [];
    } catch (error) {
        console.error('Search sessions error:', error);
        throw error;
    }
}

/**
 * Cancel an active research job
 * @param {string} jobId - Job ID to cancel
 * @returns {Promise<Object>} - Cancellation confirmation
 */
async function cancelResearch(jobId) {
    try {
        const response = await authenticatedFetch(`${RESEARCH_API_BASE}/${jobId}/cancel`, {
            method: 'POST'
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to cancel research');
        }

        return data;
    } catch (error) {
        console.error('Cancel research error:', error);
        throw error;
    }
}

// Export functions
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        startResearch,
        getResearchStatus,
        listSessions,
        getSession,
        deleteSession,
        exportSession,
        downloadFile,
        uploadFiles,
        listDocuments,
        getDocument,
        deleteDocument,
        searchSessions,
        cancelResearch
    };
}
