<script lang="ts">

  import Card from '../../components/ui/Card.svelte';
  import Input from '../../components/ui/Input.svelte';
  import Button from '../../components/ui/Button.svelte';
  import Toast from '../../components/ui/Toast.svelte';
  import Loader from '../../components/ui/Loader.svelte';
  import ErrorState from '../../components/ui/ErrorState.svelte';
  import { register } from '../../lib/api';
  import { setAuth } from '../../stores/auth';
  import { onMount } from 'svelte';

  let email = '';
  let password = '';
  let confirmPassword = '';
  let loading = false;
  let error = '';
  let success = false;
  let showToast = false;

  function validate() {
    if (!email || !/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(email)) {
      error = 'Введите корректный email';
      return false;
    }
    if (!password || password.length < 6) {
      error = 'Пароль должен быть не менее 6 символов';
      return false;
    }
    if (password !== confirmPassword) {
      error = 'Пароли не совпадают';
      return false;
    }
    error = '';
    return true;
  }

  async function handleRegister(e: Event) {
    e.preventDefault();
    if (!validate()) return;
    loading = true;
    error = '';
    success = false;
    showToast = false;
    try {
      const data = await register(email, password);
      setAuth(data.token, data.user);
      success = true;
      showToast = true;
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Ошибка регистрации. Попробуйте снова.';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    if (success) {
      setTimeout(() => {
        window.location.href = '/profile';
      }, 1500);
    }
  });
</script>

<Card>
  <form onsubmit={handleRegister} aria-label="Регистрация" autocomplete="off">
    <Input label="Email" type="email" bind:value={email} required autocomplete="email" />
    <Input label="Пароль" type="password" bind:value={password} required autocomplete="new-password" minLength={6} />
    <Input label="Повторите пароль" type="password" bind:value={confirmPassword} required autocomplete="new-password" minLength={6} />
    {#if error}
      <ErrorState message={error} />
    {/if}
    <Button type="submit" disabled={loading}>Зарегистрироваться</Button>
    {#if loading}
      <Loader />
    {/if}
  </form>
  <Toast open={showToast} message="Проверьте email для подтверждения" />
</Card>

<style>
form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
  max-width: 340px;
  margin: 0 auto;
}
</style>
