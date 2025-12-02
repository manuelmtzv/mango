<template>
  <footer class="footer">
    <div class="footer__container content flex-between-center">
      <div class="footer__left">
        <p class="footer__author">Created by <span>Manuel Martínez</span></p>
        <p>
          Made with
          <a class="footer__vue" href="https://vuejs.org/" target="_blank"
            >Vuejs</a
          >,
          <a
            class="footer__pinia"
            href="https://pinia.vuejs.org/"
            target="_blank"
            >Pinia</a
          >
          and
          <a
            class="footer__typescript"
            href="https://www.typescriptlang.org/"
            target="_blank"
            >TypeScript</a
          >.
        </p>
      </div>

      <p class="footer__right">
        Copyright © {{ new Date().getFullYear() }}
        <a
          class="important"
          href="https://github.com/manuelmtzv"
          target="_blank"
          >manumtzv</a
        >
      </p>
    </div>
  </footer>
</template>

<script setup lang="ts"></script>

<style scoped lang="css">
.footer {
  @apply py-4 border-gray-300 border-t text-sm sm:text-base;
}
.footer__container {
  @apply gap-8 items-start md:items-center;
}
.footer__left {
  @apply flex flex-col gap-2 font-semibold;
}
.footer__left p > a {
  @apply font-normal;
}
.footer__vue {
  @apply text-tech-vue-default;
}
.footer__pinia {
  @apply text-tech-pinia-default;
}
.footer__typescript {
  @apply text-tech-typescript-default font-bold;
}
.footer__vue,
.footer__pinia,
.footer__typescript {
  @apply !font-bold;
}
.footer__author {
  @apply font-semibold;
}
.footer__author > span {
  @apply font-normal;
}
.footer__right {
  @apply text-right;
}
</style>
