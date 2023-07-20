package goclitools

import (
	"fmt"
)

type Route struct {
	prefix  string
	Flags   []string
	Args    int
	Options []string
	Handler func(*Context)
}

type RouteErr struct {
	message string
}

func (e RouteErr) Error() string {
	return e.message
}

func (r *Route) match(c *Context) error {
	var err error

	err = r.matchFlags(c)

	if err != nil {
		return err
	}

	err = r.matchOptions(c)

	if err != nil {
		return err
	}

	return nil
}

func (r *Route) matchFlags(c *Context) error {
	if len(r.Flags) == 0 && len(c.Flags) > len(r.Flags) {
		return RouteErr{
			message: fmt.Sprintf("command [%s] does not accept flags,", r.prefix),
		}
	}

	exists := func(f string) bool {
		for _, f2 := range r.Flags {
			if f == f2 {
				return true
			}
		}
		return false
	}

	for k := range c.Flags {
		if !exists(k) {
			return RouteErr{
				message: fmt.Sprintf("flag [%s] does not exists", k),
			}
		}

	}

	return nil
}

func (r *Route) matchOptions(c *Context) error {
	if len(r.Options) == 0 && len(c.Options) > len(r.Options) {
		return RouteErr{
			message: fmt.Sprintf("command [%s] does not accept options,", r.prefix),
		}
	}

	exists := func(f string) bool {
		for _, f2 := range r.Options {
			if f == f2 {
				return true
			}
		}
		return false
	}

	for k := range c.Options {
		if !exists(k) {
			return RouteErr{
				message: fmt.Sprintf("option [%s] does not exists", k),
			}
		}
	}

	return nil
}
