package serializers

type Welcome struct {
	IsOpen   bool   // 是否开启
	KillMe   int    // 自销毁时间
	Desc     bool   // 是否使用描述信息
	Pinned   bool   // 是否显示置顶消息
	Template string // 自定义模版
}

type NotRobot struct {
	IsOpen  bool // 是否开启机器人验证
	Timeout int  // 验证超时时间
}

type GroupSetting struct {
	Welcome  Welcome
	NotRobot NotRobot
}
