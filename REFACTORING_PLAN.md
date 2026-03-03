# Terraform Provider Instana - Code Refactoring Plan
## Splitting Monolithic Project into Two Distinct Projects

**Document Version:** 1.0  
**Date:** 2026-03-03  
**Target Go Client Library:** `github.com/instana/go-instana-client`  
**Current Provider:** `github.com/instana/terraform-provider-instana`

---

## Executive Summary

This document outlines a comprehensive plan to refactor the Terraform Provider Instana codebase by splitting it into two distinct, well-architected projects:

1. **Instana Go Client Library** (`github.com/instana/go-instana-client`) - A standalone, reusable Go client library for Instana API interactions
2. **Terraform Provider Instana** (existing repository) - Terraform-specific implementation that consumes the Go client library

**Key Benefits:**
- Improved code reusability across multiple projects
- Better separation of concerns
- Independent versioning and release cycles
- Enhanced testability and maintainability
- Reduced coupling between Terraform logic and API client code

**Project Duration:** 14 weeks  
**Team Size:** 5-6 people (2-3 developers, 1 DevOps, 1 QA, 1 technical writer)

---

## Table of Contents

1. [Architectural Analysis](#1-architectural-analysis)
2. [Code Separation Strategy](#2-code-separation-strategy)
3. [Execution Plan](#3-execution-plan)
4. [Go Client Library Release Strategy](#4-go-client-library-release-strategy)
5. [Integration Plan](#5-integration-plan)
6. [Deployment Strategy](#6-deployment-strategy)
7. [Quality Assurance Plan](#7-quality-assurance-plan)
8. [Documentation Requirements](#8-documentation-requirements)
9. [Risk Assessment](#9-risk-assessment)
10. [Timeline & Milestones](#10-timeline--milestones)
11. [Appendices](#appendices)

---

## 1. Architectural Analysis

### 1.1 Current Codebase Overview

The current Terraform Provider Instana is a monolithic project where API client code and Terraform-specific logic are tightly coupled within the same codebase.

**Key Statistics:**
- **Total API Files:** 50+ files in `internal/restapi/`
- **Terraform Resources:** 24 resource types
- **Data Sources:** 7 data sources
- **Lines of Code:** ~15,000+ lines in API client layer
- **Dependencies:** Mix of Terraform and HTTP client libraries

### 1.2 Current Structure

```
terraform-provider-instana/
├── main.go                          # Provider entry point
├── go.mod                           # Dependencies
├── internal/
│   ├── provider/                    # Terraform provider (220 lines)
│   ├── datasources/                 # 7 data sources
│   ├── resources/                   # 24 resources
│   ├── resourcehandle/              # Resource abstraction
│   ├── restapi/                     # API CLIENT (TO EXTRACT - 50+ files)
│   │   ├── Instana-api.go          # Main API interface (193 lines)
│   │   ├── rest-client.go          # HTTP client (200+ lines)
│   │   ├── *-api.go                # API implementations
│   │   ├── *.go                    # Data models & types
│   │   └── provider-models.go      # TF-specific (KEEP)
│   ├── shared/                      # Shared utilities
│   │   └── tagfilter/              # Tag filter (EXTRACT)
│   └── util/                        # Utilities
├── testutils/                       # Test utilities
├── mocks/                          # Generated mocks
└── docs/                           # Documentation
```

### 1.3 Components to Extract

#### Core API Client (Priority: High)
- `rest-client.go` - HTTP client with rate limiting
- `Instana-api.go` - Main API factory
- `instana-rest-resource.go` - Resource interfaces
- `default-rest-resource.go` - Generic REST operations
- `read-only-rest-resource.go` - Read-only operations
- JSON unmarshallers (3 files)

#### API Implementations by Domain (Priority: High)
| Domain | Files | Models |
|--------|-------|--------|
| Alerting | 2 | AlertingChannel, AlertingConfiguration |
| Application | 2 | ApplicationConfig, ApplicationAlertConfig |
| Events | 2 | CustomEventSpecification, BuiltinEventSpecification |
| Infrastructure | 2 | InfraAlertConfig, HostAgent |
| RBAC | 4 | Group, Role, Team, User |
| SLO | 4 | SliConfig, SloConfig, SloAlertConfig, SloCorrectionConfig |
| Synthetic | 3 | SyntheticTest, SyntheticLocation, SyntheticAlertConfig |
| Website | 2 | WebsiteMonitoringConfig, WebsiteAlertConfig |
| Automation | 2 | AutomationAction, AutomationPolicy |
| Logs | 1 | LogAlertConfig |
| Maintenance | 1 | MaintenanceWindowConfig |
| Dashboard | 1 | CustomDashboard |
| Tokens | 1 | APIToken |

#### Type Definitions (Priority: Medium)
- Enums: severity, operator, aggregation, access-type, relation-type, log-level, granularity
- Scopes: application-config-scope, boundary-scope
- Evaluation types: application-alert-evaluation-type, infra-alert-evaluation-type
- Other types: threshold, tag-filter, alert-event-type, alerting-channel-type

#### Shared Components (Priority: Medium)
- `internal/shared/tagfilter/` - Tag filter parsing (6 files)

### 1.4 Components to Keep in Provider

- `main.go` - Provider entry point
- `internal/provider/` - Provider implementation
- `internal/datasources/` - All data sources
- `internal/resources/` - All resources
- `internal/resourcehandle/` - Resource abstraction
- `internal/restapi/provider-models.go` - Provider metadata
- `internal/shared/` - Terraform schema mappings
- All Terraform schema definitions
- State management logic
- Terraform validators and plan modifiers

### 1.5 Coupling Analysis

**Critical Coupling Points:**
1. **Direct Imports:** 31 files (24 resources + 7 data sources) import `internal/restapi`
2. **ProviderMeta Struct:** Couples API to Terraform context
3. **No Context Support:** API methods don't accept `context.Context`
4. **Tag Filter Split:** Logic divided between packages
5. **Shared Test Utilities:** Reference both layers

**Impact Assessment:**
- **High Impact:** Direct imports require systematic refactoring
- **Medium Impact:** Context propagation needs careful implementation
- **Low Impact:** Test utilities can be duplicated or shared

---

## 2. Code Separation Strategy

### 2.1 Target Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  Terraform Provider Instana                  │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Resources & Data Sources (Terraform Schema Layer)     │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Converters (TF Types ↔ API Models)                   │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Provider (Configuration & Client Initialization)      │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            ↓ depends on
┌─────────────────────────────────────────────────────────────┐
│              Instana Go Client Library                       │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Client Interface (API Accessor Methods)               │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  API Implementations (13 Domains)                      │ │
│  │  - Alerting  - Application  - Events  - Infrastructure│ │
│  │  - RBAC  - SLO  - Synthetic  - Website  - Automation  │ │
│  │  - Logs  - Maintenance  - Dashboard  - Tokens         │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  HTTP Client (Rate Limiting, Retry, Context Support)  │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Go Client Library Structure

```
github.com/instana/go-instana-client/
├── client/                   # Core client
│   ├── client.go            # Main interface
│   ├── config.go            # Configuration
│   ├── http_client.go       # HTTP implementation
│   └── errors.go            # Error types
├── api/                     # API by domain
│   ├── alerting/
│   ├── application/
│   ├── events/
│   ├── infrastructure/
│   ├── rbac/
│   ├── slo/
│   ├── synthetic/
│   ├── website/
│   ├── automation/
│   ├── logs/
│   ├── maintenance/
│   ├── dashboard/
│   └── tokens/
├── types/                   # Common types
│   ├── enums.go
│   ├── severity.go
│   ├── operators.go
│   └── tagfilter/
├── internal/rest/           # Internal utilities
└── examples/                # Usage examples
```

### 2.3 Key Interfaces

```go
// Client interface
type Client interface {
    Config() *Config
    Alerting() AlertingAPI
    Application() ApplicationAPI
    // ... other domains
}

// Generic resource interface
type Resource[T DataObject] interface {
    GetAll(ctx context.Context) ([]T, error)
    GetOne(ctx context.Context, id string) (T, error)
    Create(ctx context.Context, data T) (T, error)
    Update(ctx context.Context, data T) (T, error)
    Delete(ctx context.Context, id string) error
}
```

### 2.4 Migration Strategy

**Phase-based Approach:**
1. **Extract Core** → HTTP client and interfaces
2. **Extract Types** → Enums and common types
3. **Extract APIs** → Domain by domain (13 domains)
4. **Refactor Provider** → Update to use client library
5. **Test & Validate** → Comprehensive testing

**File-by-File Migration:**
- Each file migrated individually
- Tests created alongside
- Verification at each step
- No big-bang approach

---

## 3. Execution Plan

### Overview

**Total Duration:** 14 weeks  
**Phases:** 10 major phases  
**Approach:** Iterative and incremental

### Phase 1: Planning & Analysis (Week 1)

**Objectives:**
- Complete architectural analysis
- Finalize design decisions
- Get stakeholder approval

**Tasks:**
- [x] Analyze codebase structure
- [x] Identify coupling points
- [x] Design package structure
- [ ] Create migration checklist
- [ ] Set up project tracking
- [ ] Get approval from stakeholders

**Deliverables:**
- This refactoring plan
- Dependency graph
- Migration checklist
- Project timeline

**Success Criteria:**
- All stakeholders aligned
- Clear understanding of scope
- Risks identified and mitigated

### Phase 2: Design (Week 1-2)

**Objectives:**
- Design client library interfaces
- Define error handling strategy
- Plan context propagation

**Tasks:**
- [ ] Define all interfaces
- [ ] Design error types
- [ ] Plan authentication flow
- [ ] Design rate limiting
- [ ] Design retry logic
- [ ] Create API design document

**Deliverables:**
- API design document
- Interface definitions
- Error handling spec
- Authentication design

**Success Criteria:**
- Interfaces reviewed and approved
- Design supports all use cases
- No breaking changes identified

### Phase 3: Project Setup (Week 2)

**Objectives:**
- Create client library repository
- Set up development environment
- Configure CI/CD basics

**Tasks:**
- [ ] Create `github.com/instana/go-instana-client` repo
- [ ] Initialize Go module
- [ ] Set up directory structure
- [ ] Configure golangci-lint
- [ ] Create README and docs
- [ ] Set up GitHub Actions
- [ ] Configure branch protection

**Deliverables:**
- New repository
- Basic CI/CD pipeline
- Development guidelines

**Success Criteria:**
- Repository accessible
- CI/CD running
- Team can start development

### Phase 4: Core Migration (Week 3-6)

**Week 3: Core Infrastructure**

Tasks:
- [ ] Extract HTTP client → `client/http_client.go`
  - Add context support
  - Implement error wrapping
  - Add logging
- [ ] Create main client → `client/client.go`
- [ ] Implement config → `client/config.go`
- [ ] Create errors → `client/errors.go`
- [ ] Extract generic resources → `internal/rest/`
- [ ] Add comprehensive tests

**Week 4: Type System**

Tasks:
- [ ] Migrate enums to `types/`
  - severity, operator, aggregation
  - access-type, relation-type, log-level
  - granularity, scopes, evaluation types
- [ ] Migrate tag filter → `types/tagfilter/`
- [ ] Add String() methods
- [ ] Add validation methods
- [ ] Create tests

**Week 5: API Domains (Part 1)**

Tasks:
- [ ] Alerting APIs → `api/alerting/`
- [ ] Application APIs → `api/application/`
- [ ] Events APIs → `api/events/`
- [ ] Infrastructure APIs → `api/infrastructure/`
- [ ] RBAC APIs → `api/rbac/`
- [ ] Add tests for each

**Week 6: API Domains (Part 2)**

Tasks:
- [ ] SLO APIs → `api/slo/`
- [ ] Synthetic APIs → `api/synthetic/`
- [ ] Website APIs → `api/website/`
- [ ] Automation APIs → `api/automation/`
- [ ] Logs APIs → `api/logs/`
- [ ] Maintenance APIs → `api/maintenance/`
- [ ] Dashboard APIs → `api/dashboard/`
- [ ] Tokens APIs → `api/tokens/`
- [ ] Add tests for each

**Deliverables:**
- Complete client library implementation
- Unit tests for all components
- API documentation

**Success Criteria:**
- All APIs migrated
- Tests passing
- Code coverage ≥ 80%

### Phase 5: Testing (Week 7)

**Objectives:**
- Achieve comprehensive test coverage
- Validate all functionality

**Tasks:**
- [ ] Set up test infrastructure
- [ ] Create mock HTTP server
- [ ] Write unit tests (80%+ coverage)
- [ ] Create integration tests
- [ ] Add example tests
- [ ] Set up fixtures
- [ ] Configure coverage reporting
- [ ] Add benchmark tests

**Deliverables:**
- Complete test suite
- Test coverage report
- Performance benchmarks

**Success Criteria:**
- All tests passing
- Coverage ≥ 80%
- No flaky tests

### Phase 6: Provider Refactoring (Week 8-10)

**Week 8: Dependencies & Provider**

Tasks:
- [ ] Update go.mod with client library
- [ ] Remove internal/restapi
- [ ] Update provider.go
- [ ] Create new ProviderMeta
- [ ] Update error handling
- [ ] Test provider initialization

**Week 8-9: Converters**

Tasks:
- [ ] Create `internal/converters/` package
- [ ] Implement converters for all 13 domains
- [ ] Add validation in converters
- [ ] Create converter tests
- [ ] Document conversion logic

**Week 9-10: Resources & Data Sources**

Tasks:
- [ ] Update all 24 resources
- [ ] Update all 7 data sources
- [ ] Update imports
- [ ] Update GetRestResource methods
- [ ] Update state mapping
- [ ] Test each resource

**Deliverables:**
- Refactored provider
- Complete converter layer
- Updated resources

**Success Criteria:**
- Provider compiles
- All resources updated
- Tests passing

### Phase 7: Integration Testing (Week 11)

**Objectives:**
- Validate end-to-end functionality
- Ensure no regressions

**Tasks:**
- [ ] Run full test suite
- [ ] Manual testing of all resources
- [ ] Test CRUD operations
- [ ] Test state management
- [ ] Test error scenarios
- [ ] Regression testing
- [ ] Load testing
- [ ] Test with real Instana API

**Deliverables:**
- Test results
- Bug reports
- Performance metrics

**Success Criteria:**
- All tests passing
- No regressions
- Performance acceptable

### Phase 8: Documentation (Week 12)

**Objectives:**
- Complete all documentation
- Prepare for release

**Tasks:**

Client Library:
- [ ] Write README
- [ ] Create getting started guide
- [ ] Document authentication
- [ ] Create API reference
- [ ] Write usage examples
- [ ] Add godoc comments
- [ ] Create migration guide

Provider:
- [ ] Update provider docs
- [ ] Document breaking changes
- [ ] Create upgrade guide
- [ ] Update CHANGELOG

**Deliverables:**
- Complete documentation
- Migration guides
- Release notes

**Success Criteria:**
- All docs complete
- Examples working
- Clear upgrade path

### Phase 9: CI/CD (Week 12-13)

**Objectives:**
- Set up production CI/CD
- Automate releases

**Tasks:**

Client Library:
- [ ] Configure test workflows
- [ ] Set up linting
- [ ] Configure security scanning
- [ ] Set up release automation
- [ ] Configure semantic versioning

Provider:
- [ ] Update test workflows
- [ ] Update release workflow
- [ ] Test release process

**Deliverables:**
- Production CI/CD pipelines
- Release automation

**Success Criteria:**
- CI/CD working
- Releases automated
- Security scanning active

### Phase 10: Deployment (Week 13-14)

**Objectives:**
- Release both projects
- Support users during migration

**Tasks:**

Week 13:
- [ ] Create alpha releases
- [ ] Internal testing
- [ ] Create beta releases
- [ ] Extended testing
- [ ] Address feedback

Week 14:
- [ ] Release client library v0.1.0
- [ ] Release provider v4.0.0
- [ ] Publish to registries
- [ ] Announce releases
- [ ] Monitor for issues
- [ ] Provide support

**Deliverables:**
- Released client library
- Released provider
- User support

**Success Criteria:**
- Both projects released
- No critical issues
- Positive user feedback

---

## 4. Go Client Library Release Strategy

### 4.1 Versioning

**Semantic Versioning:**
- Format: `MAJOR.MINOR.PATCH`
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

**Pre-release Versions:**
- Alpha: `v0.1.0-alpha.1` (internal)
- Beta: `v0.1.0-beta.1` (public)
- RC: `v1.0.0-rc.1` (release candidate)

### 4.2 Release Process

1. **Prepare:**
   - Update CHANGELOG
   - Update version
   - Run tests
   - Update docs

2. **Create:**
   - Create Git tag
   - Push tag
   - GitHub Actions creates release

3. **Publish:**
   - Verify on pkg.go.dev
   - Update release notes
   - Announce

4. **Monitor:**
   - Watch for issues
   - Release patches if needed

### 4.3 CI/CD Pipeline

**Automated Testing:**
- Test on Go 1.22, 1.23, 1.24
- Test on Linux, macOS, Windows
- Run linting
- Security scanning
- Coverage reporting

**Release Automation:**
- Triggered on version tags
- Run full test suite
- Create GitHub release
- Publish artifacts

### 4.4 Documentation

**Required:**
- README with quick start
- API reference (godoc)
- Usage examples
- Authentication guide
- Contributing guidelines
- Security policy

---

## 5. Integration Plan

### 5.1 Provider Integration

**go.mod Configuration:**
```go
module github.com/instana/terraform-provider-instana

go 1.24

require (
    github.com/instana/go-instana-client v0.1.0
    github.com/hashicorp/terraform-plugin-framework v1.15.0
    // ... other dependencies
)
```

### 5.2 Version Pinning

**Strategy:**
- Initial: Pin to exact version `v0.1.0`
- After stabilization: Use minor constraint `v0.1` or `v0.x`

**Rationale:**
- Exact pinning ensures stability
- Minor constraints allow patch updates
- Major updates require manual intervention

### 5.3 Dependency Updates

**Regular Updates (Monthly):**
1. Monitor releases
2. Review CHANGELOG
3. Update: `go get github.com/instana/go-instana-client@v0.2.0`
4. Run tests
5. Create PR

**Major Updates:**
1. Review breaking changes
2. Create feature branch
3. Refactor code
4. Update tests
5. Update docs
6. Bump provider version
7. Release with notes

### 5.4 Breaking Changes

**Client Library:**
- 30-day advance notice
- Deprecation warnings
- Migration guide
- Support during transition

**Provider Response:**
- Evaluate impact
- Plan migration
- Test thoroughly
- Communicate to users

---

## 6. Deployment Strategy

### 6.1 Repository Structure

**Client Library:**
- `main` branch: Stable
- `develop` branch: Development
- `feature/*`: Features
- `hotfix/*`: Critical fixes

**Provider:**
- `master` branch: Stable
- `develop` branch: Development
- `feature/*`: Features
- `hotfix/*`: Critical fixes

**Branch Protection:**
- Require 2 PR reviews
- Require status checks
- No direct pushes

### 6.2 Release Cadence

**Client Library:**
- Major: As needed
- Minor: Monthly
- Patch: As needed

**Provider:**
- Major: Annually
- Minor: Quarterly
- Patch: As needed

**Coordination:**
- Client releases first
- Provider follows within 1-2 weeks

### 6.3 Backward Compatibility

**Guarantees:**
- No breaking changes in minor/patch
- Deprecation warnings
- 2-version support for deprecated features
- Clear migration paths

**Compatibility Matrix:**
| Provider | Client Library | Terraform |
|----------|---------------|-----------|
| v4.0.x   | v0.1.x       | >= 1.0    |
| v4.1.x   | v0.2.x       | >= 1.0    |
| v5.0.x   | v1.0.x       | >= 1.0    |

### 6.4 User Migration

**Timeline:**
- Week 1: Announce
- Week 2-4: Beta testing
- Week 5+: Stable release

**User Steps:**
1. Review migration guide
2. Update provider version
3. Run `terraform init -upgrade`
4. Test in non-production
5. Deploy to production

**No Config Changes:**
- Terraform configs unchanged
- State files compatible
- Internal changes only

---

## 7. Quality Assurance Plan

### 7.1 Unit Testing

**Client Library:**
- Coverage target: 80%+
- Test all API methods
- Test error scenarios
- Test edge cases

**Provider:**
- Coverage target: 75%+
- Test all resources
- Test data sources
- Test converters

### 7.2 Integration Testing

**Client Library:**
- Mock HTTP server
- Rate limiting tests
- Retry logic tests
- Timeout tests

**Provider:**
- CRUD operations
- State management
- Import functionality
- Real API tests (staging)

### 7.3 Contract Testing

- Verify API contracts
- Test request/response formats
- Validate data models
- Test error responses

### 7.4 Performance Testing

- Benchmark HTTP client
- Benchmark unmarshalling
- Load testing
- Concurrent request testing

### 7.5 Security Testing

- Dependency scanning
- Code security analysis
- API token handling
- TLS verification

---

## 8. Documentation Requirements

### 8.1 Client Library

**Essential:**
- README with quick start
- Getting started guide
- Authentication guide
- API reference (godoc)
- Usage examples (13 domains)
- Error handling guide
- Contributing guidelines
- Security policy
- Code of conduct

**Format:**
- Markdown for guides
- Godoc for API reference
- Code examples in Go

### 8.2 Provider

**Essential:**
- Provider configuration
- Resource documentation (24 resources)
- Data source documentation (7 sources)
- Migration guide v3 → v4
- Upgrade notes
- CHANGELOG
- Troubleshooting guide

**Format:**
- Terraform Registry format
- Markdown for guides

### 8.3 Migration Guides

**For Developers:**
- How to use client library
- API changes
- Code examples
- Best practices

**For Users:**
- How to upgrade provider
- What to expect
- Troubleshooting
- Support channels

---

## 9. Risk Assessment

### 9.1 Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Breaking changes | High | Medium | Comprehensive testing, beta releases |
| Performance issues | Medium | Low | Benchmarking, optimization |
| Dependency conflicts | Medium | Low | Careful dependency management |
| State corruption | High | Low | State migration testing, backups |
| API compatibility | High | Low | Contract testing, versioning |

### 9.2 Project Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Timeline delays | Medium | Medium | Buffer time, prioritization |
| Resource constraints | Medium | Low | Clear task allocation |
| Scope creep | Medium | Medium | Strict scope management |
| Stakeholder misalignment | High | Low | Regular communication |

### 9.3 Mitigation Strategies

**For Breaking Changes:**
- Extensive testing at each phase
- Beta releases for feedback
- Clear migration docs
- Support channels

**For Performance:**
- Baseline current performance
- Monitor during migration
- Optimize critical paths
- Load testing

**For Dependencies:**
- Use dependabot
- Regular security audits
- Pin critical dependencies
- Test with multiple versions

### 9.4 Rollback Plan

**If Critical Issues:**
1. Assess severity
2. Halt deployment if critical
3. Communicate to users
4. Fix or rollback
5. Re-test
6. Resume

**Rollback Steps:**
- Revert to previous version
- Document issues
- Plan fixes
- Schedule new release

---

## 10. Timeline & Milestones

### 10.1 Project Timeline

```
Week 1-2:   ████████ Planning & Design
Week 3-6:   ████████████████ Client Library Development
Week 7:     ████ Testing
Week 8-10:  ████████████ Provider Refactoring
Week 11:    ████ Integration Testing
Week 12:    ████ Documentation
Week 13:    ████ CI/CD Setup
Week 14:    ████ Deployment
```

### 10.2 Key Milestones

| Milestone | Week | Deliverable | Success Criteria |
|-----------|------|-------------|------------------|
| M1: Planning Complete | 2 | Design approved | Stakeholder sign-off |
| M2: Client Alpha | 6 | Core working | APIs functional |
| M3: Client Beta | 7 | Tests passing | 80%+ coverage |
| M4: Provider Refactored | 10 | All resources updated | Compiles successfully |
| M5: Integration Complete | 11 | All tests passing | No regressions |
| M6: Docs Complete | 12 | All docs written | Review approved |
| M7: Release Ready | 13 | CI/CD configured | Automation working |
| M8: Deployed | 14 | Both released | Users migrating |

### 10.3 Critical Path

```
Design → Client Dev → Testing → Provider Refactor → Integration → Release
```

**Dependencies:**
- Provider refactoring depends on client library completion
- Integration testing depends on provider refactoring
- Release depends on all testing complete

### 10.4 Resource Allocation

**Team:**
- 2-3 Senior Go Developers (full-time)
- 1 DevOps Engineer (50%)
- 1 QA Engineer (full-time)
- 1 Technical Writer (50%)

**Time Commitment:**
- 14 weeks primary development
- 2-4 weeks post-release support

### 10.5 Success Criteria

**Client Library:**
- ✅ All APIs implemented
- ✅ 80%+ test coverage
- ✅ Documentation complete
- ✅ Published to pkg.go.dev
- ✅ No critical bugs

**Provider:**
- ✅ All resources migrated
- ✅ All tests passing
- ✅ No regressions
- ✅ Documentation updated
- ✅ Published to Terraform Registry

**Overall:**
- ✅ Zero downtime for users
- ✅ No config changes required
- ✅ Positive user feedback
- ✅ Improved maintainability
- ✅ On time and on budget

---

## Appendices

### Appendix A: File Migration Mapping

**Core Files:**
| Current | New Location | Lines |
|---------|-------------|-------|
| `internal/restapi/Instana-api.go` | `client/client.go` | 193 |
| `internal/restapi/rest-client.go` | `client/http_client.go` | 200+ |
| `internal/restapi/instana-rest-resource.go` | `api/resources.go` | 35 |
| `internal/restapi/default-rest-resource.go` | `internal/rest/resource.go` | 160 |

**API Files (50+ files):**
See detailed mapping in project tracking system.

### Appendix B: Code Examples

**Before (Current):**
```go
// Provider initialization
instanaAPI := restapi.NewInstanaAPI(apiToken, endpoint, skipTlsVerify)
resp.ResourceData = &restapi.ProviderMeta{InstanaAPI: instanaAPI}

// Resource usage
func (r *resource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
    return api.AlertingChannels()
}
```

**After (New):**
```go
// Provider initialization
config := &client.Config{
    APIToken:      apiToken,
    Endpoint:      endpoint,
    TLSSkipVerify: skipTlsVerify,
}
instanaClient := client.NewClient(config)
resp.ResourceData = &ProviderMeta{Client: instanaClient}

// Resource usage
func (r *resource) GetRestResource(c *client.Client) api.Resource[*alerting.Channel] {
    return c.Alerting().Channels()
}
```

### Appendix C: Dependencies

**Client Library Dependencies:**
```go
require (
    gopkg.in/resty.v1 v1.12.0
    github.com/alecthomas/participle v0.7.1
    github.com/stretchr/testify v1.10.0
)
```

**Provider Dependencies:**
```go
require (
    github.com/instana/go-instana-client v0.1.0
    github.com/hashicorp/terraform-plugin-framework v1.15.0
    github.com/hashicorp/terraform-plugin-sdk/v2 v2.36.1
)
```

### Appendix D: Contact Information

**Project Team:**
- Project Lead: [Name]
- Technical Lead: [Name]
- DevOps Lead: [Name]
- QA Lead: [Name]
- Documentation Lead: [Name]

**Communication:**
- Slack: #instana-refactoring
- Email: instana-team@example.com
- Meetings: Weekly sync (Mondays 10 AM)

### Appendix E: References

- [Semantic Versioning](https://semver.org/)
- [Go Module Documentation](https://go.dev/ref/mod)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

---

**Document Status:** Final Draft  
**Last Updated:** 2026-03-03  
**Next Review:** Weekly during execution  
**Approval Required:** Yes

---

*End of Document*