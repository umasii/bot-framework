package constants

import (
	"reflect"

	"github.com/umasii/bot-framework/sites/louisvuitton"
)

var SITES = map[string]reflect.Type{
	"Louis Vuitton": reflect.TypeOf(louisvuitton.LVTask{}),
}