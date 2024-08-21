// 给予他人ATRI币
package givemoney

import (
	"github.com/FloatTech/AnimeAPI/wallet"
	"github.com/FloatTech/floatbox/math"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("givemoney", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		PublicDataFolder: "Givemoney",
		Brief:            "givemoney",
		OnEnable: func(ctx *zero.Ctx) {
			ctx.Send("插件已启用")
		},
		OnDisable: func(ctx *zero.Ctx) {
			ctx.Send("插件已禁用")
		},
	})
	engine.OnRegex(`^给予.*?(\d+).*?\s(\d+)(.*)`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			nickname := ctx.GetThisGroupMemberInfo(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]),
				false,
			).Get("nickname").Str

			omy := math.Str2Int64(ctx.State["regex_matched"].([]string)[2])
			switch {
			case omy > 1000:
				ctx.SendChain(message.Text("超出额度啦~"))
				return
			case omy <= 0:
				ctx.SendChain(message.Text("不要想着钻漏洞哦"))
				return
			case 0 < omy && omy <= 1000:
				var (
					money = int(omy)
					uid1  = ctx.Event.UserID
					uid2  = math.Str2Int64(ctx.State["regex_matched"].([]string)[1])
				)
				if wallet.GetWalletOf(uid1) < money {
					ctx.SendChain(message.Text("ATRI币不足,发送签到获取吧~"))
					return
				}
				err := wallet.InsertWalletOf(uid1, -money)
				if err != nil {
					ctx.SendChain(message.Text("ATRI币扣除失败  [ERROR at gm.go:51]:", err))
					return
				}
				err1 := wallet.InsertWalletOf(uid2, money)
				if err1 != nil {
					ctx.SendChain(message.Text("获得ATRI币出错啦,返回扣除的ATRI币  [ERROR at gm.go:55]:", err1))
					err := wallet.InsertWalletOf(uid1, money)
					if err != nil {
						ctx.SendChain(message.Text("返回ATRI币失败,您的钱被吞掉了  [ERROR at gm.go:58]:", err))
						return
					}
				}
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(nickname+"接受了你的"+ctx.State["regex_matched"].([]string)[2]+"ATRI币"))
				return
			default:
				ctx.SendChain(message.Text("发生未知错误"))
				return
			}

		})
}
