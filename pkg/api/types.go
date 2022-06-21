package api

type Town struct {
	Id   uint    `json:"id" schema: "id" gorm:"primaryKey,autoIncrement"`
	Name string  `json:"name" schema:"name" gorm:"unique"`
	Lon  float64 `json:"lon" schema:"lon"`
	Lat  float64 `json:"lat" schema:"lat"`
}

type DiffResult struct {
	Distance float64 `json:"distance"`
	Units    string  `json:"units"`
}

type ITownSvc interface {
	Get(criteria *Town) ([]Town, error)
	Save(t *Town) (*Town, error)
	Update(townName string, t *Town) (*Town, error)
	Delete(townName string) (bool, error)
	Diff(home string, destination string) (*DiffResult, error)
}

type ITownRepo interface {
	Get(criteria *Town) ([]Town, error)
	Save(t *Town) (*Town, error)
	Update(t *Town) (*Town, error)
	Delete(t *Town) (bool, error)
}
