<template>
  <div class="home-view">
    <section class="hero">
      <h1>{{ $t("home.welcome") }}</h1>
      <p>{{ $t("home.subtitle") }}</p>
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

    <!-- Potential "Today's Summary" section could go here -->
  </div>
</template>

<script setup lang="ts">
// Simple static home for now.
// Could fetch "Today's Events" from API if desired.
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
  color: #4a90e2;
  margin-bottom: 0.5rem;
}

.hero p {
  color: #666;
  font-size: 1.2rem;
  margin-bottom: 3rem;
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
  transition: transform 0.2s, box-shadow 0.2s;
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
