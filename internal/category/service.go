package category

type CategoryService interface {
	GetCategory() ([]Category, error)
	GetCategoryByName(string) (Category, error)
	GetCategoryByPostID(int64) ([]Category, error)
	CreateCategoryPostByPostID(string, int64) (error)
}

type categoryService struct {
	storage CategoryStorage
}

func NewCategoryService(storage CategoryStorage) CategoryService {
	return &categoryService{
		storage: storage,
	}
}

func (s *categoryService) GetCategory() ([]Category, error) {
	return s.storage.GetCategory()
}

func (s *categoryService) GetCategoryByName(name string) (Category, error) {
	return s.storage.GetCategoryByName(name)
}

func (s *categoryService) GetCategoryByPostID(id int64) (c []Category, err error) {
	return s.storage.GetCategoryByPostID(id)
}


func (s *categoryService) CreateCategoryPostByPostID(name string, id int64) error {
	return s.storage.CreateCategoryPostByPostID(name, id)
}