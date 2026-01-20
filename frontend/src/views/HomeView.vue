<template>
  <div class="home-view">
    <section class="hero">
      <h1>{{ $t("home.welcome") }}</h1>
      <p>{{ $t("home.subtitle") }}</p>
    </section>

    <!-- Weather Widget -->
    <section class="weather-widget" v-if="weatherData">
      <div class="weather-header">
        <h2>{{ $t("home.weatherTitle") }}</h2>
        <span class="weather-date">{{ currentDate }}</span>
      </div>
      <div class="weather-content">
        <div class="weather-icon">{{ weatherIcon }}</div>
        <div class="weather-info">
          <p class="temperature">{{ weatherData.current.temperature_2m }}°C</p>
          <p class="condition">{{ weatherDescription }}</p>
        </div>
      </div>

      <!-- AI Advice Section -->
      <div class="weather-advice">
        <h3>{{ $t("home.outfitAdvisor") }}</h3>
        <p v-if="loadingAdvice" class="advice-loading">
          {{ $t("home.advisorLoading") }}
        </p>
        <p v-else-if="advice" class="advice-text">{{ advice }}</p>
        <p v-else class="advice-error">{{ $t("home.advisorError") }}</p>
      </div>
    </section>
    <section class="weather-widget loading" v-else>
      <p>Loading weather...</p>
    </section>

    <div class="actions">
      <router-link to="/calendar" class="card">
        <div class="icon">📅</div>
        <h3>{{ $t("nav.calendar") }}</h3>
        <p>{{ $t("home.calendarDesc") }}</p>
      </router-link>

      <router-link to="/bookshelf" class="card">
        <div class="icon">📚</div>
        <h3>{{ $t("nav.bookshelf") }}</h3>
        <p>{{ $t("home.bookshelfDesc") }}</p>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue";
import { useI18n } from "vue-i18n";

const { t, locale } = useI18n();

const weatherData = ref<any>(null);
const advice = ref<string | null>(null);
const loadingAdvice = ref<boolean>(false);

const currentDate = computed(() => {
  return new Date().toLocaleDateString(locale.value, {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  });
});

const weatherCodeMap: Record<number, string> = {
  0: "☀️",
  1: "🌤️",
  2: "⛅",
  3: "☁️",
  45: "🌫️",
  48: "🌫️",
  51: "🌧️",
  53: "🌧️",
  55: "🌧️",
  61: "☔",
  63: "☔",
  65: "☔",
  71: "❄️",
  73: "❄️",
  75: "❄️",
  95: "⚡",
};

const weatherIcon = computed(() => {
  if (!weatherData.value) return "❓";
  const code = weatherData.value.current.weather_code;
  return weatherCodeMap[code] || "🌈";
});

const weatherDescription = computed(() => {
  if (!weatherData.value) return "Unknown";
  const code = weatherData.value.current.weather_code;
  return t(`home.weather.${code}`);
});

const fetchAdvice = async (data: any) => {
  loadingAdvice.value = true;
  advice.value = null;
  try {
    const payload = {
      weather_code: data.current.weather_code,
      temperature: data.current.temperature_2m,
      description: weatherDescription.value,
      language: locale.value, // Send current language
    };

    const response = await fetch("http://localhost:8080/api/advice", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!response.ok) throw new Error("API call failed");

    const result = await response.json();
    advice.value = result.advice;
  } catch (error) {
    console.error("Failed to fetch advice:", error);
    advice.value = t("home.advisorError");
  } finally {
    loadingAdvice.value = false;
  }
};

const fetchWeather = async () => {
  try {
    const response = await fetch(
      "https://api.open-meteo.com/v1/forecast?latitude=35.6895&longitude=139.6917&current=temperature_2m,weather_code&timezone=Asia%2FTokyo",
    );
    const data = await response.json();
    weatherData.value = data;
    // Fetch Advice after weather
    fetchAdvice(data);
  } catch (error) {
    console.error("Failed to fetch weather:", error);
  }
};

onMounted(() => {
  fetchWeather();
});

// Watch for language changes to re-fetch advice in new language
watch(locale, () => {
  if (weatherData.value) {
    fetchAdvice(weatherData.value);
  }
});
</script>

<style scoped>
.home-view {
  max-width: 800px;
  margin: 0 auto;
  text-align: center;
  padding: 4rem 2rem;
}

.hero h1 {
  font-size: 2.5rem;
  color: #f5f5f5; /* Light text for contrast */
  margin-bottom: 0.5rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.hero p {
  color: #e0e0e0; /* Light text for contrast */
  font-size: 1.2rem;
  margin-bottom: 3rem;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
}

.weather-widget {
  background: rgba(255, 255, 255, 0.9); /* Slightly more opaque */
  border: 2px solid #8d6e63;
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 3rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  display: inline-block;
  min-width: 300px;
  max-width: 400px;
}

.weather-header h2 {
  margin: 0;
  color: #3e2723;
  font-size: 1.5rem;
}

.weather-date {
  color: #5d4037;
  font-size: 0.9rem;
  margin-bottom: 1rem;
  display: block;
}

.weather-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1.5rem;
  margin-top: 1rem;
}

.weather-icon {
  font-size: 3.5rem;
}

.temperature {
  font-size: 2rem;
  font-weight: bold;
  color: #e65100;
  margin: 0;
}

.condition {
  color: #5d4037;
  margin: 0;
}

.weather-advice {
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px dashed #8d6e63;
}

.weather-advice h3 {
  font-size: 1.1rem;
  color: #3e2723;
  margin-bottom: 0.5rem;
}

.advice-text {
  font-style: italic;
  color: #5d4037;
}

.advice-loading {
  color: #888;
  font-size: 0.9rem;
}

.advice-error {
  color: #d32f2f;
  font-size: 0.9rem;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 2rem;
  flex-wrap: wrap;
}

.card {
  background: white;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 2rem;
  width: 250px;
  text-decoration: none;
  color: inherit;
  transition:
    transform 0.2s,
    box-shadow 0.2s;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.1);
}

.icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.card h3 {
  margin: 0 0 0.5rem;
  color: #333;
}

.card p {
  color: #888;
  font-size: 0.9rem;
  margin: 0;
}
</style>
