package render

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestMergeBlueprints(t *testing.T) {
	tests := []struct {
		name      string
		base      *Blueprint
		overlay   *Blueprint
		expected  *Blueprint
		expectErr bool
	}{
		{
			name:      "merge nil blueprints",
			base:      nil,
			overlay:   nil,
			expected:  &Blueprint{},
			expectErr: false,
		},
		{
			name: "merge with nil base",
			base: nil,
			overlay: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
			},
			expected: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
			},
			expectErr: false,
		},
		{
			name: "merge with nil overlay",
			base: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
			},
			overlay: nil,
			expected: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
			},
			expectErr: false,
		},
		{
			name: "successful merge without conflicts",
			base: &Blueprint{
				Services: []Service{
					{Name: "api", Type: ServiceTypeWeb},
				},
				Databases: []Database{
					{Name: "main-db"},
				},
			},
			overlay: &Blueprint{
				Services: []Service{
					{Name: "worker", Type: ServiceTypeWorker},
				},
				EnvVarGroups: []EnvVarGroup{
					{Name: "shared", EnvVars: []EnvVar{{Key: stringPtr("NODE_ENV"), Value: stringPtr("production")}}},
				},
			},
			expected: &Blueprint{
				Services: []Service{
					{Name: "api", Type: ServiceTypeWeb},
					{Name: "worker", Type: ServiceTypeWorker},
				},
				Databases: []Database{
					{Name: "main-db"},
				},
				EnvVarGroups: []EnvVarGroup{
					{Name: "shared", EnvVars: []EnvVar{{Key: stringPtr("NODE_ENV"), Value: stringPtr("production")}}},
				},
			},
			expectErr: false,
		},
		{
			name: "merge with conflicts should fail",
			base: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
			},
			overlay: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWorker}},
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "overlay preview configuration wins",
			base: &Blueprint{
				Previews: &Previews{Generation: "automatic"},
				PreviewsExpireAfterDays: intPtr(30),
			},
			overlay: &Blueprint{
				Previews: &Previews{Generation: "none"},
				PreviewsExpireAfterDays: intPtr(7),
			},
			expected: &Blueprint{
				Previews: &Previews{Generation: "none"},
				PreviewsExpireAfterDays: intPtr(7),
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MergeBlueprints(tt.base, tt.overlay)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !blueprintsEqual(result, tt.expected) {
				t.Errorf("result mismatch\nExpected: %+v\nGot: %+v", tt.expected, result)
			}
		})
	}
}

func TestCopyBlueprint(t *testing.T) {
	original := &Blueprint{
		Services: []Service{
			{Name: "api", Type: ServiceTypeWeb, EnvVars: []EnvVar{{Key: stringPtr("NODE_ENV"), Value: stringPtr("production")}}},
		},
		Databases: []Database{
			{Name: "main-db", Plan: planPtr(PlanBasic1GB)},
		},
		EnvVarGroups: []EnvVarGroup{
			{Name: "shared", EnvVars: []EnvVar{{Key: stringPtr("LOG_LEVEL"), Value: stringPtr("info")}}},
		},
		Previews: &Previews{Generation: "automatic"},
		PreviewsExpireAfterDays: intPtr(7),
	}

	copied := CopyBlueprint(original)

	// Test that the copy is equal to the original
	if !blueprintsEqual(copied, original) {
		t.Errorf("copied blueprint should equal original")
	}

	// Test that modifying the copy doesn't affect the original
	copied.Services[0].Name = "modified-api"
	copied.Databases[0].Name = "modified-db"
	copied.EnvVarGroups[0].Name = "modified-shared"

	if original.Services[0].Name == "modified-api" {
		t.Errorf("modifying copy affected original services")
	}
	if original.Databases[0].Name == "modified-db" {
		t.Errorf("modifying copy affected original databases")
	}
	if original.EnvVarGroups[0].Name == "modified-shared" {
		t.Errorf("modifying copy affected original env groups")
	}

	// Test copying nil blueprint
	nilCopy := CopyBlueprint(nil)
	expected := &Blueprint{}
	if !blueprintsEqual(nilCopy, expected) {
		t.Errorf("copying nil should return empty blueprint")
	}
}

func TestValidateBlueprint(t *testing.T) {
	tests := []struct {
		name     string
		bp       *Blueprint
		expected []string
	}{
		{
			name:     "nil blueprint",
			bp:       nil,
			expected: []string{"blueprint is nil"},
		},
		{
			name: "valid blueprint",
			bp: &Blueprint{
				Services: []Service{
					{Name: "api", Type: ServiceTypeWeb, Runtime: runtimePtr(RuntimeNode)},
				},
				Databases: []Database{
					{Name: "main-db"},
				},
				EnvVarGroups: []EnvVarGroup{
					{Name: "shared"},
				},
			},
			expected: []string{},
		},
		{
			name: "duplicate service names",
			bp: &Blueprint{
				Services: []Service{
					{Name: "api", Type: ServiceTypeWeb, Runtime: runtimePtr(RuntimeNode)},
					{Name: "api", Type: ServiceTypeWorker, Runtime: runtimePtr(RuntimePython)},
				},
			},
			expected: []string{"duplicate service name: api"},
		},
		{
			name: "duplicate database names",
			bp: &Blueprint{
				Databases: []Database{
					{Name: "db"},
					{Name: "db"},
				},
			},
			expected: []string{"duplicate database name: db"},
		},
		{
			name: "duplicate env group names",
			bp: &Blueprint{
				EnvVarGroups: []EnvVarGroup{
					{Name: "shared"},
					{Name: "shared"},
				},
			},
			expected: []string{"duplicate environment group name: shared"},
		},
		{
			name: "missing required fields",
			bp: &Blueprint{
				Services: []Service{
					{Name: "", Type: ServiceTypeWeb, Runtime: runtimePtr(RuntimeNode)},
					{Name: "api", Type: "", Runtime: runtimePtr(RuntimeNode)},
					{Name: "web", Type: ServiceTypeWeb}, // missing runtime
				},
				Databases: []Database{
					{Name: ""},
				},
				EnvVarGroups: []EnvVarGroup{
					{Name: ""},
				},
			},
			expected: []string{
				"service missing name",
				"service api missing type",
				"service web missing runtime",
				"database missing name",
				"environment group missing name",
			},
		},
		{
			name: "key-value service without runtime is valid",
			bp: &Blueprint{
				Services: []Service{
					{Name: "cache", Type: ServiceTypeKeyValue}, // no runtime needed
				},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateBlueprint(tt.bp)

			if len(errors) != len(tt.expected) {
				t.Errorf("expected %d errors, got %d: %v", len(tt.expected), len(errors), errors)
				return
			}

			for i, expectedError := range tt.expected {
				if !strings.Contains(errors[i], expectedError) {
					t.Errorf("expected error containing %q, got %q", expectedError, errors[i])
				}
			}
		})
	}
}

func TestFindConflicts(t *testing.T) {
	tests := []struct {
		name     string
		base     *Blueprint
		overlay  *Blueprint
		expected []string
	}{
		{
			name:     "nil blueprints",
			base:     nil,
			overlay:  nil,
			expected: []string{},
		},
		{
			name: "no conflicts",
			base: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
				Databases: []Database{{Name: "main-db"}},
				EnvVarGroups: []EnvVarGroup{{Name: "shared"}},
			},
			overlay: &Blueprint{
				Services: []Service{{Name: "worker", Type: ServiceTypeWorker}},
				Databases: []Database{{Name: "cache-db"}},
				EnvVarGroups: []EnvVarGroup{{Name: "secrets"}},
			},
			expected: []string{},
		},
		{
			name: "service name conflicts",
			base: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
			},
			overlay: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWorker}},
			},
			expected: []string{"service name conflict: api"},
		},
		{
			name: "database name conflicts",
			base: &Blueprint{
				Databases: []Database{{Name: "main-db"}},
			},
			overlay: &Blueprint{
				Databases: []Database{{Name: "main-db"}},
			},
			expected: []string{"database name conflict: main-db"},
		},
		{
			name: "env group conflicts",
			base: &Blueprint{
				EnvVarGroups: []EnvVarGroup{{Name: "shared"}},
			},
			overlay: &Blueprint{
				EnvVarGroups: []EnvVarGroup{{Name: "shared"}},
			},
			expected: []string{"environment group name conflict: shared"},
		},
		{
			name: "multiple conflicts",
			base: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
				Databases: []Database{{Name: "db"}},
			},
			overlay: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWorker}},
				Databases: []Database{{Name: "db"}},
			},
			expected: []string{
				"service name conflict: api",
				"database name conflict: db",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conflicts := FindConflicts(tt.base, tt.overlay)

			if len(conflicts) != len(tt.expected) {
				t.Errorf("expected %d conflicts, got %d: %v", len(tt.expected), len(conflicts), conflicts)
				return
			}

			// Sort both slices for comparison since order might vary
			sort.Strings(conflicts)
			sort.Strings(tt.expected)

			for i, expected := range tt.expected {
				if conflicts[i] != expected {
					t.Errorf("expected conflict %q, got %q", expected, conflicts[i])
				}
			}
		})
	}
}

func TestPrefixBlueprint(t *testing.T) {
	tests := []struct {
		name     string
		bp       *Blueprint
		prefix   string
		expected *Blueprint
	}{
		{
			name:     "nil blueprint",
			bp:       nil,
			prefix:   "test-",
			expected: &Blueprint{},
		},
		{
			name: "empty prefix",
			bp: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
			},
			prefix: "",
			expected: &Blueprint{
				Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
			},
		},
		{
			name: "prefix all resources",
			bp: &Blueprint{
				Services: []Service{
					{Name: "api", Type: ServiceTypeWeb},
					{Name: "worker", Type: ServiceTypeWorker},
				},
				Databases: []Database{
					{Name: "main-db"},
					{Name: "cache-db"},
				},
				EnvVarGroups: []EnvVarGroup{
					{Name: "shared"},
					{Name: "secrets"},
				},
			},
			prefix: "team1-",
			expected: &Blueprint{
				Services: []Service{
					{Name: "team1-api", Type: ServiceTypeWeb},
					{Name: "team1-worker", Type: ServiceTypeWorker},
				},
				Databases: []Database{
					{Name: "team1-main-db"},
					{Name: "team1-cache-db"},
				},
				EnvVarGroups: []EnvVarGroup{
					{Name: "team1-shared"},
					{Name: "team1-secrets"},
				},
			},
		},
		{
			name: "prefix with internal references",
			bp: &Blueprint{
				Services: []Service{
					{
						Name: "api",
						Type: ServiceTypeWeb,
						EnvVars: []EnvVar{
							{
								Key: stringPtr("DATABASE_URL"),
								FromDatabase: &FromDatabase{
									Name:     "main-db",
									Property: DatabasePropertyConnectionString,
								},
							},
							{
								Key: stringPtr("CACHE_URL"),
								FromService: &FromService{
									Name: "cache",
									Type: ServiceTypeKeyValue,
								},
							},
							{
								Key:       stringPtr("SHARED_VAR"),
								FromGroup: stringPtr("shared"),
							},
						},
					},
					{Name: "cache", Type: ServiceTypeKeyValue},
				},
				Databases: []Database{
					{Name: "main-db"},
				},
				EnvVarGroups: []EnvVarGroup{
					{
						Name: "shared",
						EnvVars: []EnvVar{
							{
								Key: stringPtr("DB_URL"),
								FromDatabase: &FromDatabase{
									Name:     "main-db",
									Property: DatabasePropertyConnectionString,
								},
							},
						},
					},
				},
			},
			prefix: "prod-",
			expected: &Blueprint{
				Services: []Service{
					{
						Name: "prod-api",
						Type: ServiceTypeWeb,
						EnvVars: []EnvVar{
							{
								Key: stringPtr("DATABASE_URL"),
								FromDatabase: &FromDatabase{
									Name:     "prod-main-db",
									Property: DatabasePropertyConnectionString,
								},
							},
							{
								Key: stringPtr("CACHE_URL"),
								FromService: &FromService{
									Name: "prod-cache",
									Type: ServiceTypeKeyValue,
								},
							},
							{
								Key:       stringPtr("SHARED_VAR"),
								FromGroup: stringPtr("prod-shared"),
							},
						},
					},
					{Name: "prod-cache", Type: ServiceTypeKeyValue},
				},
				Databases: []Database{
					{Name: "prod-main-db"},
				},
				EnvVarGroups: []EnvVarGroup{
					{
						Name: "prod-shared",
						EnvVars: []EnvVar{
							{
								Key: stringPtr("DB_URL"),
								FromDatabase: &FromDatabase{
									Name:     "prod-main-db",
									Property: DatabasePropertyConnectionString,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "external references unchanged",
			bp: &Blueprint{
				Services: []Service{
					{
						Name: "api",
						Type: ServiceTypeWeb,
						EnvVars: []EnvVar{
							{
								Key: stringPtr("EXTERNAL_DB_URL"),
								FromDatabase: &FromDatabase{
									Name:     "external-db", // not defined in this blueprint
									Property: DatabasePropertyConnectionString,
								},
							},
						},
					},
				},
			},
			prefix: "team-",
			expected: &Blueprint{
				Services: []Service{
					{
						Name: "team-api",
						Type: ServiceTypeWeb,
						EnvVars: []EnvVar{
							{
								Key: stringPtr("EXTERNAL_DB_URL"),
								FromDatabase: &FromDatabase{
									Name:     "external-db", // should remain unchanged
									Property: DatabasePropertyConnectionString,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PrefixBlueprint(tt.bp, tt.prefix)

			if !blueprintsEqual(result, tt.expected) {
				t.Errorf("result mismatch\nExpected: %+v\nGot: %+v", tt.expected, result)
			}

			// Ensure original blueprint is not modified (only check if prefix is non-empty)
			if tt.bp != nil && tt.prefix != "" {
				originalServices := []string{}
				for _, svc := range tt.bp.Services {
					originalServices = append(originalServices, svc.Name)
				}
				for _, name := range originalServices {
					if strings.HasPrefix(name, tt.prefix) {
						t.Errorf("original blueprint was modified: service name %q has prefix %q", name, tt.prefix)
					}
				}
			}
		})
	}
}

func TestPrefixBlueprintWithSeparator(t *testing.T) {
	bp := &Blueprint{
		Services: []Service{{Name: "api", Type: ServiceTypeWeb}},
	}

	// Test default separator
	result1 := PrefixBlueprintWithSeparator(bp, "team", "")
	expected1 := "team-api"
	if result1.Services[0].Name != expected1 {
		t.Errorf("expected service name %q, got %q", expected1, result1.Services[0].Name)
	}

	// Test custom separator
	result2 := PrefixBlueprintWithSeparator(bp, "team", "_")
	expected2 := "team_api"
	if result2.Services[0].Name != expected2 {
		t.Errorf("expected service name %q, got %q", expected2, result2.Services[0].Name)
	}
}

func TestGetAllResourceNames(t *testing.T) {
	tests := []struct {
		name               string
		bp                 *Blueprint
		expectedServices   []string
		expectedDatabases  []string
		expectedEnvGroups  []string
	}{
		{
			name:               "nil blueprint",
			bp:                 nil,
			expectedServices:   nil,
			expectedDatabases:  nil,
			expectedEnvGroups:  nil,
		},
		{
			name: "empty blueprint",
			bp:   &Blueprint{},
			expectedServices:   nil,
			expectedDatabases:  nil,
			expectedEnvGroups:  nil,
		},
		{
			name: "blueprint with all resource types",
			bp: &Blueprint{
				Services: []Service{
					{Name: "api"},
					{Name: "worker"},
				},
				Databases: []Database{
					{Name: "main-db"},
					{Name: "cache-db"},
				},
				EnvVarGroups: []EnvVarGroup{
					{Name: "shared"},
					{Name: "secrets"},
				},
			},
			expectedServices:  []string{"api", "worker"},
			expectedDatabases: []string{"main-db", "cache-db"},
			expectedEnvGroups: []string{"shared", "secrets"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			services, databases, envGroups := GetAllResourceNames(tt.bp)

			if !reflect.DeepEqual(services, tt.expectedServices) {
				t.Errorf("services mismatch. Expected: %v, Got: %v", tt.expectedServices, services)
			}
			if !reflect.DeepEqual(databases, tt.expectedDatabases) {
				t.Errorf("databases mismatch. Expected: %v, Got: %v", tt.expectedDatabases, databases)
			}
			if !reflect.DeepEqual(envGroups, tt.expectedEnvGroups) {
				t.Errorf("envGroups mismatch. Expected: %v, Got: %v", tt.expectedEnvGroups, envGroups)
			}
		})
	}
}

func TestGetExternalReferences(t *testing.T) {
	tests := []struct {
		name               string
		bp                 *Blueprint
		expectedServices   []string
		expectedDatabases  []string
		expectedEnvGroups  []string
	}{
		{
			name:               "nil blueprint",
			bp:                 nil,
			expectedServices:   nil,
			expectedDatabases:  nil,
			expectedEnvGroups:  nil,
		},
		{
			name: "no external references",
			bp: &Blueprint{
				Services: []Service{
					{
						Name: "api",
						EnvVars: []EnvVar{
							{
								Key: stringPtr("DB_URL"),
								FromDatabase: &FromDatabase{
									Name: "main-db", // defined in this blueprint
								},
							},
						},
					},
				},
				Databases: []Database{{Name: "main-db"}},
			},
			expectedServices:  []string{},
			expectedDatabases: []string{},
			expectedEnvGroups: []string{},
		},
		{
			name: "external references in services",
			bp: &Blueprint{
				Services: []Service{
					{
						Name: "api",
						EnvVars: []EnvVar{
							{
								Key: stringPtr("EXTERNAL_DB_URL"),
								FromDatabase: &FromDatabase{
									Name: "external-db", // not defined in this blueprint
								},
							},
							{
								Key: stringPtr("CACHE_URL"),
								FromService: &FromService{
									Name: "external-cache", // not defined in this blueprint
								},
							},
							{
								Key:       stringPtr("SHARED_VAR"),
								FromGroup: stringPtr("external-env"), // not defined in this blueprint
							},
						},
					},
				},
				Databases: []Database{{Name: "main-db"}}, // only this is internal
			},
			expectedServices:  []string{"external-cache"},
			expectedDatabases: []string{"external-db"},
			expectedEnvGroups: []string{"external-env"},
		},
		{
			name: "external references in env groups",
			bp: &Blueprint{
				EnvVarGroups: []EnvVarGroup{
					{
						Name: "shared",
						EnvVars: []EnvVar{
							{
								Key: stringPtr("EXTERNAL_DB_URL"),
								FromDatabase: &FromDatabase{
									Name: "external-db",
								},
							},
						},
					},
				},
			},
			expectedServices:  []string{},
			expectedDatabases: []string{"external-db"},
			expectedEnvGroups: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			services, databases, envGroups := GetExternalReferences(tt.bp)

			// Sort slices for comparison since order may vary
			sort.Strings(services)
			sort.Strings(databases)
			sort.Strings(envGroups)
			if tt.expectedServices != nil {
				sort.Strings(tt.expectedServices)
			}
			if tt.expectedDatabases != nil {
				sort.Strings(tt.expectedDatabases)
			}
			if tt.expectedEnvGroups != nil {
				sort.Strings(tt.expectedEnvGroups)
			}

			// Handle nil vs empty slice comparison
			if !slicesEqual(services, tt.expectedServices) {
				t.Errorf("external services mismatch. Expected: %v, Got: %v", tt.expectedServices, services)
			}
			if !slicesEqual(databases, tt.expectedDatabases) {
				t.Errorf("external databases mismatch. Expected: %v, Got: %v", tt.expectedDatabases, databases)
			}
			if !slicesEqual(envGroups, tt.expectedEnvGroups) {
				t.Errorf("external envGroups mismatch. Expected: %v, Got: %v", tt.expectedEnvGroups, envGroups)
			}
		})
	}
}

func TestBlueprintMethods(t *testing.T) {
	bp := &Blueprint{
		Services: []Service{
			{Name: "api", Type: ServiceTypeWeb},
			{Name: "worker", Type: ServiceTypeWorker},
		},
		Databases: []Database{
			{Name: "main-db"},
			{Name: "cache-db"},
		},
		EnvVarGroups: []EnvVarGroup{
			{Name: "shared"},
			{Name: "secrets"},
		},
	}

	// Test GetServices
	services := bp.GetServices()
	if len(services) != 2 || services[0].Name != "api" || services[1].Name != "worker" {
		t.Errorf("GetServices() failed: %v", services)
	}

	// Test GetDatabases
	databases := bp.GetDatabases()
	if len(databases) != 2 || databases[0].Name != "main-db" || databases[1].Name != "cache-db" {
		t.Errorf("GetDatabases() failed: %v", databases)
	}

	// Test GetEnvVarGroups
	envGroups := bp.GetEnvVarGroups()
	if len(envGroups) != 2 || envGroups[0].Name != "shared" || envGroups[1].Name != "secrets" {
		t.Errorf("GetEnvVarGroups() failed: %v", envGroups)
	}

	// Test FindService
	apiService := bp.FindService("api")
	if apiService == nil || apiService.Name != "api" {
		t.Errorf("FindService('api') failed: %v", apiService)
	}

	nonExistentService := bp.FindService("non-existent")
	if nonExistentService != nil {
		t.Errorf("FindService('non-existent') should return nil, got: %v", nonExistentService)
	}

	// Test FindDatabase
	mainDB := bp.FindDatabase("main-db")
	if mainDB == nil || mainDB.Name != "main-db" {
		t.Errorf("FindDatabase('main-db') failed: %v", mainDB)
	}

	nonExistentDB := bp.FindDatabase("non-existent")
	if nonExistentDB != nil {
		t.Errorf("FindDatabase('non-existent') should return nil, got: %v", nonExistentDB)
	}

	// Test FindEnvVarGroup
	sharedGroup := bp.FindEnvVarGroup("shared")
	if sharedGroup == nil || sharedGroup.Name != "shared" {
		t.Errorf("FindEnvVarGroup('shared') failed: %v", sharedGroup)
	}

	nonExistentGroup := bp.FindEnvVarGroup("non-existent")
	if nonExistentGroup != nil {
		t.Errorf("FindEnvVarGroup('non-existent') should return nil, got: %v", nonExistentGroup)
	}

	// Test methods with nil blueprint
	var nilBP *Blueprint

	if nilBP.GetServices() != nil {
		t.Errorf("nil blueprint GetServices() should return nil")
	}
	if nilBP.GetDatabases() != nil {
		t.Errorf("nil blueprint GetDatabases() should return nil")
	}
	if nilBP.GetEnvVarGroups() != nil {
		t.Errorf("nil blueprint GetEnvVarGroups() should return nil")
	}
	if nilBP.FindService("test") != nil {
		t.Errorf("nil blueprint FindService() should return nil")
	}
	if nilBP.FindDatabase("test") != nil {
		t.Errorf("nil blueprint FindDatabase() should return nil")
	}
	if nilBP.FindEnvVarGroup("test") != nil {
		t.Errorf("nil blueprint FindEnvVarGroup() should return nil")
	}
}

// Helper functions for tests

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func planPtr(p Plan) *Plan {
	return &p
}

func runtimePtr(r Runtime) *Runtime {
	return &r
}

func slicesEqual(a, b []string) bool {
	// Handle nil vs empty slice cases
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}

func blueprintsEqual(a, b *Blueprint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare services
	if len(a.Services) != len(b.Services) {
		return false
	}
	for i, svc := range a.Services {
		if !servicesEqual(svc, b.Services[i]) {
			return false
		}
	}

	// Compare databases
	if len(a.Databases) != len(b.Databases) {
		return false
	}
	for i, db := range a.Databases {
		if !databasesEqual(db, b.Databases[i]) {
			return false
		}
	}

	// Compare env groups
	if len(a.EnvVarGroups) != len(b.EnvVarGroups) {
		return false
	}
	for i, group := range a.EnvVarGroups {
		if !envGroupsEqual(group, b.EnvVarGroups[i]) {
			return false
		}
	}

	// Compare previews
	if (a.Previews == nil) != (b.Previews == nil) {
		return false
	}
	if a.Previews != nil && a.Previews.Generation != b.Previews.Generation {
		return false
	}

	// Compare preview expiration
	if (a.PreviewsExpireAfterDays == nil) != (b.PreviewsExpireAfterDays == nil) {
		return false
	}
	if a.PreviewsExpireAfterDays != nil && *a.PreviewsExpireAfterDays != *b.PreviewsExpireAfterDays {
		return false
	}

	return true
}

func servicesEqual(a, b Service) bool {
	// Compare basic fields
	if a.Name != b.Name || a.Type != b.Type {
		return false
	}

	// Compare runtime pointers
	if (a.Runtime == nil) != (b.Runtime == nil) {
		return false
	}
	if a.Runtime != nil && *a.Runtime != *b.Runtime {
		return false
	}

	// Compare env vars
	if len(a.EnvVars) != len(b.EnvVars) {
		return false
	}
	for i, envVar := range a.EnvVars {
		if !envVarsEqual(envVar, b.EnvVars[i]) {
			return false
		}
	}

	return true
}

func databasesEqual(a, b Database) bool {
	if a.Name != b.Name {
		return false
	}

	// Compare plan pointers
	if (a.Plan == nil) != (b.Plan == nil) {
		return false
	}
	if a.Plan != nil && *a.Plan != *b.Plan {
		return false
	}

	return true
}

func envGroupsEqual(a, b EnvVarGroup) bool {
	if a.Name != b.Name {
		return false
	}

	if len(a.EnvVars) != len(b.EnvVars) {
		return false
	}

	for i, envVar := range a.EnvVars {
		if !envVarsEqual(envVar, b.EnvVars[i]) {
			return false
		}
	}

	return true
}

func envVarsEqual(a, b EnvVar) bool {
	// Compare key pointers
	if (a.Key == nil) != (b.Key == nil) {
		return false
	}
	if a.Key != nil && *a.Key != *b.Key {
		return false
	}

	// Compare value pointers
	if (a.Value == nil) != (b.Value == nil) {
		return false
	}
	if a.Value != nil && *a.Value != *b.Value {
		return false
	}

	// Compare FromDatabase
	if (a.FromDatabase == nil) != (b.FromDatabase == nil) {
		return false
	}
	if a.FromDatabase != nil {
		if a.FromDatabase.Name != b.FromDatabase.Name || a.FromDatabase.Property != b.FromDatabase.Property {
			return false
		}
	}

	// Compare FromService
	if (a.FromService == nil) != (b.FromService == nil) {
		return false
	}
	if a.FromService != nil {
		if a.FromService.Name != b.FromService.Name || a.FromService.Type != b.FromService.Type {
			return false
		}
	}

	// Compare FromGroup
	if (a.FromGroup == nil) != (b.FromGroup == nil) {
		return false
	}
	if a.FromGroup != nil && *a.FromGroup != *b.FromGroup {
		return false
	}

	return true
}