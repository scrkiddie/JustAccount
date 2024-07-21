package service

import (
	"awesomeProject12/internal/adapter"
	"awesomeProject12/internal/entity"
	"awesomeProject12/internal/model"
	"awesomeProject12/internal/repository"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

type UserService struct {
	DB             *gorm.DB
	Validator      *validator.Validate
	FileStorage    *adapter.FileAdapter
	UserRepository *repository.UserRepository
	Config         *viper.Viper
}

func NewUserService(db *gorm.DB, validate *validator.Validate,
	storage *adapter.FileAdapter, userRepository *repository.UserRepository, config *viper.Viper) *UserService {
	return &UserService{db, validate, storage, userRepository, config}
}

func (s *UserService) Create(ctx context.Context, request *model.RegisterUserRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	total, err := s.UserRepository.CountByUsername(tx, request.Username)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	if total > 0 {
		return fiber.NewError(fiber.StatusConflict, "Username already exists")
	}

	total, err = s.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	if total > 0 {
		return fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	userEntity := new(entity.User)
	userEntity.FirstName = request.FirstName
	userEntity.LastName = request.LastName
	userEntity.Username = request.Username
	userEntity.Email = request.Email
	userEntity.Password = string(password)

	if err := s.UserRepository.Create(tx, userEntity); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := s.UserRepository.FindByUsername(tx, user, request.Username); err != nil {
		log.Println(err.Error())
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Username or password is incorrect")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		log.Println(err.Error())
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Username or password is incorrect")
	}

	key := s.Config.GetString("jwt.secret")
	exp := s.Config.GetInt("jwt.exp")
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Duration(exp) * time.Hour).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
	if err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &model.UserResponse{Token: token}, nil
}

func (s *UserService) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	tokenString := request.Token

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fiber.ErrInternalServerError
		}
		return []byte(s.Config.GetString("jwt.secret")), nil
	})

	if err != nil {
		log.Println("Error parsing token:", err.Error())
		return nil, fiber.ErrUnauthorized
	}

	var id int
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if subFloat64, ok := claims["sub"].(float64); ok {
			id = int(subFloat64)
		} else {
			return nil, fiber.ErrUnauthorized
		}
	} else {
		return nil, fiber.ErrUnauthorized
	}

	user := new(entity.User)
	if err := s.UserRepository.FindById(tx, user, id); err != nil {
		log.Println(err.Error())
		return nil, fiber.NewError(fiber.StatusNotFound, "Not found")
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &model.Auth{ID: user.ID}, nil
}

func (s *UserService) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := s.UserRepository.FindById(tx, user, request.ID); err != nil {
		log.Println(err.Error())
		return nil, fiber.NewError(fiber.StatusNotFound, "Not found")
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	userResponse := new(model.UserResponse)
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.Username = user.Username
	userResponse.Email = user.Email
	userResponse.ProfilePicture = user.ProfilePicture

	return userResponse, nil
}

func (s *UserService) Update(ctx context.Context, request *model.UpdateUserRequest, file *model.File) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if err := s.Validator.Struct(file); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	var storedFile *model.File
	var err error

	if file.FileHeader != nil && file.FileHeader.Filename != "" {
		storedFile, err = s.FileStorage.StoreFile(s.Config.GetString("directories.profile_pictures"), file)
		if err != nil {
			log.Println(err.Error())
			return fiber.ErrInternalServerError
		}
	}

	user := new(entity.User)
	if err := s.UserRepository.FindById(tx, user, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "Not found")
	}

	total, err := s.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	if total > 0 && user.Email != request.Email {
		return fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	user.FirstName = request.FirstName
	user.Email = request.Email
	user.LastName = request.LastName

	if storedFile != nil {
		if user.ProfilePicture != "" {
			deletedFile := new(model.File)
			deletedFile.Name = user.ProfilePicture
			s.FileStorage.DeleteFileAsync(s.Config.GetString("directories.profile_pictures"), deletedFile)
		}
		user.ProfilePicture = storedFile.Name
	}

	if err := s.UserRepository.Update(tx, user); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *UserService) UpdatePassword(ctx context.Context, request *model.UpdatePasswordRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := s.UserRepository.FindById(tx, user, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "Not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.CurrentPassword)); err != nil {
		log.Println(err.Error())
		return fiber.NewError(fiber.StatusUnauthorized, "Current password is incorrect")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	user.Password = string(password)

	if err := s.UserRepository.Update(tx, user); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
