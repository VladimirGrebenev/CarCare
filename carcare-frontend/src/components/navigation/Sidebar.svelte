<script lang="ts">
  import { goto } from '$app/navigation';

  type Props = {
    active?: string;
    onNavigate?: ((path: string) => void) | null;
    onToggleTheme?: (() => void) | null;
  };

  let { active = '', onNavigate = null, onToggleTheme = null }: Props = $props();

  type NavItem = { path: string; label: string; icon: string };
  const nav: NavItem[] = [
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

<aside class="sidebar" aria-label="Боковая навигация">
  <div class="sidebar-brand">
    <span class="brand-icon">🚘</span>
    <span class="brand-name">CarCare</span>
  </div>

  <nav class="sidebar-nav">
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

  <div class="sidebar-bottom">
    <button class="nav-item theme-toggle" onclick={() => onToggleTheme?.()} aria-label="Переключить тему">
      <span class="nav-icon">🌙</span>
      <span class="nav-label">Тема</span>
    </button>
  </div>
</aside>

<style>
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  width: var(--sidebar-width);
  background: var(--bg-sidebar);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  z-index: 100;
  padding: 0.75rem 0.5rem;
  transition: background var(--transition-slow);
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem 1.25rem;
  border-bottom: 1px solid var(--border);
  margin-bottom: 0.5rem;
}
.brand-icon { font-size: 1.5rem; }
.brand-name {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.01em;
}

.sidebar-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 0.625rem 1rem;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9375rem;
  font-weight: 500;
  transition: background var(--transition), color var(--transition);
  cursor: pointer;
  border: none;
  background: transparent;
  width: 100%;
  text-align: left;
}
.nav-item:hover {
  background: var(--accent-light);
  color: var(--text-primary);
}
.nav-item.active {
  background: var(--accent-light);
  color: var(--accent-text);
  font-weight: 600;
}
.nav-item.active .nav-icon {
  filter: none;
}

.nav-icon {
  font-size: 1.125rem;
  width: 1.5rem;
  text-align: center;
  flex-shrink: 0;
}
.nav-label { flex: 1; }

.sidebar-bottom {
  border-top: 1px solid var(--border);
  padding-top: 0.5rem;
  margin-top: 0.5rem;
}

.theme-toggle {
  font-family: var(--font);
}
</style>
