# Migration Analysis: Terraform-Specific Code in Go-Client

## Issues Found

### 1. Terraform-Specific Code in go-client Repository

#### A. `instana/rest-client.go` (Lines 154-180)
**Issue**: Contains Terraform-specific user-agent logic
```go
terraformProviderVersion := "Terraform"
file, err := os.Open(basepath + "/CHANGELOG.md")
// ... reads CHANGELOG.md to get Terraform provider version
terraformProviderVersion = "Terraform/" + strings.TrimPrefix(terraformProviderVersion, "v")
```

**Recommendation**: 
- Make user-agent configurable via client options
- Remove hardcoded "Terraform" prefix
- Allow clients to pass their own user-agent string

#### B. `instana/severity.go` (Lines 6-28)
**Issue**: Contains Terraform-specific representation
```go
type Severity struct {
    apiRepresentation       int
    terraformRepresentation string  // <- Terraform-specific
}

func (s Severity) GetTerraformRepresentation() string  // <- Terraform-specific method
```

**Recommendation**:
- Rename `terraformRepresentation` to `stringRepresentation` or `displayName`
- Rename `GetTerraformRepresentation()` to `GetStringRepresentation()` or `String()`
- This makes it generic for any client

#### C. `testutils/terraform-schema-asserts.go` (Entire File)
**Issue**: Entire file is Terraform-specific test utilities
- Imports `github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema`
- Contains `TerraformSchemaAssert` interface and implementation
- All methods are for validating Terraform schema

**Recommendation**:
- **REMOVE** this file from go-client repository
- Move it back to terraform-provider-instana repository
- It's purely for Terraform provider testing, not for API client testing

### 2. References to Old Package in Terraform Provider

#### Status: ✅ CLEAN
All references in terraform provider code have been updated to use `github.com/instana/instana-go-client/client`.

The only remaining references to `internal/restapi` are:
- Test files within `internal/restapi/` directory itself (expected, as we kept original files)
- These can be removed once migration is fully validated

## Recommended Actions

### Priority 1: Remove Terraform-Specific Test Utils from Go-Client
```bash
cd ../instana-go-client
rm testutils/terraform-schema-asserts.go
```

Then copy it back to terraform provider:
```bash
cd ../terraform-provider-instana
cp ../instana-go-client/testutils/terraform-schema-asserts.go testutils/
```

### Priority 2: Make REST Client User-Agent Configurable
Refactor `instana/rest-client.go`:
- Add `userAgent` field to `restClientImpl`
- Add option to `NewClient()` to accept custom user-agent
- Remove CHANGELOG.md reading logic from go-client
- Let terraform provider pass its own user-agent when creating client

### Priority 3: Rename Terraform-Specific Naming in Severity
Refactor `instana/severity.go`:
- Rename `terraformRepresentation` → `stringRepresentation`
- Rename `GetTerraformRepresentation()` → `String()` or `GetStringRepresentation()`
- Update all usages in both repositories

### Priority 4: Update Go-Client Dependencies
Remove terraform-specific dependencies from `go.mod`:
```bash
cd ../instana-go-client
go mod tidy
```

This should remove `github.com/hashicorp/terraform-plugin-sdk/v2` after removing terraform-schema-asserts.go

## Summary

**Current State**:
- ✅ All API models and REST client code successfully migrated
- ✅ Terraform provider successfully uses go-client
- ✅ Both repositories compile successfully
- ⚠️ 3 terraform-specific items remain in go-client (need cleanup)

**Next Steps**:
1. Remove terraform-schema-asserts.go from go-client
2. Make user-agent configurable in REST client
3. Rename terraform-specific naming in Severity type
4. Run tests to ensure everything still works
5. Consider removing old `internal/restapi/` directory after validation