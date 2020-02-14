package sconfig

import (
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type flagSetter func(
	flagSet *pflag.FlagSet,
	longFlag string,
	shortFlag string,
	info string,
	def interface{},
	ptr unsafe.Pointer,
)

func bindPFlag(
	v *viper.Viper,
	flagSet *pflag.FlagSet,
	conf interface{},
	flag string,
	def interface{},
	field reflect.StructField,
	value reflect.Value,
) error {
	info, _ := field.Tag.Lookup("info")

	longFlag, shortFlag, err := parseFlagOptions(flag)
	if err != nil {
		return err
	}

	typ := field.Type
	ptr := unsafe.Pointer(value.Addr().Pointer())

	var setFlag flagSetter

	switch typ.Kind() {
	case reflect.String:
		setFlag = setStringFlag
	case reflect.Bool:
		setFlag = setBoolFlag
	case reflect.Int:
		setFlag = setIntFlag
	case reflect.Int16:
		setFlag = setInt16Flag
	case reflect.Int32:
		setFlag = setInt32Flag
	case reflect.Int64:
		if isDurationType(typ) {
			setFlag = setDurationFlag
		} else {
			setFlag = setInt64Flag
		}
	case reflect.Uint:
		setFlag = setUintFlag
	case reflect.Uint8:
		setFlag = setUint8Flag
	case reflect.Uint16:
		setFlag = setUint16Flag
	case reflect.Uint32:
		setFlag = setUint32Flag
	case reflect.Uint64:
		setFlag = setUint64Flag
	case reflect.Float32:
		setFlag = setFloat32Flag
	case reflect.Float64:
		setFlag = setFloat64Flag
	case reflect.Slice:
		sliceKind := typ.Elem().Kind()
		switch sliceKind {
		case reflect.String:
			setFlag = setStringSliceFlag
		case reflect.Bool:
			setFlag = setBoolSliceFlag
		case reflect.Int:
			setFlag = setIntSliceFlag
		case reflect.Int64:
			if isDurationType(typ.Elem()) {
				setFlag = setDurationSliceFlag
			} else {
				return &ErrUnsupportedFieldType{Type: typ.Name()}
			}
		default:
			return &ErrUnsupportedFieldType{Type: typ.Name()}
		}
	default:
		return &ErrUnsupportedFieldType{Type: typ.Name()}
	}

	setFlag(flagSet, longFlag, shortFlag, info, def, ptr)
	v.BindPFlag(longFlag, flagSet.Lookup(longFlag))

	return nil
}

func parseFlagOptions(
	flag string,
) (longFlag string, shortFlag string, err error) {
	flags := strings.Split(flag, ",")

	if len(flags) == 1 {
		return flags[0], "", nil
	}

	if len(flags) == 2 {
		if len(flags[1]) != 1 {
			return "", "", &ErrInvalidFlagFormat{Format: flag}
		}
		return flags[0], flags[1], nil
	}

	return "", "", &ErrInvalidFlagFormat{Format: flag}
}

func isDurationType(typ reflect.Type) bool {
	return typ.PkgPath() == "time" && typ.Name() == "Duration"
}

func setStringFlag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.StringVarP((*string)(ptr), lf, sf, cast.ToString(def), info)
}

func setBoolFlag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.BoolVarP((*bool)(ptr), lf, sf, cast.ToBool(def), info)
}

func setIntFlag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.IntVarP((*int)(ptr), lf, sf, cast.ToInt(def), info)
}

func setInt16Flag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.Int16VarP((*int16)(ptr), lf, sf, cast.ToInt16(def), info)
}

func setInt32Flag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.Int32VarP((*int32)(ptr), lf, sf, cast.ToInt32(def), info)
}

func setInt64Flag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.Int64VarP((*int64)(ptr), lf, sf, cast.ToInt64(def), info)
}

func setUintFlag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.UintVarP((*uint)(ptr), lf, sf, cast.ToUint(def), info)
}

func setUint8Flag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.Uint8VarP((*uint8)(ptr), lf, sf, cast.ToUint8(def), info)
}

func setUint16Flag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.Uint16VarP((*uint16)(ptr), lf, sf, cast.ToUint16(def), info)
}

func setUint32Flag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.Uint32VarP((*uint32)(ptr), lf, sf, cast.ToUint32(def), info)
}

func setUint64Flag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.Uint64VarP((*uint64)(ptr), lf, sf, cast.ToUint64(def), info)
}

func setDurationFlag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.DurationVarP((*time.Duration)(ptr), lf, sf, cast.ToDuration(def), info)
}

func setFloat32Flag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.Float32VarP((*float32)(ptr), lf, sf, cast.ToFloat32(def), info)
}

func setFloat64Flag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.Float64VarP((*float64)(ptr), lf, sf, cast.ToFloat64(def), info)
}

func setStringSliceFlag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.StringSliceVarP((*[]string)(ptr), lf, sf, cast.ToStringSlice(def), info)
}

func setBoolSliceFlag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.BoolSliceVarP((*[]bool)(ptr), lf, sf, cast.ToBoolSlice(def), info)
}

func setIntSliceFlag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.IntSliceVarP((*[]int)(ptr), lf, sf, cast.ToIntSlice(def), info)
}

func setDurationSliceFlag(fs *pflag.FlagSet, lf, sf, info string, def interface{}, ptr unsafe.Pointer) {
	fs.DurationSliceVarP((*[]time.Duration)(ptr), lf, sf, cast.ToDurationSlice(def), info)
}
