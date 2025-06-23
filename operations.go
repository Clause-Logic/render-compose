package render

import (
	"fmt"
	"strings"
)

// MergeBlueprints combines two blueprints into one
// Returns an error if there are any name conflicts
// Use PrefixBlueprint() to avoid conflicts before merging
func MergeBlueprints(base, overlay *Blueprint) (*Blueprint, error) {
	if base == nil && overlay == nil {
		return &Blueprint{}, nil
	}
	if base == nil {
		return CopyBlueprint(overlay), nil
	}
	if overlay == nil {
		return CopyBlueprint(base), nil
	}

	// Check for conflicts first
	conflicts := FindConflicts(base, overlay)
	if len(conflicts) > 0 {
		return nil, fmt.Errorf("merge conflicts found: %s. Use PrefixBlueprint() to avoid name conflicts before merging", strings.Join(conflicts, ", "))
	}

	merged := &Blueprint{}

	// Combine services
	merged.Services = append(merged.Services, base.Services...)
	merged.Services = append(merged.Services, overlay.Services...)

	// Combine databases
	merged.Databases = append(merged.Databases, base.Databases...)
	merged.Databases = append(merged.Databases, overlay.Databases...)

	// Combine environment variable groups
	merged.EnvVarGroups = append(merged.EnvVarGroups, base.EnvVarGroups...)
	merged.EnvVarGroups = append(merged.EnvVarGroups, overlay.EnvVarGroups...)

	// Overlay wins for preview configuration
	if overlay.Previews != nil {
		merged.Previews = overlay.Previews
	} else {
		merged.Previews = base.Previews
	}

	// Overlay wins for preview expiration
	if overlay.PreviewsExpireAfterDays != nil {
		merged.PreviewsExpireAfterDays = overlay.PreviewsExpireAfterDays
	} else {
		merged.PreviewsExpireAfterDays = base.PreviewsExpireAfterDays
	}

	return merged, nil
}

// CopyBlueprint creates a deep copy of a blueprint
func CopyBlueprint(bp *Blueprint) *Blueprint {
	if bp == nil {
		return &Blueprint{}
	}

	copied := &Blueprint{}

	// Copy services
	copied.Services = make([]Service, len(bp.Services))
	copy(copied.Services, bp.Services)

	// Copy databases
	copied.Databases = make([]Database, len(bp.Databases))
	copy(copied.Databases, bp.Databases)

	// Copy environment variable groups
	copied.EnvVarGroups = make([]EnvVarGroup, len(bp.EnvVarGroups))
	copy(copied.EnvVarGroups, bp.EnvVarGroups)

	// Copy preview configuration
	if bp.Previews != nil {
		previews := *bp.Previews
		copied.Previews = &previews
	}

	// Copy preview expiration
	if bp.PreviewsExpireAfterDays != nil {
		expireDays := *bp.PreviewsExpireAfterDays
		copied.PreviewsExpireAfterDays = &expireDays
	}

	return copied
}

// ValidateBlueprint checks for common issues in a blueprint
func ValidateBlueprint(bp *Blueprint) []string {
	var errors []string

	if bp == nil {
		return []string{"blueprint is nil"}
	}

	// Check for duplicate service names
	serviceNames := make(map[string]bool)
	for _, service := range bp.Services {
		if serviceNames[service.Name] {
			errors = append(errors, fmt.Sprintf("duplicate service name: %s", service.Name))
		}
		serviceNames[service.Name] = true
	}

	// Check for duplicate database names
	dbNames := make(map[string]bool)
	for _, db := range bp.Databases {
		if dbNames[db.Name] {
			errors = append(errors, fmt.Sprintf("duplicate database name: %s", db.Name))
		}
		dbNames[db.Name] = true
	}

	// Check for duplicate environment group names
	envGroupNames := make(map[string]bool)
	for _, group := range bp.EnvVarGroups {
		if envGroupNames[group.Name] {
			errors = append(errors, fmt.Sprintf("duplicate environment group name: %s", group.Name))
		}
		envGroupNames[group.Name] = true
	}

	// Check for missing required fields
	for _, service := range bp.Services {
		if service.Name == "" {
			errors = append(errors, "service missing name")
		}
		if service.Type == "" {
			errors = append(errors, fmt.Sprintf("service %s missing type", service.Name))
		}
		// Runtime required for most service types
		if service.Runtime == nil && service.Type != ServiceTypeKeyValue {
			errors = append(errors, fmt.Sprintf("service %s missing runtime", service.Name))
		}
	}

	for _, db := range bp.Databases {
		if db.Name == "" {
			errors = append(errors, "database missing name")
		}
	}

	for _, group := range bp.EnvVarGroups {
		if group.Name == "" {
			errors = append(errors, "environment group missing name")
		}
	}

	return errors
}

// FindConflicts identifies name conflicts between two blueprints
func FindConflicts(base, overlay *Blueprint) []string {
	var conflicts []string

	if base == nil || overlay == nil {
		return conflicts
	}

	// Check service name conflicts
	baseServiceNames := make(map[string]bool)
	for _, service := range base.Services {
		baseServiceNames[service.Name] = true
	}
	for _, service := range overlay.Services {
		if baseServiceNames[service.Name] {
			conflicts = append(conflicts, fmt.Sprintf("service name conflict: %s", service.Name))
		}
	}

	// Check database name conflicts
	baseDBNames := make(map[string]bool)
	for _, db := range base.Databases {
		baseDBNames[db.Name] = true
	}
	for _, db := range overlay.Databases {
		if baseDBNames[db.Name] {
			conflicts = append(conflicts, fmt.Sprintf("database name conflict: %s", db.Name))
		}
	}

	// Check environment group name conflicts
	baseEnvGroupNames := make(map[string]bool)
	for _, group := range base.EnvVarGroups {
		baseEnvGroupNames[group.Name] = true
	}
	for _, group := range overlay.EnvVarGroups {
		if baseEnvGroupNames[group.Name] {
			conflicts = append(conflicts, fmt.Sprintf("environment group name conflict: %s", group.Name))
		}
	}

	return conflicts
}

// PrefixBlueprint adds a prefix to all named resources and updates internal references
// External references (to resources not defined in this blueprint) are left unchanged
func PrefixBlueprint(bp *Blueprint, prefix string) *Blueprint {
	if bp == nil || prefix == "" {
		return CopyBlueprint(bp)
	}

	// Create a deep copy to avoid modifying the original
	prefixed := CopyBlueprint(bp)

	// Collect all names that exist in this blueprint
	existingServiceNames := make(map[string]bool)
	existingDatabaseNames := make(map[string]bool)
	existingEnvGroupNames := make(map[string]bool)

	for _, service := range prefixed.Services {
		existingServiceNames[service.Name] = true
	}
	for _, db := range prefixed.Databases {
		existingDatabaseNames[db.Name] = true
	}
	for _, group := range prefixed.EnvVarGroups {
		existingEnvGroupNames[group.Name] = true
	}

	// Create mapping of old names to new names
	serviceNameMap := make(map[string]string)
	databaseNameMap := make(map[string]string)
	envGroupNameMap := make(map[string]string)

	for oldName := range existingServiceNames {
		serviceNameMap[oldName] = prefix + oldName
	}
	for oldName := range existingDatabaseNames {
		databaseNameMap[oldName] = prefix + oldName
	}
	for oldName := range existingEnvGroupNames {
		envGroupNameMap[oldName] = prefix + oldName
	}

	// Update service names
	for i := range prefixed.Services {
		if newName, exists := serviceNameMap[prefixed.Services[i].Name]; exists {
			prefixed.Services[i].Name = newName
		}
	}

	// Update database names
	for i := range prefixed.Databases {
		if newName, exists := databaseNameMap[prefixed.Databases[i].Name]; exists {
			prefixed.Databases[i].Name = newName
		}
	}

	// Update environment group names
	for i := range prefixed.EnvVarGroups {
		if newName, exists := envGroupNameMap[prefixed.EnvVarGroups[i].Name]; exists {
			prefixed.EnvVarGroups[i].Name = newName
		}
	}

	// Update all internal references in environment variables
	updateEnvVarReferences := func(envVars []EnvVar) {
		for i := range envVars {
			envVar := &envVars[i]

			// Update database references
			if envVar.FromDatabase != nil {
				if newName, exists := databaseNameMap[envVar.FromDatabase.Name]; exists {
					envVar.FromDatabase.Name = newName
				}
			}

			// Update service references
			if envVar.FromService != nil {
				if newName, exists := serviceNameMap[envVar.FromService.Name]; exists {
					envVar.FromService.Name = newName
				}
			}

			// Update environment group references
			if envVar.FromGroup != nil {
				if newName, exists := envGroupNameMap[*envVar.FromGroup]; exists {
					*envVar.FromGroup = newName
				}
			}
		}
	}

	// Update references in service environment variables
	for i := range prefixed.Services {
		updateEnvVarReferences(prefixed.Services[i].EnvVars)
	}

	// Update references in environment group variables
	for i := range prefixed.EnvVarGroups {
		updateEnvVarReferences(prefixed.EnvVarGroups[i].EnvVars)
	}

	// Update read replica names that reference the parent database
	for i := range prefixed.Databases {
		db := &prefixed.Databases[i]
		for j := range db.ReadReplicas {
			replica := &db.ReadReplicas[j]
			// Check if replica name follows the pattern of referencing the parent database
			if strings.HasPrefix(replica.Name, db.Name[:len(db.Name)-len(prefix)]) {
				// Update replica name to match the new database name
				oldDBName := db.Name[:len(db.Name)-len(prefix)]
				if strings.HasPrefix(replica.Name, oldDBName) {
					replica.Name = strings.Replace(replica.Name, oldDBName, db.Name, 1)
				}
			}
		}
	}

	return prefixed
}

// PrefixBlueprintWithSeparator adds a prefix with a separator to all named resources
func PrefixBlueprintWithSeparator(bp *Blueprint, prefix, separator string) *Blueprint {
	if separator == "" {
		separator = "-"
	}
	return PrefixBlueprint(bp, prefix+separator)
}

// GetAllResourceNames returns all resource names in a blueprint
func GetAllResourceNames(bp *Blueprint) (services, databases, envGroups []string) {
	if bp == nil {
		return nil, nil, nil
	}

	for _, service := range bp.Services {
		services = append(services, service.Name)
	}
	for _, db := range bp.Databases {
		databases = append(databases, db.Name)
	}
	for _, group := range bp.EnvVarGroups {
		envGroups = append(envGroups, group.Name)
	}

	return services, databases, envGroups
}

// GetExternalReferences returns references to resources not defined in this blueprint
func GetExternalReferences(bp *Blueprint) (services, databases, envGroups []string) {
	if bp == nil {
		return nil, nil, nil
	}

	// Collect all names that exist in this blueprint
	existingServiceNames := make(map[string]bool)
	existingDatabaseNames := make(map[string]bool)
	existingEnvGroupNames := make(map[string]bool)

	for _, service := range bp.Services {
		existingServiceNames[service.Name] = true
	}
	for _, db := range bp.Databases {
		existingDatabaseNames[db.Name] = true
	}
	for _, group := range bp.EnvVarGroups {
		existingEnvGroupNames[group.Name] = true
	}

	// Track external references
	externalServices := make(map[string]bool)
	externalDatabases := make(map[string]bool)
	externalEnvGroups := make(map[string]bool)

	checkEnvVarReferences := func(envVars []EnvVar) {
		for _, envVar := range envVars {
			// Check database references
			if envVar.FromDatabase != nil {
				if !existingDatabaseNames[envVar.FromDatabase.Name] {
					externalDatabases[envVar.FromDatabase.Name] = true
				}
			}

			// Check service references
			if envVar.FromService != nil {
				if !existingServiceNames[envVar.FromService.Name] {
					externalServices[envVar.FromService.Name] = true
				}
			}

			// Check environment group references
			if envVar.FromGroup != nil {
				if !existingEnvGroupNames[*envVar.FromGroup] {
					externalEnvGroups[*envVar.FromGroup] = true
				}
			}
		}
	}

	// Check references in service environment variables
	for _, service := range bp.Services {
		checkEnvVarReferences(service.EnvVars)
	}

	// Check references in environment group variables
	for _, group := range bp.EnvVarGroups {
		checkEnvVarReferences(group.EnvVars)
	}

	// Convert maps to slices
	for name := range externalServices {
		services = append(services, name)
	}
	for name := range externalDatabases {
		databases = append(databases, name)
	}
	for name := range externalEnvGroups {
		envGroups = append(envGroups, name)
	}

	return services, databases, envGroups
}

// GetServices returns all services from a blueprint (helper function)
func (bp *Blueprint) GetServices() []Service {
	if bp == nil {
		return nil
	}
	return bp.Services
}

// GetDatabases returns all databases from a blueprint (helper function)
func (bp *Blueprint) GetDatabases() []Database {
	if bp == nil {
		return nil
	}
	return bp.Databases
}

// GetEnvVarGroups returns all environment variable groups from a blueprint (helper function)
func (bp *Blueprint) GetEnvVarGroups() []EnvVarGroup {
	if bp == nil {
		return nil
	}
	return bp.EnvVarGroups
}

// FindService finds a service by name
func (bp *Blueprint) FindService(name string) *Service {
	if bp == nil {
		return nil
	}
	for i, service := range bp.Services {
		if service.Name == name {
			return &bp.Services[i]
		}
	}
	return nil
}

// FindDatabase finds a database by name
func (bp *Blueprint) FindDatabase(name string) *Database {
	if bp == nil {
		return nil
	}
	for i, db := range bp.Databases {
		if db.Name == name {
			return &bp.Databases[i]
		}
	}
	return nil
}

// FindEnvVarGroup finds an environment variable group by name
func (bp *Blueprint) FindEnvVarGroup(name string) *EnvVarGroup {
	if bp == nil {
		return nil
	}
	for i, group := range bp.EnvVarGroups {
		if group.Name == name {
			return &bp.EnvVarGroups[i]
		}
	}
	return nil
}

// findAvailableName generates a unique name by appending a number
func findAvailableName(baseName string, existingNames map[string]bool) string {
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s-%d", baseName, i)
		if !existingNames[candidate] {
			return candidate
		}
	}
}