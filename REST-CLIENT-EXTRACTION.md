# REST Client Extraction

## Overview

The Instana REST client code has been extracted from `terraform-provider-instana` into a separate repository to enable code reuse across multiple projects.

**New Repository**: https://github.com/instana/instana-go-client

## What Changed

### Before (Old Architecture)
```
terraform-provider-instana/
├── internal/
│   └── restapi/          # ❌ Removed - All REST client code
│       ├── Instana-api.go
│       ├── rest-client.go
│       ├── alerting-channels-api.go
│       └── ... (70+ files)
```

### After (New Architecture)
```
terraform-provider-instana/
├── go.mod                # ✅ Added: github.com/instana/instana-go-client v1.0.0
└── internal/
    ├── provider/         # ✅ Updated: Uses client.InstanaAPI
    ├── resources/        # ✅ Updated: Uses api.* types
    └── shared/           # ✅ Updated: Uses shared/types

External Dependency:
github.com/instana/instana-go-client/
├── client/               # REST client implementation
├── api/                  # API models (AlertingChannel, etc.)
├── config/               # Client configuration
└── shared/               # Shared types and utilities
```

## Key Changes

### 1. Import Changes
**Old:**
```go
import "github.com/instana/terraform-provider-instana/internal/restapi"
```

**New:**
```go
import (
    "github.com/instana/instana-go-client/client"
    "github.com/instana/instana-go-client/api"
    "github.com/instana/instana-go-client/shared/types"
)
```

### 2. Dependency Management
**go.mod:**
```go
require (
    github.com/instana/instana-go-client v1.0.0  // New dependency
)
```

## Development Workflow

### Adding a New Resource

**Step 1: Update instana-go-client**
```bash
# 1. Clone the client repository
git clone https://github.com/instana/instana-go-client.git

# 2. Add new API model (e.g., api/newresource.go)
# 3. Add REST resource implementation
# 4. Add tests
# 5. Commit and create PR
# 6. After merge, create a new release (e.g., v1.1.0)
```

**Step 2: Update terraform-provider-instana**
```bash
# 1. Update dependency version
go get github.com/instana/instana-go-client@v1.1.0
go mod tidy

# 2. Create resource implementation using new types
# internal/resources/newresource/resource-newresource.go

# 3. Register resource in provider.go
# 4. Add tests
# 5. Commit and create PR
```

### Version Update Process

```bash
# Update to specific version
go get github.com/instana/instana-go-client@v1.2.0

# Update to latest
go get -u github.com/instana/instana-go-client

# Verify and clean up
go mod tidy
go mod verify
```
