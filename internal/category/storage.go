package category

type CategoryStorage interface {
	GetCategory() ([]Category, error)
	GetCategoryByName(string) (Category, error)
	GetCategoryByPostID(int64) ([]Category, error)
	CreateCategoryPostByPostID(string, int64) (error)
}
