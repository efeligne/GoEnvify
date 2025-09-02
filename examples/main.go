package main

import (
	"log"
	"reflect"
	"time"
)

type Person struct {
	Name string        `default:"John Doe" description:"Person name"             env:"NAME"`
	Age  int           `default:"30"       description:"Person age"              env:"AGE"`
	I8   int8          `default:"8"        description:"Just test int8"          env:"I8"`
	Dur  time.Duration `default:"1h"       description:"Just test time.Duration" env:"DUR"`
	IsOk bool          `default:"false"    description:"Just test bool"          env:"IS_OK"`
}

type StructMeta struct {
	ReflectValue reflect.Value
	ReflectType  reflect.Type
}

func GetStructMeta[T any](s T) StructMeta {
	return StructMeta{
		ReflectValue: reflect.ValueOf(&s).Elem(),
		ReflectType:  reflect.TypeOf(s),
	}
}

type FieldInfo struct {
	EnvTag         string
	DescriptionTag string
	DefaultTag     string
	Name           string
	Value          string
	Type           string
	ReflectValue   reflect.Value
}

func GetFieldInfo(val reflect.Value, typ reflect.Type, idx int) FieldInfo {
	field := typ.Field(idx)
	value := val.Field(idx)

	return FieldInfo{
		EnvTag:         field.Tag.Get("env"),
		DescriptionTag: field.Tag.Get("description"),
		DefaultTag:     field.Tag.Get("default"),
		Name:           field.Name,
		Value:          value.String(),
		Type:           value.Type().Name(),
		ReflectValue:   value,
	}
}

//nolint:revive,add-constant
func main() { //nolint:gocognit,cyclop,funlen
	var person Person

	const (
		newName = "John Smith"
		newAge  = "34"
		newI8   = "126"
		newDur  = "2h30m23s10ms"
		newIsOk = "true"
	)

	reft := GetStructMeta(person).ReflectType
	refv := GetStructMeta(person).ReflectValue

	for i := 0; i < refv.NumField(); i++ {
		log.Printf("%#v", GetFieldInfo(refv, reft, i))
		// if envTag == "NAME" {
		// 	if value.CanSet() {
		// 		value.SetString(newName)
		// 	}
		// }
		//
		// if envTag == "AGE" {
		// 	if value.CanSet() {
		// 		iAge, _ := strconv.ParseInt(newAge, 10, 64)
		// 		value.SetInt(iAge)
		// 	}
		// }
		//
		// if envTag == "I8" {
		// 	if value.CanSet() {
		// 		iI8, err := strconv.ParseInt(newI8, 10, 8)
		// 		if err != nil {
		// 			log.Fatalln(err)
		// 		}
		//
		// 		value.SetInt(iI8)
		// 	}
		// }
		//
		// if envTag == "DUR" {
		// 	if value.CanSet() {
		// 		iDur, err := time.ParseDuration(newDur)
		// 		if err != nil {
		// 			log.Fatalln(err)
		// 		}
		//
		// 		value.SetInt(int64(iDur))
		// 	}
		// }
		//
		// if envTag == "IS_OK" {
		// 	if value.CanSet() {
		// 		boolVal, err := strconv.ParseBool(newIsOk)
		// 		if err != nil {
		// 			log.Fatalln(err)
		// 		}
		//
		// 		value.SetBool(boolVal)
		// 	}
		// }
		//
		// log.Printf(
		// 	`Поле: %v,
		// 	значение: %v,
		// 	тип: %v,
		// 	env: %v,
		// 	default: %v,
		// 	description: %v`,
		// 	field.Name,
		// 	value,
		// 	value.Type(),
		// 	envTag,
		// 	defaultTag,
		// 	descriptionTag,
		// )
	}
}
