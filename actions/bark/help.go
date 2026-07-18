package bark

const helpText = `
bark

1. 生成加密信息
2. 发送 Bark 推送通知

Usage:
  nb-action bark <title> <message>
  nb-action bark --push <title> <message>

Arguments:
  title       通知标题
  message     通知内容

Examples:
  nb-action bark "服务器报警" "Docker 服务停止"
  nb-action bark --push "服务器报警" "Docker 服务停止"

Environment:
  BARK_SERVER
  BARK_DEVICE_KEY
  BARK_AES_KEY
  BARK_AES_IV
`
