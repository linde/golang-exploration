package e2e

import (
	"io"

	pb "myapp/greeter"
)

func ReplyStreamToBuffer(replyStream pb.Greeter_SayHelloClient) ([]*pb.HelloReply, error) {

	var replies []*pb.HelloReply

	for {
		r, err := replyStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return replies, err
		}
		replies = append(replies, r)
	}

	return replies, nil
}
