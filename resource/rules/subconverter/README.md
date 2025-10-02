# Prism Rules for Subconverter

本目录包含适用于 [subconverter](https://github.com/tindy2013/subconverter) 的规则集文件。

## 使用方法

在 subconverter 的配置文件中添加以下规则集：

```ini
[custom]
; 启用规则生成
enable_rule_generator=true
overwrite_original_rules=true

; ADOBE (134 条规则)
ruleset=📦 ADOBE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/adobe.list
; AMAZON (184 条规则)
ruleset=📦 AMAZON,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/amazon.list
; APP STORE (2 条规则)
ruleset=📦 APP STORE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/app_store.list
; APPLE (26 条规则)
ruleset=📦 APPLE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/apple.list
; APPLE MUSIC (10 条规则)
ruleset=🎵 Apple Music,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/apple_music.list
; APPLE TV (8 条规则)
ruleset=📺 Apple TV,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/apple_tv.list
; BBC (19 条规则)
ruleset=📦 BBC,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bbc.list
; BILIBILI (130 条规则)
ruleset=📺 哔哩哔哩,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bilibili.list
; BILIBILI HK (1 条规则)
ruleset=📺 哔哩哔哩港澳台,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bilibili_hk.list
; BILIBILIINTL (1 条规则)
ruleset=📺 哔哩哔哩,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bilibiliintl.list
; BINANCE (12 条规则)
ruleset=📦 BINANCE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/binance.list
; BLOOMBERG (22 条规则)
ruleset=📦 BLOOMBERG,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bloomberg.list
; CCTV (37 条规则)
ruleset=📦 CCTV,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cctv.list
; CLAUDE (4 条规则)
ruleset=🤖 Claude,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/claude.list
; CLOUDFLARE (65 条规则)
ruleset=📦 CLOUDFLARE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cloudflare.list
; CNN (6 条规则)
ruleset=📦 CNN,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cnn.list
; CRYPTO.COM (199 条规则)
ruleset=📦 CRYPTO.COM,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/crypto.com.list
; CRYPTOCURRENCY (4 条规则)
ruleset=📦 CRYPTOCURRENCY,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cryptocurrency.list
; DAZN (16 条规则)
ruleset=📦 DAZN,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/dazn.list
; DEVELOP (9 条规则)
ruleset=📦 DEVELOP,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/develop.list
; DIGITALOCEAN (4 条规则)
ruleset=📦 DIGITALOCEAN,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/digitalocean.list
; DIRECT (118312 条规则)
ruleset=🎯 全球直连,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/direct.list
; DISCORD (29 条规则)
ruleset=💬 Discord,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/discord.list
; DISNEY (1 条规则)
ruleset=🎬 Disney+,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/disney.list
; DISNEY+ (174 条规则)
ruleset=📦 DISNEY+,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/disney+.list
; DOCKER (7 条规则)
ruleset=🐳 Docker,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/docker.list
; DOUYU (13 条规则)
ruleset=📦 DOUYU,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/douyu.list
; DROPBOX (17 条规则)
ruleset=📦 DROPBOX,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/dropbox.list
; EBAY (44 条规则)
ruleset=📦 EBAY,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/ebay.list
; EPIC (1 条规则)
ruleset=📦 EPIC,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/epic.list
; EPIC GAMES (14 条规则)
ruleset=📦 EPIC GAMES,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/epic_games.list
; FACEBOOK (570 条规则)
ruleset=📘 Facebook,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/facebook.list
; GAME (1 条规则)
ruleset=📦 GAME,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/game.list
; GEMINI (13 条规则)
ruleset=🤖 Gemini,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/gemini.list
; GITHUB (36 条规则)
ruleset=💻 GitHub,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/github.list
; GITLAB (6 条规则)
ruleset=💻 GitLab,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/gitlab.list
; GOOGLE (700 条规则)
ruleset=📦 GOOGLE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/google.list
; GOOGLE DRIVE (6 条规则)
ruleset=📦 GOOGLE DRIVE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/google_drive.list
; HBO (47 条规则)
ruleset=📦 HBO,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/hbo.list
; HEROKU (12 条规则)
ruleset=📦 HEROKU,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/heroku.list
; HIMALAYA (17 条规则)
ruleset=📦 HIMALAYA,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/himalaya.list
; HULU (59 条规则)
ruleset=📦 HULU,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/hulu.list
; ICLOUD (59 条规则)
ruleset=📦 ICLOUD,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/icloud.list
; INSTAGRAM (1 条规则)
ruleset=📷 Instagram,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/instagram.list
; IQIYI (58 条规则)
ruleset=📺 爱奇艺,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/iqiyi.list
; IQIYI HK (17 条规则)
ruleset=📦 IQIYI HK,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/iqiyi_hk.list
; IQIYIINTL (1 条规则)
ruleset=📦 IQIYIINTL,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/iqiyiintl.list
; LINE (23 条规则)
ruleset=📦 LINE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/line.list
; LINKEDIN (12 条规则)
ruleset=📦 LINKEDIN,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/linkedin.list
; MICROSOFT (658 条规则)
ruleset=📦 MICROSOFT,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/microsoft.list
; NETEASE MUSIC (28 条规则)
ruleset=📦 NETEASE MUSIC,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/netease_music.list
; NETFLIX (38 条规则)
ruleset=🎬 Netflix,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/netflix.list
; NEW YORK TIMES (16 条规则)
ruleset=📦 NEW YORK TIMES,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/new_york_times.list
; OKX (3 条规则)
ruleset=📦 OKX,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/okx.list
; ONEDRIVE (17 条规则)
ruleset=📦 ONEDRIVE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/onedrive.list
; OPENAI (31 条规则)
ruleset=🤖 OpenAI,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/openai.list
; ORACLE (1 条规则)
ruleset=📦 ORACLE,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/oracle.list
; PANDORA (2 条规则)
ruleset=📦 PANDORA,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/pandora.list
; PAYPAL (247 条规则)
ruleset=📦 PAYPAL,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/paypal.list
; PLAYSTATION (116 条规则)
ruleset=📦 PLAYSTATION,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/playstation.list
; PRIME VIDEO (26 条规则)
ruleset=📦 PRIME VIDEO,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/prime_video.list
; PROXY (37451 条规则)
ruleset=🚀 节点选择,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/proxy.list
; REDDIT (8 条规则)
ruleset=📦 REDDIT,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/reddit.list
; REJECT (143870 条规则)
ruleset=🛡️ 广告拦截,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/reject.list
; SCHOLAR (229 条规则)
ruleset=📦 SCHOLAR,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/scholar.list
; SHOPIFY (8 条规则)
ruleset=📦 SHOPIFY,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/shopify.list
; SOUNDCLOUD (2 条规则)
ruleset=📦 SOUNDCLOUD,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/soundcloud.list
; SPOTIFY (24 条规则)
ruleset=🎵 Spotify,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/spotify.list
; STEAM (55 条规则)
ruleset=📦 STEAM,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/steam.list
; STEAMCN (1 条规则)
ruleset=📦 STEAMCN,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/steamcn.list
; TELEGRAM (47 条规则)
ruleset=💬 Telegram,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/telegram.list
; TENCENT VIDEO (51 条规则)
ruleset=📺 腾讯视频,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/tencent_video.list
; TESTFLIGHT (2 条规则)
ruleset=📦 TESTFLIGHT,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/testflight.list
; THREADS (1 条规则)
ruleset=📦 THREADS,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/threads.list
; TIKTOK (31 条规则)
ruleset=📹 TikTok,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/tiktok.list
; TWITCH (22 条规则)
ruleset=📦 TWITCH,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/twitch.list
; TWITTER (34 条规则)
ruleset=🐦 Twitter,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/twitter.list
; VERCEL (27 条规则)
ruleset=⚡ Vercel,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/vercel.list
; VIMEO (16 条规则)
ruleset=📦 VIMEO,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/vimeo.list
; WIKIPEDIA (13 条规则)
ruleset=📦 WIKIPEDIA,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/wikipedia.list
; YOUKU (35 条规则)
ruleset=📦 YOUKU,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/youku.list
; YOUTUBE (184 条规则)
ruleset=📹 YouTube,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/youtube.list

; GEOIP 规则
ruleset=🎯 全球直连,[]GEOIP,CN
ruleset=🐟 漏网之鱼,[]FINAL
```

## 可用规则集

| 规则集 | 规则数量 | 下载链接 |
|--------|----------|----------|
| ADOBE | 134 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/adobe.list) |
| AMAZON | 184 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/amazon.list) |
| APP STORE | 2 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/app_store.list) |
| APPLE | 26 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/apple.list) |
| APPLE MUSIC | 10 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/apple_music.list) |
| APPLE TV | 8 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/apple_tv.list) |
| BBC | 19 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bbc.list) |
| BILIBILI | 130 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bilibili.list) |
| BILIBILI HK | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bilibili_hk.list) |
| BILIBILIINTL | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bilibiliintl.list) |
| BINANCE | 12 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/binance.list) |
| BLOOMBERG | 22 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/bloomberg.list) |
| CCTV | 37 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cctv.list) |
| CLAUDE | 4 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/claude.list) |
| CLOUDFLARE | 65 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cloudflare.list) |
| CNN | 6 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cnn.list) |
| CRYPTO.COM | 199 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/crypto.com.list) |
| CRYPTOCURRENCY | 4 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/cryptocurrency.list) |
| DAZN | 16 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/dazn.list) |
| DEVELOP | 9 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/develop.list) |
| DIGITALOCEAN | 4 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/digitalocean.list) |
| DIRECT | 118312 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/direct.list) |
| DISCORD | 29 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/discord.list) |
| DISNEY | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/disney.list) |
| DISNEY+ | 174 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/disney+.list) |
| DOCKER | 7 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/docker.list) |
| DOUYU | 13 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/douyu.list) |
| DROPBOX | 17 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/dropbox.list) |
| EBAY | 44 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/ebay.list) |
| EPIC | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/epic.list) |
| EPIC GAMES | 14 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/epic_games.list) |
| FACEBOOK | 570 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/facebook.list) |
| GAME | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/game.list) |
| GEMINI | 13 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/gemini.list) |
| GITHUB | 36 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/github.list) |
| GITLAB | 6 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/gitlab.list) |
| GOOGLE | 700 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/google.list) |
| GOOGLE DRIVE | 6 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/google_drive.list) |
| HBO | 47 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/hbo.list) |
| HEROKU | 12 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/heroku.list) |
| HIMALAYA | 17 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/himalaya.list) |
| HULU | 59 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/hulu.list) |
| ICLOUD | 59 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/icloud.list) |
| INSTAGRAM | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/instagram.list) |
| IQIYI | 58 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/iqiyi.list) |
| IQIYI HK | 17 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/iqiyi_hk.list) |
| IQIYIINTL | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/iqiyiintl.list) |
| LINE | 23 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/line.list) |
| LINKEDIN | 12 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/linkedin.list) |
| MICROSOFT | 658 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/microsoft.list) |
| NETEASE MUSIC | 28 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/netease_music.list) |
| NETFLIX | 38 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/netflix.list) |
| NEW YORK TIMES | 16 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/new_york_times.list) |
| OKX | 3 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/okx.list) |
| ONEDRIVE | 17 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/onedrive.list) |
| OPENAI | 31 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/openai.list) |
| ORACLE | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/oracle.list) |
| PANDORA | 2 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/pandora.list) |
| PAYPAL | 247 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/paypal.list) |
| PLAYSTATION | 116 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/playstation.list) |
| PRIME VIDEO | 26 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/prime_video.list) |
| PROXY | 37451 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/proxy.list) |
| REDDIT | 8 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/reddit.list) |
| REJECT | 143870 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/reject.list) |
| SCHOLAR | 229 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/scholar.list) |
| SHOPIFY | 8 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/shopify.list) |
| SOUNDCLOUD | 2 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/soundcloud.list) |
| SPOTIFY | 24 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/spotify.list) |
| STEAM | 55 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/steam.list) |
| STEAMCN | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/steamcn.list) |
| TELEGRAM | 47 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/telegram.list) |
| TENCENT VIDEO | 51 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/tencent_video.list) |
| TESTFLIGHT | 2 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/testflight.list) |
| THREADS | 1 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/threads.list) |
| TIKTOK | 31 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/tiktok.list) |
| TWITCH | 22 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/twitch.list) |
| TWITTER | 34 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/twitter.list) |
| VERCEL | 27 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/vercel.list) |
| VIMEO | 16 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/vimeo.list) |
| WIKIPEDIA | 13 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/wikipedia.list) |
| YOUKU | 35 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/youku.list) |
| YOUTUBE | 184 | [下载](https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/youtube.list) |

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
