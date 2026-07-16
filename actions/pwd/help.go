package pwd

const helpText = `
pwd

生成随机字符串或 WireGuard 密钥

Usage:
  nb-action pwd rand <length>
  nb-action pwd wg-keypair

Examples:
  nb-action pwd rand 32
  nb-action pwd wg-keypair

Commands:
  rand        生成随机字符串
  wg-keypair  生成 WireGuard 密钥对
`
