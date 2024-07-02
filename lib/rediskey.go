package lib

// 全局通用
func RedisKey_Table_HM_Key(Table string, PrimaryId string) string {
	return ConfRs.Conf["Name"] + "_Table_" + Table + "_HM_Key_" + PrimaryId
}

// 分布式锁
func RedisKey_Crontab_SetNx() string { // 定时器
	return ConfRs.Conf["Name"] + "_Crontab_SetNx"
}

func RedisKey_Order_CallBack_SetNx(OrderSn string) string { // 订单回调
	return ConfRs.Conf["Name"] + "_Order_CallBack_SetNx_" + OrderSn
}
