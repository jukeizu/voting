package startup

import "google.golang.org/grpc"

type Startup interface {
	Migrate() error
	RegisterServer(grpcServer *grpc.Server)
}
