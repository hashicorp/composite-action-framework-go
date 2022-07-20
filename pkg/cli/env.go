package cli

// Env should be implemented by options structs that read
// from the environment.
type Env interface {
	ReadEnv() error
}
