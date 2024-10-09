package repository

type Repository interface {
	Find(key string) (string, bool)
	Save(key, value string) error
}
