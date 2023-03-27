package controler

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

//	type TMakerControler struct {
//		Tms *service.TaskMakerS
//		Ts  *service.TasksS
//	}
type ControlMan struct {
	Tmk    map[UserName]UserInit
	Te     map[UserName]map[TaskID]ExecAble
	Tc     map[UserName]*ConcurrentLimit //并发控制
	Poster chan map[string]any
}
type UserName = string
type TaskID = string

type ConcurrentLimit struct {
	Ch     chan struct{}
	Ctx    context.Context
	Calcel context.CancelCauseFunc
	Gap    int
}
type UserInit interface {
	Login(context.Context) error
	GetPendingBooks() error
	Books2Tasks(bids []string) ([]ExecAble, error)
	Uid() string
}
type ExecAble interface {
	Init() error
	Do() (bool, error)
	Fin()

	Uid() string
	Ctx() context.Context
	NowAt() float32
}

func (cm *ControlMan) DelUser(who string) {
	delete(cm.Tmk, who)
	delete(cm.Te, who)
	cm.Tc[who].Calcel(errors.New(who + "is deleted"))
}
func (cm *ControlMan) RegisterUser(instance UserInit, concurrentNum int, gap int) {
	ctxForUser, whyCall := context.WithCancelCause(context.Background())
	cm.Tc[instance.Uid()] = &ConcurrentLimit{
		Ch:     make(chan struct{}, concurrentNum),
		Ctx:    ctxForUser,
		Calcel: whyCall,
		Gap:    gap,
	}
	_, ok := cm.Tmk[instance.Uid()]
	if ok {
		return
	}
	maxRetryTime := 3

	//login
	for maxRetryTime > 0 {
		err := instance.Login(ctxForUser)
		if err == nil {
			break
		}

		fmt.Printf("err: %v\n", err)
		maxRetryTime--
		<-time.After(3 * time.Second)
	}
	if maxRetryTime == 0 {
		cm.DelUser(instance.Uid())
		whyCall(errors.New("try login over 3times"))
		return
	}

	maxRetryTime = 5
	for maxRetryTime > 0 {
		err := instance.GetPendingBooks()
		if err == nil {
			break
		}
		fmt.Printf("err: %v\n", err)
		maxRetryTime--
	}
	if maxRetryTime == 0 {
		cm.DelUser(instance.Uid())
		whyCall(errors.New("unknow task of books"))
		return
	}
	cm.Tmk[instance.Uid()] = instance
	cm.Te[instance.Uid()] = map[TaskID]ExecAble{}
}

func (cm *ControlMan) AddTask(uid string, bids ...string) {

	maxRetryTime := 5
	who, ok := cm.Tmk[uid]
	if who == nil || !ok {
		return
	}
	var exeos []ExecAble
	//trans
	for maxRetryTime > 0 {
		var err error
		exeos, err = who.Books2Tasks(bids)
		if err == nil {
			break
		}
		fmt.Printf("err: %v\n", err)
		maxRetryTime--
	}
	//init
	// for _, ea := range exeos {
	// 	maxRetryTime = 5
	// 	for maxRetryTime > 0 {
	// 		err := ea.Init()
	// 		if err == nil {
	// 			cm.Te[uid][ea.Uid()] = ea
	// 			break
	// 		}
	// 		fmt.Printf("err: %v\n", err)
	// 		maxRetryTime--
	// 	}
	// }
	// cm.Poster <- map[string]any{
	// 	"type": "initfin",
	// 	"who":  uid,
	// }
	exeos[0].Init()
	cm.Te[uid][exeos[0].Uid()] = exeos[0]
	// run
	// for _, ea := range exeos {
	go cm.execTask(uid, exeos[0])
	// }
}
func (cm *ControlMan) DelTask(uid string, bid string) {
	cm.Te[uid][bid].Fin()
	delete(cm.Te[uid], bid)
}
func (cm *ControlMan) GetAllProcess() map[UserName]map[TaskID]ExecAble {
	return cm.Te
}
func (cm *ControlMan) execTask(who string, ea ExecAble) {
	ctl := cm.Tc[who]
	for {
		select {
		case <-ctl.Ctx.Done():
			return
		case ctl.Ch <- struct{}{}:
			done, err := ea.Do()
			if err != nil {
				cm.DelTask(who, ea.Uid())
				cm.Poster <- map[string]any{
					"type": "err",
					"who":  who,
					"uid":  ea.Uid(),
					"err":  err,
				}
				return
			}
			if done {
				<-ctl.Ch
				cm.Poster <- map[string]any{
					"type": "done",
					"who":  who,
					"bid":  ea.Uid(),
				}
				return
			}
			<-ctl.Ch
			cm.Poster <- map[string]any{
				"type": "process",
				"d":    fmt.Sprintf("%f", ea.NowAt()),
				"bid":  ea.Uid(),
				"who":  who,
			}

			<-time.After(time.Second * time.Duration((rand.Intn(ctl.Gap*2-1) + 1)))
		}
	}
}

// func (cm *ControlMan) DelUser(uid string) {

// }

// func goexec(ctx context.Context, ea ExecAble, who string, ch chan struct{}, poster chan map[string]any) {
// 	fmt.Printf("ea.Uid(): %v\n", ea.Uid())
// 	select {
// 	case <-ctx.Done():
// 		return
// 	case ch <- struct{}{}:
// 		next, err := ea.Do()
// 		<-ch

// 		if err != nil {
// 			errmap := map[string]any{
// 				"type": "error",
// 				"d":    err.Error(),
// 			}
// 			poster <- errmap

// 		}
// 		if next == nil {
// 			logmap := map[string]any{
// 				"type": "info",
// 				"d":    ea.Uid() + "fin",
// 			}

// 			poster <- logmap
// 			return
// 		}
// 		processmap := map[string]any{
// 			"type": "process",
// 			"d":    fmt.Sprintf("%f", ea.NowAt()),
// 			"bid":  ea.Uid(),
// 			"who":  who,
// 		}
// 		poster <- processmap
// 		<-time.After(3 * time.Second)
// 		go goexec(ctx, next, who, ch, poster)
// 	}

// }
