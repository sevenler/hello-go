package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"git.in.zhihu.com/zhihu/hello/utils"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"time"
)


//CREATE TABLE `user` (
//`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
//`name` varchar(36) DEFAULT NULL COMMENT 'name',
//`gender` int DEFAULT NULL COMMENT 'gender',
//`created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
//`updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//PRIMARY KEY (`id`)
//)
type DmUser struct {
	ID      *int       `json:"id"`
	Name    *string    `json:"name"`
	Gender  *int64     `json:"gender"`
	Created *time.Time `json:"created"`
	Time    *time.Time `json:"time"`
}


var ColumnNullableMap = map[reflect.Type]Nullable{
	reflect.TypeOf(int(0)): &sql.NullInt64{},
	reflect.TypeOf(int64(0)): &sql.NullInt64{},
	reflect.TypeOf(int32(0)): &sql.NullInt32{},
	reflect.TypeOf(float64(0)): &sql.NullFloat64{},
	reflect.TypeOf(""): &sql.NullString{},
	reflect.TypeOf(true): &sql.NullBool{},
	reflect.TypeOf(time.Now()): &sql.NullTime{},

	reflect.TypeOf(utils.IntPtr(0)): &sql.NullInt64{},
	reflect.TypeOf(utils.Int64Ptr(0)): &sql.NullInt64{},
	reflect.TypeOf(utils.Int32Ptr(0)): &sql.NullInt32{},
	reflect.TypeOf(utils.Float64Ptr(0)): &sql.NullFloat64{},
	reflect.TypeOf(utils.StringPtr("")): &sql.NullString{},
	reflect.TypeOf(utils.BoolPtr(true)): &sql.NullBool{},
	reflect.TypeOf(utils.TimePtr(time.Now())): &sql.NullTime{},
}

type Nullable interface {
	Value() (driver.Value, error)
}

func DynamicScan(model interface{}, sqlStr string) ([]interface{}, error){
	db, err := sql.Open("mysql","root:@/butterfly?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := db.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	columns, err := rows.Columns()
	if err != nil{
		panic("error columns")
	}

	t := reflect.TypeOf(model)
	tagFiled := map[string]reflect.StructField{}
	for i := 0; i< t.NumField();i++{
		field := t.Field(i)
		v := field.Tag.Get("json")
		tagFiled[v] = field
	}

	// sql.Nullxxx 是为了处理 db 中的值是 Null 的情况
	// 所有的 Nullxxx 是都可以调用 Value， Value 返回空表示是一个 Null 的对象
	colPtr0 := make([]interface{}, 0)
	for i := 0; i< len(columns);i++{
		name := columns[i]
		field := tagFiled[name]
		nullable, ok := ColumnNullableMap[field.Type]
		if !ok{
			return nil, fmt.Errorf("not support type %v of field %s", field.Type, name)
		}
		colPtr0 = append(colPtr0, nullable)
	}

	var slice = make([]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(colPtr0...)
		if err != nil {
			log.Fatal(err)
		}
		var obj = reflect.New(t).Elem()
		for i, col := range colPtr0 {
			rowValue, err := col.(Nullable).Value()
			filedName := columns[i]
			field := obj.FieldByName(tagFiled[filedName].Name)
			if err == nil && rowValue != nil{
				switch field.Type() {
				case reflect.TypeOf(int(0)):
					field.SetInt(rowValue.(int64))
				case reflect.TypeOf(int64(0)):
					field.SetInt(rowValue.(int64))
				case reflect.TypeOf(int32(0)):
					field.SetInt(rowValue.(int64))
				case reflect.TypeOf(float64(0)):
					field.SetFloat(rowValue.(float64))
				case reflect.TypeOf(""):
					field.SetString(rowValue.(string))
				case reflect.TypeOf(true):
					field.SetBool(rowValue.(bool))
				case reflect.TypeOf(time.Now()):
					t := rowValue.(time.Time)
					field.Set(reflect.ValueOf(t))

				case reflect.TypeOf(utils.IntPtr(0)):
					field.Set(reflect.ValueOf(utils.IntPtr(int(rowValue.(int64)))))
				case reflect.TypeOf(utils.Int64Ptr(0)):
					field.Set(reflect.ValueOf(utils.Int64Ptr(rowValue.(int64))))
				case reflect.TypeOf(utils.Int32Ptr(0)):
					field.Set(reflect.ValueOf(utils.Int32Ptr(rowValue.(int32))))
				case reflect.TypeOf(utils.Float64Ptr(0)):
					field.Set(reflect.ValueOf(utils.Float64Ptr(rowValue.(float64))))
				case reflect.TypeOf(utils.StringPtr("")):
					field.Set(reflect.ValueOf(utils.StringPtr(rowValue.(string))))
				case reflect.TypeOf(utils.BoolPtr(true)):
					field.Set(reflect.ValueOf(utils.BoolPtr(rowValue.(bool))))
				case reflect.TypeOf(utils.TimePtr(time.Now())):
					t := rowValue.(time.Time)
					field.Set(reflect.ValueOf(utils.TimePtr(t)))

				default:
					return nil, fmt.Errorf("error set row %s got type %v",
						filedName, field.Type())
				}
			}
		}
		slice = append(slice, obj.Interface())
	}
	return slice, err
}

func main(){
	// 动态 scan 结构体的数据
	slice, err := DynamicScan(DmUser{},
		"select id, name, gender, created, time from user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("%v \n", slice)
}