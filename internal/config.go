package internal

import (
	"github.com/leapkit/core/gloves"
)

var (
	// GlovesOptions are the options that will be used by the gloves
	// tool to hot reload the application.
	GlovesOptions = []gloves.Option{
		// Run the tailo watcher so when changes are made to
		// the html code it rebuilds css.
		gloves.WatchExtension(".go"),
	}
)
