Docker ベース HTTPS カレンダー＆メモアプリ 基本設計書

1. プロジェクト概要

本ドキュメントは、Docker コンテナ上で動作し、HTTPS によるセキュアな通信（一時利用含む）を前提とした Web ベースのカレンダーアプリケーションの設計仕様をまとめたものである。

1.1 コンセプト

「予定管理と記録の統合」。単なるスケジュール帳ではなく、その日の出来事や準備事項を記すメモ帳（ジャーナル）機能を有機的に結合させる。また、蓄積された記録を「知識のコレクション」として視覚的に楽しめるようにする。

1.2 主な特徴

Docker 完全対応: docker-compose up 1 つで環境構築が完了する。

HTTPS 接続: Nginx リバースプロキシを導入し、SSL 終端を行う（オレオレ証明書または Let's Encrypt 想定）。

ビュー切り替え: 月表示（Monthly）と週表示（Weekly）のサポート。

i18n 対応: 日本語/英語の UI 切り替え。

カレンダー連動メモ: 特定の日付に紐づくメモと、日付に依存しないメモの両立。

本棚風メモ一覧: 蓄積したメモを視覚的に管理できる専用ビュー。

2. システムアーキテクチャ

2.1 技術スタック構成案

Frontend: Vue.js, FullCalendar (カレンダー UI 用ライブラリ推奨)

Backend: Go (Gin)

Database: PostgreSQL

Web Server / Proxy: Nginx (HTTPS 終端および静的ファイル配信)

Container Runtime: Docker / Docker Compose

2.2 コンテナ構成図

3 つのコンテナによる構成とする。

graph TD
User((User/Browser)) -- HTTPS (443) --> Nginx[Nginx Proxy Container]
Nginx -- HTTP --> App[App Container (Go)]
App -- TCP (5432) --> DB[(PostgreSQL Container)]
App -- Volume --> MemoFiles[Optional: Attachments]

3. 機能要件

3.1 カレンダー機能

表示モード

Monthly View: 1 ヶ月のカレンダーを表示。各セルにその日のイベント件数とメモの有無アイコンを表示。

Weekly View: 時間軸（バーティカル）またはリスト形式で 1 週間の予定を表示。

イベント管理 (CRUD)

タイトル、開始日時、終了日時、場所、説明、色分けタグ。

ナビゲーション

「今日」に戻るボタン。

前月/次月（前週/次週）への移動。

3.2 メモ帳機能

メモの種類

デイリーメモ: カレンダーの日付をクリックして作成。その日の日記や議事録として機能。カレンダー上からアクセス可能。

フリーメモ: 日付に紐付かない一般的なメモ（アイデア出し、To-Do リストなど）。

カレンダー連携

カレンダー上でメモが存在する日には「📝」アイコンを表示。

カレンダーの日付クリックで、その日のイベント一覧と共にデイリーメモのエディタが開く。

メモ一覧（本棚ビュー）

ホーム画面（またはナビゲーションバー）から遷移可能な独立した画面。

作成済みのメモ（デイリー・フリー両方）を一覧表示する。

デザイン: 本棚にノートが並んでいるようなビジュアル表現。

検索機能

全文検索により、過去のメモや予定を横断検索可能。

3.3 設定機能 (多言語対応)

言語切り替え

ヘッダーまたは設定画面に「English / 日本語」トグルスイッチを配置。

i18next 等のライブラリを使用し、JSON ファイルベースで文言を管理。

選択状態はブラウザの LocalStorage またはユーザー設定として DB に保存。

4. データベース設計 (ER 図イメージ)

4.1 Users (ユーザー)

Column Name

Type

Description

id

UUID

主キー

username

VARCHAR

ログイン ID

password_hash

VARCHAR

ハッシュ化パスワード

language

VARCHAR(5)

設定言語 ('ja' or 'en')

created_at

TIMESTAMP

作成日

4.2 Events (予定)

Column Name

Type

Description

id

UUID

主キー

user_id

UUID

外部キー(Users)

title

VARCHAR

予定タイトル

start_time

TIMESTAMP

開始日時

end_time

TIMESTAMP

終了日時

description

TEXT

詳細

color

VARCHAR

表示色

4.3 Memos (メモ)

Column Name

Type

Description

id

UUID

主キー

user_id

UUID

外部キー(Users)

title

VARCHAR

追加: メモのタイトル（本棚の背表紙に表示）

content

TEXT

メモ本文（Markdown 対応推奨）

linked_date

DATE

カレンダー日付 (NULL の場合はフリーメモ)

theme_color

VARCHAR

追加: 表紙/背表紙の色コード

created_at

TIMESTAMP

作成日時

updated_at

TIMESTAMP

更新日時

5. インフラ・Docker 構成詳細

5.1 ディレクトリ構造

/
├── docker-compose.yml
├── nginx/
│ ├── Dockerfile
│ ├── nginx.conf
│ └── certs/ (server.crt, server.key を配置)
├── backend/
│ ├── Dockerfile
│ └── src/
├── frontend/
│ ├── Dockerfile
│ └── src/
└── db-data/ (Volume mount)

5.2 HTTPS 設定 (Nginx)

「一時的な接続」を想定しているため、自己署名証明書（オレオレ証明書）を使用する運用フローとする。

nginx.conf 設定例:

server {
listen 443 ssl;
server_name localhost;

    ssl_certificate /etc/nginx/certs/server.crt;
    ssl_certificate_key /etc/nginx/certs/server.key;

    location / {
        proxy_pass http://app_container:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
    }

}

# HTTP から HTTPS へのリダイレクト

server {
listen 80;
server_name localhost;
return 301 https://$host$request_uri;
}

5.3 docker-compose.yml 構成案

version: '3.8'
services:
db:
image: postgres:15-alpine
environment:
POSTGRES_USER: user
POSTGRES_PASSWORD: password
POSTGRES_DB: calendar_db
volumes: - ./db-data:/var/lib/postgresql/data

app:
build: ./app
environment:
DATABASE_URL: postgres://user:password@db:5432/calendar_db
depends_on: - db

nginx:
build: ./nginx
ports: - "80:80" - "443:443"
volumes: - ./nginx/certs:/etc/nginx/certs
depends_on: - app

6. UI/UX デザイン仕様

6.1 画面レイアウト

共通ヘッダー:

ロゴ、言語切替スイッチ。

ナビゲーションメニュー: 「📅 Calendar」（カレンダー画面へ）、「📚 Bookshelf」（メモ一覧へ）。

カレンダー画面 (Calendar View):

メインエリア: 月/週カレンダー。

サイドパネル: 選択日の詳細・簡易メモ入力。

メモ一覧画面 (Bookshelf View):

レイアウト: 木製の本棚を模した背景またはグリッドレイアウト。

アイテム表現: 各メモを「本」や「ノート」のアイコンとして配置。

背表紙: title を縦書きまたは横書きで表示。

色: theme_color に基づき、本の色を変化させる。

並び順: 更新日順、または日付順（linked_date）。

カテゴリ分け: 「Daily Log（日付あり）」の棚と、「Ideas（フリーメモ）」の棚を段で分けるなどして整理。

6.2 インタラクションフロー

言語切り替え:

スイッチ操作で UI 言語（メニュー名、ボタンラベル）を即座に反映。

メモ作成・編集:

カレンダーから: 日付クリック → サイドパネルで入力（簡易）。

本棚から: 「＋ 新しいノート」クリック → フルスクリーンエディタで執筆。表紙の色やタイトルもここで設定可能。

メモ詳細閲覧:

本棚の「本」をクリックすると、本が開くアニメーションと共に詳細画面へ遷移。
