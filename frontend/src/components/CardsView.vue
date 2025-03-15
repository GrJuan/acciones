<template>
  <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
    <StockItem v-for="item in filteredItems" :key="item.ticker" :item="item" />
  </div>
</template>

<script lang="ts" setup>
import { useDashboardData } from "@/composables/dashboardData";
import StockItem from "./StockItem.vue";

const { filteredItems } = useDashboardData();
</script>
