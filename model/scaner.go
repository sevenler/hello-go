package model

import (
sqlCon "database/sql"
"database/sql/driver"
"fmt"
"git.in.zhihu.com/zhihu/hello/utils"
_ "github.com/go-sql-driver/mysql"
"reflect"
"time"
)

var (
	TypeInt reflect.Type = reflect.TypeOf(int(0))
	TypeInt64 reflect.Type = reflect.TypeOf(int64(0))
	TypeInt32 reflect.Type = reflect.TypeOf(int32(0))
	TypeFloat64 reflect.Type = reflect.TypeOf(float64(0))
	TypeString reflect.Type = reflect.TypeOf("")
	TypeBool reflect.Type = reflect.TypeOf(true)
	TypeTime reflect.Type = reflect.TypeOf(time.Now())

	TypeIntPtr reflect.Type = reflect.TypeOf(utils.IntPtr(0))
	TypeInt64Ptr reflect.Type = reflect.TypeOf(utils.Int64Ptr(0))
	TypeInt32Ptr reflect.Type = reflect.TypeOf(utils.Int32Ptr(0))
	TypeFloat64Ptr reflect.Type = reflect.TypeOf(utils.Float64Ptr(0))
	TypeStringPtr reflect.Type = reflect.TypeOf(utils.StringPtr(""))
	TypeBoolPtr reflect.Type = reflect.TypeOf(utils.BoolPtr(true))
	TypeTimePtr reflect.Type = reflect.TypeOf(utils.TimePtr(time.Now()))
)

var ColumnNullableMap = map[reflect.Type]Nullable{
	TypeInt:     &sqlCon.NullInt64{},
	TypeInt64:   &sqlCon.NullInt64{},
	TypeInt32:   &sqlCon.NullInt32{},
	TypeFloat64: &sqlCon.NullFloat64{},
	TypeString:  &sqlCon.NullString{},
	TypeBool:    &sqlCon.NullBool{},
	TypeTime:    &sqlCon.NullTime{},

	TypeIntPtr:     &sqlCon.NullInt64{},
	TypeInt64Ptr:   &sqlCon.NullInt64{},
	TypeInt32Ptr:   &sqlCon.NullInt32{},
	TypeFloat64Ptr: &sqlCon.NullFloat64{},
	TypeStringPtr:  &sqlCon.NullString{},
	TypeBoolPtr:    &sqlCon.NullBool{},
	TypeTimePtr:    &sqlCon.NullTime{},
}

type Nullable interface {
	Value() (driver.Value, error)
}

// 将SQL查询的数据动态转化成 model 对应的结构体
func Scan(model interface{}, sql string, args ...interface{}) (interface{}, error){
	db, err := sqlCon.Open("mysql","root:@/butterfly?parseTime=true")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil{
		return nil, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil{
		return nil, err
	}

	// key 为 json 标签，value 为字段类型结构体
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr{
		t = t.Elem()
	}
	tagFiled := map[string]reflect.StructField{}
	for i := 0; i< t.NumField();i++{
		field := t.Field(i)
		v := field.Tag.Get("json")
		tagFiled[v] = field
	}

	// sql.Nullxxx 是为了处理 db 中的值是 Null 的情况
	// 所有的 Nullxxx 是都可以调用 Value， Value 返回空表示是一个 Null 的对象
	valPtr := make([]interface{}, 0)
	for i := 0; i< len(columns);i++{
		name := columns[i]
		if field, ok := tagFiled[name]; ok{
			nullable, ok := ColumnNullableMap[field.Type]
			if !ok{
				return nil, fmt.Errorf("not support type %v of field %s", field.Type, name)
			}
			valPtr = append(valPtr, nullable)
		}else{
			// 如果 columns 中的列不包含在结构体字段中，通过 column 类型示例补全 scan 数组
			ct := columnTypes[i]
			t := ct.ScanType()
			nullable := reflect.New(t).Interface()
			valPtr = append(valPtr, nullable)
		}
	}

	var slice = make([]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(valPtr...)
		if err != nil {
			return nil, err
		}
		var obj = reflect.New(t).Elem()
		for i, val := range valPtr {
			value, err := val.(Nullable).Value()
			column := columns[i]
			field := obj.FieldByName(tagFiled[column].Name)
			if err == nil && value != nil && field.IsValid(){
				switch field.Type() {
				case TypeInt:
					field.SetInt(value.(int64))
				case TypeInt64:
					field.SetInt(value.(int64))
				case TypeInt32:
					field.SetInt(value.(int64))
				case TypeFloat64:
					field.SetFloat(value.(float64))
				case TypeString:
					field.SetString(value.(string))
				case TypeBool:
					field.SetBool(value.(bool))
				case TypeTime:
					t := value.(time.Time)
					field.Set(reflect.ValueOf(t))

				case TypeIntPtr:
					field.Set(reflect.ValueOf(utils.IntPtr(int(value.(int64)))))
				case TypeInt64Ptr:
					field.Set(reflect.ValueOf(utils.Int64Ptr(value.(int64))))
				case TypeInt32Ptr:
					field.Set(reflect.ValueOf(utils.Int32Ptr(value.(int32))))
				case TypeFloat64Ptr:
					field.Set(reflect.ValueOf(utils.Float64Ptr(value.(float64))))
				case TypeStringPtr:
					field.Set(reflect.ValueOf(utils.StringPtr(value.(string))))
				case TypeBoolPtr:
					field.Set(reflect.ValueOf(utils.BoolPtr(value.(bool))))
				case TypeTimePtr:
					t := value.(time.Time)
					field.Set(reflect.ValueOf(utils.TimePtr(t)))

				default:
					return nil, fmt.Errorf("error set row %s got type %v",
						column, field.Type())
				}
			}
		}
		slice = append(slice, obj.Interface())
	}
	return slice, err
}