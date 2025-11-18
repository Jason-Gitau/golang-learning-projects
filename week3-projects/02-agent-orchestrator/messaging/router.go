package messaging

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang-learning/agent-orchestrator/models"
)

// Router handles message routing between agents and requesters
type Router struct {
	requestChan  chan *models.Request
	responseChan chan *models.Response
	pendingReqs  map[string]chan *models.Response // requestID -> response channel
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewRouter creates a new message router
func NewRouter(ctx context.Context, requestQueueSize int) *Router {
	ctx, cancel := context.WithCancel(ctx)
	return &Router{
		requestChan:  make(chan *models.Request, requestQueueSize),
		responseChan: make(chan *models.Response, requestQueueSize),
		pendingReqs:  make(map[string]chan *models.Response),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start starts the router
func (r *Router) Start() {
	r.wg.Add(1)
	go r.routeResponses()
}

// Stop stops the router
func (r *Router) Stop() {
	r.cancel()
	r.wg.Wait()
	close(r.requestChan)
	close(r.responseChan)
}

// SubmitRequest submits a request and waits for response
func (r *Router) SubmitRequest(req *models.Request) (*models.Response, error) {
	// Create a response channel for this request
	respChan := make(chan *models.Response, 1)

	r.mu.Lock()
	r.pendingReqs[req.ID] = respChan
	r.mu.Unlock()

	// Send request to queue
	select {
	case r.requestChan <- req:
		// Request queued successfully
	case <-r.ctx.Done():
		return nil, fmt.Errorf("router shutting down")
	}

	// Wait for response with timeout
	ctx, cancel := context.WithTimeout(r.ctx, req.Timeout)
	defer cancel()

	select {
	case resp := <-respChan:
		r.mu.Lock()
		delete(r.pendingReqs, req.ID)
		r.mu.Unlock()
		return resp, nil
	case <-ctx.Done():
		r.mu.Lock()
		delete(r.pendingReqs, req.ID)
		r.mu.Unlock()
		return nil, fmt.Errorf("request timeout")
	}
}

// GetRequestChannel returns the request channel for agents to consume
func (r *Router) GetRequestChannel() <-chan *models.Request {
	return r.requestChan
}

// SendResponse sends a response back to the router
func (r *Router) SendResponse(resp *models.Response) {
	select {
	case r.responseChan <- resp:
		// Response sent
	case <-r.ctx.Done():
		// Router shutting down
	}
}

// routeResponses routes responses back to waiting requesters
func (r *Router) routeResponses() {
	defer r.wg.Done()

	for {
		select {
		case resp := <-r.responseChan:
			r.mu.RLock()
			respChan, exists := r.pendingReqs[resp.RequestID]
			r.mu.RUnlock()

			if exists {
				select {
				case respChan <- resp:
					// Response delivered
				default:
					// Response channel full or closed, skip
				}
			}
		case <-r.ctx.Done():
			return
		}
	}
}

// PendingRequests returns the number of pending requests
func (r *Router) PendingRequests() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.pendingReqs)
}
