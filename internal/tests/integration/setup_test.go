package integration_test

import (
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/handler"
	"fdlp-standard-api/internal/repositories"
	"fdlp-standard-api/internal/routes"
	"fdlp-standard-api/internal/services"
	"fdlp-standard-api/internal/types"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ValidationUtil struct {
	validator *validator.Validate
}

func (cv *ValidationUtil) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Shadow structs for SQLite migration (removed postgres-specific tags)
type TestRole struct {
	ID          types.UUID `gorm:"primaryKey;type:text"`
	Name        string     `gorm:"uniqueIndex;not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TestRole) TableName() string { return "roles" }

type TestUser struct {
	ID        types.UUID `gorm:"primaryKey;type:text"`
	Username  string     `gorm:"uniqueIndex;not null"`
	Email     string     `gorm:"uniqueIndex;not null"`
	Fisrtname string
	Password  string     `gorm:"not null"`
	RoleID    types.UUID `gorm:"not null"`
	Role      TestRole   `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TestUser) TableName() string { return "users" }

// SetupIntegrationServer initializes an in-memory test environment.
// It returns the Echo instance and the GORM DB connection for seeding data.
func SetupIntegrationServer() (*echo.Echo, *gorm.DB) {
	// 1. Setup DB (SQLite In-Memory with shared cache)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 2. Migrate using Shadow structs
	err = db.AutoMigrate(&TestRole{}, &TestUser{})
	if err != nil {
		panic("migration failed: " + err.Error())
	}

	// 3. Initialize layers
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	wsService := services.NewWebSocketService()
	userService := services.NewUserService(userRepo, roleRepo, "test_secret", db)
	roleService := services.NewRoleService(roleRepo)

	userHandler := handler.NewUserHandler(userService, wsService)
	roleHandler := handler.NewRoleHandler(roleService)

	// 4. Setup Echo
	e := echo.New()
	e.Validator = &ValidationUtil{validator: validator.New()}

	// 5. Protected Routes Group
	r := e.Group("/v1")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(dto.JwtCustomClaims)
		},
		SigningKey: []byte("test_secret"),
	}
	r.Use(echojwt.WithConfig(config))

	// 6. Register Routes
	routes.RegisterUserRoutes(e, r, userHandler)
	routes.RegisterRoleRoutes(e, roleHandler)

	return e, db
}
