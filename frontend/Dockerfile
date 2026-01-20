# Node.js 20 (Alpine Linux版) をベースイメージとして使用
FROM node:20-alpine

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係定義ファイル (package.json, package-lock.json) をコピー
COPY package*.json ./

# 依存パッケージのインストール
RUN npm install

# ソースコード全体をコピー
COPY . .

# 開発サーバーのポート (5173) を公開
EXPOSE 5173

# 開発サーバーを起動 (ホストからのアクセスを許可するために --host を付与)
CMD ["npm", "run", "dev", "--", "--host"]
