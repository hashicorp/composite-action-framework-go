package cli

// Env should be implemented by options structs that read
// from the environment.
type Env interface {
	ReadEnv() error
}

func parseEnv(c *Command) error {
	e := c.Env()
	if e == nil {
		return nil
	}
	return e.ReadEnv()
}
