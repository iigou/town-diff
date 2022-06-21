package internal

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/iigou/town-diff/pkg/api"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	svc  api.ITownSvc
	repo api.ITownRepo
	data *api.Town
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.0.0"))

	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:        "sqlmock_db_0",
		DriverName: "mysql",
		Conn:       db,
	}), &gorm.Config{})

	require.NoError(s.T(), err)
	s.repo = &TownRepo{DbConnFn: func() (*gorm.DB, error) { return s.DB, nil }}
	s.svc = &TownSvc{TownRepo: s.repo}
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_Get_by_town_name() {
	var (
		id   = uint(1)
		name = "townName"
	)

	s.mock.MatchExpectationsInOrder(false)
	// s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `towns` WHERE `towns`.`name` = ?")).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))
	// s.mock.ExpectCommit()

	res, err := s.svc.Get(&api.Town{Name: name})

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal([]api.Town{{Id: id, Name: name}}, res))
}

func (s *Suite) Test_Get_by_town_id() {
	var (
		id   = uint(1)
		name = "townName"
	)

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `towns` WHERE `towns`.`id` = ?")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))
	res, err := s.svc.Get(&api.Town{Id: id})

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal([]api.Town{{Id: id, Name: name}}, res))
}

func (s *Suite) Test_Get_by_town_id_name() {
	var (
		id   = uint(1)
		name = "townName"
	)

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `towns` WHERE `towns`.`id` = ? AND `towns`.`name` = ?")).
		WithArgs(id, name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))
	res, err := s.svc.Get(&api.Town{Id: id, Name: name})

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal([]api.Town{{Id: id, Name: name}}, res))
}

func (s *Suite) Test_Get_towns() {
	var (
		townId1   = uint(1)
		townName1 = "townName"
		townId2   = uint(2)
		townName2 = "another townName"
	)

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `towns`")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(townId1, townName1).AddRow(townId2, townName2))

	res, err := s.svc.Get(&api.Town{})

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal([]api.Town{{Id: townId1, Name: townName1}, {Id: townId2, Name: townName2}}, res))
}

func (s *Suite) Test_Distance_towns() {
	var (
		homeId   = uint(1)
		homeName = "home townName"
		homeLat  = 43.580719
		homeLon  = 7.12087

		destId   = uint(2)
		destName = "dest townName"
		destLat  = 48.856613
		destLon  = 2.352222
	)

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `towns` WHERE `towns`.`name` = ?")).
		WithArgs(homeName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lat", "lon"}).
			AddRow(homeId, homeName, homeLat, homeLon))

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `towns` WHERE `towns`.`name` = ?")).
		WithArgs(destName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lat", "lon"}).
			AddRow(destId, destName, destLat, destLon))

	res, err := s.svc.Diff(homeName, destName)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&api.DiffResult{Distance: 691.5725305976358, Units: "kilometers"}, res))
}
