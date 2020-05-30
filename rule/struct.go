package rule

type Intercept struct {
	TimeStamp string              `json:"TimeStamp"` // 访问时间
	Addr      string              `json:"Addr"`      // 来源地址，一般是WAF地址
	Host      string              `json:"Host"`      // host头，区分业务
	UA        string              `json:"UA"`
	URI       string              `json:"URI"`
	Query     string              `json:"Query"`
	Rule      string              `json:"Rule"`      // 匹配上的规则，便于分类拦截
	XFF       string              `json:"XFF"`       // X-Forwarded-For
	XRI       string              `json:"XRI"`       // X-Real-IP
	Method    string              `json:"Method"`    // 请求方式
	APP       string              `json:"APP"`       // 所属应用
	Headers   map[string][]string `json:"Headers"`   // 完整头数据
	Status    int                 `json:"Status"`    // 响应状态码
	Size      int                 `json:"Size"`      // 响应包长度
	Body      string              `json:"Body"`      // post body全包
	LRegion   string              `json:"LRegion"`   // 组织
	LCountry  string              `json:"LCountry"`  // 国家
	LProvince string              `json:"LProvince"` // 省份
	LCity     string              `json:"LCity"`     // 城市
	LCounty   string              `json:"LCounty"`   // 镇
	Local     string              `json:"Local"`     // 拦截中心地址
}

type AccessLog struct {
	Host    string  `json:"host"`    // WAF字段，域名
	Status  int     `json:"status"`  // WAF字段，状态码
	XFF     string  `json:"XFF"`     // WAF字段，X-Forwarded-for
	Rule    string  `json:"rule"`    // WAF字段，触发规则名
	Size    int     `json:"size"`    // WAF字段，响应包大小
	Method  string  `json:"method"`  // WAF字段，请求方法
	URI     string  `json:"uri"`     // WAF字段，请求uri
	Reqs    string  `json:"reqs"`    // WAF字段，单个tcp链接复用次数
	Uaddr   string  `json:"uaddr"`   // WAF字段，上游地址
	Time    string  `json:"time"`    // WAF字段，请求时间的unix时间戳
	Port    string  `json:"port"`    // WAF字段，WAF使用端口
	APP     string  `json:"app"`     // WAF字段，应用名称
	CDN     string  `json:"cdn"`     // WAF字段，是否来自于公司CDN
	Addr    string  `json:"addr"`    // WAF字段，来源地址
	URT     float64 `json:"urt"`     // WAF字段，上游响应时间
	Pass    string  `json:"pass"`    // WAF字段，是否匹配白名单
	Query   string  `json:"query"`   // WAF字段，问号后的查询字段
	Remote  string  `json:"remote"`  // WAF字段，四层来源地址
	REF     string  `json:"ref"`     // WAF字段，来源地址referer
	UA      string  `json:"ua"`      // WAF字段，UA头
	Risk    string  `json:"risk"`    // WAF字段，风控风险值字段
	Uname   string  `json:"uname"`   // WAF字段，上游服务器组合名
	Conn    string  `json:"conn"`    // WAF字段，tcp链接的序列号
	Local   string  `json:"local"`   // WAF字段，请求时间的正常格式
}
