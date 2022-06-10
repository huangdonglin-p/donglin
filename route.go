package main

import "donglin.framework.use/framework"

func registerRouter(core *framework.Core) {
	core.Get("foo", FooControllerHandler)
}
