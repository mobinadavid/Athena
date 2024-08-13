package requests

type IpgCallbackRequest struct {
	MID              int    `form:"MID"`
	TerminalId       int    `form:"TerminalId"`
	RefNum           string `form:"RefNum"`
	ResNum           string `form:"ResNum"`
	State            string `form:"State"`
	TraceNo          string `form:"TraceNo"`
	Amount           int    `form:"Amount"`
	Wage             string `form:"Wage"`
	Rrn              string `form:"Rrn"`
	SecurePan        string `form:"SecurePan"`
	Status           int    `form:"Status"`
	Token            string `form:"Token"`
	HashedCardNumber string `form:"HashedCardNumber"`
}
