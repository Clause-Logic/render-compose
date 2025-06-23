package render

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// Test that validates our generated YAML against Render's official JSON schema
func TestYAMLAgainstRenderSchema(t *testing.T) {
	// Fetch the official schema
	schema, err := fetchRenderSchema()
	if err != nil {
		t.Skipf("Could not fetch Render schema: %v", err)
	}

	schemaLoader := gojsonschema.NewStringLoader(schema)

	tests := []struct {
		name      string
		blueprint func() *Blueprint
	}{
		{
			name: "Web Service with Database",
			blueprint: func() *Blueprint {
				api := NewWebService("api", RuntimeNode).
					WithDomains("api.example.com").
					WithStartCommand("npm start").
					WithGit("https://github.com/example/api", "main").
					WithBuild("npm install").
					WithPlan(PlanStarter).
					WithRegion(RegionOregon).
					WithEnvVars(
						Env("NODE_ENV", "production"),
						EnvFromDatabase("DATABASE_URL", "main-db", DatabasePropertyConnectionString),
					)

				db := NewDatabase("main-db").
					WithPlan(PlanBasic1GB).
					WithPostgreSQL(PostgreSQL16).
					WithRegion(RegionOregon)

				return NewBlueprint().
					WithServices(api).
					WithDatabases(db)
			},
		},
		{
			name: "Background Worker",
			blueprint: func() *Blueprint {
				worker := NewBackgroundWorker("worker", RuntimePython).
					WithStartCommand("python worker.py").
					WithGit("https://github.com/example/worker", "main").
					WithBuild("pip install -r requirements.txt").
					WithPlan(PlanStarter).
					WithRegion(RegionOregon)

				return NewBlueprint().WithServices(worker)
			},
		},
		{
			name: "Cron Job",
			blueprint: func() *Blueprint {
				cron := NewCronJob("cleanup", RuntimeNode, "0 2 * * *").
					WithStartCommand("npm run cleanup").
					WithGit("https://github.com/example/scripts", "main").
					WithRegion(RegionOregon)

				return NewBlueprint().WithServices(cron)
			},
		},
		{
			name: "Static Site",
			blueprint: func() *Blueprint {
				site := NewStaticSite("frontend").
					WithPublishPath("./dist").
					WithDomains("example.com", "www.example.com").
					WithGit("https://github.com/example/frontend", "main").
					WithBuild("npm run build").
					WithRegion(RegionOregon)

				return NewBlueprint().WithServices(site)
			},
		},
		{
			name: "Key-Value Service",
			blueprint: func() *Blueprint {
				cache := NewKeyValueService("cache").
					WithPlan(PlanFree).
					WithRegion(RegionOregon).
					WithPublicAccess()

				return NewBlueprint().WithServices(cache)
			},
		},
		{
			name: "Complex Multi-Service Blueprint",
			blueprint: func() *Blueprint {
				// Web service with auto-scaling
				api := NewWebService("api", RuntimeNode).
					WithDomains("api.example.com").
					WithStartCommand("npm start").
					WithGit("https://github.com/example/api", "main").
					WithBuild("npm install").
					WithAutoScaling(2, 10, 70).
					WithPlan(PlanStandard).
					WithRegion(RegionOregon).
					WithHealthCheck("/health").
					WithEnvVars(
						Env("NODE_ENV", "production"),
						EnvFromDatabase("DATABASE_URL", "main-db", DatabasePropertyConnectionString),
						EnvFromService("CACHE_URL", "cache", ServiceTypeKeyValue, ServicePropertyConnectionString),
						EnvSecret("JWT_SECRET"),
					)

				// Background worker
				worker := NewBackgroundWorker("worker", RuntimePython).
					WithStartCommand("python worker.py").
					WithGit("https://github.com/example/worker", "main").
					WithBuild("pip install -r requirements.txt").
					WithPlan(PlanStarter).
					WithRegion(RegionOregon).
					WithEnvVars(
						EnvFromDatabase("DATABASE_URL", "main-db", DatabasePropertyConnectionString),
					)

				// Database with high availability
				db := NewDatabase("main-db").
					WithPlan(PlanPro8GB).
					WithPostgreSQL(PostgreSQL16).
					WithRegion(RegionOregon).
					WithHighAvailability().
					WithReadReplicas("main-db-replica")

				// Cache
				cache := NewKeyValueService("cache").
					WithPlan(PlanStarter).
					WithRegion(RegionOregon).
					WithPublicAccess().
					WithMaxMemoryPolicy(MaxMemoryPolicyAllKeysLRU)

				// Environment variable group
				env := NewEnvVarGroup("shared").
					WithEnv("APP_NAME", "example-app").
					WithEnv("LOG_LEVEL", "info").
					WithSecret("API_KEY")

				return NewBlueprint().
					WithServices(api, worker).
					WithServices(cache).
					WithDatabases(db).
					WithEnvVarGroups(env).
					WithPreviews(PreviewGenerationAutomatic, 7)
			},
		},
		{
			name: "Docker Service",
			blueprint: func() *Blueprint {
				api := NewWebService("api", RuntimeDocker).
					WithDomains("api.example.com").
					WithDockerfile("./Dockerfile", ".").
					WithGit("https://github.com/example/api", "main").
					WithPlan(PlanStarter).
					WithRegion(RegionOregon)

				return NewBlueprint().WithServices(api)
			},
		},
		{
			name: "Private Service",
			blueprint: func() *Blueprint {
				pserv := NewPrivateService("internal-api", RuntimeNode).
					WithStartCommand("npm start").
					WithGit("https://github.com/example/internal", "main").
					WithBuild("npm install").
					WithPlan(PlanStarter).
					WithRegion(RegionOregon)

				return NewBlueprint().WithServices(pserv)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blueprint := tt.blueprint()

			// Convert blueprint to YAML string
			yamlStr, err := blueprint.ToYAMLString()
			if err != nil {
				t.Fatalf("Failed to convert blueprint to YAML: %v", err)
			}

			// Convert YAML to JSON for schema validation
			var yamlData interface{}
			if err := yaml.Unmarshal([]byte(yamlStr), &yamlData); err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			jsonData, err := json.Marshal(yamlData)
			if err != nil {
				t.Fatalf("Failed to convert to JSON: %v", err)
			}

			documentLoader := gojsonschema.NewBytesLoader(jsonData)

			// Validate against schema
			result, err := gojsonschema.Validate(schemaLoader, documentLoader)
			if err != nil {
				t.Fatalf("Schema validation failed: %v", err)
			}

			if !result.Valid() {
				t.Errorf("Generated YAML does not match Render schema:")
				for _, desc := range result.Errors() {
					t.Errorf("- %s", desc)
				}
				t.Logf("Generated YAML:\n%s", yamlStr)
			}
		})
	}
}

// Test validation of invalid configurations
func TestInvalidConfigurationsAgainstSchema(t *testing.T) {
	schema, err := fetchRenderSchema()
	if err != nil {
		t.Skipf("Could not fetch Render schema: %v", err)
	}

	schemaLoader := gojsonschema.NewStringLoader(schema)

	tests := []struct {
		name      string
		blueprint func() *Blueprint
		expectErr bool
	}{
		{
			name: "Service with invalid type",
			blueprint: func() *Blueprint {
				// Invalid service type that doesn't exist in the schema
				runtime := RuntimeNode
				bp := &Blueprint{
					Services: []Service{
						{
							Name:    "test",
							Type:    ServiceType("invalid-type"),
							Runtime: &runtime,
						},
					},
				}
				return bp
			},
			expectErr: true,
		},
		{
			name: "Invalid runtime for service type",
			blueprint: func() *Blueprint {
				// Key-value services shouldn't have runtime
				kvs := NewKeyValueService("cache")
				service := kvs.ToService()
				runtime := RuntimeNode
				service.Runtime = &runtime // This should be invalid
				
				return &Blueprint{
					Services: []Service{*service},
				}
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blueprint := tt.blueprint()

			yamlStr, err := blueprint.ToYAMLString()
			if err != nil && tt.expectErr {
				return // Expected to fail at YAML generation
			}

			var yamlData interface{}
			if err := yaml.Unmarshal([]byte(yamlStr), &yamlData); err != nil {
				if tt.expectErr {
					return // Expected to fail
				}
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			jsonData, err := json.Marshal(yamlData)
			if err != nil {
				t.Fatalf("Failed to convert to JSON: %v", err)
			}

			documentLoader := gojsonschema.NewBytesLoader(jsonData)
			result, err := gojsonschema.Validate(schemaLoader, documentLoader)
			if err != nil {
				t.Fatalf("Schema validation failed: %v", err)
			}

			if tt.expectErr && result.Valid() {
				t.Errorf("Expected validation to fail, but it passed")
				t.Logf("Generated YAML:\n%s", yamlStr)
			}

			if !tt.expectErr && !result.Valid() {
				t.Errorf("Unexpected validation failure:")
				for _, desc := range result.Errors() {
					t.Errorf("- %s", desc)
				}
				t.Logf("Generated YAML:\n%s", yamlStr)
			}
		})
	}
}

// Fetch the official Render schema
func fetchRenderSchema() (string, error) {
	resp, err := http.Get("https://render.com/schema/render.yaml.json")
	if err != nil {
		return "", fmt.Errorf("failed to fetch schema: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch schema: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read schema: %w", err)
	}

	return string(body), nil
}

// Benchmark schema validation performance
func BenchmarkSchemaValidation(b *testing.B) {
	schema, err := fetchRenderSchema()
	if err != nil {
		b.Skipf("Could not fetch Render schema: %v", err)
	}

	schemaLoader := gojsonschema.NewStringLoader(schema)

	// Create a complex blueprint for benchmarking
	api := NewWebService("api", RuntimeNode).
		WithDomains("api.example.com").
		WithAutoScaling(2, 10, 70).
		WithEnvVars(
			Env("NODE_ENV", "production"),
			EnvFromDatabase("DATABASE_URL", "db", DatabasePropertyConnectionString),
		)

	db := NewDatabase("db").
		WithPlan(PlanBasic1GB).
		WithPostgreSQL(PostgreSQL16)

	blueprint := NewBlueprint().
		WithServices(api).
		WithDatabases(db)

	yamlStr, err := blueprint.ToYAMLString()
	if err != nil {
		b.Fatalf("Failed to generate YAML: %v", err)
	}

	var yamlData interface{}
	yaml.Unmarshal([]byte(yamlStr), &yamlData)
	jsonData, _ := json.Marshal(yamlData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		documentLoader := gojsonschema.NewBytesLoader(jsonData)
		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
		if !result.Valid() {
			b.Fatalf("Generated YAML is invalid")
		}
	}
}