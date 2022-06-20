package internal

import (
	"github.com/iigou/town-diff/pkg/api"
	"gorm.io/gorm"
)

type TownSvc struct {
	DbConnFn func() (*gorm.DB, error)
}

func (t *TownSvc) Get(criteria *api.Town) ([]api.Town, error) {
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

func (t *TownSvc) Save(town *api.Town) (*api.Town, error) {

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

func (t *TownSvc) Update(townName string, town *api.Town) (*api.Town, error) {
	db, err := t.DbConnFn()
	if err != nil {
		return nil, err
	}

	towns := []api.Town{}
	db.Where(&api.Town{Name: townName}, "name").Find(&towns)
	if len(towns) == 0 {
		return town, nil
	}

	toUpdate := towns[0]
	toUpdate.Lat = town.Lat
	toUpdate.Lon = town.Lon
	toUpdate.Name = town.Name

	db.Model(&toUpdate).Updates(&toUpdate)

	return &toUpdate, nil
}

func (t *TownSvc) Delete(townName string) (bool, error) {

	db, err := t.DbConnFn()
	if err != nil {
		return false, err
	}

	towns := []api.Town{}
	db.Where(&api.Town{Name: townName}, "name").Find(&towns)

	if len(towns) == 0 {
		return true, nil
	}

	result := db.Delete(towns[0])
	return result.RowsAffected == 1, result.Error
}
