package retailcrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Name)
}

func (a *AdditionalCode) UnmarshalJSON(data []byte) error {
	var response struct {
		Code   string          `json:"code"`
		UserID json.RawMessage `json:"userId"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	a.Code = response.Code
	a.UserID = ""

	userID := bytes.TrimSpace(response.UserID)
	if len(userID) == 0 || bytes.Equal(userID, []byte("null")) {
		return nil
	}

	if userID[0] == '"' {
		return json.Unmarshal(userID, &a.UserID)
	}

	decoder := json.NewDecoder(bytes.NewReader(userID))
	decoder.UseNumber()

	var number json.Number
	if err := decoder.Decode(&number); err != nil {
		return err
	}

	if _, err := number.Int64(); err != nil {
		return err
	}

	a.UserID = number.String()

	return nil
}

func (f DeliveryDataField) MarshalJSON() ([]byte, error) {
	choices := any(nil)
	if len(f.ChoiceList) > 0 {
		choices = f.ChoiceList
	} else if len(f.Choices) > 0 {
		choices = f.Choices
	}

	return json.Marshal(struct {
		Code            string `json:"code,omitempty"`
		Label           string `json:"label,omitempty"`
		Hint            string `json:"hint,omitempty"`
		Type            string `json:"type,omitempty"`
		AutocompleteURL string `json:"autocompleteUrl,omitempty"`
		Multiple        bool   `json:"multiple,omitempty"`
		Choices         any    `json:"choices,omitempty"`
		Visible         bool   `json:"visible,omitempty"`
		Required        bool   `json:"required,omitempty"`
		AffectsCost     bool   `json:"affectsCost,omitempty"`
		Editable        bool   `json:"editable,omitempty"`
	}{
		Code:            f.Code,
		Label:           f.Label,
		Hint:            f.Hint,
		Type:            f.Type,
		AutocompleteURL: f.AutocompleteURL,
		Multiple:        f.Multiple,
		Choices:         choices,
		Visible:         f.Visible,
		Required:        f.Required,
		AffectsCost:     f.AffectsCost,
		Editable:        f.Editable,
	})
}

func (f *DeliveryDataField) UnmarshalJSON(data []byte) error {
	var response struct {
		Code            string          `json:"code"`
		Label           string          `json:"label"`
		Hint            string          `json:"hint"`
		Type            string          `json:"type"`
		AutocompleteURL string          `json:"autocompleteUrl"`
		Multiple        bool            `json:"multiple"`
		Choices         json.RawMessage `json:"choices"`
		Visible         bool            `json:"visible"`
		Required        bool            `json:"required"`
		AffectsCost     bool            `json:"affectsCost"`
		Editable        bool            `json:"editable"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	f.Code = response.Code
	f.Label = response.Label
	f.Hint = response.Hint
	f.Type = response.Type
	f.AutocompleteURL = response.AutocompleteURL
	f.Multiple = response.Multiple
	f.Choices = nil
	f.ChoiceList = nil
	f.Visible = response.Visible
	f.Required = response.Required
	f.AffectsCost = response.AffectsCost
	f.Editable = response.Editable

	choices := bytes.TrimSpace(response.Choices)
	if len(choices) == 0 || bytes.Equal(choices, []byte("null")) {
		return nil
	}

	var choiceList []DeliveryDataFieldChoice
	if err := json.Unmarshal(choices, &choiceList); err == nil {
		f.ChoiceList = choiceList
		f.Choices = make([]string, 0, len(choiceList))
		for _, choice := range choiceList {
			f.Choices = append(f.Choices, choice.Value)
		}

		return nil
	}

	return json.Unmarshal(choices, &f.Choices)
}

func (p EmbedJSPage) MarshalJSON() ([]byte, error) {
	pageHelpLink := any(nil)
	if len(p.PageHelpLinks) > 0 {
		pageHelpLink = p.PageHelpLinks
	} else if p.PageHelpLink != "" {
		pageHelpLink = p.PageHelpLink
	}

	return json.Marshal(struct {
		Code               string            `json:"code,omitempty"`
		Menu               string            `json:"menu,omitempty"`
		ParentMenuItemCode string            `json:"parentMenuItemCode,omitempty"`
		MenuItemOrdering   int               `json:"menuItemOrdering,omitempty"`
		MenuItemTitle      map[string]string `json:"menuItemTitle,omitempty"`
		PageHelpLink       any               `json:"pageHelpLink,omitempty"`
		IsSettingsMainPage bool              `json:"isSettingsMainPage,omitempty"`
	}{
		Code:               p.Code,
		Menu:               p.Menu,
		ParentMenuItemCode: p.ParentMenuItemCode,
		MenuItemOrdering:   p.MenuItemOrdering,
		MenuItemTitle:      p.MenuItemTitle,
		PageHelpLink:       pageHelpLink,
		IsSettingsMainPage: p.IsSettingsMainPage,
	})
}

func (p *EmbedJSPage) UnmarshalJSON(data []byte) error {
	var response struct {
		Code               string            `json:"code"`
		Menu               string            `json:"menu"`
		ParentMenuItemCode string            `json:"parentMenuItemCode"`
		MenuItemOrdering   int               `json:"menuItemOrdering"`
		MenuItemTitle      map[string]string `json:"menuItemTitle"`
		PageHelpLink       json.RawMessage   `json:"pageHelpLink"`
		IsSettingsMainPage bool              `json:"isSettingsMainPage"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	p.Code = response.Code
	p.Menu = response.Menu
	p.ParentMenuItemCode = response.ParentMenuItemCode
	p.MenuItemOrdering = response.MenuItemOrdering
	p.MenuItemTitle = response.MenuItemTitle
	p.PageHelpLink = ""
	p.PageHelpLinks = nil
	p.IsSettingsMainPage = response.IsSettingsMainPage

	pageHelpLink := bytes.TrimSpace(response.PageHelpLink)
	if len(pageHelpLink) == 0 || bytes.Equal(pageHelpLink, []byte("null")) {
		return nil
	}

	if pageHelpLink[0] == '"' {
		return json.Unmarshal(pageHelpLink, &p.PageHelpLink)
	}

	var translations ConfigurationTranslation
	if err := json.Unmarshal(pageHelpLink, &translations); err != nil {
		return err
	}

	p.PageHelpLinks = translations
	p.PageHelpLink = firstTranslation(translations)

	return nil
}

func firstTranslation(translations ConfigurationTranslation) string {
	for _, lang := range []string{"ru", "en", "es"} {
		if value := translations[lang]; value != "" {
			return value
		}
	}

	for _, value := range translations {
		return value
	}

	return ""
}

func (v *StringOrNumber) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*v = ""
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return err
	}

	switch value := value.(type) {
	case string:
		*v = StringOrNumber(value)
	case json.Number:
		*v = StringOrNumber(value.String())
	default:
		return fmt.Errorf("string or number: expected string or number, got %T", value)
	}

	return nil
}

func (r *IntegrationModuleEditResponse) UnmarshalJSON(data []byte) error {
	var response struct {
		Success bool            `json:"success"`
		Info    json.RawMessage `json:"info"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	r.Success = response.Success
	r.Info = ResponseInfo{}
	r.InfoList = nil
	r.InfoRaw = nil

	info := bytes.TrimSpace(response.Info)
	if len(info) == 0 || bytes.Equal(info, []byte("null")) {
		return nil
	}

	r.InfoRaw = append(r.InfoRaw, response.Info...)

	switch info[0] {
	case '[':
		var infoList []ResponseInfo
		if err := json.Unmarshal(info, &infoList); err != nil {
			return err
		}

		r.InfoList = infoList
		if len(infoList) > 0 {
			r.Info = infoList[0]
		}

		return nil
	case '{':
		return json.Unmarshal(info, &r.Info)
	default:
		return fmt.Errorf("integration module edit info: expected object or array, got %q", info[0])
	}
}

func (r *ResponseInfo) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*r = ResponseInfo{}
		return nil
	}

	if len(data) > 0 && data[0] == '[' {
		var values []json.RawMessage
		if err := json.Unmarshal(data, &values); err != nil {
			return err
		}

		if len(values) == 0 {
			*r = ResponseInfo{}
			return nil
		}
	}

	type responseInfo ResponseInfo
	return json.Unmarshal(data, (*responseInfo)(r))
}

func (a *APIErrorsList) UnmarshalJSON(data []byte) error {
	var i interface{}
	var m APIErrorsList
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		m = make(APIErrorsList, len(e))
		for idx, val := range e {
			m[idx] = fmt.Sprint(val)
		}
	case []interface{}:
		m = make(APIErrorsList, len(e))
		for idx, val := range e {
			m[strconv.Itoa(idx)] = fmt.Sprint(val)
		}
	}

	*a = m
	return nil
}

func (l *StringMap) UnmarshalJSON(data []byte) error {
	var i interface{}
	var m StringMap
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		m = make(StringMap, len(e))
		for idx, val := range e {
			m[idx] = fmt.Sprint(val)
		}
	case []interface{}:
		m = make(StringMap, len(e))
		for idx, val := range e {
			m[strconv.Itoa(idx)] = fmt.Sprint(val)
		}
	}

	*l = m
	return nil
}

func (l *CustomFieldMap) UnmarshalJSON(data []byte) error {
	var i interface{}
	var items CustomFieldMap
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		items = make(CustomFieldMap, len(e))
		for idx, val := range e {
			items[idx] = val
		}
	case []interface{}:
		items = make(CustomFieldMap, len(e))
		for idx, val := range e {
			items[strconv.Itoa(idx)] = val
		}
	}

	*l = items
	return nil
}

func (p *OrderPayments) UnmarshalJSON(data []byte) error {
	var i interface{}
	var m OrderPayments
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		m = make(OrderPayments, len(e))
		for idx, val := range e {
			var res OrderPayment
			err := unmarshalMap(val.(map[string]interface{}), &res)
			if err != nil {
				return err
			}
			m[idx] = res
		}
	case []interface{}:
		m = make(OrderPayments, len(e))
		for idx, val := range e {
			var res OrderPayment
			err := unmarshalMap(val.(map[string]interface{}), &res)
			if err != nil {
				return err
			}
			m[strconv.Itoa(idx)] = res
		}
	}

	*p = m
	return nil
}

func (p *Properties) UnmarshalJSON(data []byte) error {
	var i interface{}
	var m Properties
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		m = make(Properties, len(e))
		for idx, val := range e {
			var res Property
			err := unmarshalMap(val.(map[string]interface{}), &res)
			if err != nil {
				return err
			}
			m[idx] = res
		}
	case []interface{}:
		m = make(Properties, len(e))
		for idx, val := range e {
			var res Property
			err := unmarshalMap(val.(map[string]interface{}), &res)
			if err != nil {
				return err
			}
			m[strconv.Itoa(idx)] = res
		}
	}

	*p = m
	return nil
}

func unmarshalMap(m map[string]interface{}, v interface{}) (err error) {
	var data []byte
	data, err = json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
