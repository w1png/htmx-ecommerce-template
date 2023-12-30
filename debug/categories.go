package debug

import (
	"github.com/w1png/htmx-template/models"
	"github.com/w1png/htmx-template/storage"
)

func InitFakeCategoriesIfLessThenN(n int) error {
	categories_count, err := storage.StorageInstance.GetCategoriesCount()
	if err != nil {
		return err
	}

	if categories_count > n {
		return nil
	}

	new_categories := n - categories_count
	for i := 0; i < new_categories; i++ {

		go func() {
			r := GenerateRandomString()
			category := models.NewCategory(r, r, r, r, 0)
			storage.StorageInstance.CreateCategory(category)
		}()
	}

	return nil
}
