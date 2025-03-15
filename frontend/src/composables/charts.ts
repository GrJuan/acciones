import { ref, computed } from "vue";
import { useDashboardData } from "@/composables/dashboardData";
import { Chart, registerables } from "chart.js";
import "chartjs-adapter-date-fns";

Chart.register(...registerables);

export function useCharts() {
    const { sortedPortfolioStocks } = useDashboardData();

    // Rango de fechas para filtrar datos
    const startDate = ref("");
    const endDate = ref("");

    // Filtrar acciones por rango de fechas
    const filteredStocks = computed(() => {
        if (!startDate.value || !endDate.value) return sortedPortfolioStocks.value;

        return sortedPortfolioStocks.value.filter((stock) => {
            const stockDate = new Date(stock.time);
            return (
                stockDate >= new Date(startDate.value) &&
                stockDate <= new Date(endDate.value)
            );
        });
    });

    // 📈 Datos para el gráfico de Precio Objetivo Promedio
    const priceData = computed(() => {
        return {
            labels: filteredStocks.value.map((s) => s.ticker),
            datasets: [
                {
                    label: "Precio Objetivo",
                    data: filteredStocks.value.map((s) =>
                        parseFloat(s.target_to.replace("$", ""))
                    ),
                    borderColor: "#4F46E5",
                    backgroundColor: "rgba(79, 70, 229, 0.2)",
                    fill: true,
                    tension: 0.3,
                },
            ],
        };
    });

    // 📊 Datos para el gráfico de Distribución de Ratings
    const ratingData = computed(() => {
        const ratingCount: Record<string, number> = {};
        filteredStocks.value.forEach((s) => {
            ratingCount[s.rating_to] = (ratingCount[s.rating_to] || 0) + 1;
        });

        return {
            labels: Object.keys(ratingCount),
            datasets: [
                {
                    label: "Recomendaciones",
                    data: Object.values(ratingCount),
                    backgroundColor: ["#16a34a", "#eab308", "#dc2626"],
                },
            ],
        };
    });

    // 🍕 Datos para el gráfico de Corretajes Activos
    const brokerageData = computed(() => {
        const brokerageCount: Record<string, number> = {};
        filteredStocks.value.forEach((s) => {
            brokerageCount[s.brokerage] = (brokerageCount[s.brokerage] || 0) + 1;
        });

        return {
            labels: Object.keys(brokerageCount),
            datasets: [
                {
                    label: "Número de Análisis",
                    data: Object.values(brokerageCount),
                    backgroundColor: ["#ef4444", "#f97316", "#eab308", "#22c55e", "#3b82f6"],
                },
            ],
        };
    });

    // 🔻 Empresas para Vender | 🔺 Empresas para Comprar
    const buySellData = computed(() => {
        const buyCompanies = filteredStocks.value.filter((s) => s.rating_to === "Buy");
        const sellCompanies = filteredStocks.value.filter((s) => s.rating_to === "Sell");

        return {
            buy: {
                labels: buyCompanies.map((s) => s.ticker),
                datasets: [
                    {
                        label: "Comprar",
                        data: buyCompanies.map((s) => parseFloat(s.target_to.replace("$", ""))),
                        backgroundColor: "#22c55e",
                    },
                ],
            },
            sell: {
                labels: sellCompanies.map((s) => s.ticker),
                datasets: [
                    {
                        label: "Vender",
                        data: sellCompanies.map((s) => parseFloat(s.target_to.replace("$", ""))),
                        backgroundColor: "#ef4444",
                    },
                ],
            },
        };
    });

    function resetZoom() {
        startDate.value = "";
        endDate.value = "";
    }

    return {
        startDate,
        endDate,
        filteredStocks,
        priceData,
        ratingData,
        brokerageData,
        buySellData,
        resetZoom,
    };
}
