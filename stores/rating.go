package stores

import (
	"sync"
)

type InMemoryRatingStore struct {
	// There will be concurrent requests to write
	// a laptop score to memory, so a mutex is needed.
	m sync.RWMutex // multiple readers, one writer
	// key: laptop ID; value: Rating.
	rating map[string]*Rating
}

func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		rating: make(map[string]*Rating),
	}
}

func (st *InMemoryRatingStore) Add(laptopId string, score float64) (*Rating, error) {
	st.m.Lock() // locking for writing. Also reads are blocked.
	defer st.m.Unlock()

	rating, alreadyExists := st.rating[laptopId]

	if alreadyExists {
		// Update rating.
		rating.Count++
		rating.Sum += score
	} else {
		// First rating!
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
		st.rating[laptopId] = rating
	}

	return rating, nil
}
