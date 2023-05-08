package transport

import "math/rand"

var Sessions = map[string]*Session{}

type Session struct {
	Login string
}

func NewSession(login string) (newSessionID string) {

	var newSession = &Session{Login: login}

	newSessionID = generateNewSessionID()
	Sessions[newSessionID] = newSession

	return newSessionID
}

func generateNewSessionID() string {
	var letters = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	var lenLetters = len(letters)
	var newID = make([]rune, 16)
	for i, _ := range newID {
		newID[i] = letters[rand.Intn(lenLetters)]
	}
	return string(newID)
}