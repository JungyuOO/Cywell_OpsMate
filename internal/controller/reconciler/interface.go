package reconciler

type Reconciler interface {
	GetNamespace() string
	GetLightspeedServiceImage() string
	GetPostgresImage() string
	GetConsolePluginImage() string
}
