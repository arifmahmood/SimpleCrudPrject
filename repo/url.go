package repo

import "simple-crud-project/model"

// Url represents Url repository interface
type Url interface {
	EnsureIndices(url *model.Url) error
	Fetch(username string) (*model.Url, error)
	Create(user *model.Url) error
}
