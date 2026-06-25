document.addEventListener('DOMContentLoaded', () => {
            function safeCreateIcons() {
                if (window.lucide && typeof window.lucide.createIcons === 'function') {
                    window.lucide.createIcons();
                    document.body.classList.remove('icons-unavailable');
                    return;
                }

                document.body.classList.add('icons-unavailable');
            }

            // Инициализация графики/иконок
            safeCreateIcons();

            // --- Логика переключения темы (Темная/Светлая) ---
            const themeSwitch = document.getElementById('theme-switch-checkbox');
            const themeLabel = document.getElementById('theme-toggle-label');
            
            // Восстановление темы из localStorage
            const savedTheme = localStorage.getItem('taskflow_theme') || 'dark';
            if (savedTheme === 'light') {
                document.body.classList.add('light-theme');
                if (themeSwitch) themeSwitch.checked = true;
                if (themeLabel) themeLabel.textContent = 'Светлый режим';
            } else {
                if (themeSwitch) themeSwitch.checked = false;
                if (themeLabel) themeLabel.textContent = 'Темный режим';
            }

            if (themeSwitch) {
                themeSwitch.addEventListener('change', (e) => {
                    const isLight = e.target.checked;
                    if (isLight) {
                        document.body.classList.add('light-theme');
                        localStorage.setItem('taskflow_theme', 'light');
                        if (themeLabel) themeLabel.textContent = 'Светлый режим';
                    } else {
                        document.body.classList.remove('light-theme');
                        localStorage.setItem('taskflow_theme', 'dark');
                        if (themeLabel) themeLabel.textContent = 'Темный режим';
                    }
                });
            }

            // --- Логика кастомного выбора цвета акцента ---
            const svCanvas = document.getElementById('sv-canvas');
            const svPointer = document.getElementById('sv-pointer');
            const hueSlider = document.getElementById('hue-slider');
            const colorPreview = document.getElementById('color-preview');
            const colorHex = document.getElementById('color-hex');
            const resetColorBtn = document.getElementById('reset-color-btn');

            // Значения по умолчанию (соответствуют индиго #6366f1)
            const DEFAULT_COLOR = { h: 238, s: 0.59, v: 0.95 };
            let currentColor = { ...DEFAULT_COLOR };

            // HSV в HSL
            function hsvToHsl(h, s, v) {
                const l = v * (1 - s / 2);
                const sHsl = (l === 0 || l === 1) ? 0 : (v - l) / Math.min(l, 1 - l);
                return {
                    h: Math.round(h),
                    s: Math.round(sHsl * 100),
                    l: Math.round(l * 100)
                };
            }

            // HSL в HEX
            function hslToHex(h, s, l) {
                l /= 100;
                const a = (s * Math.min(l, 1 - l)) / 100;
                const f = n => {
                    const k = (n + h / 30) % 12;
                    const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
                    return Math.round(255 * color).toString(16).padStart(2, '0');
                };
                return `#${f(0)}${f(8)}${f(4)}`;
            }

            // HEX в RGB строку (для --primary-glow)
            function hexToRgb(hex) {
                const r = parseInt(hex.slice(1, 3), 16);
                const g = parseInt(hex.slice(3, 5), 16);
                const b = parseInt(hex.slice(5, 7), 16);
                return `${r}, ${g}, ${b}`;
            }

            // Применить цвет глобально
            function applyAccentColor(h, s, v) {
                currentColor = { h, s, v };
                const hsl = hsvToHsl(h, s, v);
                const hex = hslToHex(hsl.h, hsl.s, hsl.l);

                // Расчет контрастного цвета текста для кнопок на основе яркости
                function getContrastTextColor(hex) {
                    const r = parseInt(hex.slice(1, 3), 16);
                    const g = parseInt(hex.slice(3, 5), 16);
                    const b = parseInt(hex.slice(5, 7), 16);
                    
                    // Вычисление относительной яркости (Relative Luminance)
                    const y = 0.2126 * r + 0.7152 * g + 0.0722 * b;
                    
                    // Коэффициент "светлоты" от 0 (темный) до 1 (белый)
                    const threshold = 140; // порог, при котором текст начинает темнеть
                    let w = 0;
                    if (y > threshold) {
                        w = (y - threshold) / (255 - threshold);
                    }
                    w = Math.max(0, Math.min(1, w)); // зажим от 0 до 1
                    
                    // Интерполяция между белым (255, 255, 255) и темно-серым (17, 24, 39)
                    const textR = Math.round(255 - w * (255 - 17));
                    const textG = Math.round(255 - w * (255 - 24));
                    const textB = Math.round(255 - w * (255 - 39));
                    
                    return `rgb(${textR}, ${textG}, ${textB})`;
                }

                // Установка CSS переменных на корневом уровне и уровне body для корректной работы в светлой теме
                document.documentElement.style.setProperty('--primary', hex);
                document.documentElement.style.setProperty('--primary-glow', `rgba(${hexToRgb(hex)}, 0.15)`);
                document.documentElement.style.setProperty('--text-on-primary', getContrastTextColor(hex));

                document.body.style.setProperty('--primary', hex);
                document.body.style.setProperty('--primary-glow', `rgba(${hexToRgb(hex)}, 0.15)`);
                document.body.style.setProperty('--text-on-primary', getContrastTextColor(hex));

                // Обновление элементов предпросмотра
                if (colorPreview) colorPreview.style.backgroundColor = hex;
                if (colorHex) colorHex.textContent = hex.toUpperCase();

                // Обновление положения ползунка Hue
                if (hueSlider) {
                    hueSlider.value = h;
                }
                
                // Обновление фона полотна SV
                if (svCanvas) {
                    svCanvas.style.backgroundColor = `hsl(${h}, 100%, 50%)`;
                }

                // Обновление положения указателя на полотне
                if (svPointer) {
                    svPointer.style.left = `${s * 100}%`;
                    svPointer.style.top = `${(1 - v) * 100}%`;
                }

                // Сохранение в localStorage
                localStorage.setItem('taskflow_accent', JSON.stringify(currentColor));
            }

            // Обработка клика/перетаскивания по полотну SV
            let isDraggingSV = false;

            function handleSVMouse(e) {
                if (!svCanvas) return;
                const rect = svCanvas.getBoundingClientRect();
                
                // Рассчитываем координаты x и y в пределах от 0 до 1
                let x = (e.clientX - rect.left) / rect.width;
                let y = (e.clientY - rect.top) / rect.height;

                // Ограничиваем рамками
                x = Math.max(0, Math.min(1, x));
                y = Math.max(0, Math.min(1, y));

                const s = x;
                const v = 1 - y;
                const h = parseInt(hueSlider.value);

                applyAccentColor(h, s, v);
            }

            if (svCanvas) {
                svCanvas.addEventListener('mousedown', (e) => {
                    isDraggingSV = true;
                    handleSVMouse(e);
                });

                // Сенсорная поддержка для мобильных
                svCanvas.addEventListener('touchstart', (e) => {
                    isDraggingSV = true;
                    if (e.touches[0]) handleSVMouse(e.touches[0]);
                }, { passive: true });
            }

            window.addEventListener('mousemove', (e) => {
                if (isDraggingSV) handleSVMouse(e);
            });

            window.addEventListener('touchmove', (e) => {
                if (isDraggingSV && e.touches[0]) handleSVMouse(e.touches[0]);
            }, { passive: true });

            window.addEventListener('mouseup', () => {
                isDraggingSV = false;
            });
            window.addEventListener('touchend', () => {
                isDraggingSV = false;
            });

            // Обработка ползунка Hue
            if (hueSlider) {
                hueSlider.addEventListener('input', (e) => {
                    const h = parseInt(e.target.value);
                    applyAccentColor(h, currentColor.s, currentColor.v);
                });
            }

            // Кнопка сброса цвета к дефолтному
            if (resetColorBtn) {
                resetColorBtn.addEventListener('click', () => {
                    applyAccentColor(DEFAULT_COLOR.h, DEFAULT_COLOR.s, DEFAULT_COLOR.v);
                });
            }

            // Инициализация при загрузке страницы
            const savedAccent = localStorage.getItem('taskflow_accent');
            if (savedAccent) {
                try {
                    const parsed = JSON.parse(savedAccent);
                    applyAccentColor(parsed.h, parsed.s, parsed.v);
                } catch (err) {
                    console.error('Ошибка восстановления цвета:', err);
                    applyAccentColor(DEFAULT_COLOR.h, DEFAULT_COLOR.s, DEFAULT_COLOR.v);
                }
            } else {
                applyAccentColor(DEFAULT_COLOR.h, DEFAULT_COLOR.s, DEFAULT_COLOR.v);
            }

            // --- Константы API ---
            const API_BASE_URL = '/api/v1';

            // Параметры пагинации
            let tasksPage = 1;
            let tasksLimit = 20;
            let usersPage = 1;
            let usersLimit = 20;
            let statsChartPage = 1;
            let statsChartLimit = 10;
            
            // Фильтры для статистики
            let statsFilterUserId = null;
            let statsFilterFromDate = null;
            let statsFilterToDate = null;

            // Фильтры для задач
            let tasksFilterUserId = null;

            // Состояние приложения (загружается только с бэкенда)
            let appState = {
                tasks: [],              // Задачи на текущей странице
                tasksHasNextPage: false,
                allTasks: [],           // Все задачи для графиков статистики
                users: [],              // Все пользователи для выпадающих списков и графиков
                usersHasNextPage: false,
                paginatedUsers: [],     // Пользователи на текущей странице
                statistics: {
                    tasks_created: 0,
                    task_completed: 0,
                    task_completed_rate: 0,
                    task_average_completion_time: null
                }
            };

            const offlineOverlay = document.getElementById('offline-overlay');
            const retryBtn = document.getElementById('retry-connection-btn');
            let activeLoadController = null;
            let loadRequestId = 0;
            let activeStatsController = null;
            let statsLoadRequestId = 0;

            const USERS_DIRECTORY_TTL_MS = 60 * 1000;
            const STATS_TASKS_TTL_MS = 60 * 1000;
            let usersDirectoryLoadedAt = 0;
            let statsTasksLoadedAt = 0;
            let statsTasksFilterKey = null;

            function buildApiUrl(endpoint, { limit, offset, params = {} } = {}) {
                const searchParams = new URLSearchParams();

                if (limit !== undefined) {
                    searchParams.set('limit', String(limit));
                }

                if (offset !== undefined) {
                    searchParams.set('offset', String(offset));
                }

                Object.entries(params).forEach(([key, value]) => {
                    if (value !== null && value !== undefined && value !== '') {
                        searchParams.set(key, String(value));
                    }
                });

                return searchParams.size > 0
                    ? `${API_BASE_URL}${endpoint}?${searchParams.toString()}`
                    : `${API_BASE_URL}${endpoint}`;
            }

            async function fetchJson(url, { signal } = {}) {
                const response = await fetch(url, { signal });
                if (!response.ok) {
                    throw new Error(`Backend returned error status for ${url}`);
                }

                return response.json();
            }

            async function fetchAllPages(endpoint, { limit = 100, params = {}, signal } = {}) {
                const items = [];
                let offset = 0;

                while (true) {
                    const pageItems = await fetchJson(buildApiUrl(endpoint, {
                        limit,
                        offset,
                        params
                    }), { signal });

                    if (!Array.isArray(pageItems)) {
                        break;
                    }

                    items.push(...pageItems);

                    if (pageItems.length < limit) {
                        break;
                    }

                    offset += limit;
                }

                return items;
            }

            async function fetchPaginatedSlice(endpoint, { page, limit, params = {}, signal } = {}) {
                const items = await fetchJson(buildApiUrl(endpoint, {
                    limit: limit + 1,
                    offset: (page - 1) * limit,
                    params
                }), { signal });

                const normalizedItems = Array.isArray(items) ? items : [];

                return {
                    items: normalizedItems.slice(0, limit),
                    hasNextPage: normalizedItems.length > limit
                };
            }

            function isStatsSectionActive() {
                const statsSection = document.getElementById('section-stats');
                return Boolean(statsSection && statsSection.classList.contains('active'));
            }

            function getStatsTasksFilterKey() {
                return tasksFilterUserId ? `user:${tasksFilterUserId}` : 'all';
            }

            function invalidateUsersDirectory() {
                usersDirectoryLoadedAt = 0;
            }

            function invalidateStatsTasksData() {
                appState.allTasks = [];
                statsTasksLoadedAt = 0;
                statsTasksFilterKey = null;
            }

            async function loadUsersDirectory({ force = false, signal } = {}) {
                const hasFreshCache = !force
                    && usersDirectoryLoadedAt > 0
                    && (Date.now() - usersDirectoryLoadedAt) < USERS_DIRECTORY_TTL_MS;

                if (hasFreshCache) {
                    return appState.users;
                }

                const users = await fetchAllPages('/users', { signal });
                appState.users = Array.isArray(users) ? users : [];
                usersDirectoryLoadedAt = Date.now();

                return appState.users;
            }

            async function loadStatsTasksData({ force = false, signal } = {}) {
                const filterKey = getStatsTasksFilterKey();
                const hasFreshCache = !force
                    && statsTasksLoadedAt > 0
                    && statsTasksFilterKey === filterKey
                    && (Date.now() - statsTasksLoadedAt) < STATS_TASKS_TTL_MS;

                if (hasFreshCache) {
                    return appState.allTasks;
                }

                const tasks = await fetchAllPages('/tasks', {
                    params: {
                        user_id: tasksFilterUserId
                    },
                    signal
                });

                appState.allTasks = Array.isArray(tasks) ? tasks : [];
                statsTasksLoadedAt = Date.now();
                statsTasksFilterKey = filterKey;

                return appState.allTasks;
            }

            function renderStatsChartPlaceholder(message) {
                const chartContainer = document.querySelector('.chart-mock-container');
                if (!chartContainer) return;

                chartContainer.innerHTML = `<p style="color: var(--text-muted); text-align: center; padding: 20px;">${escapeHtml(message)}</p>`;
            }

            async function ensureStatsChartData({ force = false } = {}) {
                if (!isStatsSectionActive()) {
                    return;
                }

                const filterKey = getStatsTasksFilterKey();
                const hasFreshCache = !force
                    && statsTasksLoadedAt > 0
                    && statsTasksFilterKey === filterKey
                    && (Date.now() - statsTasksLoadedAt) < STATS_TASKS_TTL_MS;

                if (hasFreshCache) {
                    return;
                }

                const requestId = ++statsLoadRequestId;

                if (activeStatsController) {
                    activeStatsController.abort();
                }

                const controller = new AbortController();
                activeStatsController = controller;

                renderStatsChartPlaceholder('Загрузка диаграммы...');

                try {
                    await loadStatsTasksData({ force: true, signal: controller.signal });

                    if (requestId !== statsLoadRequestId) {
                        return;
                    }

                    updateStatistics();
                } catch (e) {
                    if (e.name === 'AbortError') {
                        return;
                    }

                    console.error('Ошибка загрузки данных для диаграммы:', e);
                    renderStatsChartPlaceholder('Не удалось загрузить диаграмму');
                } finally {
                    if (activeStatsController === controller) {
                        activeStatsController = null;
                    }
                }
            }

            // Функция проверки соединения и загрузки данных
            async function loadAppData({ refreshUsersDirectory = false } = {}) {
                const requestId = ++loadRequestId;

                if (activeLoadController) {
                    activeLoadController.abort();
                }

                const controller = new AbortController();
                activeLoadController = controller;

                try {
                    let nextTasksPage = tasksPage;
                    let nextUsersPage = usersPage;
                    let shouldRefreshUsersDirectory = refreshUsersDirectory;
                    let usersDirectory = appState.users;
                    let paginatedUsers = { items: [], hasNextPage: false };
                    let tasks = { items: [], hasNextPage: false };
                    let stats = null;

                    while (true) {
                        const statsSearchParams = new URLSearchParams();

                        if (statsFilterUserId) statsSearchParams.set('user_id', statsFilterUserId);
                        if (statsFilterFromDate) statsSearchParams.set('from', statsFilterFromDate);
                        if (statsFilterToDate) statsSearchParams.set('to', statsFilterToDate);

                        const statsUrl = statsSearchParams.size > 0
                            ? `${API_BASE_URL}/statistics?${statsSearchParams.toString()}`
                            : `${API_BASE_URL}/statistics`;

                        [usersDirectory, paginatedUsers, tasks, stats] = await Promise.all([
                            loadUsersDirectory({
                                force: shouldRefreshUsersDirectory,
                                signal: controller.signal
                            }),
                            fetchPaginatedSlice('/users', {
                                page: nextUsersPage,
                                limit: usersLimit,
                                signal: controller.signal
                            }),
                            fetchPaginatedSlice('/tasks', {
                                page: nextTasksPage,
                                limit: tasksLimit,
                                params: {
                                    user_id: tasksFilterUserId
                                },
                                signal: controller.signal
                            }),
                            fetchJson(statsUrl, { signal: controller.signal })
                        ]);
                        shouldRefreshUsersDirectory = false;

                        const tasksList = Array.isArray(tasks.items) ? tasks.items : [];
                        const usersList = Array.isArray(paginatedUsers.items) ? paginatedUsers.items : [];

                        // Если страница пустая после удаления элементов и мы не на первой странице, сбрасываем страницу на предыдущую
                        if (tasksList.length === 0 && nextTasksPage > 1) {
                            nextTasksPage--;
                            continue;
                        }

                        if (usersList.length === 0 && nextUsersPage > 1) {
                            nextUsersPage--;
                            continue;
                        }

                        break;
                    }

                    if (requestId !== loadRequestId) {
                        return;
                    }

                    tasksPage = nextTasksPage;
                    usersPage = nextUsersPage;

                    // Обновляем состояние
                    appState.users = Array.isArray(usersDirectory) ? usersDirectory : [];
                    appState.paginatedUsers = Array.isArray(paginatedUsers.items) ? paginatedUsers.items : [];
                    appState.usersHasNextPage = Boolean(paginatedUsers.hasNextPage);
                    appState.tasks = Array.isArray(tasks.items) ? tasks.items : [];
                    appState.tasksHasNextPage = Boolean(tasks.hasNextPage);
                    appState.statistics = stats || {
                        tasks_created: 0,
                        task_completed: 0,
                        task_completed_rate: 0,
                        task_average_completion_time: null
                    };

                    // Соединение успешно
                    offlineOverlay.classList.remove('active');

                    // Рендерим интерфейс
                    renderAll();

                    if (isStatsSectionActive()) {
                        ensureStatsChartData().catch(console.error);
                    }

                } catch (e) {
                    if (e.name === 'AbortError') {
                        return;
                    }

                    console.error('Ошибка подключения к бэкенду:', e);
                    // Бэкенд недоступен — блокируем приложение
                    offlineOverlay.classList.add('active');
                } finally {
                    if (activeLoadController === controller) {
                        activeLoadController = null;
                    }
                }
            }

            async function refreshAfterTaskChange({ resetTasksPage = false } = {}) {
                if (resetTasksPage) {
                    tasksPage = 1;
                }

                invalidateStatsTasksData();
                await loadAppData();
            }

            async function refreshAfterUserChange({ resetUsersPage = false } = {}) {
                if (resetUsersPage) {
                    usersPage = 1;
                }

                invalidateUsersDirectory();
                await loadAppData({ refreshUsersDirectory: true });
            }

            // Повторная попытка подключения по кнопке с визуальной анимацией загрузки
            if (retryBtn) {
                retryBtn.addEventListener('click', async () => {
                    const btnText = retryBtn.querySelector('span');
                    
                    retryBtn.classList.add('btn-spinning');
                    retryBtn.disabled = true;
                    if (btnText) btnText.textContent = 'Проверка соединения...';
                    
                    // Небольшая задержка для реалистичного эффекта подключения
                    await new Promise(resolve => setTimeout(resolve, 800));
                    
                    await loadAppData();
                    
                    // Если соединение все еще отсутствует, восстанавливаем кнопку
                    retryBtn.classList.remove('btn-spinning');
                    retryBtn.disabled = false;
                    if (btnText) btnText.textContent = 'Повторить попытку подключения';
                });
            }



            // --- Функционал навигации ---
            const navItems = document.querySelectorAll('.nav-item');
            const sections = document.querySelectorAll('.content-section');
            const pageTitle = document.getElementById('page-title');

            navItems.forEach(item => {
                item.addEventListener('click', () => {
                    navItems.forEach(nav => nav.classList.remove('active'));
                    item.classList.add('active');

                    sections.forEach(section => section.classList.remove('active'));
                    
                    const targetId = `section-${item.dataset.target}`;
                    const targetSection = document.getElementById(targetId);
                    if (targetSection) {
                        targetSection.classList.add('active');
                    }

                    if (item.dataset.target === 'stats') {
                        updateStatistics();
                        ensureStatsChartData().catch(console.error);
                    }

                    const itemText = item.querySelector('span').textContent;
                    pageTitle.textContent = itemText;
                });
            });

            // --- Рендеринг задач ---
            const tasksContainer = document.getElementById('tasks-container');
            const newTaskUserSelect = document.getElementById('new-task-user');

            // Состояние поиска и пагинации кастомного выпадающего списка
            let dropdownSearchQuery = '';
            let dropdownPage = 1;
            const dropdownLimit = 5;

            // Вспомогательная функция для обновления отображения выбранного пользователя в кастомном триггере
            function newTaskSelectUserVal(user) {
                const trigger = document.getElementById('custom-user-select-trigger');
                if (!trigger) return;

                if (!user) {
                    trigger.innerHTML = '<span class="custom-select-placeholder">Выберите исполнителя...</span>';
                    return;
                }

                const firstLetter = user.full_name ? user.full_name.charAt(0).toUpperCase() : 'U';
                trigger.innerHTML = `
                    <div style="display: flex; align-items: center; gap: 12px;">
                        <div class="custom-select-user-avatar" style="width: 24px; height: 24px; font-size: 10px;">${firstLetter}</div>
                        <div style="display: flex; flex-direction: column; text-align: left;">
                            <span class="custom-select-user-name" style="font-size: 13px; font-weight: 500;">${escapeHtml(user.full_name)}</span>
                            <span class="custom-select-user-id" style="font-size: 11px; margin-top: -2px;">ID: ${user.id}</span>
                        </div>
                    </div>
                `;
            }

            function populateUserSelect() {
                if (!newTaskUserSelect) return;
                
                const selectedValue = newTaskUserSelect.value;
                
                // 1. Очищаем и заполняем оригинальный select
                newTaskUserSelect.innerHTML = '<option value="" disabled selected>Исполнитель...</option>';
                
                appState.users.forEach(user => {
                    const option = document.createElement('option');
                    option.value = user.id;
                    option.textContent = user.full_name;
                    newTaskUserSelect.appendChild(option);
                });

                if (selectedValue && appState.users.some(u => u.id == selectedValue)) {
                    newTaskUserSelect.value = selectedValue;
                } else {
                    newTaskSelectUserVal(null);
                }

                // Заполняем фильтр пользователей в статистике
                const statsUserFilter = document.getElementById('stats-user-filter');
                if (statsUserFilter) {
                    const currentStatsUser = statsUserFilter.value;
                    statsUserFilter.innerHTML = '<option value="">Все пользователи</option>';
                    appState.users.forEach(user => {
                        const option = document.createElement('option');
                        option.value = user.id;
                        option.textContent = user.full_name;
                        statsUserFilter.appendChild(option);
                    });
                    if (currentStatsUser && appState.users.some(u => u.id == currentStatsUser)) {
                        statsUserFilter.value = currentStatsUser;
                    }
                }

                // Заполняем фильтр пользователей в задачах
                const tasksUserFilter = document.getElementById('tasks-user-filter');
                if (tasksUserFilter) {
                    const currentTasksUser = tasksUserFilter.value;
                    tasksUserFilter.innerHTML = '<option value="">Все пользователи</option>';
                    appState.users.forEach(user => {
                        const option = document.createElement('option');
                        option.value = user.id;
                        option.textContent = user.full_name;
                        tasksUserFilter.appendChild(option);
                    });
                    if (currentTasksUser && appState.users.some(u => u.id == currentTasksUser)) {
                        tasksUserFilter.value = currentTasksUser;
                    }
                }

                // Сбрасываем фильтры при обновлении данных о пользователях
                dropdownSearchQuery = '';
                dropdownPage = 1;
                const searchInput = document.getElementById('custom-user-select-search');
                if (searchInput) searchInput.value = '';

                // 2. Вызываем рендеринг элементов кастомного списка
                renderDropdownItems();
            }

            // Функция отрисовки элементов и пагинации внутри кастомного списка
            function renderDropdownItems() {
                const itemsContainer = document.getElementById('custom-user-select-items');
                const paginationContainer = document.getElementById('custom-user-select-pagination');
                const customSelectContainer = document.getElementById('custom-user-select-container');
                if (!itemsContainer || !paginationContainer || !customSelectContainer) return;

                // 1. Фильтруем пользователей по поисковому запросу
                const filteredUsers = appState.users.filter(user => 
                    user.full_name && user.full_name.toLowerCase().includes(dropdownSearchQuery.toLowerCase())
                );

                // 2. Рассчитываем пагинацию
                const totalDropdownPages = Math.ceil(filteredUsers.length / dropdownLimit) || 1;
                if (dropdownPage > totalDropdownPages) {
                    dropdownPage = totalDropdownPages;
                }

                const startIndex = (dropdownPage - 1) * dropdownLimit;
                const paginatedUsers = filteredUsers.slice(startIndex, startIndex + dropdownLimit);

                // 3. Рендерим элементы
                itemsContainer.innerHTML = '';

                if (filteredUsers.length === 0) {
                    itemsContainer.innerHTML = '<div style="padding: 12px; font-size: 13px; color: var(--text-secondary); text-align: center;">Ничего не найдено</div>';
                    paginationContainer.innerHTML = '';
                    return;
                }

                paginatedUsers.forEach(user => {
                    const optionDiv = document.createElement('div');
                    optionDiv.className = 'custom-select-option';
                    if (user.id == newTaskUserSelect.value) {
                        optionDiv.classList.add('selected');
                        newTaskSelectUserVal(user);
                    }

                    const firstLetter = user.full_name ? user.full_name.charAt(0).toUpperCase() : 'U';

                    optionDiv.innerHTML = `
                        <div class="custom-select-user-avatar">${firstLetter}</div>
                        <div class="custom-select-user-info">
                            <span class="custom-select-user-name">${escapeHtml(user.full_name)}</span>
                            <span class="custom-select-user-id">ID: ${user.id}</span>
                        </div>
                    `;

                    optionDiv.addEventListener('click', (e) => {
                        e.stopPropagation();
                        
                        itemsContainer.querySelectorAll('.custom-select-option').forEach(opt => {
                            opt.classList.remove('selected');
                        });
                        
                        optionDiv.classList.add('selected');
                        newTaskUserSelect.value = user.id;
                        newTaskSelectUserVal(user);
                        customSelectContainer.classList.remove('active');
                    });

                    itemsContainer.appendChild(optionDiv);
                });

                // 4. Рендерим пагинацию
                paginationContainer.innerHTML = '';

                // Кнопка назад
                const prevBtn = document.createElement('button');
                prevBtn.className = 'custom-select-pagination-btn';
                prevBtn.type = 'button';
                prevBtn.innerHTML = '‹';
                prevBtn.disabled = dropdownPage === 1;
                prevBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    if (dropdownPage > 1) {
                        dropdownPage--;
                        renderDropdownItems();
                    }
                });
                paginationContainer.appendChild(prevBtn);

                // Номера страниц
                const numbersDiv = document.createElement('div');
                numbersDiv.className = 'custom-select-pagination-numbers';
                
                for (let i = 1; i <= totalDropdownPages; i++) {
                    const pageNum = document.createElement('button');
                    pageNum.className = 'custom-select-pagination-num';
                    pageNum.type = 'button';
                    if (i === dropdownPage) {
                        pageNum.classList.add('active');
                    }
                    pageNum.textContent = i;
                    pageNum.addEventListener('click', (e) => {
                        e.stopPropagation();
                        if (dropdownPage !== i) {
                            dropdownPage = i;
                            renderDropdownItems();
                        }
                    });
                    numbersDiv.appendChild(pageNum);
                }
                paginationContainer.appendChild(numbersDiv);

                // Кнопка вперед
                const nextBtn = document.createElement('button');
                nextBtn.className = 'custom-select-pagination-btn';
                nextBtn.type = 'button';
                nextBtn.innerHTML = '›';
                nextBtn.disabled = dropdownPage >= totalDropdownPages;
                nextBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    if (dropdownPage < totalDropdownPages) {
                        dropdownPage++;
                        renderDropdownItems();
                    }
                });
                paginationContainer.appendChild(nextBtn);
            }

            // Вспомогательная функция для форматирования даты в формат "10 november 2026"
            function formatDate(dateString) {
                if (!dateString) return '';
                const date = new Date(dateString);
                if (isNaN(date.getTime())) return '';
                
                const day = date.getDate();
                const calendarMonths = [
                    'january', 'february', 'march', 'april', 'may', 'june',
                    'july', 'august', 'september', 'october', 'november', 'december'
                ];
                const month = calendarMonths[date.getMonth()];
                const year = date.getFullYear();
                
                return `${day} ${month} ${year}`;
            }

            // Форматирование для бейджей в стиле референса
            function formatBadgeDate(dateString) {
                if (!dateString) return '';
                const date = new Date(dateString);
                if (isNaN(date.getTime())) return '';
                const calendarMonths = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
                return `${calendarMonths[date.getMonth()]} ${date.getDate()}, ${date.getFullYear()}`;
            }

            function formatBadgeDateTime(dateString) {
                if (!dateString) return '';
                const date = new Date(dateString);
                if (isNaN(date.getTime())) return '';
                
                const calendarMonths = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
                const datePart = `${calendarMonths[date.getMonth()]} ${date.getDate()}`;
                
                let hours = date.getHours();
                const minutes = date.getMinutes().toString().padStart(2, '0');
                const ampm = hours >= 12 ? 'PM' : 'AM';
                hours = hours % 12;
                hours = hours ? hours : 12;
                
                return `${datePart} ${hours}:${minutes} ${ampm}`;
            }

            // Вспомогательная функция для расчета времени нахождения задачи в работе
            function getDurationInWork(createdStr, completedStr, isCompleted) {
                if (!createdStr) return '';
                
                const createdDate = new Date(createdStr);
                if (isNaN(createdDate.getTime())) return '';
                
                let endDate = new Date();
                if (isCompleted && completedStr) {
                    const compDate = new Date(completedStr);
                    if (!isNaN(compDate.getTime())) {
                        endDate = compDate;
                    }
                }
                
                const diffMs = endDate - createdDate;
                if (diffMs < 0) return '0 мин';
                
                const diffMins = Math.floor(diffMs / (1000 * 60));
                const hours = Math.floor(diffMins / 60);
                const mins = diffMins % 60;
                
                if (hours > 0) {
                    return `${hours} ч ${mins} мин`;
                } else {
                    return `${mins} мин`;
                }
            }

            function renderTasks() {
                tasksContainer.innerHTML = '';
                
                if (appState.tasks.length === 0) {
                    tasksContainer.innerHTML = `
                        <div style="text-align: center; padding: 40px; color: var(--text-muted);">
                            <i data-lucide="check-circle" width="48" height="48" style="margin-bottom: 12px; stroke-width: 1.5"></i>
                            <p>Список задач пуст.</p>
                        </div>
                    `;
                    safeCreateIcons();
                    return;
                }

                // Сортировка: сначала невыполненные, затем новые
                const sortedTasks = [...appState.tasks].sort((a, b) => {
                    if (a.completed !== b.completed) {
                        return a.completed ? 1 : -1;
                    }
                    return b.id - a.id;
                });

                sortedTasks.forEach(task => {
                    const taskCard = document.createElement('div');
                    taskCard.className = `task-card ${task.completed ? 'completed' : ''}`;
                    
                    const author = appState.users.find(u => u.id === task.author_user_id);
                    const authorName = author ? author.full_name : `Пользователь #${task.author_user_id}`;

                    taskCard.innerHTML = `
                        <div class="task-card-left" style="width: 100%;">
                            <label class="checkbox-container" style="margin-right: 16px;">
                                <input type="checkbox" ${task.completed ? 'checked' : ''} data-id="${task.id}" class="task-toggle">
                                <span class="checkmark" style="width: 28px; height: 28px; border-radius: 50%; border: none; background-color: ${task.completed ? 'var(--success)' : 'rgba(255, 255, 255, 0.1)'}; display: flex; align-items: center; justify-content: center; box-shadow: ${task.completed ? '0 0 12px rgba(16, 185, 129, 0.4)' : 'none'}">
                                    ${task.completed ? '<i data-lucide="check" style="color: white; stroke-width: 3px;" width="16" height="16"></i>' : ''}
                                </span>
                            </label>
                            <div style="display: flex; flex-direction: column; flex: 1;">
                                <div style="display: flex; justify-content: space-between; align-items: flex-start;">
                                    <span class="task-text" style="font-size: 16px; font-weight: 500; color: ${task.completed ? 'var(--text-muted)' : 'var(--text-primary)'}; text-decoration: ${task.completed ? 'line-through' : 'none'};">${escapeHtml(task.title)}</span>
                                    <div class="task-actions" style="display: flex; gap: 8px; opacity: 0.7;">
                                        <button class="btn-edit-task" data-id="${task.id}" title="Редактировать" style="background: none; border: none; color: var(--text-secondary); cursor: pointer; padding: 4px;">
                                            <i data-lucide="pencil" width="16" height="16"></i>
                                        </button>
                                        <button class="btn-delete" data-id="${task.id}" title="Удалить" style="background: none; border: none; color: var(--text-secondary); cursor: pointer; padding: 4px;">
                                            <i data-lucide="trash-2" width="16" height="16"></i>
                                        </button>
                                    </div>
                                </div>
                                ${task.description ? `<span style="font-size: 14px; color: var(--text-secondary); margin-top: 6px; margin-bottom: 4px;">${escapeHtml(task.description)}</span>` : ''}
                                
                                <div style="display: flex; align-items: center; gap: 8px; font-size: 12px; margin-top: 12px; flex-wrap: wrap;">
                                    <!-- Пользователь -->
                                    <span style="display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; background-color: rgba(255, 255, 255, 0.05); border-radius: 16px; color: var(--text-secondary);">
                                        <i data-lucide="user" width="14" height="14" style="opacity: 0.7;"></i>
                                        ${escapeHtml(authorName)}
                                        <span style="color: rgba(255,255,255,0.3); font-size: 11px; margin-left: 2px;">ID: ${task.author_user_id}</span>
                                    </span>

                                    <!-- Дата создания -->
                                    <span style="display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; background-color: rgba(255, 255, 255, 0.05); border-radius: 16px; color: var(--text-secondary);">
                                        <i data-lucide="calendar" width="14" height="14" style="opacity: 0.7;"></i>
                                        ${formatBadgeDate(task.created_at)}
                                    </span>

                                    <!-- Статус и время -->
                                    ${task.completed ? `
                                        <span style="display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; background-color: rgba(16, 185, 129, 0.1); border-radius: 16px; color: var(--success);">
                                            <i data-lucide="check" width="14" height="14"></i>
                                            Готово • ${formatBadgeDateTime(task.completed_at)}
                                        </span>
                                        <span style="display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; background-color: rgba(59, 130, 246, 0.15); border-radius: 16px; color: #60a5fa;">
                                            <i data-lucide="clock" width="14" height="14"></i>
                                            ${getDurationInWork(task.created_at, task.completed_at, true)}
                                        </span>
                                    ` : `
                                        <span style="display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; background-color: rgba(245, 158, 11, 0.1); border-radius: 16px; color: var(--warning);">
                                            <i data-lucide="clock" width="14" height="14"></i>
                                            В работе: ${getDurationInWork(task.created_at, task.completed_at, false)}
                                        </span>
                                    `}
                                </div>
                            </div>
                        </div>
                    `;

                    // Событие редактирования задачи
                    taskCard.querySelector('.btn-edit-task').addEventListener('click', () => {
                        openEditTaskModal(task);
                    });

                    // Событие удаления задачи
                    taskCard.querySelector('.btn-delete').addEventListener('click', () => {
                        deleteTask(task.id);
                    });

                    // Событие переключения выполнения
                    taskCard.querySelector('.task-toggle').addEventListener('change', (e) => {
                        toggleTask(task.id, e.target.checked);
                    });

                    tasksContainer.appendChild(taskCard);
                });

                safeCreateIcons();
                renderTasksPagination();
            }

            // Вспомогательная функция для генерации номеров страниц с многоточием
            function getPageNumbers(currentPage, totalPages) {
                const pages = [];
                if (totalPages <= 7) {
                    for (let i = 1; i <= totalPages; i++) {
                        pages.push(i);
                    }
                } else {
                    if (currentPage <= 4) {
                        pages.push(1, 2, 3, 4, 5, '...', totalPages);
                    } else if (currentPage >= totalPages - 3) {
                        pages.push(1, '...', totalPages - 4, totalPages - 3, totalPages - 2, totalPages - 1, totalPages);
                    } else {
                        pages.push(1, '...', currentPage - 1, currentPage, currentPage + 1, '...', totalPages);
                    }
                }
                return pages;
            }

            function getIncrementalPageNumbers(currentPage, hasNextPage) {
                const pages = [];

                if (currentPage > 1) {
                    pages.push(1);
                }

                if (currentPage > 3) {
                    pages.push('...');
                }

                if (currentPage > 2) {
                    pages.push(currentPage - 1);
                }

                if (currentPage > 1) {
                    pages.push(currentPage);
                } else {
                    pages.push(1);
                }

                if (hasNextPage) {
                    pages.push(currentPage + 1);
                }

                return [...new Set(pages)];
            }

            function renderTasksPagination() {
                const prevBtn = document.getElementById('tasks-prev-btn');
                const nextBtn = document.getElementById('tasks-next-btn');
                const pageNumbersContainer = document.getElementById('tasks-page-numbers');
                const paginationContainer = document.getElementById('tasks-pagination');

                if (!prevBtn || !nextBtn || !pageNumbersContainer || !paginationContainer) return;
                
                // Очистка старых номеров страниц
                pageNumbersContainer.innerHTML = '';
                
                // Генерация массива номеров страниц
                const pages = getIncrementalPageNumbers(tasksPage, appState.tasksHasNextPage);
                
                pages.forEach(page => {
                    if (page === '...') {
                        const ellipsis = document.createElement('span');
                        ellipsis.className = 'pagination-ellipsis';
                        ellipsis.textContent = '...';
                        pageNumbersContainer.appendChild(ellipsis);
                    } else {
                        const btn = document.createElement('button');
                        btn.className = 'pagination-number-btn';
                        if (page === tasksPage) {
                            btn.classList.add('active');
                        }
                        btn.textContent = page;
                        btn.addEventListener('click', async () => {
                            if (tasksPage !== page) {
                                tasksPage = page;
                                await loadAppData();
                            }
                        });
                        pageNumbersContainer.appendChild(btn);
                    }
                });
                
                prevBtn.disabled = tasksPage === 1;
                nextBtn.disabled = !appState.tasksHasNextPage;

                if (tasksPage === 1 && !appState.tasksHasNextPage) {
                    paginationContainer.style.display = 'none';
                } else {
                    paginationContainer.style.display = 'flex';
                }
            }

            // --- События пагинации задач ---
            const tasksPrevBtn = document.getElementById('tasks-prev-btn');
            const tasksNextBtn = document.getElementById('tasks-next-btn');
            const tasksLimitSelect = document.getElementById('tasks-limit-select');

            if (tasksPrevBtn) {
                tasksPrevBtn.addEventListener('click', async () => {
                    if (tasksPage > 1) {
                        tasksPage--;
                        await loadAppData();
                    }
                });
            }

            if (tasksNextBtn) {
                tasksNextBtn.addEventListener('click', async () => {
                    if (appState.tasksHasNextPage) {
                        tasksPage++;
                        await loadAppData();
                    }
                });
            }

            if (tasksLimitSelect) {
                tasksLimitSelect.addEventListener('change', async (e) => {
                    tasksLimit = parseInt(e.target.value);
                    tasksPage = 1;
                    await loadAppData();
                });
            }

            // --- Операции над задачами (API) ---
            const addTaskForm = document.getElementById('add-task-form');
            const newTaskInput = document.getElementById('new-task-input');
            const newTaskDesc = document.getElementById('new-task-desc');

            // --- Управление модальным окном добавления задачи ---
            const addTaskModal = document.getElementById('add-task-modal');
            const openAddTaskModalBtn = document.getElementById('open-add-task-modal');
            const closeAddTaskModalIcon = document.getElementById('close-add-task-modal');
            const cancelAddTaskModalBtn = document.getElementById('cancel-add-task-modal');

            function toggleAddTaskModal(show) {
                if (show) {
                    addTaskModal.classList.add('active');
                } else {
                    addTaskModal.classList.remove('active');
                    addTaskForm.reset();
                    // Сбрасываем кастомный выпадающий список
                    newTaskSelectUserVal(null);
                    const customSelectContainer = document.getElementById('custom-user-select-container');
                    if (customSelectContainer) {
                        customSelectContainer.classList.remove('active');
                    }
                    // Сброс поиска и страницы в выпадающем списке
                    dropdownSearchQuery = '';
                    dropdownPage = 1;
                    const searchInput = document.getElementById('custom-user-select-search');
                    if (searchInput) searchInput.value = '';
                }
            }

            if (openAddTaskModalBtn) openAddTaskModalBtn.addEventListener('click', () => toggleAddTaskModal(true));
            if (closeAddTaskModalIcon) closeAddTaskModalIcon.addEventListener('click', () => toggleAddTaskModal(false));
            if (cancelAddTaskModalBtn) cancelAddTaskModalBtn.addEventListener('click', () => toggleAddTaskModal(false));
            
            if (addTaskModal) {
                addTaskModal.addEventListener('click', (e) => {
                    if (e.target === addTaskModal) toggleAddTaskModal(false);
                });
            }

            // --- Кастомный красивый выпадающий список ---
            const customSelectContainer = document.getElementById('custom-user-select-container');
            const customSelectTrigger = document.getElementById('custom-user-select-trigger');
            
            if (customSelectTrigger && customSelectContainer) {
                customSelectTrigger.addEventListener('click', (e) => {
                    e.stopPropagation();
                    customSelectContainer.classList.toggle('active');
                });
            }

            // Инициализация поиска в кастомном выпадающем списке
            const customSearchInput = document.getElementById('custom-user-select-search');
            if (customSearchInput) {
                customSearchInput.addEventListener('input', (e) => {
                    dropdownSearchQuery = e.target.value;
                    dropdownPage = 1; // Сброс на первую страницу при поиске
                    renderDropdownItems();
                });
                
                // Предотвращаем закрытие списка при клике на поле поиска
                customSearchInput.addEventListener('click', (e) => {
                    e.stopPropagation();
                });
            }

            document.addEventListener('click', () => {
                if (customSelectContainer) {
                    customSelectContainer.classList.remove('active');
                }
            });

            if (addTaskForm) {
                addTaskForm.addEventListener('submit', async (e) => {
                    e.preventDefault();
                    
                    const title = newTaskInput.value.trim();
                    const desc = newTaskDesc.value.trim();
                    const userId = newTaskUserSelect.value;

                    if (!userId) {
                        alert('Пожалуйста, выберите исполнителя.');
                        return;
                    }

                    try {
                        const response = await fetch(`${API_BASE_URL}/tasks`, {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({
                                title,
                                description: desc ? desc : null,
                                author_user_id: parseInt(userId)
                            })
                        });

                        if (!response.ok) {
                            const errData = await response.json();
                            throw new Error(errData.message || 'Ошибка сервера');
                        }

                        addTaskForm.reset();
                        
                        // Закрываем модальное окно добавления задачи
                        toggleAddTaskModal(false);
                        
                        await refreshAfterTaskChange({ resetTasksPage: true });

                    } catch (err) {
                        console.error('Ошибка добавления задачи:', err);
                        alert(`Не удалось добавить задачу: ${err.message}`);
                    }
                });
            }

            async function toggleTask(id, completed) {
                try {
                    const response = await fetch(`${API_BASE_URL}/tasks/${id}`, {
                        method: 'PATCH',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            completed: completed
                        })
                    });

                    if (!response.ok) {
                        throw new Error('Ошибка обновления статуса задачи');
                    }

                    await refreshAfterTaskChange();
                } catch (err) {
                    console.error('Ошибка переключения задачи:', err);
                    alert(err.message);
                    // Возвращаем интерфейс в исходное состояние
                    renderTasks();
                }
            }

            async function deleteTask(id) {
                if (!confirm('Вы действительно хотите удалить эту задачу?')) return;
                try {
                    const response = await fetch(`${API_BASE_URL}/tasks/${id}`, {
                        method: 'DELETE'
                    });

                    if (!response.ok) {
                        throw new Error('Ошибка удаления задачи');
                    }

                    await refreshAfterTaskChange();
                } catch (err) {
                    console.error('Ошибка удаления задачи:', err);
                    alert(err.message);
                }
            }

            // --- Управление модальным окном редактирования задач (API) ---
            const editTaskModal = document.getElementById('edit-task-modal');
            const closeEditTaskModalIcon = document.getElementById('close-edit-task-modal');
            const cancelEditTaskModalBtn = document.getElementById('cancel-edit-task-modal');
            const editTaskForm = document.getElementById('edit-task-form');

            function toggleEditTaskModal(show) {
                if (show) {
                    editTaskModal.classList.add('active');
                } else {
                    editTaskModal.classList.remove('active');
                    editTaskForm.reset();
                }
            }

            if (closeEditTaskModalIcon) closeEditTaskModalIcon.addEventListener('click', () => toggleEditTaskModal(false));
            if (cancelEditTaskModalBtn) cancelEditTaskModalBtn.addEventListener('click', () => toggleEditTaskModal(false));
            
            if (editTaskModal) {
                editTaskModal.addEventListener('click', (e) => {
                    if (e.target === editTaskModal) toggleEditTaskModal(false);
                });
            }

            function openEditTaskModal(task) {
                document.getElementById('edit-task-id').value = task.id;
                document.getElementById('edit-task-input-title').value = task.title;
                document.getElementById('edit-task-input-desc').value = task.description || '';
                toggleEditTaskModal(true);
            }

            if (editTaskForm) {
                editTaskForm.addEventListener('submit', async (e) => {
                    e.preventDefault();
                    
                    const id = document.getElementById('edit-task-id').value;
                    const title = document.getElementById('edit-task-input-title').value.trim();
                    const desc = document.getElementById('edit-task-input-desc').value.trim();

                    try {
                        const response = await fetch(`${API_BASE_URL}/tasks/${id}`, {
                            method: 'PATCH',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({
                                title: title,
                                description: desc ? desc : null
                            })
                        });
                        
                        if (!response.ok) {
                            const errData = await response.json();
                            throw new Error(errData.message || errData.error || 'Ошибка сервера при обновлении задачи');
                        }

                        toggleEditTaskModal(false);
                        await refreshAfterTaskChange();

                    } catch (err) {
                        console.error('Ошибка редактирования задачи:', err);
                        alert(`Не удалось изменить задачу: ${err.message}`);
                    }
                });
            }

            // --- Рендеринг пользователей ---
            const usersContainer = document.getElementById('users-container');

            function renderUsers() {
                usersContainer.innerHTML = '';
                
                if (appState.paginatedUsers.length === 0) {
                    usersContainer.innerHTML = `
                        <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: var(--text-muted);">
                            <i data-lucide="users" width="48" height="48" style="margin-bottom: 12px; stroke-width: 1.5"></i>
                            <p>Список пользователей пуст. Добавьте первого участника!</p>
                        </div>
                    `;
                    safeCreateIcons();
                    renderUsersPagination();
                    return;
                }

                appState.paginatedUsers.forEach(user => {
                    const userCard = document.createElement('div');
                    userCard.className = 'user-card';
                    
                    const firstLetter = user.full_name ? user.full_name.charAt(0).toUpperCase() : 'U';
                    const phone = user.phone_number ? user.phone_number : 'Телефон не указан';
                    
                    userCard.innerHTML = `
                        <div class="user-status-dot"></div>
                        <div class="user-card-avatar">${firstLetter}</div>
                        <h3 class="user-card-name">${escapeHtml(user.full_name)}</h3>
                        <span class="user-card-email">${escapeHtml(phone)}</span>
                        <span class="user-card-role">ID: ${user.id}</span>
                        <div class="user-card-actions">
                            <button class="btn-user-action btn-edit-user" title="Редактировать">
                                <i data-lucide="pencil" width="14" height="14"></i>
                            </button>
                            <button class="btn-user-action btn-delete-user" title="Удалить">
                                <i data-lucide="trash-2" width="14" height="14"></i>
                            </button>
                        </div>
                    `;
                    
                    // Навешиваем обработчики событий
                    userCard.querySelector('.btn-edit-user').addEventListener('click', () => {
                        openEditUserModal(user);
                    });
                    
                    userCard.querySelector('.btn-delete-user').addEventListener('click', () => {
                        deleteUser(user.id, user.full_name);
                    });

                    usersContainer.appendChild(userCard);
                });
                safeCreateIcons();
                renderUsersPagination();
            }

            function renderUsersPagination() {
                const prevBtn = document.getElementById('users-prev-btn');
                const nextBtn = document.getElementById('users-next-btn');
                const pageNumbersContainer = document.getElementById('users-page-numbers');
                const paginationContainer = document.getElementById('users-pagination');

                if (!prevBtn || !nextBtn || !pageNumbersContainer || !paginationContainer) return;
                
                // Очистка старых номеров страниц
                pageNumbersContainer.innerHTML = '';
                
                // Генерация массива номеров страниц
                const pages = getIncrementalPageNumbers(usersPage, appState.usersHasNextPage);
                
                pages.forEach(page => {
                    if (page === '...') {
                        const ellipsis = document.createElement('span');
                        ellipsis.className = 'pagination-ellipsis';
                        ellipsis.textContent = '...';
                        pageNumbersContainer.appendChild(ellipsis);
                    } else {
                        const btn = document.createElement('button');
                        btn.className = 'pagination-number-btn';
                        if (page === usersPage) {
                            btn.classList.add('active');
                        }
                        btn.textContent = page;
                        btn.addEventListener('click', async () => {
                            if (usersPage !== page) {
                                usersPage = page;
                                await loadAppData();
                            }
                        });
                        pageNumbersContainer.appendChild(btn);
                    }
                });
                
                prevBtn.disabled = usersPage === 1;
                nextBtn.disabled = !appState.usersHasNextPage;

                if (usersPage === 1 && !appState.usersHasNextPage) {
                    paginationContainer.style.display = 'none';
                } else {
                    paginationContainer.style.display = 'flex';
                }
            }

            // --- События пагинации пользователей ---
            const usersPrevBtn = document.getElementById('users-prev-btn');
            const usersNextBtn = document.getElementById('users-next-btn');
            const usersLimitSelect = document.getElementById('users-limit-select');

            if (usersPrevBtn) {
                usersPrevBtn.addEventListener('click', async () => {
                    if (usersPage > 1) {
                        usersPage--;
                        await loadAppData();
                    }
                });
            }

            if (usersNextBtn) {
                usersNextBtn.addEventListener('click', async () => {
                    if (appState.usersHasNextPage) {
                        usersPage++;
                        await loadAppData();
                    }
                });
            }

            if (usersLimitSelect) {
                usersLimitSelect.addEventListener('change', async (e) => {
                    usersLimit = parseInt(e.target.value);
                    usersPage = 1;
                    await loadAppData();
                });
            }

            // --- События пагинации диаграммы статистики ---
            const statsChartPrevBtn = document.getElementById('stats-chart-prev-btn');
            const statsChartNextBtn = document.getElementById('stats-chart-next-btn');
            const statsChartLimitSelect = document.getElementById('stats-chart-limit-select');

            if (statsChartPrevBtn) {
                statsChartPrevBtn.addEventListener('click', () => {
                    if (statsChartPage > 1) {
                        statsChartPage--;
                        updateStatistics();
                    }
                });
            }

            if (statsChartNextBtn) {
                statsChartNextBtn.addEventListener('click', () => {
                    const totalUsers = appState.users.length;
                    const totalPages = Math.ceil(totalUsers / statsChartLimit) || 1;
                    if (statsChartPage < totalPages) {
                        statsChartPage++;
                        updateStatistics();
                    }
                });
            }

            if (statsChartLimitSelect) {
                statsChartLimitSelect.addEventListener('change', (e) => {
                    statsChartLimit = parseInt(e.target.value);
                    statsChartPage = 1;
                    updateStatistics();
                });
            }

            // --- Управление модальным окном пользователей (API) ---
            const modal = document.getElementById('add-user-modal');
            const openModalBtn = document.getElementById('open-add-user-modal');
            const closeModalIcon = document.getElementById('close-user-modal');
            const cancelModalBtn = document.getElementById('cancel-user-modal');
            const addUserForm = document.getElementById('add-user-form');

            function toggleModal(show) {
                if (show) {
                    modal.classList.add('active');
                } else {
                    modal.classList.remove('active');
                    addUserForm.reset();
                }
            }

            if (openModalBtn) openModalBtn.addEventListener('click', () => toggleModal(true));
            if (closeModalIcon) closeModalIcon.addEventListener('click', () => toggleModal(false));
            if (cancelModalBtn) cancelModalBtn.addEventListener('click', () => toggleModal(false));
            
            if (modal) {
                modal.addEventListener('click', (e) => {
                    if (e.target === modal) toggleModal(false);
                });
            }

            if (addUserForm) {
                addUserForm.addEventListener('submit', async (e) => {
                    e.preventDefault();
                    
                    const name = document.getElementById('user-input-name').value.trim();
                    const phone = document.getElementById('user-input-phone').value.trim();

                    try {
                        const response = await fetch(`${API_BASE_URL}/users`, {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({
                                full_name: name,
                                phone_number: phone ? phone : null
                            })
                        });
                        
                        if (!response.ok) {
                            const errData = await response.json();
                            throw new Error(errData.message || 'Ошибка сервера при создании пользователя');
                        }

                        toggleModal(false);
                        
                        await refreshAfterUserChange({ resetUsersPage: true });

                    } catch (err) {
                        console.error('Ошибка добавления пользователя:', err);
                        alert(`Не удалось добавить пользователя: ${err.message}`);
                    }
                });
            }

            // --- Управление модальным окном редактирования пользователей (API) ---
            const editModal = document.getElementById('edit-user-modal');
            const closeEditModalIcon = document.getElementById('close-edit-user-modal');
            const cancelEditModalBtn = document.getElementById('cancel-edit-user-modal');
            const editUserForm = document.getElementById('edit-user-form');

            function toggleEditModal(show) {
                if (show) {
                    editModal.classList.add('active');
                } else {
                    editModal.classList.remove('active');
                    editUserForm.reset();
                }
            }

            if (closeEditModalIcon) closeEditModalIcon.addEventListener('click', () => toggleEditModal(false));
            if (cancelEditModalBtn) cancelEditModalBtn.addEventListener('click', () => toggleEditModal(false));
            
            if (editModal) {
                editModal.addEventListener('click', (e) => {
                    if (e.target === editModal) toggleEditModal(false);
                });
            }

            function openEditUserModal(user) {
                document.getElementById('edit-user-id').value = user.id;
                document.getElementById('edit-user-input-name').value = user.full_name;
                document.getElementById('edit-user-input-phone').value = user.phone_number || '';
                toggleEditModal(true);
            }

            if (editUserForm) {
                editUserForm.addEventListener('submit', async (e) => {
                    e.preventDefault();
                    
                    const id = document.getElementById('edit-user-id').value;
                    const name = document.getElementById('edit-user-input-name').value.trim();
                    const phone = document.getElementById('edit-user-input-phone').value.trim();

                    try {
                        const response = await fetch(`${API_BASE_URL}/users/${id}`, {
                            method: 'PATCH',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({
                                full_name: name,
                                phone_number: phone ? phone : null
                            })
                        });
                        
                        if (!response.ok) {
                            const errData = await response.json();
                            throw new Error(errData.message || errData.error || 'Ошибка сервера при обновлении пользователя');
                        }

                        toggleEditModal(false);
                        await refreshAfterUserChange();

                    } catch (err) {
                        console.error('Ошибка редактирования пользователя:', err);
                        alert(`Не удалось изменить пользователя: ${err.message}`);
                    }
                });
            }

            // --- Удаление пользователя (API) ---
            async function deleteUser(id, name) {
                if (!confirm(`Вы действительно хотите удалить пользователя "${name}"?`)) {
                    return;
                }
                try {
                    const response = await fetch(`${API_BASE_URL}/users/${id}`, {
                        method: 'DELETE'
                    });

                    if (!response.ok) {
                        const errData = await response.json();
                        throw new Error(errData.message || errData.error || 'Ошибка сервера при удалении пользователя');
                    }

                    await refreshAfterUserChange();
                } catch (err) {
                    console.error('Ошибка удаления пользователя:', err);
                    alert(`Не удалось удалить пользователя: ${err.message}. Возможно, к нему привязаны задачи.`);
                }
            }

            function renderStatsChartPagination(totalItems) {
                const prevBtn = document.getElementById('stats-chart-prev-btn');
                const nextBtn = document.getElementById('stats-chart-next-btn');
                const pageNumbersContainer = document.getElementById('stats-chart-page-numbers');
                const paginationContainer = document.getElementById('stats-chart-pagination');

                if (!prevBtn || !nextBtn || !pageNumbersContainer || !paginationContainer) return;

                const totalPages = Math.ceil(totalItems / statsChartLimit) || 1;

                // Очистка старых номеров страниц
                pageNumbersContainer.innerHTML = '';

                // Генерация массива номеров страниц
                const pages = getPageNumbers(statsChartPage, totalPages);

                pages.forEach(page => {
                    if (page === '...') {
                        const ellipsis = document.createElement('span');
                        ellipsis.className = 'pagination-ellipsis';
                        ellipsis.textContent = '...';
                        pageNumbersContainer.appendChild(ellipsis);
                    } else {
                        const btn = document.createElement('button');
                        btn.className = 'pagination-number-btn';
                        if (page === statsChartPage) {
                            btn.classList.add('active');
                        }
                        btn.textContent = page;
                        btn.addEventListener('click', () => {
                            if (statsChartPage !== page) {
                                statsChartPage = page;
                                updateStatistics();
                            }
                        });
                        pageNumbersContainer.appendChild(btn);
                    }
                });

                prevBtn.disabled = statsChartPage === 1;
                nextBtn.disabled = statsChartPage >= totalPages;

                if (totalItems === 0) {
                    paginationContainer.style.display = 'none';
                } else {
                    paginationContainer.style.display = 'flex';
                }
            }

            // --- Расчет и отрисовка статистики ---
            function updateStatistics() {
                const stats = appState.statistics;
                
                const totalTasks = stats.tasks_created;
                const completedTasks = stats.task_completed;
                const pendingTasks = totalTasks - completedTasks;
                const totalUsers = appState.users.length;

                // Заполнение KPI карточек
                document.getElementById('stat-total-tasks').textContent = totalTasks;
                document.getElementById('stat-completed-tasks').textContent = completedTasks;
                document.getElementById('stat-pending-tasks').textContent = pendingTasks;
                document.getElementById('stat-total-users').textContent = totalUsers;

                // Скрыть или показать панель "Задачи по исполнителям"
                const tasksByPerformersPanel = document.getElementById('tasks-by-performers-panel');
                if (tasksByPerformersPanel) {
                    if (statsFilterUserId) {
                        tasksByPerformersPanel.style.display = 'none';
                    } else {
                        tasksByPerformersPanel.style.display = 'block';
                    }
                }

                // Задачи по пользователям (динамический рендеринг с локальной пагинацией)
                const chartContainer = document.querySelector('.chart-mock-container');
                if (chartContainer) {
                    chartContainer.innerHTML = '';
                    
                    if (appState.users.length === 0) {
                        chartContainer.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: 20px;">Нет пользователей для отображения</p>';
                        renderStatsChartPagination(0);
                    } else if (!statsFilterUserId && statsTasksLoadedAt === 0) {
                        chartContainer.innerHTML = '<p style="color: var(--text-muted); text-align: center; padding: 20px;">Диаграмма загрузится при открытии раздела статистики</p>';
                        renderStatsChartPagination(0);
                    } else {
                        // Считаем количество задач на каждого пользователя
                        const userTasksCount = {};
                        appState.users.forEach(u => { userTasksCount[u.id] = 0; });
                        appState.allTasks.forEach(t => {
                            if (userTasksCount[t.author_user_id] !== undefined) {
                                userTasksCount[t.author_user_id]++;
                            }
                        });

                        const maxCount = Math.max(...Object.values(userTasksCount), 1);

                        // Сортируем пользователей по количеству задач
                        const sortedUsers = [...appState.users].sort((a, b) => {
                            return (userTasksCount[b.id] || 0) - (userTasksCount[a.id] || 0);
                        });

                        const totalPages = Math.ceil(sortedUsers.length / statsChartLimit) || 1;
                        if (statsChartPage > totalPages) {
                            statsChartPage = totalPages;
                        }

                        const startIndex = (statsChartPage - 1) * statsChartLimit;
                        const endIndex = startIndex + statsChartLimit;
                        const paginatedUsersForChart = sortedUsers.slice(startIndex, endIndex);

                        paginatedUsersForChart.forEach(user => {
                            const count = userTasksCount[user.id] || 0;
                            const percentage = (count / maxCount) * 100;
                            
                            const row = document.createElement('div');
                            row.className = 'chart-row';
                            row.innerHTML = `
                                <span class="chart-label" style="text-overflow: ellipsis; overflow: hidden; white-space: nowrap;" title="${escapeHtml(user.full_name)}">
                                    ${escapeHtml(user.full_name)}
                                </span>
                                <div class="chart-bar-bg">
                                    <div class="chart-bar-fill" style="width: ${percentage}%;"></div>
                                </div>
                                <span class="chart-value">${count}</span>
                            `;
                            chartContainer.appendChild(row);
                        });

                        // Обновляем пагинацию диаграммы
                        renderStatsChartPagination(sortedUsers.length);
                    }
                }

                // Круговой индикатор выполнения
                const completionPercentage = stats.task_completed_rate !== null ? Math.max(0, Math.min(100, Math.round(stats.task_completed_rate))) : 0;
                document.getElementById('radial-percentage').textContent = `${completionPercentage}%`;
                
                const radialFill = document.getElementById('radial-progress-fill');
                if (radialFill) {
                    const circumference = 314;
                    const strokeDashoffset = circumference - (completionPercentage / 100) * circumference;
                    radialFill.style.strokeDashoffset = strokeDashoffset;
                }

                // Среднее время выполнения
                const avgDuration = stats.task_average_completion_time;
                const avgDurationEl = document.getElementById('stat-avg-duration');
                if (avgDurationEl) {
                    avgDurationEl.textContent = formatDuration(avgDuration);
                }
            }

            // --- Вспомогательные функции ---
            function formatDuration(durationStr) {
                if (!durationStr) return '—';
                
                // Если длительность меньше секунды (миллисекунды, микросекунды, наносекунды)
                if (durationStr.includes('ms') || durationStr.includes('µs') || durationStr.includes('us') || durationStr.includes('ns')) {
                    return '0 сек';
                }
                
                // Регулярное выражение для разбора формата Go duration (например, 1h28m10.729049s, 28m10.729049s, 10.729049s)
                const regex = /^(?:(\d+)h)?(?:(\d+)m)?(?:(\d+(?:\.\d+)?)s)?$/;
                const match = durationStr.match(regex);
                
                if (!match) return durationStr;
                
                const hours = match[1] ? parseInt(match[1]) : 0;
                const minutes = match[2] ? parseInt(match[2]) : 0;
                const secondsFloat = match[3] ? parseFloat(match[3]) : 0;
                
                let seconds = Math.round(secondsFloat);
                let totalMinutes = minutes + hours * 60;
                
                if (seconds >= 60) {
                    totalMinutes += 1;
                    seconds -= 60;
                }
                
                let result = '';
                if (totalMinutes > 0) {
                    result += `${totalMinutes} мин `;
                }
                result += `${seconds} сек`;
                
                return result.trim();
            }

            function escapeHtml(string) {
                const map = {
                    '&': '&amp;',
                    '<': '&lt;',
                    '>': '&gt;',
                    '"': '&quot;',
                    "'": '&#039;'
                };
                return String(string).replace(/[&<>"']/g, function(m) { return map[m]; });
            }

            // Обработчики фильтров статистики
            const btnApplyStatsFilters = document.getElementById('btn-apply-stats-filters');
            const btnResetStatsFilters = document.getElementById('btn-reset-stats-filters');
            const statsUserFilterInput = document.getElementById('stats-user-filter');
            const statsDateFromInput = document.getElementById('stats-date-from');
            const statsDateToInput = document.getElementById('stats-date-to');

            if (btnApplyStatsFilters) {
                btnApplyStatsFilters.addEventListener('click', () => {
                    statsFilterUserId = statsUserFilterInput.value || null;
                    statsFilterFromDate = statsDateFromInput.value || null;
                    statsFilterToDate = statsDateToInput.value || null;
                    loadAppData();
                });
            }

            if (btnResetStatsFilters) {
                btnResetStatsFilters.addEventListener('click', () => {
                    statsUserFilterInput.value = '';
                    statsDateFromInput.value = '';
                    statsDateToInput.value = '';
                    statsFilterUserId = null;
                    statsFilterFromDate = null;
                    statsFilterToDate = null;
                    loadAppData();
                });
            }

            // Обработчик фильтра задач
            const tasksUserFilterSelect = document.getElementById('tasks-user-filter');
            if (tasksUserFilterSelect) {
                tasksUserFilterSelect.addEventListener('change', () => {
                    tasksFilterUserId = tasksUserFilterSelect.value || null;
                    tasksPage = 1; // Сброс на первую страницу при фильтрации
                    invalidateStatsTasksData();
                    loadAppData();
                });
            }

            // Кнопка обновления данных
            const btnRefreshData = document.getElementById('btn-refresh-data');
            if (btnRefreshData) {
                btnRefreshData.addEventListener('click', () => {
                    btnRefreshData.classList.add('is-loading');
                    btnRefreshData.disabled = true;
                    
                    // Запускаем загрузку данных, но анимацию останавливаем ровно через 500 мс
                    invalidateUsersDirectory();
                    invalidateStatsTasksData();
                    loadAppData({ refreshUsersDirectory: true }).catch(console.error);
                    
                    setTimeout(() => {
                        btnRefreshData.classList.remove('is-loading');
                        btnRefreshData.disabled = false;
                    }, 500);
                });
            }

            function renderAll() {
                populateUserSelect();
                renderTasks();
                renderUsers();
                updateStatistics();
            }

            // Первоначальная загрузка данных при запуске приложения
            loadAppData();

            // Периодический опрос API для проверки соединения и автообновления данных (каждые 15 секунд)
            // Автоматически ставится на паузу, когда вкладка в фоне, чтобы не нагружать бэкенд
            let pollInterval = setInterval(loadAppData, 15000);

            document.addEventListener('visibilitychange', () => {
                if (document.hidden) {
                    clearInterval(pollInterval);
                } else {
                    loadAppData();
                    pollInterval = setInterval(loadAppData, 15000);
                }
            });
        });
