<template>
  <div class="calendar-view">
    <div class="controls">
      <h2>{{ currentTitle }}</h2>
      <div class="view-switcher">
        <button
          class="btn"
          :class="{ active: viewMode === 'notebook' }"
          @click="switchView('notebook')"
        >
          Notebook
        </button>
        <button
          class="btn"
          :class="{ active: viewMode === 'calendar' }"
          @click="switchView('calendar')"
        >
          {{ $t("calendar.month") }}
        </button>
      </div>
    </div>

    <!-- Custom Notebook View -->
    <div v-if="viewMode === 'notebook'">
      <div class="notebook-nav">
        <button class="btn sm" @click="addWeeks(-1)">&lt;</button>
        <button class="btn sm" @click="resetDate">
          {{ $t("calendar.today") }}
        </button>
        <button class="btn sm" @click="addWeeks(1)">&gt;</button>
      </div>
      <WeeklyNotebook
        :currentDate="notebookDate"
        :events="calendarEvents"
        :memos="memoList"
        @dayClick="handleNotebookDateClick"
        @timeClick="handleNotebookTimeClick"
        @memoClick="handleNotebookMemoClick"
      />
    </div>

    <!-- Standard FullCalendar (Wrapped for visual contrast) -->
    <div v-else class="calendar-paper">
      <FullCalendar ref="fullCalendar" :options="calendarOptions" />
    </div>

    <!-- Simple Modal for creating Event/Memo -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <h3>
          {{ newEvent.type === "event" ? $t("modal.event") : $t("modal.memo") }}
        </h3>

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

        <!-- Start/End Inputs for Events -->
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
import { ref, reactive, onMounted, computed, watch, nextTick } from "vue";
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
});

// Separate state for Events and Memos
const calendarEvents = ref<any[]>([]);
const memoList = ref<any[]>([]);
const fullCalendarEvents = computed(() => [
  ...calendarEvents.value,
  ...memoList.value,
]);

const viewMode = ref("notebook");
const notebookDate = ref(new Date());
const fullCalendar = ref(null);

// Computed title
const currentTitle = computed(() => {
  if (viewMode.value === "notebook") {
    const d = notebookDate.value;
    return d.getFullYear() + " / " + (d.getMonth() + 1);
  }
  return t("calendar.title");
});

// Navigation
const addWeeks = (n: number) => {
  const d = new Date(notebookDate.value);
  d.setDate(d.getDate() + n * 7);
  notebookDate.value = d;
};
const resetDate = () => {
  notebookDate.value = new Date();
};

// Helpers for datetime-local format (YYYY-MM-DDTHH:mm)
const toLocalISO = (d: Date) => {
  const offset = d.getTimezoneOffset() * 60000;
  const local = new Date(d.getTime() - offset);
  return local.toISOString().slice(0, 16);
};

const handleNotebookDateClick = (dateStr: string) => {
  // Default 9am - 10am for day clicks
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

const handleNotebookTimeClick = (isoDateStr: string) => {
  const d = new Date(isoDateStr);
  const end = new Date(d);
  end.setHours(end.getHours() + 1); // Default 1 hour duration

  selectedDate.value = isoDateStr;
  newEvent.title = "";
  newEvent.content = "";
  newEvent.type = "event";
  newEvent.start = toLocalISO(d);
  newEvent.end = toLocalISO(end);
  showModal.value = true;
};

const handleNotebookMemoClick = (dateStr: string) => {
  const d = new Date(dateStr);
  d.setHours(9, 0, 0, 0);

  selectedDate.value = dateStr;
  newEvent.title = "";
  newEvent.content = "";
  newEvent.type = "memo"; // Default to Memo
  newEvent.start = toLocalISO(d);
  newEvent.end = "";
  showModal.value = true;
};

const switchView = (mode: string) => {
  viewMode.value = mode;
  if (mode === "calendar") {
    nextTick(() => {
      if (fullCalendar.value) (fullCalendar.value as any).getApi().updateSize();
    });
  }
};

// Calendar Options (FullCalendar)
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
      (arg.event.extendedProps.description || "")
  );
}

function closeModal() {
  showModal.value = false;
}

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
    // Memo
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

async function fetchData() {
  try {
    const [evRes, memoRes] = await Promise.all([
      axios.get("/api/events"),
      axios.get("/api/memos"),
    ]);

    calendarEvents.value = evRes.data.map((e: any) => ({
      id: e.ID,
      title: e.Title,
      start: e.StartTime,
      color: "#4a90e2",
      extendedProps: { description: e.Description },
    }));

    memoList.value = memoRes.data
      .filter((m: any) => m.LinkedDate)
      .map((m: any) => ({
        id: m.ID,
        title: "📝 " + m.Title,
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
  background: #eee;
  color: #333;
}
.view-switcher .btn.active {
  background: #4a90e2;
  color: white;
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
