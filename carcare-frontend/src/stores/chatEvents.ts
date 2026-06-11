import { writable } from 'svelte/store';

// Инкрементируется каждый раз когда чат успешно создаёт/редактирует/удаляет сущность.
// Страницы подписываются и обновляют свои списки при изменении.
export const chatActionCount = writable(0);

export function notifyChatAction() {
  chatActionCount.update(n => n + 1);
}
