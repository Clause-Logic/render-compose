package render

// String enum types
type ServiceType string
type Runtime string
type Plan string
type Region string
type PreviewGeneration string
type RouteType string
type MaxMemoryPolicy string
type DatabaseProperty string
type ServiceProperty string
type PostgreSQLVersion string

// Service Types
const (
	ServiceTypeWeb      ServiceType = "web"
	ServiceTypeWorker   ServiceType = "worker"
	ServiceTypePServ    ServiceType = "pserv"
	ServiceTypeCron     ServiceType = "cron"
	ServiceTypeKeyValue ServiceType = "keyvalue"
	ServiceTypeRedis    ServiceType = "redis" // deprecated alias
)

// Runtimes
const (
	RuntimeNode   Runtime = "node"
	RuntimePython Runtime = "python"
	RuntimeRuby   Runtime = "ruby"
	RuntimeGo     Runtime = "go"
	RuntimeRust   Runtime = "rust"
	RuntimeDocker Runtime = "docker"
	RuntimeStatic Runtime = "static"
	RuntimeImage  Runtime = "image"
)

// Plans
const (
	// Service plans
	PlanStarter     Plan = "starter"
	PlanStandard    Plan = "standard"
	PlanStandard2x  Plan = "standard-2x"
	PlanStandard4x  Plan = "standard-4x"
	PlanPro         Plan = "pro"
	PlanPro2x       Plan = "pro-2x"
	PlanPro4x       Plan = "pro-4x"
	PlanProMax      Plan = "pro-max"
	
	// Database plans
	PlanBasic256MB Plan = "basic-256mb"
	PlanBasic1GB   Plan = "basic-1gb"
	PlanBasic4GB   Plan = "basic-4gb"
	PlanPro8GB     Plan = "pro-8gb"
	PlanPro16GB    Plan = "pro-16gb"
	
	// Key Value plans
	PlanFree Plan = "free"
)

// Regions
const (
	RegionOregon    Region = "oregon"
	RegionVirginia  Region = "virginia"
	RegionFrankfurt Region = "frankfurt"
	RegionSingapore Region = "singapore"
)

// Preview Generation
const (
	PreviewGenerationAutomatic PreviewGeneration = "automatic"
	PreviewGenerationNone      PreviewGeneration = "none"
)

// Route Types
const (
	RouteTypeRedirect RouteType = "redirect"
	RouteTypeRewrite  RouteType = "rewrite"
)

// Max Memory Policies
const (
	MaxMemoryPolicyAllKeysLRU     MaxMemoryPolicy = "allkeys-lru"
	MaxMemoryPolicyAllKeysRandom  MaxMemoryPolicy = "allkeys-random"
	MaxMemoryPolicyVolatileLRU    MaxMemoryPolicy = "volatile-lru"
	MaxMemoryPolicyVolatileRandom MaxMemoryPolicy = "volatile-random"
	MaxMemoryPolicyVolatileTTL    MaxMemoryPolicy = "volatile-ttl"
	MaxMemoryPolicyNoEviction     MaxMemoryPolicy = "noeviction"
)

// Database Properties
const (
	DatabasePropertyConnectionString         DatabaseProperty = "connectionString"
	DatabasePropertyInternalConnectionString DatabaseProperty = "internalConnectionString"
	DatabasePropertyHost                     DatabaseProperty = "host"
	DatabasePropertyPort                     DatabaseProperty = "port"
	DatabasePropertyUser                     DatabaseProperty = "user"
	DatabasePropertyPassword                 DatabaseProperty = "password"
	DatabasePropertyDatabase                 DatabaseProperty = "database"
)

// Service Properties
const (
	ServicePropertyHost                     ServiceProperty = "host"
	ServicePropertyPort                     ServiceProperty = "port"
	ServicePropertyConnectionString         ServiceProperty = "connectionString"
	ServicePropertyInternalConnectionString ServiceProperty = "internalConnectionString"
)

// PostgreSQL Versions
const (
	PostgreSQL13 PostgreSQLVersion = "13"
	PostgreSQL14 PostgreSQLVersion = "14"
	PostgreSQL15 PostgreSQLVersion = "15"
	PostgreSQL16 PostgreSQLVersion = "16"
)

// Root Blueprint structure
type Blueprint struct {
	Services                []Service     `yaml:"services,omitempty"`
	Databases               []Database    `yaml:"databases,omitempty"`
	EnvVarGroups            []EnvVarGroup `yaml:"envVarGroups,omitempty"`
	Previews                *Previews     `yaml:"previews,omitempty"`
	PreviewsExpireAfterDays *int          `yaml:"previewsExpireAfterDays,omitempty"`
}

// Service types
type Service struct {
	// Essential fields
	Name string      `yaml:"name"`
	Type ServiceType `yaml:"type"`
	
	// Runtime (required unless keyvalue/redis)
	Runtime *Runtime `yaml:"runtime,omitempty"`
	
	// Instance type
	Plan *Plan `yaml:"plan,omitempty"`
	
	// Preview configuration
	Previews    *ServicePreviews `yaml:"previews,omitempty"`
	PreviewPlan *Plan            `yaml:"previewPlan,omitempty"`
	
	// Build commands
	BuildCommand     *string `yaml:"buildCommand,omitempty"`
	StartCommand     *string `yaml:"startCommand,omitempty"`
	PreDeployCommand *string `yaml:"preDeployCommand,omitempty"`
	
	// Git configuration
	Repo   *string `yaml:"repo,omitempty"`
	Branch *string `yaml:"branch,omitempty"`
	
	// Deployment
	AutoDeploy              *bool `yaml:"autoDeploy,omitempty"`
	MaxShutdownDelaySeconds *int  `yaml:"maxShutdownDelaySeconds,omitempty"`
	
	// Web service specific
	Domains []string `yaml:"domains,omitempty"`
	
	// Region
	Region *Region `yaml:"region,omitempty"`
	
	// Scaling
	NumInstances *int     `yaml:"numInstances,omitempty"`
	Scaling      *Scaling `yaml:"scaling,omitempty"`
	
	// Environment variables
	EnvVars []EnvVar `yaml:"envVars,omitempty"`
	
	// Docker specific
	DockerCommand      *string             `yaml:"dockerCommand,omitempty"`
	DockerfilePath     *string             `yaml:"dockerfilePath,omitempty"`
	DockerContext      *string             `yaml:"dockerContext,omitempty"`
	Image              *DockerImage        `yaml:"image,omitempty"`
	RegistryCredential *RegistryCredential `yaml:"registryCredential,omitempty"`
	
	// Build configuration
	BuildFilter *BuildFilter `yaml:"buildFilter,omitempty"`
	RootDir     *string      `yaml:"rootDir,omitempty"`
	
	// Persistent disk
	Disk *Disk `yaml:"disk,omitempty"`
	
	// Static site specific
	StaticPublishPath *string  `yaml:"staticPublishPath,omitempty"`
	Headers           []Header `yaml:"headers,omitempty"`
	Routes            []Route  `yaml:"routes,omitempty"`
	
	// Cron specific
	Schedule *string `yaml:"schedule,omitempty"`
	
	// Key Value specific
	IPAllowList     []IPAllow        `yaml:"ipAllowList,omitempty"`
	MaxMemoryPolicy *MaxMemoryPolicy `yaml:"maxmemoryPolicy,omitempty"`
	
	// Health check
	HealthCheckPath *string `yaml:"healthCheckPath,omitempty"`
}

// Database configuration
type Database struct {
	// Essential
	Name string `yaml:"name"`
	
	// Instance configuration
	Plan              *Plan   `yaml:"plan,omitempty"`
	PreviewPlan       *Plan   `yaml:"previewPlan,omitempty"`
	DiskSizeGB        *int    `yaml:"diskSizeGB,omitempty"`
	PreviewDiskSizeGB *int    `yaml:"previewDiskSizeGB,omitempty"`
	Region            *Region `yaml:"region,omitempty"`
	
	// PostgreSQL specific
	PostgresMajorVersion *PostgreSQLVersion `yaml:"postgresMajorVersion,omitempty"`
	DatabaseName         *string            `yaml:"databaseName,omitempty"`
	User                 *string            `yaml:"user,omitempty"`
	
	// Access control
	IPAllowList []IPAllow `yaml:"ipAllowList,omitempty"`
	
	// High availability and replicas
	ReadReplicas     []ReadReplica     `yaml:"readReplicas,omitempty"`
	HighAvailability *HighAvailability `yaml:"highAvailability,omitempty"`
}

// Environment variable configuration
type EnvVar struct {
	Key           *string        `yaml:"key,omitempty"`
	Value         *string        `yaml:"value,omitempty"`
	GenerateValue *bool          `yaml:"generateValue,omitempty"`
	Sync          *bool          `yaml:"sync,omitempty"`
	FromDatabase  *FromDatabase  `yaml:"fromDatabase,omitempty"`
	FromService   *FromService   `yaml:"fromService,omitempty"`
	FromGroup     *string        `yaml:"fromGroup,omitempty"`
}

// Environment variable group
type EnvVarGroup struct {
	Name    string   `yaml:"name"`
	EnvVars []EnvVar `yaml:"envVars,omitempty"`
}

// Reference to database property
type FromDatabase struct {
	Name     string           `yaml:"name"`
	Property DatabaseProperty `yaml:"property"`
}

// Reference to service property
type FromService struct {
	Name      string           `yaml:"name"`
	Type      ServiceType      `yaml:"type"`
	Property  *ServiceProperty `yaml:"property,omitempty"`
	EnvVarKey *string          `yaml:"envVarKey,omitempty"`
}

// Scaling configuration
type Scaling struct {
	MinInstances         *int `yaml:"minInstances,omitempty"`
	MaxInstances         *int `yaml:"maxInstances,omitempty"`
	TargetMemoryPercent  *int `yaml:"targetMemoryPercent,omitempty"`
	TargetCPUPercent     *int `yaml:"targetCPUPercent,omitempty"`
}

// Docker image configuration
type DockerImage struct {
	URL          string            `yaml:"url"`
	Credentials  *ImageCredentials `yaml:"credentials,omitempty"`
}

type ImageCredentials struct {
	FromRegistryCreds *RegistryCredsRef `yaml:"fromRegistryCreds,omitempty"`
}

type RegistryCredsRef struct {
	Name string `yaml:"name"`
}

// Registry credential
type RegistryCredential struct {
	FromRegistryCreds *RegistryCredsRef `yaml:"fromRegistryCreds,omitempty"`
}

// Build filter
type BuildFilter struct {
	Paths        []string `yaml:"paths,omitempty"`
	IgnoredPaths []string `yaml:"ignoredPaths,omitempty"`
}

// Persistent disk
type Disk struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
	SizeGB    *int   `yaml:"sizeGB,omitempty"`
}

// HTTP headers for static sites
type Header struct {
	Path  string `yaml:"path"`
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Routes for static sites
type Route struct {
	Type        string `yaml:"type"` // redirect, rewrite
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

// IP allow list entry
type IPAllow struct {
	Source      string  `yaml:"source"`
	Description *string `yaml:"description,omitempty"`
}

// Read replica configuration
type ReadReplica struct {
	Name string `yaml:"name"`
}

// High availability configuration
type HighAvailability struct {
	Enabled bool `yaml:"enabled"`
}

// Preview environment configuration
type Previews struct {
	Generation string `yaml:"generation"` // automatic, none
}

// Service-specific preview configuration
type ServicePreviews struct {
	Generation string `yaml:"generation"` // automatic, none
}