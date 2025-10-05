package comp

type Component interface {
	GetName() string
	CanAddItem() bool
	SetProp(string, interface{}) error
	GetProp(string) (interface{}, error)
}
