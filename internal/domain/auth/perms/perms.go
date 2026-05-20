package perms

import "github.com/alan-b-lima/almodon/internal/domain/auth"

var (
	None = auth.Allow()
	All  = auth.Allow(auth.Unlogged)
)

var (
	StockRead   = auth.Allow(auth.User)
	StockMgmt   = auth.Allow(auth.Admin)
	StockChange = auth.Allow(auth.Promoted)
)

var (
	UserMgmt = auth.Allow(auth.Chief)
	UserSelf = auth.Allow(auth.User)
)
