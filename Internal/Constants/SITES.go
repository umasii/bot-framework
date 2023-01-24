package Constants

import (
	"reflect"

	"github.com/cicadaaio/LVBot/Sites/louisvuitton"
)

var SITES = map[string]reflect.Type{
	"Louis Vuitton": reflect.TypeOf(louisvuitton.LVTask{}),
}