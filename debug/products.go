package debug

import (
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
)

func InitFakeProductsIfLessThenN(n int) error {
	products_count, err := storage.StorageInstance.GetProductsCount()
	if err != nil {
		return err
	}

	if products_count > n {
		return nil
	}

	new_products := n - products_count
	for i := 0; i < new_products; i++ {

		go func() {
			r := GenerateRandomString()
			category := models.NewProduct(r, r, r, 0, models.StockTypeInStock, r, 0, []string{r, r, r})
			storage.StorageInstance.CreateProduct(category)
		}()
	}

	return nil
}
