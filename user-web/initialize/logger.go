package initialize

func InitLogger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}
