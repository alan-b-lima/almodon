package inventory

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Material struct {
	uuid                  uuid.UUID
	name                  string
	siads                 int
	catmat                int
	ecampus               int
	description           string
	measurementUnit       string
	currentTotalQuantity  float64
	minAcceptableQuantity float64
	createdAt             time.Time
	updatedAt             time.Time 
}

type ItemUnit struct {
	uuid           uuid.UUID
	batch          uuid.UUID
	material       uuid.UUID
	quantity       float64
	expirationDate time.Time 
	createdAt      time.Time  
}

type ItemBatch struct {
	uuid           uuid.UUID  
	material       uuid.UUID  
	batchNumber    string     
	totalQuantity  float64    
	invoiceNumber  string     
	unitPrice      float64    
	isActive       bool       
	expirationDate time.Time 
	purchaseDate   time.Time  
	createdAt      time.Time  
}


