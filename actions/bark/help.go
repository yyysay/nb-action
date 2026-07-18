package bark

const helpText = `
bark

发送 Bark 推送通知

Usage:
  nb-action bark <title> <message>

Arguments:
  title       通知标题
  message     通知内容

Examples:
  nb-action bark "服务器报警" "Docker 服务停止"

Environment:
  BARK_SERVER
  BARK_KEY
`
