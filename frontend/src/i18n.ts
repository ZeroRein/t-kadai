import { createI18n } from 'vue-i18n'

const messages = {
  en: {
    nav: {
      calendar: 'Calendar',
      bookshelf: 'Bookshelf',
      lang: 'English'
    },
    calendar: {
      title: 'Calendar',
      month: 'Month',
      week: 'Week',
      today: 'Today'
    },
    bookshelf: {
      title: 'Bookshelf',
      empty: 'No memos found. Go to Calendar to add some!'
    },
    modal: {
      titlePlaceholder: 'Event/Memo Title',
      contentPlaceholder: 'Details or Memo Content',
      type: 'Type:',
      event: 'Event',
      memo: 'Daily Memo',
      save: 'Save',
      cancel: 'Cancel'
    },
    event: {
        alert: 'Event: '
    },
    home: {
      welcome: 'Welcome to Calendar',
      subtitle: 'Manage your schedule and memories in one place.',
      calendarDesc: 'View your monthly and weekly schedule.',
      bookshelfDesc: 'Browse your collection of daily memos.'
    }
  },
  ja: {
    nav: {
      calendar: 'カレンダー',
      bookshelf: '本棚',
      lang: '日本語'
    },
    calendar: {
      title: 'カレンダー',
      month: '月',
      week: '週',
      today: '今日'
    },
    bookshelf: {
      title: '本棚メモ',
      empty: 'メモがありません。カレンダーから追加してください！'
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
      bookshelfDesc: '日々のメモのコレクションを閲覧します。'
    }
  }
}

const i18n = createI18n({
  legacy: false, // you must set `false`, to use Composition API
  locale: 'en', // set locale
  fallbackLocale: 'en', // set fallback locale
  messages, // set locale messages
})

export default i18n
