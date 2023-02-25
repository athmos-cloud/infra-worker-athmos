package plugin

import (
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestFactoryBuilder(t *testing.T) {
	builder := FactoryBuilder("terraform")
	assert.Equal(t, reflect.TypeOf(builder), reflect.TypeOf(&terraform.Builder{}))
}
