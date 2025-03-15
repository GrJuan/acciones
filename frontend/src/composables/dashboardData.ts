// src/composables/dashboardData.ts
import { ref, computed, nextTick, watch } from 'vue';
import { chatConversations, addMessage, conversationKey } from '@/composables/chatData';

let dashboardDataInstance: ReturnType<typeof createDashboardData> | null = null;

function createDashboardData() {
    // Datos para "Acciones disponibles" y portafolio
    const availableStocks = ref<any[]>([]);
    const portfolioStocks = ref<any[]>([]);
    // Usaremos nextPage para almacenar el número de página siguiente (como string)
    const nextPage = ref<string | null>(null);
    const loading = ref(false); // Para fetchData
    const loadingRecommendations = ref(false); // Para fetchRecommendations
    const searchQuery = ref('');
    const selectedRating = ref('');
    const connectedUsers = ref<any[]>([]);
    const socket = ref<WebSocket | null>(null);
    const balance = ref(10000);

    function updateBalance(amount: number) {
        balance.value += amount;
    }

    watch(
        () => availableStocks.value.length,
        (newLength, oldLength) => {
            if (newLength === 0 && oldLength > 0) {
                nextPage.value = null; // Reinicia la paginación
                fetchData(); // Vuelve a cargar las acciones
            }
        }
    );
    // Conexión a WebSocket (sin cambios)
    function connectWebSocket() {
        const currentUserId = localStorage.getItem("current_user_id") || "";
        const authToken = localStorage.getItem("auth_token") || "";
        socket.value = new WebSocket(`ws://localhost:8081/ws?user_id=${currentUserId}&token=${authToken}`);
        socket.value.addEventListener("message", (event) => {
            const data = JSON.parse(event.data);
            const currentUserId = localStorage.getItem("current_user_id") || "Anónimo";
            if (data.event === "user-connected") {
                if (data.currentUsers) {
                    connectedUsers.value = data.currentUsers.map((user: any) => ({
                        id: user.id,
                        avatar: user.avatar,
                        isOnline: true,
                    }));
                } else {
                    connectedUsers.value.push({
                        id: data.id,
                        avatar: data.avatar,
                        isOnline: true,
                    });
                }
            } else if (data.event === "user-disconnected") {
                connectedUsers.value = connectedUsers.value.filter(
                    (user) => user.id !== data.id
                );
            } else if (data.event === "private-message") {
                const currentUserId = localStorage.getItem('current_user_id') || '';
                if (data.recipient === currentUserId) {
                    const key = conversationKey(data.sender, data.recipient);
                    addMessage(key, {
                        sender: data.sender,
                        text: data.text,
                        timestamp: Date.now(),
                    }, true);
                }
            }
        });
    }

    // Función para obtener "Acciones disponibles" usando el nuevo endpoint del backend

    async function fetchData() {
        const token = localStorage.getItem('auth_token');
        if (!token || loading.value) return;
        loading.value = true;

        try {
            const page = nextPage.value ? parseInt(nextPage.value) : 1;
            const limit = 10;
            const url = `http://localhost:8081/api/stocks?page=${page}&limit=${limit}`;

            console.log("📡 FetchData URL:", url);

            const response = await fetch(url, {
                method: "GET",
                headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
            });

            if (response.ok) {
                const data = await response.json();
                console.log("🔹 Nuevas acciones recibidas:", data.items);

                if (data.items.length > 0) {
                    availableStocks.value = [...availableStocks.value, ...data.items];
                    nextPage.value = data.next_page || null;
                    console.log("✅ availableStocks actualizado:", availableStocks.value);
                } else {
                    console.warn("⚠ No hay más acciones para cargar");
                }
            } else {
                console.error("❌ Error en fetchData:", response.statusText);
            }
        } catch (error) {
            console.error("❌ Error en fetchData:", error);
        } finally {
            loading.value = false;
        }
    }

    // Variables para la recomendación del día y las sugeridas por ML
    const dailyRecommendation = ref<any>(null);
    const mlRecommendations = ref<any[]>([]);

    // fetchRecommendations llama al endpoint que devuelve las dos llaves:
    // "daily_recommendation" y "ml_recommendations"
    async function fetchRecommendations() {
        console.log("Obteniendo recomendaciones...");
        const token = localStorage.getItem("auth_token");

        if (!token) {
            console.warn("No hay token; se omite fetchRecommendations");
            return null;
        }

        console.log("Token:", token);
        console.log("Valor de loadingRecommendations antes:", loadingRecommendations.value);

        if (loadingRecommendations.value) return null;

        loadingRecommendations.value = true;
        let attempts = 0;
        const maxAttempts = 3;

        while (attempts < maxAttempts) {
            try {
                console.log(`Intento ${attempts + 1} de ${maxAttempts}...`);
                const url = "http://localhost:8081/api/daily-recommendation";
                const response = await fetch(url, {
                    method: "GET",
                    headers: {
                        Authorization: `Bearer ${token}`,
                        "Content-Type": "application/json",
                    },
                });

                if (response.ok) {
                    const data = await response.json();
                    console.log("Recomendación diaria:", data.daily_recommendation);
                    console.log("Recomendaciones ML:", data.ml_recommendations);

                    // Verifica si los datos son válidos
                    if (!data.daily_recommendation || !data.ml_recommendations) {
                        throw new Error("Datos incompletos");
                    }

                    dailyRecommendation.value = data.daily_recommendation;
                    mlRecommendations.value = data.ml_recommendations;
                    loadingRecommendations.value = false;
                    return data;
                } else {
                    console.error("Error al obtener recomendaciones:", response.statusText);
                }
            } catch (error) {
                console.error("Error en fetchRecommendations:", error);
            }

            attempts++;
            if (attempts < maxAttempts) {
                console.warn(`Reintentando en 3 segundos...`);
                await new Promise((resolve) => setTimeout(resolve, 3000));
            }
        }

        console.error("No se pudieron obtener las recomendaciones después de varios intentos.");
        loadingRecommendations.value = false;
        return null;
    }

    function initLazyLoading(mainContent: HTMLElement, loadMore: HTMLElement) {
        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting && nextPage.value && !loading.value) {
                        console.log("🟡 Scroll detectado, cargando más datos...");
                        fetchData();
                    }
                });
            },
            {
                root: mainContent,
                rootMargin: "200px", // Asegura que cargue antes de llegar al final
                threshold: 0.1,
            }
        );

        observer.observe(loadMore);
    }

    /*function initLazyLoading(mainContent: HTMLElement, loadMore: HTMLElement) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting && nextPage.value && !loading.value) {
                    fetchData();
                }
            });
        }, {
            root: mainContent,
            rootMargin: '0px',
            threshold: 0.1,
        });
        observer.observe(loadMore);
    } */

    function formatDate(dateString: string) {
        const date = new Date(dateString);
        return date.toLocaleDateString('es-ES', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    }

    const filteredAvailableStocks = computed(() => {
        return availableStocks.value.filter(item => {
            const matchesSearch =
                item.ticker.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                item.company.toLowerCase().includes(searchQuery.value.toLowerCase());
            const matchesRating = selectedRating.value ? item.rating_to === selectedRating.value : true;
            return matchesSearch && matchesRating;
        });
    });

    const sortedAvailableStocks = computed(() => {
        return [...filteredAvailableStocks.value].sort(
            (a, b) => new Date(b.time).getTime() - new Date(a.time).getTime()
        );
    });

    const kanbanGroups = computed(() => {
        return filteredAvailableStocks.value.reduce((groups, item) => {
            const key = item.rating_to || 'Sin Rating';
            if (!groups[key]) groups[key] = [];
            groups[key].push(item);
            return groups;
        }, {} as Record<string, any[]>);
    });

    function getKanbanHeaderClass(rating: string) {
        const r = rating.toLowerCase();
        if (r === 'buy') return 'bg-green-600 text-white';
        else if (r === 'sell') return 'bg-red-600 text-white';
        else if (r === 'neutral') return 'bg-yellow-500 text-black';
        else if (r === 'outperform') return 'bg-blue-600 text-white';
        else return 'bg-gray-500 text-white';
    }

    function getRecommendationTemperature(item: { target_from: string; target_to: string }) {
        const from = parseFloat(item.target_from.replace("$", ""));
        const to = parseFloat(item.target_to.replace("$", ""));
        if (!from || !to) return { label: "", icon: "" };
        const diffPercentage = ((to - from) / from) * 100;
        if (diffPercentage >= 20) {
            return { label: "Hot", icon: "🔥" };
        } else if (diffPercentage >= 10) {
            return { label: "Warm", icon: "🌤️" };
        } else {
            return { label: "Cold", icon: "❄️" };
        }
    }

    function buyStock(stock: any) {
        const index = availableStocks.value.findIndex(item => item.ticker === stock.ticker);
        if (index > -1) {
            availableStocks.value.splice(index, 1);
        }
        // Marca la acción como comprada
        stock.purchased = true;
        portfolioStocks.value.push({ ...stock, time: new Date().toISOString() });
        const price = parseFloat(stock.target_to.replace("$", ""));
        updateBalance(-price);
    }


    function sellStock(stock: any) {
        const index = portfolioStocks.value.findIndex(item => item.ticker === stock.ticker);
        if (index > -1) {
            portfolioStocks.value.splice(index, 1);
            // Restablecer la propiedad "purchased" en el objeto vendido
            stock.purchased = false;
            // Si la acción vendida es la recomendación del día, actualízala
            if (dailyRecommendation.value && dailyRecommendation.value.ticker === stock.ticker) {
                dailyRecommendation.value.purchased = false;
            }
            // Si la acción vendida está en mlRecommendations, actualízala
            mlRecommendations.value = mlRecommendations.value.map(item => {
                if (item.ticker === stock.ticker) {
                    item.purchased = false;
                }
                return item;
            });
            availableStocks.value.push(stock);
            const price = parseFloat(stock.target_to.replace("$", ""));
            updateBalance(price);
        }
    }

    /* function sellStock(stock: any) {
        const index = portfolioStocks.value.findIndex(item => item.ticker === stock.ticker);
        if (index > -1) {
            portfolioStocks.value.splice(index, 1);

            // Agrega la acción de nuevo a "Acciones disponibles"
            stock.purchased = false;
            availableStocks.value.unshift(stock); // 🔹 Agrega al inicio para que aparezca enseguida
            nextPage.value = null; // Reinicia la paginación
            fetchData(); // Vuelve a cargar las acciones
        }
    } */



    const sortedPortfolioStocks = computed(() => {
        return [...portfolioStocks.value].sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime());
    });

    return {
        availableStocks,
        portfolioStocks,
        balance,
        updateBalance,
        connectedUsers,
        socket,
        fetchData,
        connectWebSocket,
        initLazyLoading,
        formatDate,
        filteredAvailableStocks,
        sortedAvailableStocks,
        kanbanGroups,
        getKanbanHeaderClass,
        getRecommendationTemperature,
        buyStock,
        sellStock,
        sortedPortfolioStocks,
        searchQuery,
        selectedRating,
        loading,
        fetchRecommendations,
        dailyRecommendation,
        mlRecommendations,
    };
}


export function useDashboardData() {
    if (!dashboardDataInstance) {
        dashboardDataInstance = createDashboardData();
    }
    return dashboardDataInstance;
}
