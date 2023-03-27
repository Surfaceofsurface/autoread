package depman

import (
	"reflect"
)

func Equip[T any](obj T) T {
	//思路：利用递归遍历结构体
	//获得空接口对应的类型,因为后续要修改该对象,所以obj必须以指针的方式传入
	_typ := reflect.TypeOf(obj)
	if _typ.Kind() != reflect.Pointer {
		return obj
	}
	//传入的不是结构体则不需要注入依赖
	_typ = _typ.Elem()
	if _typ.Kind() != reflect.Struct {
		return obj
	}

	//Elem获取的是空接口指针里所指向的值
	_val := reflect.ValueOf(obj).Elem()
	//开始依赖注入
	equipStruct(_val, _typ)
	return obj
}

func equipField(v reflect.Value, t reflect.StructField) {
	//判断field是否在container里
	if item, ok := GetItem(t.Name); ok && item != nil {
		//item不能是空指针,否则报错
		v.Set(reflect.ValueOf(item))
		return
	}

	if t.Type.Kind() == reflect.Struct {
		equipStruct(v, t.Type)
		return
	}
	/*如果要装载的字段是指针的话,就会在内部New一个指针所指向的对象出来
	然后装配New出来的对象后再将结构体字段的指针指向New出来的对象
	*/
	if t.Type.Kind() == reflect.Pointer {
		ptoT := t.Type.Elem()
		if ptoT.Kind() != reflect.Struct {
			return
		}
		//利用反射创建新对象,new出来的是指针
		ptoVP := reflect.New(ptoT)
		//获取指针所指对象的真实值
		ptoV := reflect.Indirect(ptoVP)
		//装配该真实值
		equipStruct(ptoV, ptoT)
		//装配完成后,修改原来的对象指针,使其指向刚刚装配的新对象
		v.Set(ptoVP)
	}
}
func equipStruct(v reflect.Value, t reflect.Type) {
	_num := t.NumField()
	for i := 0; i < _num; i++ {
		equipField(v.Field(i), t.Field(i))
	}
}
