package session

import (
	"fmt"
	"sync"

	"github.com/artshah-coder/go-shopql/internal/utils/randutils"
)

var (
	ErrNoSession  = fmt.Errorf("no session for this token")
	ErrNoSessions = fmt.Errorf("no sessions for this uid")
)

type SessionsStorage interface {
	Add(*Session) (Token, error)
	Get(Token) (*Session, error)
	DeleteByToken(Token) error
	DeleteAllSessionsByUID(uint32) error
}

type SessionStMem struct {
	mu       *sync.Mutex
	Sessions map[Token]*Session
}

func NewSessionStMem() *SessionStMem {
	return &SessionStMem{
		mu:       &sync.Mutex{},
		Sessions: make(map[Token]*Session),
	}
}

func (sessSt *SessionStMem) Add(session *Session) (Token, error) {
	token := Token(randutils.RandString(32))
	sessSt.mu.Lock()
	sessSt.Sessions[token] = session
	sessSt.mu.Unlock()
	return token, nil
}

func (sessSt *SessionStMem) Get(token Token) (*Session, error) {
	sessSt.mu.Lock()
	session, ok := sessSt.Sessions[token]
	sessSt.mu.Unlock()
	if ok {
		return session, nil
	}
	return nil, ErrNoSession
}

func (sessSt *SessionStMem) DeleteByToken(token Token) error {
	sessSt.mu.Lock()
	_, ok := sessSt.Sessions[token]
	sessSt.mu.Unlock()
	if ok {
		delete(sessSt.Sessions, token)
		return nil
	}
	return ErrNoSession
}

func (sessSt *SessionStMem) DeleteAllSessionsByUID(uid uint32) error {
	count := 0
	sessSt.mu.Lock()
	for token, session := range sessSt.Sessions {
		if session.UserID == uid {
			delete(sessSt.Sessions, token)
			count++
		}
	}
	sessSt.mu.Unlock()
	if count == 0 {
		return ErrNoSessions
	}

	return nil
}
