<template>
  <div
    class="recommendations bg-gray-800 border-l border-gray-700 w-[320px] h-screen overflow-y-auto p-4 shadow-md"
    @mouseenter="pauseRecommendations"
    @mouseleave="resumeRecommendations"
  >
    <h2 class="text-xl font-bold text-gray-100 mb-4 flex items-center">
      <Icon
        name="mdi:lightbulb-on-outline"
        class="mr-2 text-lg text-yellow-400"
      />
      Recomendaciones de Inversión
    </h2>

    <!-- Sección: Recomendación del Día -->
    <div class="mb-6">
      <h3 class="text-lg font-semibold text-gray-200 mb-2 flex items-center">
        <Icon
          name="mdi:star-circle-outline"
          class="mr-2 text-yellow-500 text-xl"
        />
        Recomendación del Día
      </h3>
      <div
        v-if="dailyRecommendation"
        :class="[
          'recommendation-item border border-gray-700 rounded-lg p-4 transition-all',
          dailyRecommendation.purchased
            ? 'bg-green-500/50'
            : 'hover:bg-gray-700 shadow-lg',
        ]"
      >
        <h3 class="text-lg font-semibold text-yellow-400">
          {{ dailyRecommendation.ticker }} - {{ dailyRecommendation.company }}
        </h3>
        <p class="text-sm text-gray-300">
          <strong>Rating:</strong> {{ dailyRecommendation.rating_to }}
        </p>
        <p class="text-sm text-gray-300">
          <strong>Objetivo:</strong>
          <span class="text-green-400">{{
            dailyRecommendation.target_from
          }}</span>
          -
          <span class="text-red-400">{{ dailyRecommendation.target_to }}</span>
        </p>
        <p class="text-sm text-gray-300">
          <strong>Razón:</strong> {{ dailyRecommendation.reason }}
        </p>
        <p class="text-xs text-gray-400">{{ dailyRecommendation.date }}</p>

        <canvas
          ref="dailyScaleChartRef"
          width="100"
          height="40"
          class="scale-chart-canvas mx-auto mt-2"
        ></canvas>

        <!-- Botones de opciones -->
        <div class="mt-2 space-y-2">
          <button
            v-if="!dailyRecommendation.purchased"
            class="btn-action bg-green-600 hover:bg-green-700"
            style="color: #fff"
            @click="comprarAccionDaily(dailyRecommendation)"
          >
            Comprar Acción
          </button>
          <button
            v-if="!dailyRecommendation.purchased"
            class="btn-action bg-blue-600 hover:bg-blue-700"
            style="color: #fff"
            @click="pregustarDaily(dailyRecommendation)"
          >
            Preguntar sobre esta acción
          </button>
        </div>

        <div
          v-if="dailyRecommendation.purchased"
          class="absolute inset-0 flex items-center justify-center comprado-label"
        >
          <Icon name="mdi:check-circle" class="text-green-400 text-4xl" />
        </div>
      </div>
    </div>

    <!-- Sección: Recomendaciones de ML -->
    <div>
      <h3 class="text-lg font-semibold text-gray-200 mb-2 flex items-center">
        <Icon name="mdi:robot-outline" class="mr-2 text-blue-400 text-xl" />
        Recomendaciones de Machine Learning
      </h3>

      <div v-if="mlRecommendations && mlRecommendations.length">
        <div
          v-for="(item, idx) in mlRecommendations"
          :key="item.ticker"
          :class="[
            'recommendation-item border border-gray-700 rounded-lg p-4 transition-all',
            item.purchased ? 'bg-green-500/50' : 'hover:bg-gray-700 shadow-lg',
          ]"
        >
          <h3 class="text-lg font-semibold text-indigo-400">
            {{ item.ticker }} - {{ item.company }}
          </h3>
          <p class="text-sm text-gray-300">
            <strong>Rating:</strong> {{ item.rating_to }}
          </p>
          <p class="text-sm text-gray-300">
            <strong>Objetivo:</strong>
            <span class="text-green-400">{{ item.target_from }}</span> -
            <span class="text-red-400">{{ item.target_to }}</span>
          </p>
          <p class="text-sm text-gray-300">
            <strong>Razón:</strong> {{ item.reason }}
          </p>
          <p class="text-xs text-gray-400">{{ item.date }}</p>

          <canvas
            :ref="(el) => (mlScaleChartRefs[idx] = el)"
            width="100"
            height="40"
            class="scale-chart-canvas mx-auto mt-2"
          ></canvas>

          <div class="mt-2 space-y-2">
            <button
              v-if="!item.purchased"
              class="btn-action bg-green-600 hover:bg-green-700"
              style="color: #fff"
              @click="comprarAccionML(item)"
            >
              Comprar Acción
            </button>
            <button
              v-if="!item.purchased"
              style="color: #fff"
              class="btn-action bg-blue-600 hover:bg-blue-700"
              @click="pregustarML(item)"
            >
              Preguntar sobre esta acción
            </button>
          </div>

          <div
            v-if="item.purchased"
            class="absolute inset-0 flex items-center justify-center comprado-label"
          >
            <Icon name="mdi:check-circle" class="text-green-400 text-4xl" />
          </div>
        </div>
      </div>
      <div v-else>
        <p class="text-sm text-gray-300">Cargando recomendaciones de ML...</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, nextTick, watch } from "vue";
import { useDashboardData } from "@/composables/dashboardData";
import { Chart, registerables } from "chart.js";
import "chartjs-adapter-date-fns";
import "chart.js/auto";
import confetti from "canvas-confetti";
import { useMotionVariants } from "@vueuse/motion";

Chart.register(...registerables);

const {
  dailyRecommendation,
  mlRecommendations,
  getRecommendationTemperature,
  pauseRecommendations,
  resumeRecommendations,
  buyStock,
  fetchRecommendations,
} = useDashboardData();

// Variables para animación
const animationVisible = ref<Record<string, boolean>>({});
const animationMessage = ref<Record<string, string>>({});

// Funciones wrapper para la sección de Recomendación del día
function comprarAccionDaily(item: any) {
  buyStock(item);
  console.log("Se compró la acción (día):", item.ticker);
  dailyRecommendation.value.purchased = true;
  animationMessage.value[item.ticker] = "Bought!";
  animationVisible.value[item.ticker] = true;
  const salePrice = parseFloat(item.target_to.replace("$", ""));
  if (salePrice > 200) {
    confetti({
      particleCount: 100,
      spread: 70,
      origin: { y: 0.6 },
    });
  }
  setTimeout(() => {
    animationVisible.value[item.ticker] = false;
  }, 1500);
}

function pregustarDaily(item: any) {
  console.log("Pregustar (día) sobre la acción:", item.ticker);
}

// Funciones wrapper para la sección de Recomendación de ML
function comprarAccionML(item: any) {
  buyStock(item);
  // Actualizamos solo el item comprado dentro de mlRecommendations
  mlRecommendations.value = mlRecommendations.value.map((rec) => {
    if (rec.ticker === item.ticker) {
      rec.purchased = true;
    }
    return rec;
  });
  animationMessage.value[item.ticker] = "Bought!";
  animationVisible.value[item.ticker] = true;
  const salePrice = parseFloat(item.target_to.replace("$", ""));
  if (salePrice > 200) {
    confetti({
      particleCount: 100,
      spread: 70,
      origin: { y: 0.6 },
    });
  }
  setTimeout(() => {
    animationVisible.value[item.ticker] = false;
  }, 1500);
}

function pregustarML(item: any) {
  console.log("Pregustar (ML) sobre la acción:", item.ticker);
}

// Flags para mostrar opciones
const showDailyOptions = ref(false);
const activeOptionIndexML = ref<number | null>(null);
const activeOptionIndex = ref<number | null>(null);

function toggleDailyOptions() {
  showDailyOptions.value = !showDailyOptions.value;
}

function toggleMLOptions(idx: number) {
  activeOptionIndexML.value = activeOptionIndexML.value === idx ? null : idx;
}

// Referencias para los canvas de gráficos
const dailyScaleChartRef = ref<HTMLCanvasElement | null>(null);
const mlScaleChartRefs = ref<(HTMLCanvasElement | null)[]>([]);
const gaugeChartDaily = ref<Chart | null>(null);
const gaugeChartsML = ref<Chart[]>([]);

function initDailyChart() {
  if (!dailyScaleChartRef.value || !dailyRecommendation.value) return;
  if (gaugeChartDaily.value) {
    gaugeChartDaily.value.destroy();
    gaugeChartDaily.value = null;
  }
  const item = dailyRecommendation.value;
  const fromVal = parseFloat(item.target_from.replace("$", "")) || 0;
  const toVal = parseFloat(item.target_to.replace("$", "")) || 0;
  const valley = Math.min(fromVal, toVal) * 0.5;
  const dataPoints = [0, fromVal, valley, toVal];
  gaugeChartDaily.value = new Chart(dailyScaleChartRef.value.getContext("2d"), {
    type: "line",
    data: {
      labels: ["Inicio", "Objetivo", "Valle", "Final"],
      datasets: [
        {
          data: dataPoints,
          borderColor: "#60a5fa",
          backgroundColor: "transparent",
          fill: false,
          tension: 0.3,
          pointRadius: 0,
          segment: {
            borderColor: (ctx) => {
              const chartCtx = ctx.chart.ctx;
              const gradient = chartCtx.createLinearGradient(
                ctx.p0.x,
                ctx.p0.y,
                ctx.p1.x,
                ctx.p1.y
              );
              if (ctx.p1.parsed.y > ctx.p0.parsed.y) {
                gradient.addColorStop(0, "#ffcccc");
                gradient.addColorStop(1, "#ff0000");
              } else {
                gradient.addColorStop(0, "#ffe5b4");
                gradient.addColorStop(1, "#ff8c00");
              }
              return gradient;
            },
          },
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: { enabled: false },
      },
      scales: {
        x: { display: false },
        y: { display: false },
      },
      animation: { duration: 500 },
    },
  });
}

function initMLCharts() {
  gaugeChartsML.value.forEach((chart) => chart.destroy());
  gaugeChartsML.value = [];
  mlScaleChartRefs.value.forEach((canvas, idx) => {
    if (!canvas) return;
    const item = mlRecommendations.value[idx];
    if (!item) return;
    const fromVal = parseFloat(item.target_from.replace("$", "")) || 0;
    const toVal = parseFloat(item.target_to.replace("$", "")) || 0;
    const valley = Math.min(fromVal, toVal) * 0.5;
    const dataPoints = [0, fromVal, valley, toVal];
    const chartInstance = new Chart(canvas.getContext("2d"), {
      type: "line",
      data: {
        labels: ["Inicio", "Objetivo", "Valle", "Final"],
        datasets: [
          {
            data: dataPoints,
            borderColor: "#60a5fa",
            backgroundColor: "transparent",
            fill: false,
            tension: 0.3,
            pointRadius: 0,
            segment: {
              borderColor: (ctx) => {
                const chartCtx = ctx.chart.ctx;
                const gradient = chartCtx.createLinearGradient(
                  ctx.p0.x,
                  ctx.p0.y,
                  ctx.p1.x,
                  ctx.p1.y
                );
                if (ctx.p1.parsed.y > ctx.p0.parsed.y) {
                  gradient.addColorStop(0, "#ffcccc");
                  gradient.addColorStop(1, "#ff0000");
                } else {
                  gradient.addColorStop(0, "#ffe5b4");
                  gradient.addColorStop(1, "#ff8c00");
                }
                return gradient;
              },
            },
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false },
          tooltip: { enabled: false },
        },
        scales: {
          x: { display: false },
          y: { display: false },
        },
        animation: { duration: 500 },
      },
    });
    gaugeChartsML.value.push(chartInstance);
  });
}

watch(
  () => dailyRecommendation.value,
  async (newVal) => {
    if (newVal) {
      await nextTick();
      setTimeout(() => {
        initDailyChart();
      }, 300);
    }
  },
  { immediate: true }
);

watch(
  () => mlRecommendations.value,
  async (newVal) => {
    if (newVal && newVal.length > 0) {
      await nextTick();
      setTimeout(() => {
        initMLCharts();
      }, 300);
    }
  },
  { immediate: true }
);

onMounted(async () => {
  await fetchRecommendations();
});
</script>

<style scoped>
.scale-chart-canvas {
  display: block;
  margin: 0 auto;
  height: 40px !important;
  max-height: 40px !important;
}

.btn-action {
  width: 100%;
  text-align: center;
  padding: 8px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.2s;
}

.comprado-label {
  position: absolute;
  top: 50%;
  left: 50%;
  color: white;
  font-size: 1.5rem;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 1px;
  border-radius: 5px;
  padding: 4px 8px;
  z-index: 50; /* Asegura que está sobre el contenido */
  pointer-events: none; /* Evita que bloquee clics */
}
</style>
