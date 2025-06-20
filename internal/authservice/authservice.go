package authservice

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vizurth/eventify/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthService struct {
	db     *pgxpool.Pool
	router *gin.Engine
	secret []byte
}

// NewAuthService создаем сервис аунтификации
func NewAuthService(db *pgxpool.Pool, router *gin.Engine, secretKey []byte) *AuthService {
	return &AuthService{
		db:     db,
		router: router,
		secret: secretKey,
	}
}

// GenerateToken генерируем JWT из username, email, password, role
func (a *AuthService) GenerateToken(username, email, password, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"password": password,
		"email":    email,
		"role":     role,
	})

	return token.SignedString(a.secret)
}

// HashPassword хэширует пароль для дальнейшего использования
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// RegisterHandler обработчик запроса /auth/register
func (a *AuthService) RegisterHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var registerStruct models.RegisterStruct
	// проверяем корректность запроса JSON
	if err := c.ShouldBindJSON(&registerStruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// делаем запрос к базе данных и проверяем существует ли user которого пытаемся зарегистрировать
	var userExists int
	err := a.db.QueryRow(ctx, "SELECT COUNT(*) FROM schema_name.users WHERE username = $1 OR email = $2", registerStruct.Username, registerStruct.Email).Scan(&userExists)
	// ошибка запроса
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// ошибка существования юзера
	if userExists > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := HashPassword(registerStruct.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// добавляем в базу данных информацию о нашем пользователе с хэшированным паролем
	_, err = a.db.Exec(ctx, "INSERT INTO schema_name.users (username, email, password_hash, role) VALUES ($1, $2, $3, $4)", registerStruct.Username, registerStruct.Email, hashedPassword, registerStruct.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// выводим команду о успешном создании user
	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (a *AuthService) LoginHandler(c *gin.Context) {
	//получаем контекст для работы с базой данных
	ctx := c.Request.Context()

	// проверяем коректность JSON запроса
	var loginStruct models.LoginStruct
	if err := c.ShouldBindJSON(&loginStruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверяем таблицу на наличие user's
	var hashedPassword string
	if err := a.db.QueryRow(ctx, "SELECT password_hash FROM schema_name.users WHERE username = $1 OR email = $2", loginStruct.Username, loginStruct.Email).Scan(&hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// смотрим корректность пароля
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginStruct.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	//генерируем JWT для дальнейшей работы
	token, err := a.GenerateToken(loginStruct.Username, loginStruct.Email, loginStruct.Password, loginStruct.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// AuthMiddleWare мидлварь которая проверяет имеет ли user role admin
func (a *AuthService) AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// смотрим авторизованны ли мы
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		// получаем token через claims
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return a.secret, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
		}
		// просматриваем role user чтобы узнать давать ему доступ к действию или нет
		role := claims["role"].(string)
		if role != "admin" {
			c.JSON(http.StatusLocked, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RegisterRoutes собираем все хендлеры в одну функцию
func (a *AuthService) RegisterRoutes() {
	auth := a.router.Group("/auth")
	auth.POST("/login", a.LoginHandler)
	auth.POST("/register", a.RegisterHandler)
}
