package internal

import (
	"math"

	"github.com/iigou/town-diff/pkg/api"
)

/*
TownSvc is the service exposing CRUD functionalities for a Town struct, as well as calculating the distance between two towns
*/
type TownSvc struct {
	TownRepo api.ITownRepo
}

/*
Get the towns matching the provided criteria
*/
func (t *TownSvc) Get(criteria *api.Town) ([]api.Town, error) {
	return t.TownRepo.Get(criteria)
}

/*
Save the provided town
*/
func (t *TownSvc) Save(town *api.Town) (*api.Town, error) {
	return t.TownRepo.Save(town)
}

/*
Update the town matching the townName, with the given town struct
*/
func (t *TownSvc) Update(townName string, town *api.Town) (*api.Town, error) {
	towns, err := t.TownRepo.Get(&api.Town{Name: townName})
	if err != nil {
		return nil, err
	}

	if len(towns) == 0 {
		return town, nil
	}

	toUpdate := towns[0]
	toUpdate.Lat = town.Lat
	toUpdate.Lon = town.Lon
	toUpdate.Name = town.Name

	return t.TownRepo.Update(&toUpdate)
}

/*
Delete the town matching the townName
*/
func (t *TownSvc) Delete(townName string) (bool, error) {
	towns, err := t.TownRepo.Get(&api.Town{Name: townName})
	if err != nil {
		return false, err
	}

	if len(towns) == 0 {
		return true, nil
	}
	return t.TownRepo.Delete(&towns[0])
}

/*
Diff calculates the kilometric distance between the given home and destination
*/
func (t *TownSvc) Diff(home string, destination string) (*api.DiffResult, error) {

	homeTowns, err := t.Get(&api.Town{Name: home})
	if err != nil {
		return nil, err
	}

	if len(homeTowns) == 0 {
		return nil, nil
	}

	destTowns, err := t.Get(&api.Town{Name: destination})
	if err != nil {
		return nil, err
	}

	if len(destTowns) == 0 {
		return nil, nil
	}

	return t.distance(&homeTowns[0], &destTowns[0]), nil
}

func (t *TownSvc) distance(home *api.Town, destination *api.Town) *api.DiffResult {
	radHomeLat := float64(math.Pi * home.Lat / 180)
	radDestLat := float64(math.Pi * destination.Lat / 180)

	theta := float64(home.Lon - destination.Lon)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radHomeLat)*math.Sin(radDestLat) + math.Cos(radHomeLat)*math.Cos(radDestLat)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	return &api.DiffResult{Distance: dist * 1.609344, Units: "kilometers"}
}
