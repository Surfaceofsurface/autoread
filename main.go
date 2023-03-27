package main

import (
	"fmt"
	"net"
	"private/autoread/reactor/controler"
	oapi "private/autoread/surface/controler"
)

func main() {
	ch := make(chan struct{}, 1)
	cman := controler.ControlMan{
		Tmk:    map[string]controler.UserInit{},
		Te:     map[string]map[string]controler.ExecAble{},
		Tc:     map[string]*controler.ConcurrentLimit{},
		Poster: make(chan map[string]any, 1),
	}
	go func() {
		logIP()
		oapi.X(&cman)
		ch <- struct{}{}
	}()
	// viper.AddConfigPath("./")
	// viper.SetConfigType("json")
	// viper.SetConfigName("config")
	// viper.ReadInConfig()
	// psw := viper.GetStringSlice("psw")
	// uid := viper.GetStringSlice("uid")
	// fmt.Printf("psw: %v\n", psw)
	// for i, v := range psw {
	// 	u := &service.TaskMakerS{
	// 		TMaker: &model.TaskMaker{
	// 			Password: v,
	// 			Username: uid[i],
	// 			SchoolID: "1cdceffd0000020bce"},
	// 	}
	// 	cman.RegisterUser(u)
	// 	xx := []string{}
	// 	for _, b := range u.TMaker.PendingBook {
	// 		xx = append(xx, b.BookID)
	// 	}
	// 	// cman.AddTask(uid[i], xx...)
	// }
	<-ch
}
func logIP() {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("UI running at one of them:")

	for _, address := range addrs {
		// 检查ip地址的合法性
		if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.IsPrivate() {
			if ipnet.IP.To4() != nil {
				fmt.Println("http://" + ipnet.IP.String() + ":44551/h")
			}
		}
	}

}
