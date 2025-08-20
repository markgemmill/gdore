package scraper

import (
	"fmt"
	"gdore/broker"
	"time"

	"github.com/go-rod/rod"
)

func Pause(seconds int) {
	fmt.Println("pausing...")
	time.Sleep(time.Second * time.Duration(seconds))
	fmt.Println("continue...")
}

func ExtractInputValue(el *rod.Element) string {
	prop, err := el.Property("value")
	if err != nil {
		fmt.Printf("error access value attribute: '%s'\n", err)
		return ""
	}
	value := prop.String()
	fmt.Printf("el.value='%s'\n", value)
	return value
}

func GetInputText(el *rod.Element, selector string) string {
	fmt.Printf("%s\n", el)
	fmt.Printf("get element: %s\n", selector)
	inputEl := el.MustElement(selector)
	if inputEl == nil {
		fmt.Printf("selector '%s' returned nil\n", selector)
		return ""
	}

	prop, err := inputEl.Property("value")
	fmt.Printf("property: %s\n", prop)

	if err != nil {
		fmt.Printf("error access value property: '%s'\n", err)
		return ""
	}

	value := prop.String()

	fmt.Printf("el.value='%s'\n", value)
	return value
}

func GetInputBool(el *rod.Element, selector string) bool {
	fmt.Printf("%s\n", el)
	fmt.Printf("get element: %s\n", selector)
	inputEl := el.MustElement(selector)
	if inputEl == nil {
		fmt.Printf("selector '%s' returned nil\n", selector)
		return false
	}

	prop, err := inputEl.Property("checked")
	fmt.Printf("property: %s\n", prop)

	if err != nil {
		fmt.Printf("error access value property: '%s'\n", err)
		return false
	}

	value := prop.Bool()

	fmt.Printf("el.value=%t\n", value)
	return value
}

func GetElementText(el *rod.Element, selector string) string {
	value, err := el.MustElement(selector).Text()
	if err != nil {
		fmt.Printf("%s\n", err)
		return ""
	}
	return value
}

func GetText(el *rod.Element) string {
	value, err := el.Text()
	if err != nil {
		fmt.Printf("%s\n", err)
		return ""
	}
	return value
}

var messanger *broker.MessageBroker = broker.NewMessageBroker()
