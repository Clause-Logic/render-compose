package render

// Configuration abstractions for related fields

// DockerConfig groups Docker-related configuration
type DockerConfig struct {
	DockerCommand      *string             `yaml:"dockerCommand,omitempty"`
	DockerfilePath     *string             `yaml:"dockerfilePath,omitempty"`
	DockerContext      *string             `yaml:"dockerContext,omitempty"`
	Image              *DockerImage        `yaml:"image,omitempty"`
	RegistryCredential *RegistryCredential `yaml:"registryCredential,omitempty"`
}

// GitConfig groups Git repository configuration
type GitConfig struct {
	Repo   *string `yaml:"repo,omitempty"`
	Branch *string `yaml:"branch,omitempty"`
}

// BuildConfig groups build and deployment configuration
type BuildConfig struct {
	BuildCommand     *string      `yaml:"buildCommand,omitempty"`
	PreDeployCommand *string      `yaml:"preDeployCommand,omitempty"`
	BuildFilter      *BuildFilter `yaml:"buildFilter,omitempty"`
	RootDir          *string      `yaml:"rootDir,omitempty"`
	AutoDeploy       *bool        `yaml:"autoDeploy,omitempty"`
}

// ScalingConfig groups scaling-related configuration
type ScalingConfig struct {
	NumInstances *int     `yaml:"numInstances,omitempty"`
	Scaling      *Scaling `yaml:"scaling,omitempty"`
}

// PreviewConfig groups preview environment configuration
type PreviewConfig struct {
	Previews    *ServicePreviews `yaml:"previews,omitempty"`
	PreviewPlan *Plan            `yaml:"previewPlan,omitempty"`
}

// StaticSiteConfig groups static site specific configuration
type StaticSiteConfig struct {
	StaticPublishPath string   `yaml:"staticPublishPath"`
	Headers           []Header `yaml:"headers,omitempty"`
	Routes            []Route  `yaml:"routes,omitempty"`
}

// KeyValueConfig groups key-value store specific configuration
type KeyValueConfig struct {
	IPAllowList     []IPAllow        `yaml:"ipAllowList"`
	MaxMemoryPolicy *MaxMemoryPolicy `yaml:"maxmemoryPolicy,omitempty"`
}

// WebService represents a web service with HTTP endpoints
type WebService struct {
	// Essential
	Name    string  `yaml:"name"`
	Runtime Runtime `yaml:"runtime"`

	// Web-specific
	Domains         []string `yaml:"domains,omitempty"`
	HealthCheckPath *string  `yaml:"healthCheckPath,omitempty"`
	StartCommand    *string  `yaml:"startCommand,omitempty"`

	// Infrastructure
	Plan   *Plan   `yaml:"plan,omitempty"`
	Region *Region `yaml:"region,omitempty"`

	// Configuration groups
	Git                     *GitConfig     `yaml:",inline,omitempty"`
	Build                   *BuildConfig   `yaml:",inline,omitempty"`
	Docker                  *DockerConfig  `yaml:",inline,omitempty"`
	Scaling                 *ScalingConfig `yaml:",inline,omitempty"`
	Preview                 *PreviewConfig `yaml:",inline,omitempty"`
	EnvVars                 []EnvVar       `yaml:"envVars,omitempty"`
	MaxShutdownDelaySeconds *int           `yaml:"maxShutdownDelaySeconds,omitempty"`
	Disk                    *Disk          `yaml:"disk,omitempty"`
}

// ToService converts WebService to generic Service
func (ws *WebService) ToService() *Service {
	service := &Service{
		Name:                    ws.Name,
		Type:                    ServiceTypeWeb,
		Runtime:                 &ws.Runtime,
		Plan:                    ws.Plan,
		StartCommand:            ws.StartCommand,
		Domains:                 ws.Domains,
		Region:                  ws.Region,
		EnvVars:                 ws.EnvVars,
		MaxShutdownDelaySeconds: ws.MaxShutdownDelaySeconds,
		Disk:                    ws.Disk,
		HealthCheckPath:         ws.HealthCheckPath,
	}

	// Apply Git configuration
	if ws.Git != nil {
		service.Repo = ws.Git.Repo
		service.Branch = ws.Git.Branch
	}

	// Apply Build configuration
	if ws.Build != nil {
		service.BuildCommand = ws.Build.BuildCommand
		service.PreDeployCommand = ws.Build.PreDeployCommand
		service.BuildFilter = ws.Build.BuildFilter
		service.RootDir = ws.Build.RootDir
		service.AutoDeploy = ws.Build.AutoDeploy
	}

	// Apply Docker configuration
	if ws.Docker != nil {
		service.DockerCommand = ws.Docker.DockerCommand
		service.DockerfilePath = ws.Docker.DockerfilePath
		service.DockerContext = ws.Docker.DockerContext
		service.Image = ws.Docker.Image
		service.RegistryCredential = ws.Docker.RegistryCredential
	}

	// Apply Scaling configuration
	if ws.Scaling != nil {
		service.NumInstances = ws.Scaling.NumInstances
		service.Scaling = ws.Scaling.Scaling
	}

	// Apply Preview configuration
	if ws.Preview != nil {
		service.Previews = ws.Preview.Previews
		service.PreviewPlan = ws.Preview.PreviewPlan
	}

	return service
}

// BackgroundWorker represents a background worker service
type BackgroundWorker struct {
	// Essential
	Name         string  `yaml:"name"`
	Runtime      Runtime `yaml:"runtime"`
	StartCommand *string `yaml:"startCommand,omitempty"`

	// Infrastructure
	Plan   *Plan   `yaml:"plan,omitempty"`
	Region *Region `yaml:"region,omitempty"`

	// Configuration groups
	Git                     *GitConfig     `yaml:",inline,omitempty"`
	Build                   *BuildConfig   `yaml:",inline,omitempty"`
	Docker                  *DockerConfig  `yaml:",inline,omitempty"`
	Preview                 *PreviewConfig `yaml:",inline,omitempty"`
	EnvVars                 []EnvVar       `yaml:"envVars,omitempty"`
	MaxShutdownDelaySeconds *int           `yaml:"maxShutdownDelaySeconds,omitempty"`
	Disk                    *Disk          `yaml:"disk,omitempty"`
}

// ToService converts BackgroundWorker to generic Service
func (bw *BackgroundWorker) ToService() *Service {
	service := &Service{
		Name:                    bw.Name,
		Type:                    ServiceTypeWorker,
		Runtime:                 &bw.Runtime,
		Plan:                    bw.Plan,
		StartCommand:            bw.StartCommand,
		Region:                  bw.Region,
		EnvVars:                 bw.EnvVars,
		MaxShutdownDelaySeconds: bw.MaxShutdownDelaySeconds,
		Disk:                    bw.Disk,
	}

	// Apply Git configuration
	if bw.Git != nil {
		service.Repo = bw.Git.Repo
		service.Branch = bw.Git.Branch
	}

	// Apply Build configuration
	if bw.Build != nil {
		service.BuildCommand = bw.Build.BuildCommand
		service.PreDeployCommand = bw.Build.PreDeployCommand
		service.BuildFilter = bw.Build.BuildFilter
		service.RootDir = bw.Build.RootDir
		service.AutoDeploy = bw.Build.AutoDeploy
	}

	// Apply Docker configuration
	if bw.Docker != nil {
		service.DockerCommand = bw.Docker.DockerCommand
		service.DockerfilePath = bw.Docker.DockerfilePath
		service.DockerContext = bw.Docker.DockerContext
		service.Image = bw.Docker.Image
		service.RegistryCredential = bw.Docker.RegistryCredential
	}

	// Apply Preview configuration
	if bw.Preview != nil {
		service.Previews = bw.Preview.Previews
		service.PreviewPlan = bw.Preview.PreviewPlan
	}

	return service
}

// PrivateService represents a private service
type PrivateService struct {
	// Essential
	Name         string  `yaml:"name"`
	Runtime      Runtime `yaml:"runtime"`
	StartCommand *string `yaml:"startCommand,omitempty"`

	// Infrastructure
	Plan   *Plan   `yaml:"plan,omitempty"`
	Region *Region `yaml:"region,omitempty"`

	// Configuration groups
	Git                     *GitConfig     `yaml:",inline,omitempty"`
	Build                   *BuildConfig   `yaml:",inline,omitempty"`
	Docker                  *DockerConfig  `yaml:",inline,omitempty"`
	Preview                 *PreviewConfig `yaml:",inline,omitempty"`
	EnvVars                 []EnvVar       `yaml:"envVars,omitempty"`
	MaxShutdownDelaySeconds *int           `yaml:"maxShutdownDelaySeconds,omitempty"`
	Disk                    *Disk          `yaml:"disk,omitempty"`
}

// ToService converts PrivateService to generic Service
func (ps *PrivateService) ToService() *Service {
	service := &Service{
		Name:                    ps.Name,
		Type:                    ServiceTypePServ,
		Runtime:                 &ps.Runtime,
		Plan:                    ps.Plan,
		StartCommand:            ps.StartCommand,
		Region:                  ps.Region,
		EnvVars:                 ps.EnvVars,
		MaxShutdownDelaySeconds: ps.MaxShutdownDelaySeconds,
		Disk:                    ps.Disk,
	}

	// Apply Git configuration
	if ps.Git != nil {
		service.Repo = ps.Git.Repo
		service.Branch = ps.Git.Branch
	}

	// Apply Build configuration
	if ps.Build != nil {
		service.BuildCommand = ps.Build.BuildCommand
		service.PreDeployCommand = ps.Build.PreDeployCommand
		service.BuildFilter = ps.Build.BuildFilter
		service.RootDir = ps.Build.RootDir
		service.AutoDeploy = ps.Build.AutoDeploy
	}

	// Apply Docker configuration
	if ps.Docker != nil {
		service.DockerCommand = ps.Docker.DockerCommand
		service.DockerfilePath = ps.Docker.DockerfilePath
		service.DockerContext = ps.Docker.DockerContext
		service.Image = ps.Docker.Image
		service.RegistryCredential = ps.Docker.RegistryCredential
	}

	// Apply Preview configuration
	if ps.Preview != nil {
		service.Previews = ps.Preview.Previews
		service.PreviewPlan = ps.Preview.PreviewPlan
	}

	return service
}

// CronJob represents a scheduled job
type CronJob struct {
	// Essential
	Name         string  `yaml:"name"`
	Runtime      Runtime `yaml:"runtime"`
	Schedule     string  `yaml:"schedule"`
	StartCommand *string `yaml:"startCommand,omitempty"`

	// Infrastructure
	Region *Region `yaml:"region,omitempty"`

	// Configuration groups (no scaling for cron jobs)
	Git     *GitConfig     `yaml:",inline,omitempty"`
	Build   *BuildConfig   `yaml:",inline,omitempty"`
	Docker  *DockerConfig  `yaml:",inline,omitempty"`
	Preview *PreviewConfig `yaml:",inline,omitempty"`
	EnvVars []EnvVar       `yaml:"envVars,omitempty"`
}

// ToService converts CronJob to generic Service
func (cj *CronJob) ToService() *Service {
	service := &Service{
		Name:         cj.Name,
		Type:         ServiceTypeCron,
		Runtime:      &cj.Runtime,
		StartCommand: cj.StartCommand,
		Region:       cj.Region,
		EnvVars:      cj.EnvVars,
		Schedule:     &cj.Schedule,
	}

	// Apply Git configuration
	if cj.Git != nil {
		service.Repo = cj.Git.Repo
		service.Branch = cj.Git.Branch
	}

	// Apply Build configuration
	if cj.Build != nil {
		service.BuildCommand = cj.Build.BuildCommand
		service.PreDeployCommand = cj.Build.PreDeployCommand
		service.BuildFilter = cj.Build.BuildFilter
		service.RootDir = cj.Build.RootDir
		service.AutoDeploy = cj.Build.AutoDeploy
	}

	// Apply Docker configuration
	if cj.Docker != nil {
		service.DockerCommand = cj.Docker.DockerCommand
		service.DockerfilePath = cj.Docker.DockerfilePath
		service.DockerContext = cj.Docker.DockerContext
		service.Image = cj.Docker.Image
		service.RegistryCredential = cj.Docker.RegistryCredential
	}

	// Apply Preview configuration
	if cj.Preview != nil {
		service.Previews = cj.Preview.Previews
		service.PreviewPlan = cj.Preview.PreviewPlan
	}

	return service
}

// StaticSite represents a static website
type StaticSite struct {
	// Essential
	Name string `yaml:"name"`

	// Infrastructure
	Region *Region `yaml:"region,omitempty"`

	// Configuration groups
	Git        *GitConfig        `yaml:",inline,omitempty"`
	Build      *BuildConfig      `yaml:",inline,omitempty"`
	StaticSite *StaticSiteConfig `yaml:",inline"`
	Preview    *PreviewConfig    `yaml:",inline,omitempty"`
	Domains    []string          `yaml:"domains,omitempty"`
}

// ToService converts StaticSite to generic Service
func (ss *StaticSite) ToService() *Service {
	runtime := RuntimeStatic
	service := &Service{
		Name:    ss.Name,
		Type:    ServiceTypeWeb,
		Runtime: &runtime,
		Domains: ss.Domains,
		Region:  ss.Region,
	}

	// Apply Git configuration
	if ss.Git != nil {
		service.Repo = ss.Git.Repo
		service.Branch = ss.Git.Branch
	}

	// Apply Build configuration
	if ss.Build != nil {
		service.BuildCommand = ss.Build.BuildCommand
		service.PreDeployCommand = ss.Build.PreDeployCommand
		service.BuildFilter = ss.Build.BuildFilter
		service.RootDir = ss.Build.RootDir
		service.AutoDeploy = ss.Build.AutoDeploy
	}

	// Apply Static Site configuration
	if ss.StaticSite != nil {
		service.StaticPublishPath = &ss.StaticSite.StaticPublishPath
		service.Headers = ss.StaticSite.Headers
		service.Routes = ss.StaticSite.Routes
	}

	// Apply Preview configuration
	if ss.Preview != nil {
		service.Previews = ss.Preview.Previews
		service.PreviewPlan = ss.Preview.PreviewPlan
	}

	return service
}

// MarshalYAML implements custom YAML marshaling for StaticSite to match Render schema
func (ss *StaticSite) MarshalYAML() (interface{}, error) {
	// Create a structure that matches the staticService schema exactly
	result := map[string]interface{}{
		"name":    ss.Name,
		"type":    "web",
		"runtime": "static",
	}

	// Add domains
	if len(ss.Domains) > 0 {
		result["domains"] = ss.Domains
	}

	// Add region
	if ss.Region != nil {
		result["region"] = *ss.Region
	}

	// Add Git configuration
	if ss.Git != nil {
		if ss.Git.Repo != nil {
			result["repo"] = *ss.Git.Repo
		}
		if ss.Git.Branch != nil {
			result["branch"] = *ss.Git.Branch
		}
	}

	// Add Build configuration
	if ss.Build != nil {
		if ss.Build.BuildCommand != nil {
			result["buildCommand"] = *ss.Build.BuildCommand
		}
		if ss.Build.PreDeployCommand != nil {
			result["preDeployCommand"] = *ss.Build.PreDeployCommand
		}
		if ss.Build.BuildFilter != nil {
			result["buildFilter"] = ss.Build.BuildFilter
		}
		if ss.Build.RootDir != nil {
			result["rootDir"] = *ss.Build.RootDir
		}
		if ss.Build.AutoDeploy != nil {
			result["autoDeploy"] = *ss.Build.AutoDeploy
		}
	}

	// Add Static Site specific configuration
	if ss.StaticSite != nil {
		if ss.StaticSite.StaticPublishPath != "" {
			result["staticPublishPath"] = ss.StaticSite.StaticPublishPath
		}
		if len(ss.StaticSite.Headers) > 0 {
			result["headers"] = ss.StaticSite.Headers
		}
		if len(ss.StaticSite.Routes) > 0 {
			result["routes"] = ss.StaticSite.Routes
		}
	}

	// Add Preview configuration
	if ss.Preview != nil {
		if ss.Preview.Previews != nil {
			result["previews"] = ss.Preview.Previews
		}
		// Note: previewPlan is not supported for static sites in the schema
	}

	return result, nil
}

// KeyValueService represents a Redis/Key-Value store
type KeyValueService struct {
	// Essential
	Name string `yaml:"name"`

	// Infrastructure
	Plan   *Plan   `yaml:"plan,omitempty"`
	Region *Region `yaml:"region,omitempty"`

	// Configuration groups
	KeyValue *KeyValueConfig `yaml:",inline"`
	Preview  *PreviewConfig  `yaml:",inline,omitempty"`
}

// ToService converts KeyValueService to generic Service
func (kvs *KeyValueService) ToService() *Service {
	service := &Service{
		Name:   kvs.Name,
		Type:   ServiceTypeKeyValue,
		Plan:   kvs.Plan,
		Region: kvs.Region,
	}

	// Apply KeyValue configuration
	if kvs.KeyValue != nil {
		service.IPAllowList = kvs.KeyValue.IPAllowList
		service.MaxMemoryPolicy = kvs.KeyValue.MaxMemoryPolicy
	}

	// Apply Preview configuration
	if kvs.Preview != nil {
		service.Previews = kvs.Preview.Previews
		service.PreviewPlan = kvs.Preview.PreviewPlan
	}

	return service
}

// ServiceBuilder interface for all service types
type ServiceBuilder interface {
	ToService() *Service
}

// Convenience function to build a Blueprint from specific service types
func NewBlueprintFromServices(services []ServiceBuilder, databases []Database, envGroups []EnvVarGroup) *Blueprint {
	genericServices := make([]Service, len(services))
	for i, svc := range services {
		genericServices[i] = *svc.ToService()
	}

	return &Blueprint{
		Services:     genericServices,
		Databases:    databases,
		EnvVarGroups: envGroups,
	}
}

// Builder functions and fluent interface methods

// NewWebService creates a new WebService with the given name and runtime
func NewWebService(name string, runtime Runtime) *WebService {
	return &WebService{
		Name:    name,
		Runtime: runtime,
	}
}

// WithDomains sets the domains for the web service
func (ws *WebService) WithDomains(domains ...string) *WebService {
	ws.Domains = append(ws.Domains, domains...)
	return ws
}

// WithHealthCheck sets the health check path
func (ws *WebService) WithHealthCheck(path string) *WebService {
	ws.HealthCheckPath = &path
	return ws
}

// WithStartCommand sets the start command
func (ws *WebService) WithStartCommand(cmd string) *WebService {
	ws.StartCommand = &cmd
	return ws
}

// WithPlan sets the instance plan
func (ws *WebService) WithPlan(plan Plan) *WebService {
	ws.Plan = &plan
	return ws
}

// WithRegion sets the region
func (ws *WebService) WithRegion(region Region) *WebService {
	ws.Region = &region
	return ws
}

// WithGit configures Git repository settings
func (ws *WebService) WithGit(repo string, branch ...string) *WebService {
	ws.Git = &GitConfig{Repo: &repo}
	if len(branch) > 0 {
		ws.Git.Branch = &branch[0]
	}
	return ws
}

// WithBuild configures build settings
func (ws *WebService) WithBuild(buildCmd string) *WebService {
	if ws.Build == nil {
		ws.Build = &BuildConfig{}
	}
	ws.Build.BuildCommand = &buildCmd
	return ws
}

// WithPreDeploy sets the pre-deploy command
func (ws *WebService) WithPreDeploy(cmd string) *WebService {
	if ws.Build == nil {
		ws.Build = &BuildConfig{}
	}
	ws.Build.PreDeployCommand = &cmd
	return ws
}

// WithAutoDeploy sets auto-deploy flag
func (ws *WebService) WithAutoDeploy(enabled bool) *WebService {
	if ws.Build == nil {
		ws.Build = &BuildConfig{}
	}
	ws.Build.AutoDeploy = &enabled
	return ws
}

// WithDocker configures Docker settings
func (ws *WebService) WithDocker(config *DockerConfig) *WebService {
	ws.Docker = config
	return ws
}

// WithDockerfile sets Docker build configuration
func (ws *WebService) WithDockerfile(dockerfilePath string, context ...string) *WebService {
	if ws.Docker == nil {
		ws.Docker = &DockerConfig{}
	}
	ws.Docker.DockerfilePath = &dockerfilePath
	if len(context) > 0 {
		ws.Docker.DockerContext = &context[0]
	}
	return ws
}

// WithDockerImage sets prebuilt Docker image configuration
func (ws *WebService) WithDockerImage(imageURL string) *WebService {
	if ws.Docker == nil {
		ws.Docker = &DockerConfig{}
	}
	ws.Docker.Image = &DockerImage{URL: imageURL}
	return ws
}

// WithScaling configures manual scaling
func (ws *WebService) WithScaling(numInstances int) *WebService {
	if ws.Scaling == nil {
		ws.Scaling = &ScalingConfig{}
	}
	ws.Scaling.NumInstances = &numInstances
	return ws
}

// WithAutoScaling configures auto scaling
func (ws *WebService) WithAutoScaling(min, max int, targetCPU ...int) *WebService {
	if ws.Scaling == nil {
		ws.Scaling = &ScalingConfig{}
	}
	scaling := &Scaling{
		MinInstances: &min,
		MaxInstances: &max,
	}
	if len(targetCPU) > 0 {
		scaling.TargetCPUPercent = &targetCPU[0]
	}
	ws.Scaling.Scaling = scaling
	return ws
}

// WithEnvVars adds environment variables
func (ws *WebService) WithEnvVars(envVars ...EnvVar) *WebService {
	ws.EnvVars = append(ws.EnvVars, envVars...)
	return ws
}

// WithEnv adds a simple key-value environment variable
func (ws *WebService) WithEnv(key, value string) *WebService {
	ws.EnvVars = append(ws.EnvVars, EnvVar{
		Key:   &key,
		Value: &value,
	})
	return ws
}

// WithDisk configures persistent disk
func (ws *WebService) WithDisk(name, mountPath string, sizeGB ...int) *WebService {
	disk := &Disk{
		Name:      name,
		MountPath: mountPath,
	}
	if len(sizeGB) > 0 {
		disk.SizeGB = &sizeGB[0]
	}
	ws.Disk = disk
	return ws
}

// NewBackgroundWorker creates a new BackgroundWorker
func NewBackgroundWorker(name string, runtime Runtime) *BackgroundWorker {
	return &BackgroundWorker{
		Name:    name,
		Runtime: runtime,
	}
}

// NewPrivateService creates a new PrivateService
func NewPrivateService(name string, runtime Runtime) *PrivateService {
	return &PrivateService{
		Name:    name,
		Runtime: runtime,
	}
}

// WithStartCommand sets the start command for the worker
func (bw *BackgroundWorker) WithStartCommand(cmd string) *BackgroundWorker {
	bw.StartCommand = &cmd
	return bw
}

// WithPlan sets the instance plan for the worker
func (bw *BackgroundWorker) WithPlan(plan Plan) *BackgroundWorker {
	bw.Plan = &plan
	return bw
}

// WithRegion sets the region for the worker
func (bw *BackgroundWorker) WithRegion(region Region) *BackgroundWorker {
	bw.Region = &region
	return bw
}

// WithGit configures Git repository settings for the worker
func (bw *BackgroundWorker) WithGit(repo string, branch ...string) *BackgroundWorker {
	bw.Git = &GitConfig{Repo: &repo}
	if len(branch) > 0 {
		bw.Git.Branch = &branch[0]
	}
	return bw
}

// WithBuild configures build settings for the worker
func (bw *BackgroundWorker) WithBuild(buildCmd string) *BackgroundWorker {
	if bw.Build == nil {
		bw.Build = &BuildConfig{}
	}
	bw.Build.BuildCommand = &buildCmd
	return bw
}

// WithEnvVars adds environment variables to the worker
func (bw *BackgroundWorker) WithEnvVars(envVars ...EnvVar) *BackgroundWorker {
	bw.EnvVars = append(bw.EnvVars, envVars...)
	return bw
}

// WithEnv adds a simple key-value environment variable to the worker
func (bw *BackgroundWorker) WithEnv(key, value string) *BackgroundWorker {
	bw.EnvVars = append(bw.EnvVars, EnvVar{
		Key:   &key,
		Value: &value,
	})
	return bw
}

// WithStartCommand sets the start command for the private service
func (ps *PrivateService) WithStartCommand(cmd string) *PrivateService {
	ps.StartCommand = &cmd
	return ps
}

// WithPlan sets the instance plan for the private service
func (ps *PrivateService) WithPlan(plan Plan) *PrivateService {
	ps.Plan = &plan
	return ps
}

// WithRegion sets the region for the private service
func (ps *PrivateService) WithRegion(region Region) *PrivateService {
	ps.Region = &region
	return ps
}

// WithGit configures Git repository settings for the private service
func (ps *PrivateService) WithGit(repo string, branch ...string) *PrivateService {
	ps.Git = &GitConfig{Repo: &repo}
	if len(branch) > 0 {
		ps.Git.Branch = &branch[0]
	}
	return ps
}

// WithBuild configures build settings for the private service
func (ps *PrivateService) WithBuild(buildCmd string) *PrivateService {
	if ps.Build == nil {
		ps.Build = &BuildConfig{}
	}
	ps.Build.BuildCommand = &buildCmd
	return ps
}

// WithEnvVars adds environment variables to the private service
func (ps *PrivateService) WithEnvVars(envVars ...EnvVar) *PrivateService {
	ps.EnvVars = append(ps.EnvVars, envVars...)
	return ps
}

// WithEnv adds a simple key-value environment variable to the private service
func (ps *PrivateService) WithEnv(key, value string) *PrivateService {
	ps.EnvVars = append(ps.EnvVars, EnvVar{
		Key:   &key,
		Value: &value,
	})
	return ps
}

// NewCronJob creates a new CronJob
func NewCronJob(name string, runtime Runtime, schedule string) *CronJob {
	return &CronJob{
		Name:     name,
		Runtime:  runtime,
		Schedule: schedule,
	}
}

// WithStartCommand sets the start command for the cron job
func (cj *CronJob) WithStartCommand(cmd string) *CronJob {
	cj.StartCommand = &cmd
	return cj
}

// WithRegion sets the region for the cron job
func (cj *CronJob) WithRegion(region Region) *CronJob {
	cj.Region = &region
	return cj
}

// WithGit configures Git repository settings for the cron job
func (cj *CronJob) WithGit(repo string, branch ...string) *CronJob {
	cj.Git = &GitConfig{Repo: &repo}
	if len(branch) > 0 {
		cj.Git.Branch = &branch[0]
	}
	return cj
}

// WithBuild configures build settings for the cron job
func (cj *CronJob) WithBuild(buildCmd string) *CronJob {
	if cj.Build == nil {
		cj.Build = &BuildConfig{}
	}
	cj.Build.BuildCommand = &buildCmd
	return cj
}

// WithEnvVars adds environment variables to the cron job
func (cj *CronJob) WithEnvVars(envVars ...EnvVar) *CronJob {
	cj.EnvVars = append(cj.EnvVars, envVars...)
	return cj
}

// NewStaticSite creates a new StaticSite
func NewStaticSite(name string) *StaticSite {
	return &StaticSite{
		Name: name,
	}
}

// WithPublishPath sets the static publish path
func (ss *StaticSite) WithPublishPath(path string) *StaticSite {
	if ss.StaticSite == nil {
		ss.StaticSite = &StaticSiteConfig{}
	}
	ss.StaticSite.StaticPublishPath = path
	return ss
}

// WithDomains sets the domains for the static site
func (ss *StaticSite) WithDomains(domains ...string) *StaticSite {
	ss.Domains = append(ss.Domains, domains...)
	return ss
}

// WithHeaders adds HTTP headers
func (ss *StaticSite) WithHeaders(headers ...Header) *StaticSite {
	if ss.StaticSite == nil {
		ss.StaticSite = &StaticSiteConfig{}
	}
	ss.StaticSite.Headers = append(ss.StaticSite.Headers, headers...)
	return ss
}

// WithRoutes adds routing rules
func (ss *StaticSite) WithRoutes(routes ...Route) *StaticSite {
	if ss.StaticSite == nil {
		ss.StaticSite = &StaticSiteConfig{}
	}
	ss.StaticSite.Routes = append(ss.StaticSite.Routes, routes...)
	return ss
}

// WithRegion sets the region for the static site
func (ss *StaticSite) WithRegion(region Region) *StaticSite {
	ss.Region = &region
	return ss
}

// WithGit configures Git repository settings for the static site
func (ss *StaticSite) WithGit(repo string, branch ...string) *StaticSite {
	ss.Git = &GitConfig{Repo: &repo}
	if len(branch) > 0 {
		ss.Git.Branch = &branch[0]
	}
	return ss
}

// WithBuild configures build settings for the static site
func (ss *StaticSite) WithBuild(buildCmd string) *StaticSite {
	if ss.Build == nil {
		ss.Build = &BuildConfig{}
	}
	ss.Build.BuildCommand = &buildCmd
	return ss
}

// NewKeyValueService creates a new KeyValueService
func NewKeyValueService(name string) *KeyValueService {
	return &KeyValueService{
		Name: name,
	}
}

// WithPlan sets the plan for the key-value service
func (kvs *KeyValueService) WithPlan(plan Plan) *KeyValueService {
	kvs.Plan = &plan
	return kvs
}

// WithRegion sets the region for the key-value service
func (kvs *KeyValueService) WithRegion(region Region) *KeyValueService {
	kvs.Region = &region
	return kvs
}

// WithIPAllowList sets the IP allow list
func (kvs *KeyValueService) WithIPAllowList(allowList ...IPAllow) *KeyValueService {
	if kvs.KeyValue == nil {
		kvs.KeyValue = &KeyValueConfig{}
	}
	kvs.KeyValue.IPAllowList = append(kvs.KeyValue.IPAllowList, allowList...)
	return kvs
}

// WithPublicAccess allows access from anywhere
func (kvs *KeyValueService) WithPublicAccess() *KeyValueService {
	return kvs.WithIPAllowList(IPAllow{
		Source:      "0.0.0.0/0",
		Description: func() *string { s := "public access"; return &s }(),
	})
}

// WithMaxMemoryPolicy sets the eviction policy
func (kvs *KeyValueService) WithMaxMemoryPolicy(policy MaxMemoryPolicy) *KeyValueService {
	if kvs.KeyValue == nil {
		kvs.KeyValue = &KeyValueConfig{}
	}
	kvs.KeyValue.MaxMemoryPolicy = &policy
	return kvs
}

// Helper functions for creating environment variables

// Env creates a simple environment variable
func Env(key, value string) EnvVar {
	return EnvVar{Key: &key, Value: &value}
}

// EnvFromDatabase creates an environment variable from a database property
func EnvFromDatabase(key, dbName string, property DatabaseProperty) EnvVar {
	return EnvVar{
		Key: &key,
		FromDatabase: &FromDatabase{
			Name:     dbName,
			Property: property,
		},
	}
}

// EnvFromService creates an environment variable from a service property
func EnvFromService(key, serviceName string, serviceType ServiceType, property ServiceProperty) EnvVar {
	return EnvVar{
		Key: &key,
		FromService: &FromService{
			Name:     serviceName,
			Type:     serviceType,
			Property: &property,
		},
	}
}

// EnvSecret creates a secret environment variable that prompts for input
func EnvSecret(key string) EnvVar {
	sync := false
	return EnvVar{
		Key:  &key,
		Sync: &sync,
	}
}

// EnvGenerated creates an auto-generated environment variable
func EnvGenerated(key string) EnvVar {
	generate := true
	return EnvVar{
		Key:           &key,
		GenerateValue: &generate,
	}
}