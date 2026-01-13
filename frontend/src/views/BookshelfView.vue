<template>
  <div class="bookshelf-view">
    <div class="library-header">
      <h2>{{ $t("bookshelf.title") }}</h2>
    </div>

    <!-- Bookshelf -->
    <div class="bookcase">
      <div class="shelf-row">
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
          <div class="spine-decor bottom"></div>
        </div>
      </div>
      <div class="shelf-plank"></div>

      <div v-if="memos.length === 0" class="empty-shelf">
        {{ $t("bookshelf.empty") }}
      </div>
    </div>

    <!-- Open Book Modal -->
    <div v-if="selectedMemo" class="book-overlay" @click.self="closeBook">
      <div class="open-book-container">
        <div class="open-book">
          <!-- Left Page -->
          <div class="page left-page">
            <div class="page-content">
              <h2>{{ selectedMemo.Title }}</h2>
              <div class="meta-info">
                <span
                  >📅
                  {{
                    formatDate(
                      selectedMemo.LinkedDate || selectedMemo.UpdatedAt
                    )
                  }}</span
                >
              </div>
              <hr class="separator" />
              <div class="decor-stamp">📝</div>
            </div>
          </div>

          <!-- Right Page -->
          <div class="page right-page">
            <div class="page-content">
              <p class="handwritten-text">{{ selectedMemo.Content }}</p>
            </div>
            <div class="page-number">- 1 -</div>
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
const selectedMemo = ref<any>(null);

async function fetchMemos() {
  try {
    const res = await axios.get("/api/memos");
    memos.value = res.data;
  } catch (e) {
    console.error(e);
  }
}

// Visual randomization for books (Dark Antique Style)
const COLORS = [
  "#212121", // Very Dark Grey
  "#3e2723", // Dark Brown
  "#263238", // Dark Blue Grey
  "#1b0000", // Dark Red/Black
  "#00101b", // Deep Navy
  "#1c1c1c", // Black
  "#2d241e", // Cocoa
];

function getBookStyle(memo: any) {
  const seed = (memo.ID || 0) + (memo.Title ? memo.Title.length : 0);
  const color = COLORS[seed % COLORS.length];
  // Size to match reference image (Small/Short books on large shelf)
  const height = 120 + (seed % 40) + "px";
  const width = 40 + (seed % 15) + "px";

  return {
    backgroundColor: color,
    height: height,
    width: width,
    backgroundImage: `linear-gradient(to right, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0.1) 10%, rgba(0,0,0,0) 50%, rgba(0,0,0,0.3) 100%)`, // Simpler dark binding
  };
}

function openBook(memo: any) {
  selectedMemo.value = memo;
}

function closeBook() {
  selectedMemo.value = null;
}

function formatDate(dateStr: string) {
  if (!dateStr) return "";
  return new Date(dateStr).toLocaleDateString();
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
  transition: transform 0.2s, margin 0.2s;
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
  box-shadow: 0 1px 0 rgba(255, 215, 0, 0.5) inset,
    0 -1px 0 rgba(255, 215, 0, 0.5) inset;
}
.spine-decor.bottom {
  bottom: 30px;
  border-top: 1px solid rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  height: 10px;
  background: transparent;
  box-shadow: 0 1px 0 rgba(255, 215, 0, 0.5) inset,
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
  perspective: 1500px;
}

.open-book-container {
  animation: bookEntrance 0.8s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
}

.open-book {
  width: 800px;
  height: 500px;
  background: #fdf5e6;
  display: flex;
  border-radius: 5px;
  box-shadow: 10px 10px 30px rgba(0, 0, 0, 0.5);
  position: relative;
}
/* Binding in the middle */
.open-book::before {
  content: "";
  position: absolute;
  top: 0;
  bottom: 0;
  left: 50%;
  width: 40px;
  margin-left: -20px;
  background: linear-gradient(
    to right,
    rgba(0, 0, 0, 0.1),
    rgba(0, 0, 0, 0.2) 40%,
    rgba(0, 0, 0, 0.2) 60%,
    rgba(0, 0, 0, 0.1)
  );
  z-index: 10;
  border-radius: 5px;
}

.page {
  flex: 1;
  padding: 40px;
  color: #333;
  position: relative;
  background: linear-gradient(to right, #fffbf0 0%, #fff 100%);
  box-shadow: inset -5px 0 10px rgba(0, 0, 0, 0.05); /* Paper curve */
}
.left-page {
  border-radius: 5px 0 0 5px;
  text-align: right;
  border-right: 1px solid #eee;
  background: linear-gradient(to left, #f5f5f5 0%, #fffbf0 100%);
}
.right-page {
  border-radius: 0 5px 5px 0;
  text-align: left;
}

.page-content {
  height: 100%;
  overflow-y: auto;
}

.handwritten-text {
  font-family: "Nanum Pen Script", cursive, sans-serif;
  font-size: 1.5rem;
  line-height: 1.6;
  color: #2c3e50;
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

@keyframes bookEntrance {
  0% {
    transform: translateY(100px) rotateX(60deg) scale(0.5);
    opacity: 0;
  }
  100% {
    transform: translateY(0) rotateX(0) scale(1);
    opacity: 1;
  }
}
</style>
