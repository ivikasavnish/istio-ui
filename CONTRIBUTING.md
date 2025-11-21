# Contributing to MeshControl Center

Thank you for your interest in contributing to MeshControl Center! This document provides guidelines and instructions for contributing.

## Code of Conduct

This project follows a standard code of conduct. Be respectful, constructive, and collaborative.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- kubectl configured with cluster access
- Git

### Setting Up Development Environment

1. Fork and clone the repository:
```bash
git clone https://github.com/ivikasavnish/istio-ui.git
cd istio-ui
```

2. Set up the backend:
```bash
cd backend
go mod download
```

3. Set up the frontend:
```bash
cd frontend
npm install
```

## Development Workflow

### Backend Development

1. Make your changes in the `backend/` directory
2. Run the server locally:
```bash
cd backend
go run cmd/server/main.go
```

3. Test your changes:
```bash
go test ./...
```

4. Format your code:
```bash
go fmt ./...
```

### Frontend Development

1. Make your changes in the `frontend/` directory
2. Start the development server:
```bash
cd frontend
npm start
```

3. Run tests:
```bash
npm test
```

4. Check for linting issues:
```bash
npm run lint
```

## Code Style Guidelines

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors explicitly
- Use structured logging

Example:
```go
// ListVirtualServices returns all VirtualServices from the specified namespace
func ListVirtualServices(client *istio.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        namespace := c.DefaultQuery("namespace", "")
        
        vsList, err := client.ListVirtualServices(c.Request.Context(), namespace)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, vsList.Items)
    }
}
```

### React Code Style

- Use functional components with hooks
- Keep components small and reusable
- Use PropTypes or TypeScript for type safety
- Follow Material-UI patterns
- Use meaningful component and variable names

Example:
```javascript
import React, { useState, useEffect } from 'react';

export default function ResourceList({ apiClient, resourceType }) {
  const [resources, setResources] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadResources();
  }, []);

  const loadResources = async () => {
    try {
      const response = await apiClient.list();
      setResources(response.data);
    } finally {
      setLoading(false);
    }
  };

  // ... rest of component
}
```

## Submitting Changes

### Pull Request Process

1. Create a new branch for your feature:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes and commit with clear messages:
```bash
git add .
git commit -m "Add feature: description of your changes"
```

3. Push to your fork:
```bash
git push origin feature/your-feature-name
```

4. Create a Pull Request on GitHub with:
   - Clear title describing the change
   - Detailed description of what was changed and why
   - Reference to any related issues
   - Screenshots for UI changes

### Commit Message Guidelines

Follow conventional commits format:

```
type(scope): subject

body (optional)

footer (optional)
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(api): add endpoint for listing service entries

fix(ui): resolve pagination issue in virtual services table

docs(readme): update installation instructions
```

## Testing

### Backend Tests

Create tests alongside your code:

```go
// virtualservice_test.go
package api

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestListVirtualServices(t *testing.T) {
    // Test implementation
}
```

Run tests:
```bash
cd backend
go test ./...
```

### Frontend Tests

Create tests using React Testing Library:

```javascript
// VirtualServices.test.js
import { render, screen } from '@testing-library/react';
import VirtualServices from './VirtualServices';

test('renders virtual services page', () => {
  render(<VirtualServices />);
  expect(screen.getByText(/Virtual Services/i)).toBeInTheDocument();
});
```

Run tests:
```bash
cd frontend
npm test
```

## Adding New Features

### Adding a New Istio Resource Type

1. **Backend**: Add client methods in `internal/istio/client.go`
2. **Backend**: Create API handlers in `internal/api/`
3. **Backend**: Add routes in `cmd/server/main.go`
4. **Frontend**: Add API client methods in `services/api.js`
5. **Frontend**: Create page component in `pages/`
6. **Frontend**: Add route in `App.js`
7. **Frontend**: Add navigation item in `components/Layout.js`
8. **Documentation**: Update API docs

### Adding a New API Endpoint

1. Define the handler function
2. Add route to the router
3. Update API documentation
4. Add tests
5. Update frontend API client if needed

## Documentation

- Update README.md for user-facing changes
- Update API.md for API changes
- Update ARCHITECTURE.md for architectural changes
- Add inline code comments for complex logic
- Create examples for new features

## Issue Guidelines

### Reporting Bugs

Include:
- Clear description of the bug
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, Node version, etc.)
- Error messages or logs
- Screenshots if applicable

### Suggesting Features

Include:
- Clear description of the feature
- Use cases and benefits
- Proposed implementation approach
- Any potential drawbacks

## Community

- Be respectful and inclusive
- Help others in issues and discussions
- Share your use cases and feedback
- Contribute to documentation
- Review pull requests

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

If you have questions about contributing:
- Open an issue with the "question" label
- Check existing documentation
- Review closed issues for similar questions

Thank you for contributing to MeshControl Center! 🎉
