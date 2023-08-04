package ports

import (
	"context"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

type Producer interface {
	Write(context.Context, domain.Ship) error
}
