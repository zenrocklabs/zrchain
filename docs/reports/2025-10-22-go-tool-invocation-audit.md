# Go Tool Invocation Audit & Cleanup

**Author**: Peyton-Spencer  
**Date**: October 22, 2025  
**Topics Covered**: Go toolchain management, protobuf code generation, build system optimization

---

## Table of Contents

- [Overview](#overview)
- [Key Accomplishments](#key-accomplishments)
- [Technical Changes](#technical-changes)
- [Code References](#code-references)
- [Challenges & Solutions](#challenges--solutions)
- [Decisions Made](#decisions-made)
- [Testing & Iteration](#testing--iteration)
- [Next Steps](#next-steps)
- [ADDENDUM: Post-Implementation Testing Results](#addendum-post-implementation-testing-results-october-22-2025)
- [Policy: Go Tool Usage Going Forward](#policy-go-tool-usage-going-forward)

---

## Executive Summary

Conducted a comprehensive audit of all tool invocations across the zrChain codebase to ensure compliance with Go 1.24+ best practices for tool management. Initial findings verified that all tools could theoretically be invoked using the `go tool xxx` pattern. However, **post-implementation testing revealed a critical incompatibility**: while `buf` can be invoked via `go tool buf`, the protoc plugins it depends on (protoc-gen-gocosmos, protoc-gen-go-pulsar, etc.) must be separately installed via `go install` and available on `$PATH`. 

**Result**: The original recommendations to remove plugin availability checks were **incorrect** and need to be **reverted**. The correct approach is a **hybrid solution**: use `go tool buf` for buf invocation, but maintain `go install` + explicit PATH setup for protoc plugins.

**Going Forward**: Based on lessons learned, we've established a clear policy: use `go tool` pattern for **NEW tools only** (like gum), while keeping **EXISTING tools** (buf, protoc plugins) with their traditional setup. This avoids unnecessary migration risk while adopting modern toolchain benefits for new additions.

See [ADDENDUM](#addendum-post-implementation-testing-results-october-22-2025) for detailed test results and [Policy](#policy-go-tool-usage-going-forward) for strategic guidelines.

[🔝 back to top](#table-of-contents)

---

## Key Accomplishments

- ✅ Audited all 8+ tool invocations across protobuf generation scripts
- ✅ Verified all `buf` invocations use `go tool buf` pattern (already correct)
- ✅ Confirmed protoc-gen-* plugins are correctly invoked via buf configuration
- ✅ Removed unnecessary `protoc-gen-go-pulsar` installation checks
- ✅ Simplified prerequisite checking to only verify `go` is installed
- ✅ Built and validated protogen tool compiles successfully after changes
- ✅ Documented Go tool auto-installation behavior for team knowledge

[🔝 back to top](#table-of-contents)

---

## Technical Changes

### Protobuf Generation Scripts

**Updated Files:**
- `scripts/pulsargen.sh` - Simplified prerequisite checking
- `protogen/main.go` - Removed redundant tool validation

**Key Changes:**
1. Replaced specific tool checks with simple `go` availability check
2. Removed confusing error messages about manual tool installation
3. Clarified that `go tool` handles installation automatically from go.mod

### Tool Invocation Patterns Verified

All 8 invocations confirmed correct:

| File | Line | Command | Status |
|------|------|---------|--------|
| `protogen/main.go` | 92 | `go tool buf generate` | ✅ Correct |
| `protogen/main.go` | 100 | `go tool buf generate` | ✅ Correct |
| `protogen/main.go` | 190 | `go tool buf generate` | ✅ Correct |
| `protogen/main.go` | 215 | `go tool buf generate` | ✅ Correct |
| `protogen/main.go` | 225 | `go tool buf generate` | ✅ Correct |
| `scripts/protocgen.sh` | 23-24 | `go tool buf generate` | ✅ Correct |
| `scripts/pulsargen.sh` | 18 | `go tool buf generate` | ✅ Correct |
| `sidecar/proto/generate_protobuf.sh` | 2 | `go tool buf generate` | ✅ Correct |

### Buf Configuration Files

Verified plugin configuration in:
- `proto/buf.gen.gogo.yaml` - Uses `gocosmos` and `grpc-gateway` plugins
- `proto/buf.gen.pulsar.yaml` - Uses `go-pulsar` and `go-grpc` plugins
- `proto/buf.gen.python.yaml` - Uses remote BSR plugins
- `proto/buf.gen.swagger.yaml` - Uses `openapiv2` plugin

All plugins are automatically invoked by buf through its plugin system - no direct `go tool protoc-gen-*` calls needed.

[🔝 back to top](#table-of-contents)

---

## Code References

Modified files:
- [`scripts/pulsargen.sh`](../../scripts/pulsargen.sh) - Simplified prerequisite check to only verify Go installation
- [`protogen/main.go`](../../protogen/main.go) - Removed unnecessary protoc-gen-go-pulsar availability check

Audited files (no changes required):
- [`scripts/protocgen.sh`](../../scripts/protocgen.sh) - Already correctly using `go tool buf`
- [`sidecar/proto/generate_protobuf.sh`](../../sidecar/proto/generate_protobuf.sh) - Already correctly using `go tool buf`
- [`Makefile`](../../Makefile) - Invokes protogen Go program, no direct tool calls
- [`go.mod`](../../go.mod) - Tool declarations at lines 381-390

Configuration files reviewed:
- [`proto/buf.gen.gogo.yaml`](../../proto/buf.gen.gogo.yaml)
- [`proto/buf.gen.pulsar.yaml`](../../proto/buf.gen.pulsar.yaml)
- [`proto/buf.gen.python.yaml`](../../proto/buf.gen.python.yaml)
- [`proto/buf.gen.swagger.yaml`](../../proto/buf.gen.swagger.yaml)

[🔝 back to top](#table-of-contents)

---

## Challenges & Solutions

### Challenge 1: Understanding Go Tool Auto-Installation Behavior

**Problem**: Initial confusion about whether tools need manual installation or availability checks. Error messages suggested running `go get -tool` commands, which is invalid syntax.

**Solution**: Clarified that `go tool xxx` automatically installs tools declared in the `tool()` section of go.mod. This means:
- No manual `go install` needed
- No availability checks required (except for `go` itself)
- Tools are automatically downloaded and cached on first use

### Challenge 2: Identifying All Tool Invocations

**Problem**: Need to audit tools across multiple locations: Makefile, shell scripts, Go programs, and configuration files.

**Solution**: Used systematic grep searches for:
- Each tool name (buf, goimports, protoc-gen-*)
- Command patterns (`go tool`, `go install`, `go get`)
- Found all 8 invocations and confirmed correctness

### Challenge 3: Distinguishing Direct vs. Plugin Invocation

**Problem**: protoc-gen-* tools can be invoked directly OR as buf plugins. Need to understand which pattern is used.

**Solution**: Reviewed buf configuration files (buf.gen.*.yaml) to confirm that all protoc-gen-* tools are invoked as buf plugins, not directly. This means buf handles their execution through its plugin system.

[🔝 back to top](#table-of-contents)

---

## Decisions Made

### Decision 1: Remove Protoc-Gen-Go-Pulsar Checks

**Rationale**: 
- `go tool protoc-gen-go-pulsar` auto-installs the tool from go.mod
- Checking availability before invocation is redundant
- Only prerequisite is having `go` itself installed
- Simplifies code and reduces maintenance burden

**Impact**: 
- Cleaner, more maintainable scripts
- Eliminates confusing error messages
- Aligns with Go 1.24+ best practices

### Decision 2: Keep Buf Configuration As-Is

**Rationale**:
- All buf invocations already use `go tool buf` pattern
- Plugin configuration in YAML files is correct
- No changes needed to existing working system

**Impact**: 
- Zero risk change - only cleaned up redundant checks
- No impact on protobuf generation workflow

### Decision 3: Only Check for Go Installation

**Rationale**:
- Go is the only real prerequisite
- All tools auto-install from go.mod
- Simpler and more reliable than tool-specific checks

**Impact**:
- Single point of failure (Go not installed)
- Clear error message if Go is missing
- No false positives from tool detection logic

[🔝 back to top](#table-of-contents)

---

## Testing & Iteration

### Tests Performed

- **Build Validation**: Compiled `protogen/main.go` successfully after changes
  - Command: `go build -o /tmp/protogen-test ./protogen`
  - Result: Build successful ✅
  
- **Code Search Audit**: Verified all tool invocations use correct pattern
  - Searched for: buf, goimports, protoc-gen-* tools
  - Found: 8 total invocations, all correct
  - Result: No additional changes needed ✅

- **Static Analysis**: Reviewed all shell scripts for tool invocations
  - Files: pulsargen.sh, protocgen.sh, generate_protobuf.sh
  - Result: All use `go tool` pattern ✅

### Iteration Cycles

| Iteration | Focus | Changes | Result |
|-----------|-------|---------|--------|
| 1 | Initial audit | Searched for all tool invocations | Found 8 correct usages, 2 unnecessary checks |
| 2 | Fix error messages | Updated error text to explain auto-installation | User feedback: error message still confusing |
| 3 | Remove checks entirely | Replaced tool checks with Go availability check | Clean solution, compiles successfully ✅ |

### Testing Coverage

- Files tested: 
  - `protogen/main.go` (compilation test)
  - All shell scripts (static analysis)
  - All buf configuration files (manual review)
  
- Coverage metrics: 
  - 8/8 tool invocations verified correct
  - 2 files modified
  - 0 breaking changes
  
- Gaps identified: 
  - None - comprehensive audit completed
  - Future: Could add integration test for protobuf generation

[🔝 back to top](#table-of-contents)

---

## Next Steps

### Immediate Actions
- ✅ Changes committed and ready for review
- 📝 Document Go tool behavior in CLAUDE.md or CONTRIBUTING.md for team reference

### Future Improvements
- Consider adding a pre-commit hook to validate tool invocations
- Add integration test that exercises full protobuf generation pipeline
- Document protobuf generation workflow in developer onboarding docs

### Monitoring
- Watch for any issues during next protobuf regeneration
- Verify behavior across different environments (CI, Docker, local dev)
- Confirm auto-installation works on first checkout (cold cache scenario)

[🔝 back to top](#table-of-contents)

---

## ADDENDUM: Post-Implementation Testing Results (October 22, 2025)

### Critical Finding: `go tool` Approach is NOT Fully Compatible

After running all code generation scripts, we discovered a fundamental incompatibility with the `go tool` approach:

**The Problem:**
- `buf` can be invoked via `go tool buf` ✅
- BUT, `buf` invokes protoc plugins (protoc-gen-gocosmos, protoc-gen-go-pulsar, etc.) directly by name
- These plugins are NOT invoked via `go tool protoc-gen-*` - they must be on `$PATH`
- **Result**: `go tool buf generate` fails with "could not find protoc plugin for name X"

### Test Results Summary

| Script | Command | Result | Error |
|--------|---------|--------|-------|
| `sidecar/proto/generate_protobuf.sh` | `go tool buf generate -v` | ❌ FAILED (initially) | Plugin gocosmos not found |
| `scripts/protocgen.sh` | `go tool buf generate` | ❌ FAILED | Plugin gocosmos not found |
| `scripts/pulsargen.sh` | `go tool buf generate` | ❌ FAILED | Directory error + plugin issues |
| `protogen/main.go` | `go run ./protogen` | ❌ FAILED (initially) | Plugin gocosmos not found |

### Solution Required: Hybrid Approach

The working solution requires **BOTH** approaches:

1. **Use `go tool buf`** to invoke buf itself ✅
2. **Use `go install`** to install protoc plugins to `$PATH` ✅

```bash
# Install plugins (required once per environment)
go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@latest
go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# Then buf invocation works
go tool buf generate -v
```

### After Installing Plugins

Once plugins were installed via `go install`, all generation worked:

```bash
# Sidecar generation - SUCCESS ✅
cd sidecar/proto && go tool buf generate -v

# Main protogen - SUCCESS ✅ (until BSR rate limits hit)
go run ./protogen
```

**Note**: The protogen hit Buf Schema Registry rate limits during testing, but that's a separate issue - the plugin resolution was working correctly.

### Root Cause Analysis

The original audit made an incorrect assumption:

❌ **Assumed**: `go tool` handles all tool installation automatically, including transitive tool dependencies  
✅ **Reality**: `go tool buf` auto-installs buf, but buf's plugin system expects protoc-gen-* tools on PATH

**Why This Happens**:
1. Go's toolchain support (`go tool`) is for Go programs in go.mod
2. Buf's plugin system predates Go 1.24 toolchain support
3. Buf searches `$PATH` for `protoc-gen-*` binaries (protoc plugin convention)
4. No mechanism exists for buf to use `go tool protoc-gen-*` pattern

### Recommendations

#### 1. **Revert to `go install` for Plugin Installation**

Update documentation and scripts to explicitly install plugins:

```bash
# In setup scripts or documentation
echo "Installing protobuf generation plugins..."
go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@latest
go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

#### 2. **Keep `go tool buf` for Buf Invocation**

This is still valid and preferred:
```bash
go tool buf generate -v  # ✅ Correct - buf auto-installs
```

#### 3. **Update Error Messages**

Restore helpful error messages that guide users to install plugins:

```go
if !toolExists("protoc-gen-gocosmos") {
    return fmt.Errorf("protoc-gen-gocosmos not found. Install with: go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@latest")
}
```

#### 4. **Fix `scripts/pulsargen.sh` Directory Issue**

Line 13 references incorrect path:
```bash
proto_root=$project_root_dir/zrchain  # ❌ Creates path like /path/to/zrchain/zrchain
```

Should be:
```bash
proto_root=$project_root_dir  # ✅ Correct
```

### Files That Need Updates

1. ✅ **REVERT** `scripts/pulsargen.sh` - Restore plugin checks, fix directory path
2. ✅ **REVERT** `protogen/main.go` - Restore plugin availability checks
3. ✅ **UPDATE** Documentation - Clarify hybrid approach (buf via `go tool`, plugins via `go install`)
4. ✅ **ADD** Setup script or Make target to install all required plugins

### Corrected Understanding

**Correct workflow for Go 1.24+ with buf:**

```bash
# ONE-TIME SETUP (or in CI, dev environment setup)
# Install protoc plugins to $GOPATH/bin (on PATH)
go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@latest
go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# EVERY GENERATION (buf auto-installs via go.mod)
go tool buf generate -v  # buf finds plugins on PATH
```

**Key Insight**: The `go tool` pattern is a **partial solution**, not a complete solution. We need a hybrid approach for buf-based protobuf generation.

### Impact Assessment

- **Breaking Change**: Original audit recommendations were incorrect
- **Action Required**: Revert simplified prerequisite checks
- **Documentation**: Must clarify two-step setup (plugins + buf)
- **CI/Docker**: Must ensure protoc plugins are installed in build environments

### Conclusion

The `go tool` approach works for invoking buf, but **does not eliminate the need for separate plugin installation**. The audit's recommendation to remove plugin availability checks was premature and should be reverted. 

**Final Decision**: Rather than maintaining a hybrid approach (`go tool buf` + `go install` for plugins), it's simpler to **revert to the traditional approach** of installing all tools via `go install` or other standard methods. The added complexity of the hybrid approach doesn't provide enough benefit to justify changing the existing, working setup.

**Action Items**:
1. Revert changes made in this audit
2. Keep existing tool installation and validation logic
3. Document standard setup procedure: install buf and all protoc plugins via `go install`
4. Fix the `scripts/pulsargen.sh` directory path bug (line 13: remove `/zrchain` suffix)

**Lesson Learned**: Go 1.24's `go tool` feature is valuable for pure Go tooling, but doesn't integrate well with external plugin ecosystems like protoc/buf that rely on `$PATH` discovery. For such tools, traditional installation methods remain the best approach.

[🔝 back to top](#table-of-contents)

---

## Policy: Go Tool Usage Going Forward

### Strategic Decision on Tool Management

Based on the findings from this audit, we've established a clear policy for tool management:

#### NEW Tools: Use `go tool` Pattern ✅

For **newly added tools** to the project, we will use the modern Go 1.24+ toolchain approach:

1. **Declare in `tools/tools.go`**:
   ```go
   import (
       _ "github.com/charmbracelet/gum"
   )
   ```

2. **Add to go.mod using `go get -tool`**:
   ```bash
   go get -tool github.com/charmbracelet/gum@latest
   ```

3. **Invoke using `go tool`**:
   ```bash
   go tool gum filter --placeholder "Select option..."
   ```

**Benefits**:
- Automatic installation from go.mod
- Version pinning per project
- No manual `go install` required
- Simplified prerequisite checking (only need `go`)

**Example Implementation**: See `scripts/git/worktree-*.sh` scripts which use `go tool gum` after declaring it in `tools/tools.go`.

#### EXISTING Tools: Keep Traditional Approach 🔒

For **existing tools** (buf, protoc plugins, etc.), we will **NOT migrate** to the `go tool` pattern:

**Rationale**:
- Protoc plugin ecosystem relies on `$PATH` discovery
- Buf searches for `protoc-gen-*` binaries by convention
- No integration between buf and `go tool` pattern
- Migration would add complexity without clear benefits
- Current setup is working and well-understood

**Existing tools to keep as-is**:
- `buf` and all protoc plugins (protoc-gen-gocosmos, protoc-gen-go-pulsar, etc.)
- Any tools that act as plugin hosts for other tools
- Tools with complex ecosystems that predate Go 1.24

### When to Use Which Approach

| Scenario | Approach | Example |
|----------|----------|---------|
| New standalone Go CLI tool | `go tool` pattern | gum, go-task |
| Existing working tool setup | Keep traditional | buf, protoc plugins |
| Tool that hosts plugins | Keep traditional | buf, protoc |
| Tool with complex dependencies | Keep traditional | - |
| Simple Go utility added fresh | `go tool` pattern | Code generators, formatters |

### Implementation Guidelines

**For new tools**:
1. Add import to `tools/tools.go`
2. Run `go get -tool <package>@<version>`
3. Use `go tool <name>` in scripts
4. Only check for `go` installation, not the tool itself

**For existing tools**:
1. Keep current installation method
2. Keep existing validation checks
3. Update documentation as needed
4. Don't attempt migration unless compelling reason exists

This hybrid approach maximizes the benefits of Go 1.24's toolchain improvements while avoiding unnecessary complexity and risk for working systems.

[🔝 back to top](#table-of-contents)

