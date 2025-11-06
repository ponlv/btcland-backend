package usercol

import "fmt"

type UserType string

const (
	UserTypeGeneral   UserType = "GENERAL"
	UserTypeStartup   UserType = "STARTUP"
	UserTypeMentor    UserType = "MENTOR"
	UserTypeIncubator UserType = "INCUBATOR"
	UserTypeFund      UserType = "FUND"
)

func StringToUserType(t string) (UserType, error) {
	switch t {
	case "GENERAL":
		return UserTypeGeneral, nil
	case "STARTUP":
		return UserTypeStartup, nil
	case "MENTOR":
		return UserTypeMentor, nil
	case "INCUBATOR":
		return UserTypeIncubator, nil
	case "FUND":
		return UserTypeFund, nil
	default:
		return "", fmt.Errorf("unknown UserType: %s", t)
	}
}

func (i UserType) String() string {
	return string(i)
}

func (i UserType) Text() string {
	switch i {
	case UserTypeGeneral:
		return "Người dùng thường"
	case UserTypeStartup:
		return "Startup"
	case UserTypeMentor:
		return "Tư vấn viên"
	case UserTypeIncubator:
		return "Vườn ươm"
	case UserTypeFund:
		return "Quỹ đầu tư"
	default:
		return ""
	}
}

type Role string

const (
	RoleEmployee          Role = "employee"
	RoleManager           Role = "manager"
	RoleLeader            Role = "leader"
	RoleAssistantDirector Role = "assistant_director"
)

func StringToRole(r string) (Role, error) {
	switch r {
	case "employee":
		return RoleEmployee, nil
	case "manager":
		return RoleManager, nil
	case "leader":
		return RoleLeader, nil
	case "assistant_director":
		return RoleAssistantDirector, nil
	default:
		return "", fmt.Errorf("unknown Role: %s", r)
	}
}

func (r Role) String() string {
	return string(r)
}

func (r Role) Text() string {
	switch r {
	case RoleEmployee:
		return "Nhân viên"
	case RoleManager:
		return "Quản lý"
	case RoleLeader:
		return "Lãnh đạo"
	case RoleAssistantDirector:
		return "Trợ lý giám đốc"
	default:
		return ""
	}
}
