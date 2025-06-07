package repositories

import (
	"context"

	"app/internal/domain/model"
)

// UserRepository define las operaciones de persistencia para los usuarios
type UserRepository interface {
	// Create crea un nuevo usuario en la base de datos
	Create(ctx context.Context, user *model.User) error
	
	// FindByID busca un usuario por su ID
	FindByID(ctx context.Context, id uint) (*model.User, error)
	
	// FindByEmail busca un usuario por su email
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	
	// FindByUsername busca un usuario por su nombre de usuario
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	
	// FindByVerifyToken busca un usuario por su token de verificación
	FindByVerifyToken(ctx context.Context, token string) (*model.User, error)
	
	// Update actualiza un usuario existente
	Update(ctx context.Context, user *model.User) error
	
	// Delete elimina un usuario de la base de datos
	Delete(ctx context.Context, id uint) error
} 