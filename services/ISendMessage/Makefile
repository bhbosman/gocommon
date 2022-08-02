all: DeleteAllGeneratedFiles  ISendMessage

DeleteAllGeneratedFiles:
	del *_InterfaceMethods.go




ISendMessage:
	mockgen -package ISendMessage -generateWhat ddd -destination ISendMessage_InterfaceMethods.go . ISendMessage,IMultiSendMessage
