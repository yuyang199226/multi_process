package main

type Status uint8

const (
	Init           Status = 0
	Running        Status = 1
	Upgrading      Status = 2
	Exiting        Status = 3
	Exited         Status = 4
	BizStart       Status = 5
	UpgradeFailure Status = 6
	Rollback       Status = 7
	SystemFault    Status = 8
	PowerOff       Status = 11
)

func (s Status) String() string {
	switch s {
	case Init:
		return "未启动"
	case Running:
		return "正常运行"
	case Upgrading:
		return "升级中"
	case Exiting:
		return "业务退出"
	case Exited:
		return "退出完成"
	case BizStart:
		return "业务启动"
	case UpgradeFailure:
		return "升级失败"
	case Rollback:
		return "回滚启动"
	case SystemFault:
		return "系统故障"
	case PowerOff:
		return "关机退出"

	}
	return ""
}
