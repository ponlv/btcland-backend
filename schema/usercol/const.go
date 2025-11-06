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
