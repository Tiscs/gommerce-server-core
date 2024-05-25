package secure

import (
	"context"
	"strings"

	"github.com/expr-lang/expr"
)

var (
	anonymous = &Identity{schema: "", token: (*Token)(nil)}
)

// AuthFunc is a function that authorizes an identity.
type AuthFunc func(user *Identity) error

// Authorize authorizes the identity in the given context with the given auth functions.
func Authorize(ctx context.Context, auth ...AuthFunc) error {
	user := IdentityFromContext(ctx)
	if user == nil {
		user = anonymous
	}
	for _, f := range auth {
		if err := f(user); err != nil {
			return err
		}
	}
	return nil
}

// AuthFuncExpression returns an AuthFunc that evaluates the given expression.
// The expression must return a boolean.
// The expression is evaluated with the identity as the context.
// Example:
//
//	AuthFuncExpression(`Token().Realm() == "default"`)
func AuthFuncExpression(script string) AuthFunc {
	program, err := expr.Compile(script, expr.Env(&Identity{}))
	if err != nil {
		panic(err) // panic here because it's a developer error
	}
	return func(user *Identity) error {
		if output, err := expr.Run(program, user); err != nil {
			return err
		} else if result, ok := output.(bool); !ok {
			return ErrInvalidAuthExprOutput
		} else if !result {
			return ErrPermissionDenied
		}
		return nil
	}
}

// AuthFuncAuthenticated returns an AuthFunc that requires the identity is authenticated.
func AuthFuncAuthenticated(user *Identity) error {
	if user.token == nil || user.token.Subject() == "" {
		return ErrUnauthenticated
	}
	return nil
}

// AuthFuncRequireSchema returns an AuthFunc that requires the identity has the given schema.
func AuthFuncRequireSchema(schema string) AuthFunc {
	return func(user *Identity) error {
		if !strings.EqualFold(user.schema, schema) {
			return ErrPermissionDenied
		}
		return nil
	}
}

// AuthFuncRequireScope returns an AuthFunc that requires the identity has the given scope.
func AuthFuncRequireScope(scope string) AuthFunc {
	return func(user *Identity) error {
		if !user.token.HasScope(scope) {
			return ErrPermissionDenied
		}
		return nil
	}
}

// AuthFuncRequireRealm returns an AuthFunc that requires the identity has the given realm.
func AuthFuncRequireRealm(realm string) AuthFunc {
	return func(user *Identity) error {
		if !strings.EqualFold(user.token.Realm(), realm) {
			return ErrPermissionDenied
		}
		return nil
	}
}
