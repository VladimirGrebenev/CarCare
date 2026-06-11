<script lang="ts">
  import { onMount } from 'svelte';
  import { notifyChatAction } from '../../stores/chatEvents';

  interface Message {
    role: 'user' | 'assistant';
    text: string;
  }

  let isOpen = $state(false);
  let messages = $state<Message[]>([]);
  let inputText = $state('');
  let isLoading = $state(false);
  let chatContainer = $state<HTMLDivElement | null>(null);
  let inputRef = $state<HTMLTextAreaElement | null>(null);
  let isAuthenticated = $state(false);

  onMount(() => {
    const token = localStorage.getItem('authToken');
    isAuthenticated = !!token;
  });

  function getAuthHeaders(): Record<string, string> {
    const token = localStorage.getItem('authToken');
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  function focusInput() {
    requestAnimationFrame(() => {
      if (inputRef) inputRef.focus();
    });
  }

  function openChat() {
    if (!isAuthenticated) {
      messages = [{
        role: 'assistant',
        text: '🔒 Пожалуйста, войдите в аккаунт, чтобы использовать чат с AI-помощником.'
      }];
      isOpen = true;
      focusInput();
      return;
    }

    if (messages.length === 0) {
      messages = [{
        role: 'assistant',
        text: '👋 Привет! Я AI-помощник CarCare. Я могу:\n\n• Ответить на вопросы по приложению\n• Помочь создать автомобиль, заправку, ТО или штраф\n• Получить актуальные штрафы с Госуслуг по всем вашим автомобилям\n\nЧем могу помочь?'
      }];
    }
    isOpen = true;
    focusInput();
  }

  async function sendMessage() {
    const text = inputText.trim();
    if (!text || isLoading) return;

    messages = [...messages, { role: 'user', text }];
    inputText = '';
    isLoading = true;

    requestAnimationFrame(() => {
      if (chatContainer) chatContainer.scrollTop = chatContainer.scrollHeight;
    });

    try {
      const history = messages.slice(0, -1).map(m => ({ role: m.role, text: m.text }));

      const res = await fetch('/api/ai/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({ message: text, history })
      });

      if (res.status === 401) {
        messages = [...messages, {
          role: 'assistant',
          text: '🔒 Сессия истекла. Пожалуйста, войдите в аккаунт снова.'
        }];
        return;
      }

      if (!res.ok) throw new Error('Network error');

      const data = await res.json();
      messages = [...messages, { role: 'assistant', text: data.reply }];
      if (data.reply && data.reply.includes('✅')) {
        notifyChatAction();
      }
    } catch {
      messages = [...messages, {
        role: 'assistant',
        text: '😔 Извините, произошла ошибка. Пожалуйста, попробуйте ещё раз позже.'
      }];
    } finally {
      isLoading = false;
      requestAnimationFrame(() => {
        if (chatContainer) chatContainer.scrollTop = chatContainer.scrollHeight;
      });
      focusInput();
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  }

  $effect(() => {
    if (chatContainer && messages.length > 0) {
      requestAnimationFrame(() => {
        chatContainer.scrollTop = chatContainer.scrollHeight;
      });
    }
  });
</script>

<!-- Кнопка открытия чата -->
<button class="chat-fab" onclick={openChat} aria-label="Открыть чат">
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
  </svg>
</button>

<!-- Окно чата -->
{#if isOpen}
  <div class="chat-overlay" onclick={() => isOpen = false}></div>
  <div class="chat-window">
    <!-- Заголовок -->
    <div class="chat-header">
      <div class="chat-header-info">
        <span class="chat-avatar">🤖</span>
        <div>
          <span class="chat-title">AI-помощник CarCare</span>
          <span class="chat-status">Онлайн</span>
        </div>
      </div>
      <button class="chat-close" onclick={() => isOpen = false} aria-label="Закрыть чат">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
    </div>

    <!-- Сообщения -->
    <div class="chat-messages" bind:this={chatContainer}>
      {#each messages as msg}
        <div class="message" class:message-user={msg.role === 'user'} class:message-assistant={msg.role === 'assistant'}>
          {#if msg.role === 'assistant'}
            <span class="msg-avatar">🤖</span>
          {/if}
          <div class="msg-bubble">
            <p class="msg-text">{msg.text}</p>
          </div>
        </div>
      {/each}
      {#if isLoading}
        <div class="message message-assistant">
          <span class="msg-avatar">🤖</span>
          <div class="msg-bubble">
            <div class="typing-indicator">
              <span></span><span></span><span></span>
            </div>
          </div>
        </div>
      {/if}
    </div>

    <!-- Поле ввода -->
    <div class="chat-input-area">
      <textarea
        class="chat-input"
        bind:value={inputText}
        onkeydown={handleKeydown}
        placeholder={isAuthenticated ? "Напишите сообщение..." : "Войдите в аккаунт для доступа к чату"}
        rows="1"
        disabled={!isAuthenticated || isLoading}
        bind:this={inputRef}
      ></textarea>
      <button class="chat-send" onclick={sendMessage} disabled={!inputText.trim() || !isAuthenticated || isLoading} aria-label="Отправить">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/>
        </svg>
      </button>
    </div>
  </div>
{/if}

<style>
  /* Кнопка-фаб */
  .chat-fab {
    position: fixed;
    bottom: 24px;
    right: 24px;
    width: 56px;
    height: 56px;
    border-radius: 50%;
    background: var(--accent, #0078d4);
    color: #fff;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
    z-index: 1000;
    transition: transform 0.2s, box-shadow 0.2s;
  }
  .chat-fab:hover {
    transform: scale(1.05);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.25);
  }

  /* Оверлей */
  .chat-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.3);
    z-index: 1000;
  }

  /* Окно чата */
  .chat-window {
    position: fixed;
    bottom: 90px;
    right: 24px;
    width: 380px;
    height: 560px;
    max-height: calc(100vh - 120px);
    background: var(--bg-base, #1a1a2e);
    border: 1px solid var(--border, rgba(255,255,255,0.1));
    border-radius: 16px;
    display: flex;
    flex-direction: column;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    z-index: 1001;
    overflow: hidden;
    animation: slideUp 0.3s ease-out;
  }

  @keyframes slideUp {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }

  /* Заголовок */
  .chat-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid var(--border, rgba(255,255,255,0.1));
    background: var(--bg-overlay, rgba(255,255,255,0.05));
  }
  .chat-header-info {
    display: flex;
    align-items: center;
    gap: 10px;
  }
  .chat-avatar {
    font-size: 1.5rem;
  }
  .chat-title {
    font-weight: 600;
    font-size: 0.9375rem;
    color: var(--text-primary, #fff);
    display: block;
  }
  .chat-status {
    font-size: 0.75rem;
    color: #4caf50;
  }
  .chat-close {
    background: none;
    border: none;
    color: var(--text-secondary, #999);
    cursor: pointer;
    padding: 4px;
    border-radius: 8px;
    transition: background 0.2s;
  }
  .chat-close:hover {
    background: rgba(255,255,255,0.1);
  }

  /* Сообщения */
  .chat-messages {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .message {
    display: flex;
    gap: 8px;
    max-width: 85%;
  }
  .message-user {
    align-self: flex-end;
    flex-direction: row-reverse;
  }
  .message-assistant {
    align-self: flex-start;
  }
  .msg-avatar {
    font-size: 1.25rem;
    flex-shrink: 0;
    margin-top: 4px;
  }
  .msg-bubble {
    padding: 10px 14px;
    border-radius: 12px;
    font-size: 0.875rem;
    line-height: 1.5;
    color: var(--text-primary, #fff);
  }
  .message-user .msg-bubble {
    background: var(--accent, #0078d4);
    border-bottom-right-radius: 4px;
  }
  .message-assistant .msg-bubble {
    background: var(--bg-overlay, rgba(255,255,255,0.08));
    border-bottom-left-radius: 4px;
  }
  .msg-text {
    margin: 0;
    white-space: pre-wrap;
    word-break: break-word;
  }

  /* Индикатор печатания */
  .typing-indicator {
    display: flex;
    gap: 4px;
    padding: 4px 0;
  }
  .typing-indicator span {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--text-secondary, #999);
    animation: typing 1.4s infinite ease-in-out;
  }
  .typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
  .typing-indicator span:nth-child(3) { animation-delay: 0.4s; }
  @keyframes typing {
    0%, 60%, 100% { opacity: 0.3; transform: scale(0.8); }
    30% { opacity: 1; transform: scale(1); }
  }

  /* Поле ввода */
  .chat-input-area {
    display: flex;
    align-items: flex-end;
    gap: 8px;
    padding: 12px 16px;
    border-top: 1px solid var(--border, rgba(255,255,255,0.1));
    background: var(--bg-overlay, rgba(255,255,255,0.03));
  }
  .chat-input {
    flex: 1;
    background: var(--bg-overlay, rgba(255,255,255,0.05));
    border: 1px solid var(--border, rgba(255,255,255,0.1));
    border-radius: 12px;
    padding: 10px 14px;
    color: var(--text-primary, #fff);
    font-size: 0.875rem;
    font-family: inherit;
    resize: none;
    outline: none;
    max-height: 120px;
    transition: border-color 0.2s;
  }
  .chat-input:focus {
    border-color: var(--accent, #0078d4);
  }
  .chat-input::placeholder {
    color: var(--text-secondary, #666);
  }
  .chat-send {
    background: var(--accent, #0078d4);
    color: #fff;
    border: none;
    border-radius: 12px;
    padding: 10px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.2s, opacity 0.2s;
    flex-shrink: 0;
  }
  .chat-send:hover:not(:disabled) {
    background: var(--accent-hover, #005a9e);
  }
  .chat-send:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  /* Адаптивность для мобильных */
  @media (max-width: 480px) {
    .chat-window {
      right: 8px;
      left: 8px;
      bottom: 80px;
      width: auto;
      height: calc(100vh - 100px);
    }
    .chat-fab {
      bottom: 80px;
      right: 16px;
    }
  }
</style>