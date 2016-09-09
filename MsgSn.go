package Common

import (
	"Common/logger"
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

//TODO:
//slice
func MsgUsn(r *bytes.Buffer, data interface{}) error {
	//	log.Println("111111")
	// s := reflect.TypeOf(data)
	// log.Println(reflect.ValueOf(data).CanSet())

	var dataKind reflect.Kind

	tmData, err := data.(reflect.Value)
	if err {
		dataKind = tmData.Kind()
		//		log.Println("first import", dataKind)
		//		log.Println(tmData.CanSet())
	} else {
		dataKind = reflect.TypeOf(data).Elem().Kind()
		//		log.Println("second import", dataKind)
		tmData = reflect.ValueOf(data).Elem()
	}

	if !tmData.CanSet() {
		logger.Error("data not canset")
		return errors.New("data not canset")
	}

	switch dataKind {
	case reflect.Int8:
		var _t int8
		binary.Read(r, binary.LittleEndian, &_t)
		tmData.SetInt(int64(_t))
	case reflect.Int16:
		var _t int16
		binary.Read(r, binary.LittleEndian, &_t)
		tmData.SetInt(int64(_t))
	case reflect.Int32, reflect.Int:
		var _t int32
		binary.Read(r, binary.LittleEndian, &_t)
		tmData.SetInt(int64(_t))
	case reflect.Int64:
		var _t int64
		binary.Read(r, binary.LittleEndian, &_t)
		tmData.SetInt(int64(_t))
	case reflect.Uint8:
		var _t uint8
		binary.Read(r, binary.LittleEndian, &_t)
		tmData.SetUint(uint64(_t))
	case reflect.Uint16:
		var _t uint16
		binary.Read(r, binary.LittleEndian, &_t)
		tmData.SetUint(uint64(_t))
	case reflect.Uint32, reflect.Uint:
		var _t uint32
		binary.Read(r, binary.LittleEndian, &_t)
		tmData.SetUint(uint64(_t))
	case reflect.Uint64:
		var _t uint64
		binary.Read(r, binary.LittleEndian, &_t)
		tmData.SetUint(uint64(_t))
	case reflect.Slice:
		if tmData.Type().String() == string("[]uint8") {
			var len uint16
			binary.Read(r, binary.LittleEndian, &len)
			_uint8_buff := make([]byte, len)
			binary.Read(r, binary.LittleEndian, &_uint8_buff)
			tmData.SetBytes(_uint8_buff)
		} else {
			logger.Error("interface not support ", dataKind, " ", tmData.Type())
			logger.Info(reflect.SliceOf(reflect.TypeOf(data)))
		}

		// tmData.SetLen(2)
		// for i := 0; i < tmData.Len(); i++ {
		// 	ReadParam(r, tmData.Index(i))
		// }
	case reflect.Array:
		_uint8_buff := make([]byte, tmData.Len())
		binary.Read(r, binary.LittleEndian, &_uint8_buff)
		for i := 0; i < tmData.Len(); i++ {
			tmData.Index(i).SetUint(uint64(_uint8_buff[i]))
		}
	case reflect.String:
		var len uint8
		binary.Read(r, binary.LittleEndian, &len)
		str := make([]byte, len)
		binary.Read(r, binary.LittleEndian, &str)
		tmData.SetString(string(str))
	case reflect.Struct:
		//		log.Println(tmData.NumField())
		for i := 0; i < tmData.NumField(); i++ {
			f := tmData.Field(i)
			//			log.Println(f.CanSet())
			//			log.Println(f.CanInterface())
			//			log.Println(f.Kind())
			// readParamKind(r, f, f.Kind())
			MsgUsn(r, f)
		}
	default:
		logger.Error("interface not support ", dataKind)
		return errors.New("interface not support")
	}
	return nil
}

//TODO:
//slice
func MsgSn(w *bytes.Buffer, data interface{}, l int) error {
	s := reflect.TypeOf(data)
	//	log.Println(s.Kind())

	switch s.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		binary.Write(w, binary.LittleEndian, data)
	case reflect.Int:
		i, _ := data.(int)
		binary.Write(w, binary.LittleEndian, int32(i))
	case reflect.Uint:
		i, _ := data.(uint)
		binary.Write(w, binary.LittleEndian, uint32(i))
	case reflect.String:
		s, _ := data.(string)
		if l == 1 && len(s) > 255 {
			logger.Error("string len too big", len(s))
			return errors.New("string len too big")
		}
		if l == 1 {
			binary.Write(w, binary.LittleEndian, byte(len(s)))
		} else if l == 4 {
			binary.Write(w, binary.LittleEndian, int32(len(s)))
		} else {
			binary.Write(w, binary.LittleEndian, int64(len(s)))
		}

		if len(s) > 0 {
			binary.Write(w, binary.LittleEndian, []byte(s))
		}
	case reflect.Slice:
		if s.String() == string("[]uint8") {
			sli, _ := data.([]byte)
			_str_len := uint16(len(sli))
			//	log.Println(">>>>>>>>>>>> []byte = ", _str_len)
			binary.Write(w, binary.LittleEndian, _str_len)
			if len(sli) > 0 {
				binary.Write(w, binary.LittleEndian, sli)
			}
		} else {
			logger.Error("interface not support ", s.Kind(), " type:", s.String())
			return errors.New("interface not support")
		}
		// log.Println(reflect.ValueOf(data).Len())
		// log.Println(reflect.ValueOf(data).Index(0))
		// for i := 0; i < reflect.ValueOf(data).Len(); i++ {
		// 	MsgSn(w, reflect.ValueOf(data).Index(i).Interface())
		// }
	case reflect.Array:
		v := reflect.ValueOf(data)
		for i := 0; i < v.Len(); i++ {
			binary.Write(w, binary.LittleEndian, v.Index(i).Interface().(uint8))
		}
	case reflect.Struct:
		for i := 0; i < s.NumField(); i++ {
			f := reflect.ValueOf(data).Field(i)
			if f.CanInterface() {
				MsgSn(w, f.Interface(), l)
			} else {
				MsgSn(w, f, l)
			}
		}
	default:
		logger.Error("interface not support ", s.Kind())
		return errors.New("interface not support")
	}

	return nil
}
