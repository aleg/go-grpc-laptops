package stores

import (
	"sync"

	"github.com/aleg/go-grpc-laptops/users"
)

type InMemoryUserStore struct {
	// There will be concurrent requests to write
	// a laptop to memory, so a mutex is needed.
	m sync.RWMutex // multiple readers, one writer
	// key: username; value: User object.
	users map[string]*users.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*users.User),
	}
}

func (st *InMemoryUserStore) Save(user *users.User) error {
	st.m.Lock()
	defer st.m.Unlock()

	if _, alreadyExist := st.users[user.Username]; alreadyExist {
		return ErrorAlreadyExists
	}

	st.users[user.Username] = user.Clone()
	return nil
}

func (st *InMemoryUserStore) Find(username string) (*users.User, error) {
	st.m.RLock()
	defer st.m.RUnlock()

	user, found := st.users[username]
	if !found {
		return nil, nil
	}

	return user.Clone(), nil
}
