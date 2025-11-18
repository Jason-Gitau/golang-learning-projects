/**
 * Authentication Module
 * Handles user registration, login, logout, and token management
 */

const AUTH_API_BASE = '/api/v1/auth';
const TOKEN_KEY = 'research_agent_token';
const USER_KEY = 'research_agent_user';

/**
 * User login
 * @param {string} email - User email
 * @param {string} password - User password
 * @returns {Promise<Object>} - Login response with token
 */
async function login(email, password) {
    try {
        const response = await fetch(`${AUTH_API_BASE}/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Login failed');
        }

        // Store token and user info
        localStorage.setItem(TOKEN_KEY, data.token);
        if (data.user) {
            localStorage.setItem(USER_KEY, JSON.stringify(data.user));
        }

        return data;
    } catch (error) {
        console.error('Login error:', error);
        throw error;
    }
}

/**
 * User registration
 * @param {string} email - User email
 * @param {string} password - User password
 * @param {string} name - User's full name
 * @returns {Promise<Object>} - Registration response
 */
async function register(email, password, name) {
    try {
        const response = await fetch(`${AUTH_API_BASE}/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password, name })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Registration failed');
        }

        // Store token and user info
        if (data.token) {
            localStorage.setItem(TOKEN_KEY, data.token);
        }
        if (data.user) {
            localStorage.setItem(USER_KEY, JSON.stringify(data.user));
        }

        return data;
    } catch (error) {
        console.error('Registration error:', error);
        throw error;
    }
}

/**
 * User logout
 * Clears local storage and redirects to login page
 */
function logout() {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    // Clear any active WebSocket connections
    if (window.activeWebSockets) {
        window.activeWebSockets.forEach(ws => ws.close());
        window.activeWebSockets = [];
    }
}

/**
 * Get stored authentication token
 * @returns {string|null} - JWT token or null
 */
function getToken() {
    return localStorage.getItem(TOKEN_KEY);
}

/**
 * Check if user is authenticated
 * @returns {boolean} - True if token exists
 */
function isAuthenticated() {
    const token = getToken();
    if (!token) return false;

    // Optional: Check token expiration if JWT
    try {
        const payload = parseJWT(token);
        if (payload.exp) {
            const now = Math.floor(Date.now() / 1000);
            return payload.exp > now;
        }
        return true;
    } catch (error) {
        return false;
    }
}

/**
 * Get current user info from storage
 * @returns {Object|null} - User object or null
 */
function getCurrentUser() {
    const userStr = localStorage.getItem(USER_KEY);
    if (!userStr) return null;

    try {
        return JSON.parse(userStr);
    } catch (error) {
        console.error('Error parsing user data:', error);
        return null;
    }
}

/**
 * Fetch current user profile from API
 * @returns {Promise<Object>} - User profile
 */
async function fetchUserProfile() {
    try {
        const token = getToken();
        if (!token) {
            throw new Error('No authentication token');
        }

        const response = await fetch(`${AUTH_API_BASE}/profile`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to fetch profile');
        }

        // Update stored user info
        localStorage.setItem(USER_KEY, JSON.stringify(data.user));

        return data.user;
    } catch (error) {
        console.error('Fetch profile error:', error);
        throw error;
    }
}

/**
 * Parse JWT token (basic parsing, doesn't verify signature)
 * @param {string} token - JWT token
 * @returns {Object} - Decoded payload
 */
function parseJWT(token) {
    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));

        return JSON.parse(jsonPayload);
    } catch (error) {
        console.error('Error parsing JWT:', error);
        return {};
    }
}

/**
 * Make authenticated API request
 * @param {string} url - API endpoint
 * @param {Object} options - Fetch options
 * @returns {Promise<Response>} - Fetch response
 */
async function authenticatedFetch(url, options = {}) {
    const token = getToken();
    if (!token) {
        throw new Error('Not authenticated');
    }

    const headers = {
        ...options.headers,
        'Authorization': `Bearer ${token}`
    };

    const response = await fetch(url, { ...options, headers });

    // Handle token expiration
    if (response.status === 401) {
        logout();
        window.location.href = '/';
        throw new Error('Session expired. Please login again.');
    }

    return response;
}

// Export functions for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        login,
        register,
        logout,
        getToken,
        isAuthenticated,
        getCurrentUser,
        fetchUserProfile,
        authenticatedFetch
    };
}
