package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"git.in.zhihu.com/zhihu/hello/utils"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"time"
)

//
//CREATE TABLE `user0` (
//`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
//`name` varchar(36) DEFAULT NULL COMMENT 'name',
//`gender` varchar(36) DEFAULT NULL COMMENT 'gender',
//`created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
//`updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//PRIMARY KEY (`id`)
//)

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
	TypeInt:     &sql.NullInt64{},
	TypeInt64:   &sql.NullInt64{},
	TypeInt32:   &sql.NullInt32{},
	TypeFloat64: &sql.NullFloat64{},
	TypeString:  &sql.NullString{},
	TypeBool:    &sql.NullBool{},
	TypeTime:    &sql.NullTime{},

	TypeIntPtr:     &sql.NullInt64{},
	TypeInt64Ptr:   &sql.NullInt64{},
	TypeInt32Ptr:   &sql.NullInt32{},
	TypeFloat64Ptr: &sql.NullFloat64{},
	TypeStringPtr:  &sql.NullString{},
	TypeBoolPtr:    &sql.NullBool{},
	TypeTimePtr:    &sql.NullTime{},
}

type Nullable interface {
	Value() (driver.Value, error)
}

// 将SQL查询的数据动态转化成 model 对应的结构体
func DynamicScan(model interface{}, sqlStr string) ([]interface{}, error){
	db, err := sql.Open("mysql","root:@/butterfly?parseTime=true")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sqlStr)
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

func main(){
	type DmUser struct {
		ID      *int       `json:"id"`
		Name    *string    `json:"name"`
		Gender  *int64     `json:"gender"`
		Created *time.Time `json:"created"`
		Time    *time.Time `json:"time"`
	}

	// 动态 scan 结构体的数据
	slice, err := DynamicScan(DmUser{}, "select id, name, gender, created, time from user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("动态 scan 结构体的数据: %v \n", slice)

	slice, err = DynamicScan(DmUser{}, "select id, name, gender, created, time from user where 1 = 2")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("动态 scan 结构体的数据无数据: %v \n", slice)

	// 结构体指针
	slice, err = DynamicScan(&DmUser{}, "select id, name, gender, created, time from user")
	if err != nil{
		panic(err)
	}

	// 字段非指针类型
	type DmUser0 struct {
		ID      int       `json:"id"`
		Name    string    `json:"name"`
		Gender  int64     `json:"gender"`
		Created time.Time `json:"created"`
		Time    time.Time `json:"time"`
	}
	slice, err = DynamicScan(DmUser0{}, "select id, name, gender, created, time from user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("字段非指针类型: %v \n", slice)

	// 实际的字段类型不匹配
	type DmUser1 struct {
		ID      float64   `json:"id"`
		Name    string    `json:"name"`
		Gender  string    `json:"gender"`
		Created time.Time `json:"created"`
		Time    time.Time `json:"time"`
	}
	slice, err = DynamicScan(DmUser1{}, "select id, name, gender, created, time from user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("实际的字段类型不匹配: %v \n", slice)

	// SQL多余的列不解析
	type DmUser2 struct {
		ID      float64   `json:"id"`
		Name    string    `json:"name"`
		Gender  string    `json:"gender"`
	}
	slice, err = DynamicScan(DmUser2{}, "select id, name, gender, created, time from user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("SQL多余的列不解析: %v \n", slice)

	// SQL 的列不足结构体会使用默认值
	type DmUser3 struct {
		ID      int       `json:"id"`
		Name    string    `json:"name"`
		Gender  int64     `json:"gender"`
		Created time.Time `json:"created"`
		Time    time.Time `json:"time"`
	}
	slice, err = DynamicScan(DmUser3{}, "select id, name, gender from user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("SQL 的列不足结构体会使用默认值: %v \n", slice)

	// TODO 待覆盖所有的值类型 mysql 值类型
}