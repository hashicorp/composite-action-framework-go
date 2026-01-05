// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package cli

// Env should be implemented by options structs that read
// from the environment.
type Env interface {
	ReadEnv() error
}

func ReadEnvAll(objs ...Env) error {
	var err error
	for _, e := range objs {
		if err = e.ReadEnv(); err != nil {
			return err
		}
	}
	return nil
}

func parseEnv(c *Command) error {
	e := c.Env()
	if e == nil {
		return nil
	}
	return e.ReadEnv()
}
