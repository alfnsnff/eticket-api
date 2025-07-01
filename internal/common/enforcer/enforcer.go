package enforcer

type Enforcer interface {
	Enforce(sub, obj, act string) (bool, error)
}
