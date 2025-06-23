package render

import (
	"gopkg.in/yaml.v3"
)

// ServiceMarshalable interface for services that need custom marshaling
type ServiceMarshalable interface {
	MarshalYAML() (interface{}, error)
}

// CustomBlueprint wraps Blueprint to provide custom marshaling
type CustomBlueprint struct {
	*Blueprint
	RawServices []ServiceMarshalable `yaml:"-"`
}

// MarshalYAML implements custom YAML marshaling for Blueprint to handle different service types
func (bp *Blueprint) MarshalYAML() (interface{}, error) {
	type Alias Blueprint
	
	// Create a map to hold the final structure
	result := make(map[string]interface{})
	
	// Marshal the blueprint without services first
	temp := &Alias{
		Databases:               bp.Databases,
		EnvVarGroups:            bp.EnvVarGroups,
		Previews:                bp.Previews,
		PreviewsExpireAfterDays: bp.PreviewsExpireAfterDays,
	}
	
	// Convert to map
	data, err := yaml.Marshal(temp)
	if err != nil {
		return nil, err
	}
	
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	
	// Handle services separately
	if len(bp.Services) > 0 {
		services := make([]interface{}, len(bp.Services))
		for i, service := range bp.Services {
			// Check if this is a static site (web + static runtime + has staticPublishPath)
			if service.Type == ServiceTypeWeb && 
			   service.Runtime != nil && 
			   *service.Runtime == RuntimeStatic && 
			   service.StaticPublishPath != nil {
				// Marshal as staticService format
				staticData := map[string]interface{}{
					"name":    service.Name,
					"type":    "web",
					"runtime": "static",
				}
				
				// Add optional fields
				if service.BuildCommand != nil {
					staticData["buildCommand"] = *service.BuildCommand
				}
				if service.StaticPublishPath != nil {
					staticData["staticPublishPath"] = *service.StaticPublishPath
				}
				if service.Repo != nil {
					staticData["repo"] = *service.Repo
				}
				if service.Branch != nil {
					staticData["branch"] = *service.Branch
				}
				if len(service.Domains) > 0 {
					staticData["domains"] = service.Domains
				}
				// Note: region is not supported for static services in the Render schema
				if len(service.Headers) > 0 {
					staticData["headers"] = service.Headers
				}
				if len(service.Routes) > 0 {
					staticData["routes"] = service.Routes
				}
				if service.AutoDeploy != nil {
					staticData["autoDeploy"] = *service.AutoDeploy
				}
				if service.BuildFilter != nil {
					staticData["buildFilter"] = service.BuildFilter
				}
				if service.RootDir != nil {
					staticData["rootDir"] = *service.RootDir
				}
				if len(service.EnvVars) > 0 {
					staticData["envVars"] = service.EnvVars
				}
				if service.Previews != nil {
					staticData["previews"] = service.Previews
				}
				
				services[i] = staticData
			} else {
				// Marshal as regular service
				services[i] = service
			}
		}
		result["services"] = services
	}
	
	return result, nil
}