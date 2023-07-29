package ports

import "github.com/mikeewhite/ship-locator/backend/internal/core/domain"

type Publisher interface {
	Publish(data *domain.Ship) error
}

type Consumer interface {
	Consume()
}
