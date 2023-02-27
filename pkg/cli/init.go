// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

type Init interface {
	Init() error
}

func InitAll(objs ...Init) error {
	var err error
	for _, o := range objs {
		if err = o.Init(); err != nil {
			return err
		}
	}
	return nil
}

func initOpts(c *Command) error {
	i := c.Init()
	if i == nil {
		return nil
	}
	return i.Init()
}
