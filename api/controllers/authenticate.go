package controllers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/Deku/passwords"
	"github.com/rafinhacuri/Deku/sqlite"
	"github.com/rafinhacuri/Deku/utils"
)

type user struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *user) validate() error {
	if u.Email == "" {
		return errors.New("the field 'email' is required")
	}
	if u.Password == "" {
		return errors.New("the field 'password' is required")
	}
	if err := utils.ValidatePassword(u.Password); err != nil {
		return errors.New("the field 'password' must be at least 6 characters long")
	}
	if err := utils.ValidateEmail(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

func (u *user) login(ctx context.Context) (token string, err error) {
	row := sqlite.SQL.QueryRowContext(ctx, "SELECT password FROM users WHERE email = ?", u.Email)
	var hashedPassword string
	if err := row.Scan(&hashedPassword); err != nil {
		return "", err
	}

	if !passwords.VerifyBCrypt(u.Password, hashedPassword) {
		return "", errors.New("invalid credentials")
	}

	token, err = utils.GenerateJWT(u.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

func Auth(c *gin.Context) {
	var u user
	if err := c.ShouldBindJSON(&u); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Invalid request payload"})
		return
	}

	if err := u.validate(); err != nil {
		log.Println("Validation error:", err)
		c.AbortWithStatusJSON(400, gin.H{"message": "Username or password incorrect"})
		return
	}

	ctxReq, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	token, err := u.login(ctxReq)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Username or password incorrect"})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful", "token": token})
}
