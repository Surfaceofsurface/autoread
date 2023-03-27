package global

const (
	LOCAL_PROXY = "http://localhost:8888"
	LOGIN       = "https://m.cxstar.com/api/auth/login"
	G_READING   = "https://m.cxstar.com/api/user/readings" //?page=1&size=10&state=1
	G_BKDETAIL  = "https://m.cxstar.com/api/books/%s"
	G_CATALOG   = "https://m.cxstar.com/api/books/%s/catalog"            //?filetype=0
	PROCESS_3   = "https://m.cxstar.com/api/books/%s/paragraph/progress" //段落
	PROCESS_0   = "https://m.cxstar.com/api/books/%s/read/progress"      //图片
	PROCESS_N   = "https://m.cxstar.com/api/books/%s/read"
	P_TRACK     = "https://m.cxstar.com/api/user/ReadTrack"
	G_STATE     = "https://m.cxstar.com/api/books/%s/state"
)

//https://m.cxstar.com/api/books/211de93d000001XXXX/read?page=334&token=fbef0643-28af-4ba1-b493-eed8ab748300&pinst=1cdceffd0000020bce&nonce=34a339fc-52ab-4233-bba0-ae87b20c8aab&stime=1677945088&sign=7E93CC5CF8C5CE1611E95C9668DB44FB
// MD5(atob("MTIzNDU2") + nonce + stime)
//nonce=random="972fc563-4585-4d28-a09b-6ff2556b100f"
//stime=now - 95000

//fileType=1 'cxbf'
//=0 'pdf'
//=3 default

//https://m.cxstar.com/api/books/211de93d000001XXXX/status?lengthTime=1&fileType=0

// https://m.cxstar.com/api/books/24026b13000001XXXX/read?page=219&from=default&pinst=1cdceffd0000020bce&nonce=fd4ebb4d-091a-4d6f-8c7a-4fc0412924cc&stime=1679628076&sign=CFCE904CD3851F922F6F58E2065BD223
