package dmesg

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

// Stream creates a goroutine that, every sampleTime ticks, will send
// new dmesg messages to out.  It also listens on stop, in case you
// need to abort the goroutine.  If there is an error setting up the
// initial state, it is returned, but otherwise errors are logged and
// otherwise ignored.
func Stream(out chan<- *Message, stop <-chan bool, sampleTime time.Duration) error {
	state, err := New()
	if err != nil {
		return err
	}
	go doStream(state, out, stop, sampleTime)
	return nil
}

func doStream(state *State, out chan<- *Message, stop <-chan bool, sampleTime time.Duration) {
	var lastMessage *Message
	ticker := time.NewTicker(sampleTime)
	defer ticker.Stop()
	for _ = range ticker.C {
		select {
		case <-stop:
			log.Debug("Terminating as requested")
			return
		default:
		}
		messages, err := state.Messages()
		if err == nil {
			if lastMessage == nil {
				for _, message := range messages {
					lastMessage = message
					out <- message
				}
			} else {
				hasSeenLast := false
				for _, message := range messages {
					if message.Timestamp.After(lastMessage.Timestamp) {
						log.Debug("missed some, resuming where available")
						hasSeenLast = true
					} else if message == lastMessage {
						hasSeenLast = true
					}
					if hasSeenLast {
						lastMessage = message
						out <- message
					}
				}
			}
		} else {
			log.WithError(err).Warning("Messages returned error; hoping it clears up")
		}
	}
}
