package domain

import "strings"

type RolePermission struct {
	rolePermissions map[string][]string
}

func (p RolePermission) IsAuthorizedFor(role string, routeName string) bool {
	perms := p.rolePermissions[role]
	for _, r := range perms {
		if r == strings.TrimSpace(routeName) {
			return true
		}
	}
	return false
}

func GetRolePermissions() RolePermission {
	return RolePermission{map[string][]string{
		"admin": {"NewUser"},
		"user":  {"NewUser"},
	}}
}
