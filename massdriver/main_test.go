package massdriver

import (
	mocks "xo/utils/mocks"
)

func init() {
	Client = &mocks.MockClient{}
}
