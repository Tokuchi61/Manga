package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
)

var (
	ErrNotFound = errors.New("shop_repository_not_found")
)

// Store defines shop persistence boundary.
type Store interface {
	UpsertProductDefinition(ctx context.Context, product entity.ProductDefinition) error
	GetProductDefinition(ctx context.Context, productID string) (entity.ProductDefinition, error)
	ListProductDefinitions(ctx context.Context, state string, limit int, offset int) ([]entity.ProductDefinition, error)

	UpsertOfferDefinition(ctx context.Context, offer entity.OfferDefinition) error
	GetOfferDefinition(ctx context.Context, offerID string) (entity.OfferDefinition, error)
	ListOfferDefinitions(ctx context.Context, productID string, visibility string, activeOnly bool, limit int, offset int) ([]entity.OfferDefinition, error)

	CreatePurchaseIntent(ctx context.Context, intent entity.PurchaseIntent) error
	GetPurchaseIntent(ctx context.Context, intentID string) (entity.PurchaseIntent, error)
	UpdatePurchaseIntent(ctx context.Context, intent entity.PurchaseIntent) error
	ListPurchaseIntentsByUser(ctx context.Context, userID string, limit int, offset int) ([]entity.PurchaseIntent, error)
	FindLatestPurchaseByUserProduct(ctx context.Context, userID string, productID string) (entity.PurchaseIntent, error)

	GetPurchaseDedup(ctx context.Context, dedupKey string) (entity.PurchaseIntent, error)
	PutPurchaseDedup(ctx context.Context, dedupKey string, intent entity.PurchaseIntent) error

	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error
}
