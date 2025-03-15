<template>
  <div class="bg-gray-800 rounded-lg shadow p-4 overflow-auto">
    <table class="min-w-full border border-gray-700 text-gray-200">
      <thead class="bg-gray-700">
        <tr>
          <th class="px-4 py-2 text-left">Ticker</th>
          <th class="px-4 py-2 text-left">Empresa</th>
          <th class="px-4 py-2 text-left">Rating Anterior</th>
          <th class="px-4 py-2 text-left">Rating Actual</th>
          <th class="px-4 py-2 text-left">Objetivo</th>
          <th class="px-4 py-2 text-left">Acción</th>
          <th class="px-4 py-2 text-left">Corretaje</th>
          <th class="px-4 py-2 text-left">Fecha</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="item in filteredItems"
          :key="item.ticker"
          class="border-b last:border-0 hover:bg-gray-700"
        >
          <td class="px-4 py-2">{{ item.ticker }}</td>
          <td class="px-4 py-2">{{ item.company }}</td>
          <td class="px-4 py-2">{{ item.rating_from }}</td>
          <td class="px-4 py-2">{{ item.rating_to }}</td>
          <td class="px-4 py-2">
            {{ item.target_from }} - {{ item.target_to }}
          </td>
          <td class="px-4 py-2">{{ item.action }}</td>
          <td class="px-4 py-2">{{ item.brokerage }}</td>
          <td class="px-4 py-2">{{ formatDate(item.time) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts" setup>
import { useDashboardData } from "@/composables/dashboardData";
const { filteredItems, formatDate } = useDashboardData();
</script>
