package logging

import (
	"fmt"
	"reflect"

	"github.com/anacrolix/missinggo/iter"
)

type Msg struct {
	Format  string
	Args    []interface{}
	Printer func(Msg) string
	Skip_   int
	Values_ []interface{}
	Level   Level
	Name    string
}

func (me Msg) String() string {
	return me.Text()
}

func fmsgPrinter(m Msg) string {
	return fmt.Sprintf(m.Format, m.Args...)
}

func Fmsg(format string, a ...interface{}) Msg {
	return Msg{
		Format:  format,
		Args:    a,
		Printer: fmsgPrinter,
	}
}

var Fstr = Fmsg

func strPrinter(m Msg) string {
	return m.Format
}

func Str(s string) (m Msg) {
	return Msg{
		Format:  s,
		Printer: strPrinter,
	}
}

func (m Msg) Skip(skip int) Msg {
	m.Skip_ += skip
	return m
}

type item struct {
	key, value interface{}
}

// rename sink
func (m Msg) Log(l *Logger) Msg {
	if l.IsEnabledFor(m.Level) {
		l.Handle(m.Skip(1))
	}
	return m
}

func (m Msg) SinkNew(l *Logger) Msg {
	l.Handle(m)
	return m
}

// TODO: What ordering should be applied to the values here, per MsgImpl.Values. For now they're
// traversed in order of the slice.
func (m Msg) WithValues(v ...interface{}) Msg {
	m.Values_ = append(m.Values_, v...)
	return m
}

func (m Msg) AddValues(v ...interface{}) Msg {
	return m.WithValues(v...)
}

func (m Msg) With(key, value interface{}) Msg {
	return m.WithValues(item{key, value})
}

func (m Msg) Add(key, value interface{}) Msg {
	return m.With(key, value)
}

func (m Msg) SetLevel(level Level) Msg {
	return m.With(levelKey, level)
}

func (m Msg) GetByKey(key interface{}) (value interface{}, ok bool) {
	m.Values(func(i interface{}) bool {
		if keyValue, isKeyValue := i.(item); isKeyValue && keyValue.key == key {
			value = keyValue.value
			ok = true
		}
		return !ok
	})
	return
}

func (m Msg) GetLevel() (l Level, ok bool) {
	v, ok := m.GetByKey(levelKey)
	if ok {
		l = v.(Level)
	}
	return
}

func (m Msg) HasValue(v interface{}) (has bool) {
	m.Values(func(i interface{}) bool {
		if i == v {
			has = true
		}
		return !has
	})
	return
}

func (m Msg) AddValue(v interface{}) Msg {
	return m.AddValues(v)
}

func (m Msg) GetValueByType(p interface{}) bool {
	pve := reflect.ValueOf(p).Elem()
	t := pve.Type()
	return !iter.All(func(i interface{}) bool {
		iv := reflect.ValueOf(i)
		if iv.Type() == t {
			pve.Set(iv)
			return false
		}
		return true
	}, m.Values)
}

func (m Msg) WithText(f func(Msg) string) Msg {
	m_ := m
	m_.Printer = func(Msg) string {
		return m.Printer(m)
	}
	return m_
}

func (m Msg) withName(name string) Msg {
	m.Name = name
	return m
}
