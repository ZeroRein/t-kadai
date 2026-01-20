import { createI18n } from 'vue-i18n'

// 言語リソースの定義
// 'en' (英語) と 'ja' (日本語) の翻訳データをここで管理します。
const messages = {
  // 英語の翻訳データ
  en: {
    // ナビゲーションバーのテキスト
    nav: {
      calendar: 'Calendar',
      bookshelf: 'Bookshelf',
      lang: 'English',
      home: 'Home'
    },
    // カレンダービュー関連のテキスト
    calendar: {
      title: 'Calendar',
      month: 'Month',
      week: 'Week',
      today: 'Today',
      notebook: 'Notebook'
    },
    // 本棚（メモ一覧）ビュー関連のテキスト
    bookshelf: {
      title: 'Bookshelf',
      empty: 'No memos found. Go to Calendar to add some!',
      new: 'NEW',
      save: 'SAVE',
      titlePlaceholder: 'Title...',
      contentPlaceholder: 'Write your memo here...'
    },
    // モーダルダイアログ（予定/メモ作成）のテキスト
    modal: {
      titlePlaceholder: 'Event/Memo Title',
      contentPlaceholder: 'Details or Memo Content',
      type: 'Type:',
      event: 'Event',
      memo: 'Daily Memo',
      save: 'Save',
      cancel: 'Cancel'
    },
    // イベント表示用
    event: {
        alert: 'Event: '
    },
    // ホーム画面のテキストと天気予報の翻訳
    home: {
      welcome: 'Welcome to Calendar',
      subtitle: 'Manage your schedule and memories in one place.',
      calendarDesc: 'View your monthly and weekly schedule.',
      bookshelfDesc: 'Browse your collection of daily memos.',
      weatherTitle: 'Today in Tokyo',
      outfitAdvisor: '👗 Outfit Advisor',
      advisorLoading: 'Asking the stylist...',
      advisorError: 'Stylist is taking a nap.',
      // WMO Weather Codes の翻訳
      weather: {
        0: 'Clear Sky',
        1: 'Mainly Clear',
        2: 'Partly Cloudy',
        3: 'Overcast',
        45: 'Fog',
        48: 'Depositing Rime Fog',
        51: 'Light Drizzle',
        53: 'Moderate Drizzle',
        55: 'Dense Drizzle',
        61: 'Slight Rain',
        63: 'Moderate Rain',
        65: 'Heavy Rain',
        71: 'Slight Snow',
        73: 'Moderate Snow',
        75: 'Heavy Snow',
        58: 'Thunderstorm',
        95: 'Thunderstorm'
      }
    }
  },
  // 日本語の翻訳データ
  ja: {
    nav: {
      calendar: 'カレンダー',
      bookshelf: '本棚',
      lang: '日本語',
      home: 'ホーム'
    },
    calendar: {
      title: 'カレンダー',
      month: '月',
      week: '週',
      today: '今日',
      notebook: '手帳'
    },
    bookshelf: {
      title: '本棚メモ',
      empty: 'メモがありません。カレンダーから追加してください！',
      new: '新規',
      save: '保存',
      titlePlaceholder: 'タイトル...',
      contentPlaceholder: 'ここにメモを書いてください...'
    },
    modal: {
      titlePlaceholder: 'タイトルを入力',
      contentPlaceholder: '詳細またはメモ内容',
      type: '種類:',
      event: '予定',
      memo: 'デイリーメモ',
      save: '保存',
      cancel: 'キャンセル'
    },
    event: {
        alert: '予定: '
    },
    home: {
      welcome: 'Calendarへようこそ',
      subtitle: '予定と記録をシンプルに管理しましょう。',
      calendarDesc: '月次・週次のスケジュールを確認します。',
      bookshelfDesc: '日々のメモのコレクションを閲覧します。',
      weatherTitle: '今日の東京',
      outfitAdvisor: '👗 服装アドバイザー',
      advisorLoading: 'スタイリストに相談中...',
      advisorError: 'スタイリストは休憩中です。',
      weather: {
        0: '快晴',
        1: '晴れ',
        2: '一部曇り',
        3: '曇り',
        45: '霧',
        48: '着氷性の霧',
        51: '霧雨（弱）',
        53: '霧雨（中）',
        55: '霧雨（強）',
        61: '小雨',
        63: '雨',
        65: '大雨',
        71: '小雪',
        73: '雪',
        75: '大雪',
        58: '雷雨',
        95: '雷雨'
      }
    }
  }
}

// i18nインスタンスの作成
const i18n = createI18n({
  legacy: false, // Composition APIを使用するためにfalseに設定
  locale: 'en', // デフォルトの言語設定
  fallbackLocale: 'en', // 翻訳が見つからない場合のフォールバック言語
  messages, // 上記で定義した翻訳データ
})

export default i18n

