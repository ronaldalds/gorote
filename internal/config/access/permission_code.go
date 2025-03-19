package access

import "reflect"

type PermissionCode struct {
	SuperUser string `value:"super_user"`
	// Permissões de usuários
	CreateUser           string `value:"create_user"`
	ViewUser             string `value:"view_user"`
	UpdateUser           string `value:"update_user"`
	DeleteUser           string `value:"delete_user"`
	EditePermissionsUser string `value:"edite_permissions_user"`
	// Permissões de roles
	CreateRole string `value:"create_role"`
	ViewRole   string `value:"view_role"`
	UpdateRole string `value:"update_role"`
	DeleteRole string `value:"delete_role"`
}

func SetValuesFromTags(s *PermissionCode) *PermissionCode {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		valueTag := field.Tag.Get("value")

		if valueTag != "" && field.Type.Kind() == reflect.String {
			v.Field(i).SetString(valueTag)
		}
	}
	return s
}

var Permissions PermissionCode
