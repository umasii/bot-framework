package cardidentification

import (
	"errors"
	"strings"

	ccID "github.com/durango/go-credit-card"
	"github.com/cicadaaio/LVBot/Internal/Errors"
	"github.com/cicadaaio/LVBot/Internal/Profiles"
)

func CreditCardType(profile *Profiles.Profile) (string, error) {
	card := ccID.Card{Number: strings.ReplaceAll(profile.Billing.CardNumber, " ", ""), Cvv: profile.Billing.CVC, Month: profile.Billing.ExpMonth, Year: profile.Billing.ExpYear}
	err := card.Method()

	if err != nil {
		return "", Errors.Handler(errors.New("Failed to ID card"))
	}

	return card.Company.Short, nil

}
