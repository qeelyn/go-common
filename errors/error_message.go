package errors

import (
	"fmt"
	"strings"
)

// support error template for self defined output
type ErrorMessage struct {
	IsDebug                bool
	Templates              *map[string]ErrorTemplate
	ErrorDescriptionHandle func(error) ErrorDescription
}

type ErrorDescription struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Params map[string]interface{}

type ErrorTemplate struct {
	Code    int    `yaml:"code"`
	Message string `yaml:"message"`
	Debug   string `yaml:"debug"`
}

// getMessage returns the error message by replacing placeholders in the error template with the actual parameters.
func (e ErrorTemplate) getCode() int {
	return e.Code
}

// getMessage returns the error message by replacing placeholders in the error template with the actual parameters.
func (e ErrorTemplate) getMessage(params Params) string {
	return replacePlaceholders(e.Message, params)
}

// getDeveloperMessage returns the developer message by replacing placeholders in the error template with the actual parameters.
func (e ErrorTemplate) getDebugMessage(params Params) string {
	return replacePlaceholders(e.Debug, params)
}

func replacePlaceholders(message string, params Params) string {
	if len(message) == 0 {
		return ""
	}
	for key, value := range params {
		message = strings.Replace(message, "{"+key+"}", fmt.Sprint(value), -1)
	}
	return message
}

// LoadMessages reads a YAML file containing error templates.
// func LoadMessages(file string) error {
// 		bytes, err := ioutil.ReadFile(file)
//		if err != nil {
//			return err
//		}
//		templates = map[string]ErrorTemplate{}
//		yaml.Unmarshal(bytes, &templates)
//      NewErrorHandle(templates)
//	}
func NewErrorMessage(templates *map[string]ErrorTemplate) *ErrorMessage {
	return &ErrorMessage{
		Templates: templates,
	}
}

//const errorDescription map[string]string =
func (t *ErrorMessage) GetErrorDescription(err error) *ErrorDescription {
	if t == nil || t.Templates == nil {
		return nil
	}
	ins, ok := (*t.Templates)[err.Error()]
	if !ok {
		return &ErrorDescription{Message: err.Error()}
	}
	if t.IsDebug {
		msg := ins.getDebugMessage(Params{"error": err.Error()})
		return &ErrorDescription{Code: ins.getCode(), Message: msg}
	}
	msg := ins.getMessage(Params{"error": err.Error()})
	return &ErrorDescription{Code: ins.getCode(), Message: msg}
}
