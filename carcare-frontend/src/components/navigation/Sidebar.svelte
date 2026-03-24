<script lang="ts">
  export let active: string = '';
  export let onNavigate: ((path: string) => void) | null = null;
  type NavItem = { path: string; label: string; icon: string };
  const nav: NavItem[] = [
    { path: '/', label: 'Главная', icon: '🏠' },
    { path: '/cars', label: 'Машины', icon: '🚗' },
    { path: '/fuel', label: 'Топливо', icon: '⛽' },
    { path: '/maintenance', label: 'Сервис', icon: '🛠️' },
    { path: '/fines', label: 'Штрафы', icon: '💸' },
    { path: '/reports', label: 'Отчёты', icon: '📊' },
    { path: '/profile', label: 'Профиль', icon: '👤' }
  ];
  function go(path: string) {
    if (onNavigate) onNavigate(path);
  }
</script>
<aside class="sidebar glassmorphism" aria-label="Боковая навигация">
  {#each nav as item}
    <a
      href={(item as NavItem).path}
      class:active={active === (item as NavItem).path}
      aria-current={active === (item as NavItem).path ? 'page' : undefined}
      onclick={(event: MouseEvent) => {
        event.preventDefault();
        go((item as NavItem).path);
      }}
    >
      <span class="sidebar-icon">{(item as NavItem).icon}</span>
      <span class="sidebar-label">{(item as NavItem).label}</span>
    </a>
  {/each}
</aside>
<style>
.sidebar {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  width: 80px;
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  z-index: 100;
  border-radius: 0 1.5rem 1.5rem 0;
  padding: 2rem 0.5rem;
  align-items: center;
}
.sidebar a {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: var(--text);
  text-decoration: none;
  font-size: 1rem;
  font-weight: 500;
  padding: 0.5rem 0.75rem;
  border-radius: 0.75rem;
  transition: background 0.2s;
}
.sidebar a.active, .sidebar a:focus {
  background: rgba(125,226,252,0.18);
  color: var(--accent);
}
.sidebar-icon {
  font-size: 1.3rem;
}
:global(.dark) .sidebar {
  --glass-bg: rgba(30, 30, 40, 0.7);
  --glass-shadow: 0 2px 24px 0 rgba(0,0,0,0.18);
  --text: #fff;
  --accent: #7de2fc;
}
</style>
