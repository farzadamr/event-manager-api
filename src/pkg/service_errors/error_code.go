package service_errors

const (
	// Token
	UnExpectedError     = "Expected error"
	ClaimsNotFound      = "Claims not found"
	TokenRequired       = "token required"
	TokenExpired        = "token expired"
	TokenInvalid        = "token invalid"
	InvalidRefreshToken = "invalid refresh token"

	// User
	EmailExists               = "Email exists"
	StudentNumberExists       = "Student number exists"
	PermissionDenied          = "Permission denied"
	UsernameOrPasswordInvalid = "username or password invalid"
	InvalidRolesFormat        = "invalid roles format"

	// DB
	RecordNotFound = "record not found"

	// Attendance
	AttendanceIsNotPresent = "attendance is not present"

	// Certificate
	CertificateInvalid = "certificate invalid"

	// Registration
	NoCapacity = "not enough capacity"
)
