// Package givemoney 给予他人ATRI币
package givemoney

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/FloatTech/AnimeAPI/wallet"
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
		Help: "给予他人ATRI币(<1000)\n" +
			"- 给予 [QQ号]|[@xxx] [金额]",
		OnEnable: func(ctx *zero.Ctx) {
			ctx.Send("插件已启用")
		},
		OnDisable: func(ctx *zero.Ctx) {
			ctx.Send("插件已禁用")
		},
	})
	var (
		uidre    int64
		uidgi    int64
		nickname = ""
		err      error
		uidStr   string
		money    int
	)
	engine.OnPrefix("给予", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			param := strings.TrimSpace(ctx.State["args"].(string))
			re := regexp.MustCompile(`^[+]?\d+$`)
			money, err = strconv.Atoi(re.FindString(param))
			uidgi = ctx.Event.UserID
			if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
				uidStr = ctx.Event.Message[1].Data["qq"]
			} else {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("获取对方信息出错"))
				return
			}

			uidre, err = strconv.ParseInt(uidStr, 10, 64)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("QQ号处理出错"))
				return
			}
			nickname = ctx.GetThisGroupMemberInfo(
				uidgi,
				false,
			).Get("nickname").Str
			if nickname == "" {
				nickname = ctx.GetStrangerInfo(
					uidgi,
					false,
				).Get("nickname").Str

			}
			switch {
			case money > 1000:
				ctx.SendChain(message.At(uidgi), message.Text("一次最多给予1000ATRI币哦~"))
				return
			case money <= 0:
				ctx.SendChain(message.At(uidgi), message.Text("ATRI币数额必须大于0哦~"))
				return
			case 0 < money && money <= 1000:

				if wallet.GetWalletOf(uidgi) < money {
					ctx.SendChain(message.Reply(uidgi), message.At(uidgi), message.Text("ATRI币不足,发送签到获取吧~"))
					return
				}
				err := wallet.InsertWalletOf(uidgi, -money)
				if err != nil {
					ctx.SendChain(message.Reply(uidgi), message.At(uidgi), message.Text("ATRI币扣除失败  [ERROR at gm.go:51]:", err))
					return
				}
				err1 := wallet.InsertWalletOf(uidre, money)
				if err1 != nil {
					ctx.SendChain(message.Reply(uidgi), message.At(uidgi), message.Text("接受ATRI币出错啦,返回扣除的ATRI币  [ERROR at gm.go:55]:", err1))
					err := wallet.InsertWalletOf(uidgi, money)
					if err != nil {
						ctx.SendChain(message.Reply(uidgi), message.At(uidgi), message.Text("返回ATRI币失败,您的钱被吞掉了~  [ERROR at gm.go:58]:", err))
						return
					}
				}
				ctx.SendChain(message.Reply(uidgi), message.Text("成功给予"), message.At(uidre), message.Text(money, "ATRI币"))
				return
			default:
				ctx.SendChain(message.Text("发生未知错误"))
				return
			}
		})
}
