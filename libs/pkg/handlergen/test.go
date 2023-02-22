package handlergen

import (
	"context"

	"github.com/LSDXXX/libs/pkg/handlergen/helper"
)

type serviceTest struct{}

func (s *serviceTest) WithContext(ctx context.Context) *serviceTest {
	return s
}

//@RequestMapping(/test/haha)
//@GenerateType(server)
type TestHandler interface {
	//@RequestParam(haah=@id)
	//@RequestMapping(/user, GET)
	TestMethod(id *int) (int, error)

	//@RequestMapping(/user1, GET)
	TestMethod2(id *int) (int, error)

	//@RequestMapping(/user2, GET)
	TestNoRes(id *int) error

	helper.InjectServices1[*serviceTest]
}
