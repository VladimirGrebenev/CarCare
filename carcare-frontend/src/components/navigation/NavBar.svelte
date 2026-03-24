<script lang="ts">
  export let active: string = '';
  export let onNavigate: ((path: string) => void) | null = null;
  const nav = [
    { path: '/', label: 'Главная', icon: '🏠' },
    { path: '/cars', label: 'Машины', icon: '🚗' },
    { path: '/fuel', label: 'Топливо', icon: '⛽' },
    { path: '/maintenance', label: 'Сервис', icon: '🛠️' },
    { path: '/fines', label: 'Штрафы', icon: '💸' },
    { path: '/reports', label: 'Отчёты', icon: '📊' },
    { path: '/profile', label: 'Профиль', icon: '👤' }
  ];
  function go(path) {
    if (onNavigate) onNavigate(path);
  }
</script>
<nav class="navbar glassmorphism" aria-label="Основная навигация">
  {#each nav as item}
    <a
      href={item.path}
      class:active={active === item.path}
      aria-current={active === item.path ? 'page' : undefined}
      onclick={(event) => {
        event.preventDefault();
        go(item.path);
      }}
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
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  z-index: 100;
  border-radius: 1.5rem 1.5rem 0 0;
  padding: 0 1rem;
}
.navbar a {
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
.navbar a.active, .navbar a:focus {
  background: rgba(125,226,252,0.18);
  color: var(--accent);
}
.nav-icon {
  font-size: 1.3rem;
}
:global(.dark) .navbar {
  --glass-bg: rgba(30, 30, 40, 0.7);
  --glass-shadow: 0 2px 24px 0 rgba(0,0,0,0.18);
  --text: #fff;
  --accent: #7de2fc;
}
</style>
