package config

// DatabaseURL is the env variable name for the database url
const DatabaseURL = "DATABASE_URL"

// AuthServerPort is the env variable name for the port to use for the auth server
const AuthServerPort = "AUTH_SERVER_PORT"

// SessionCookieName is the env variable name used to set the cookie for sessions
const SessionKey = "DISCUSSION_APP_SESSION_KEY"

// SessionCookieName is the env variable name used to set the cookie for sessions
const SessionCookieName = "DISCUSSION_APP_SESSION_COOKIE"

// SessionExpiration is the time in seconds when a token will expire
const SessionExpiration = 3600 * 24 * 7


// MinEntropyBits is the minimum number of bits of entropy required for a password.
const MinEntropyBits = 64

// MaxLoginAttempts is the maximum number of times user can attempt to enter the correct password
// before their account is temporarily locked
const MaxLoginAttempts = 5

// AccountLockoutLength is the time in seconds that an account will be locked
const AccountLockoutLength = 60 * 15

