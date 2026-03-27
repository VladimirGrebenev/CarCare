<script lang="ts">
  import { goto } from '$app/navigation';

  type Props = {
    active?: string;
    onNavigate?: ((path: string) => void) | null;
    onToggleTheme?: (() => void) | null;
  };

  let { active = '', onNavigate = null, onToggleTheme = null }: Props = $props();

  const nav = [
    { path: '/',            label: 'Главная',   icon: '⊞' },
    { path: '/cars',        label: 'Машины',    icon: '🚗' },
    { path: '/fuel',        label: 'Топливо',   icon: '⛽' },
    { path: '/maintenance', label: 'Сервис',    icon: '🔧' },
    { path: '/fines',       label: 'Штрафы',    icon: '💸' },
    { path: '/reports',     label: 'Отчёты',    icon: '📊' },
    { path: '/profile',     label: 'Профиль',   icon: '👤' },
  ];

  function go(path: string) {
    if (onNavigate) onNavigate(path);
    else goto(path);
  }
</script>

<nav class="navbar" aria-label="Основная навигация">
  {#each nav as item}
    <a
      href={item.path}
      class="nav-item"
      class:active={active === item.path}
      aria-current={active === item.path ? 'page' : undefined}
      onclick={(e) => { e.preventDefault(); go(item.path); }}
    >
      <span class="nav-icon">{item.icon}</span>
      <span class="nav-label">{item.label}</span>
    </a>
  {/each}
</nav>

<style>
.navbar {
  display: flex;
  justify-content: space-around;
  align-items: center;
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100vw;
  height: 64px;
  background: var(--bg-sidebar);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-top: 1px solid var(--border);
  z-index: 100;
  padding: 0 0.5rem;
  overflow-x: auto;
  gap: 2px;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.6875rem;
  font-weight: 500;
  padding: 0.375rem 0.5rem;
  border-radius: var(--radius-md);
  transition: background var(--transition), color var(--transition);
  min-width: 48px;
  flex-shrink: 0;
}
.nav-item:hover { background: var(--accent-light); color: var(--text-primary); }
.nav-item.active { color: var(--accent-text); }

.nav-icon { font-size: 1.25rem; }
.nav-label { white-space: nowrap; }
</style>
