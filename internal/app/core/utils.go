package core

import (
	"fmt"
	"unicode"

	"github.com/ronaldalds/base-go-api/internal/config/handlers"
	"golang.org/x/crypto/bcrypt"
)

func ExtractNameRolesByUser(user User) []uint {
	var data []uint
	for _, role := range user.Roles {
		data = append(data, role.ID)
	}
	return data
}

func ExtractCodePermissionsByUser(user *User) []string {
	var codePermissions []string
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			codePermissions = append(codePermissions, permission.Code)
		}
	}
	return codePermissions
}

func ContainsAll(listX, listY []Role) bool {
	// Criar um mapa para os itens de X
	itemMap := make(map[uint]bool)
	for _, item := range listX {
		itemMap[item.ID] = true
	}

	// Verificar se todos os itens de Y estão no mapa de X
	for _, item := range listY {
		if !itemMap[item.ID] {
			return false // Item de Y não está em X
		}
	}

	return true // Todos os itens de Y estão em X
}

func HashPassword(password string) (string, error) {
	// Exemplo usando bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ValidatePassword(password string) *handlers.ErrHandler {
	errors := handlers.NewError()

	// Verificar se contém uma letra maiúscula
	hasUpper := false
	hasSymbol := false

	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsSymbol(r) || unicode.IsPunct(r) { // Símbolos e pontuações
			hasSymbol = true
		}
	}

	// Adicionar erros caso os critérios não sejam atendidos
	if !hasUpper {
		errors.AddDetailErr("uppercase", "password must contain at least one uppercase letter.")
	}
	if !hasSymbol {
		errors.AddDetailErr("symbol", "password must contain at least one symbol.")
	}

	// Retornar nil se não houver erros
	if hasSymbol || hasUpper {
		return nil
	}

	return errors
}
