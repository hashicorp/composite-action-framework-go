package cli

type Init interface {
	Init() error
}

func initOpts(c Command) error {
	i := c.Init()
	if i == nil {
		return nil
	}
	return i.Init()
}
