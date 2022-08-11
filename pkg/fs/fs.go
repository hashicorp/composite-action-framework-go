package fs

type FS struct {
	Settings
}

func New(opts ...Option) *FS {
	return &FS{
		Settings: newSettings(opts),
	}
}
