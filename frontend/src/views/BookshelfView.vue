<template>
  <div class="bookshelf-view">
    <div class="library-header">
      <h2>{{ $t("bookshelf.title") }}</h2>
    </div>

    <!-- 本棚エリア -->
    <div class="bookcase">
      <div class="shelf-row">
        <!-- 既存のメモを本の背表紙として表示 -->
        <div
          v-for="memo in memos"
          :key="memo.ID"
          class="book-spine"
          :style="getBookStyle(memo)"
          @click="openBook(memo)"
          :title="memo.Title"
        >
          <div class="spine-decor top"></div>
          <div class="spine-title">{{ memo.Title }}</div>
          <div class="spine-date">
            {{ formatDateShort(memo.LinkedDate || memo.UpdatedAt) }}
          </div>
          <div class="spine-decor bottom"></div>
        </div>

        <!-- 新しいメモ作成用の「新しい本」 -->
        <div class="book-spine new-book" @click="openNewBook">
          <div class="spine-decor top"></div>
          <div class="spine-title">+</div>
          <div class="spine-date">{{ $t("bookshelf.new") }}</div>
          <div class="spine-decor bottom"></div>
        </div>
      </div>
      <div class="shelf-plank"></div>

      <!-- メモがない場合のメッセージ -->
      <div v-if="memos.length === 0" class="empty-shelf">
        {{ $t("bookshelf.empty") }}
      </div>
    </div>

    <!-- 本が開いた状態のオーバーレイ -->
    <div v-if="selectedMemo" class="book-overlay" @click.self="closeBook">
      <div class="open-book-container">
        <div class="open-book">
          <!-- 右ページ (静的): 本文の編集エリア -->
          <div class="page right-page">
            <div class="page-content">
              <textarea
                v-model="editingContent"
                class="handwritten-textarea"
                :placeholder="$t('bookshelf.contentPlaceholder')"
              ></textarea>
            </div>
            <div class="page-actions">
              <button class="btn-save" @click="saveMemo">
                {{ $t("bookshelf.save") }}
              </button>
            </div>
            <div class="page-number">- 1 -</div>
          </div>

          <!-- アニメーションする左ページアセンブリ (表紙と左ページ) -->
          <div class="page left-page-assembly">
            <!-- 表面: 表紙デザイン -->
            <div class="face front-face" :style="getCoverStyle(selectedMemo)">
              <div class="cover-design">
                <div class="spine-decor top"></div>
                <div class="cover-title">{{ selectedMemo.Title }}</div>
                <div class="spine-decor bottom"></div>
              </div>
            </div>

            <!-- 裏面: 左ページ (タイトルや日付情報) -->
            <div class="face back-face">
              <div class="page-content">
                <!-- 新規作成時はタイトル入力可能 -->
                <div v-if="!selectedMemo.ID" class="new-title-input-wrapper">
                  <input
                    v-model="selectedMemo.Title"
                    :placeholder="$t('bookshelf.titlePlaceholder')"
                    class="new-title-input"
                  />
                </div>
                <h2 v-else>{{ selectedMemo.Title }}</h2>

                <div class="meta-info">
                  <span
                    >📅
                    {{
                      formatDate(
                        selectedMemo.LinkedDate ||
                          selectedMemo.UpdatedAt ||
                          new Date(),
                      )
                    }}</span
                  >
                </div>
                <hr class="separator" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import axios from "axios";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const memos = ref<any[]>([]);
const selectedMemo = ref<any>(null); // 現在開いているメモ

const editingContent = ref("");

// メモ一覧の取得
async function fetchMemos() {
  try {
    const res = await axios.get("/api/memos");
    memos.value = res.data;
  } catch (e) {
    console.error(e);
  }
}

// 本の色パターン (アンティーク調)
const COLORS = [
  "#212121", // Very Dark Grey
  "#3e2723", // Dark Brown
  "#263238", // Dark Blue Grey
  "#1b0000", // Dark Red/Black
  "#00101b", // Deep Navy
  "#1c1c1c", // Black
  "#2d241e", // Cocoa
];

// 背表紙のスタイル生成 (色などをランダムに決定)
function getBookStyle(memo: any) {
  const seed = (memo.ID || 0) + (memo.Title ? memo.Title.length : 0);
  const color = COLORS[seed % COLORS.length];
  const height = "140px";
  const width = "45px";

  return {
    backgroundColor: color,
    height: height,
    width: width,
    backgroundImage: `linear-gradient(to right, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0.1) 10%, rgba(0,0,0,0) 50%, rgba(0,0,0,0.3) 100%)`,
  };
}

// 表紙のスタイル生成
function getCoverStyle(memo: any) {
  if (!memo.ID && memo.Title === "New Memo") {
    // 新規作成用の特別なスタイル (革のような質感)
    return {
      backgroundColor: "#5d4037",
      backgroundImage: `url("https://www.transparenttextures.com/patterns/wood-pattern.png")`,
      border: "2px solid #3e2723",
    };
  }
  const seed = (memo.ID || 0) + (memo.Title ? memo.Title.length : 0);
  const color = COLORS[seed % COLORS.length];
  return {
    backgroundColor: color,
    backgroundImage: `linear-gradient(to bottom right, rgba(255,255,255,0.1), rgba(0,0,0,0.2))`,
    boxShadow: "inset 0 0 20px rgba(0,0,0,0.5)",
  };
}

// 本を開く処理
function openBook(memo: any) {
  selectedMemo.value = { ...memo }; // コピーを作成して直接変更を防ぐ
  editingContent.value = memo.Content;
}

// 新規メモ作成として本を開く処理
function openNewBook() {
  selectedMemo.value = {
    Title: "New Memo",
    Content: "",
    ThemeColor: "#8d6e63",
  };
  editingContent.value = "";
}

// 本を閉じる処理
function closeBook() {
  selectedMemo.value = null;
  editingContent.value = "";
}

// メモの保存処理
async function saveMemo() {
  if (!selectedMemo.value) return;

  try {
    // 新規作成 (IDがない場合)
    if (!selectedMemo.value.ID) {
      await axios.post("/api/memos", {
        Title: selectedMemo.value.Title,
        Content: editingContent.value,
        // 新規の場合は今日の日付をデフォルトに
        LinkedDate: new Date().toISOString(),
        ThemeColor: "#8d6e63",
      });
    }
    // 更新
    else {
      await axios.put(`/api/memos/${selectedMemo.value.ID}`, {
        ...selectedMemo.value,
        Content: editingContent.value,
      });
    }

    alert("Saved!");
    closeBook();
    fetchMemos(); // リストを再取得
  } catch (e) {
    console.error(e);
    alert("Failed to save");
  }
}

// 日付フォーマットヘルパー
function formatDate(dateStr: string) {
  if (!dateStr) return "";
  return new Date(dateStr).toLocaleDateString();
}

function formatDateShort(dateStr: string) {
  if (!dateStr) return "";
  const d = new Date(dateStr);
  return `${d.getMonth() + 1}/${d.getDate()}`;
}

onMounted(() => {
  fetchMemos();
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Cinzel:wght@400;700&family=Nanum+Pen+Script&display=swap");

.bookshelf-view {
  min-height: 100vh;
  background-color: #3e2723;
  background-image: url("https://www.transparenttextures.com/patterns/wood-pattern.png");
  padding: 40px;
  color: white;
  font-family: "Noto Serif JP", serif;
}

.library-header h2 {
  text-align: center;
  font-family: "Cinzel", serif;
  font-size: 2.5rem;
  color: #d7ccc8;
  text-shadow: 2px 2px 4px black;
  margin-bottom: 40px;
  border-bottom: 2px solid #8d6e63;
  padding-bottom: 10px;
  display: inline-block;
}
.library-header {
  text-align: center;
}

.bookcase {
  max-width: 1000px;
  margin: 0 auto;
  background: #4e342e;
  box-shadow: inset 0 0 50px rgba(0, 0, 0, 0.8);
  padding: 20px 40px 0 40px;
  border: 20px solid #3e2723;
  border-radius: 4px;
}

.shelf-row {
  display: flex;
  align-items: flex-end;
  justify-content: flex-start;
  flex-wrap: wrap;
  gap: 2px;
  min-height: 350px;
  padding-bottom: 0;
}

.shelf-plank {
  height: 25px;
  background: linear-gradient(to bottom, #5d4037 0%, #3e2723 100%);
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.5);
  margin-bottom: 20px;
  border-radius: 2px;
}

.book-spine {
  position: relative;
  border-radius: 3px 3px 3px 3px;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.5);
  cursor: pointer;
  transition:
    transform 0.2s,
    margin 0.2s;
  color: #d4af37; /* Gold Text */
  font-family: "Noto Serif JP", serif;
  writing-mode: vertical-rl;
  text-orientation: upright; /* Upright characters for JP */
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  border-bottom: 1px solid rgba(0, 0, 0, 0.8);
}

.book-spine:hover {
  transform: translateY(-5px) scale(1.02);
  z-index: 10;
  box-shadow: 5px 5px 15px rgba(0, 0, 0, 0.6);
}

.spine-title {
  font-size: 1.2rem;
  letter-spacing: 2px;
  text-shadow: 1px 1px 1px black;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-height: 80%;
}

.spine-decor {
  position: absolute;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(to right, transparent, gold, transparent);
  opacity: 0.7;
}
.spine-decor.top {
  top: 30px;
  border-top: 1px solid rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  height: 10px;
  background: transparent;
  box-shadow:
    0 1px 0 rgba(255, 215, 0, 0.5) inset,
    0 -1px 0 rgba(255, 215, 0, 0.5) inset;
}
.spine-date {
  position: absolute;
  bottom: 5px; /* Above the bottom decor */
  font-size: 0.7rem;
  writing-mode: horizontal-tb; /* Horizontal text */
  width: 100%;
  text-align: center;
  color: #a1887f;
  font-family: sans-serif;
  letter-spacing: 0;
  z-index: 5;
}
.spine-decor.bottom {
  bottom: 30px;
  border-top: 1px solid rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  height: 10px;
  background: transparent;
  box-shadow:
    0 1px 0 rgba(255, 215, 0, 0.5) inset,
    0 -1px 0 rgba(255, 215, 0, 0.5) inset;
}

.empty-shelf {
  color: #a1887f;
  text-align: center;
  padding: 50px;
  font-style: italic;
}

/* Open Book Animation */
.book-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  perspective: 2500px; /* Higher perspective for better top-down feel */
}

.open-book-container {
  animation: containerAppear 0.5s ease-out forwards;
}

.open-book {
  width: 800px;
  height: 500px;
  position: relative;
  /* No transform-style on container needed unless we rotate it whole */
}

/* Base Page Style */
.page {
  position: absolute;
  top: 0;
  right: 0; /* Align to right side (since closed book is just the right stack) */
  width: 50%;
  height: 100%;
}

/* STATICE RIGHT PAGE */
.right-page {
  z-index: 1;
  background: linear-gradient(to right, #e8dfc5 0%, #fff 100%);
  border-radius: 0 5px 5px 0;
  box-shadow: 5px 5px 15px rgba(0, 0, 0, 0.3);
  padding: 40px;
  display: flex;
  flex-direction: column;
}

/* ANIMATING LEFT ASSEMBLY (The Cover/Left Page) */
.left-page-assembly {
  z-index: 10;
  transform-style: preserve-3d;
  transform-origin: left center; /* Pivot on Spine (Left edge of Right-aligned element) */
  animation: coverFlip 1.2s cubic-bezier(0.25, 1, 0.5, 1) forwards;
  /* Start closed (0deg) */
}

.face {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  backface-visibility: hidden;
  border-radius: 0 5px 5px 0; /* Matches Right Page shape initially */
}

/* FRONT FACE: The Cover (Visible initially) */
.front-face {
  z-index: 5;
  border-radius: 0 5px 5px 0; /* Matches right page shape */
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  /* Styling provided by inline style */
}
.cover-design {
  text-align: center;
  width: 80%;
  border: 2px solid rgba(255, 215, 0, 0.5);
  padding: 20px;
  background: rgba(0, 0, 0, 0.2);
}
.cover-title {
  font-family: "Cinzel", serif;
  font-size: 2rem;
  color: #e6c229; /* Gold */
  text-shadow: 2px 2px 4px black;
}

/* BACK FACE: The Inner Left Page (Visible after flip) */
.back-face {
  z-index: 4;
  transform: rotateY(180deg); /* Facing opposite way */
  background: linear-gradient(to left, #e8dfc5 0%, #fffbf0 100%);
  border-radius: 5px 0 0 5px; /* Rounded on LEFT now */
  padding: 40px;
  color: #333;
  box-shadow: inset -5px 0 10px rgba(0, 0, 0, 0.05);
}

.page-content {
  flex: 1;
  min-height: 0; /* Firefox flex fix */
  overflow-y: auto;
}

/* Animation Keyframes */
@keyframes containerAppear {
  from {
    transform: scale(0.8);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes coverFlip {
  0% {
    transform: rotateY(0deg);
  }
  100% {
    transform: rotateY(-180deg);
  }
}

.handwritten-textarea {
  width: 100%;
  height: 100%;
  border: none;
  background: transparent;
  font-family: "Nanum Pen Script", cursive, sans-serif;
  font-size: 1.5rem;
  line-height: 1.6;
  color: #2c3e50;
  outline: none;
  resize: none;
  background-image: linear-gradient(transparent 95%, #eee 95%);
  background-size: 100% 1.6em;
  padding: 14px 10px 0 55px; /* Moved further away from center (spine) */
  box-sizing: border-box;
}

.page-actions {
  margin-top: 10px;
  text-align: right;
}

.btn-save {
  background: #5d4037;
  color: #d7ccc8;
  border: 1px solid #3e2723;
  padding: 5px 15px;
  font-family: "Cinzel", serif;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.2s;
}
.btn-save:hover {
  background: #3e2723;
  color: white;
}

.separator {
  border: 0;
  border-top: 1px dashed #bbb;
  margin: 20px 0;
}

.page-number {
  position: absolute;
  bottom: 20px;
  width: 100%;
  text-align: center;
  color: #999;
  font-size: 0.8rem;
  left: 0;
}

@keyframes containerZoom {
  0% {
    transform: scale(0.1) translateY(200px);
    opacity: 0;
  }
  100% {
    transform: scale(1) translateY(0);
    opacity: 1;
  }
}

@keyframes openLeft {
  0% {
    transform: rotateY(-110deg);
    box-shadow: 10px 0 20px rgba(0, 0, 0, 0.5);
  }
  100% {
    transform: rotateY(0deg);
    box-shadow: -5px 5px 15px rgba(0, 0, 0, 0.2);
  }
}

@keyframes openRight {
  0% {
    transform: rotateY(110deg);
    box-shadow: -10px 0 20px rgba(0, 0, 0, 0.5);
  }
  100% {
    transform: rotateY(0deg);
    box-shadow: 5px 5px 15px rgba(0, 0, 0, 0.2);
  }
}

@keyframes contentFade {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.new-book {
  background-color: #5d4037;
  border: 1px dashed #d7ccc8;
  color: #ffe0b2;
  opacity: 0.8;
  transition: opacity 0.2s;
  height: 140px; /* Fixed height since we didn't use getBookStyle for it */
  width: 45px;
}
.new-book:hover {
  opacity: 1;
  transform: translateY(-5px) scale(1.02);
}
.new-title-input {
  width: 100%;
  font-size: 1.5rem;
  font-family: inherit;
  border: none;
  border-bottom: 2px dashed #aaa;
  background: transparent;
  text-align: right;
  outline: none;
}
</style>
