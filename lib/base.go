package lib

import (
	"github.com/Qesy/qesygo"
)

// Config 系统配置文件
type Config struct {
	Db    map[string]string `json:"Db"`
	Cache map[string]string `json:"Cache"`
	Conf  map[string]string `json:"Conf"`
	Amqp  map[string]string `json:"Amqp"`
}

type ApiRet struct {
	Code int               `json:"Code"`
	Msg  string            `json:"Msg"`
	Data map[string]string `json:"Data"`
}

// 配置文件
var ConfRs Config

// RedisCr Redis进程池
var RedisCr qesygo.CacheRedis

// 雪花算法 Worker
var SnowWorker *qesygo.Worker

// 服务器ID
var ServerID int64

const Key string = "1234!@#$" //Key 网站加密种子
const LogFilePath string = "./static/log/error"
const LogActPath string = "./logData/log"

var StaticFiles = []string{
	"codeError",
	"RoleLevel",
	"Gate",
	"GateTemplate",
	"ItemNew",
	"RobotName",
	"common",
	"qualityUltimate",
	"ultimateBasic",
	"ultimateSpecial",
	"ultimateQualityAttribut",
	"equipmentSetting",
	"equipmentUpdateSetting",
	"equipmentConstantSetting",
	"treasureBox",
	"treasureTask",
	"drop",
	"dropgroup",
	"dailyTaskSetting",
	"weeklyTaskSetting",
	"achievementTaskSetting",
	"TaskActiveSetting",
	"HangupReward",
	"resetTask",
	"reset",
	"battleitems",
	"RoleProp",
	"entryTranslate",
	"shopItemsSetting",
	"shopBasicSetting",
	"shopProductSetting",
	"skillBank",
	"activitySetting",
	"privilegeCard",
	"skillLevel",
	"soulSkills",
	"adMonthlyCard",
	"soulBasic",
	"ActiveSkills",
	//"ActiveSkillsExtra",
	"mailSetting",
	"FirstRecharge",
	"treasureTaskLoop",
	"sevenDayTag",
	"developEffect",
	"DailyReward",
	"foundation",
	"foundationBasic",
	"PassiveSkills",
	"heroBasic",
	"heroLevel",
	"heroSoul",
	"EquipmentConversion",
}

const COINFREE string = "1101"        //免费币
const COINPAY string = "1201"         //付费币
const EXP string = "1301"             //经验
const CHEST_POINTS string = "1401"    //宝箱积分
const TASK_SCORE_DAY string = "1501"  //日活跃值
const TASK_SCORE_WEEK string = "1601" //周活跃值
const TASK_SCORE_7DAY string = "1602" //七日活动活跃值
const REBIRTH string = "1701"         //复活券
const USERSETCARD string = "2501"     //改名卡
const FURY string = "2209"            //狂暴之力

const CHEST string = "2101"      //默认闯关发放宝箱（后期可能不需要）
const CHEST_LOOP_ID int = 100001 //循环宝箱任务起始ID
const ACTIVITY_7DAY_ID string = "4001"

const ItemChannelSGm string = "-1"                //后台发放
const ItemChannelCGm string = "-2"                //客户端调试发放
const ItemChannelTest string = "-3"               //调试发放
const ItemChannelOneKeySend string = "-4"         //一键调试发放
const ItemChannelGate string = "1"                //关卡奖励
const ItemChannelEquipUpgrade string = "2"        //装备升级
const ItemChannelChest string = "3"               //开宝箱
const ItemChannelChestTask string = "4"           //宝箱任务
const ItemChannelAchievement string = "5"         //成就领取
const ItemChannelActive string = "6"              //成就活跃领取
const ItemChannelHang string = "7"                //挂机奖励
const ItemChannelRebirth string = "8"             //战斗内复活
const ItemChannelEquipReset string = "9"          //装备重置
const ItemChannelEquipCompose string = "10"       //装备合成
const ItemChannelSweeping string = "11"           //扫荡
const ItemChannelUserReset string = "12"          //转生
const ItemChannelShopFlush string = "13"          //商店刷新
const ItemChannelShopBuy string = "14"            //商店购买
const ItemChannelSetUserInfo string = "15"        //设置用户信息（改名卡）
const ItemChannelSkillBank string = "16"          //天书激活
const ItemChannelOrderPay string = "17"           //订单支付
const ItemChannelPrivilegeCard string = "18"      //特权卡购买时领取
const ItemChannelPrivilegeCardDaily string = "19" //特权卡每日领取
const ItemChannelSkillSoul string = "20"          //激活元神
const ItemChannelMonthCard string = "21"          //月卡领取
const ItemChannelAccountCreate string = "22"      //创角领取
const ItemChannelOpen string = "23"               //打开道具获得
const ItemChannelMail string = "24"               //邮件领取
const ItemChannelChestTen string = "25"           //十连抽宝箱
const ItemChannelGateFury string = "26"           //狂暴之力
const ItemChannelCpGift string = "27"             //CP礼包兑换
const ItemChannelLevelUp string = "28"            //升级奖励
const ItemChannelBindOneKey string = "29"         //一键绑定
const ItemChannelShare string = "30"              //分享
const ItemChannelFirstRecharge string = "31"      //首充奖励领取
const ItemChannelUserResetTask string = "32"      //转生
const ItemChannelGateRobot string = "33"          //关卡机器人提交
const ItemChannelGatePickup string = "34"         //关卡拾取
const ItemChannelGateDouble string = "35"         //关卡看广告（双倍奖励）
const ItemChannelFoundation string = "36"         //基金
const ItemChannelChallenge string = "37"          //挑战关卡领取
const ItemChannelHeroUpgrade string = "38"        //英雄升级
const ItemChannelEquipConvert string = "39"       //装备转换
const ItemChannelHeroUpDown string = "40"         //英雄上下场
const ItemChannelHeroSoul string = "41"           //英雄命魂升级

const EquipActAdd string = "1"    //添加
const EquipActUpdate string = "2" //修改
const EquipActDel string = "3"    //删除

var UserAttr []string = []string{COINFREE, COINPAY, EXP, CHEST_POINTS, TASK_SCORE_DAY, TASK_SCORE_WEEK, TASK_SCORE_7DAY}
var PlatformMap map[string]string = map[string]string{
	"1": "Test",
	"2": "KingNet",
}
