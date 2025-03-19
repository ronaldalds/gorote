package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ronaldalds/base-go-api/internal/config/access"
	"github.com/ronaldalds/base-go-api/internal/config/databases"
	"github.com/ronaldalds/base-go-api/internal/config/envs"
	"gorm.io/gorm"
)

type HealthHandler struct {
	Sql   map[string]string `json:"sql"`
	Redis map[string]string `json:"redis"`
	NoSql map[string]string `json:"nosql"`
}

type Service struct {
	GormStore  *databases.GormStore
	RedisStore *databases.RedisStore
	MongoStore *databases.MongoStore
}

func NewService() *Service {
	service := &Service{
		GormStore:  databases.DB.GormStore,
		RedisStore: databases.DB.RedisStore,
		MongoStore: databases.DB.MongoStore,
	}
	// Executar as Migrations
	service.GormStore.DB.AutoMigrate(&User{}, &Role{}, &Permission{})
	// Executar as Seeds
	if err := service.SeedUserAdmin(); err != nil {
		fmt.Println(err.Error())
	}
	if err := service.SeedPermissions(); err != nil {
		fmt.Println(err.Error())
	}
	return service
}

func (s *Service) Login(req *Login) (*User, error) {
	var user User
	result := s.GormStore.DB.
		Preload("Roles.Permissions").
		Where("username = ? OR email = ?", req.Username, req.Username).
		First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to login: username or password is incorrect")
	}
	if !CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("failed to login: username or password is incorrect")
	}
	if !user.Active {
		return nil, fmt.Errorf("failed to login: user is inactive")
	}
	return &user, nil
}

func (s *Service) SeedUserAdmin() error {
	var user User
	err := s.GormStore.DB.Where("username = ?", envs.Env.SuperUsername).First(&user).Error
	if err == nil {
		return fmt.Errorf("admin already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check admin existence: %s", err.Error())
	}
	hashPassword, err := HashPassword(envs.Env.SuperPass)
	if err != nil {
		return fmt.Errorf("failed to create admin: %s", err.Error())
	}
	admin := &User{
		FirstName:   envs.Env.SuperName,
		LastName:    "Admin",
		Username:    envs.Env.SuperUsername,
		Email:       envs.Env.SuperEmail,
		Password:    hashPassword,
		Active:      true,
		IsSuperUser: true,
		Phone1:      envs.Env.SuperPhone,
	}
	if err := s.GormStore.DB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create user: %s", err.Error())
	}
	log.Println("admin created successfully")
	return nil
}

func (s *Service) SeedPermissions() error {
	v := reflect.ValueOf(access.Permissions)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).String() != "" && t.Field(i).Type.Kind() == reflect.String {
			var item Permission
			err := s.GormStore.DB.Where("code = ?", v.Field(i)).First(&item).Error
			if err == nil {
				log.Printf("permission with code '%s' already exists \n", v.Field(i).String())
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to check permission existence: %s", err.Error())
			}

			permission := &Permission{
				Name:        t.Field(i).Name,
				Code:        v.Field(i).String(),
				Description: t.Field(i).Tag.Get("description"),
			}
			if err := s.GormStore.DB.Create(&permission).Error; err != nil {
				return fmt.Errorf("failed to create permission: %s", err.Error())
			}
		}
	}
	log.Println("permissions created successfully")
	return nil
}

func (s *Service) GetKey(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := s.RedisStore.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", fmt.Errorf("failed to get key: %v", err)
	}
	return val, nil
}

func (s *Service) SetKeyHash(key string, value map[string]any, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.RedisStore.Client.HSet(ctx, key, value).Err(); err != nil {
		return fmt.Errorf("failed to set key: %v", err)
	}

	if err := s.RedisStore.Client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set expiration for key: %v", err)
	}

	return nil
}

func (s *Service) SetToken(id uint, access string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.RedisStore.Client.Set(
		ctx,
		fmt.Sprintf("%d", id),
		access,
		ttl,
	).Err(); err != nil {
		return fmt.Errorf("failed to set key: %v", err)
	}
	return nil
}

func (s *Service) ListPermission(permissions *[]Permission) error {
	result := s.GormStore.DB.
		Find(permissions)
	if result.Error != nil {
		return fmt.Errorf("failed to query database: %w", result.Error)
	}
	return nil
}

func (s *Service) GetPermissionByIds(permissions *[]Permission, ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("no permission IDs provided")
	}

	// Buscar as permissões pelos IDs fornecidos
	if err := s.GormStore.DB.Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return fmt.Errorf("failed to fetch permissions: %s", err.Error())
	}

	// Verificar se todas as permissões foram encontradas
	if len(*permissions) != len(ids) {
		return fmt.Errorf("permissions not found for IDs")
	}

	return nil
}

func (s *Service) ListRole(roles *[]Role) error {
	result := s.GormStore.DB.
		Preload("Permissions").
		Find(roles)
	if result.Error != nil {
		return fmt.Errorf("failed to query database: %w", result.Error)
	}
	return nil
}

func (s *Service) CreateRole(role *Role, req *CreateRole) error {
	var permissions []Permission
	if err := s.GetPermissionByIds(&permissions, req.Permissions); err != nil {
		return fmt.Errorf("permission with ids '%v' does not exist", req.Permissions)
	}

	role.Name = req.Name
	role.Permissions = permissions // Associar permissões à role
	role.Description = req.Description

	if err := s.GormStore.DB.Create(role).Error; err != nil {
		return fmt.Errorf("failed to create role: %s", err.Error())
	}
	return nil
}

func (s *Service) GetRoleByIds(ids []uint) ([]Role, error) {
	var roles []Role
	if len(ids) == 0 {
		return roles, nil
	}

	// Buscar as permissões pelos IDs fornecidos
	if err := s.GormStore.DB.Where("id IN ?", ids).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %s", err.Error())
	}

	// Verificar se todas as permissões foram encontradas
	if len(roles) != len(ids) {
		return nil, fmt.Errorf("roles not found for IDs")
	}

	return roles, nil
}

func (s *Service) GetUserByID(id uint) (*User, error) {
	var user User
	result := s.GormStore.DB.
		Preload("Roles.Permissions").
		Where("id = ?", id).
		First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no record found for id: %d", id)
		}
		return nil, fmt.Errorf("failed to query database: %w", result.Error)
	}
	return &user, nil
}

func (s *Service) CreateUser(creatorID uint, req *CreateUser) (*User, error) {
	// Buscar o criador do usuário
	creator, err := s.GetUserByID(creatorID)
	if err != nil {
		return nil, fmt.Errorf("user with id '%v' does not exist", creatorID)
	}

	// Buscar as roles pelo ID
	roles, err := s.GetRoleByIds(req.Roles)
	if err != nil {
		return nil, fmt.Errorf("group with ids '%v' does not exist", req.Roles)
	}

	// Validar se o criador possui as roles necessárias ou é superusuário
	if !creator.IsSuperUser && !ContainsAll(creator.Roles, roles) {
		return nil, fmt.Errorf("failed to create user: creator does not have all required roles")
	}

	// Criar o usuário (apenas em memória)
	user := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Active:    req.Active,
		IsSuperUser: func() bool {
			if req.IsSuperUser {
				if creator.IsSuperUser {
					return true
				}
				panic(fmt.Errorf("only superusers can create other superusers"))
			}
			return false
		}(),
		Phone1: req.Phone1,
		Phone2: req.Phone2,
	}

	// Persistir o usuário no banco de dados
	if err := s.GormStore.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	// Associar as roles ao usuário
	if err := s.GormStore.DB.Model(&user).Association("Roles").Replace(roles); err != nil {
		return nil, fmt.Errorf("failed to set roles for user")
	}

	// Retornar o usuário criado
	return &user, nil
}

func (s *Service) UpdateUser(editorID uint, id uint, req *UserSchema) (*User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user modified with id '%v' does not exists", id)
	}
	// var editor models.User
	editor, err := s.GetUserByID(editorID)
	if err != nil {
		return nil, fmt.Errorf("user editor with id '%v' does not exist", editorID)
	}
	// Verificar se o editor pode atualizar o usuário
	canUpdateFull := editor.IsSuperUser || slices.Contains(ExtractCodePermissionsByUser(editor), access.Permissions.UpdateUser)
	// Se o editor não tiver permissão para atualizar o usuário, ele só pode atualizar a si mesmo
	if !canUpdateFull && user.ID != editorID {
		return nil, fmt.Errorf("user editor with id '%v' does not have permission to update user with id '%v'", editorID, id)
	}
	// Se o editor tiver permissão para atualizar o usuário, atualize o usuário
	if canUpdateFull {
		if err := s.UpdateFullUser(user, editor, req); err != nil {
			return nil, err
		}
	} else {
		// Caso contrário, o editor só pode atualizar a si mesmo
		if err := s.UpdateSimpleUser(user, req); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *Service) UpdateSimpleUser(user *User, req *UserSchema) error {
	// Atualizar outros campos do usuário
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Active = req.Active
	user.Phone1 = req.Phone1
	user.Phone2 = req.Phone2

	// Salvar as alterações
	if err := s.GormStore.DB.Model(user).Updates(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %s", err.Error())
	}
	return nil
}

func (s *Service) UpdateFullUser(editor *User, user *User, req *UserSchema) error {
	// Atualizar as Roles somente se permitido
	if len(req.Roles) > 0 {
		// Buscar as roles especificadas na atualização
		roles, err := s.GetRoleByIds(req.Roles)
		if err != nil {
			return fmt.Errorf("role with ids '%v' does not exist", req.Roles)
		}
		// Validar se o criador possui as roles necessárias ou é superusuário
		if !editor.IsSuperUser {
			if !ContainsAll(editor.Roles, roles) {
				return fmt.Errorf("failed to update user: editor does not have all required roles")
			}
		}

		// Atualizar as roles do usuário
		if err := s.GormStore.DB.Model(user).Association("Roles").Replace(roles); err != nil {
			return fmt.Errorf("failed to set roles for user: %v", err)
		}
	}

	// Atualizar outros campos do usuário
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Username = req.Username
	user.Email = req.Email
	user.Active = req.Active
	if req.IsSuperUser {
		if !editor.IsSuperUser {
			return fmt.Errorf("only superusers can update other superusers")
		}
		user.IsSuperUser = true
	}
	user.Phone1 = req.Phone1
	user.Phone2 = req.Phone2

	// Salvar as alterações
	if err := s.GormStore.DB.Model(user).Updates(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %s", err.Error())
	}
	return nil
}

func (s *Service) ListUser() ([]User, error) {
	var users []User
	result := s.GormStore.DB.
		Preload("Roles.Permissions").
		Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query database list: %w", result.Error)
	}
	return users, nil
}
