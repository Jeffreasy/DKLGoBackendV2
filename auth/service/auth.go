package service

import (
	"dklautomationgo/database/repository"
	"dklautomationgo/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials   = errors.New("ongeldige inloggegevens")
	ErrUserNotFound         = errors.New("gebruiker niet gevonden")
	ErrUserNotActive        = errors.New("gebruiker is niet actief")
	ErrInvalidToken         = errors.New("ongeldige token")
	ErrTokenExpired         = errors.New("token is verlopen")
	ErrPasswordResetExpired = errors.New("wachtwoord reset link is verlopen")
	ErrPasswordTooWeak      = errors.New("wachtwoord voldoet niet aan de vereisten")
)

// AuthService bevat de business logic voor authenticatie
type AuthService struct {
	userRepo     *repository.UserRepository
	tokenService *TokenService
}

// NewAuthService maakt een nieuwe AuthService
func NewAuthService(userRepo *repository.UserRepository, tokenService *TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Login authenticeert een gebruiker en geeft tokens terug
func (s *AuthService) Login(email, password string) (*models.TokenResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("[AuthService] Error finding user by email: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Controleer wachtwoord
	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	// Controleer of gebruiker actief is
	if user.Status != models.StatusActive {
		return nil, ErrUserNotActive
	}

	// Update laatste login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		log.Printf("[AuthService] Error updating last login: %v", err)
		// Niet fataal, ga door
	}

	// Genereer tokens
	return s.generateTokens(user)
}

// RefreshToken vernieuwt een access token met een refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenResponse, error) {
	// Valideer refresh token
	token, err := s.userRepo.FindRefreshToken(refreshToken)
	if err != nil {
		log.Printf("[AuthService] Error finding refresh token: %v", err)
		return nil, err
	}
	if token == nil {
		return nil, ErrInvalidToken
	}

	// Controleer of token verlopen is
	if token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// Haal gebruiker op
	user, err := s.userRepo.FindByID(token.UserID)
	if err != nil {
		log.Printf("[AuthService] Error finding user by ID: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Controleer of gebruiker actief is
	if user.Status != models.StatusActive {
		return nil, ErrUserNotActive
	}

	// Herroep oude token
	if err := s.userRepo.RevokeRefreshToken(refreshToken); err != nil {
		log.Printf("[AuthService] Error revoking refresh token: %v", err)
		// Niet fataal, ga door
	}

	// Genereer nieuwe tokens
	return s.generateTokens(user)
}

// Logout logt een gebruiker uit door refresh token te herroepen
func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepo.RevokeRefreshToken(refreshToken)
}

// LogoutAll logt een gebruiker uit op alle apparaten
func (s *AuthService) LogoutAll(userID uuid.UUID) error {
	return s.userRepo.RevokeAllUserRefreshTokens(userID)
}

// CreateUser maakt een nieuwe gebruiker aan
func (s *AuthService) CreateUser(email, password string, role models.UserRole) (*models.User, error) {
	// Controleer of email al bestaat
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("[AuthService] Error checking existing user: %v", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("email is al in gebruik")
	}

	// Valideer wachtwoord
	if err := s.validatePassword(password); err != nil {
		return nil, err
	}

	// Maak nieuwe gebruiker
	user := &models.User{
		Email:  email,
		Role:   role,
		Status: models.StatusPending,
	}

	// Set password
	if err := user.SetPassword(password); err != nil {
		log.Printf("[AuthService] Error setting password: %v", err)
		return nil, err
	}

	// Sla gebruiker op
	if err := s.userRepo.Create(user); err != nil {
		log.Printf("[AuthService] Error creating user: %v", err)
		return nil, err
	}

	return user, nil
}

// ApproveUser keurt een gebruiker goed
func (s *AuthService) ApproveUser(userID, approverID uuid.UUID) error {
	// Controleer of gebruiker bestaat
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[AuthService] Error finding user: %v", err)
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Controleer of approver bestaat
	approver, err := s.userRepo.FindByID(approverID)
	if err != nil {
		log.Printf("[AuthService] Error finding approver: %v", err)
		return err
	}
	if approver == nil {
		return ErrUserNotFound
	}

	// Controleer of approver een beheerder is
	if approver.Role != models.RoleBeheerder {
		return fmt.Errorf("alleen beheerders kunnen gebruikers goedkeuren")
	}

	// Keur gebruiker goed
	return s.userRepo.ApproveUser(userID, approverID)
}

// UpdateUser werkt een gebruiker bij
func (s *AuthService) UpdateUser(userID uuid.UUID, updates *models.UpdateUserRequest) error {
	// Controleer of gebruiker bestaat
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[AuthService] Error finding user: %v", err)
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Update velden
	if updates.Email != nil {
		// Controleer of email al bestaat
		if *updates.Email != user.Email {
			existingUser, err := s.userRepo.FindByEmail(*updates.Email)
			if err != nil {
				log.Printf("[AuthService] Error checking existing user: %v", err)
				return err
			}
			if existingUser != nil {
				return fmt.Errorf("email is al in gebruik")
			}
			user.Email = *updates.Email
		}
	}

	if updates.Role != nil {
		user.Role = *updates.Role
	}

	if updates.Status != nil {
		user.Status = *updates.Status
	}

	// Sla gebruiker op
	return s.userRepo.Update(user)
}

// ChangePassword wijzigt het wachtwoord van een gebruiker
func (s *AuthService) ChangePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	// Controleer of gebruiker bestaat
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[AuthService] Error finding user: %v", err)
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Controleer huidig wachtwoord
	if !user.CheckPassword(currentPassword) {
		return ErrInvalidCredentials
	}

	// Valideer nieuw wachtwoord
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// Set nieuw wachtwoord
	if err := user.SetPassword(newPassword); err != nil {
		log.Printf("[AuthService] Error setting new password: %v", err)
		return err
	}

	// Sla gebruiker op
	if err := s.userRepo.Update(user); err != nil {
		log.Printf("[AuthService] Error updating user: %v", err)
		return err
	}

	// Herroep alle refresh tokens
	return s.userRepo.RevokeAllUserRefreshTokens(userID)
}

// ForgotPassword start het wachtwoord reset proces
func (s *AuthService) ForgotPassword(email string) (string, error) {
	// Controleer of gebruiker bestaat
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("[AuthService] Error finding user: %v", err)
		return "", err
	}
	if user == nil {
		// Geef geen fout om privacy redenen
		return "", nil
	}

	// Genereer reset token
	token := uuid.New()
	expires := time.Now().Add(24 * time.Hour) // 24 uur geldig

	// Sla token op
	if err := s.userRepo.SetPasswordResetToken(user.ID, token, expires); err != nil {
		log.Printf("[AuthService] Error setting password reset token: %v", err)
		return "", err
	}

	return token.String(), nil
}

// ResetPassword reset het wachtwoord met een token
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Controleer of token geldig is
	user, err := s.userRepo.FindByPasswordResetToken(token)
	if err != nil {
		log.Printf("[AuthService] Error finding user by reset token: %v", err)
		return err
	}
	if user == nil {
		return ErrPasswordResetExpired
	}

	// Valideer nieuw wachtwoord
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// Set nieuw wachtwoord
	if err := user.SetPassword(newPassword); err != nil {
		log.Printf("[AuthService] Error setting new password: %v", err)
		return err
	}

	// Wis reset token
	if err := s.userRepo.ClearPasswordResetToken(user.ID); err != nil {
		log.Printf("[AuthService] Error clearing password reset token: %v", err)
		// Niet fataal, ga door
	}

	// Sla gebruiker op
	if err := s.userRepo.Update(user); err != nil {
		log.Printf("[AuthService] Error updating user: %v", err)
		return err
	}

	// Herroep alle refresh tokens
	return s.userRepo.RevokeAllUserRefreshTokens(user.ID)
}

// GetUserByID haalt een gebruiker op op ID
func (s *AuthService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByEmail haalt een gebruiker op op email
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

// GetAllUsers haalt alle gebruikers op
func (s *AuthService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

// DeleteUser verwijdert een gebruiker
func (s *AuthService) DeleteUser(userID uuid.UUID, deleterID uuid.UUID) error {
	// Controleer of gebruiker bestaat
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[AuthService] Error finding user: %v", err)
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Controleer of deleter bestaat
	deleter, err := s.userRepo.FindByID(deleterID)
	if err != nil {
		log.Printf("[AuthService] Error finding deleter: %v", err)
		return err
	}
	if deleter == nil {
		return ErrUserNotFound
	}

	// Controleer of deleter een beheerder is
	if deleter.Role != models.RoleBeheerder {
		return fmt.Errorf("alleen beheerders kunnen gebruikers verwijderen")
	}

	// Verwijder gebruiker
	return s.userRepo.DeleteByID(userID)
}

// AdminChangePassword stelt een beheerder in staat om het wachtwoord van een gebruiker te wijzigen
// zonder het huidige wachtwoord te kennen
func (s *AuthService) AdminChangePassword(userID uuid.UUID, adminID uuid.UUID, newPassword string) error {
	// Controleer of gebruiker bestaat
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[AuthService] Error finding user: %v", err)
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Controleer of admin bestaat
	admin, err := s.userRepo.FindByID(adminID)
	if err != nil {
		log.Printf("[AuthService] Error finding admin: %v", err)
		return err
	}
	if admin == nil {
		return ErrUserNotFound
	}

	// Controleer of admin een beheerder is
	if admin.Role != models.RoleBeheerder {
		return fmt.Errorf("alleen beheerders kunnen wachtwoorden wijzigen")
	}

	// Valideer nieuw wachtwoord
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// Set nieuw wachtwoord
	if err := user.SetPassword(newPassword); err != nil {
		log.Printf("[AuthService] Error setting new password: %v", err)
		return err
	}

	// Sla gebruiker op
	if err := s.userRepo.Update(user); err != nil {
		log.Printf("[AuthService] Error updating user: %v", err)
		return err
	}

	// Herroep alle refresh tokens
	return s.userRepo.RevokeAllUserRefreshTokens(userID)
}

// GetUserRepository geeft de user repository terug
func (s *AuthService) GetUserRepository() *repository.UserRepository {
	return s.userRepo
}

// Interne hulpfuncties

// generateTokens genereert access en refresh tokens
func (s *AuthService) generateTokens(user *models.User) (*models.TokenResponse, error) {
	// Genereer access token
	accessToken, err := s.tokenService.GenerateAccessToken(user)
	if err != nil {
		log.Printf("[AuthService] Error generating access token: %v", err)
		return nil, err
	}

	// Bepaal refresh token expiratie
	refreshExpiry := getRefreshTokenExpiry()

	// Genereer refresh token
	refreshTokenString := uuid.New().String()
	refreshExpires := time.Now().Add(refreshExpiry)

	// Sla refresh token op
	_, err = s.userRepo.CreateRefreshToken(user.ID, refreshTokenString, refreshExpires)
	if err != nil {
		log.Printf("[AuthService] Error creating refresh token: %v", err)
		return nil, err
	}

	// Maak token response
	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int(getAccessTokenExpiry().Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// validatePassword valideert een wachtwoord
func (s *AuthService) validatePassword(password string) error {
	minLength := getPasswordMinLength()
	if len(password) < minLength {
		return fmt.Errorf("wachtwoord moet minimaal %d karakters bevatten", minLength)
	}

	// Controleer op hoofdletter
	if getPasswordRequireUppercase() && !containsUppercase(password) {
		return fmt.Errorf("wachtwoord moet minimaal één hoofdletter bevatten")
	}

	// Controleer op kleine letter
	if getPasswordRequireLowercase() && !containsLowercase(password) {
		return fmt.Errorf("wachtwoord moet minimaal één kleine letter bevatten")
	}

	// Controleer op cijfer
	if getPasswordRequireNumber() && !containsNumber(password) {
		return fmt.Errorf("wachtwoord moet minimaal één cijfer bevatten")
	}

	// Controleer op speciaal teken
	if getPasswordRequireSpecial() && !containsSpecial(password) {
		return fmt.Errorf("wachtwoord moet minimaal één speciaal teken bevatten")
	}

	return nil
}

// Hulpfuncties voor wachtwoord validatie
func containsUppercase(s string) bool {
	for _, r := range s {
		if 'A' <= r && r <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, r := range s {
		if 'a' <= r && r <= 'z' {
			return true
		}
	}
	return false
}

func containsNumber(s string) bool {
	for _, r := range s {
		if '0' <= r && r <= '9' {
			return true
		}
	}
	return false
}

func containsSpecial(s string) bool {
	specials := "!@#$%^&*()_+-=[]{}|;:,.<>?/~`"
	for _, r := range s {
		for _, sp := range specials {
			if r == sp {
				return true
			}
		}
	}
	return false
}

// Configuratie helpers
func getAccessTokenExpiry() time.Duration {
	expiryStr := os.Getenv("JWT_ACCESS_TOKEN_EXPIRY")
	if expiryStr == "" {
		return 15 * time.Minute // Default: 15 minuten
	}

	duration, err := time.ParseDuration(expiryStr)
	if err != nil {
		log.Printf("[AuthService] Error parsing JWT_ACCESS_TOKEN_EXPIRY: %v, using default", err)
		return 15 * time.Minute
	}

	return duration
}

func getRefreshTokenExpiry() time.Duration {
	expiryStr := os.Getenv("JWT_REFRESH_TOKEN_EXPIRY")
	if expiryStr == "" {
		return 7 * 24 * time.Hour // Default: 7 dagen
	}

	duration, err := time.ParseDuration(expiryStr)
	if err != nil {
		log.Printf("[AuthService] Error parsing JWT_REFRESH_TOKEN_EXPIRY: %v, using default", err)
		return 7 * 24 * time.Hour
	}

	return duration
}

func getPasswordMinLength() int {
	minLengthStr := os.Getenv("PASSWORD_MIN_LENGTH")
	if minLengthStr == "" {
		return 8 // Default: 8 karakters
	}

	minLength, err := strconv.Atoi(minLengthStr)
	if err != nil {
		log.Printf("[AuthService] Error parsing PASSWORD_MIN_LENGTH: %v, using default", err)
		return 8
	}

	return minLength
}

func getPasswordRequireUppercase() bool {
	return getEnvBool("PASSWORD_REQUIRE_UPPERCASE", true)
}

func getPasswordRequireLowercase() bool {
	return getEnvBool("PASSWORD_REQUIRE_LOWERCASE", true)
}

func getPasswordRequireNumber() bool {
	return getEnvBool("PASSWORD_REQUIRE_NUMBER", true)
}

func getPasswordRequireSpecial() bool {
	return getEnvBool("PASSWORD_REQUIRE_SPECIAL", true)
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("[AuthService] Error parsing %s: %v, using default", key, err)
		return defaultValue
	}

	return boolValue
}
