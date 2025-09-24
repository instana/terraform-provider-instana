package restapi

// This file contains extensions to the InstanaAPI interface

// LogAlertConfig implementation of InstanaAPI interface
func (api *baseInstanaAPI) LogAlertConfig() RestResource[*LogAlertConfig] {
	return NewCreatePOSTUpdatePOSTRestResource(LogAlertConfigResourcePath, NewCustomPayloadFieldsUnmarshallerAdapter(NewDefaultJSONUnmarshaller(&LogAlertConfig{})), api.client)
}

// Made with Bob
