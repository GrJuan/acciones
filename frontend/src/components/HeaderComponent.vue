<template>
  <header
    v-motion="headerMotion"
    style="margin-top: 1.5rem"
    class="mb-6 p-4 bg-gray-900 rounded-lg shadow-lg border border-gray-700 flex items-center justify-between"
  >
    <div>
      <h1 class="text-3xl font-bold text-gray-100 flex items-center">
        <Icon name="mdi:finance" class="mr-3 text-green-400 text-4xl" />
        Dashboard de Acciones
      </h1>
      <p class="text-gray-400 text-sm">
        Visualiza tus finanzas y estadísticas en un ambiente profesional.
      </p>

      <!-- Pequeñas estadísticas clave -->
      <div class="mt-2 flex space-x-4 text-gray-300 text-sm">
        <div class="flex items-center">
          <Icon name="mdi:calendar-month-outline" class="mr-2 text-blue-400" />
          <span>{{ currentDate }}</span>
        </div>
        <div class="flex items-center">
          <Icon
            name="mdi:chart-line"
            class="mr-2 text-yellow-400 animate-bounce"
          />
          <span class="relative flex items-center">
            <span class="mr-2">Mercado en vivo</span>
            <span class="relative flex h-3 w-3">
              <span
                class="animate-ping absolute inline-flex h-full w-full rounded-full bg-[red] opacity-75"
              ></span>
              <span
                class="relative inline-flex rounded-full h-3 w-3 bg-[red]"
              ></span>
            </span>
          </span>
        </div>
      </div>
    </div>
  </header>
</template>

<script lang="ts" setup>
import { ref } from "vue";

const currentDate = ref(
  new Date().toLocaleDateString("es-ES", {
    weekday: "long",
    day: "numeric",
    month: "long",
  })
);

const headerMotion = {
  initial: { opacity: 0, y: -20 },
  enter: { opacity: 1, y: 0, transition: { duration: 0.5 } },
};

function openSettings() {
  console.log("Abriendo configuración...");
}
</script>

<style scoped>
/* Efecto de pulso en el punto verde */
@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
  100% {
    opacity: 1;
  }
}
</style>
