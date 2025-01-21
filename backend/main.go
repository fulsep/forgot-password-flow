package main

import (
	"ff-flow/lib"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/matthewhartstonge/argon2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FormForgotPassword struct {
	Email string `form:"email" binding:"required"`
}

type FormResetPassword struct {
	Token           string `form:"token" binding:"required"`
	NewPassword     string `form:"password"`
	ConfirmPassword string `form:"confirm-password" binding:"eqfield=NewPassword"`
}

type FormRegister struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required,min=8"`
}

type FormLogin struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Results any    `json:"results,omitempty"`
}

type User struct {
	gorm.Model
	ID        int    `gorm:"primaryKey,autoIncrement"`
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type ForgotPassword struct {
	gorm.Model
	ID        int `gorm:"primaryKey,autoIncrement"`
	Token     string
	UserID    int
	User      User
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func main() {
	r := gin.Default()

	db, err := gorm.Open(sqlite.Open("backend.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&ForgotPassword{})

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	r.POST("/forgot-password", func(ctx *gin.Context) {
		var form FormForgotPassword
		err := ctx.ShouldBind(&form)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}

		var findRequest ForgotPassword
		tx := db.Joins("User").First(&findRequest, "User__email = ?", form.Email)

		if tx.Error != nil {
			if !strings.Contains(tx.Error.Error(), "record not found") {
				ctx.JSON(http.StatusBadRequest, Response{
					Success: false,
					Message: fmt.Sprintf("Error: %s", tx.Error.Error()),
				})
				return
			}
		}

		if findRequest != (ForgotPassword{}) {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "We have sent an email, please check it first",
			})
			return
		}

		// db.Create()

		var user User

		tx = db.First(&user, "email = ?", form.Email)

		if tx.Error != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", tx.Error.Error()),
			})
			return
		}

		rand := rand.Intn(999999)
		tx = db.Create(&ForgotPassword{
			Token:  strconv.Itoa(rand),
			UserID: user.ID,
		})

		if tx.Error != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", tx.Error.Error()),
			})
			return
		}

		lib.SendMail(form.Email, "Forgot Password", lib.ResetPasswordTemplate(strconv.Itoa(rand)))
		ctx.JSON(http.StatusOK, Response{
			Success: true,
			Message: "Your request has been sent to email, please check it",
		})
	})

	r.POST("/reset-password", func(ctx *gin.Context) {
		var form FormResetPassword
		err := ctx.ShouldBind(&form)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		var findRequest ForgotPassword
		tx := db.Joins("User").First(&findRequest, "token = ?", form.Token)

		if tx.Error != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", tx.Error.Error()),
			})
			return
		}
		cfg := argon2.DefaultConfig()

		raw, err := cfg.Hash([]byte(form.NewPassword), nil)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		newPassword := raw.Encode()

		updatedUser := User{
			ID: findRequest.User.ID,
		}
		db.Model(&updatedUser).Update("password", newPassword)
		db.Delete(&findRequest, "id = ?", findRequest.ID)

		ctx.JSON(http.StatusOK, Response{
			Success: true,
			Message: "Your password has been reset",
		})
	})

	r.POST("/login", func(ctx *gin.Context) {
		var form FormLogin
		err := ctx.ShouldBind(&form)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		var user User
		tx := db.First(&user, "email = ?", form.Email)

		if tx.Error != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", tx.Error.Error()),
			})
			return
		}

		ok, _ := argon2.VerifyEncoded([]byte(form.Password), []byte(user.Password))

		if !ok {
			ctx.JSON(http.StatusUnauthorized, Response{
				Success: false,
				Message: "Wrong email or password",
			})
			return
		}

		ctx.JSON(http.StatusOK, Response{
			Success: true,
			Message: "Login success",
		})
	})

	r.POST("/register", func(ctx *gin.Context) {
		var form FormRegister
		err := ctx.ShouldBind(&form)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		cfg := argon2.DefaultConfig()

		raw, err := cfg.Hash([]byte(form.Password), nil)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		password := raw.Encode()

		tx := db.Create(&User{
			Email:    form.Email,
			Password: string(password),
		})
		if tx.Error != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: fmt.Sprintf("Error: %s", tx.Error.Error()),
			})
			return
		}

		ctx.JSON(http.StatusCreated, Response{
			Success: true,
			Message: "Register success",
		})
	})

	r.Run(":8888")
}
