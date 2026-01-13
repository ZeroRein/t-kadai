<template>
  <header>
    <div class="logo">📅 Calendar</div>
    <nav>
      <router-link to="/calendar">{{ $t("nav.calendar") }}</router-link>
      <router-link to="/bookshelf">{{ $t("nav.bookshelf") }}</router-link>
    </nav>
    <div class="lang-switch">
      <button @click="toggleLang" class="btn">{{ $t("nav.lang") }}</button>
    </div>
  </header>
  <main>
    <router-view />
  </main>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";

const { t, locale } = useI18n();

const toggleLang = () => {
  locale.value = locale.value === "en" ? "ja" : "en";
};
</script>
