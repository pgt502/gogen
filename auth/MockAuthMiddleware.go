package auth

import (
	context "context"
	http "net/http"

	grpc "google.golang.org/grpc"
)

type MockAuthMiddleware struct {
	CopyCall struct {
		Receives struct {
			Param0 string
			Param1 []interface{}
		}
		Returns struct {
			Ret0 int
			Ret1 []interface{}
		}
		GetsCalled struct {
			Times int
		}
	}

	GetIdentityCall struct {
		Receives struct {
			Param0 http.Handler
		}
		Returns struct {
			Ret0 http.Handler
		}
		GetsCalled struct {
			Times int
		}
	}

	GrpcGetIdentityCall struct {
		Receives struct {
			Param0 context.Context
			Param1 interface{}
			Param2 grpc.UnaryServerInfo
			Param3 grpc.UnaryHandler
		}
		Returns struct {
			Ret0 interface{}
			Ret1 error
		}
		GetsCalled struct {
			Times int
		}
	}

	LogCall struct {
		Receives struct {
			Param0 http.Handler
		}
		Returns struct {
			Ret0 http.Handler
		}
		GetsCalled struct {
			Times int
		}
	}
}

func (m *MockAuthMiddleware) Copy(p0 string, p1 ...interface{}) (r0 int, r1 []interface{}) {
	m.CopyCall.GetsCalled.Times++
	m.CopyCall.Receives.Param0 = p0
	m.CopyCall.Receives.Param1 = p1

	r0 = m.CopyCall.Returns.Ret0
	r1 = m.CopyCall.Returns.Ret1

	return
}

func (m *MockAuthMiddleware) GetIdentity(p0 http.Handler) (r0 http.Handler) {
	m.GetIdentityCall.GetsCalled.Times++
	m.GetIdentityCall.Receives.Param0 = p0

	r0 = m.GetIdentityCall.Returns.Ret0

	return
}

func (m *MockAuthMiddleware) GrpcGetIdentity(p0 context.Context, p1 interface{}, p2 grpc.UnaryServerInfo, p3 grpc.UnaryHandler) (r0 interface{}, r1 error) {
	m.GrpcGetIdentityCall.GetsCalled.Times++
	m.GrpcGetIdentityCall.Receives.Param0 = p0
	m.GrpcGetIdentityCall.Receives.Param1 = p1
	m.GrpcGetIdentityCall.Receives.Param2 = p2
	m.GrpcGetIdentityCall.Receives.Param3 = p3

	r0 = m.GrpcGetIdentityCall.Returns.Ret0
	r1 = m.GrpcGetIdentityCall.Returns.Ret1

	return
}

func (m *MockAuthMiddleware) Log(p0 http.Handler) (r0 http.Handler) {
	m.LogCall.GetsCalled.Times++
	m.LogCall.Receives.Param0 = p0

	r0 = m.LogCall.Returns.Ret0

	return
}
