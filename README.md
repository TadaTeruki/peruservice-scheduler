# peruservice-scheduler

スケジュール管理サービスです。

APIの仕様は`OpenApi.yml`を参照してください。

## 環境変数

`MODE`: サーバーのモード (`DEVELOPMENT`|`PRODUCTION`)<br>
`SCHEDULER_PORT`: サーバーのポート番号<br>
`SCHEDULER_ALLOW_ORIGINS`: `PRODUCTION`モード時、サーバーとの通信を許可するオリジン (コンマ区切り) (例: `http://localhost:3000`)<br>
`PUBLIC_KEY_FILE`: JWT復号用の公開鍵へのファイルパス (例: `/path/to/key.pub`)<br>
`CONFIG_JSON_FILE`: サービスの設定ファイル名とパス (例: `/path/to/config.json`)<br>
`DB_DIRECTORY`: DBへのパス<br>
`DB_PORT`: DBのポート番号<br>
`DB_HOST`: DBのホスト<br>
`DB_USER`: DBのログイン用ユーザー名<br>
`DB_PASSWORD`: DBのログイン用パスワード<br>
`DB_NAME`: DBの名前<br>