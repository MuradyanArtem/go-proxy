package repository

type Scanner interface {
	MakeInject(req *string) ([]string, error)
}
