package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"private/autoread/myutils"
	"private/autoread/reactor/controler"
	"private/autoread/reactor/global"
	"private/autoread/reactor/model"
	"private/autoread/reactor/service/dto"
	"sync"
)

type TaskMakerS struct {
	TMaker *model.TaskMaker
	WithToken
}

func (s *TaskMakerS) Login(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		loginJson := dto.LoginReq{
			Username: s.TMaker.Username,
			Password: s.TMaker.Password,
			Pinst:    s.TMaker.SchoolID,
		}
		data, err := json.Marshal(loginJson)
		if err != nil {
			return err
		}

		bytearr, err := s.P(global.LOGIN, data)
		if err != nil {
			return err
		}

		loginRes := &dto.LoginRes{}
		err = json.Unmarshal(bytearr, loginRes)
		if err != nil {
			return err
		}
		if loginRes.Code != "1" {
			return errors.New(loginRes.Text)
		}
		s.Token = loginRes.Token

		//

		if err != nil {
			return err
		}
		return nil
	}

}

func (s *TaskMakerS) GetPendingBooks() error {
	ctx, cancelCall := context.WithCancelCause(context.Background())
	wp := sync.WaitGroup{}
	wp.Add(2)
	res := dto.GBookPendingRes{}
	//读完的书目放到resf里,result_finish
	resf := dto.GBookPendingRes{}
	go func(ctx context.Context) {
		defer wp.Done()
		//同步正在读的书目（正在读的书包括已经读的书）
		// 参数形式 : ?page=1&size=10&state=1
		d, err := s.G(global.G_READING, map[string]string{
			"page":  "1",
			"size":  "32",
			"state": "1", //正在读state=1;读完state=0
		})
		if err != nil {
			fmt.Printf("err: %v\n", err)
			cancelCall(err)
			return
		}
		err = json.Unmarshal(d, &res)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			cancelCall(err)
			return
		}

	}(ctx)
	go func(ctx context.Context) {
		defer wp.Done()
		d, err := s.G(global.G_READING, map[string]string{
			"page":  "1",
			"size":  "32",
			"state": "0", //正在读state=1;读完state=0
		})
		if err != nil {
			fmt.Printf("err: %v\n", err)
			cancelCall(err)
			return
		}
		err = json.Unmarshal(d, &resf)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			cancelCall(err)
			return
		}

	}(ctx)
	wp.Wait()
	//获取已读的书目后去重
	margin := myutils.Diff(res.Data, resf.Data)
	s.TMaker.PendingBook = s.fromPendingToModel(margin)
	return ctx.Err()
}

func (s *TaskMakerS) Books2Tasks(bids []string) ([]controler.ExecAble, error) {
	res := make([]controler.ExecAble, len(bids))
	for i, id := range bids {
		fi := s.find(id)
		if fi == -1 {
			fmt.Printf("\"no book in tasks\": %v\n", "no book in tasks")
			return nil, errors.New("no book in tasks")
		}
		res[i] = &TasksS{

			TP: &model.TProcess{
				Meta: model.TProcessMeta{
					Ctx:     context.Background(),
					Target:  &s.TMaker.PendingBook[fi],
					Creater: s.TMaker,
				},
			},
			WithToken: s.WithToken,
		}
		// res[i].P = s.P
		// res[i].G = s.G
		// res[i].TP.Meta.Ctx = context.Background()
		// res[i].TP.Meta.Target = &s.TMaker.PendingBook[fi]
		// res[i].TP.Meta.Creater = s.TMaker
	}
	return res, nil
}

// 为了实现Task和TaskMaker解耦
func (s *TaskMakerS) DeriveTask(bookIds ...string) ([]model.TProcess, chan error) {
	//&&**这些都是用来并发控制的
	var waitGp sync.WaitGroup
	errc := make(chan error, 1)
	waitGp.Add(len(bookIds))
	//&&**这些都是用来并发控制的

	res := make([]model.TProcess, len(bookIds))

	//这部分for循环是为了请求获得fileType
	{
		for i, id := range bookIds {
			fi := s.find(id)
			if fi == -1 {
				fmt.Printf("\"no book in tasks\": %v\n", "no book in tasks")
				errc <- errors.New("no book in tasks")
			}
			res[i].Meta.Target = &s.TMaker.PendingBook[fi]
			res[i].Meta.Creater = s.TMaker

			//request for FileType(GET res[i].Target.FileType)
			go func(i int) {
				defer waitGp.Done()
				d, e := s.G(fmt.Sprintf(global.G_BKDETAIL, res[i].Meta.Target.BookID), map[string]string{"pinst": s.TMaker.SchoolID})
				if e != nil {
					fmt.Printf("e: %v\n", e)
					errc <- e
				}
				m := map[string]any{}
				e = json.Unmarshal(d, &m)
				if e != nil {
					fmt.Printf("err: %v\n", e)
					errc <- e
				}
				res[i].Meta.Target.FileType = int(m["fileType"].(float64))
			}(i)

		}
		fmt.Printf("\"waitGpafend\": %v\n", "waitGpafend")
		waitGp.Wait()

	}
	/*
		这部分for循环是为了设置以下字段:
		res[i].StartAt 需要请求
		res[i].NowAt=res[i].StartAt
		res[i].CreateT = time.Now() 请求后
	*/

	/*
		// {
		// 	for _, v := range res {
		// 		switch v.Meta.Target.FileType {
		// 		default:
		// 			fmt.Printf("unknown v.Target.FileType: %v\n", v.Meta.Target.FileType)
		// 			return nil, errors.New("unknown v.Target.FileType")
		// 		case 0:
		// 			waitGp.Add(1)
		// 			dc, ec := s.G(fmt.Sprintf(global.G0_PROCESS, v.Meta.Target.BookID), nil)
		// 			go func(v *model.TProcess, dc <-chan []byte, er <-chan error) {
		// 				defer waitGp.Done()

		// 				select {
		// 				case <-ctx.Done():
		// 					return
		// 				case e := <-ec:
		// 					fmt.Printf("e: %v\n", e)
		// 					cancelCall(e)
		// 					return
		// 				case d := <-dc:
		// 					m := map[string]int{}
		// 					err := json.Unmarshal(d, &m)
		// 					if err != nil {
		// 						fmt.Printf("err: %v\n", err)
		// 					}
		// 					v.Meta.StartAt = m["paragraph"]
		// 					v.NowAt = m["paragraph"]
		// 					v.Meta.CreateT = time.Now()
		// 				}
		// 			}(&v, dc, ec)
		// 		case 3:
		// 			waitGp.Add(1)
		// 			dc, ec := s.G(fmt.Sprintf(global.G3_PROCESS, v.Meta.Target.BookID), nil)
		// 			go func(v *model.TProcess, dc <-chan []byte, er <-chan error) {
		// 				defer waitGp.Done()

		// 				select {
		// 				case <-ctx.Done():
		// 					return
		// 				case e := <-ec:
		// 					fmt.Printf("e: %v\n", e)
		// 					cancelCall(e)
		// 					return
		// 				case d := <-dc:
		// 					m := map[string]int{}
		// 					err := json.Unmarshal(d, &m)
		// 					if err != nil {
		// 						fmt.Printf("err: %v\n", err)
		// 					}
		// 					v.Meta.StartAt = m["page"]
		// 					v.NowAt = m["page"]
		// 					v.Meta.CreateT = time.Now()
		// 				}
		// 			}(&v, dc, ec)
		// 		}
		// 	}
		// 	waitGp.Wait()
		// 	if ctx.Err() != nil {
		// 		return nil, ctx.Err()
		// 	}
		// }
	*/
	//这部分为了得到目录
	fmt.Printf("\"wait done\": %v\n", "wait done")
	errc <- nil
	return res, errc
}
func (*TaskMakerS) fromPendingToModel(d []dto.PendingBook) []model.Book {
	//为解决传输对象与模型不一致的问题而做出转换
	res := make([]model.Book, len(d))
	for k, v := range d {
		res[k].Title = v.Title
		res[k].BookID = v.BookID
		res[k].FileType = -1
	}
	return res
}

func (t *TaskMakerS) find(bid string) int {
	for k, v := range t.TMaker.PendingBook {
		if v.BookID == bid {
			return k
		}
	}
	return -1
}
func (t *TaskMakerS) Uid() string {
	return t.TMaker.Username
}

// func (s *TaskMakerS) GetAllBookType() error {
// 	var waitGp sync.WaitGroup
// 	ctx, cancelCall := context.WithCancelCause(context.Background())
// 	waitGp.Add(len(s.TMaker.PendingBook))
// 	for k, v := range s.TMaker.PendingBook {
// 		go func(ctx context.Context, bookIndex int, v model.Book) {
// 			defer waitGp.Done()
// 			d, e := s.G(fmt.Sprintf(global.G_BKDETAIL, v.BookID), map[string]string{"pinst": s.TMaker.SchoolID})
// 			if e != nil {
// 				fmt.Printf("e: %v\n", e)
// 				cancelCall(e)
// 				return
// 			}
// 			m := map[string]any{}
// 			e = json.Unmarshal(d, &m)
// 			if e != nil {
// 				fmt.Printf("err: %v\n", e)
// 				cancelCall(e)
// 			}
// 			s.TMaker.PendingBook[bookIndex].FileType = int(m["fileType"].(float64))
// 		}(ctx, k, v)
// 	}

// 	waitGp.Wait()

//		return ctx.Err()
//	}
