package namespace

import (
	"log"

	"api/services/socketio/server/contract"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func DigitalSignNamespace(socketIO contract.Server) error {
	np := socketIO.AddNamespace("digital-sign")

	return np.OnConnect(func(s *socket.Socket) {
		np.AddConnection(s)

		// events
		{
			if err := s.On("message", func(args ...interface{}) {
				//requestID := msg.(string)
				//s.Emit("prepare-response", requestID)
				msg := args[len(args)-1].(string)

				if err := s.Emit("test", msg+"123"); err != nil {
					log.Println("error while sending message:", err)
				}
			}); err != nil {
				log.Println("error while binding event:", err)
			}

			if err := s.On("ping", func(_ ...interface{}) {
				if err := s.Emit("pong"); err != nil {
					log.Println("error while sending pong:", err)
				}
			}); err != nil {
				log.Println("error while binding event ping:", err)
			}
		}
	})

	//np.OnDisconnect(func(s *socket.Socket, reason string) {
	//	np.DeleteConnection(s.ID())
	//
	//	err := s.Close()
	//	if err != nil {
	//		log.Println("error while closing socket connection:", err)
	//	}
	//	log.Println("socket closed", reason)
	//})
	//
	//np.OnError(func(s *socket.Socket, e error) {
	//	log.Println("meet socket error:", e)
	//})
}
