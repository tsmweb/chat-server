package handler

import (
	"context"
	"fmt"
	"github.com/tsmweb/chat-service/chat"
	"github.com/tsmweb/chat-service/chat/message"
	"github.com/tsmweb/chat-service/pkg/concurrent"
	"github.com/tsmweb/chat-service/pkg/ebus"
	"github.com/tsmweb/chat-service/pkg/topic"
	"time"
)

// GroupMessage handles messages destined for a group of users.
type GroupMessage struct {
	repository chat.Repository
	executor   concurrent.ExecutorService
}

func (gm *GroupMessage) HandleGroupMessage() {
	bus := ebus.Instance()
	sub := bus.Subscribe(topic.GroupMessage)
	defer sub.Unsubscribe()

	for event := range sub.Event {
		gm.executor.Schedule(func(ctx context.Context) {
			msg := event.Data.(message.Message)

			members, err := gm.repository.GetGroupMembers(msg.Group)
			if err != nil {
				bus.Publish(topic.ErrorMessage,
					fmt.Errorf("%s; GroupMessage.HandleGroupMessage(): %v", msg.Group, err.Error()))
				return
			}

			for _, member := range members {
				msg.To = member
				bus.Publish(topic.Message, msg)
			}
		})
	}
}

// PresenceStatus handles user presence messages.
type PresenceStatus struct {
	repository chat.Repository
	executor   concurrent.ExecutorService
	host       string
}

func (ps *PresenceStatus) HandlePresenceMessage() {
	bus := ebus.Instance()
	sub := bus.Subscribe(topic.UserStatus)
	defer sub.Unsubscribe()

	for event := range sub.Event {
		ps.executor.Schedule(func(ctx context.Context) {
			user := event.Data.(chat.UserPresence)

			if err := ps.setStatus(user.ID, user.Status); err != nil {
				bus.Publish(topic.ErrorMessage,
					fmt.Sprintf("%s:%s; PresenceStatus.HandlePresenceMessage(): %v",
						user.ID, user.Status, err.Error()))
			}

			//TODO Notify Presence

			if user.Status == chat.UserOnline {
				if err := ps.sendOfflineMessage(user.ID); err != nil {
					bus.Publish(topic.ErrorMessage,
						fmt.Sprintf("%s:%s; PresenceStatus.HandlePresenceMessage(): %v",
							user.ID, user.Status, err.Error()))
				}
			}
		})
	}
}

func (ps *PresenceStatus) setStatus(userID string, userStatus chat.UserStatus) error {
	if userStatus == chat.UserOnline {
		if err := ps.repository.AddUserOnline(userID, ps.host, time.Now().UTC()); err != nil {
			return err
		}
	} else {
		if err := ps.repository.DeleteUserOnline(userID); err != nil {
			return err
		}
	}

	return nil
}

func (ps *PresenceStatus) sendOfflineMessage(userID string) error {
	messages, err := ps.repository.GetMessagesOffline(userID)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		ebus.Instance().Publish(topic.Message, msg)
	}

	return nil
}
