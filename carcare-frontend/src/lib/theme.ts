const STORAGE_KEY = 'carcare-theme';

export type Theme = 'dark' | 'light';

function createTheme() {
  let current: Theme = 'dark';

  function apply(t: Theme) {
    current = t;
    if (typeof document !== 'undefined') {
      if (t === 'light') {
        document.documentElement.classList.add('light');
      } else {
        document.documentElement.classList.remove('light');
      }
    }
  }

  return {
    get current() { return current; },
    toggle() {
      const next: Theme = current === 'dark' ? 'light' : 'dark';
      apply(next);
      try { localStorage.setItem(STORAGE_KEY, next); } catch {}
    },
    init() {
      if (typeof localStorage === 'undefined') return;
      const saved = localStorage.getItem(STORAGE_KEY) as Theme | null;
      const preferred = window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
      apply(saved ?? preferred);
    }
  };
}

export const theme = createTheme();
