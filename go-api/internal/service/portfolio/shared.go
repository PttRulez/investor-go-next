package portfolio

import (
	"context"

	"github.com/pttrulez/investor-go/internal/domain"
)

func (s *Service) checkSecurity(ctx context.Context, e domain.Exchange,
	st domain.SecurityType, t string) (decimalCount int, err error) {
	if e == domain.EXCHMoex && st == domain.STShare {
		var share domain.Share
		share, err = s.moexService.GetShareByTicker(ctx, t)
		decimalCount = share.PriceDecimals
	} else if e == domain.EXCHMoex && st == domain.STBond {
		_, err = s.moexService.GetBondByTicker(ctx, t)
		decimalCount = 2
	}

	return
}
