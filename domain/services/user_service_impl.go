package services

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/chud-lori/go-boilerplate/adapters/transport"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
)

type UserServiceImpl struct {
	ports.UserRepository
}

// provider or constructor
func NewUserService(userRepository ports.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}

func generatePasscode() string {
	// get current ms
	curMs := time.Now().Nanosecond() / 1000

	// convert ms to str and get first 4 char
	msStr := strconv.Itoa(curMs)[:4]

	// generate random char between A and Z
	var alphb []int
	for i := 0; i < 4; i++ {
		alphb = append(alphb, rand.Intn(26)+65)
	}

	// Convert ascii values to character and join them
	var alphChar []string
	for _, a := range alphb {
		alphChar = append(alphChar, string(rune(a)))
	}
	alphStr := strings.Join(alphChar, "")

	// combine alphabet string and ms string
	return alphStr + msStr
}

func (service *UserServiceImpl) Save(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error) {
	user := entities.User{
		Id:         "",
		Email:      request.Email,
		Passcode:   generatePasscode(),
		Created_at: time.Now(),
	}
	user_result, error := service.UserRepository.Save(ctx, &user)

	if error != nil {

		panic(error)
	}

	user_response := &transport.UserResponse{
		Id:         user_result.Id,
		Email:      user_result.Email,
		Created_at: user_result.Created_at,
	}

	return user_response, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error) {
	user := entities.User{
		Id:         "",
		Email:      request.Email,
		Created_at: time.Now(),
	}

	user_result, error := service.UserRepository.Update(ctx, &user)
	if error != nil {
		panic(error)
	}

	user_response := &transport.UserResponse{
		Id:         user_result.Id,
		Email:      user_result.Email,
		Created_at: user_result.Created_at,
	}

	return user_response, nil
}

func (service *UserServiceImpl) Delete(ctx context.Context, id string) error {

	err := service.UserRepository.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) FindById(ctx context.Context, id string) (*transport.UserResponse, error) {
	user_result, err := s.UserRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	user_response := &transport.UserResponse{
		Id:         user_result.Id,
		Email:      user_result.Email,
		Created_at: user_result.Created_at,
	}

	return user_response, nil
}

func (service *UserServiceImpl) FindAll(ctx context.Context) ([]*transport.UserResponse, error) {

	users_result, err := service.UserRepository.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	users_response := make([]*transport.UserResponse, len(users_result))

	for i, user := range users_result {
		users_response[i] = &transport.UserResponse{
			Id:         user.Id,
			Email:      user.Email,
			Created_at: user.Created_at,
		}
	}

	return users_response, nil
}
