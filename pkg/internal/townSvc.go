package internal

import (
	"github.com/iigou/town-diff/pkg/api"
)

type TownSvc struct {
}

func (t *TownSvc) Get(criteria *api.Town) ([]api.Town, error) {
	criteria.Name = "changed from svc!"
	return []api.Town{*criteria}, nil
}

func (t *TownSvc) Save(town *api.Town) (api.Town, error) {

	return *town, nil
}

func (t *TownSvc) Update(townID string, town *api.Town) (api.Town, error) {

	return *town, nil
}

func (t *TownSvc) Delete(townID string) (bool, error) {

	return false, nil
}
