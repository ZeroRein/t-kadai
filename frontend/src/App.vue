<template>
  <header>
    <div class="logo">
      <span class="logo-icon">⚜️</span>
      <span class="logo-text">Calendar</span>
    </div>
    <nav>
      <!-- ナビゲーションリンク: カレンダーと本棚への遷移 -->
      <router-link to="/calendar">{{ $t("nav.calendar") }}</router-link>
      <router-link to="/bookshelf">{{ $t("nav.bookshelf") }}</router-link>
    </nav>
    <div class="header-actions">
      <!-- ホームボタン -->
      <router-link to="/" class="btn home-btn"
        >🏠 {{ $t("nav.home") }}</router-link
      >
      <!-- 言語切り替えボタン -->
      <button @click="toggleLang" class="btn">{{ $t("nav.lang") }}</button>
    </div>
  </header>
  <main>
    <!-- router-view: 現在のルートに応じたコンポーネントがここに表示されます -->
    <router-view />
  </main>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";

const { t, locale } = useI18n();

// 言語切り替え関数: 英語(en)と日本語(ja)をトグルします
const toggleLang = () => {
  locale.value = locale.value === "en" ? "ja" : "en";
};
</script>

<style>
/* Global Logo Styles */
.logo {
  display: flex;
  align-items: center;
  font-size: 1.8rem;
  color: #ffca28; /* Gold */
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.8);
}
.logo-icon {
  margin-right: 10px;
  font-size: 2rem;
  filter: drop-shadow(0 0 2px #000);
}
.logo-text {
  font-family: "Cinzel", serif;
  border-bottom: 2px solid #ffca28;
  padding-bottom: 2px;
}
</style>
