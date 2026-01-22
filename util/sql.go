// sql语句拼接工具类
package util

import (
	"regexp"
	"strconv"
	"strings"
)

type SqlCompletion struct {
	initSql      string
	initCountSql string
	whereSql     strings.Builder
	groupSql     string
	havingSql    string
	orderSql     strings.Builder
	limitSql     string
	paramIndex   int
	whereParams  []any
	limitParams  []any
}

// 设置初始sql语句
func (s *SqlCompletion) InitSql(sql string) *SqlCompletion {
	s.initSql = sql
	r, _ := regexp.Compile("(?i)(?s)select(.*?)from")
	s.initCountSql = r.ReplaceAllString(sql, "select count(*) as count from")
	return s
}

// 设置初始sql语句和总行数语句
func (s *SqlCompletion) InitSqlAndCount(sql, countSql string) *SqlCompletion {
	s.initSql = sql
	s.initCountSql = countSql
	return s
}

// 获取sql语句
func (s *SqlCompletion) GetSql() string {
	return s.initSql + s.whereSql.String() + s.groupSql + s.havingSql + s.orderSql.String() + s.limitSql
}

// 获取行数sql语句
func (s *SqlCompletion) GetCountSql() string {
	return s.initCountSql + s.whereSql.String() + s.groupSql + s.havingSql
}

// 获取条件列表
func (s *SqlCompletion) GetParams() []any {
	return append(s.whereParams, s.limitParams...)
}

// 获取行数条件列表
func (s *SqlCompletion) GetCountParams() []any {
	return s.whereParams
}

// =
func (s *SqlCompletion) Eq(field string, param any) *SqlCompletion {
	s.where(field, param, "=", false)
	return s
}

// !=
func (s *SqlCompletion) Ne(field string, param any) *SqlCompletion {
	s.where(field, param, "!=", false)
	return s
}

// >
func (s *SqlCompletion) Gt(field string, param any) *SqlCompletion {
	s.where(field, param, ">", false)
	return s
}

// <
func (s *SqlCompletion) Lt(field string, param any) *SqlCompletion {
	s.where(field, param, "<", false)
	return s
}

// >=
func (s *SqlCompletion) Ge(field string, param any) *SqlCompletion {
	s.where(field, param, ">=", false)
	return s
}

// <=
func (s *SqlCompletion) Le(field string, param any) *SqlCompletion {
	s.where(field, param, "<=", false)
	return s
}

// like
func (s *SqlCompletion) Like(field string, param string) *SqlCompletion {
	if param != "" {
		s.where(field, param, "like", true)
	}
	return s
}

// in
func (s *SqlCompletion) In(field string, params []string) *SqlCompletion {
	s.whereIn(field, params, "in")
	return s
}

// not in
func (s *SqlCompletion) NotIn(field string, params []string) *SqlCompletion {
	s.whereIn(field, params, "not in")
	return s
}

// is null
func (s *SqlCompletion) IsNull(field string) *SqlCompletion {
	s.whereHyphen()
	s.whereSql.WriteString(field)
	s.whereSql.WriteString(" is null")
	return s
}

// is not null
func (s *SqlCompletion) IsNotNull(field string) *SqlCompletion {
	s.whereHyphen()
	s.whereSql.WriteString(field)
	s.whereSql.WriteString(" is not null")
	return s
}

// = or，用以单输入框筛选全部字段，如调用只能用or相关条件，数据量大会有性能问题
func (s *SqlCompletion) EqOr(field string, param any) *SqlCompletion {
	s.whereOr(field, param, "=", false)
	return s
}

// like or，用以单输入框筛选全部字段，如调用只能用or相关条件，数据量大会有性能问题
func (s *SqlCompletion) LikeOr(fields []string, param string) *SqlCompletion {
	if param != "" {
		for _, field := range fields {
			s.whereOr(field, param, "like", true)
		}
	}
	return s
}

// 分组，只设置最后一次
func (s *SqlCompletion) Group(fields string) *SqlCompletion {
	s.groupSql = " group by " + fields
	return s
}

// having条件，只设置最后一次
func (s *SqlCompletion) Having(fields string) *SqlCompletion {
	s.havingSql = " having " + fields
	return s
}

// 排序
func (s *SqlCompletion) Order(field string, isAsc bool) *SqlCompletion {
	needComma := true
	if s.orderSql.Len() == 0 {
		s.orderSql.WriteString(" order by ")
		needComma = false
	}
	if needComma {
		s.orderSql.WriteString(",")
	}
	s.orderSql.WriteString(field)
	if !isAsc {
		s.orderSql.WriteString(" desc")
	}
	return s
}

// 分页，只设置最后一次，需最后调用
func (s *SqlCompletion) Limit(page, size int) *SqlCompletion {
	if page > 0 && size > 0 {
		offset := size * (page - 1)
		s.paramIndex = s.paramIndex + 1
		paramPlaceholder1 := "$" + strconv.Itoa(s.paramIndex)
		s.paramIndex = s.paramIndex + 1
		paramPlaceholder2 := "$" + strconv.Itoa(s.paramIndex)
		s.limitSql = " limit " + paramPlaceholder1 + " offset " + paramPlaceholder2
		s.limitParams = []any{size, offset}
	}
	return s
}

// 添加where条件
func (s *SqlCompletion) where(field string, param any, symbol string, isLike bool) {
	s.paramIndex = s.paramIndex + 1
	paramPlaceholder := "$" + strconv.Itoa(s.paramIndex)
	s.whereHyphen()
	s.whereSql.WriteString(field)
	s.whereSql.WriteString(" ")
	s.whereSql.WriteString(symbol)
	if isLike {
		s.whereSql.WriteString(" '%'||" + paramPlaceholder + "||'%'")
	} else {
		s.whereSql.WriteString(" " + paramPlaceholder)
	}
	s.whereParams = append(s.whereParams, param)
}

// 添加where条件（in / not in）
func (s *SqlCompletion) whereIn(field string, params []string, symbol string) {
	if len(params) == 0 {
		return
	}
	s.whereHyphen()
	s.whereSql.WriteString(field)
	s.whereSql.WriteString(" ")
	s.whereSql.WriteString(symbol)
	s.whereSql.WriteString(" (")
	for i, v := range params {
		if i != 0 {
			s.whereSql.WriteString(",")
		}
		s.paramIndex = s.paramIndex + 1
		paramPlaceholder := "$" + strconv.Itoa(s.paramIndex)
		s.whereSql.WriteString(paramPlaceholder)
		s.whereParams = append(s.whereParams, v)
	}
	s.whereSql.WriteString(")")
}

// 添加where or条件
func (s *SqlCompletion) whereOr(field string, param any, symbol string, isLike bool) {
	s.paramIndex = s.paramIndex + 1
	paramPlaceholder := "$" + strconv.Itoa(s.paramIndex)
	s.whereHyphenOr()
	s.whereSql.WriteString(field)
	s.whereSql.WriteString(" ")
	s.whereSql.WriteString(symbol)
	if isLike {
		s.whereSql.WriteString(" '%'||" + paramPlaceholder + "||'%'")
	} else {
		s.whereSql.WriteString(" " + paramPlaceholder)
	}
	s.whereParams = append(s.whereParams, param)
}

// where条件，拼接where和and
func (s *SqlCompletion) whereHyphen() {
	needHyphen := true
	if s.whereSql.Len() == 0 {
		s.whereSql.WriteString(" where ")
		needHyphen = false
	}
	if needHyphen {
		s.whereSql.WriteString(" and ")
	}
}

// where条件，拼接where和or
func (s *SqlCompletion) whereHyphenOr() {
	needHyphen := true
	if s.whereSql.Len() == 0 {
		s.whereSql.WriteString(" where ")
		needHyphen = false
	}
	if needHyphen {
		s.whereSql.WriteString(" or ")
	}
}
