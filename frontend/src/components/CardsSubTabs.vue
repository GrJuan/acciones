<template>
  <div class="flex space-x-2 mb-6">
    <!-- Acciones Disponibles -->
    <button @click="updateView('available')" :class="buttonClass('available')">
      <ChartBarIcon class="w-5 h-5 mr-1 inline-block" />
      Acciones Disponibles
    </button>

    <!-- Estos botones solo se muestran si hay acciones en el portafolio -->
    <template v-if="portfolioStocks.length > 0">
      <button
        @click="updateView('portfolio')"
        :class="buttonClass('portfolio')"
      >
        <FolderIcon class="w-5 h-5 mr-1 inline-block" />
        Mis Acciones
      </button>
      <button @click="updateView('table')" :class="buttonClass('table')">
        <TableCellsIcon class="w-5 h-5 mr-1 inline-block" />
        Tabla
      </button>
      <button @click="updateView('timeline')" :class="buttonClass('timeline')">
        <ClockIcon class="w-5 h-5 mr-1 inline-block" />
        Línea de Tiempo
      </button>
    </template>

    <!-- Kanban -->
    <button @click="updateView('kanban')" :class="buttonClass('kanban')">
      <ViewColumnsIcon class="w-5 h-5 mr-1 inline-block" />
      Kanban
    </button>
  </div>
</template>

<script lang="ts" setup>
import { defineProps, defineEmits, watchEffect, nextTick } from "vue";
import { useDashboardData } from "@/composables/dashboardData";
// Importamos los íconos de Heroicons
import {
  ChartBarIcon,
  FolderIcon,
  TableCellsIcon,
  ClockIcon,
  ViewColumnsIcon,
} from "@heroicons/vue/24/solid";

const props = defineProps<{ modelValue: string }>();
const emit = defineEmits<{ (e: "update:modelValue", value: string): void }>();

const { portfolioStocks } = useDashboardData();

function updateView(value: string) {
  emit("update:modelValue", value);
}

// 🔹 Observamos el estado del portafolio y volvemos a "Acciones Disponibles" si se vacía
watchEffect(() => {
  if (portfolioStocks.value.length === 0 && props.modelValue !== "available") {
    nextTick(() => {
      emit("update:modelValue", "available");
    });
  }
});

function buttonClass(value: string) {
  return [
    "px-4 py-2 rounded-md flex items-center space-x-2",
    props.modelValue === value
      ? "bg-gray-700 text-white"
      : "bg-gray-600 text-gray-300 hover:bg-gray-500",
  ];
}
</script>
