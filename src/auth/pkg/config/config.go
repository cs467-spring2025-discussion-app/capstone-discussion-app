package config

// SessionCookieName is the env variable name used to set the cookie for sessions
const SessionKey = "SESSION_SECRET"

// SessionCookieName is the env variable name used to set the cookie for sessions
const SessionCookieName = "SESSION_COOKIE_NAME"

// SessionExpiration is the time in seconds when a token will expire
const SessionExpiration = 3600 * 24 * 7


// MinEntropyBits is the minimum number of bits of entropy required for a password.
const MinEntropyBits = 64

// MaxLoginAttempts is the maximum number of times user can attempt to enter the correct password
// before their account is temporarily locked
const MaxLoginAttempts = 5

// AccountLockoutLength is the time in seconds that an account will be locked
const AccountLockoutLength = 60 * 15

