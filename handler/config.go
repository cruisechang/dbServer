package handler

const (
	timeFormat = "2006-01-02 15:04:05" //不管格式如何，壹定要是這個時間
	//"Jan 2, 2006 at 3:04pm (MST)"
	//"2006-Jan-02"

	CodeSuccess        = 0
	CodeErrorUndefined = 1
	CodeIPError        = 2
	CodeRouteError     = 3
	CodeHeaderError    = 4
	CodeMethodError    = 5

	CodePathError = 10

	CodeBodyError0    = 20
	CodeBodyNil       = 21
	CodeBodyReadError = 22
	CodeBodyReadNil   = 23

	CodePanic = 40

	CodeRequestDataUnmarshalError = 100
	CodeRequestDataError          = 101
	CodeRequestPostDataIllegal    = 102
	CodeRequestPathError          = 103

	CodeResponseDataMarshalError = 200

	CodeDBPrepareError = 300
	CodeDBQueryError   = 301
	CodeDBScanError    = 302
	CodeDBRowsError    = 303

	CodeDBExecError       = 310 //insert duplicate value into unique column,
	CodeDBExecResultError = 311
	CodeDBExecLastIDError = 312

	//insert  插入了重複的資料到unique欄位，會錯
	//update  更新跟資料褲內相同東西，不會錯，但affected column count=0
	//delete  刪除不存在的資料，不會錯，但affected column count=0

	HeadAPIKey = "qwerASDFzxcv!@#$"
)

type responseData struct {
	Code    int         `json:"code"`
	Count   int         `json:"count"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//return data
type timeParam struct {
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

type userIDData struct {
	UserID uint64 `json:"userID"`
}

type passwordData struct {
	Password string `json:"password"`
}

type partnerIDData struct {
	PartnerID uint64 `json:"partnerID"`
}

type aesKeyData struct {
	AESKey string `json:"aesKey"`
}
type prefixData struct {
	Prefix string `json:"prefix"`
}
type apiBindIPData struct {
	APIBindIP string `json:"apiBindIP"` //json array ["xxx.xxx.xxx.xxx","xxx.xxx.xxx.xxx"]
}
type cmsBindIPData struct {
	CMSBindIP string `json:"cmsBindIP"` //json array ["xxx.xxx.xxx.xxx","xxx.xxx.xxx.xxx"]
}
type accessTokenData struct {
	AccessToken string `json:"accessToken"`
}
type userAccessTokenData struct {
	Found       int    `json:"found"`
	Active      int    `json:"active"`
	AccessToken string `json:"accessToken"`
}
type transferIDData struct {
	TransferID uint64 `json:"transferID"`
}

type creditData struct {
	Credit float32 `json:"credit"`
}
type loginData struct {
	Login uint `json:"login"`
}
type nameData struct {
	Name string `json:"name"`
}
type hlsURLData struct {
	HLSURL string `json:"hlsURL"`
}

type activeData struct {
	Active uint `json:"active"`
}
type bootData struct {
	Boot uint `json:"boot"`
}
type roundIDData struct {
	Round uint64 `json:"round"`
}

type statusData struct {
	Status uint `json:"status"`
}

type betCountdownData struct {
	BetCountdown uint `json:"betCountdown"`
}
type dealerIDData struct {
	DealerID uint `json:"dealerID"`
}

type hallIDData struct {
	HallID uint `json:"hallID"`
}
type roomIDData struct {
	RoomID uint `json:"roomID"`
}

type IDData struct {
	ID uint64 `json:"ID"`
}
type betIDData struct {
	BetID uint64 `json:"betID"`
}
type managerIDData struct {
	ManagerID uint `json:"managerID"`
}
type roleIDData struct {
	RoleID uint `json:"roleID"`
}
type broadcastIDData struct {
	BroadcastID uint64 `json:"broadcastID"`
}
type bannerIDData struct {
	BannerID uint64 `json:"bannerID"`
}

type checkLoginData struct {
	ManagerID uint   `json:"managerID"`
	Password  string `json:"password"`
	Active    uint   `json:"active"`
}

type dealerLoginData struct {
	DealerID int `json:"dealerID"`
	Active   int `json:"active"`
}

//user
//usersGetParam is users get parameters
type userGetParam struct {
	PartnerID int64  `json:"partnerID"`
	Active    int    `json:"active"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	OrderBy   string `json:"orderBy"`
	Order     string `json:"order"`
}

//userData is for  data response
type userData struct {
	UserID     uint64  `json:"userID"`
	PartnerID  uint64  `json:"partnerID"`
	Account    string  `json:"account"`
	Name       string  `json:"name"`
	Credit     float32 `json:"credit"`
	Level      int     `json:"level"`
	Category   int     `json:"category"`
	Active     int     `json:"active"`
	IP         string  `json:"ip"`
	Platform   int     `json:"platform"`
	Login      int     `json:"login"`
	CreateDate string  `json:"createDate"`
}

//usersPostParam si users post param
type userPostParam struct {
	PartnerID uint64 `json:"partnerID"`
	Account   string `json:"account"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Platform  int    `json:"platform"`
	//UserID    uint64 `json:"userID"` //後來自己取unique加入，post用
}

//userDB is for marshal  db rows
type userDB struct {
	user_id     uint64
	partner_id  uint64
	account     string
	name        string
	credit      float32
	level       int
	category    int
	active      int
	ip          string
	platform    int
	login       int
	create_date string
}

type userAccountGetParam struct {
	PartnerID int64  `json:"partnerID"` //有負的
	Password  string `json:"password"`
}

//for game server
type userAccessTokenGetParam struct {
	AccessToken string `json:"accessToken"`
}
type userAccessTokenDB struct {
	user_id             int64
	account             string
	credit              float32
	name                string
	partner_id          int64
	active              int
	access_token_expire string
}
type userAccessTokenGetData struct {
	UserID            int64   `json:"userID"`
	Account           string  `json:"account"`
	Credit            float32 `json:"credit"`
	Name              string  `json:"name"`
	PartnerID         int64   `json:"partnerID"`
	Active            int     `json:"active"`
	AccessTokenExpire string  `json:"accessTokenExpire"`
}

//user log
type userLogDB struct {
	log_id      uint64
	user_id     uint64
	account     string
	name        string
	category    int
	ip          string
	platform    int
	create_date string
}
type userLogData struct {
	LogID      uint64 `json:"logID"`
	UserID     uint64 `json:"userID"`
	Account    string `json:"account"`
	Name       string `json:"name"`
	Category   int    `json:"category"`
	IP         string `json:"ip"`
	Platform   int    `json:"platform"`
	CreateDate string `json:"createDate"`
}

//partner
//partnersGetParam is users get parameters
type partnerGetParam struct {
	Active  int    `json:"active"`
	OrderBy string `json:"orderBy"`
	Order   string `json:"order"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
}

//partnerData is for response
type partnerData struct {
	PartnerID  uint64 `json:"partnerID"`
	Account    string `json:"account"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
	Category   int    `json:"category"`
	Active     int    `json:"active"`
	APIBindIP  string `json:"apiBindIP"`
	CMSBindIP  string `json:"cmsBindIP"`
	CreateDate string `json:"createDate"`
}
type partnerDB struct {
	partner_id  uint64
	account     string
	name        string
	level       int
	category    int
	active      int
	api_bind_ip string
	cms_bind_ip string
	create_date string
}

//partnersPostParam si users post param
type partnerPostParam struct {
	Account     string `json:"account"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Category    int    `json:"catetory"`
	AESKey      string `json:"aesKey"`
	AccessToken string `json:"accessToken"`
	APIBindIP   string `json:"apiBindIP"`
	CMSBindIP   string `json:"cmsBindIP"`
	//PartnerID   uint64 `json:"partnerID"` //自己補入，給post用
}
type partnerPatchParam struct {
	Password    string `json:"password"`
	Name        string `json:"name"`
	Level       uint   `json:"level"`
	Category    uint   `json:"catetory"`
	AESKey      string `json:"aesKey"`
	AccessToken string `json:"accessToken"`
	APIBindIP   string `json:"apiBindIP"`
	CMSBindIP   string `json:"cmsBindIP"`
	Active      uint   `json:"active"`
}
type partnerAccountGetParam struct {
	Password string `json:"password"`
}

//partner log
type partnerLogDB struct {
	log_id      uint64
	partner_id  uint64
	account     string
	name        string
	category    int
	create_date string
}
type partnerLogData struct {
	LogID      uint64 `json:"logID"`
	PartnerID  uint64 `json:"partnerID"`
	Account    string `json:"account"`
	Name       string `json:"name"`
	Category   int    `json:"category"`
	CreateDate string `json:"createDate"`
}

//hall
type hallPostParam struct {
	HallID uint   `json:"hallID"`
	Name   string `json:"name"`
}
type hallData struct {
	HallID     uint   `json:"hallID"`
	Name       string `json:"name"`
	Active     uint   `json:"active"`
	CreateDate string `json:"createDate"`
}

type hallDB struct {
	hall_id     uint
	name        string
	active      uint
	create_date string
}

type hallPatchParam struct {
	HallID uint   `json:"hallID"`
	Name   string `json:"name"`
	Active uint   `json:"active"`
}

//room
type roomPostParam struct {
	RoomID       uint   `json:"roomID"`
	HallID       uint   `json:"hallID"`
	Name         string `json:"name"`
	RoomType     uint   `json:"roomType"`
	BetCountdown uint   `json:"betCountdown"`
	HLSURL       string `json:"hlsURL"`
	//DealerID     uint   `json:"dealerID"`
	LimitationID uint `json:"limitationID"`
}

type roomData struct {
	RoomID       uint   `json:"roomID"`
	HallID       uint   `json:"hallID"`
	Name         string `json:"name"`
	RoomType     uint   `json:"roomType"`
	Active       uint   `json:"active"`
	HLSURL       string `json:"hlsURL"`
	Boot         uint   `json:"boot"`
	RoundID      uint64 `json:"round"`
	Status       int    `json:"status"`
	BetCountdown uint   `json:"betCountdown"`
	DealerID     uint   `json:"dealerID"`
	LimitationID uint   `json:"limitationID"`
	CreateDate   string `json:"createDate"`
}
type roomDB struct {
	room_id       uint
	hall_id       uint
	name          string
	room_type     uint
	active        uint
	hls_url       string
	boot          uint
	round_id      uint64
	status        int
	bet_countdown uint
	dealer_id     uint
	limitation_id uint
	create_date   string
}

type roomPatchParam struct {
	RoomID       uint   `json:"roomID"`
	HallID       uint   `json:"hallID"`
	Name         string `json:"name"`
	RoomType     uint   `json:"roomType"`
	Active       uint   `json:"active"`
	HLSURL       string `json:"hlsURL"`
	BetCountdown uint   `json:"betCountdown"`
	LimitationID uint   `json:"limitationID"`
}

type roomNewRoundPatchParam struct {
	Boot    uint   `json:"boot"`
	RoundID uint64 `json:"round"`
	Status  int    `json:"status"`
}

//dealer
type dealerData struct {
	DealerID    uint   `json:"dealerID"`
	Name        string `json:"name"`
	Account     string `json:"account"`
	Active      int    `json:"active"`
	PortraitURL string `json:"portraitUrl"`
	CreateDate  string `json:"createDate"`
}
type dealerDB struct {
	dealer_id    uint
	name         string
	account      string
	active       int
	portrait_url string
	create_date  string
}

type dealerPostParam struct {
	Name        string `json:"name"`
	Account     string `json:"account"`
	Password    string `json:"password"`
	PortraitURL string `json:"portraitUrl"`
}
type dealerPatchParam struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	PortraitURL string `json:"portraitUrl"`
	Active      int    `json:"active"`
}
type dealerAccountGetParam struct {
	Password string `json:"password"`
}

//limitation
type limitationData struct {
	LimitationID uint   `json:"limitationID"`
	Limitation   string `json:"limitation"`
}
type limitationDB struct {
	limitation_id uint
	limitation    string
}

//transfer
type transferPostParam struct {
	//TransferID        uint64  `json:"transferID"` //自己產生unique 補入
	PartnerTransferID string  `json:"partnerTransferID"`
	PartnerID         uint64  `json:"partnerID"`
	UserID            uint64  `json:"userID"`
	Category          uint    `json:"category"`
	TransferCredit    float32 `json:"transferCredit"`
	Credit            float32 `json:"credit"`
	Status            uint    `json:"status"`
}
type transferGetParam struct {
	PartnerID int64  `json:"partnerID"` //有-1
	UserID    int64  `json:"userID"`    //有-1
	Category  int    `json:"category"`  //有-1
	Status    int    `json:"status"`    //有-1
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

//from db
type transferDB struct {
	transfer_id         uint64
	partner_transfer_id string
	partner_id          uint64
	user_id             uint64
	category            uint
	transfer_credit     float32
	credit              float32
	status              uint
	create_date         string
	account             string
	name                string
}

//response
type transferData struct {
	TransferID        uint64  `json:"transferID"`
	PartnerTransferID string  `json:"partnerTransferID"`
	PartnerID         uint64  `json:"partnerID"`
	UserID            uint64  `json:"userID"`
	Category          uint    `json:"category"`
	TransferCredit    float32 `json:"transferCredit"`
	Credit            float32 `json:"credit"`
	Status            uint    `json:"status"`
	CreateDate        string  `json:"createDate"`
	Account           string  `json:"account"`
	Name              string  `json:"name"`
}

//bet
type betPostParam struct {
	//BetID     uint64 `json:"betID"`
	PartnerID uint64 `json:"partnerID"`
	UserID    uint64 `json:"userID"`
	RoomID    uint   `json:"category"`
	RoomType  uint   `json:"roomType"`
	RoundID   uint64 `json:"round"`
	SeatID    int    `json:"seatID"`

	BetCredit      float32 `json:"betCredit"`
	ActiveCredit   float32 `json:"activeCredit"`
	PrizeCredit    float32 `json:"prizeCredit"`
	ResultCredit   float32 `json:"resultCredit"`
	BalanceCredit  float32 `json:"balanceCredit"`
	OriginalCredit float32 `json:"originalCredit"`
	Record         string  `json:"record"`
	Status         uint    `json:"status"`
}

type betGetParam struct {
	PartnerID int64  `json:"partnerID"` //有-1
	UserID    int64  `json:"userID"`    //有-1
	RoomID    int    `json:"category"`  //有-1
	RoomType  int    `json:"roomType"`  //有-1
	RoundID   int64  `json:"round"`     //有-1
	Status    int    `json:"status"`    //有-1
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

//form db
type betDB struct {
	bet_id          uint64
	partner_id      string
	user_id         uint64
	room_id         uint
	room_type       uint
	round_id        uint
	seat_id         int
	bet_credit      float32
	active_credit   float32
	prize_credit    float32
	result_credit   float32
	balance_credit  float32
	original_credit float32
	record          string
	status          uint
	create_date     string
	account         string
	name            string
}

//response
type betData struct {
	BetID          uint64  `json:"betID"`
	PartnerID      string  `json:"partnerID"`
	UserID         uint64  `json:"userID"`
	RoomID         uint    `json:"roomID"`
	RoomType       uint    `json:"roomType"`
	RoundID        uint    `json:"round"`
	SeatID         int     `json:"seatID"`
	BetCredit      float32 `json:"betCredit"`
	ActiveCredit   float32 `json:"activeCredit"`
	PrizeCredit    float32 `json:"prizeCredit"`
	ResultCredit   float32 `json:"resultCredit"`
	BalanceCredit  float32 `json:"balanceCredit"`
	OriginalCredit float32 `json:"originalCredit"`
	Record         string  `json:"record"`
	Status         uint    `json:"status"`
	CreateDate     string  `json:"createDate"`
	Account        string  `json:"account"`
	Name           string  `json:"name"`
}

//round
type roundPostParam struct {
	HallID   uint   `json:"hallID"`
	RoomID   uint   `json:"round"`
	RoomType uint   `json:"roomType"`
	Brief    string `json:"brief"`
	Record   string `json:"record"`
	Status   uint   `json:"status"`
}

type roundGetParam struct {
	HallID    int    `json:"round"`    //有-1
	RoomID    int    `json:"roomID"`   //有-1
	RoomType  int    `json:"roomType"` //有-1
	Status    int    `json:"status"`   //有-1
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}
type roundDB struct {
	round_id    uint
	hall_id     uint
	room_id     uint
	room_type   uint
	brief       string
	record      string
	status      uint
	create_date string
	end_datea   string
	name        string
}

type roundData struct {
	RoundID    uint   `json:"round"`
	HallID     uint   `json:"hallID"`
	RoomID     uint   `json:"roomID"`
	RoomType   uint   `json:"roomType"`
	Brief      string `json:"brief"`
	Record     string `json:"record"`
	Status     uint   `json:"status"`
	CreateDate string `json:"createDate"`
	EndDate    string `json:"endDate"`
	Name       string `json:"name"`
}

type roundPatchParam struct {
	Brief  string `json:"brief"`
	Record string `json:"record"`
	Status uint   `json:"status"`
}

//official cms manager
type officialCMSManagerPostParam struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	RoleID   uint   `json:"roleID"`
}

type officialCMSManagerDB struct {
	manager_id  uint
	account     string
	active      uint
	role_id     uint
	login       uint
	create_date string
}

type officialCMSManagerData struct {
	ManagerID  uint   `json:"managerID"`
	Account    string `json:"account"`
	Active     uint   `json:"active"`
	RoleID     uint   `json:"roleID"`
	Login      uint   `json:"login"`
	CreateDate string `json:"createDate"`
}

type officialCMSManagerPatchParam struct {
	Password string `json:"password"`
	RoleID   uint   `json:"roleID"`
	Active   uint   `json:"active"`
}

//official cms role
type officialCMSRolePostParam struct {
	Permission string `json:"permission"`
}
type officialCMSRoleDB struct {
	role_id     uint
	permission  string
	create_date string
}
type officialCMSRoleData struct {
	RoleID     uint   `json:"roleID"`
	Permission string `json:"permission"`
	CreateDate string `json:"createDate"`
}

type officialCMSRolePatchParam struct {
	Permission string `json:"permission"`
}

//broadcast
//收到的資料，要比較嚴格
type broadcastPostParam struct {
	Content     string `json:"content"`
	Internal    uint    `json:"internal"`
	RepeatTimes uint    `json:"repeatTimes"`
	Active      uint    `json:"active"`
}
type broadcastData struct {
	BroadcastID int    `json:"broadcastID"`
	Content     string `json:"content"`
	Internal    int    `json:"internal"`
	RepeatTimes int    `json:"repeatTimes"`
	Active      uint   `json:"active"`
	CreateDate  string `json:"createDate"`
}

type broadcastDB struct {
	broadcast_id int
	content      string
	internal     int
	repeat_times int
	active       uint
	create_date  string
}

//banner

type bannerPostParam struct {
	PicURL      string `json:"picURL"`
	LinkURL     string `json:"linkURL"`
	Description string `json:"description"`
	Platform    uint   `json:"platform"`
	Active      uint   `json:"active"`
}
type bannerDB struct {
	banner_id   int
	pic_url     string
	link_url    string
	description string
	platform    uint
	active      uint
	create_date string
}
type bannerData struct {
	BannerID    int    `json:"bannerID"`
	PicURL      string `json:"picURL"`
	LinkURL     string `json:"linkURL"`
	Description string `json:"description"`
	Platform    uint    `json:"platform"`
	Active      uint    `json:"active"`
	CreateDate  string `json:"createDate"`
}

