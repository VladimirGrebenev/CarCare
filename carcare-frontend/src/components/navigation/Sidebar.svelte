<script lang="ts">
  import { goto } from '$app/navigation';

  type Props = {
    active?: string;
    onNavigate?: ((path: string) => void) | null;
    onToggleTheme?: (() => void) | null;
    onCollapse?: ((collapsed: boolean) => void) | null;
  };

  let { active = '', onNavigate = null, onToggleTheme = null, onCollapse = null }: Props = $props();

  const STORAGE_KEY = 'sidebar-collapsed';

  function readCollapsed(): boolean {
    if (typeof localStorage === 'undefined') return false;
    return localStorage.getItem(STORAGE_KEY) === 'true';
  }

  let collapsed = $state(readCollapsed());

  function toggleCollapsed() {
    collapsed = !collapsed;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(STORAGE_KEY, String(collapsed));
    }
    onCollapse?.(collapsed);
  }

  type NavItem = { path: string; label: string; icon: string };
  const nav: NavItem[] = [
    {
      path: '/',
      label: 'Главная',
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>`
    },
    {
      path: '/cars',
      label: 'Машины',
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M19 17H5a2 2 0 0 1-2-2V9a2 2 0 0 1 .85-1.61L7 5h10l3.15 2.39A2 2 0 0 1 21 9v6a2 2 0 0 1-2 2z"/><circle cx="7.5" cy="17" r="1.5"/><circle cx="16.5" cy="17" r="1.5"/></svg>`
    },
    {
      path: '/fuel',
      label: 'Топливо',
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="7" width="12" height="14" rx="2"/><path d="M14 9h2a2 2 0 0 1 2 2v4a1 1 0 0 0 2 0v-5l-3-4"/><path d="M6 7V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v3"/><line x1="6" y1="12" x2="10" y2="12"/></svg>`
    },
    {
      path: '/maintenance',
      label: 'Сервис',
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>`
    },
    {
      path: '/fines',
      label: 'Штрафы',
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`
    },
    {
      path: '/reports',
      label: 'Отчёты',
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>`
    },
    {
      path: '/profile',
      label: 'Профиль',
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>`
    },
  ];


  const themeIcon = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>`;

  /* Toggle button double-chevron icons */
  const chevronsLeftIcon = `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="11 17 6 12 11 7"/><polyline points="18 17 13 12 18 7"/></svg>`;
  const chevronsRightIcon = `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="13 17 18 12 13 7"/><polyline points="6 17 11 12 6 7"/></svg>`;

  function go(path: string) {
    if (onNavigate) onNavigate(path);
    else goto(path);
  }
</script>

<!-- eslint-disable svelte/no-at-html-tags -->
<aside class="sidebar" class:collapsed aria-label="Боковая навигация">
  <!-- Brand row -->
  <div class="sidebar-brand">
    <span class="brand-icon">
      <svg class="logo-svg" width="28" height="28" viewBox="0 0 280 260" xmlns="http://www.w3.org/2000/svg">
        <circle class="logo-outer" cx="132" cy="130" r="90" fill="none" stroke="currentColor" stroke-width="26"
          stroke-dasharray="503 62" stroke-dashoffset="62" stroke-linecap="butt"/>
        <circle class="logo-inner" cx="132" cy="130" r="54" fill="none" stroke="currentColor" stroke-width="24"
          stroke-dasharray="322 42" stroke-dashoffset="40" stroke-linecap="butt"/>
      </svg>
    </span>
    {#if !collapsed}
      <span class="brand-name"><span class="brand-car">Car</span><span class="brand-care">Care</span></span>
    {/if}
  </div>

  <nav class="sidebar-nav">
    {#each nav as item}
      <a
        href={item.path}
        class="nav-item spotlight"
        class:active={active === item.path}
        aria-current={active === item.path ? 'page' : undefined}
        title={item.label}
        onclick={(e) => { e.preventDefault(); go(item.path); }}
        data-tooltip={item.label}
      >
        <span class="nav-icon">{@html item.icon}</span>
        {#if !collapsed}
          <span class="nav-label">{item.label}</span>
        {/if}
      </a>
    {/each}
  </nav>

  <div class="sidebar-bottom">
    <button
      class="nav-item spotlight theme-toggle"
      onclick={() => onToggleTheme?.()}
      aria-label="Переключить тему"
      title="Тема"
      data-tooltip="Тема"
    >
      <span class="nav-icon">{@html themeIcon}</span>
      {#if !collapsed}
        <span class="nav-label">Тема</span>
      {/if}
    </button>
    <button
      class="nav-item spotlight collapse-toggle"
      onclick={toggleCollapsed}
      aria-label={collapsed ? 'Полный режим' : 'Компактный режим'}
      data-tooltip={collapsed ? 'Полный режим' : 'Компактный режим'}
    >
      <span class="nav-icon">{@html collapsed ? chevronsRightIcon : chevronsLeftIcon}</span>
      {#if !collapsed}
        <span class="nav-label">Компактный режим</span>
      {/if}
    </button>
  </div>
</aside>

<style>
/* ── Sidebar shell ── */
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
  /* smooth width transition */
  transition:
    width 220ms cubic-bezier(0.4, 0, 0.2, 1),
    background var(--transition-slow);
  overflow: hidden;
}

.sidebar.collapsed {
  width: var(--sidebar-width-collapsed);
}

/* ── Brand ── */
.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem 1.25rem;
  border-bottom: 1px solid var(--border);
  margin-bottom: 0.5rem;
  overflow: hidden;
  white-space: nowrap;
}

.collapsed .sidebar-brand {
  padding-left: 0;
  padding-right: 0;
  justify-content: center;
  gap: 0;
}

.brand-icon {
  display: flex;
  align-items: center;
  color: var(--accent-text);
  flex-shrink: 0;
  cursor: pointer;
}

/* Анимация логотипа при наведении */
.logo-svg .logo-outer {
  transform-origin: 132px 130px;
  transition: transform 600ms cubic-bezier(0.4, 0, 0.2, 1);
}
.logo-svg .logo-inner {
  transform-origin: 132px 130px;
  transition: transform 600ms cubic-bezier(0.4, 0, 0.2, 1) 80ms;
}
.sidebar-brand:hover .logo-outer {
  transform: rotate(90deg);
}
.sidebar-brand:hover .logo-inner {
  transform: rotate(-90deg);
}

.brand-name {
  font-size: 1.125rem;
  color: var(--text-primary);
  letter-spacing: -0.01em;
  transition: opacity 150ms ease, width 150ms ease;
}
.brand-car { font-weight: 700; }
.brand-care { font-weight: 300; }

/* ── Nav ── */
.sidebar-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow-y: auto;
  overflow-x: hidden;
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
  transition:
    background var(--transition),
    color var(--transition),
    padding var(--transition);
  cursor: pointer;
  border: none;
  background: transparent;
  width: 100%;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  position: relative;
}

.collapsed .nav-item {
  padding: 0.625rem;
  justify-content: center;
  gap: 0;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-primary);
}
.nav-item.active {
  background: var(--accent-light);
  color: var(--accent-text);
  font-weight: 600;
}

/* SVG icon transitions */
.nav-icon :global(svg) {
  display: block;
  transition: stroke var(--transition), color var(--transition);
}
.nav-item.active .nav-icon :global(svg) {
  stroke: var(--accent-text);
}
.nav-item:hover:not(.active) .nav-icon :global(svg) {
  stroke: var(--text-primary);
  opacity: 1;
}

.nav-icon {
  display: flex;
  align-items: center;
  width: 1.25rem;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.nav-label {
  flex: 1;
  position: relative;
  z-index: 1;
}

/* ── Tooltip shown in collapsed state ── */
.collapsed .nav-item[data-tooltip] {
  overflow: visible;
}

.collapsed .nav-item[data-tooltip]::after {
  content: attr(data-tooltip);
  position: absolute;
  left: calc(100% + 12px);
  top: 50%;
  transform: translateY(-50%);
  background: var(--bg-layer);
  color: var(--text-primary);
  font-size: 0.8125rem;
  font-weight: 500;
  padding: 0.3rem 0.6rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-md);
  white-space: nowrap;
  pointer-events: none;
  opacity: 0;
  transition: opacity 120ms ease;
  z-index: 200;
}

.collapsed .nav-item[data-tooltip]:hover::after {
  opacity: 1;
}

/* ── Bottom ── */
.sidebar-bottom {
  border-top: 1px solid var(--border);
  padding-top: 0.5rem;
  margin-top: 0.5rem;
}

.theme-toggle {
  font-family: var(--font);
}
</style>
