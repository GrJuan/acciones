<template>
  <div class="relative pl-6 border-l-4 border-gray-700">
    <h2 class="text-2xl font-bold text-white mb-6">
      📅 Mis Acciones - Línea de Tiempo
    </h2>

    <div
      v-for="(item, index) in sortedPortfolioStocks"
      :key="item.ticker + index"
      class="mb-8 relative group"
    >
      <!-- Línea con conector -->
      <span
        class="absolute -left-4 top-5 w-4 h-4 rounded-full border-2 border-white flex items-center justify-center"
        :class="getActionColor(item.action)"
      >
        <component
          :is="getActionIcon(item.action)"
          class="w-3 h-3 text-white"
        />
      </span>

      <!-- Contenedor del evento -->
      <div
        class="timeline-card bg-gray-800 p-5 rounded-lg shadow-lg border border-gray-600 transition-all duration-300 transform group-hover:scale-100"
      >
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-semibold text-gray-100">
            {{ item.ticker }} - {{ item.company }}
          </h3>
          <span class="text-sm text-gray-400">{{ formatDate(item.time) }}</span>
        </div>

        <p class="text-sm text-gray-300 mt-2">
          <strong>Rating:</strong>
          <span class="text-yellow-400">{{ item.rating_to }}</span>
        </p>
        <p class="text-sm text-gray-300">
          <strong>Objetivo:</strong>
          <span class="text-green-400">{{ item.target_from }}</span> -
          <span class="text-red-400">{{ item.target_to }}</span>
        </p>
        <p class="text-sm text-gray-300">
          <strong>Acción:</strong>
          <span :class="getActionTextColor(item.action)">{{
            item.action
          }}</span>
        </p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useDashboardData } from "@/composables/dashboardData";
import {
  ArrowUpIcon,
  ArrowDownIcon,
  CheckCircleIcon,
} from "@heroicons/vue/24/solid";

const { sortedPortfolioStocks, formatDate } = useDashboardData();

// Función para cambiar el color del punto según la acción
function getActionColor(action: string) {
  switch (action.toLowerCase()) {
    case "buy":
      return "bg-green-600";
    case "sell":
      return "bg-red-600";
    default:
      return "bg-gray-600";
  }
}

// Función para asignar colores al texto de la acción
function getActionTextColor(action: string) {
  switch (action.toLowerCase()) {
    case "buy":
      return "text-green-400";
    case "sell":
      return "text-red-400";
    default:
      return "text-gray-400";
  }
}

// Función para cambiar el ícono según la acción
function getActionIcon(action: string) {
  switch (action.toLowerCase()) {
    case "buy":
      return ArrowUpIcon;
    case "sell":
      return ArrowDownIcon;
    default:
      return CheckCircleIcon;
  }
}
</script>

<style scoped>
/* Estilo de animación para la línea de tiempo */
.timeline-card {
  position: relative;
  overflow: hidden;
}

.timeline-card::before {
  content: "";
  position: absolute;
  left: -10px;
  top: 0;
  bottom: 0;
  width: 4px;
  background: linear-gradient(to bottom, #16a34a, #dc2626);
  opacity: 0.7;
}
</style>
