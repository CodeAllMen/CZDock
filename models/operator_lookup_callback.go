/**
  create by yy on 2020/5/26
*/

package models

type OperatorLookupCallback struct {
	Result OperatorLookupCallBackResult `xml:"result"`
}

type OperatorLookupCustomParameter struct {
	Key   string `xml:"key" json:"key"`
	Value string `xml:"value" json:"value"`
}

type OperatorLookupCustomParameters struct {
	CustomParameter OperatorLookupCustomParameter `xml:"custom_parameter" json:"custom_parameter"`
	Country         string                        `xml:"country" json:"country"`
	Id              int                           `xml:"id" json:"id"`
	Language        string                        `xml:"language" json:"language"`
	Msisdn          string                        `xml:"msisdn" json:"msisdn"`
	Operator        string                        `xml:"operator" json:"operator"`
}

type OperatorLookupCallBackResult struct {
	ActionResult     OperatorLookupActionResult     `xml:"action_result" json:"action_result"`
	CustomParameters OperatorLookupCustomParameters `xml:"custom_parameters" json:"custom_parameters"`
	Reference        string                         `xml:"reference" json:"reference"`
	RequestId        string                         `xml:"request_id" json:"request_id"`
}
