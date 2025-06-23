# Render Compose

> **Note**: This project is not affiliated with, endorsed by, or sponsored by Render Services Inc. Render and render.yaml are trademarks of Render Services Inc.

**Infrastructure as Code for Render.com with Go**

Render Compose is a Go library that lets you build, organize, and compose [Render Blueprint](https://render.com/docs/infrastructure-as-code) infrastructure using type-safe Go code instead of massive YAML files.

## Why Render Compose?

**Before:** One giant `render.yaml` file with hundreds of lines, no type safety, no team boundaries.

**After:** Modular infrastructure components that teams can own, compose, and deploy safely.

```go
// Build infrastructure with code, not YAML
api := render.NewWebService("user-api", render.RuntimeNode).
    WithDomains("api.myapp.com").
    WithAutoScaling(3, 10, 70).
    WithEnvVars(
        render.EnvFromDatabase("DATABASE_URL", "user-db", render.DatabasePropertyConnectionString),
        render.EnvSecret("JWT_SECRET"),
    )

db := render.NewDatabase("user-db").
    WithPlan(render.PlanBasic1GB).
    WithHighAvailability().
    WithReadReplicas("user-db-replica")

blueprint := render.NewBlueprint().
    WithServices(api).
    WithDatabases(db)

// Write to render.yaml
blueprint.WriteRenderYAML()
```

## Key Features

### 🔒 **Type Safety**
- All Render Blueprint properties mapped to Go structs with proper enums
- Compile-time validation prevents invalid configurations
- IDE autocomplete for all fields and values

### 🧩 **Modular Architecture** 
- Organize infrastructure across multiple files and directories
- Team-owned components that compose into complete deployments
- Load → Merge → Deploy workflow

### 🎯 **Focused Service Types**
- `WebService`, `BackgroundWorker`, `CronJob`, `StaticSite`, `KeyValueService`
- Each type only exposes relevant configuration options
- No more "bag of everything" structs

### 🔧 **Builder Patterns**
- Fluent APIs for constructing infrastructure iteratively
- Logical grouping of related configuration (Docker, Git, Scaling, etc.)
- Helper functions for common patterns

### 🔄 **Composition Functions**
- Merge blueprints from different teams/environments
- Prefix resources to avoid naming conflicts
- Validation and conflict detection

### 📁 **File Organization**
- Read/write YAML files from any directory structure
- Compose distributed infrastructure into single `render.yaml`
- Enable team ownership while maintaining deployment simplicity

## Installation

```bash
go get github.com/Clause-Logic/render-compose
```

## Quick Start

### 1. Define Your Infrastructure

```go
package main

import "github.com/Clause-Logic/render-compose"

func main() {
    // Create a web service
    api := render.NewWebService("user-api", render.RuntimeNode).
        WithDomains("api.myapp.com").
        WithGit("https://github.com/company/user-service", "main").
        WithBuild("npm install").
        WithStartCommand("npm start").
        WithAutoScaling(2, 10, 70)

    // Create a database
    db := render.NewDatabase("user-db").
        WithPlan(render.PlanBasic1GB).
        WithPostgreSQL(render.PostgreSQL16).
        WithPublicAccess()

    // Create environment variables
    env := render.NewEnvVarGroup("shared").
        WithEnv("NODE_ENV", "production").
        WithEnvVars(
            render.EnvFromDatabase("DATABASE_URL", "user-db", render.DatabasePropertyConnectionString),
            render.EnvSecret("JWT_SECRET"),
        )

    // Build the blueprint
    blueprint := render.NewBlueprint().
        WithServices(api).
        WithDatabases(db).
        WithEnvVarGroups(env)

    // Write to render.yaml
    blueprint.WriteRenderYAML()
}
```

### 2. Run and Deploy

```bash
go run main.go
git add render.yaml
git commit -m "Update infrastructure"
git push  # Render deploys automatically
```

## Team Workflow

### Organize by Team

```
infrastructure/
├── teams/
│   ├── frontend/
│   │   └── services.yaml      # Frontend team owns this
│   ├── backend/
│   │   └── services.yaml      # Backend team owns this  
│   └── data/
│       └── databases.yaml     # Data team owns this
├── environments/
│   ├── staging.yaml           # Environment-specific config
│   └── production.yaml
├── shared/
│   └── common.yaml            # Shared infrastructure
└── main.go                    # Composition logic
```

### Compose Everything

```go
// main.go - Composition logic
func main() {
    // Load team components
    frontend, _ := render.LoadRenderYAMLFrom("teams/frontend")
    backend, _ := render.LoadRenderYAMLFrom("teams/backend") 
    data, _ := render.LoadRenderYAMLFrom("teams/data")
    shared, _ := render.LoadRenderYAMLFrom("shared")

    // Load environment config
    env := os.Getenv("ENVIRONMENT")
    envConfig, _ := render.LoadRenderYAMLFrom("environments/" + env)

    // Compose the stack
    base, _ := render.MergeBlueprints(shared, frontend)
    base, _ = render.MergeBlueprints(base, backend)
    base, _ = render.MergeBlueprints(base, data)

    // Apply environment-specific changes
    final, _ := render.MergeBlueprints(base, envConfig)

    // Write final render.yaml
    final.WriteRenderYAML()
}
```

### Build for Different Environments

```bash
ENVIRONMENT=staging go run main.go   # Generates staging render.yaml
ENVIRONMENT=production go run main.go # Generates production render.yaml
```

## Advanced Patterns

### Avoid Naming Conflicts

```go
// Each team prefixes their resources
frontendServices := render.PrefixBlueprint(frontend, "fe-")
backendServices := render.PrefixBlueprint(backend, "be-")

// Now they can safely merge without conflicts
combined, _ := render.MergeBlueprints(frontendServices, backendServices)
```

### Environment-Specific Scaling

```go
// Base infrastructure
base := render.NewWebService("api", render.RuntimeNode).
    WithDomains("api.myapp.com")

// Environment-specific modifications
if os.Getenv("ENVIRONMENT") == "production" {
    base = base.WithAutoScaling(5, 20, 70).WithPlan(render.PlanStandard)
} else {
    base = base.WithScaling(1).WithPlan(render.PlanStarter)
}
```

### Validation and Safety

```go
// Validate before deploying
if errors := render.ValidateBlueprint(blueprint); len(errors) > 0 {
    log.Fatalf("Invalid blueprint: %v", errors)
}

// Check for conflicts before merging
if conflicts := render.FindConflicts(base, overlay); len(conflicts) > 0 {
    log.Printf("Conflicts found: %v", conflicts)
    overlay = render.PrefixBlueprint(overlay, "team2-")
}

final, err := render.MergeBlueprints(base, overlay)
if err != nil {
    log.Fatal(err)
}
```

## Service Types

### Web Services
```go
api := render.NewWebService("api", render.RuntimeNode).
    WithDomains("api.myapp.com").
    WithHealthCheck("/health").
    WithAutoScaling(2, 10, 70).
    WithDocker(&render.DockerConfig{
        DockerfilePath: func() *string { s := "./Dockerfile"; return &s }(),
    })
```

### Background Workers
```go
worker := render.NewBackgroundWorker("worker", render.RuntimePython).
    WithStartCommand("python worker.py").
    WithEnvVars(
        render.EnvFromDatabase("DATABASE_URL", "db", render.DatabasePropertyConnectionString),
    )
```

### Cron Jobs
```go
cleanup := render.NewCronJob("cleanup", render.RuntimeNode, "0 2 * * *").
    WithStartCommand("npm run cleanup").
    WithEnvVars(render.EnvFromDatabase("DATABASE_URL", "db", render.DatabasePropertyConnectionString))
```

### Static Sites
```go
frontend := render.NewStaticSite("frontend").
    WithPublishPath("./dist").
    WithDomains("myapp.com", "www.myapp.com").
    WithBuild("npm run build")
```

### Key-Value Stores
```go
cache := render.NewKeyValueService("cache").
    WithPlan(render.PlanStarter).
    WithPublicAccess().
    WithMaxMemoryPolicy(render.MaxMemoryPolicyAllKeysLRU)
```

## Library Structure

The library is organized into focused modules:

- **`types.go`** - Core Blueprint structs and all string enums (ServiceType, Runtime, Plan, etc.)
- **`services.go`** - Service-specific types and builders (WebService, BackgroundWorker, CronJob, etc.)
- **`resources.go`** - Database and EnvVarGroup builders, Blueprint composition functions
- **`operations.go`** - Blueprint utilities (MergeBlueprints, PrefixBlueprint, ValidateBlueprint, etc.)
- **`io.go`** - File I/O operations (WriteToFile, LoadFromFile, ToYAMLString, etc.)

## API Reference

### Core Types

```go
// String enums for type safety
type ServiceType string
type Runtime string
type Plan string
type Region string

// Core structures
type Blueprint struct { ... }
type Service struct { ... }
type Database struct { ... }
type EnvVarGroup struct { ... }
```

### Builder Functions

```go
// Service builders
func NewWebService(name string, runtime Runtime) *WebService
func NewBackgroundWorker(name string, runtime Runtime) *BackgroundWorker
func NewCronJob(name string, runtime Runtime, schedule string) *CronJob
func NewStaticSite(name string) *StaticSite
func NewKeyValueService(name string) *KeyValueService

// Resource builders
func NewDatabase(name string) *Database
func NewEnvVarGroup(name string) *EnvVarGroup
func NewBlueprint() *Blueprint

// Environment variable helpers
func Env(key, value string) EnvVar
func EnvFromDatabase(key, dbName string, property DatabaseProperty) EnvVar
func EnvFromService(key, serviceName string, serviceType ServiceType, property ServiceProperty) EnvVar
func EnvSecret(key string) EnvVar
func EnvGenerated(key string) EnvVar
```

### Operations

```go
// Blueprint operations
func MergeBlueprints(base, overlay *Blueprint) (*Blueprint, error)
func CopyBlueprint(bp *Blueprint) *Blueprint
func PrefixBlueprint(bp *Blueprint, prefix string) *Blueprint
func ValidateBlueprint(bp *Blueprint) []string
func FindConflicts(base, overlay *Blueprint) []string

// I/O operations
func (bp *Blueprint) WriteToFile(path string) error
func (bp *Blueprint) WriteRenderYAML() error
func (bp *Blueprint) ToYAMLString() (string, error)
func LoadFromFile(path string) (*Blueprint, error)
func LoadRenderYAML() (*Blueprint, error)
```

## Why Not Just YAML?

| YAML | Render Compose |
|------|----------------|
| ❌ No type safety | ✅ Compile-time validation |
| ❌ No IDE support | ✅ Full autocomplete |
| ❌ Single massive file | ✅ Modular organization |
| ❌ No team boundaries | ✅ Team-owned components |
| ❌ Manual conflict resolution | ✅ Automatic conflict detection |
| ❌ Copy-paste reuse | ✅ Programmatic composition |
| ❌ Hard to validate | ✅ Built-in validation |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

MIT License - see LICENSE file for details.

---

**Built with ❤️ for teams who want to manage Render infrastructure like code, not like config files.**