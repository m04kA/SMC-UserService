package domain

// Role представляет роль пользователя в системе
type Role string

const (
	RoleClient     Role = "client"     // Обычный клиент автомойки
	RoleManager    Role = "manager"    // Менеджер компании (автомойки)
	RoleSuperUser  Role = "superuser"  // Суперпользователь с полным доступом
)

// IsValid проверяет, является ли роль валидной
func (r Role) IsValid() bool {
	switch r {
	case RoleClient, RoleManager, RoleSuperUser:
		return true
	default:
		return false
	}
}

// CanAccessUser проверяет, может ли пользователь с данной ролью получить доступ к другому пользователю
func (r Role) CanAccessUser(targetUserID, requestUserID int64) bool {
	switch r {
	case RoleSuperUser:
		return true // Суперпользователь может видеть всех
	case RoleManager, RoleClient:
		return targetUserID == requestUserID // Клиент и менеджер могут видеть только себя
	default:
		return false
	}
}

// CanModifyUser проверяет, может ли пользователь с данной ролью модифицировать другого пользователя
func (r Role) CanModifyUser(targetUserID, requestUserID int64) bool {
	switch r {
	case RoleSuperUser:
		return true // Суперпользователь может изменять всех
	case RoleManager, RoleClient:
		return targetUserID == requestUserID // Клиент и менеджер могут изменять только себя
	default:
		return false
	}
}
