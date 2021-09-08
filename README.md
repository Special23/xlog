# xlog
log library for golang.

# log level
const (
   XLOG_FLAG_DEBUG = 1
   XLOG_FLAG_TRACE = 2
   XLOG_FLAG_INFO  = 4
   XLOG_FLAG_ERROR = 8
   XLOG_FLAG_ALL   = XLOG_FLAG_DEBUG | XLOG_FLAG_TRACE | XLOG_FLAG_INFO | XLOG_FLAG_ERROR
)

# Use example

log := xlog.NewXLog("../logpath/", "", XLOG_FLAG_INFO)
log.Debug("test debug log %d", 123)
log.Info("test info log %d", 123)
log.Error("test error log %d", 123)
