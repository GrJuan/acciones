<template>
  <div class="space-y-6">
    <!-- Gráficos -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="(chart, index) in chartConfig"
        :key="index"
        class="bg-gray-800 rounded-lg shadow p-4 border border-gray-700"
      >
        <h2 class="text-lg font-semibold text-gray-100 mb-2">
          {{ chart.title }}
        </h2>
        <canvas
          :ref="(el) => setChartRef(el, index)"
          class="w-full"
          style="height: 300px"
        ></canvas>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, nextTick, watchEffect } from "vue";
import { Chart } from "chart.js";
import { useDashboardData } from "@/composables/dashboardData";

// 📌 Datos de acciones compradas
const { sortedPortfolioStocks } = useDashboardData();

// 📌 Variables reactivas
const chartRefs = ref<(HTMLCanvasElement | null)[]>([]);
const chartsInstances = ref<Chart[]>([]); // 🔹 Guardamos instancias para poder destruirlas

// 📌 Configuración de gráficos
const chartConfig = ref([
  { type: "line", title: "Precio Objetivo Promedio" },
  { type: "bar", title: "Distribución de Ratings" },
  { type: "pie", title: "Acciones Compradas por Empresa" },
  { type: "bar", title: "Empresas con Mayor Crecimiento" },
]);

// 📌 Función para asignar referencias de canvas correctamente
function setChartRef(el: HTMLCanvasElement | null, index: number) {
  if (el) chartRefs.value[index] = el;
}

// 📌 Función para formatear datos de acciones compradas
function getChartData() {
  if (!sortedPortfolioStocks.value.length) return null;

  const labels = sortedPortfolioStocks.value.map((stock) => stock.ticker);
  const targetPrices = sortedPortfolioStocks.value.map((stock) =>
    parseFloat(stock.target_to.replace("$", ""))
  );
  const ratings = sortedPortfolioStocks.value.map((stock) => stock.rating_to);

  return {
    labels,
    targetPrices,
    ratings,
  };
}

// 📌 Función para renderizar los gráficos
async function renderCharts() {
  await nextTick(); // ✅ Asegura que los canvas estén en el DOM
  const chartData = getChartData();
  if (!chartData) {
    return;
  }

  // 🔹 Destruir gráficos previos antes de renderizar nuevos
  chartsInstances.value.forEach((chart) => chart.destroy());
  chartsInstances.value = [];

  chartRefs.value.forEach((canvas, index) => {
    if (canvas instanceof HTMLCanvasElement) {
      const ctx = canvas.getContext("2d");
      if (!ctx) {
        return;
      }

      let dataset = {
        label: chartConfig.value[index].title,
        data:
          index === 0
            ? chartData.targetPrices
            : chartData.ratings.map(() => Math.random() * 50), // 🔹 Simulación de valores para ratings
        backgroundColor: ["#f87171", "#60a5fa", "#34d399"],
        borderColor: "rgba(75, 192, 192, 1)",
        borderWidth: 1,
      };

      if (chartConfig.value[index].type === "pie") {
        dataset.backgroundColor = ["#ef4444", "#3b82f6", "#10b981"];
      }

      const chartInstance = new Chart(ctx, {
        type: chartConfig.value[index].type,
        data: {
          labels: chartData.labels,
          datasets: [dataset],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
        },
      });

      chartsInstances.value.push(chartInstance); // Guardamos la instancia para limpiar después
    } else {
      console.warn(
        `⚠️ El ref de ${chartConfig.value[index].title} no es un canvas válido`,
        canvas
      );
    }
  });
}

// 📌 Observamos cambios en los datos para renderizar automáticamente
watchEffect(() => {
  renderCharts();
});

onMounted(() => {
  renderCharts();
});
</script>
