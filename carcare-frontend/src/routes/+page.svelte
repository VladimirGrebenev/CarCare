<!-- src/routes/+page.svelte -->

<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { authToken, bootstrapAuth } from '../stores/auth';
  import { fetchFuelHistory } from '../lib/api';
  import Card from '../components/ui/Card.svelte';

  // --- State ---
  let isAuthenticated = $state(false);
  let loading = $state(true);

  // Cost stats
  let totalFuel = $state<number | null>(null);
  let totalMaintenance = $state<number | null>(null);
  let totalFines = $state<number | null>(null);

  // Recent fuel events
  type FuelEvent = {
    id: string;
    date: string;
    volume?: number;
    liters?: number;
    cost?: number;
    total_cost?: number;
  };
  let recentFuel = $state<FuelEvent[]>([]);

  // --- Helpers ---
  function formatRub(value: number | null): string {
    if (value === null) return '—';
    return value.toLocaleString('ru-RU') + ' ₽';
  }

  function formatDate(dateStr: string): string {
    try {
      return new Date(dateStr).toLocaleDateString('ru-RU');
    } catch {
      return dateStr;
    }
  }

  function getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('authToken');
  }

  function authHeaders(): Record<string, string> {
    const token = getToken();
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  async function fetchReports(): Promise<void> {
    try {
      const res = await fetch('/api/reports', {
        credentials: 'include',
        headers: authHeaders()
      });
      if (!res.ok) return;
      const data = await res.json().catch(() => null);
      if (!data) return;
      totalFuel = typeof data.total_fuel_cost === 'number' ? data.total_fuel_cost : null;
      totalMaintenance = typeof data.total_maintenance_cost === 'number' ? data.total_maintenance_cost : null;
      totalFines = typeof data.total_fines_amount === 'number' ? data.total_fines_amount : null;
    } catch {
      // Graceful degradation — leave values as null
    }
  }

  async function loadDashboard(): Promise<void> {
    try {
      const [fuelRaw] = await Promise.all([
        fetchFuelHistory().catch(() => [] as FuelEvent[]),
        fetchReports()
      ]);

      const fuelArr = Array.isArray(fuelRaw) ? (fuelRaw as FuelEvent[]) : [];
      recentFuel = fuelArr
        .slice()
        .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
        .slice(0, 3);
    } catch {
      recentFuel = [];
    }
  }

  // --- Lifecycle ---
  onMount(() => {
    let disposed = false;

    (async () => {
      const auth = await bootstrapAuth();
      if (disposed) return;

      isAuthenticated = auth;
      loading = false;

      if (auth) {
        await loadDashboard();
      }
    })();

    return () => {
      disposed = true;
    };
  });

  // --- Derived ---
  const totalAll = $derived(
    totalFuel !== null || totalMaintenance !== null || totalFines !== null
      ? (totalFuel ?? 0) + (totalMaintenance ?? 0) + (totalFines ?? 0)
      : null
  );

  const quickLinks = [
    { label: 'Мои автомобили', icon: '🚗', href: '/cars' },
    { label: 'Заправки', icon: '⛽', href: '/fuel' },
    { label: 'Техобслуживание', icon: '🔧', href: '/maintenance' },
    { label: 'Штрафы', icon: '📋', href: '/fines' },
    { label: 'Отчёты', icon: '📊', href: '/reports' }
  ];
</script>

{#if loading}
  <main class="dashboard-root">
    <div class="loader-wrap">
      <div class="spinner"></div>
    </div>
  </main>
{:else if !isAuthenticated}
  <main class="dashboard-root">
    <section class="hero glass">
      <h1 class="hero-title">CarCare</h1>
      <p class="hero-sub">Управляйте расходами на автомобиль — топливо, ТО, штрафы.</p>
      <button class="btn-accent" onclick={() => goto('/welcome')}>Войти</button>
    </section>
  </main>
{:else}
  <main class="dashboard-root">
    <h1 class="page-title">Обзор</h1>

    <!-- Row 1: Cost stats -->
    <section class="section">
      <h2 class="section-title">Расходы</h2>
      <div class="grid-4">
        <div class="stat-card glass">
          <span class="stat-label">Топливо</span>
          <span class="stat-value">{formatRub(totalFuel)}</span>
        </div>
        <div class="stat-card glass">
          <span class="stat-label">ТО</span>
          <span class="stat-value">{formatRub(totalMaintenance)}</span>
        </div>
        <div class="stat-card glass">
          <span class="stat-label">Штрафы</span>
          <span class="stat-value">{formatRub(totalFines)}</span>
        </div>
        <div class="stat-card glass accent-card">
          <span class="stat-label">Всего</span>
          <span class="stat-value accent-value">{formatRub(totalAll)}</span>
        </div>
      </div>
    </section>

    <!-- Row 2: Quick links -->
    <section class="section">
      <h2 class="section-title">Быстрый доступ</h2>
      <div class="grid-links">
        {#each quickLinks as link}
          <a href={link.href} class="link-card glass">
            <span class="link-icon">{link.icon}</span>
            <span class="link-label">{link.label}</span>
          </a>
        {/each}
      </div>
    </section>

    <!-- Row 3: Recent fuel events -->
    <section class="section">
      <h2 class="section-title">Последние заправки</h2>
      {#if recentFuel.length === 0}
        <div class="empty-state glass">
          <span>Нет записей о заправках</span>
          <a href="/fuel" class="link-inline">Добавить заправку</a>
        </div>
      {:else}
        <Card>
          {#snippet children()}
            <table class="events-table">
              <thead>
                <tr>
                  <th>Дата</th>
                  <th>Объём</th>
                  <th>Сумма</th>
                </tr>
              </thead>
              <tbody>
                {#each recentFuel as event}
                  <tr>
                    <td>{formatDate(event.date)}</td>
                    <td>{event.volume ?? event.liters ?? '—'} л</td>
                    <td>{formatRub(event.cost ?? event.total_cost ?? null)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/snippet}
        </Card>
      {/if}
    </section>
  </main>
{/if}

<style>
.dashboard-root {
  max-width: 1100px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Loading */
.loader-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
}
.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Hero (unauthenticated) */
.hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.25rem;
  padding: 4rem 2rem;
  border-radius: var(--radius-xl);
  text-align: center;
}
.hero-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.03em;
}
.hero-sub {
  font-size: 1.1rem;
  color: var(--text-secondary);
  max-width: 480px;
}
.btn-accent {
  padding: 0.65rem 2rem;
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition);
  font-family: var(--font);
}
.btn-accent:hover { background: var(--accent-hover); }

/* Page heading */
.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}

/* Sections */
.section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.section-title {
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-secondary);
}

/* Glassmorphism base */
.glass {
  background: var(--bg-overlay);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-md);
}

/* Stat cards grid */
.grid-4 {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}
.stat-card {
  border-radius: var(--radius-lg);
  padding: 1.25rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  transition: transform var(--transition), box-shadow var(--transition);
}
.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}
.stat-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}
.accent-card {
  border-color: rgba(0, 120, 212, 0.3);
}
.accent-value {
  color: var(--accent-text);
}

/* Quick links grid */
.grid-links {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 16px;
}
.link-card {
  border-radius: var(--radius-lg);
  padding: 1.25rem 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.625rem;
  text-decoration: none;
  color: var(--text-primary);
  transition: transform var(--transition), box-shadow var(--transition), background var(--transition);
}
.link-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-lg);
  background: var(--accent-light);
  border-color: rgba(0, 120, 212, 0.25);
}
.link-icon {
  font-size: 2rem;
  line-height: 1;
}
.link-label {
  font-size: 0.875rem;
  font-weight: 600;
  text-align: center;
  color: var(--text-primary);
}

/* Empty state */
.empty-state {
  border-radius: var(--radius-lg);
  padding: 2rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  color: var(--text-secondary);
  font-size: 0.9375rem;
}
.link-inline {
  color: var(--accent-text);
  font-weight: 500;
  text-decoration: none;
}
.link-inline:hover { text-decoration: underline; }

/* Events table */
.events-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9375rem;
}
.events-table th {
  text-align: left;
  padding: 0.5rem 0.75rem;
  color: var(--text-secondary);
  font-size: 0.8125rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-bottom: 1px solid var(--border);
}
.events-table td {
  padding: 0.75rem;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border);
}
.events-table tr:last-child td {
  border-bottom: none;
}
.events-table tr:hover td {
  background: var(--accent-light);
}
</style>
