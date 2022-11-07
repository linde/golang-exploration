package cmd

import (
	"log"
	gc "myapp/greeterclient"
	pb "myapp/helloservice"

	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "minimal grpc client for greeter service",
	Run:   doClientRun,
}

var name, host string
var times, rest int64
var timeoutSecs int

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&host, "host", "s", "localhost", "server host")
	clientCmd.Flags().StringVarP(&name, "name", "n", "world", "whom to greet")
	clientCmd.Flags().Int64VarP(&times, "times", "t", 1, "times to greet them")
	clientCmd.Flags().Int64VarP(&rest, "rest", "r", 0, "seconds to sleep before serving")
	clientCmd.Flags().IntVarP(&timeoutSecs, "timeout", "x", 60, "timeout (in seconds)")
}

func handleReply(r *pb.HelloReply, err error) {
	if err != nil {
		log.Fatalf("client.SayHello failed: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func doClientRun(cmd *cobra.Command, args []string) {

	port, _ := cmd.Flags().GetInt("port") // TODO should this be with the others?
	client := gc.New(host, port)

	client.Call(name, times, rest, timeoutSecs, handleReply)

}
