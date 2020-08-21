package rule

type AccessLog struct {
	RemoteAddr   string `ngx:"remote_addr"`            // 用户ip
	RemoteUser   string `ngx:"remote_user"`            // 用户xx
	TimeLocal    string `ngx:"time_local"`             // 访问时间
	Request      string `ngx:"request"`                // http包第一行
	Status       int    `ngx:"status"`                 // 状态码
	BodyByteSent int    `ngx:"body_bytes_sent"`        // 请求体大小
	HttpReferer  string `ngx:"http_referer"`           // referer
	RequestTime  string `ngx:"request_time"`           //
	URT          string `ngx:"upstream_response_time"` // 上游响应时间
	UA           string `ngx:"http_user_agent"`        // UA
	XFF          string `ngx:"http_x_forwarded_for"`   // XFF
	RequestBody  string `ngx:"request_body"`           // 请求体
	AccessToken  string `ngx:"http_accesstoken"`       //
	UidGot       string `ngx:"uid_got"`                //
	UidSet       string `ngx:"uid_set"`                //
	Host         string `ngx:"host"`                   // host字段头
	Cookie       string `ngx:"http_cookie"`            // 用户cookie
}
