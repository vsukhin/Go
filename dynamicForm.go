package main

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MyForm struct {
	UserName     string        `required:"true" field:"name" name:"Имя пользователя" type:"text"`
	UserPassword string        `required:"true" field:"password" name:"Пароль пользователя" type:"password"`
	Customer     string        `required:"false" field:"customer" name:"Имя и Фамилия" type:"textarea" default="false"`
	Resident     bool          `field:"resident" type:"radio" radio:"1;checked" name:"Резидент РФ"`
	NoResident   bool          `field:"resident" type:"radio" radio:"2" name:"Нерезидент РФ"`
	Gender       string        `field:"gender" name:"Пол" type:"select" select:"Неизвестный=3;selected,Мужской=1,Женский=2"`
	Salary       float64       `field:"salary" name:"Зарплата" type:"text" default:"true"`
	Age          int64         `field:"age" name:"Возраст" type:"text" default:"true"`
	Token        string        `field:"token" type:"hidden" default:"true"`
	Token_1      uint          `field:"" type:"hidden" default:"true"`
	Token_2      uint64        `field:"-" type:"hidden" default:"true"`
	Subscription bool          `field:"subscription" type:"checkbox" checkbox:"1;checked" name:"Подписка"`
	Agreement    bool          `field:"agreement" type:"checkbox" checkbox:"1" name:"Согласие"`
	Secret       string        `field:"secret" name:"" type:"hidden" default:"true"`
	Kids         int           `field:"kids" name:"Количество детей" default:"false"`
	ClickMe      string        `field:"clickme" type:"button" name:"Для дополнительной информации" default:"true"`
	Timestamp    time.Duration `field:"timestamp" name:"Метка времени" type:"hidden" default:"true"`
}

func FormRead(formData *MyForm, request *http.Request) (err error) {
	if formData == nil {
		return errors.New("MyForm structure can't be nil")
	}
	if request == nil {
		return errors.New("Server request can't be nil")
	}
	fmt.Println("Reading data from form...")

	structAddr := reflect.ValueOf(formData).Elem()
	typeOfStruct := structAddr.Type()
	for i := 0; i < structAddr.NumField(); i++ {
		field := structAddr.Field(i)
		fmt.Println("Analyzing structure field")

		fieldStruct := typeOfStruct.Field(i).Name
		fmt.Println("Name: ", fieldStruct)

		tag := typeOfStruct.Field(i).Tag

		fieldTag := tag.Get("field")
		fmt.Println("Field tag: ", fieldTag)
		if fieldTag == "" || fieldTag == "-" {
			continue
		}
		matched, err := regexp.MatchString("^[A-Za-z]*$", fieldTag)
		if err != nil || !matched {
			return fmt.Errorf("Wrong structure field value", fieldTag)
		}

		fieldValue := field.Interface()
		fmt.Println("Value: ", fieldValue)

		postValue := request.FormValue(fieldTag)
		fmt.Println("Posted form field: ", postValue)

		fieldRequired := tag.Get("required")
		fmt.Println("Required tag: ", fieldRequired)
		if fieldRequired != "true" && fieldRequired != "false" && fieldRequired != "" {
			return fmt.Errorf("Wrong structure required value", fieldRequired)
		}

		if fieldRequired == "true" && postValue == "" {
			return fmt.Errorf("Required value can't be empty", postValue)
		}

		fieldType := field.Type().Kind()
		fmt.Println("Type: ", fieldType)

		typeField := tag.Get("type")
		fmt.Println("Type tag: ", typeField)

		switch fieldType {
		case reflect.Bool:
			if typeField == "radio" {
				fieldRadio := tag.Get("radio")
				fmt.Println("Tag radio", fieldRadio)

				values := strings.Split(fieldRadio, ";")
				if len(values) < 1 || len(values) > 2 {
					return fmt.Errorf("Wrong structure radio value ", fieldRadio)
				}
				if len(values) == 2 {
					if values[1] != "checked" {
						return fmt.Errorf("Wrong structure radio value part", values[1])
					}
				}
				if postValue == values[0] {
					field.SetBool(true)
				} else {
					field.SetBool(false)
				}
			} else {
				if postValue != "" {
					field.SetBool(true)
				} else {
					field.SetBool(false)
				}
			}
		case reflect.Float64:
			if postValue != "" {
				f, err := strconv.ParseFloat(postValue, 64)
				if err != nil {
					return fmt.Errorf("Form data mistmaches structure field type", err)
				} else {
					field.SetFloat(f)
				}
			}
		case reflect.Int:
			if postValue != "" {
				i, err := strconv.ParseInt(postValue, 0, 0)
				if err != nil {
					return fmt.Errorf("Form data mistmaches structure field type", err)
				} else {
					field.SetInt(i)
				}
			}
		case reflect.Int64:
			if postValue != "" {
				i, err := strconv.ParseInt(postValue, 0, 64)
				if err != nil {
					return fmt.Errorf("Form data mistmaches structure field type", err)
				} else {
					field.SetInt(i)
				}
			}
		case reflect.Uint:
			if postValue != "" {
				u, err := strconv.ParseUint(postValue, 0, 0)
				if err != nil {
					return fmt.Errorf("Form data mistmaches structure field type", err)
				} else {
					field.SetUint(u)
				}
			}
		case reflect.Uint64:
			if postValue != "" {
				u, err := strconv.ParseUint(postValue, 0, 64)
				if err != nil {
					return fmt.Errorf("Form data mistmaches structure field type", err)
				} else {
					field.SetUint(u)
				}
			}
		case reflect.String:
			field.SetString(postValue)
		default:
			return fmt.Errorf("Wrong structure field type", fieldType)
		}
	}

	return nil
}

func processCheckedValues(tag *reflect.StructTag, name string) (input string, err error) {
	field := tag.Get(name)
	fmt.Println("Tag %s: %s", name, field)

	input = " type='" + name + "'"

	values := strings.Split(field, ";")
	if len(values) < 1 || len(values) > 2 {
		return "", fmt.Errorf("Wrong structure %s value %s", name, field)
	}
	input += " value='" + html.EscapeString(values[0]) + "'"
	if len(values) == 2 {
		if values[1] == "checked" {
			input += " checked"
		} else {
			return "", fmt.Errorf("Wrong structure %s value part %s", name, values[1])
		}
	}

	return input, nil
}

func FormCreate(formData *MyForm) (form string, err error) {
	form = ""
	if formData == nil {
		return form, errors.New("MyForm structure can't be nil")
	}
	fmt.Println("Generating form from data...")

	form = "<html><head><title>Dynamicly generated form</title></head><body><form action='/form' method='post' enctype='multipart/form-data'>"

	structAddr := reflect.ValueOf(formData).Elem()
	typeOfStruct := structAddr.Type()
	for i := 0; i < structAddr.NumField(); i++ {
		field := structAddr.Field(i)
		fmt.Println("Analyzing structure field")

		fieldStruct := typeOfStruct.Field(i).Name
		fmt.Println("Name: ", fieldStruct)

		fieldValue := field.Interface()
		fmt.Println("Value: ", fieldValue)

		tag := typeOfStruct.Field(i).Tag

		fieldTag := tag.Get("field")
		fmt.Println("Field tag: ", fieldTag)
		if fieldTag == "" || fieldTag == "-" {
			continue
		}
		matched, err := regexp.MatchString("^[A-Za-z]*$", fieldTag)
		if err != nil || !matched {
			return "", fmt.Errorf("Wrong structure field value", fieldTag)
		}

		inputTag := " name='" + fieldTag + "'"

		fieldType := field.Type().Kind()
		fmt.Println("Type: ", fieldType)

		stringValue := ""
		switch fieldType {
		case reflect.Bool, reflect.Float64, reflect.String:
			stringValue = fmt.Sprint(fieldValue)
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
			stringValue = fmt.Sprintf("%d", fieldValue)
		default:
			return "", fmt.Errorf("Wrong structure field type", fieldType)
		}

		fieldRequired := tag.Get("required")
		fmt.Println("Required tag: ", fieldRequired)
		if fieldRequired != "true" && fieldRequired != "false" && fieldRequired != "" {
			return "", fmt.Errorf("Wrong structure required value", fieldRequired)
		}

		if fieldRequired == "true" {
			inputTag += " required"
		}

		fieldDefault := tag.Get("default")
		fmt.Println("Default tag: ", fieldDefault)
		if fieldDefault != "true" && fieldDefault != "false" && fieldDefault != "" {
			return "", fmt.Errorf("Wrong structure default value", fieldDefault)
		}

		typeField := tag.Get("type")
		fmt.Println("Type tag: ", typeField)
		if typeField == "radio" {
			input, err := processCheckedValues(&tag, "radio")
			if err != nil {
				return "", fmt.Errorf("Error during processing radio checked values", err)
			}
			inputTag = "<input" + inputTag + input
		} else if typeField == "checkbox" {
			input, err := processCheckedValues(&tag, "checkbox")
			if err != nil {
				return "", fmt.Errorf("Error during processing checkbox checked values", err)
			}
			inputTag = "<input" + inputTag + input
		} else if typeField == "select" {
			fieldSelect := tag.Get("select")
			fmt.Println("Select tag: ", fieldSelect)

			values := strings.Split(fieldSelect, ",")
			if len(values) == 0 {
				return "", fmt.Errorf("Wrong structure select value", fieldSelect)
			}
			inputTag += ">"
			for _, v := range values {
				options := strings.Split(v, ";")
				if len(options) < 1 || len(options) > 2 {
					return "", fmt.Errorf("Wrong structure option value", v)
				}

				items := strings.Split(options[0], "=")
				if len(items) != 2 {
					return "", fmt.Errorf("Wrong structure option item value", options[0])
				}
				inputTag += "<option value='" + html.EscapeString(items[1]) + "'"
				if len(options) == 2 {
					if options[1] == "selected" {
						inputTag += " selected"
					} else {
						return "", fmt.Errorf("Wrong structure option value part", options[1])
					}
				}
				inputTag += ">" + items[0] + "</option>"
			}
			inputTag = "<select" + inputTag + "</select"
		} else if typeField == "textarea" {
			inputTag = "<textarea" + inputTag + ">"

			if fieldDefault == "true" {
				inputTag += stringValue
			}
			inputTag += "</textarea"
		} else if typeField == "text" || typeField == "password" || typeField == "hidden" || typeField == "button" || typeField == "" {
			if typeField == "" {
				typeField = "text"
			}
			inputTag = "<input" + inputTag + " type='" + typeField + "'"

			if fieldDefault == "true" {
				inputTag += " value='" + html.EscapeString(stringValue) + "'"
			}
		} else {
			return "", fmt.Errorf("Wrong structure type value", typeField)
		}

		fieldName := tag.Get("name")
		fmt.Println("Name tag: ", fieldName)

		inputTag += ">"

		if typeField != "hidden" {
			inputTag = "<label for name='" + fieldTag + "'>" + fieldName + "</label>" + inputTag + "<br>"
		}

		form += inputTag
	}
	form += "<input type='submit' value='Submit'></form></body></html>"

	return form, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var fdOut *MyForm = new(MyForm)

		err := FormRead(fdOut, r)

		if err != nil {
			fmt.Println("Can't read form for the provided structure: ", err)
		} else {
			fmt.Println("Read posted form data: ", fdOut)
		}
	case "GET":
		var fdIn *MyForm = &MyForm{
			Age:       18,
			Token:     "345625145123451234123412342345",
			Timestamp: 1234567890,
			Salary:    10.50,
			ClickMe:   "Нажми меня",
		}

		form, err := FormCreate(fdIn)

		if err != nil {
			fmt.Println("Can't create form from the provided structure: ", err)
		} else {
			fmt.Fprintln(w, form)
		}
	default:
		fmt.Println("Method is not supported: ", r.Method)
	}
}

func main() {
	fmt.Println("Starting http server at: ", time.Now())

	http.HandleFunc("/form", handler)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Can't launch http server: ", err)
	}
}
