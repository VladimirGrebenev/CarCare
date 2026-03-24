<script lang="ts">
  import Card from '../../components/ui/Card.svelte';
  import Input from '../../components/ui/Input.svelte';
  import Button from '../../components/ui/Button.svelte';
  import Toast from '../../components/ui/Toast.svelte';
  import Loader from '../../components/ui/Loader.svelte';
  import ErrorState from '../../components/ui/ErrorState.svelte';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { login, oauthLogin } from '../../lib/api';
  import { bootstrapAuth, setAuth } from '../../stores/auth';

  // State
  let email = $state('');
  let password = $state('');
  let loading = $state(false);
  let error = $state('');
  let success = $state(false);
  let showToast = $state(false);

  // Basic validation
  function validate() {
    if (!email || !/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(email)) {
      error = 'Введите корректный email';
      return false;
    }
    if (!password || password.length < 6) {
      error = 'Пароль должен быть не менее 6 символов';
      return false;
    }
    error = '';
    return true;
  }

  // Submit handler (template for backend integration)

  async function handleLogin(e: Event) {
    e.preventDefault();
    if (!validate()) return;
    loading = true;
    error = '';
    success = false;
    showToast = false;
    try {
      const data = await login(email, password);
      setAuth(data.token, data.user);
      success = true;
      showToast = true;
      setTimeout(() => {
        goto('/profile', { replaceState: true });
      }, 1000);
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Ошибка входа. Попробуйте снова.';
    } finally {
      loading = false;
    }
  }


  function handleOAuth(provider: string) {
    oauthLogin(provider);
  }

  onMount(() => {
    let disposed = false;

    (async () => {
      const isAuthenticated = await bootstrapAuth();
      if (!disposed && isAuthenticated) {
        goto('/profile', { replaceState: true });
      }
    })();

    return () => {
      disposed = true;
    };
  });

  $effect(() => {
    if (showToast && success) {
      setTimeout(() => (showToast = false), 2000);
    }
  });
</script>
<Card className="login-card glass minimal dark-mode" aria-label="Login form">
  <form class="login-form" onsubmit={handleLogin} autocomplete="on" novalidate>
    <h2 class="login-title">Вход</h2>
    <Input
      label="Email"
      type="email"
      name="email"
      bind:value={email}
      required
      autocomplete="email"
      autofocus
      aria-invalid={!!error && !email}
      aria-describedby="email-error"
      class="login-input"
    />
    <Input
      label="Пароль"
      type="password"
      name="password"
      bind:value={password}
      required
      autocomplete="current-password"
      aria-invalid={!!error && !password}
      aria-describedby="password-error"
      class="login-input"
    />
    {#if error}
      <ErrorState id="email-error" message={error} class="login-error" />
    {/if}
    <Button type="submit" class="login-btn" disabled={loading}>
      {#if loading}
        <Loader size={20} />
        Входим...
      {:else}
        Войти
      {/if}
    </Button>
    <div class="oauth-buttons">
      <Button type="button" class="oauth-btn google" onclick={() => handleOAuth('google')} aria-label="Войти через Google">
        <svg aria-hidden="true" width="20" height="20" viewBox="0 0 20 20"><g><circle fill="#fff" cx="10" cy="10" r="10"/><path d="M17.64 10.2c0-.68-.06-1.36-.18-2H10v3.8h4.28c-.18 1-.73 1.85-1.56 2.42v2h2.52c1.48-1.36 2.34-3.36 2.34-5.22z" fill="#4285F4"/><path d="M10 18c2.16 0 3.98-.72 5.3-1.96l-2.52-2c-.7.48-1.6.76-2.78.76-2.14 0-3.96-1.44-4.6-3.38H2.8v2.12C4.12 16.16 6.88 18 10 18z" fill="#34A853"/><path d="M5.4 10.42c-.16-.48-.26-.98-.26-1.5s.1-1.02.26-1.5V5.3H2.8A7.98 7.98 0 0 0 2 10c0 1.28.3 2.5.8 3.58l2.6-2.16z" fill="#FBBC05"/><path d="M10 4.38c1.18 0 2.24.4 3.08 1.18l2.3-2.3C13.98 1.72 12.16 1 10 1 6.88 1 4.12 2.84 2.8 5.3l2.6 2.12C6.04 6.22 7.86 4.38 10 4.38z" fill="#EA4335"/></g></svg>
        Google
      </Button>
      <Button type="button" class="oauth-btn yandex" onclick={() => handleOAuth('yandex')} aria-label="Войти через Яндекс">
        <svg aria-hidden="true" width="20" height="20" viewBox="0 0 20 20"><g><circle fill="#ffcc00" cx="10" cy="10" r="10"/><path d="M10.7 5.5h-1.4l-3 8.9h1.5l.7-2.2h2.9l.7 2.2h1.5l-2.9-8.9zm-1.1 5.7l1.2-3.7 1.2 3.7h-2.4z" fill="#222"/></g></svg>
        Яндекс
      </Button>
    </div>
  </form>
  {#if showToast && success}
    <Toast type="success" message="Успешный вход!" open={showToast} />
  {/if}
</Card>
<style>
/* Glassmorphism, minimalism, dark mode, accessibility */
:global(.login-card) {
  max-width: 370px;
  margin: 4rem auto;
  padding: 2.5rem 2rem 2rem 2rem;
  border-radius: 1.5rem;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.37);
  background: rgba(34, 40, 49, 0.7);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255,255,255,0.18);
  color: #f3f6fa;
}
.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
}
.login-title {
  font-size: 2rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  text-align: center;
}
:global(.login-input):focus-within {
  outline: 2px solid #00bfae;
  outline-offset: 2px;
}
:global(.login-btn) {
  margin-top: 0.5rem;
  background: linear-gradient(90deg, #00bfae 0%, #3a86ff 100%);
  color: #fff;
  font-weight: 600;
  border-radius: 0.75rem;
  min-height: 44px;
  font-size: 1.1rem;
  transition: background 0.2s;
}
:global(.login-btn):focus {
  box-shadow: 0 0 0 3px #00bfae55;
}
.oauth-buttons {
  display: flex;
  gap: 0.75rem;
  margin-top: 0.5rem;
  justify-content: center;
}
:global(.oauth-btn) {
  flex: 1;
  min-width: 0;
  background: #23272f;
  color: #fff;
  border-radius: 0.75rem;
  border: 1px solid #444;
  font-size: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-height: 44px;
  transition: background 0.2s, border 0.2s;
}
:global(.oauth-btn):focus {
  box-shadow: 0 0 0 3px #3a86ff55;
}
:global(.oauth-btn.google) {
  background: #fff;
  color: #222;
  border: 1px solid #e0e0e0;
}
:global(.oauth-btn.yandex) {
  background: #ffcc00;
  color: #222;
  border: 1px solid #ffe066;
}
:global(.login-error) {
  color: #ff6b6b;
  font-size: 1rem;
  margin-top: -0.5rem;
  margin-bottom: 0.5rem;
  text-align: center;
}
:global(.dark-mode) {
  --background: #222831;
  --foreground: #f3f6fa;
  --primary: #00bfae;
  --error: #ff6b6b;
  --card-bg: rgba(34, 40, 49, 0.7);
  --card-border: rgba(255,255,255,0.18);
  --input-bg: #23272f;
  --input-border: #444;
  --input-focus: #00bfae;
  --button-bg: linear-gradient(90deg, #00bfae 0%, #3a86ff 100%);
  --button-focus: #00bfae55;
}
@media (max-width: 480px) {
  :global(.login-card) {
    padding: 1.2rem 0.5rem;
    margin: 2rem 0.25rem;
  }
}
</style>
