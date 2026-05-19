package retailcrm

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTag_MarshalJSON(t *testing.T) {
	tags := []Tag{
		{"first", "#3e89b6", false},
		{"second", "#ffa654", false},
	}
	names := []byte(`["first","second"]`)
	str, err := json.Marshal(tags)

	if err != nil {
		t.Errorf("%v", err.Error())
	}

	if !reflect.DeepEqual(str, names) {
		t.Errorf("Marshaled: %#v\nExpected: %#v\n", str, names)
	}
}

func TestStringOrNumber_UnmarshalJSON(t *testing.T) {
	var value StringOrNumber

	require.NoError(t, json.Unmarshal([]byte(`"10"`), &value))
	assert.Equal(t, StringOrNumber("10"), value)

	require.NoError(t, json.Unmarshal([]byte(`20`), &value))
	assert.Equal(t, StringOrNumber("20"), value)

	require.NoError(t, json.Unmarshal([]byte(`null`), &value))
	assert.Empty(t, value)

	require.Error(t, json.Unmarshal([]byte(`true`), &value))
}

func TestStringOrNumber_MarshalJSON(t *testing.T) {
	store := Store{Code: "main", Ordering: "10"}

	data, err := json.Marshal(store)

	require.NoError(t, err)
	assert.JSONEq(t, `{"code":"main","ordering":"10"}`, string(data))
}

func TestResponseInfo_UnmarshalJSON(t *testing.T) {
	var fromObject ResponseInfo
	require.NoError(t, json.Unmarshal([]byte(`{
		"deliveryType": {
			"id": 38,
			"code": "sdek-v-2-podkliuchenie-1"
		},
		"billingInfo": {
			"price": 0,
			"currency": {
				"name": "Рубль",
				"shortName": "руб.",
				"code": "RUB"
			},
			"billingType": "fixed"
		}
	}`), &fromObject))

	assert.Equal(t, 38, fromObject.DeliveryTypeInfo.ID)
	assert.Equal(t, "sdek-v-2-podkliuchenie-1", fromObject.DeliveryTypeInfo.Code)
	require.NotNil(t, fromObject.BillingInfo)
	require.NotNil(t, fromObject.BillingInfo.Currency)
	assert.Equal(t, "RUB", fromObject.BillingInfo.Currency.Code)

	var fromArray ResponseInfo
	require.NoError(t, json.Unmarshal([]byte(`[]`), &fromArray))
	assert.Equal(t, ResponseInfo{}, fromArray)

	fromArray = ResponseInfo{
		BillingInfo: &BillingInfo{BillingType: "fixed"},
	}
	require.NoError(t, json.Unmarshal([]byte(`null`), &fromArray))
	assert.Equal(t, ResponseInfo{}, fromArray)
}

func TestIntegrationModuleEditResponse_UnmarshalJSON(t *testing.T) {
	var fromObject IntegrationModuleEditResponse
	require.NoError(t, json.Unmarshal([]byte(`{
		"success": true,
		"info": {
			"billingInfo": {
				"billingType": "fixed"
			}
		}
	}`), &fromObject))

	assert.True(t, fromObject.Success)
	assert.Equal(t, "fixed", fromObject.Info.BillingInfo.BillingType)
	assert.JSONEq(t, `{"billingInfo":{"billingType":"fixed"}}`, string(fromObject.InfoRaw))
	assert.Empty(t, fromObject.InfoList)

	var fromEmptyArray IntegrationModuleEditResponse
	require.NoError(t, json.Unmarshal([]byte(`{"success": true, "info": []}`), &fromEmptyArray))

	assert.True(t, fromEmptyArray.Success)
	assert.Equal(t, ResponseInfo{}, fromEmptyArray.Info)
	assert.JSONEq(t, `[]`, string(fromEmptyArray.InfoRaw))
	assert.Empty(t, fromEmptyArray.InfoList)

	var fromArray IntegrationModuleEditResponse
	require.NoError(t, json.Unmarshal([]byte(`{
		"success": true,
		"info": [
			{
				"deliveryType": {
					"id": 38,
					"code": "sdek-v-2-podkliuchenie-1"
				}
			}
		]
	}`), &fromArray))

	require.Len(t, fromArray.InfoList, 1)
	assert.Equal(t, 38, fromArray.InfoList[0].DeliveryTypeInfo.ID)
	assert.Equal(t, fromArray.InfoList[0], fromArray.Info)
}

func TestIntegrationModule_MarshalNewAndLegacyFields(t *testing.T) {
	active := true

	module := IntegrationModule{
		Code:               "module-code",
		IntegrationCode:    "integration-code",
		Active:             &active,
		AvailableCountries: []string{"RU"},
		Integrations: &Integrations{
			Delivery: &Delivery{
				DeliveryDataFieldList: []DeliveryDataField{
					{
						Code: "terminal",
						Type: "choice",
						ChoiceList: []DeliveryDataFieldChoice{
							{
								Value: "terminal-1",
								Label: "Terminal 1",
							},
						},
					},
				},
			},
			Payment: &PaymentModule{
				Actions:      PaymentModuleActions{Create: "/create"},
				InvoiceTypes: []string{"link"},
				Shops: []PaymentModuleShop{
					{
						Code:   "main",
						Name:   "Main shop",
						Active: true,
					},
				},
			},
			EmbedJS: &EmbedJS{
				Pages: []EmbedJSPage{
					{
						Code:          "settings",
						MenuItemTitle: map[string]string{"ru": "Настройки"},
						PageHelpLinks: ConfigurationTranslation{"ru": "https://example.com/help"},
					},
				},
			},
		},
	}

	data, err := json.Marshal(module)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"code": "module-code",
		"integrationCode": "integration-code",
		"active": true,
		"availableCountries": ["RU"],
		"integrations": {
			"delivery": {
				"deliveryDataFieldList": [
					{
						"code": "terminal",
						"type": "choice",
						"choices": [
							{
								"value": "terminal-1",
								"label": "Terminal 1"
							}
						]
					}
				]
			},
			"payment": {
				"actions": {
					"create": "/create"
				},
				"invoiceTypes": ["link"],
				"shops": [
					{
						"code": "main",
						"name": "Main shop",
						"active": true
					}
				]
			},
			"embedJs": {
				"pages": [
					{
						"code": "settings",
						"menuItemTitle": {
							"ru": "Настройки"
						},
						"pageHelpLink": {
							"ru": "https://example.com/help"
						}
					}
				]
			}
		}
	}`, string(data))

	var decoded IntegrationModule
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.NotNil(t, decoded.Integrations)
	require.NotNil(t, decoded.Integrations.Delivery)
	require.NotNil(t, decoded.Integrations.EmbedJS)

	deliveryField := decoded.Integrations.Delivery.DeliveryDataFieldList[0]
	assert.Equal(t, []string{"terminal-1"}, deliveryField.Choices)
	require.Len(t, deliveryField.ChoiceList, 1)
	assert.Equal(t, "Terminal 1", deliveryField.ChoiceList[0].Label)

	page := decoded.Integrations.EmbedJS.Pages[0]
	assert.Equal(t, "https://example.com/help", page.PageHelpLink)
	assert.Equal(t, ConfigurationTranslation{"ru": "https://example.com/help"}, page.PageHelpLinks)

	legacy := IntegrationModule{
		Integrations: &Integrations{
			Telephony: &Telephony{
				AdditionalCodes: []AdditionalCode{
					{
						UserID: "10",
						Code:   "101",
					},
				},
			},
			Delivery: &Delivery{
				DeliveryDataFieldList: []DeliveryDataField{
					{
						Code:    "terminal",
						Choices: []string{"terminal-1"},
					},
				},
			},
			EmbedJS: &EmbedJS{
				Pages: []EmbedJSPage{
					{
						Code:          "settings",
						MenuItemTitle: map[string]string{"ru": "Настройки"},
						PageHelpLink:  "https://example.com/help",
					},
				},
			},
		},
	}

	data, err = json.Marshal(legacy)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"integrations": {
			"telephony": {
				"additionalCodes": [
					{
						"userId": "10",
						"code": "101"
					}
				]
			},
			"delivery": {
				"deliveryDataFieldList": [
					{
						"code": "terminal",
						"choices": ["terminal-1"]
					}
				]
			},
			"embedJs": {
				"pages": [
					{
						"code": "settings",
						"menuItemTitle": {
							"ru": "Настройки"
						},
						"pageHelpLink": "https://example.com/help"
					}
				]
			}
		}
	}`, string(data))
}

func TestAdditionalCode_UnmarshalUserIDNumber(t *testing.T) {
	var additionalCode AdditionalCode

	require.NoError(t, json.Unmarshal([]byte(`{"userId":10,"code":"101"}`), &additionalCode))
	assert.Equal(t, "10", additionalCode.UserID)
	assert.Equal(t, "101", additionalCode.Code)
}

func TestStoreOrdering_UnmarshalStringOrNumber(t *testing.T) {
	var fromString Store
	require.NoError(t, json.Unmarshal([]byte(`{"ordering":"10"}`), &fromString))
	assert.Equal(t, StringOrNumber("10"), fromString.Ordering)

	var fromNumber Store
	require.NoError(t, json.Unmarshal([]byte(`{"ordering":20}`), &fromNumber))
	assert.Equal(t, StringOrNumber("20"), fromNumber.Ordering)
}

func TestAPIErrorsList_UnmarshalJSON(t *testing.T) {
	var list APIErrorsList

	require.NoError(t, json.Unmarshal([]byte(`["first", "second"]`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["0"], "first")
	assert.Equal(t, list["1"], "second")

	require.NoError(t, json.Unmarshal([]byte(`{"a": "first", "b": "second"}`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["a"], "first")
	assert.Equal(t, list["b"], "second")

	require.NoError(t, json.Unmarshal([]byte(`[]`), &list))
	assert.Len(t, list, 0)
}

func TestCustomFieldsList_UnmarshalJSON(t *testing.T) {
	var list CustomFieldMap

	require.NoError(t, json.Unmarshal([]byte(`["first", "second"]`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["0"], "first")
	assert.Equal(t, list["1"], "second")

	require.NoError(t, json.Unmarshal([]byte(`{"a": "first", "b": "second"}`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["a"], "first")
	assert.Equal(t, list["b"], "second")

	require.NoError(t, json.Unmarshal([]byte(`{"a": ["first", "second"], "b": "second"}`), &list))
	assert.Len(t, list, 2)
	assert.Len(t, list["a"].([]interface{}), 2)
	assert.Equal(t, list["a"].([]interface{})[0], "first")
	assert.Equal(t, list["a"].([]interface{})[1], "second")
	assert.Equal(t, list["b"], "second")

	require.NoError(t, json.Unmarshal([]byte(`[]`), &list))
	assert.Len(t, list, 0)
}

func TestOrderPayments_UnmarshalJSON(t *testing.T) {
	var list OrderPayments

	require.NoError(t, json.Unmarshal([]byte(`[{"id": 1}, {"id": 2}]`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["0"], OrderPayment{ID: 1})
	assert.Equal(t, list["1"], OrderPayment{ID: 2})

	require.NoError(t, json.Unmarshal([]byte(`{"a": {"id": 1}, "b": {"id": 2}}`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["a"], OrderPayment{ID: 1})
	assert.Equal(t, list["b"], OrderPayment{ID: 2})

	require.NoError(t, json.Unmarshal([]byte(`[]`), &list))
	assert.Len(t, list, 0)
}

func TestProperties_UnmarshalJSON(t *testing.T) {
	var list Properties

	require.NoError(t, json.Unmarshal([]byte(`[{"code": "first"}, {"code": "second"}]`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["0"], Property{Code: "first"})
	assert.Equal(t, list["1"], Property{Code: "second"})

	require.NoError(t, json.Unmarshal([]byte(`{"a": {"code": "first"}, "b": {"code": "second"}}`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["a"], Property{Code: "first"})
	assert.Equal(t, list["b"], Property{Code: "second"})

	require.NoError(t, json.Unmarshal([]byte(`[]`), &list))
	assert.Len(t, list, 0)
}
