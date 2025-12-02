<script setup lang="ts">
import FooterComponent from "./FooterComponent.vue";
import HeaderComponent from "./HeaderComponent.vue";
</script>

<template>
  <HeaderComponent />

  <main class="principle content">
    <router-view></router-view>
  </main>

  <FooterComponent />
</template>
