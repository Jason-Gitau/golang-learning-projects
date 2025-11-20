# Docker Setup Guide - Deep Research Agent

Complete guide for running the Deep Research Agent in Docker with monitoring.

## 📦 What's Included

This Docker setup includes:
- **Research Agent**: Main application (API + WebSocket + Web UI)
- **Prometheus**: Metrics collection and monitoring
- **Grafana**: Beautiful dashboards for visualization

## 🚀 Quick Start

### 1. Prerequisites

Make sure you have installed:
- Docker Desktop (or Docker Engine + Docker Compose)
- At least 2GB free disk space
- Ports available: 8080, 3000, 9090

Check your installation:
```bash
docker --version
docker-compose --version
```

### 2. Setup Environment Variables

```bash
# Copy the example environment file
cp .env.example .env

# Generate a secure JWT secret
openssl rand -base64 32

# Edit .env and replace JWT_SECRET with the generated value
nano .env  # or use your favorite editor
```

### 3. Start the Stack

```bash
# Build and start all services
docker-compose up -d

# Watch the logs
docker-compose logs -f research-agent
```

### 4. Access the Services

Once all containers are running:

| Service | URL | Description |
|---------|-----|-------------|
| **Research Agent** | http://localhost:8080 | Main API + Web UI |
| **API Health** | http://localhost:8080/health | Health check |
| **Prometheus** | http://localhost:9090 | Metrics database |
| **Grafana** | http://localhost:3000 | Dashboards (admin/admin) |

## 📁 Directory Structure

After first run, your project will look like this:

```
deep-research-agent/
├── Dockerfile                 # Application container definition
├── docker-compose.yml         # Multi-container orchestration
├── .dockerignore             # Files to exclude from build
├── .env                      # Your environment variables
│
├── docker-data/              # Persistent data (AUTO-CREATED)
│   ├── database/             # SQLite database
│   ├── uploads/              # Uploaded PDF/DOCX files
│   ├── exports/              # Generated reports
│   └── config/               # Configuration files
│
└── monitoring/               # Monitoring configuration
    ├── prometheus.yml        # Prometheus config
    └── grafana/
        ├── provisioning/     # Auto-provisioning configs
        └── dashboards/       # Dashboard definitions
```

## 🔧 Common Commands

### Starting and Stopping

```bash
# Start all services
docker-compose up -d

# Start with rebuild
docker-compose up -d --build

# Stop all services
docker-compose down

# Stop and remove volumes (⚠️ deletes data!)
docker-compose down -v
```

### Viewing Logs

```bash
# View all logs
docker-compose logs

# Follow logs in real-time
docker-compose logs -f

# View specific service logs
docker-compose logs -f research-agent
docker-compose logs -f prometheus
docker-compose logs -f grafana
```

### Service Management

```bash
# Restart a service
docker-compose restart research-agent

# Stop a specific service
docker-compose stop research-agent

# Start a specific service
docker-compose start research-agent

# View service status
docker-compose ps
```

### Troubleshooting

```bash
# Check container health
docker-compose ps

# Enter research-agent container
docker-compose exec research-agent sh

# View container resource usage
docker stats

# Rebuild from scratch
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## 📊 Using the Monitoring Stack

### Prometheus (Metrics Database)

**Access**: http://localhost:9090

**What you can do:**
- View raw metrics data
- Run queries (PromQL)
- Check target health

**Example queries:**
```promql
# Check if research agent is up
up{job="research-agent"}

# Request rate
rate(http_requests_total[5m])
```

### Grafana (Visualization)

**Access**: http://localhost:3000
**Default credentials**: admin / admin

**First-time setup:**
1. Login with admin/admin
2. Change password (or skip)
3. Navigate to "Dashboards" → "Deep Research Agent - Overview"

**What you can see:**
- Application health status
- Request rates and latencies
- System resource usage
- Custom metrics (when implemented)

## 🧪 Testing the Setup

### 1. Health Check

```bash
curl http://localhost:8080/health
# Should return: {"status":"ok","timestamp":"..."}
```

### 2. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123",
    "name": "Test User"
  }'

# Save the token from the response
```

### 3. Start Research

```bash
# Replace YOUR_TOKEN with the token from registration
curl -X POST http://localhost:8080/api/v1/research/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "query": "Go programming language",
    "depth": "shallow",
    "max_sources": 3
  }'
```

### 4. Use the Web UI

Open http://localhost:8080 in your browser and interact with the built-in UI.

## 🔒 Security Notes

### For Development

✅ Current setup is fine:
- Simple JWT secret (change in .env)
- Admin credentials in plain text
- Anonymous Grafana access enabled

### For Production

⚠️ You MUST:
- [ ] Set strong JWT_SECRET (min 32 characters)
- [ ] Change Grafana admin password
- [ ] Disable Grafana anonymous access
- [ ] Use HTTPS (add nginx with SSL)
- [ ] Enable firewall rules
- [ ] Use Docker secrets instead of .env
- [ ] Regular backups of docker-data/

## 💾 Data Backup

Your important data is in `docker-data/`:

```bash
# Backup data directory
tar -czf research-agent-backup-$(date +%Y%m%d).tar.gz docker-data/

# Restore from backup
tar -xzf research-agent-backup-20240101.tar.gz
```

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Change port in docker-compose.yml
ports:
  - "8081:8080"  # Use 8081 instead
```

### Container Won't Start

```bash
# Check logs
docker-compose logs research-agent

# Common issues:
# 1. JWT_SECRET not set → Check .env file
# 2. Port conflict → Change port mapping
# 3. Build failed → Run with --build flag
```

### Can't Connect to API

```bash
# 1. Check if container is running
docker-compose ps

# 2. Check container logs
docker-compose logs research-agent

# 3. Test from inside container
docker-compose exec research-agent curl http://localhost:8080/health
```

### Database Errors

```bash
# If database is corrupted:
# ⚠️ This will delete all data!
rm -rf docker-data/database/*
docker-compose restart research-agent
```

### Grafana Not Showing Data

```bash
# 1. Check Prometheus is scraping
# Visit http://localhost:9090/targets

# 2. Verify datasource in Grafana
# Configuration → Data Sources → Prometheus

# 3. Restart Grafana
docker-compose restart grafana
```

## 📈 Performance Tuning

### Allocate More Resources

Edit `docker-compose.yml`:

```yaml
services:
  research-agent:
    # ... existing config ...
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

### Increase Worker Pool

Edit `.env`:
```bash
RESEARCH_AGENT_MAX_WORKERS=20
RESEARCH_AGENT_CONCURRENT_TOOLS=10
```

Then restart:
```bash
docker-compose restart research-agent
```

## 🔄 Updating the Application

```bash
# Pull latest code changes
git pull

# Rebuild and restart
docker-compose down
docker-compose up -d --build

# Check if running properly
docker-compose logs -f research-agent
```

## 🧹 Cleanup

### Remove Everything (including data)

```bash
# Stop and remove containers, networks, volumes
docker-compose down -v

# Remove built images
docker rmi deep-research-agent_research-agent

# Remove data directory (⚠️ PERMANENT!)
rm -rf docker-data/
```

### Remove Only Containers (keep data)

```bash
# Stop and remove containers
docker-compose down

# Data in docker-data/ is preserved
```

## 📚 Next Steps

### Learn More

- **Prometheus queries**: https://prometheus.io/docs/prometheus/latest/querying/basics/
- **Grafana dashboards**: https://grafana.com/docs/grafana/latest/dashboards/
- **Docker compose**: https://docs.docker.com/compose/

### Extend the Setup

- Add **cAdvisor** for container metrics
- Add **Node Exporter** for system metrics
- Add **nginx** for SSL termination
- Add **Redis** for caching
- Set up **automated backups**

### Connect Your UI

The Research Agent API is accessible at `http://localhost:8080/api/v1`

See the API documentation:
- Endpoints: Check `api/routes.go`
- Integration guide: `web/INTEGRATION_GUIDE.md`

## 💡 Tips

1. **Development workflow**: Use `docker-compose logs -f` in a separate terminal
2. **Quick rebuild**: `docker-compose up -d --build research-agent` (only rebuilds one service)
3. **Shell access**: `docker-compose exec research-agent sh` for debugging
4. **Environment changes**: After editing `.env`, run `docker-compose up -d` to apply
5. **Monitor resources**: Use `docker stats` to watch CPU/memory usage

## 🎓 Learning Opportunities

Since you're doing this for learning, explore:

1. **Prometheus Queries** (http://localhost:9090)
   - Try different PromQL queries
   - Understand metric types (counter, gauge, histogram)

2. **Grafana Dashboards** (http://localhost:3000)
   - Customize the provided dashboard
   - Create new panels and visualizations
   - Set up alerts

3. **Container Internals**
   - `docker-compose exec research-agent sh` - Explore the container
   - Check running processes: `ps aux`
   - View resource usage: `top`

4. **Logs and Debugging**
   - Watch application logs in real-time
   - Understand container lifecycle
   - Practice troubleshooting

## ✅ Checklist

Before considering this production-ready:

- [ ] Changed JWT_SECRET to secure value
- [ ] Changed Grafana admin password
- [ ] Tested all API endpoints
- [ ] Verified data persistence (restart container)
- [ ] Set up automated backups
- [ ] Documented any custom changes
- [ ] Tested failure scenarios
- [ ] Monitored resource usage under load

---

**Questions?** Check the main README.md or open an issue.

**Happy Containerizing! 🐳**
