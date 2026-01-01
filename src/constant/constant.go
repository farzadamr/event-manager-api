package constant

const (
	// User
	AdminRoleName          string = "admin"
	AdminRoleDisplayName   string = "Administrator"
	DefaultStudentNumber   string = "4010000000"
	DefaultRoleName        string = "default"
	DefaultRoleDisplayName string = "Default User"
	TeacherRoleName        string = "teacher"
	TeacherRoleDisplayName string = "Teacher User"

	// Claims
	AuthorizationHeaderKey string = "Authorization"
	UserIdKey              string = "UserId"
	FirstNameKey           string = "FirstName"
	LastNameKey            string = "LastName"
	StudentNumberKey       string = "StudentNumber"
	EmailKey               string = "Email"
	MobileNumberKey        string = "MobileNumber"
	RolesKey               string = "Roles"
	ExpireTimeKey          string = "Exp"

	// JWT
	RefreshTokenCookieName string = "refresh_token"
)
