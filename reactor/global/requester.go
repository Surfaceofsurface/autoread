package global

type Requester map[string]ReqMethod

var R Requester

func init() {
	R = make(Requester)
}
func (r Requester) Register(uid string, opt EnvOpt) {
	r[uid] = NewReqEnv(opt)
}
func (r Requester) GetEnv(uid string) ReqMethod {
	return r[uid]
}
