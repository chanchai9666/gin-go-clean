package usecase

import (
	"errors"

	"arczed/internal/entities/models"
	"arczed/internal/entities/schemas"
	"arczed/internal/server/repositories"
)

type UsersService interface {
	FindUsers(req *schemas.FindUsersReq) ([]models.Users, error)
	FindUsersAll() ([]models.Users, error)
	CreateUsers(req *schemas.AddUsers) error
	FindUsersByUserId(req *schemas.FindUsersByUserIdReq) (*models.Users, error)
	UpdateUsers(req *schemas.AddUsers) error
	DeleteUsers(req *schemas.AddUsers) error
}

type userRequest struct {
	repo repositories.UsersRepository
}

func NewUserService(repo repositories.UsersRepository) UsersService {
	return &userRequest{
		repo: repo,
	}
}

func (s *userRequest) FindUsers(req *schemas.FindUsersReq) ([]models.Users, error) {
	data, err := s.repo.FindUsers(req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *userRequest) FindUsersByUserId(req *schemas.FindUsersByUserIdReq) (*models.Users, error) {
	data, err := s.FindUsers(&schemas.FindUsersReq{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("user not found")
	}
	return &data[0], nil
}

func (s *userRequest) FindUsersAll() ([]models.Users, error) {
	data, err := s.repo.FindUsers(&schemas.FindUsersReq{})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *userRequest) CreateUsers(req *schemas.AddUsers) error {
	return s.repo.CreateUsers(req)
}

func (s *userRequest) UpdateUsers(req *schemas.AddUsers) error {
	return s.repo.UpdateUser(req)
}

func (s *userRequest) DeleteUsers(req *schemas.AddUsers) error {
	return s.repo.DeletedUser(&req.UserId)
}

func (s *userRequest) Login(req *schemas.LoginReq) error {
	result, err := s.FindUsersByUserId(&schemas.FindUsersByUserIdReq{UserId: req.UserId})
	if err != nil {
		return err
	}
	if result.UserId != "" && result.Password != "" {

	}
	return err
}
