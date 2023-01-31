package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"tests_app/internal/config"
	"tests_app/internal/handler"
	"tests_app/internal/repository"
	"tests_app/internal/service"
	"tests_app/pkg/hash"
	"tests_app/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

var conf config.Config

func init() {
	conf = config.Config{
		DB: config.DatabaseConfig{
			Host:     "127.0.0.1",
			Port:     "6001",
			Username: "username",
			DBName:   "test",
			SSLMode:  "disable",
			Password: "qwerty123",
		},
		JWT: config.JWTConfig{
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 720 * time.Hour,
			SigningKey:      "112233",
		},
		PasswordSalt: "123",
	}
}

type APITestSuite struct {
	suite.Suite

	db       *sqlx.DB
	handlers *handler.Handler
	services *service.Service
	repos    *repository.Repository
	router   *gin.Engine

	hasher       hash.PasswordHasher
	tokenManager *token.Manager
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func (s *APITestSuite) SetupTest() {
	// connect to database
	var err error
	s.db, err = repository.NewPostgresDB(repository.Config{
		Host:     conf.DB.Host,
		Port:     conf.DB.Port,
		Username: conf.DB.Username,
		DBName:   conf.DB.DBName,
		SSLMode:  conf.DB.SSLMode,
		Password: conf.DB.Password,
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	s.hasher = hash.NewSHA512Hasher(conf.PasswordSalt)
	s.tokenManager = token.NewManager(conf.JWT.AccessTokenTTL, conf.JWT.SigningKey)

	s.repos = repository.New(s.db)
	s.services = service.New(s.repos, s.tokenManager, s.hasher)
	s.handlers = handler.New(s.services, conf)
	s.router = s.handlers.InitApiRoutes()
}

func (s *APITestSuite) TearDownTest() {
	_, err := s.db.Exec(`TRUNCATE TABLE question_answers, questions, test_answers, tests, refresh_sessions, users RESTART IDENTITY CASCADE;`)
	s.NoError(err)

	s.db.Close()
}

/*func (s *APITestSuite) BeforeTest(_, _ string) {}*/
