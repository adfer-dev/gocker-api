package storage

type Storage interface {
	Get(int) (interface{}, error)
	Create(interface{}) error
	Update(interface{}) error
	Delete(interface{}) error
}
