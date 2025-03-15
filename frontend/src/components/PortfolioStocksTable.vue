<template>
  <div class="bg-gray-800 rounded-lg shadow p-4 overflow-auto">
    <h2 class="text-xl font-bold text-gray-100 mb-4 flex items-center">
      <Icon name="ph:table-light" class="mr-2 text-lg text-blue-400" />
      Mis Acciones - Tabla
    </h2>
    <table
      class="min-w-full border border-gray-700 text-gray-200 rounded-lg overflow-hidden"
    >
      <thead class="bg-gray-700 text-gray-300 uppercase text-sm tracking-wider">
        <tr>
          <th class="px-4 py-3 text-left">Ticker</th>
          <th class="px-4 py-3 text-left">Empresa</th>
          <th class="px-4 py-3 text-left">Rating</th>
          <th class="px-4 py-3 text-left">Objetivo</th>
          <th class="px-4 py-3 text-left">Fecha de Compra</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="item in sortedPortfolioStocks"
          :key="item.ticker"
          class="border-b last:border-0 transition-all hover:bg-gray-700"
        >
          <td class="px-4 py-3 font-semibold flex items-center space-x-2">
            <Icon name="mdi:chart-line" class="text-blue-400" />
            <span>{{ item.ticker }}</span>
          </td>
          <td class="px-4 py-3">{{ item.company }}</td>
          <td class="px-4 py-3 flex items-center space-x-2">
            <Icon
              :name="getRatingIcon(item.rating_to)"
              :class="getRatingColor(item.rating_to)"
            />
            <span>{{ item.rating_to }}</span>
          </td>
          <td class="px-4 py-3 font-mono">
            <span class="text-green-400">{{ item.target_from }}</span> -
            <span class="text-red-400">{{ item.target_to }}</span>
          </td>
          <td class="px-4 py-3 text-sm text-gray-400">
            {{ formatDate(item.time) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts" setup>
import { useDashboardData } from "@/composables/dashboardData";

// Importamos la librería de iconos
import { Icon } from "@iconify/vue";

const { sortedPortfolioStocks, formatDate } = useDashboardData();

// Función para asignar un icono según el rating
function getRatingIcon(rating: string) {
  switch (rating.toLowerCase()) {
    case "buy":
      return "mdi:thumb-up";
    case "sell":
      return "mdi:thumb-down";
    case "neutral":
      return "mdi:pause-circle-outline";
    case "outperform":
      return "mdi:star-circle";
    default:
      return "mdi:help-circle";
  }
}

// Función para asignar colores según el rating
function getRatingColor(rating: string) {
  switch (rating.toLowerCase()) {
    case "buy":
      return "text-green-400";
    case "sell":
      return "text-red-400";
    case "neutral":
      return "text-yellow-400";
    case "outperform":
      return "text-blue-400";
    default:
      return "text-gray-400";
  }
}
</script>

<style scoped>
table {
  border-collapse: collapse;
  width: 100%;
}

th {
  font-weight: 600;
  padding: 12px;
  text-align: left;
}

td {
  padding: 12px;
}
</style>
