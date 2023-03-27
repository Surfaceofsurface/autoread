package model

import (
	"context"
	"time"
)

//	type Tasks struct {
//		PendingBook []Book
//		Creater     *TaskMaker
//	}
type Book struct {
	Title    string
	BookID   string
	FileType int
}
type TProcess struct {
	Meta             TProcessMeta
	StartPage        int
	NowAt            int
	NodeIndex        int
	ReadingStartTime time.Time
	SessionToken     string `json:"-"`
}
type TProcessMeta struct {
	CreateT time.Time
	Creater *TaskMaker `json:"-"`
	Target  *Book
	Catalog []CataLog `json:"-"`
	StartAt int
	Total   int
	Ctx     context.Context `json:"-"`
}
type CataLog struct {
	ParaIndex int
	Page      int
}

// 必须拥有Meta后才能调用
// func (t *TProcess) ProduceEgo() error {
// 	//创建请求环境
// 	_Env := global.R.GetEnv(t.Meta.Creater.Username)
// 	if _Env.G == nil {
// 		return errors.New("not found t.Meta.Creater.Username")
// 	}
// 	G := _Env.G
// 	// P := _Env.G

// 	/*
// 		这部分for循环是为了设置以下字段:
// 		res[i].StartAt 需要请求
// 		res[i].NowAt=res[i].StartAt
// 		res[i].CreateT = time.Now() 请求后
// 	*/

// 	switch t.Meta.Target.FileType {
// 	default:
// 		fmt.Printf("unknown v.Target.FileType: %v\n", t.Meta.Target.FileType)
// 		return errors.New("unknown v.Target.FileType")
// 	case 0:
// 		d, e := G(fmt.Sprintf(global.PROCESS_0, t.Meta.Target.BookID), nil)
// 		if e != nil {
// 			fmt.Printf("e: %v\n", e)
// 			return e
// 		}
// 		m := map[string]int{}
// 		e = json.Unmarshal(d, &m)
// 		fmt.Printf("map: %v\n", m)
// 		if e != nil {
// 			return e
// 		}
// 		t.Meta.StartAt = m["page"]
// 		t.NowAt = m["page"]

// 	case 3:
// 		d, e := G(fmt.Sprintf(global.PROCESS_3, t.Meta.Target.BookID), nil)
// 		if e != nil {
// 			fmt.Printf("e: %v\n", e)
// 			return e
// 		}
// 		m := map[string]int{}
// 		e = json.Unmarshal(d, &m)
// 		if e != nil {
// 			fmt.Printf("err: %v\n", e)
// 			return e
// 		}
// 		t.Meta.StartAt = m["paragraph"]
// 		t.NowAt = m["paragraph"]
// 	}

//		/*
//			这部分for循环是为了设置以下字段:
//			t.Meta.Catalog/Total 需要请求
//			res[i].CreateT = time.Now() 请求后
//		*/
//		d, e := G(fmt.Sprintf(global.G_CATALOG, t.Meta.Target.BookID), map[string]string{"filetype": strconv.Itoa(t.Meta.Target.FileType)})
//		if e != nil {
//			fmt.Printf("e: %v\n", e)
//			return e
//		}
//		var DTO dto.CataLogDTO
//		e = json.Unmarshal(d, &DTO)
//		if e != nil {
//			fmt.Printf("e: %v\n", e)
//			return e
//		}
//		t.Meta.Total = DTO.Paragraph
//		t.Meta.Catalog = fromCataLogDTO2Model(DTO.Data)
//		t.Meta.CreateT = time.Now()
//		return nil
//	}
// func fromCataLogDTO2Model(DTO []dto.InData) []CataLog {

// 	_res := make([]CataLog, 0)
// 	res := &_res
// 	fromCataLogDTO2ModelRecrusion(DTO, res)
// 	return *res
// }
// func fromCataLogDTO2ModelRecrusion(DTO []dto.InData, res *([]CataLog)) {
// 	if len(DTO) == 0 {
// 		return
// 	}
// 	for _, v := range DTO {
// 		_tmp := CataLog{
// 			ParaIndex: v.ParaIndex,
// 			Page:      v.Page,
// 		}
// 		*res = append(*res, _tmp)
// 		fromCataLogDTO2ModelRecrusion(v.Children, res)
// 	}
// }
