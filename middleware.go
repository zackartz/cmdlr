package cmdlr

type Middleware struct {
	Trigger func(ctx Ctx)
}
