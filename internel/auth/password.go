package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword хеширует пароль перед сохранением в БД
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword сверяет введённый пароль с сохранённым хешем
func CheckPassword(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false

	}
	return true

}
