package hashing

type Hasher interface {
	Hash(password string) (string, error)
	CompareTo(hash string, password string) error
}
