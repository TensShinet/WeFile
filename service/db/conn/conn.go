package conn

import "github.com/TensShinet/WeFile/logging"

var logger = logging.GetLogger("db_service_conn")

func Init() {
	dbInit()
	redisInit()
}
