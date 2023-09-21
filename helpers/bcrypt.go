package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) string {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hashPassword)
}

func ComparePassword(pw, hashedPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw))
	return err == nil
}
