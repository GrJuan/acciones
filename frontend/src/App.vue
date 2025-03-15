<template>
  <router-view v-if="isLoginPage" />
  <div v-else-if="isAuthenticated">
    <div class="flex h-screen bg-gray-900 home-container">
      <Sidebar />
      <div ref="mainContentRef" class="flex-1 overflow-y-auto p-6">
        <StockTicker />
        <HeaderComponent />
        <PortfolioHeader />
        <TabsComponent
          :activeTab="activeTab"
          @update:activeTab="activeTab = $event"
        />
        <!-- Subbotones de la pestaña "cards" -->
        <CardsSubTabs v-if="activeTab === 'cards'" v-model="viewMode" />
        <!-- Vista dinámica según activeTab y viewMode -->
        <component :is="currentView" />
        <!-- Lazy loading solo se activa cuando se muestran acciones disponibles -->
        <div
          v-if="activeTab === 'cards' && viewMode === 'available'"
          ref="loadMoreRef"
          class="mt-6 text-center text-gray-400"
        >
          <span v-if="loading">Cargando más datos...</span>
        </div>
      </div>
      <RecommendationsColumn />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, nextTick, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import ChatWindow from "./components/ChatWindow.vue";
import Sidebar from "./components/Sidebar.vue";
import StockTicker from "./components/StockTicker.vue";
import HeaderComponent from "./components/HeaderComponent.vue";
import PortfolioHeader from "./components/PortfolioHeader.vue";
import TabsComponent from "./components/TabsComponent.vue";
import CardsSubTabs from "./components/CardsSubTabs.vue";
import AvailableStocks from "./components/AvailableStocks.vue";
import PortfolioStocks from "./components/PortfolioStocks.vue";
import PortfolioStocksTable from "./components/PortfolioStocksTable.vue";
import PortfolioStocksTimeline from "./components/PortfolioStocksTimeline.vue";
import KanbanView from "./components/KanbanView.vue";
import ChartsView from "./components/ChartsView.vue";
import RecommendationsColumn from "./components/RecommendationsColumn.vue";
import { useDashboardData } from "@/composables/dashboardData";

const router = useRouter();
const route = useRoute();
const isAuthenticated = ref(false);

const { connectWebSocket, fetchData, initLazyLoading, loading } =
  useDashboardData();

const activeTab = ref("cards");
const viewMode = ref("available");

// Si el usuario está en la página de login, no mostramos el layout completo
const isLoginPage = computed(() => route.name === "Login");

watch([activeTab, viewMode], async () => {
  await nextTick();
  if (mainContentRef.value && loadMoreRef.value) {
    initLazyLoading(mainContentRef.value, loadMoreRef.value);
  }
});

const currentView = computed(() => {
  if (activeTab.value === "cards") {
    switch (viewMode.value) {
      case "available":
        return AvailableStocks;
      case "portfolio":
        return PortfolioStocks;
      case "table":
        return PortfolioStocksTable;
      case "timeline":
        return PortfolioStocksTimeline;
      case "kanban":
        return KanbanView;
      default:
        return AvailableStocks;
    }
  } else if (activeTab.value === "charts") {
    return ChartsView;
  }
  return AvailableStocks;
});

const mainContentRef = ref<HTMLElement | null>(null);
const loadMoreRef = ref<HTMLElement | null>(null);

onMounted(() => {
  const token = localStorage.getItem("auth_token");
  if (!token) {
    router.push({ name: "Login" });
  } else {
    isAuthenticated.value = true;
    connectWebSocket();
    fetchData();
    nextTick(() => {
      if (mainContentRef.value && loadMoreRef.value) {
        initLazyLoading(mainContentRef.value, loadMoreRef.value);
      }
    });
  }
});
</script>
