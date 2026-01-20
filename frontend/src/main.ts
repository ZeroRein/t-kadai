import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import HomeView from './views/HomeView.vue'
import CalendarView from './views/CalendarView.vue'
import BookshelfView from './views/BookshelfView.vue'
import i18n from './i18n'

import './style.css' // グローバルスタイルのインポート

// ルーターの設定: URLパスとコンポーネントの対応を定義
const router = createRouter({
  history: createWebHistory(), // ブラウザの履歴管理モードを使用
  routes: [
    { path: '/', component: HomeView }, // ホーム画面
    { path: '/calendar', component: CalendarView }, // カレンダー画面
    { path: '/bookshelf', component: BookshelfView }, // 本棚（メモ一覧）画面
  ]
})

// Vueアプリケーションの作成とマウント
const app = createApp(App)
app.use(router) // ルータープラグインの使用
app.use(i18n)   // i18n（多言語対応）プラグインの使用
app.mount('#app') // index.htmlの #app 要素にマウント
