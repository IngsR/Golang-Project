package services

import (
	"errors"
	"goproject/models"
	"goproject/repositories"
)

type AuthService interface {
	Login(email, password string) (*models.User, error)
	Register(name, email, password, confirmPassword string) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

// NewAuthService creates a new auth service instance mapping to UserRepository.
func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(email, password string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("Email dan password harus diisi")
	}

	if len(password) < 6 {
		return nil, errors.New("Password minimal 6 karakter")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil || !user.CheckPassword(password) {
		return nil, errors.New("Email atau password salah")
	}

	return user, nil
}

func (s *authService) Register(name, email, password, confirmPassword string) (*models.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("Semua field harus diisi")
	}

	if len(password) < 6 {
		return nil, errors.New("Password minimal 6 karakter")
	}

	if password != confirmPassword {
		return nil, errors.New("Password dan konfirmasi password tidak cocok")
	}

	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("Email sudah terdaftar")
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := user.HashPassword(); err != nil {
		return nil, errors.New("Gagal memproses registrasi")
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("Gagal menyimpan data pengguna")
	}

	return user, nil
}
