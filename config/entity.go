package config

type System struct {
	HasInit  bool `json:"hasInit"`
	HasLogin bool `json:"hasLogin"`
}

type Manage struct {
	Port    int    `json:"port"`
	WebPath string `json:"webPath"`
}

type Api struct {
	Port int `json:"port"`
}

type SetInfo struct {
	HasInit   bool    `json:"hasInit"`
	DbType    string  `json:"dbType"`
	CachePath string  `json:"cachePath"`
	Sqlite    *Sqlite `json:"sqlite"`
	Mysql     *Mysql  `json:"mysql"`
	Admin     *Admin  `json:"admin"`
	Manage    *Manage `json:"manage"`
	Api       *Api    `json:"api"`
}

var defaultSetInfo = &SetInfo{
	HasInit:   false,
	DbType:    "sqlite",
	CachePath: ".cache",
	Sqlite:    &Sqlite{Filename: "d-mail.db"},
	Mysql:     &Mysql{Host: "", Port: 3306, Dbname: "d-main", Username: "", Password: "", Charset: "utf-8"},
	Admin:     &Admin{Username: "", Password: ""},
	Manage:    &Manage{WebPath: "web", Port: 12566},
	Api:       &Api{Port: 12566},
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
