package smtp

import (
	"github.com/chuccp/smtp2http/db"
	"testing"
)

func TestSendAPIMail(t *testing.T) {

	err := SendAPIMail(&db.Schedule{Url: "https://www.baidu.com/"}, &db.SMTP{}, []*db.Mail{})
	if err != nil {
		return
	}

}
