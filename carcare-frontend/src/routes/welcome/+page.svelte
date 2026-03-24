<script lang="ts">

  import Button from '../../components/ui/Button.svelte';
  import { goto } from '$app/navigation';
  import { loadAuth, user } from '../../stores/auth';
  import { onMount } from 'svelte';

  let isAuthenticated = false;

  onMount(() => {
    loadAuth();
    const unsub = user.subscribe(u => {
      isAuthenticated = !!u;
      if (isAuthenticated) {
        goto('/profile');
      }
    });
    return unsub;
  });

  function handleLogin() {
    goto('/login');
  }
  function handleRegister() {
    goto('/register');
  }
</script>

<main class="welcome">
  <h1>Добро пожаловать в CarCare!</h1>
  <p>Управляйте своими авто, расходами, заправками, ТО и штрафами в одном месте.</p>
  <div class="welcome-actions">
    <Button onclick={handleLogin}>Войти</Button>
    <Button onclick={handleRegister} variant="secondary">Зарегистрироваться</Button>
  </div>
</main>

<style>
.welcome {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 70vh;
  text-align: center;
}
.welcome-actions {
  margin-top: 2rem;
  display: flex;
  gap: 1.5rem;
  justify-content: center;
}
</style>
