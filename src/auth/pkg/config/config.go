package config

// JwtCookieName is the env variable name that is also used to set the JWT
// cookie for sessions
const JwtCookieName = "JWT_SECRET"

// TokenExpiration is the time in seconds when a token will expire
const TokenExpiration = 3600 * 24 * 7


// MinEntropyBits is the minimum number of bits of entropy required for a password.
const MinEntropyBits = 64

// MaxLoginAttempts is the maximum number of times user can attempt to enter the correct password
// before their account is temporarily locked
const MaxLoginAttempts = 5

// AccountLockoutLength is the time in seconds that an account will be locked
const AccountLockoutLength = 60 * 15

