<template>
  <div class="calendar-view">
    <div class="controls">
      <!-- 現在の表示年月またはタイトル -->
      <h2>{{ currentTitle }}</h2>
      <div class="view-switcher">
        <!-- 手帳ビューへの切り替えボタン -->
        <button
          class="btn"
          :class="{ active: viewMode === 'notebook' }"
          @click="switchView('notebook')"
        >
          {{ $t("calendar.notebook") }}
        </button>
        <!-- カレンダービューへの切り替えボタン -->
        <button
          class="btn"
          :class="{ active: viewMode === 'calendar' }"
          @click="switchView('calendar')"
        >
          {{ $t("calendar.month") }}
        </button>
      </div>
    </div>

    <!-- 手帳（週間）ビューエリア -->
    <div v-if="viewMode === 'notebook'">
      <div class="notebook-nav">
        <!-- 週の移動ナビゲーション -->
        <button class="btn sm" @click="addWeeks(-1)">&lt;</button>
        <button class="btn sm" @click="resetDate">
          {{ $t("calendar.today") }}
        </button>
        <button class="btn sm" @click="addWeeks(1)">&gt;</button>
      </div>
      <!-- WeeklyNotebookコンポーネント: 週間スケジュールとメモを表示 -->
      <WeeklyNotebook
        :currentDate="notebookDate"
        :events="calendarEvents"
        :memos="memoList"
        @dayClick="handleNotebookDateClick"
        @timeClick="handleNotebookTimeClick"
        @memoClick="handleNotebookMemoClick"
        @createMemo="handleCreateMemo"
        @updateMemo="handleUpdateMemo"
      />
    </div>

    <!-- 月間カレンダーエリア (FullCalendar) -->
    <div v-else class="calendar-paper">
      <FullCalendar ref="fullCalendar" :options="calendarOptions" />
    </div>

    <!-- イベント/メモ作成・編集用モーダル -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <h3>
          {{ newEvent.type === "event" ? $t("modal.event") : $t("modal.memo") }}
        </h3>

        <!-- タイプ選択 (予定 or メモ) -->
        <div style="margin-bottom: 10px">
          <label>{{ $t("modal.type") }}: </label>
          <select v-model="newEvent.type">
            <option value="event">{{ $t("modal.event") }}</option>
            <option value="memo">{{ $t("modal.memo") }}</option>
          </select>
        </div>

        <input
          v-model="newEvent.title"
          :placeholder="$t('modal.titlePlaceholder')"
          style="display: block; width: 100%; margin-bottom: 10px; padding: 5px"
        />

        <!-- 予定の場合の開始・終了時刻入力 -->
        <div
          v-if="newEvent.type === 'event'"
          style="margin-bottom: 10px; display: flex; gap: 10px"
        >
          <div style="flex: 1">
            <label style="font-size: 0.8rem; font-weight: bold">Start</label>
            <input
              type="datetime-local"
              v-model="newEvent.start"
              style="width: 100%; padding: 5px"
            />
          </div>
          <div style="flex: 1">
            <label style="font-size: 0.8rem; font-weight: bold">End</label>
            <input
              type="datetime-local"
              v-model="newEvent.end"
              style="width: 100%; padding: 5px"
            />
          </div>
        </div>

        <textarea
          v-model="newEvent.content"
          :placeholder="$t('modal.contentPlaceholder')"
          rows="5"
          style="display: block; width: 100%; margin-bottom: 10px; padding: 5px"
        ></textarea>

        <button class="btn" @click="saveItem">{{ $t("modal.save") }}</button>
        <button
          class="btn"
          style="background: #ccc; margin-left: 10px"
          @click="closeModal"
        >
          {{ $t("modal.cancel") }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, nextTick } from "vue";
import FullCalendar from "@fullcalendar/vue3";
import dayGridPlugin from "@fullcalendar/daygrid";
import interactionPlugin from "@fullcalendar/interaction";
import jaLocale from "@fullcalendar/core/locales/ja";
import axios from "axios";
import { useI18n } from "vue-i18n";
import WeeklyNotebook from "../components/WeeklyNotebook.vue";

const { t, locale } = useI18n();

const showModal = ref(false);
const selectedDate = ref("");
const newEvent = reactive({
  title: "",
  content: "",
  type: "event",
  start: "",
  end: "",
  // note: ID is missing here for updates, simple create focused
});

// イベントとメモの状態管理
const calendarEvents = ref<any[]>([]);
const memoList = ref<any[]>([]);

// FullCalendar表示用にイベントとメモをマージ
const fullCalendarEvents = computed(() => [
  ...calendarEvents.value,
  ...memoList.value,
]);

// ビューモード: 'notebook' (手帳) または 'calendar' (月間カレンダー)
const viewMode = ref("notebook");
const notebookDate = ref(new Date());
const fullCalendar = ref(null);

// タイトルの計算 (手帳モード時は年月を表示)
const currentTitle = computed(() => {
  if (viewMode.value === "notebook") {
    const d = notebookDate.value;
    return d.getFullYear() + " / " + (d.getMonth() + 1);
  }
  return t("calendar.title");
});

// 週ナビゲーション
const addWeeks = (n: number) => {
  const d = new Date(notebookDate.value);
  d.setDate(d.getDate() + n * 7);
  notebookDate.value = d;
};
const resetDate = () => {
  notebookDate.value = new Date();
};

// datetime-local 形式への変換ヘルパー (YYYY-MM-DDTHH:mm)
const toLocalISO = (d: Date) => {
  const offset = d.getTimezoneOffset() * 60000;
  const local = new Date(d.getTime() - offset);
  return local.toISOString().slice(0, 16);
};

// ノートブックでの日付クリックハンドラ
const handleNotebookDateClick = (dateStr: string) => {
  // デフォルトで9時から10時を設定
  const d = new Date(dateStr);
  d.setHours(9, 0, 0, 0);
  const end = new Date(d);
  end.setHours(10, 0, 0, 0);

  selectedDate.value = dateStr;
  newEvent.title = "";
  newEvent.content = "";
  newEvent.type = "event";
  newEvent.start = toLocalISO(d);
  newEvent.end = toLocalISO(end);
  showModal.value = true;
};

// ノートブックでの時間枠クリックハンドラ
const handleNotebookTimeClick = (isoDateStr: string) => {
  const d = new Date(isoDateStr);
  const end = new Date(d);
  end.setHours(end.getHours() + 1); // 1時間枠

  selectedDate.value = isoDateStr;
  newEvent.title = "";
  newEvent.content = "";
  newEvent.type = "event";
  newEvent.start = toLocalISO(d);
  newEvent.end = toLocalISO(end);
  showModal.value = true;
};

// ノートブックでのメモクリックハンドラ
const handleNotebookMemoClick = (dateStr: string) => {
  const d = new Date(dateStr);
  d.setHours(9, 0, 0, 0);

  selectedDate.value = dateStr;
  newEvent.title = "";
  newEvent.content = "";
  newEvent.type = "memo"; // デフォルトでメモを選択
  newEvent.start = toLocalISO(d);
  newEvent.end = "";
  showModal.value = true;
};

// クイックメモ作成ハンドラ
const handleCreateMemo = async (payload: { date: string; content: string }) => {
  try {
    await axios.post("/api/memos", {
      Title: payload.content, // クイックメモでは内容をタイトルとして使用
      Content: payload.content,
      LinkedDate: new Date(payload.date).toISOString(),
      ThemeColor: "#ffba00",
    });
    fetchData(); // リストを更新
  } catch (e) {
    console.error(e);
    alert("Failed to save memo");
  }
};

// メモ更新ハンドラ
const handleUpdateMemo = async (payload: {
  id: number;
  content: string;
  date: string;
}) => {
  try {
    // 既存のメモを探して更新 (PUT)
    // 本来はIDで取得して更新するのが確実ですが、ここではリストから検索
    const memo = memoList.value.find((m) => m.id === payload.id);
    if (memo) {
      await axios.put(`/api/memos/${payload.id}`, {
        ID: payload.id,
        Title: payload.content,
        Content: payload.content,
        LinkedDate: new Date(payload.date).toISOString(),
        ThemeColor: "#ffba00",
      });
      fetchData();
    }
  } catch (e) {
    console.error(e);
  }
};

// ビュー切り替え
const switchView = (mode: string) => {
  viewMode.value = mode;
  // FullCalendarのサイズ再計算
  if (mode === "calendar") {
    nextTick(() => {
      if (fullCalendar.value) (fullCalendar.value as any).getApi().updateSize();
    });
  }
};

// FullCalendarの設定
const calendarOptions = reactive({
  plugins: [dayGridPlugin, interactionPlugin],
  initialView: "dayGridMonth",
  locales: [jaLocale],
  locale: locale,
  dateClick: (arg: any) => handleNotebookDateClick(arg.dateStr),
  eventClick: handleEventClick,
  events: fullCalendarEvents,
  headerToolbar: {
    left: "prev,next today",
    center: "title",
    right: "dayGridMonth",
  },
});

function handleEventClick(arg: any) {
  alert(
    t("event.alert") +
      arg.event.title +
      "\n" +
      (arg.event.extendedProps.description || ""),
  );
}

function closeModal() {
  showModal.value = false;
}

// アイテム（予定/メモ）の保存
async function saveItem() {
  if (newEvent.type === "event") {
    await axios.post("/api/events", {
      Title: newEvent.title,
      Description: newEvent.content,
      StartTime: new Date(newEvent.start).toISOString(),
      EndTime: new Date(newEvent.end).toISOString(),
      Color: "#4a90e2",
    });
  } else {
    // メモ
    await axios.post("/api/memos", {
      Title: newEvent.title,
      Content: newEvent.content,
      LinkedDate: new Date(newEvent.start).toISOString(),
      ThemeColor: "#ffba00",
    });
  }
  closeModal();
  fetchData();
}

// データ取得 (予定とメモを並列で取得)
async function fetchData() {
  try {
    const [evRes, memoRes] = await Promise.all([
      axios.get("/api/events"),
      axios.get("/api/memos"),
    ]);

    // カレンダー用イベント形式に変換
    calendarEvents.value = evRes.data.map((e: any) => ({
      id: e.ID,
      title: e.Title,
      start: e.StartTime,
      color: "#4a90e2",
      extendedProps: { description: e.Description },
    }));

    // メモをカレンダーイベント形式に変換 (日付があるもの)
    memoList.value = memoRes.data
      .filter((m: any) => m.LinkedDate)
      .map((m: any) => ({
        id: m.ID,
        title: m.Title,
        start: m.LinkedDate,
        color: "#ff9f43",
        extendedProps: { description: m.Content },
      }));
  } catch (e) {
    console.error("Failed to fetch data", e);
  }
}

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.calendar-view {
  max-width: 1200px;
  margin: 0 auto;
}
.controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.view-switcher .btn {
  margin-left: 10px;
  background: linear-gradient(to bottom, #5d4037, #3e2723);
  color: #d7ccc8;
  border: 1px solid #2e1c16;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 2px 4px rgba(0, 0, 0, 0.4);
  font-family: "Cinzel", serif;
  text-transform: uppercase;
  letter-spacing: 1px;
}
.view-switcher .btn.active {
  background: linear-gradient(to bottom, #ffb300, #ff6f00);
  color: #3e2723;
  border: 1px solid #ff6f00;
  box-shadow:
    inset 0 0 10px rgba(0, 0, 0, 0.2),
    0 1px 2px rgba(0, 0, 0, 0.3);
  font-weight: bold;
  text-shadow: 0 1px 0 rgba(255, 255, 255, 0.4);
}
.notebook-nav {
  display: flex;
  justify-content: center;
  margin-bottom: 15px;
  gap: 10px;
}
.btn.sm {
  padding: 0.2rem 0.8rem;
  font-size: 0.9rem;
}

/* Paper Style for FullCalendar */
.calendar-paper {
  background-color: #fffdf0; /* Cream paper */
  padding: 20px;
  border-radius: 4px;
  box-shadow: 2px 2px 10px rgba(0, 0, 0, 0.1);
  font-family: "Noto Serif JP", serif;
}
</style>
