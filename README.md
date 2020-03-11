# 코로나19 푸시알리미(COVID19 Push)

[질병관리본부의 코로나19 홈페이지](http://ncov.mohw.go.kr/)를 주기적으로 확인하여, 발생 현황에 변동이 있거나 새 공지사항이 있는 경우 푸시알림으로 알려주는 웹서비스 입니다.

- [서비스 이용하기](https://covid19push.youngbin.xyz)  
- [Telegram 채널 구독](https://t.me/covid19push)
> 웹 푸시알림을 구독하려면 Service Worker, Web Push, Web Notification 사용 가능한 웹 브라우저가 필요합니다.
> Telegram 채널을 통해서도 알림을 받아보실 수 있습니다.

## 서비스 구축에 사용된 기술

- 프론트엔드
  - JavaScript
  - React.js
  - Material UI
  - Firebase SDK
  - Web Push, Web Notification
  - Service Worker
  - Google Analytics
- 백엔드
  - Go
  - Echo
  - Firebase Admin SDK
  - Gorm
  - SQLite
  - Chromedp
  - GoQuery
- 서비스 운영/배포
  - AWS EC2
  - Docker(Swarm Mode)
  - Caddy

## 프로젝트 관리
한영빈 (sukso96100@gmail.com)

## License

Copyright 2020-Present Youngbin Han and all COVID19-Push project contributors  
Licensed under MIT License

Check [LICENSE](LICENSE) for more

## 프로젝트에 기여하기
- 아래와 같은 경우, 이슈를 새로 생성해 주시거나 이메일로 연락해 주세요
  - 프로젝트에 관해 의견이 있는 경우
  - 사용 중 버그를 발견한 경우
- Pull Request 받습니다. 다만 승인까지 시간이 오래 걸릴 수 있습니다.
- Pull Request 열기 전 잘 작동하는지 테스트 후 제출하세요.
