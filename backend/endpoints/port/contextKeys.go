package port

import "context"

type contextKey int

const (
	contextKeyDatastore contextKey = iota
)

// SetDatastore sets the datastore adapter in the context
func SetDatastore(ctx context.Context, db Datastore) context.Context {
	return context.WithValue(ctx, contextKeyDatastore, db)
}

// GetDatastore retrieves a datastore adapter from the given context
func GetDatastore(ctx context.Context) Datastore {
	return ctx.Value(contextKeyDatastore).(Datastore)
}
