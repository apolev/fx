package fx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestDecorate(t *testing.T) {
	type expectation [2]interface{}
	tests := []struct {
		name string
		test func() (*fxtest.App, []expectation)
	}{
		{
			name: "simple",
			test: func() (*fxtest.App, []expectation) {
				var populate string
				return fxtest.New(t,
						fx.Provide(func() string { return "thing" }),
						fx.Decorate(func(s string) string { return "decorated: " + s }),
						fx.Populate(&populate),
					), []expectation{
						{"decorated: thing", populate},
					}
			},
		},
		{
			name: "decorate no decorate parent module",
			test: func() (*fxtest.App, []expectation) {
				var inner, outer string
				return fxtest.New(t,
						fx.Module("foo",
							fx.Provide(func() string { return "thing" }),
							fx.Decorate(func(s string) string { return "decorated: " + s }),
							fx.Populate(&inner),
						),
						fx.Populate(&outer),
					), []expectation{
						{"decorated: thing", inner},
						{"thing", outer},
					}
			},
		},
		{
			name: "decorate decorates child module",
			test: func() (*fxtest.App, []expectation) {
				var inner, outer string
				return fxtest.New(t,
						fx.Module("foo",
							fx.Provide(func() string { return "thing" }),
							fx.Populate(&inner),
						),
						fx.Decorate(func(s string) string { return "decorated: " + s }),
						fx.Populate(&outer),
					), []expectation{
						{"decorated: thing", inner},
						{"decorated: thing", outer},
					}
			},
		},
		{
			name: "complex decorate",
			test: func() (*fxtest.App, []expectation) {
				var inner, middle, outer string
				return fxtest.New(t,
						fx.Module("foo",
							fx.Module("bar",
								fx.Provide(func() string { return "thing" }),
								fx.Populate(&inner),
							),
							fx.Decorate(func(s string) string { return "decorated: " + s }),
							fx.Populate(&middle),
						),
						fx.Populate(&outer),
					), []expectation{
						{"decorated: thing", inner},
						{"decorated: thing", middle},
						{"thing", outer},
					}
			},
		},
		{
			name: "multiple decorate decorates",
			test: func() (*fxtest.App, []expectation) {
				var inner, middle, outer string
				return fxtest.New(t,
						fx.Module("foo",
							fx.Module("bar",
								fx.Provide(func() string { return "thing" }),
								fx.Decorate(func(s string) string { return "decorate in: " + s }),
								fx.Populate(&inner),
							),
							fx.Decorate(func(s string) string { return "decorate out: " + s }),
							fx.Populate(&middle),
						),
						fx.Populate(&outer),
					), []expectation{
						{"decorated: thing", inner},
						{"decorate out: decorate in: thing", middle},
						{"thing", outer},
					}
			},
		},
		{
			name: "multiple decorate same module",
			test: func() (*fxtest.App, []expectation) {
				var out string
				return fxtest.New(t,
						fx.Provide(func() string { return "thing" }),
						fx.Decorate(func(s string) string { return "decorate one: " + s }),
						fx.Decorate(func(s string) string { return "decorate two: " + s }),
						fx.Populate(&out),
					), []expectation{
						{"decorate two: decorate one: thing", out},
					}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, expectations := tt.test()
			app.RequireStart().RequireStop()
			for _, e := range expectations {
				assert.EqualValues(t, e[0], e[1])
			}
		})
	}
}