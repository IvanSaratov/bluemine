package handlers

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

func loginUser(userLogin, userPassword string) (sessionId string, err error) {
	var password string

	if passwordHash(userPassword) != password {
		err = errors.New("Incorrect password")
		return
	}

	return
}

func passwordHash(password string) string {
	sh := sha1.New()
	io.WriteString(sh, password)

	md := md5.New()
	io.WriteString(md, password)

	return fmt.Sprintf("%x:%x", sh.Sum(nil), md.Sum(nil))
}
