package main

import (
	"errors"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	// todo

	if reflect.ValueOf(out).Kind() != reflect.Ptr {
		return errors.New("error")
	}

	reflectDataT := reflect.TypeOf(data)
	switch reflectDataT.Kind() {
	case reflect.Map:
		if reflect.ValueOf(out).Elem().Kind() != reflect.Struct {
			return errors.New("error")
		}
		err := getOut(out, data.(map[string]interface{}))
		if err != nil {
			return err
		}
	case reflect.Slice:
		if reflect.ValueOf(out).Elem().Kind() != reflect.Slice {
			return errors.New("error")
		}
		for _, value := range data.([]interface{}) {

			p := reflect.New(reflect.ValueOf(out).Elem().Type().Elem())
			i := p.Elem().Addr().Interface()

			e := i2s(value, i)
			if e != nil {
				return e
			}

			s := reflect.Append(reflect.ValueOf(out).Elem(), reflect.ValueOf(i).Elem())
			reflect.ValueOf(out).Elem().Set(s)
		}
	default:
		return errors.New("error")
	}

	return nil
}

func getOut(out interface{}, in map[string]interface{}) error {

	reflectOutV := reflect.ValueOf(out).Elem() // внутренности out

	isError := false
	var err error

	for key, value := range in {

		if isError {
			break
		}

		field, ok := reflectOutV.Type().FieldByName(key) // поле из out по ключу из базовой структуры
		if !ok {
			continue
		}

		sourceT := reflect.ValueOf(value).Type().Kind()
		destT := field.Type.Kind()

		switch destT {
		case reflect.Map:
			e := i2s(value, reflectOutV.FieldByName(key).Addr().Interface())
			if e != nil {
				isError = true
				err = e
			}
		case reflect.Slice:
			e := i2s(value, reflectOutV.FieldByName(key).Addr().Interface())
			if e != nil {
				isError = true
				err = e
			}
		case reflect.Int:
			if sourceT == reflect.Float64 {
				reflectOutV.FieldByName(key).Set(reflect.ValueOf(int(value.(float64))))
			} else if sourceT == reflect.Int64 {
				reflectOutV.FieldByName(key).Set(reflect.ValueOf(value.(int64)))
			} else {
				isError = true
				err = errors.New("")
			}
		case reflect.Float64:
			if sourceT == reflect.Int {
				reflectOutV.FieldByName(key).Set(reflect.ValueOf(value.(int64)))
			} else if sourceT == reflect.Float64 {
				reflectOutV.FieldByName(key).Set(reflect.ValueOf(int(value.(float64))))
			} else {
				isError = true
				err = errors.New("")
			}
		case reflect.Struct:
			if sourceT == destT {
				reflectOutV.FieldByName(key).Set(reflect.ValueOf(value))
			} else if sourceT == reflect.Map || sourceT == reflect.Slice {
				e := i2s(value, reflectOutV.FieldByName(key).Addr().Interface())
				if e != nil {
					isError = true
					err = e
				}
			} else {
				isError = true
				err = errors.New("")
			}
		case reflect.Bool:
			if sourceT == destT {
				reflectOutV.FieldByName(key).Set(reflect.ValueOf(value))
			} else {
				isError = true
				err = errors.New("")
			}
		case reflect.String:
			if sourceT == destT {
				reflectOutV.FieldByName(key).Set(reflect.ValueOf(value))
			} else {
				isError = true
				err = errors.New("")
			}
		default:
			reflectOutV.FieldByName(key).Set(reflect.ValueOf(reflect.ValueOf(value).Elem()))
		}

	}

	return err
}
