package message

type BodyStructure struct {
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

type Body struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func message(status bool, data interface{}, msg string) BodyStructure {
	return BodyStructure{Status: status, Data: data, Message: msg}
}

func StatusOkMessage(body interface{}, msg string) BodyStructure {
	return message(true, body, msg)
}
func StatusBadRequestMessage(data interface{}, msg string) BodyStructure {
	return message(false, data, msg)
}
func StatusInternalServerErrorMessage() BodyStructure {
	return message(false, nil, "processing your request encountered some error!")
}
