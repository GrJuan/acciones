<template>
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
    v-if="visible"
  >
    <div class="bg-gray-900 rounded-lg shadow-lg w-full max-w-3xl p-6 relative">
      <!-- Botón para cerrar -->
      <button
        @click="close"
        class="absolute top-4 right-4 text-gray-400 hover:text-gray-200"
      >
        &times;
      </button>
      <h2 class="text-2xl font-bold text-gray-100 mb-4">Asistente de IA</h2>
      <!-- Pestañas -->
      <div class="flex space-x-4 mb-4">
        <button
          @click="activeTab = 'comparativo'"
          :class="tabButtonClass('comparativo')"
        >
          Comparativo
        </button>
        <button
          @click="activeTab = 'rebalanceo'"
          :class="tabButtonClass('rebalanceo')"
        >
          Rebalanceo
        </button>
        <button
          @click="activeTab = 'sentimiento'"
          :class="tabButtonClass('sentimiento')"
        >
          Sentimiento
        </button>
      </div>
      <!-- Área de contenido con efecto de escritura -->
      <div
        class="bg-gray-800 rounded-lg p-4 h-64 overflow-y-auto text-gray-200 font-mono"
      >
        <p v-if="loading" class="text-center">Generando recomendación...</p>
        <div v-else v-html="displayedHtml"></div>
      </div>
      <!-- Botón para solicitar recomendación -->
      <div class="mt-4 flex justify-end">
        <button
          @click="getRecommendation"
          class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition-colors"
        >
          Obtener Recomendación
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { marked } from "marked";
// Importa la función real si en el futuro deseas usarla
import { getDeepSeekRecommendation } from "@/services/deepseekService";
import { useDashboardData } from "@/composables/dashboardData";

const { availableStocks, portfolioStocks } = useDashboardData();

const visible = ref(false);
const activeTab = ref("comparativo");
const loading = ref(false);
const displayedText = ref("");
const displayedHtml = ref("");
let typingInterval: number | null = null;

// Flag para pruebas: si es true, no se llama al API real
const TEST_MODE = true;

function buildPrompt(): string {
  if (activeTab.value === "comparativo") {
    if (availableStocks.value.length === 0) {
      return "No hay acciones disponibles en este momento.";
    }
    const list = availableStocks.value
      .map((stock) => {
        const price = parseFloat(stock.price).toFixed(2);
        return `${stock.ticker}: precio $${price}, rating ${stock.rating_to}, objetivo ${stock.target_from} - ${stock.target_to}`;
      })
      .join("; ");
    return `### Acciones Disponibles\nAnaliza las siguientes acciones: ${list}.\n\n**¿Cuáles recomendarías comprar hoy y por qué?**`;
  } else if (activeTab.value === "rebalanceo") {
    if (portfolioStocks.value.length === 0) {
      return "No tienes acciones en tu portafolio.";
    }
    const list = portfolioStocks.value
      .map((stock) => {
        const price = parseFloat(stock.price).toFixed(2);
        return `${stock.ticker}: precio $${price}, rating ${stock.rating_to}`;
      })
      .join("; ");
    return `### Portafolio Actual\nMi portafolio tiene: ${list}.\n\n**¿Qué cambios recomendarías para rebalancear mi cartera y mejorar mi rendimiento a corto plazo?**`;
  } else if (activeTab.value === "sentimiento") {
    if (availableStocks.value.length === 0) {
      return "No hay datos disponibles para evaluar el sentimiento del mercado.";
    }
    const sentiments = availableStocks.value
      .map((stock) => {
        return `${stock.ticker}: ${stock.changesPercentage}%`;
      })
      .join("; ");
    return `### Sentimiento del Mercado\nDatos: ${sentiments}.\n\n**¿Cuál es el sentimiento general del mercado hoy y qué acciones parecen prometedoras?**`;
  }
  return "";
}

function typeText(newText: string) {
  if (typingInterval) {
    clearInterval(typingInterval);
  }
  displayedText.value = "";
  displayedHtml.value = "";
  let index = 0;
  typingInterval = window.setInterval(() => {
    if (index < newText.length) {
      displayedText.value += newText[index];
      displayedHtml.value = marked.parse(displayedText.value);
      index++;
    } else {
      clearInterval(typingInterval);
    }
  }, 50);
}

async function getRecommendation() {
  loading.value = true;
  displayedText.value = "";
  displayedHtml.value = "";
  try {
    const prompt = buildPrompt();
    // Si el prompt indica que no hay datos, lo mostramos directamente
    if (prompt.startsWith("No hay") || prompt.startsWith("No tienes")) {
      typeText(prompt);
      return;
    }
    if (TEST_MODE) {
      // Respuesta fija para pruebas (Markdown)
      const fixedResponse = `
### Análisis Comparativo
**Recomendación:** Se sugiere comprar **AAPL** por su potencial de subida del 12% basado en análisis técnico y fundamentales sólidos.
      
1. **AAPL**: Precio $150, Rating Buy, Objetivo $170.
2. **GOOGL**: Precio $2800, Rating Sell, Objetivo $2900.

### Conclusión
El mercado muestra tendencias alcistas en tecnología. Se recomienda diversificar en sectores complementarios.`;
      typeText(fixedResponse);
    } else {
      const apiKey = "TU_DEEPSEEK_API_KEY"; // Reemplaza con tu API key real
      const data = await getDeepSeekRecommendation(prompt, apiKey);
      const fullText = data.choices[0].message.content;
      typeText(fullText);
    }
  } catch (error) {
    typeText("Error al obtener recomendación.");
  } finally {
    loading.value = false;
  }
}

function tabButtonClass(tab: string) {
  return [
    "px-4 py-2 rounded-md",
    activeTab.value === tab
      ? "bg-gray-700 text-white"
      : "bg-gray-600 text-gray-300 hover:bg-gray-500",
  ];
}

function open() {
  visible.value = true;
}

function close() {
  visible.value = false;
}

defineExpose({ open });
</script>

<style scoped>
/* Estilos adicionales o animaciones pueden añadirse aquí */
</style>
