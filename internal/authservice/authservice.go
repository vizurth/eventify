package authservice

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthService struct {
	db     *pgxpool.Pool
	router *gin.Engine
	secret []byte
}

type LoginStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthService(db *pgxpool.Pool, router *gin.Engine, secretKey []byte) *AuthService {
	return &AuthService{
		db:     db,
		router: router,
		secret: secretKey,
	}
}

func (a *AuthService) GenerateToken(username, email, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"password": password,
		"email":    email,
	})

	return token.SignedString(a.secret)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func (a *AuthService) RegisterHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var registerStruct RegisterStruct
	if err := c.ShouldBindJSON(&registerStruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userExists int
	err := a.db.QueryRow(ctx, "SELECT COUNT(*) FROM schema_name.users WHERE username = $1 OR email = $2", registerStruct.Username, registerStruct.Email).Scan(&userExists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if userExists > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := HashPassword(registerStruct.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = a.db.Exec(ctx, "INSERT INTO schema_name.users (username, email, password_hash) VALUES ($1, $2, $3)", registerStruct.Username, registerStruct.Email, hashedPassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (a *AuthService) LoginHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var loginStruct LoginStruct
	if err := c.ShouldBindJSON(&loginStruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var hashedPassword string
	if err := a.db.QueryRow(ctx, "SELECT password_hash FROM schema_name.users WHERE username = $1 OR email = $2", loginStruct.Username, loginStruct.Email).Scan(&hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginStruct.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	token, err := a.GenerateToken(loginStruct.Username, loginStruct.Email, loginStruct.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *AuthService) AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}
		c.Set("username", claims["username"])
		c.Set("email", claims["email"])
		c.Next()
	}
}

func (a *AuthService) RegisterRoutes() {
	a.router.POST("/auth/login", a.LoginHandler)
	a.router.POST("/auth/register", a.RegisterHandler)
}
