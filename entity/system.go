package entity

type System struct {
	HasInit  bool `json:"hasInit"`
	HasLogin bool `json:"hasLogin"`
}

type SetInfo struct {
	HasInit bool    `json:"hasInit"`
	DbType  string  `json:"dbType"`
	Sqlite  *Sqlite `json:"sqlite"`
	Mysql   *Mysql  `json:"mysql"`
	Admin   *Admin  `json:"admin"`
}

type Sqlite struct {
	Filename string `json:"filename"`
}
type Mysql struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Dbname   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
}

type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
