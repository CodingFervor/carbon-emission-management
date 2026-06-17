package password

import "golang.org/x/crypto/bcrypt"

// Hash returns a bcrypt hash of the plain-text password.
func Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare reports whether the plain-text password matches the stored hash.
func Compare(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
