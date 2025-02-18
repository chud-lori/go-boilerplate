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
	db ports.Database
	ports.UserRepository
	ctxTimeout time.Duration
}

// provider or constructor
func NewUserService(db ports.Database, userRepository ports.UserRepository, ctxTimeout time.Duration) *UserServiceImpl {
	return &UserServiceImpl{
		db: db,
		UserRepository: userRepository,
		ctxTimeout: ctxTimeout,
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

func (s *UserServiceImpl) Save(c context.Context, request *transport.UserRequest) (*transport.UserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.ctxTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx)
    if err != nil {
        return nil, err
    }
	// handle panic gracefully
    defer func() {
        if r := recover(); r != nil || err != nil {
            tx.Rollback()
        }
    }()

	user := entities.User{
		Id:         "",
		Email:      request.Email,
		Passcode:   generatePasscode(),
		Created_at: time.Now(),
	}
	user_result, err := s.UserRepository.Save(ctx, tx, &user)

	if err != nil {
		return nil, err
	}

	user_response := &transport.UserResponse{
		Id:         user_result.Id,
		Email:      user_result.Email,
		Created_at: user_result.Created_at,
	}

	if err := tx.Commit(); err != nil {
        return nil, err
    }

	return user_response, nil
}

func (s *UserServiceImpl) Update(c context.Context, request *transport.UserRequest) (*transport.UserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.ctxTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx)
    if err != nil {
        return nil, err
    }
	// handle panic gracefully
    defer func() {
        if r := recover(); r != nil || err != nil {
            tx.Rollback()
        }
    }()

	user := entities.User{
		Id:         "",
		Email:      request.Email,
		Created_at: time.Now(),
	}

	user_result, err := s.UserRepository.Update(ctx, tx, &user)
	if err != nil {
		return nil, err
	}

	user_response := &transport.UserResponse{
		Id:         user_result.Id,
		Email:      user_result.Email,
		Created_at: user_result.Created_at,
	}

	if err := tx.Commit(); err != nil {
        return nil, err
    }

	return user_response, nil
}

func (s *UserServiceImpl) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, s.ctxTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx)
    if err != nil {
        return err
    }
	// handle panic gracefully
    defer func() {
        if r := recover(); r != nil || err != nil {
            tx.Rollback()
        }
    }()

	err = s.UserRepository.Delete(ctx, tx, id)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
        return err
    }

	return nil
}

func (s *UserServiceImpl) FindById(c context.Context, id string) (*transport.UserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.ctxTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx)
    if err != nil {
        return nil, err
    }
	// handle panic gracefully
    defer func() {
        if r := recover(); r != nil || err != nil {
            tx.Rollback()
        }
    }()

	user_result, err := s.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	user_response := &transport.UserResponse{
		Id:         user_result.Id,
		Email:      user_result.Email,
		Created_at: user_result.Created_at,
	}

	if err := tx.Commit(); err != nil {
        return nil, err
    }

	return user_response, nil
}

func (s *UserServiceImpl) FindAll(c context.Context) ([]*transport.UserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.ctxTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx)
    if err != nil {
        return nil, err
    }
	// handle panic gracefully
    defer func() {
        if r := recover(); r != nil || err != nil {
            tx.Rollback()
        }
    }()

	users_result, err := s.UserRepository.FindAll(ctx, tx)

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

	if err := tx.Commit(); err != nil {
        return nil, err
    }

	return users_response, nil
}
