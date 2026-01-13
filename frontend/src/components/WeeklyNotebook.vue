<template>
  <div class="notebook-container">
    <div class="notebook-header">
      <div class="month-label">
        <span>{{ currentMonthLabel }}</span>
      </div>
      <div class="year-label">{{ currentYear }}</div>
    </div>

    <div class="notebook-page">
      <!-- Time Header -->
      <div class="time-header">
        <div class="date-col-header"></div>
        <div class="timeline-col-header">
          <!-- Show labels every 3 hours -->
          <div
            v-for="h in hourSlots"
            :key="h"
            class="time-label"
            :class="{ visible: h % 3 === 0 }"
            :style="{ left: getTimePos(h) + '%' }"
          >
            {{ h % 3 === 0 ? h : "" }}
          </div>
        </div>
        <div class="memo-col-header">MEMO</div>
      </div>

      <div
        v-for="day in weekDays"
        :key="day.dateStr"
        class="day-row"
        :class="{
          'weekend-sun': day.dayOfWeek === 0,
          'weekend-sat': day.dayOfWeek === 6,
        }"
      >
        <!-- Date Column -->
        <div class="date-col">
          <div class="date-num">{{ day.dayOfMonth }}</div>
          <div class="day-name">{{ day.dayName }}</div>
        </div>

        <!-- Schedule Column (Timeline) -->
        <div class="timeline-col">
          <!-- Grid Lines -->
          <div class="grid-lines">
            <div
              v-for="h in hourSlots"
              :key="h"
              class="grid-line"
              :class="{ 'major-line': h % 3 === 0 }"
              :style="{ left: getTimePos(h) + '%' }"
            ></div>
          </div>

          <!-- Clickable Time Slots Layer -->
          <div class="time-slots-layer">
            <div
              v-for="h in hourSlots"
              :key="h"
              class="time-slot"
              :style="{
                left: getTimePos(h) + '%',
                width: 100 / TOTAL_HOURS + '%',
              }"
              @click="handleSlotClick(day.dateStr, h)"
              :title="h + ':00'"
            ></div>
          </div>

          <!-- Events (Horizontal Bars) -->
          <div class="events-container">
            <div
              v-for="event in day.events"
              :key="event.id"
              class="time-event"
              :style="getEventStyle(event)"
              :title="event.title"
            >
              {{ event.title }}
            </div>
          </div>
        </div>

        <!-- Memo Column -->
        <div class="memo-col" @click="$emit('memoClick', day.dateStr)">
          <div class="memo-lines">
            <div class="memo-line" v-for="n in 3" :key="n"></div>
          </div>
          <div class="memos-list">
            <div v-for="memo in day.memos" :key="memo.id" class="notebook-memo">
              {{ memo.title }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "vue-i18n";

const props = defineProps<{
  currentDate: Date;
  events: any[];
  memos: any[];
}>();

const emit = defineEmits(["dayClick", "timeClick", "memoClick"]);
const { t, locale } = useI18n();

// Config for Timeline
const START_HOUR = 6;
const END_HOUR = 24;
const TOTAL_HOURS = END_HOUR - START_HOUR;
// Slots 6..23 (for clicking).
const hourSlots = Array.from({ length: TOTAL_HOURS }, (_, i) => START_HOUR + i);

const getTimePos = (hour: number) => {
  return ((hour - START_HOUR) / TOTAL_HOURS) * 100;
};

const handleSlotClick = (dateStr: string, hour: number) => {
  const d = new Date(dateStr);
  d.setHours(hour, 0, 0, 0);
  // Emit ISO string to parent
  emit("timeClick", d.toISOString());
};

const getEventStyle = (event: any) => {
  const start = new Date(event.start);
  let h = start.getHours() + start.getMinutes() / 60;

  if (h < START_HOUR) h = START_HOUR;

  let duration = 1.0;

  let left = ((h - START_HOUR) / TOTAL_HOURS) * 100;
  let width = (duration / TOTAL_HOURS) * 100;

  if (left + width > 100) width = 100 - left;

  return {
    left: left + "%",
    width: width + "%",
    backgroundColor: event.color || "#4a90e2",
  };
};

// Formatting helpers
const getWeekDays = (baseDate: Date) => {
  const days = [];
  const startOfWeek = new Date(baseDate);
  const day = startOfWeek.getDay();
  const diff = startOfWeek.getDate() - day + (day === 0 ? -6 : 1);
  const monday = new Date(startOfWeek.setDate(diff));

  for (let i = 0; i < 7; i++) {
    const d = new Date(monday);
    d.setDate(monday.getDate() + i);

    const dayEvents = props.events.filter((e) =>
      isSameDay(new Date(e.start), d)
    );
    const dayMemos = props.memos.filter((m) => isSameDay(new Date(m.start), d));

    days.push({
      dateStr: d.toISOString(),
      dayOfMonth: d.getDate(),
      dayName: d
        .toLocaleDateString(locale.value, { weekday: "short" })
        .toUpperCase(),
      dayOfWeek: d.getDay(),
      events: dayEvents,
      memos: dayMemos,
    });
  }
  return days;
};

const isSameDay = (d1: Date, d2: Date) => {
  return (
    d1.getFullYear() === d2.getFullYear() &&
    d1.getMonth() === d2.getMonth() &&
    d1.getDate() === d2.getDate()
  );
};

const weekDays = computed(() => getWeekDays(props.currentDate));

const currentMonthLabel = computed(() => {
  return props.currentDate
    .toLocaleDateString(locale.value, { month: "long" })
    .toUpperCase();
});

const currentYear = computed(() => props.currentDate.getFullYear());
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Cinzel:wght@400;700&family=Noto+Serif+JP:wght@400;700&display=swap");

.notebook-container {
  background-color: #fffdf0; /* Cream paper */
  color: #333;
  font-family: "Noto Serif JP", serif;
  padding: 20px 20px 20px 50px; /* Extra left padding for binding */
  border: 1px solid #dcdcdc;
  /* Book Look: Rounded right corners, shadow */
  border-radius: 2px 10px 10px 2px;
  box-shadow: 5px 5px 15px rgba(0, 0, 0, 0.3),
    /* Main shadow */ 1px 0px 0px #aaa inset; /* Inner edge */
  max-width: 100%;
  position: relative;
}

/* Simulated Binding (Spiral or Spined) */
.notebook-container::before {
  content: "";
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  width: 40px;
  background: linear-gradient(
    to right,
    #5e4b35 0%,
    #8b6c4e 10%,
    #5e4b35 20%,
    #3e3223 100%
  ); /* Leather spine look */
  border-radius: 2px 0 0 2px;
  box-shadow: inset -2px 0 5px rgba(0, 0, 0, 0.5);
  z-index: 10;
}

.notebook-header {
  display: flex;
  justify-content: space-between;
  border-bottom: 2px solid #333;
  margin-bottom: 10px;
  padding-bottom: 5px;
  font-family: "Cinzel", serif;
}

.time-header {
  display: flex;
  font-size: 0.7rem;
  color: #888;
  border-bottom: 1px solid #ccc;
  padding-bottom: 5px;
  margin-bottom: 5px;
}
.date-col-header {
  width: 60px;
  flex-shrink: 0;
}
.timeline-col-header {
  flex-grow: 1;
  position: relative;
  height: 15px;
  margin-right: 10px;
}
.memo-col-header {
  width: 30%;
  flex-shrink: 0;
  text-align: center;
  font-weight: bold;
  border-left: 1px dotted #ccc;
}

.time-label {
  position: absolute;
  transform: translateX(-50%);
  bottom: 0;
  display: none;
}
.time-label.visible {
  display: block;
}

.day-row {
  display: flex;
  border-bottom: 1px solid #333;
  min-height: 80px;
}

/* Date Column */
.date-col {
  width: 60px;
  border-right: 1px solid #333;
  padding: 10px 5px;
  text-align: center;
  background: #fffdf5;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.date-num {
  font-size: 1.2rem;
  font-weight: bold;
}
.day-name {
  font-size: 0.8rem;
}

/* Timeline Column */
.timeline-col {
  flex-grow: 1;
  position: relative;
  border-right: 1px dotted #333;
  margin-right: 0;
  padding: 5px 0;
}
.grid-lines {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}
.grid-line {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 1px;
  background: #eee;
  border-right: 1px dashed #ddd;
}
.grid-line.major-line {
  background: #ccc;
  width: 2px;
}

/* Time Slots Layer - explicit clickable blocks */
.time-slots-layer {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
}
.time-slot {
  position: absolute;
  top: 0;
  bottom: 0;
  cursor: pointer;
  z-index: 5;
}
.time-slot:hover {
  background-color: rgba(74, 144, 226, 0.1);
}

/* Events Layer - on top of slots */
.events-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 10;
}
.time-event {
  position: absolute;
  height: 20px;
  background: #4a90e2;
  color: white;
  font-size: 0.7rem;
  border-radius: 4px;
  overflow: hidden;
  white-space: nowrap;
  padding: 0 4px;
  line-height: 20px;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.2);
  top: 50%;
  transform: translateY(-50%);
  pointer-events: auto;
}

/* Memo Column */
.memo-col {
  width: 30%;
  flex-shrink: 0;
  padding: 10px;
  position: relative;
  background: #fafcf5;
}
.memo-lines {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;
  pointer-events: none;
}
.memo-line {
  border-bottom: 1px dashed #dcdcdc;
  margin: 0 10px;
  width: calc(100% - 20px);
}
.notebook-memo {
  font-size: 0.85rem;
  color: #555;
  background: rgba(255, 255, 255, 0.8);
  margin-bottom: 4px;
  padding: 2px 4px;
  border-radius: 2px;
}

.weekend-sun .date-col {
  color: #e74c3c;
}
.weekend-sat .date-col {
  color: #3498db;
}
</style>
