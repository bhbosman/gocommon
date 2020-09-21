package multiBlock

import (
	"context"
	"github.com/bhbosman/gologging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"testing"
)

func TestWriterFactoryService(t *testing.T) {
	var sut IReaderWriterFactoryService
	app := fxtest.New(
		t,
		gologging.ProvideLogFactoryForTesting(t, nil),
		ProvideReaderWriterFactoryService(),
		fx.Populate(&sut),
	)
	if !assert.NoError(t, app.Err()) {
		return
	}
	if !assert.NoError(t, app.Start(context.TODO())) {
		return
	}
	defer func() {
		assert.NoError(t, app.Start(context.TODO()))
	}()
	controller := gomock.NewController(t)
	defer controller.Finish()
}
