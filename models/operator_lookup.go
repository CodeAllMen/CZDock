/**
  create by yy on 2020/5/26
*/

package models

type OperatorLookup struct {
	Result OperatorLookupResult `xml:"result" json:"result"`
}

type OperatorLookupActionResult struct {
	Code   int    `xml:"code" json:"code"`
	Detail string `xml:"detail" json:"detail"`
	Url    string `xml:"url" json:"url"`
	Status int    `xml:"status" json:"status"`
}

type OperatorLookupResult struct {
	ActionResult OperatorLookupActionResult `xml:"action_result" json:"action_result"`
	Reference    string                     `xml:"reference" json:"reference"`
	RequestId    string                     `xml:"request_id" json:"request_id"`
}
