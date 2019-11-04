package components

import (
	"github.com/bluest-eel/common/components"
	// "github.com/bluest-eel/state/api"
	"github.com/bluest-eel/state/components/config"
	"github.com/bluest-eel/state/components/db"
	// "google.golang.org/grpc"
)

// Base component collection
type Base struct {
	Config *config.Config
	components.BaseLogger
}

// BaseApp ...
type BaseApp struct {
	components.BaseApp
}

// BaseCLI ...
type BaseCLI struct {
	components.BaseCLI
}

// BaseDB ...
type BaseDB struct {
	DB *db.DB
}

// // BaseStateClient ...
// type BaseKVClient struct {
// 	KVConn   *grpc.ClientConn
// 	KVClient api.StateServiceClient
// }

// Default component collection
type Default struct {
	Base
	components.BaseGRPC
	// db db.Database
}

// TestBase component that keeps stdout clean
type TestBase struct {
	Config *config.Config
	components.TestBase
}

// TestGRPC ...
type TestGRPC struct {
	Config *config.Config
	components.TestBaseGRPC
}
