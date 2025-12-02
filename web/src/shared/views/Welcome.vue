<script setup lang="ts">
import pen from "@/assets/imgs/png/pen.png";
import devices from "@/assets/imgs/png/devices.png";
import tags from "@/assets/imgs/png/tags.png";
</script>

<template>
  <div class="py-6 mb-4 grid gap-8">
    <div class="welcome__container">
      <div class="welcome">
        <div class="welcome__content">
          <h1 class="text-3xl sm:text-4xl font-bold hyphens-manual">
            Hello, welcome to Mangocat&shy;notes!
          </h1>

          <div class="flex flex-col gap-2">
            <p class="welcome__message">
              Mangocatnotes is a free note management web app. It allows you to
              create and manage your notes from any device with a web browser.
            </p>

            <p class="welcome__message text-balance">
              Sign up now to start organizing your notes. Access them from
              anywhere! 🐱🥭
            </p>
          </div>
        </div>

        <div class="welcome__register">
          <picture>
            <source
              srcset="@/assets/imgs/webp/mangocat_computer.webp"
              type="image/webp"
            />
            <img
              class="max-w-[8rem] sm:max-w-[10rem] md:max-w-[14rem] w-full md:w-auto"
              src="@/assets/imgs/png/mangocat_computer.png"
              alt="Mangocat using a computer"
            />
          </picture>

          <RouterLink
            :to="{ name: 'register' }"
            class="block button submit mx-auto whitespace-nowrap font-semibold text-lg lg:hidden !px-5"
          >
            Start creating notes!
          </RouterLink>
        </div>
      </div>

      <RouterLink
        :to="{ name: 'register' }"
        class="hidden button submit mx-auto whitespace-nowrap text-lg font-semibold lg:block !px-4 mt-2"
      >
        Start creating notes!
      </RouterLink>
    </div>

    <hr />

    <div class="features__container">
      <h2 class="features__title">Highlighted Features:</h2>

      <div class="features">
        <figure class="feature">
          <picture>
            <source srcset="@/assets/imgs/webp/pen.webp" type="image/webp" />
            <img :src="pen" alt="Writing pencil icon" />
          </picture>

          <p class="feature__message">
            <span>Free Note Management:</span> Create and manage your notes
            without any cost.
          </p>
        </figure>

        <figure class="feature">
          <picture>
            <source
              srcset="@/assets/imgs/webp/devices.webp"
              type="image/webp"
            />
            <img :src="devices" alt="Writing pencil icon" />
          </picture>

          <p class="feature__message">
            <span>Universal Access:</span> Log in from any device with a web
            browser for seamless access.
          </p>
        </figure>

        <figure class="feature feature--recently-added">
          <picture>
            <source srcset="@/assets/imgs/webp/tags.webp" type="image/webp" />
            <img :src="tags" alt="Writing pencil icon" />
          </picture>

          <p class="feature__message">
            <span>Custom Tags:</span> Organize your notes efficiently by adding
            personalized tags.
          </p>
        </figure>
      </div>
    </div>
  </div>
</template>

<style scoped>
.welcome__container {
  @apply flex flex-col gap-4 mt-2 sm:mt-6;
}
.welcome {
  @apply flex flex-col gap-4 md:gap-8 md:flex-row items-center justify-between;
}
.welcome__content {
  @apply flex flex-col gap-4 max-w-4xl;
}
.welcome__register {
  @apply flex flex-col w-fit gap-4 items-center sm:flex-row sm:gap-8 md:flex-col md:gap-2;
}
.welcome__message {
  @apply text-lg font-medium text-justify sm:text-left sm:text-xl md:text-2xl hyphens-auto;
}

.features__container {
  @apply flex flex-col items-start gap-4;
}
.features__title {
  @apply text-2xl font-bold mb-4;
}
.features {
  @apply grid grid-cols-1 w-full gap-6 min-[500px]:grid-cols-2 md:grid-cols-3;
}
.feature {
  @apply flex flex-col flex-1 h-full items-center justify-center gap-4 place-self-center w-full p-4 bg-gray-100 rounded-lg shadow-md border min-[500px]:last-of-type:col-span-2 md:last-of-type:col-span-1;
}
.feature--recently-added {
  @apply relative after:content-['Now_available!'] after:block after:absolute after:top-0 after:right-0 after:z-10 after:bg-navigation-default after:rounded-bl-lg  after:shadow-md after:py-2 after:text-center after:text-lg after:font-bold after:opacity-90 after:w-fit after:px-4;
}
.feature img {
  @apply w-full max-w-[5rem] md:max-w-[7rem];
}
.feature__message {
  @apply text-center text-lg font-medium md:text-left;
}
.feature__message span {
  @apply font-bold;
}
</style>
