package action

import (
	"fmt"
	"reflect"
)

type DBLimit struct {
	maxParams  int // 最大参数数量（0 = 不限制）
	maxRowSize int // 单次插入的包大小限制 ( 0 = 不限制)
}

// mysql 默认4MB
// postgres 默认1GB
// Clickhouse 默认不限制，也按1GB算
var DBLimits = map[string]DBLimit{
	"postgres":   {maxParams: 65535, maxRowSize: 4190000},
	"mysql":      {maxParams: 65535, maxRowSize: 1073740000},
	"clickhouse": {maxParams: 0, maxRowSize: 1073740000},
}

// 计算安全批大小
func calcBatchCount(rows interface{}, dbType string) (int, error) {
	rv := reflect.ValueOf(rows)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		return 0, fmt.Errorf("rows must be a slice, got %T", rows)
	}
	if rv.Len() == 0 {
		return 0, fmt.Errorf("rows is empty")
	}

	elemType := rv.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	var columnCount int
	switch elemType.Kind() {
	case reflect.Struct:
		columnCount = countGormColumns(elemType)
	case reflect.Map:
		columnCount = rv.Index(0).Len()
	default:
		return 0, fmt.Errorf("unsupported element type: %s", elemType.Kind())
	}

	limit, _ := DBLimits[dbType]
	b := int(float64(limit.maxParams/columnCount) * 0.95)
	return b, nil
}

func countGormColumns(t reflect.Type) int {
	count := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}

		gormTag := field.Tag.Get("gorm")

		// 跳过 gorm:"-" 标记的字段
		if gormTag == "-" {
			continue
		}

		// 处理嵌套结构体（如 gorm.Model）
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct && field.Anonymous {
			// 不是时间类型的嵌套结构体，递归计算
			if fieldType.String() != "time.Time" {
				count += countGormColumns(fieldType)
				continue
			}
		}

		count++
	}
	return count
}
