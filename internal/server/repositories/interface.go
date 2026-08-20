package repositories

type Repository interface {
	// Source возвращает инкапсулированное хранилище (источник).
	Source() any
}
