package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"app/internal/domain/model"
	"app/internal/domain/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase implementa la lógica de negocio relacionada con usuarios
// (registro, autenticación, verificación, etc.)
//
//	Esta implementación es la que requieren los handlers y middlewares actuales.
type UserUseCase struct {
	userRepo     repositories.UserRepository
	emailService EmailService
}

// EmailService es la interfaz del servicio de envío de correos (Mailer)
// Se define aquí para desacoplar el caso de uso de la implementación concreta.
type EmailService interface {
	SendVerificationEmail(to, token, username string) error
	SendPasswordResetEmail(to, token, username string) error
}

// NewUserUseCase devuelve una nueva instancia del caso de uso de usuarios.
func NewUserUseCase(userRepo repositories.UserRepository, emailSvc EmailService) *UserUseCase {
	return &UserUseCase{
		userRepo:     userRepo,
		emailService: emailSvc,
	}
}

// AuthenticateUser verifica email + contraseña y devuelve el usuario si es correcto.
func (uc *UserUseCase) AuthenticateUser(ctx context.Context, email, password string) (*model.User, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	return user, nil
}

// CreateUser crea un nuevo usuario y genera hash de contraseña y token de verificación.
func (uc *UserUseCase) CreateUser(ctx context.Context, user *model.User, plainPassword string) error {
	// Validaciones básicas
	if user.Email == "" || user.Username == "" || plainPassword == "" {
		return errors.New("datos incompletos")
	}

	// ¿Existe email?
	if existing, _ := uc.userRepo.FindByEmail(ctx, user.Email); existing != nil {
		return errors.New("el correo ya está registrado")
	}
	// ¿Existe username?
	if existing, _ := uc.userRepo.FindByUsername(ctx, user.Username); existing != nil {
		return errors.New("el nombre de usuario ya está en uso")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al generar contraseña: %w", err)
	}

	user.PasswordHash = string(hash)
	user.VerifyToken = uuid.New().String()
	user.Verified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

// SendVerificationEmail envía el correo de verificación.
func (uc *UserUseCase) SendVerificationEmail(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("usuario nulo")
	}
	return uc.emailService.SendVerificationEmail(user.Email, user.VerifyToken, user.Username)
}

// VerifyUser comprueba el token y marca al usuario como verificado.
func (uc *UserUseCase) VerifyUser(ctx context.Context, token string) (*model.User, error) {
	usr, err := uc.userRepo.FindByVerifyToken(ctx, token)
	if err != nil {
		return nil, errors.New("token inválido")
	}
	usr.Verified = true
	usr.VerifyToken = ""
	usr.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, usr); err != nil {
		return nil, err
	}
	return usr, nil
}

// GetUserByID devuelve el usuario por su ID.
func (uc *UserUseCase) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

// GetUserByEmail devuelve el usuario por su Email.
func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return uc.userRepo.FindByEmail(ctx, email)
}

// UpdateUserSettings actualiza las configuraciones del usuario
func (uc *UserUseCase) UpdateUserSettings(ctx context.Context, userID uint, emailNotifications bool) error {
	// Buscar el usuario
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Actualizar el campo de notificaciones por email
	user.EmailNotifications = emailNotifications
	user.UpdatedAt = time.Now()

	// Guardar los cambios
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error al actualizar configuración: %w", err)
	}

	return nil
}

// InitiatePasswordReset genera token y envía correo de restablecimiento
func (uc *UserUseCase) InitiatePasswordReset(ctx context.Context, userID uint) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("[ERROR] No se pudo encontrar el usuario ID=%d para restablecimiento: %v", userID, err)
		return fmt.Errorf("usuario no encontrado: %w", err)
	}

	resetToken := uuid.New().String()
	user.VerifyToken = resetToken // reutilizamos campo existente
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		log.Printf("[ERROR] No se pudo actualizar el token de restablecimiento para usuario ID=%d: %v", userID, err)
		return fmt.Errorf("error al guardar token de restablecimiento: %w", err)
	}

	log.Printf("[INFO] Enviando correo de restablecimiento a usuario ID=%d email=%s", userID, user.Email)
	// Enviar correo con el enlace para restablecer contraseña
	if err := uc.emailService.SendPasswordResetEmail(user.Email, resetToken, user.Username); err != nil {
		log.Printf("[ERROR] No se pudo enviar correo de restablecimiento a usuario ID=%d email=%s: %v", userID, user.Email, err)
		return fmt.Errorf("error al enviar correo de restablecimiento: %w", err)
	}
	log.Printf("[INFO] Correo de restablecimiento enviado exitosamente a usuario ID=%d email=%s", userID, user.Email)

	return nil
}

// ResetPassword verifica el token y restablece la contraseña
func (uc *UserUseCase) ResetPassword(ctx context.Context, token, newPassword string) (*model.User, error) {
	user, err := uc.userRepo.FindByVerifyToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token inválido o expirado: %w", err)
	}

	// Hash de la nueva contraseña
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al generar hash de contraseña: %w", err)
	}

	// Actualizar contraseña y limpiar token
	user.PasswordHash = string(passwordHash)
	user.VerifyToken = ""
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	return user, nil
}

// ChangePassword cambia la contraseña del usuario después de verificar la contraseña actual
func (uc *UserUseCase) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	log.Printf("[INFO] UserUseCase.ChangePassword - Iniciando cambio de contraseña para userID: %d", userID)

	// Buscar el usuario
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("[ERROR] UserUseCase.ChangePassword - Error al buscar usuario ID %d: %v", userID, err)
		return fmt.Errorf("usuario no encontrado: %w", err)
	}
	log.Printf("[DEBUG] UserUseCase.ChangePassword - Usuario encontrado: %+v", user)

	// Verificar la contraseña actual
	log.Printf("[DEBUG] UserUseCase.ChangePassword - Verificando contraseña actual para userID: %d. Hash almacenado: %s", userID, user.PasswordHash)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		log.Printf("[ERROR] UserUseCase.ChangePassword - Contraseña actual incorrecta para userID %d. Error de comparación: %v", userID, err)
		return errors.New("la contraseña actual es incorrecta")
	}
	log.Printf("[INFO] UserUseCase.ChangePassword - Contraseña actual verificada correctamente para userID: %d", userID)

	// Hash de la nueva contraseña
	log.Printf("[DEBUG] UserUseCase.ChangePassword - Generando hash para nueva contraseña para userID: %d", userID)
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] UserUseCase.ChangePassword - Error al generar hash para nueva contraseña para userID %d: %v", userID, err)
		return fmt.Errorf("error al generar hash de contraseña: %w", err)
	}
	log.Printf("[DEBUG] UserUseCase.ChangePassword - Nuevo hash generado para userID: %d", userID)

	// Actualizar contraseña
	user.PasswordHash = string(newPasswordHash)
	user.UpdatedAt = time.Now()

	log.Printf("[DEBUG] UserUseCase.ChangePassword - Intentando actualizar usuario en DB para userID: %d con nuevo hash.", userID)
	if err := uc.userRepo.Update(ctx, user); err != nil {
		log.Printf("[ERROR] UserUseCase.ChangePassword - Error al actualizar contraseña en DB para userID %d: %v", userID, err)
		return fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	log.Printf("[INFO] UserUseCase.ChangePassword - Contraseña actualizada exitosamente para userID: %d", userID)
	return nil
}

// DeleteAccount elimina la cuenta del usuario después de verificar la contraseña
func (uc *UserUseCase) DeleteAccount(ctx context.Context, userID uint, password string) error {
	// Buscar el usuario
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Verificar la contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return errors.New("contraseña incorrecta")
	}

	// Eliminar el usuario
	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("error al eliminar la cuenta: %w", err)
	}

	return nil
}

// Register crea un nuevo usuario en el sistema
func (uc *UserUseCase) Register(ctx context.Context, username, email, password string, notifyByEmail bool) (*model.User, error) {
	// Generar un token de verificación
	verifyToken := uuid.New().String()

	// Crear el usuario
	user := &model.User{
		Username:           username,
		Email:              email,
		VerifyToken:        verifyToken,
		Verified:           false,
		EmailNotifications: notifyByEmail,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Generar hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al generar hash de contraseña: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	// Intentar crear el usuario
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error al crear usuario: %w", err)
	}

	return user, nil
}
