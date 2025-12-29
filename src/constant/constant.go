package constant

const (
	// User
	AdminRoleName          string = "admin"
	AdminRoleDisplayName   string = "Administrator"
	DefaultStudentNumber   string = "4010000000"
	DefaultRoleName        string = "default"
	DefaultRoleDisplayName string = "Default User"

	// Claims
	AuthorizationHeaderKey string = "Authorization"
	UserIdKey              string = "UserId"
	FirstNameKey           string = "FirstName"
	LastNameKey            string = "LastName"
	UsernameKey            string = "Username"
	EmailKey               string = "Email"
	MobileNumberKey        string = "MobileNumber"
	RolesKey               string = "Roles"
	ExpireTimeKey          string = "Exp"

	// JWT
	RefreshTokenCookieName string = "refresh_token"
)
