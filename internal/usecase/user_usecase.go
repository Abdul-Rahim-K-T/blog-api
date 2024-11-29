package usecase

import (
	"blog-api/internal/entity"
	"blog-api/internal/repository/mysql"
	"blog-api/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(user *entity.User) error
	Login(username, password string) (string, error)
	GetByUsernameOrEmail(username, email string) (*entity.User, error)
	GenerateJWTToken(userID int, role string) (string, error)
}

type userUsecase struct {
	userRepo  *mysql.UserRepository
	jwtSecret string
}

func NewUserUsecase(userRepo *mysql.UserRepository, jwtSecret string) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) Register(user *entity.User) error {
	return u.userRepo.Create(user)
}

func (u *userUsecase) Login(username, password string) (string, error) {

	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return "", err
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	token, err := jwt.GenerateJWTToken(user.ID, user.Role, u.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUsecase) GetByUsernameOrEmail(username, email string) (*entity.User, error) {
	return u.userRepo.GetByUsernameOrEmail(username, email)
}

func (u *userUsecase) GenerateJWTToken(userID int, role string) (string, error) {
	return jwt.GenerateJWTToken(userID, role, u.jwtSecret)
}
