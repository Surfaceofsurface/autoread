package service

import (
	"context"
	"crypto"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"private/autoread/reactor/global"
	"private/autoread/reactor/model"
	"private/autoread/reactor/service/dto"
	"strconv"
	"strings"
	"time"
)

type TasksS struct {
	TP *model.TProcess
	WithToken
}

func (t TasksS) Ctx() context.Context {
	return context.Background()
}
func (t TasksS) Uid() string {
	return t.TP.Meta.Target.Title
}

func (t TasksS) Init() error {
	d, e := t.G(fmt.Sprintf(global.G_BKDETAIL, t.TP.Meta.Target.BookID), map[string]string{"pinst": t.TP.Meta.Creater.SchoolID})
	if e != nil {
		fmt.Printf("e: %v\n", e)
		return e
	}
	m := map[string]any{}
	e = json.Unmarshal(d, &m)
	if e != nil {
		fmt.Printf("err: %v\n", e)
		return e
	}
	t.TP.Meta.Target.FileType = int(m["fileType"].(float64))

	switch t.TP.Meta.Target.FileType {
	default:
		fmt.Printf("unknown v.Target.FileType: %v\n", t.TP.Meta.Target.FileType)
		return errors.New("unknown v.Target.FileType")
	case 0:
		d, e := t.G(fmt.Sprintf(global.PROCESS_0, t.TP.Meta.Target.BookID), nil)
		if e != nil {
			fmt.Printf("e: %v\n", e)
			return e
		}
		m := map[string]int{}
		e = json.Unmarshal(d, &m)
		if e != nil {
			return e
		}
		t.TP.Meta.StartAt = m["page"]
		t.TP.StartPage = m["page"]
		t.TP.NowAt = m["page"] + 1

	case 3:
		d, e := t.G(fmt.Sprintf(global.PROCESS_3, t.TP.Meta.Target.BookID), nil)
		if e != nil {
			fmt.Printf("e: %v\n", e)
			return e
		}
		m := map[string]int{}
		e = json.Unmarshal(d, &m)
		if e != nil {
			fmt.Printf("err: %v\n", e)
			return e
		}
		t.TP.Meta.StartAt = m["paragraph"]
		t.TP.NowAt = m["paragraph"] + 1
		t.TP.StartPage = m["paragraph"]
	}

	d, e = t.G(fmt.Sprintf(global.G_CATALOG, t.TP.Meta.Target.BookID), map[string]string{"filetype": strconv.Itoa(t.TP.Meta.Target.FileType)})
	if e != nil {
		fmt.Printf("e: %v\n", e)
		return e
	}

	var DTO dto.CataLogDTO
	e = json.Unmarshal(d, &DTO)
	if e != nil {
		fmt.Printf("e: %v\n", e)
		return e
	}
	t.TP.Meta.Catalog = fromCataLogDTO2Model(DTO.Data)

	switch t.TP.Meta.Target.FileType {
	case 0:
		t.TP.Meta.Total = t.TP.Meta.Catalog[len(t.TP.Meta.Catalog)-1].Page

	case 3:
		t.TP.Meta.Total = DTO.Paragraph
	}
	if t.TP.NowAt >= t.TP.Meta.Total {
		t.TP.NowAt = 2
	}
	t.TP.Meta.Catalog = fromCataLogDTO2Model(DTO.Data)
	t.TP.ReadingStartTime = time.Now()
	t.TP.Meta.CreateT = time.Now()
	t.TP.NodeIndex = binarySearch(t.TP.Meta.Catalog, t.TP.NowAt, t.TP.Meta.Target.FileType)
	return nil
}

// func (t TasksS) Init() error {
// 	d, e := t.G(fmt.Sprintf(global.G_BKDETAIL, t.TP.Meta.Target.BookID), map[string]string{"pinst": t.TP.Meta.Creater.SchoolID})
// 	if e != nil {
// 		fmt.Printf("e: %v\n", e)
// 		return e
// 	}
// 	m := map[string]any{}
// 	e = json.Unmarshal(d, &m)
// 	if e != nil {
// 		fmt.Printf("err: %v\n", e)
// 		return e
// 	}
// 	t.TP.Meta.Target.FileType = int(m["fileType"].(float64))

// 	d, e = t.G(fmt.Sprintf(global.G_CATALOG, t.TP.Meta.Target.BookID), map[string]string{"filetype": strconv.Itoa(t.TP.Meta.Target.FileType)})
// 	if e != nil {
// 		fmt.Printf("e: %v\n", e)
// 		return e
// 	}

// 	var DTO dto.CataLogDTO
// 	e = json.Unmarshal(d, &DTO)
// 	if e != nil {
// 		fmt.Printf("e: %v\n", e)
// 		return e
// 	}
// 	t.TP.Meta.Catalog = fromCataLogDTO2Model(DTO.Data)

//		switch t.TP.Meta.Target.FileType {
//		case 0:
//			t.TP.Meta.Total = t.TP.Meta.Catalog[len(t.TP.Meta.Catalog)-1].Page
//		case 3:
//			t.TP.Meta.Total = DTO.Paragraph
//		}
//		t.TP.NowAt = 2
//		t.TP.StartPage = 1
//		t.TP.Meta.Catalog = fromCataLogDTO2Model(DTO.Data)
//		t.TP.ReadingStartTime = time.Now()
//		t.TP.Meta.CreateT = time.Now()
//		t.TP.NodeIndex = 0
//		// t.TP.NodeIndex = binarySearch(t.TP.Meta.Catalog, t.TP.NowAt, t.TP.Meta.Target.FileType)
//		return nil
//	}
func (t TasksS) Do() (bool, error) {
	switch t.TP.Meta.Target.FileType {
	case 0:
		//结束的信号
		if t.TP.NodeIndex+1 >= len(t.TP.Meta.Catalog)-1 {
			return true, nil
		}
		if t.TP.NowAt > t.TP.Meta.Total {
			return true, nil
		}

		return false, t.type0do()
	case 3:
		//结束的信号
		if t.TP.NodeIndex+1 >= len(t.TP.Meta.Catalog)-1 {
			return true, nil
		}
		if t.TP.NowAt > t.TP.Meta.Total {
			return true, nil
		}
		return false, t.type3do()
	}

	return false, nil
}
func (t TasksS) Fin() {
	t.TP.NodeIndex = len(t.TP.Meta.Catalog)
}
func (t TasksS) NowAt() float32 {
	if float32(t.TP.Meta.Total) == 0 {
		return 0
	}
	return float32(t.TP.NowAt) / float32(t.TP.Meta.Total)
}
func (t *TasksS) type0do() error {
	var readEndTime time.Time
	if t.TP.Meta.Catalog[t.TP.NodeIndex+1].Page == t.TP.Meta.Catalog[t.TP.NodeIndex].Page {
		t.TP.NodeIndex++
		return nil
	}
	t.type0WithToken()
	tokenMap := t.AdtionalMap(map[string]string{
		"token": t.TP.SessionToken,
	})
	if t.TP.SessionToken == "" {
		delete(tokenMap, "token")
	}

	bytearr, e := t.G(fmt.Sprintf(global.PROCESS_N, t.TP.Meta.Target.BookID), tokenMap)
	if e != nil {
		fmt.Printf("e1: %v\n", e)
		return e
	}
	readdto := dto.ReadGetDTO{}
	e = json.Unmarshal(bytearr, &readdto)
	if e != nil {
		fmt.Printf("e2: %v\n", e)
		return e
	}
	_, e = t.G(readdto.Content, t.AdtionalMap(map[string]string{
		"page": fmt.Sprintf("%d", t.TP.NowAt+1),
	}))
	if e != nil {
		fmt.Printf("e3: %v\n", e)
		return e
	}
	endpage := t.TP.Meta.Catalog[t.TP.NodeIndex+1].Page

	//只有达到每节的末尾时,才PostTrack
	if endpage == t.TP.NowAt {
		readEndTime = time.Now()
		bodymap := map[string]any{
			"networktype": "online",
			"lists": [1]map[string]any{
				{
					"startpage":        t.TP.StartPage,
					"endpage":          t.TP.NowAt - 1,
					"bookruid":         t.TP.Meta.Target.BookID,
					"filetype":         t.TP.Meta.Target.FileType,
					"pageno":           t.TP.Meta.Catalog[t.TP.NodeIndex].Page,
					"readingstarttime": t.TP.ReadingStartTime.Format("2006-01-02 15:04:05"),
					"readingendtime":   readEndTime.Format("2006-01-02 15:04:05"),
					"readingcount":     t.TP.NowAt - t.TP.StartPage,
					"typecode":         "ebook",
				},
			},
		}
		body, e := json.Marshal(bodymap)
		if e != nil {
			return e
		}
		_, e = t.P(global.P_TRACK, body)
		if e != nil {
			return e
		}

		t.TP.StartPage = t.TP.NowAt
		t.TP.ReadingStartTime = readEndTime
		t.TP.NodeIndex++
	}
	progressmap := map[string]any{
		"charIndex": 0,
		"page":      float32(t.TP.NowAt),
		"paragraph": 0,
		"percent":   100 * float32(t.TP.NowAt-1) / float32(t.TP.Meta.Total),
	}
	b, e := json.Marshal(progressmap)
	if e != nil {
		return e
	}

	_, e = t.P(fmt.Sprintf(global.PROCESS_0, t.TP.Meta.Target.BookID), b)
	if e != nil {
		fmt.Printf("e: %v\n", e)

		return e
	}
	t.TP.NowAt++

	return nil
}
func (t *TasksS) type3do() error {
	if t.TP.Meta.Catalog[t.TP.NodeIndex+1].ParaIndex == t.TP.Meta.Catalog[t.TP.NodeIndex].ParaIndex {
		t.TP.NodeIndex++
		return nil
	}
	var readEndTime time.Time

	endpage := t.TP.Meta.Catalog[t.TP.NodeIndex+1].ParaIndex
	//只有达到每节的末尾时,才PostTrack
	if endpage == t.TP.NowAt {
		readEndTime = time.Now()
		bodymap := map[string]any{
			"networktype": "online",
			"lists": [1]map[string]any{
				{
					"startpage":        t.TP.StartPage,
					"endpage":          endpage,
					"bookruid":         t.TP.Meta.Target.BookID,
					"filetype":         2,
					"pageno":           t.TP.Meta.Catalog[t.TP.NodeIndex].ParaIndex,
					"readingstarttime": t.TP.ReadingStartTime.Format("2006-01-02 15:04:05"),
					"readingendtime":   readEndTime.Format("2006-01-02 15:04:05"),
					"readingcount":     t.TP.NowAt - t.TP.StartPage,
					"typecode":         "ebook",
				},
			},
		}
		body, e := json.Marshal(bodymap)
		if e != nil {
			return e
		}
		_, e = t.P(global.P_TRACK, body)
		if e != nil {
			return e
		}
		t.TP.StartPage = t.TP.NowAt
		t.TP.ReadingStartTime = readEndTime
		t.TP.NodeIndex++
	}
	progressmap := map[string]any{
		"charIndex": 7,
		"page":      0,
		"paragraph": float32(t.TP.NowAt),
		"percent":   100 * float32(t.TP.NowAt) / float32(t.TP.Meta.Total),
	}
	b, e := json.Marshal(progressmap)
	if e != nil {
		return e
	}
	_, e = t.P(fmt.Sprintf(global.PROCESS_3, t.TP.Meta.Target.BookID), b)
	if e != nil {
		return e
	}
	t.TP.NowAt++
	if t.TP.NowAt > t.TP.Meta.Total {
		return nil
	}
	return nil
}
func fromCataLogDTO2Model(DTO []dto.InData) []model.CataLog {

	_res := make([]model.CataLog, 0)
	res := &_res
	fromCataLogDTO2ModelRecrusion(DTO, res)
	return *res
}
func fromCataLogDTO2ModelRecrusion(DTO []dto.InData, res *([]model.CataLog)) {
	if len(DTO) == 0 {
		return
	}

	for _, v := range DTO {
		_tmp := model.CataLog{
			ParaIndex: v.ParaIndex,
			Page:      v.Page,
		}
		if len(*res) == 0 {
			*res = append(*res, _tmp)
		}
		if _tmp.Page > (*res)[len(*res)-1].Page {
			*res = append(*res, _tmp)
		}
		fromCataLogDTO2ModelRecrusion(v.Children, res)
	}
}

func binarySearch(arr []model.CataLog, who int, fileType int) int {
	maxi := len(arr) - 1
	mini := 0
	midi := (maxi + mini) / 2
	{
		switch fileType {
		default:
			return -1
		case 0:
			for i := 0; maxi-mini > 1 && i < len(arr); i++ {
				if who == arr[midi].Page {
					break
				}
				if who < arr[midi].Page {
					maxi = midi
				} else {
					mini = midi
				}
				midi = (maxi + mini) / 2
			}
		case 3:
			for i := 0; maxi-mini > 1 && i < len(arr); i++ {
				if who == arr[midi].ParaIndex {
					break
				}
				if who < arr[midi].ParaIndex {
					maxi = midi
				} else {
					mini = midi
				}
				midi = (maxi + mini) / 2
			}
		}
		return midi

	}
}
func (t *TasksS) AdtionalMap(opt map[string]string) map[string]string {

	res := map[string]string{
		"nonce": GenNonce(),
		"stime": strconv.FormatInt(time.Now().Unix(), 10),
		"pinst": t.TP.Meta.Creater.SchoolID,
		"page":  fmt.Sprintf("%d", t.TP.NowAt-1),
	}
	for k, v := range opt {
		res[k] = v
	}
	hash := crypto.MD5.New()
	hash.Write([]byte("123456" + res["nonce"] + res["stime"]))
	res["sign"] = strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))

	return res
}
func (t *TasksS) type0WithToken() error {
	if t.TP.SessionToken == "" {
		statearr, e := t.G(fmt.Sprintf(global.G_STATE, t.TP.Meta.Target.BookID), map[string]string{
			"now":     strconv.FormatInt(time.Now().UnixMilli(), 10),
			"ipToken": "null",
			"pinst":   t.TP.Meta.Creater.SchoolID,
		})
		if e != nil {
			return e
		}
		state := dto.StateDTO{}
		json.Unmarshal(statearr, &state)
		bytearr, e := t.G(fmt.Sprintf(global.PROCESS_N, t.TP.Meta.Target.BookID), t.AdtionalMap(map[string]string{
			"from":  "default",
			"stime": fmt.Sprintf("%d", int64(math.Floor(float64(state.Time)/1000))),
		}))
		if e != nil {
			return e
		}
		tokenMap := map[string]string{}
		json.Unmarshal(bytearr, &tokenMap)
		t.TP.SessionToken = tokenMap["token"]

		_, e = t.G(tokenMap["content"], t.AdtionalMap(nil))

		if e != nil {
			return e
		}
	}

	return nil
}

func GenNonce() string {
	s := "xxxxxxxx-xxxx-4xxx-xxxx-xxxxxxxxxxxx"
	for range s {
		s = strings.Replace(s, "x", fmt.Sprintf("%x", rand.Int31n(16)), 1)
	}
	return s
}
