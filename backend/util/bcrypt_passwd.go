package util

import "golang.org/x/crypto/bcrypt"

func GenPasswd(passwd string) string {
	hashPasswd, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(hashPasswd)
}

func ComparePasswd(hashPasswd string, passwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPasswd), []byte(passwd)); err != nil {
		return err
	}

	return nil
}
