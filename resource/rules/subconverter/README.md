# Prism Rules for Subconverter

本目录包含适用于 [subconverter](https://github.com/tindy2013/subconverter) 的规则集文件。

## 使用方法

在 subconverter 的配置文件中添加以下规则集：

```ini
[custom]
; 启用规则生成
enable_rule_generator=true
overwrite_original_rules=true

; Reject (143870 条规则)
ruleset=🛡️ 广告拦截,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/reject.list
; Direct (118312 条规则)
ruleset=🎯 全球直连,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/direct.list
; Proxy (37451 条规则)
ruleset=🚀 节点选择,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/proxy.list
; Adobe (134 条规则)
ruleset=📦 Adobe,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/adobe.list
; Amazon (184 条规则)
ruleset=📦 Amazon,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/amazon.list
; Apple (26 条规则)
ruleset=📦 Apple,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/apple.list
; BBC (19 条规则)
ruleset=📦 BBC,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bbc.list
; Bilibili (130 条规则)
ruleset=📺 哔哩哔哩,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bilibili.list
; Binance (12 条规则)
ruleset=📦 Binance,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/binance.list
; Bloomberg (22 条规则)
ruleset=📦 Bloomberg,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bloomberg.list
; CCTV (37 条规则)
ruleset=📦 CCTV,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cctv.list
; CNN (6 条规则)
ruleset=📦 CNN,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cnn.list
; Claude (4 条规则)
ruleset=🤖 Claude,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/claude.list
; Cloudflare (65 条规则)
ruleset=📦 Cloudflare,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cloudflare.list
; Cryptocurrency (4 条规则)
ruleset=📦 Cryptocurrency,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cryptocurrency.list
; DAZN (16 条规则)
ruleset=📦 DAZN,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/dazn.list
; DigitalOcean (4 条规则)
ruleset=📦 DigitalOcean,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/digitalocean.list
; Discord (29 条规则)
ruleset=💬 Discord,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/discord.list
; Disney+ (1 条规则)
ruleset=📦 Disney+,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/disney+.list
; Docker (7 条规则)
ruleset=🐳 Docker,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/docker.list
; Douyu (13 条规则)
ruleset=📦 Douyu,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/douyu.list
; Dropbox (17 条规则)
ruleset=📦 Dropbox,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/dropbox.list
; Epic Games (1 条规则)
ruleset=📦 Epic Games,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/epic_games.list
; Facebook (570 条规则)
ruleset=📘 Facebook,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/facebook.list
; Gemini (13 条规则)
ruleset=🤖 Gemini,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/gemini.list
; GitHub (36 条规则)
ruleset=💻 GitHub,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/github.list
; GitLab (6 条规则)
ruleset=💻 GitLab,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/gitlab.list
; Google (700 条规则)
ruleset=📦 Google,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/google.list
; HBO (47 条规则)
ruleset=📦 HBO,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/hbo.list
; Heroku (12 条规则)
ruleset=📦 Heroku,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/heroku.list
; Himalaya (17 条规则)
ruleset=📦 Himalaya,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/himalaya.list
; Hulu (59 条规则)
ruleset=📦 Hulu,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/hulu.list
; Instagram (1 条规则)
ruleset=📷 Instagram,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/instagram.list
; Line (23 条规则)
ruleset=📦 Line,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/line.list
; LinkedIn (12 条规则)
ruleset=📦 LinkedIn,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/linkedin.list
; Microsoft (658 条规则)
ruleset=📦 Microsoft,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/microsoft.list
; Netflix (38 条规则)
ruleset=🎬 Netflix,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/netflix.list
; OKX (3 条规则)
ruleset=📦 OKX,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/okx.list
; OneDrive (17 条规则)
ruleset=📦 OneDrive,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/onedrive.list
; OpenAI (31 条规则)
ruleset=🤖 OpenAI,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/openai.list
; Pandora (2 条规则)
ruleset=📦 Pandora,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/pandora.list
; PayPal (247 条规则)
ruleset=📦 PayPal,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/paypal.list
; PlayStation (116 条规则)
ruleset=📦 PlayStation,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/playstation.list
; Reddit (8 条规则)
ruleset=📦 Reddit,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/reddit.list
; Scholar (229 条规则)
ruleset=📦 Scholar,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/scholar.list
; Shopify (8 条规则)
ruleset=📦 Shopify,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/shopify.list
; SoundCloud (2 条规则)
ruleset=📦 SoundCloud,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/soundcloud.list
; Spotify (24 条规则)
ruleset=🎵 Spotify,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/spotify.list
; Steam (55 条规则)
ruleset=📦 Steam,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/steam.list
; Telegram (47 条规则)
ruleset=💬 Telegram,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/telegram.list
; TestFlight (2 条规则)
ruleset=📦 TestFlight,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/testflight.list
; Threads (1 条规则)
ruleset=📦 Threads,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/threads.list
; TikTok (31 条规则)
ruleset=📹 TikTok,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/tiktok.list
; Twitch (22 条规则)
ruleset=📦 Twitch,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/twitch.list
; Twitter (34 条规则)
ruleset=🐦 Twitter,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/twitter.list
; Unknown (565 条规则)
ruleset=📦 Unknown,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/unknown.list
; Vercel (27 条规则)
ruleset=⚡ Vercel,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/vercel.list
; Vimeo (16 条规则)
ruleset=📦 Vimeo,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/vimeo.list
; Wikipedia (13 条规则)
ruleset=📦 Wikipedia,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/wikipedia.list
; YouTube (184 条规则)
ruleset=📹 YouTube,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/youtube.list
; Youku (35 条规则)
ruleset=📦 Youku,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/youku.list
; eBay (44 条规则)
ruleset=📦 eBay,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/ebay.list
; iCloud (59 条规则)
ruleset=📦 iCloud,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/icloud.list
; iQIYI (58 条规则)
ruleset=📺 爱奇艺,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/iqiyi.list

; GEOIP 规则
ruleset=🎯 全球直连,[]GEOIP,CN
ruleset=🐟 漏网之鱼,[]FINAL
```

## 可用规则集

| 规则集 | 规则数量 | 下载链接 |
|--------|----------|----------|
| Reject | 143870 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/reject.list) |
| Direct | 118312 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/direct.list) |
| Proxy | 37451 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/proxy.list) |
| Adobe | 134 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/adobe.list) |
| Amazon | 184 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/amazon.list) |
| Apple | 26 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/apple.list) |
| BBC | 19 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bbc.list) |
| Bilibili | 130 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bilibili.list) |
| Binance | 12 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/binance.list) |
| Bloomberg | 22 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bloomberg.list) |
| CCTV | 37 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cctv.list) |
| CNN | 6 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cnn.list) |
| Claude | 4 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/claude.list) |
| Cloudflare | 65 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cloudflare.list) |
| Cryptocurrency | 4 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cryptocurrency.list) |
| DAZN | 16 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/dazn.list) |
| DigitalOcean | 4 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/digitalocean.list) |
| Discord | 29 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/discord.list) |
| Disney+ | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/disney+.list) |
| Docker | 7 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/docker.list) |
| Douyu | 13 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/douyu.list) |
| Dropbox | 17 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/dropbox.list) |
| Epic Games | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/epic_games.list) |
| Facebook | 570 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/facebook.list) |
| Gemini | 13 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/gemini.list) |
| GitHub | 36 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/github.list) |
| GitLab | 6 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/gitlab.list) |
| Google | 700 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/google.list) |
| HBO | 47 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/hbo.list) |
| Heroku | 12 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/heroku.list) |
| Himalaya | 17 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/himalaya.list) |
| Hulu | 59 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/hulu.list) |
| Instagram | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/instagram.list) |
| Line | 23 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/line.list) |
| LinkedIn | 12 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/linkedin.list) |
| Microsoft | 658 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/microsoft.list) |
| Netflix | 38 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/netflix.list) |
| OKX | 3 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/okx.list) |
| OneDrive | 17 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/onedrive.list) |
| OpenAI | 31 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/openai.list) |
| Pandora | 2 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/pandora.list) |
| PayPal | 247 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/paypal.list) |
| PlayStation | 116 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/playstation.list) |
| Reddit | 8 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/reddit.list) |
| Scholar | 229 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/scholar.list) |
| Shopify | 8 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/shopify.list) |
| SoundCloud | 2 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/soundcloud.list) |
| Spotify | 24 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/spotify.list) |
| Steam | 55 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/steam.list) |
| Telegram | 47 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/telegram.list) |
| TestFlight | 2 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/testflight.list) |
| Threads | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/threads.list) |
| TikTok | 31 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/tiktok.list) |
| Twitch | 22 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/twitch.list) |
| Twitter | 34 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/twitter.list) |
| Unknown | 565 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/unknown.list) |
| Vercel | 27 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/vercel.list) |
| Vimeo | 16 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/vimeo.list) |
| Wikipedia | 13 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/wikipedia.list) |
| YouTube | 184 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/youtube.list) |
| Youku | 35 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/youku.list) |
| eBay | 44 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/ebay.list) |
| iCloud | 59 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/icloud.list) |
| iQIYI | 58 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/iqiyi.list) |

## 规则格式

所有规则文件遵循 subconverter 的 `.list` 格式：

```
DOMAIN,example.com
DOMAIN-SUFFIX,example.com
DOMAIN-KEYWORD,keyword
IP-CIDR,192.168.0.0/16,no-resolve
IP-CIDR6,2001:db8::/32,no-resolve
```

## 更新频率

规则集每日自动更新，来源于：
- [blackmatrix7/ios_rule_script](https://github.com/blackmatrix7/ios_rule_script)
- [Loyalsoldier/clash-rules](https://github.com/Loyalsoldier/clash-rules)
- [ACL4SSR/ACL4SSR](https://github.com/ACL4SSR/ACL4SSR)

## 许可证

本项目采用 GPL-3.0 许可证。详见 [LICENSE](../../LICENSE) 文件。
