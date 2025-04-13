package config

// MinEntropyBits is the minimum number of bits of entropy required for a password.
const MinEntropyBits = 64

// MaxLoginAttempts is the maximum number of times user can attempt to enter the correct password
// before their account is temporarily locked
const MaxLoginAttempts = 5

// AccountLockoutLength is the time in seconds that an account will be locked
const AccountLockoutLength = 60 * 15

