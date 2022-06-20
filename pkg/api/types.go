package api

type Town struct {
	Id   uint    `json:"id" schema: "id" gorm:"primaryKey,autoIncrement"`
	Name string  `json:"name" schema:"name" gorm:"unique"`
	Lon  float64 `json:"lon" schema:"lon"`
	Lat  float64 `json:"lat" schema:"lat"`
}

type ITownSvc interface {
	Get(criteria *Town) ([]Town, error)
	Save(t *Town) (*Town, error)
	Update(townName string, t *Town) (*Town, error)
	Delete(townName string) (bool, error)
}

type ITownDiffSvc interface {
	Diff(home *Town, destination *Town) (float64, error)
}
