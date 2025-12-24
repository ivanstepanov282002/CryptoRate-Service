# CryptoRate-Service

A cryptocurrency rates API service with Telegram bot integration that provides real-time cryptocurrency rates and automated notifications.

## 🚀 CI/CD Pipeline

This project has a comprehensive CI/CD pipeline set up with GitHub Actions and Docker for automated testing, building, and deployment.

### 🛠️ Pipeline Features

- **Automated Testing**: Runs on every push and pull request
- **Code Quality**: Linting and security scanning
- **Docker Integration**: Multi-stage Docker builds
- **Coverage Reports**: Automated code coverage reporting
- **Multi-environment Deployment**: Staging and production deployments

### 📊 GitHub Actions Workflows

#### 1. CI/CD Pipeline (`ci-cd.yml`)
- Runs on `main` and `develop` branches
- Tests with multiple Go versions (1.21, 1.22)
- Security scanning
- Docker image building and pushing to DockerHub
- Production deployment

#### 2. Tests (`test.yml`)
- Runs on all pull requests
- Linting and vetting
- Unit and integration tests
- Coverage reporting

#### 3. Manual Deploy (`deploy.yml`)
- Manual deployment workflow
- Supports staging and production environments
- Configurable image tags

### 🐳 Docker Configuration

The service is containerized with separate images for:
- **API**: REST API server
- **Bot**: Telegram bot
- **Worker**: Background job processor

### 🔧 Build & Test Commands

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Generate HTML coverage report
make test-cover-html

# Run benchmarks
make bench

# Build the project
make build

# Lint the code
make lint

# Format the code
make fmt

# Run in Docker
make docker-up

# Stop Docker containers
make docker-down
```

### 🚀 Deployment

#### Production Deployment
```bash
# Using Docker Compose
make docker-up-prod

# Or using the manual workflow in GitHub
```

#### Environment Variables
Create a `.env` file with the following variables:
```env
POSTGRES_USER=crypto_user
POSTGRES_PASSWORD=secure_password_123
POSTGRES_DB=crypto_db
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
DOCKERHUB_USERNAME=your_dockerhub_username
```

### 🛡️ Security Features

- **Static Code Analysis**: Using golangci-lint
- **Security Scanning**: Using gosec
- **Dependency Scanning**: Automated dependency vulnerability checks
- **Container Security**: Multi-stage builds with non-root users

### 📈 Code Quality

- **Test Coverage**: >90% for critical components
- **Code Review**: Required for all pull requests
- **Automated Linting**: Enforces code standards
- **Security Scanning**: Automated security checks

### 🚦 Deployment Environments

1. **Development**: Automated on `develop` branch
2. **Staging**: Manual deployment with review
3. **Production**: Manual deployment with approval

### 🔐 Required Secrets

For the CI/CD pipeline to work, you need to configure these secrets in GitHub:

- `DOCKERHUB_USERNAME`: Your DockerHub username
- `DOCKERHUB_TOKEN`: Your DockerHub access token
- `TELEGRAM_BOT_TOKEN`: Telegram bot token (for bot service)

### 📦 Docker Compose

The project includes Docker Compose files for different environments:

- `docker-compose.yml`: Development environment
- `docker-compose.prod.yml`: Production environment

### 🧪 Testing Strategy

- **Unit Tests**: Core logic and business rules
- **Integration Tests**: Database and API interactions
- **End-to-End Tests**: Full system testing
- **Performance Tests**: Benchmarks and load testing

### 📊 Monitoring & Observability

- **Health Checks**: Built-in health check endpoints
- **Logging**: Structured logging with levels
- **Metrics**: Performance and usage metrics

### 🔄 Continuous Integration

1. Code pushed to `main` or `develop`
2. Automated tests run on multiple Go versions
3. Code quality checks (linting, security)
4. Docker images built and pushed
5. Deployment to staging/production (if on main)

### 🚀 Continuous Deployment

- **Staging**: Automated after successful tests
- **Production**: Manual approval required
- **Rollback**: Automated rollback on health check failure

### 🛠️ Development Workflow

1. Create feature branch
2. Write code and tests
3. Submit pull request
4. Automated tests run
5. Code review and approval
6. Merge to main
7. Automated deployment

### 📋 Pre-commit Checks

Before committing, ensure:
- All tests pass
- Code is properly formatted (`make fmt`)
- Linting passes (`make lint`)
- Coverage remains high
- Documentation is updated