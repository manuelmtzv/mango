<script setup lang="ts"></script>

<template>
  <section class="not-found">
    <figure class="not-found__logo">
      <h2 class="not-found__404">404</h2>

      <h2 class="not-found__title">Page not found.</h2>

      <picture>
        <source srcset="/404.webp" type="image/webp" />
        <img class="not-found__icon" src="/404.png" alt="Not found cat icon" />
      </picture>
    </figure>

    <p class="not-found__description">
      <span class="important text-3xl">Oh no!</span> It seems like you've
      stumbled upon uncharted territory. Our beloved yellow cat appears to be a
      bit perplexed and unsure about this
      <span class="important">404 page</span>.
    </p>

    <p class="not-found__description">
      Don't worry, though! We assure you that our feline friend is perfectly
      fine and eager to help you find your way back to your notes.
    </p>

    <RouterLink
      :to="{ path: '/app' }"
      class="button not-found__return navigation"
    >
      Return to Home
    </RouterLink>
  </section>
</template>

<style scoped lang="css">
.not-found {
  @apply flex flex-col gap-5 py-12 mb-4;
}
.not-found__logo {
  @apply flex flex-col gap-6 items-center;
}
.not-found__icon {
  @apply w-full max-w-[10rem];
}
.not-found__404 {
  @apply text-[7rem] [line-height:1] font-bold;
}
.not-found__title {
  @apply text-4xl font-semibold;
}
.not-found__description {
  @apply text-xl text-justify mx-auto max-w-xl;
}
.not-found__return {
  @apply mx-auto hover:bg-yellow-200 text-lg font-semibold;
}
</style>
