/**
 * WebSocket Module
 * Handles real-time updates for research progress
 */

// Track active WebSocket connections
if (!window.activeWebSockets) {
    window.activeWebSockets = [];
}

/**
 * WebSocket connection manager for a research job
 */
class ResearchWebSocket {
    constructor(jobId, token) {
        this.jobId = jobId;
        this.token = token;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000; // Start with 1 second
        this.callbacks = {
            onMessage: [],
            onProgress: [],
            onStepStarted: [],
            onStepCompleted: [],
            onComplete: [],
            onError: [],
            onClose: []
        };
    }

    /**
     * Connect to WebSocket
     */
    connect() {
        try {
            // Determine WebSocket protocol based on current page protocol
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const host = window.location.host;
            const url = `${protocol}//${host}/api/v1/research/${this.jobId}/stream?token=${this.token}`;

            this.ws = new WebSocket(url);

            this.ws.onopen = () => {
                console.log(`WebSocket connected for job ${this.jobId}`);
                this.reconnectAttempts = 0;
                this.reconnectDelay = 1000;
            };

            this.ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    this.handleMessage(data);
                } catch (error) {
                    console.error('Error parsing WebSocket message:', error);
                }
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.trigger('onError', { error: 'WebSocket connection error' });
            };

            this.ws.onclose = (event) => {
                console.log(`WebSocket closed for job ${this.jobId}`, event.code, event.reason);
                this.trigger('onClose', { code: event.code, reason: event.reason });

                // Attempt to reconnect if not a normal closure
                if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
                    this.reconnect();
                }
            };

            // Track this connection
            window.activeWebSockets.push(this);

        } catch (error) {
            console.error('Failed to connect WebSocket:', error);
            this.trigger('onError', { error: 'Failed to establish connection' });
        }
    }

    /**
     * Handle incoming message
     * @param {Object} data - Message data
     */
    handleMessage(data) {
        // Trigger generic message callback
        this.trigger('onMessage', data);

        // Handle specific message types
        switch (data.type) {
            case 'progress':
                this.trigger('onProgress', {
                    percentage: data.data.percentage,
                    current: data.data.current,
                    total: data.data.total,
                    message: data.data.message
                });
                break;

            case 'step_started':
                this.trigger('onStepStarted', {
                    step: data.data.step,
                    stepNumber: data.data.step_number,
                    totalSteps: data.data.total_steps
                });
                break;

            case 'step_completed':
                this.trigger('onStepCompleted', {
                    step: data.data.step,
                    duration: data.data.duration,
                    result: data.data.result
                });
                break;

            case 'tool_execution':
                this.trigger('onMessage', {
                    type: 'tool',
                    tool: data.data.tool_name,
                    status: data.data.status,
                    duration: data.data.duration
                });
                break;

            case 'research_completed':
                this.trigger('onComplete', {
                    sessionId: data.data.session_id,
                    results: data.data.results,
                    duration: data.data.total_duration
                });
                // Close connection after completion
                this.disconnect();
                break;

            case 'error':
                this.trigger('onError', {
                    error: data.data.message,
                    details: data.data.details
                });
                break;

            default:
                console.log('Unknown message type:', data.type);
        }
    }

    /**
     * Register callback for message type
     * @param {string} event - Event name
     * @param {Function} callback - Callback function
     */
    on(event, callback) {
        if (this.callbacks[event]) {
            this.callbacks[event].push(callback);
        }
    }

    /**
     * Trigger callbacks for event
     * @param {string} event - Event name
     * @param {Object} data - Event data
     */
    trigger(event, data) {
        if (this.callbacks[event]) {
            this.callbacks[event].forEach(callback => {
                try {
                    callback(data);
                } catch (error) {
                    console.error(`Error in ${event} callback:`, error);
                }
            });
        }
    }

    /**
     * Reconnect with exponential backoff
     */
    reconnect() {
        this.reconnectAttempts++;
        const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);

        console.log(`Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

        setTimeout(() => {
            if (this.ws.readyState === WebSocket.CLOSED) {
                this.connect();
            }
        }, delay);
    }

    /**
     * Send message through WebSocket
     * @param {Object} message - Message to send
     */
    send(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.error('WebSocket is not open');
        }
    }

    /**
     * Disconnect WebSocket
     */
    disconnect() {
        if (this.ws) {
            this.ws.close(1000, 'Client disconnect');
            // Remove from active connections
            const index = window.activeWebSockets.indexOf(this);
            if (index > -1) {
                window.activeWebSockets.splice(index, 1);
            }
        }
    }

    /**
     * Check if WebSocket is connected
     * @returns {boolean}
     */
    isConnected() {
        return this.ws && this.ws.readyState === WebSocket.OPEN;
    }
}

/**
 * Create and connect WebSocket for a research job
 * @param {string} jobId - Research job ID
 * @param {string} token - Authentication token
 * @returns {ResearchWebSocket} - WebSocket instance
 */
function connectResearchWebSocket(jobId, token) {
    const ws = new ResearchWebSocket(jobId, token);
    ws.connect();
    return ws;
}

/**
 * Disconnect all active WebSockets
 */
function disconnectAllWebSockets() {
    window.activeWebSockets.forEach(ws => {
        ws.disconnect();
    });
    window.activeWebSockets = [];
}

/**
 * Fallback polling function for when WebSocket is not available
 * @param {string} jobId - Job ID to poll
 * @param {Function} callback - Callback for status updates
 * @param {number} interval - Polling interval in ms (default: 5000)
 * @returns {number} - Interval ID
 */
function pollResearchStatus(jobId, callback, interval = 5000) {
    const pollInterval = setInterval(async () => {
        try {
            const status = await getResearchStatus(jobId);
            callback(status);

            // Stop polling if completed or failed
            if (status.status === 'completed' || status.status === 'failed') {
                clearInterval(pollInterval);
            }
        } catch (error) {
            console.error('Polling error:', error);
            // Continue polling even on error
        }
    }, interval);

    return pollInterval;
}

/**
 * Create WebSocket or fallback to polling
 * @param {string} jobId - Job ID
 * @param {string} token - Auth token
 * @param {Object} callbacks - Callback functions
 * @returns {Object} - Connection object with disconnect method
 */
function createResearchConnection(jobId, token, callbacks = {}) {
    // Try WebSocket first
    if ('WebSocket' in window) {
        const ws = connectResearchWebSocket(jobId, token);

        // Register callbacks
        Object.keys(callbacks).forEach(event => {
            if (callbacks[event]) {
                ws.on(event, callbacks[event]);
            }
        });

        return {
            type: 'websocket',
            connection: ws,
            disconnect: () => ws.disconnect()
        };
    } else {
        // Fallback to polling
        console.warn('WebSocket not supported, falling back to polling');
        const intervalId = pollResearchStatus(jobId, (status) => {
            if (callbacks.onProgress) {
                callbacks.onProgress({
                    percentage: status.progress || 0,
                    message: status.current_step || 'Processing...'
                });
            }

            if (status.status === 'completed' && callbacks.onComplete) {
                callbacks.onComplete({
                    sessionId: status.session_id,
                    results: status.results
                });
            }

            if (status.status === 'failed' && callbacks.onError) {
                callbacks.onError({
                    error: status.error || 'Research failed'
                });
            }
        });

        return {
            type: 'polling',
            intervalId,
            disconnect: () => clearInterval(intervalId)
        };
    }
}

// Export functions
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        ResearchWebSocket,
        connectResearchWebSocket,
        disconnectAllWebSockets,
        pollResearchStatus,
        createResearchConnection
    };
}
