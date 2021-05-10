package core

type MessageHandler struct {
	repository               Repository
	offlineMessageDispatcher *OfflineMessageDispatcher
	groupMessageDispatcher   *GroupMessageDispatcher
}

func NewMessageHandler(
	r Repository,
	omd *OfflineMessageDispatcher,
	gmd *GroupMessageDispatcher,
) *MessageHandler {
	return &MessageHandler{
		repository:               r,
		offlineMessageDispatcher: omd,
		groupMessageDispatcher:   gmd,
	}
}

func (mh *MessageHandler) HandleMessage(msg *Message) error {
	if msg.IsGroupMessage() {
		return mh.groupMessageDispatcher.Send(msg)
	}

	host, ok, err := mh.repository.GetUserOnline(msg.To)
	if err != nil {
		return err
	}
	if !ok {
		return mh.offlineMessageDispatcher.Send(msg)
	}

	return mh.sendMessageBygRPC(host, msg)
}

func (mh *MessageHandler) SendMessageOffline(userID string, chMessage chan<- *Message) error {
	messages, err := mh.repository.GetMessagesOffline(userID)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		chMessage <- msg
	}

	return nil
}

func (mh *MessageHandler) sendMessageBygRPC(host string, msg *Message) error {
	// TODO
	return nil
}
