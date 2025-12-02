<script setup lang="ts"></script>

<template>
  <TransitionGroup name="list-transition-wrapper" tag="ul">
    <slot />
  </TransitionGroup>
</template>

<style>
.list-transition-wrapper-move,
.list-transition-wrapper-enter-active,
.list-transition-wrapper-leave-active {
  transition: all 350ms ease;
}

.list-transition-wrapper-enter-from,
.list-transition-wrapper-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.list-transition-wrapper-leave-active {
  position: absolute;
}
</style>
