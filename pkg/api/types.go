package api

type Town struct {
	Id   string  `json:"id" schema: "id"`
	Name string  `json:"name" schema:"name"`
	Lon  float64 `json:"lon" schema:"lon"`
	Lat  float64 `json:"lat" schema:"lat"`
}

type ITownSvc interface {
	Get(criteria *Town) ([]Town, error)
	Save(t *Town) (Town, error)
	Update(townID string, t *Town) (Town, error)
	Delete(townID string) (bool, error)
}

type ITownDiffSvc interface {
	Diff(home *Town, destination *Town) (float64, error)
}
