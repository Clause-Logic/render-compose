package render

// Database builder functions

// NewDatabase creates a new Database with the given name
func NewDatabase(name string) *Database {
	return &Database{
		Name: name,
	}
}

// WithPlan sets the database plan
func (db *Database) WithPlan(plan Plan) *Database {
	db.Plan = &plan
	return db
}

// WithPreviewPlan sets the preview environment plan
func (db *Database) WithPreviewPlan(plan Plan) *Database {
	db.PreviewPlan = &plan
	return db
}

// WithRegion sets the database region
func (db *Database) WithRegion(region Region) *Database {
	db.Region = &region
	return db
}

// WithPostgreSQL sets the PostgreSQL version
func (db *Database) WithPostgreSQL(version PostgreSQLVersion) *Database {
	db.PostgresMajorVersion = &version
	return db
}

// WithDatabaseName sets the database name (different from service name)
func (db *Database) WithDatabaseName(name string) *Database {
	db.DatabaseName = &name
	return db
}

// WithUser sets the database user
func (db *Database) WithUser(user string) *Database {
	db.User = &user
	return db
}

// WithDiskSize sets the disk size in GB
func (db *Database) WithDiskSize(sizeGB int) *Database {
	db.DiskSizeGB = &sizeGB
	return db
}

// WithPreviewDiskSize sets the preview environment disk size
func (db *Database) WithPreviewDiskSize(sizeGB int) *Database {
	db.PreviewDiskSizeGB = &sizeGB
	return db
}

// WithIPAllowList adds IP allow list entries
func (db *Database) WithIPAllowList(entries ...IPAllow) *Database {
	db.IPAllowList = append(db.IPAllowList, entries...)
	return db
}

// WithIPAccess adds a single IP allow entry
func (db *Database) WithIPAccess(source string, description ...string) *Database {
	entry := IPAllow{Source: source}
	if len(description) > 0 {
		entry.Description = &description[0]
	}
	return db.WithIPAllowList(entry)
}

// WithPublicAccess allows access from anywhere (0.0.0.0/0)
func (db *Database) WithPublicAccess() *Database {
	return db.WithIPAccess("0.0.0.0/0", "public access")
}

// WithPrivateAccess blocks all external connections (empty IP allow list)
func (db *Database) WithPrivateAccess() *Database {
	db.IPAllowList = []IPAllow{} // Empty list
	return db
}

// WithReadReplicas adds read replicas
func (db *Database) WithReadReplicas(names ...string) *Database {
	for _, name := range names {
		db.ReadReplicas = append(db.ReadReplicas, ReadReplica{Name: name})
	}
	return db
}

// WithHighAvailability enables high availability
func (db *Database) WithHighAvailability() *Database {
	db.HighAvailability = &HighAvailability{Enabled: true}
	return db
}

// EnvVarGroup builder functions

// NewEnvVarGroup creates a new environment variable group
func NewEnvVarGroup(name string) *EnvVarGroup {
	return &EnvVarGroup{
		Name:    name,
		EnvVars: []EnvVar{},
	}
}

// WithEnvVars adds environment variables to the group
func (evg *EnvVarGroup) WithEnvVars(envVars ...EnvVar) *EnvVarGroup {
	evg.EnvVars = append(evg.EnvVars, envVars...)
	return evg
}

// WithEnv adds a simple key-value environment variable
func (evg *EnvVarGroup) WithEnv(key, value string) *EnvVarGroup {
	return evg.WithEnvVars(Env(key, value))
}

// WithSecret adds a secret environment variable (sync: false)
func (evg *EnvVarGroup) WithSecret(key string) *EnvVarGroup {
	return evg.WithEnvVars(EnvSecret(key))
}

// WithGenerated adds an auto-generated environment variable
func (evg *EnvVarGroup) WithGenerated(key string) *EnvVarGroup {
	return evg.WithEnvVars(EnvGenerated(key))
}

// Blueprint builder functions

// NewBlueprint creates a new empty Blueprint
func NewBlueprint() *Blueprint {
	return &Blueprint{
		Services:     []Service{},
		Databases:    []Database{},
		EnvVarGroups: []EnvVarGroup{},
	}
}

// WithServices adds services to the blueprint
func (bp *Blueprint) WithServices(services ...ServiceBuilder) *Blueprint {
	for _, svc := range services {
		bp.Services = append(bp.Services, *svc.ToService())
	}
	return bp
}

// WithDatabases adds databases to the blueprint
func (bp *Blueprint) WithDatabases(databases ...*Database) *Blueprint {
	for _, db := range databases {
		bp.Databases = append(bp.Databases, *db)
	}
	return bp
}

// WithEnvVarGroups adds environment variable groups to the blueprint
func (bp *Blueprint) WithEnvVarGroups(groups ...*EnvVarGroup) *Blueprint {
	for _, group := range groups {
		bp.EnvVarGroups = append(bp.EnvVarGroups, *group)
	}
	return bp
}

// WithPreviews configures preview environments
func (bp *Blueprint) WithPreviews(generation PreviewGeneration, expireAfterDays ...int) *Blueprint {
	bp.Previews = &Previews{Generation: string(generation)}
	if len(expireAfterDays) > 0 {
		bp.PreviewsExpireAfterDays = &expireAfterDays[0]
	}
	return bp
}

// Helper function for environment variable references

// EnvFromGroup creates an environment variable reference to a group
func EnvFromGroup(groupName string) EnvVar {
	return EnvVar{
		FromGroup: &groupName,
	}
}