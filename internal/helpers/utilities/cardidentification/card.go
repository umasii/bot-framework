package cardidentification

import (
	goErrors "errors"
	"strings"

	ccID "github.com/durango/go-credit-card"
	"github.com/umasii/bot-framework/internal/errors"
	profiles "github.com/umasii/bot-framework/internal/profiles"
)

func CreditCardType(profile *profiles.Profile) (string, error) {
	card := ccID.Card{Number: strings.ReplaceAll(profile.Billing.CardNumber, " ", ""), Cvv: profile.Billing.CVC, Month: profile.Billing.ExpMonth, Year: profile.Billing.ExpYear}
	err := card.Method()

	if err != nil {
		return "", errors.Handler(goErrors.New("Failed to ID card"))
	}

	return card.Company.Short, nil

}
