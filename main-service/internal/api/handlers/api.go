package handlers

type APIHandler struct {
	*UserHandler
	*ContentHandler
}

func NewAPIHandler(userHandler *UserHandler, contentHandler *ContentHandler) *APIHandler {
	return &APIHandler{
		UserHandler:    userHandler,
		ContentHandler: contentHandler,
	}
}
