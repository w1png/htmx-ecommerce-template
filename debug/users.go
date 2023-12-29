package debug

import (
	"github.com/w1png/htmx-template/models"
	"github.com/w1png/htmx-template/storage"
)

func InitFakeUsersIfLessThenN(n int) error {
	users_count, err := storage.StorageInstance.GetUsersCount()
	if err != nil {
		return err
	}

	if users_count > n {
		return nil
	}

	new_users := n - users_count
	for i := 0; i < new_users; i++ {

		go func() {
			r := GenerateRandomString()
			user, err := models.NewUser(r, r, false)
			if err == nil {
				storage.StorageInstance.CreateUser(user)
			}
		}()
	}

	return nil
}
