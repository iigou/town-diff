package internal

import (
	"github.com/iigou/town-diff/pkg/api"
	"gorm.io/gorm"
)

/*
TownRepo is the repository handling the town struct
*/
type TownRepo struct {
	DbConnFn func() (*gorm.DB, error)
}

/*
Get the towns matching the provided criteria
*/
func (t *TownRepo) Get(criteria *api.Town) ([]api.Town, error) {
	db, err := t.DbConnFn()
	if err != nil {
		return nil, err
	}

	results := []api.Town{}
	fields := []string{}
	if len(criteria.Name) > 0 {
		fields = append(fields, "name")

	}

	if criteria.Id > 0 {
		fields = append(fields, "id")
	}

	if len(fields) > 0 {
		db.Where(criteria, fields).Find(&results)
	} else {
		db.Find(&results)
	}

	return results, nil
}

/*
Save the provided town
*/
func (t *TownRepo) Save(town *api.Town) (*api.Town, error) {

	db, err := t.DbConnFn()
	if err != nil {
		return nil, err
	}
	result := db.Create(town)
	if result.Error != nil {
		return nil, result.Error
	}
	return town, nil
}

/*
Update the town matching the townName, with the given town struct
*/
func (t *TownRepo) Update(toUpdate *api.Town) (*api.Town, error) {
	db, err := t.DbConnFn()
	if err != nil {
		return nil, err
	}

	db.Model(toUpdate).Updates(toUpdate)

	return toUpdate, nil
}

/*
Delete the town matching the townName
*/
func (t *TownRepo) Delete(town *api.Town) (bool, error) {

	db, err := t.DbConnFn()
	if err != nil {
		return false, err
	}

	result := db.Delete(town)
	return result.RowsAffected == 1, result.Error
}
