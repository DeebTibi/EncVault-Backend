package ClientDirect

import "google.golang.org/grpc"

type ClientDirect[T any] struct {
	Address      string
	CreateClient func(cc grpc.ClientConnInterface) T
}

func (obj *ClientDirect[T]) Connect() (disconnect func(), client T, err error) {
	conn, err := grpc.NewClient(obj.Address)
	if err != nil {
		var client T
		return nil, client, err
	}
	client = obj.CreateClient(conn)
	return func() { conn.Close() }, client, nil
}
