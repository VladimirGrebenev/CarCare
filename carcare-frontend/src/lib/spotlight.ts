// Spotlight effect — Windows 11 Fluent Design
// Добавляет свечение следующее за курсором на элементах с классом spotlight

export function initSpotlight() {
  if (typeof window === 'undefined') return;

  function handleMouseMove(e: MouseEvent) {
    const target = e.currentTarget as HTMLElement;
    const rect = target.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    target.style.setProperty('--spotlight-x', `${x}px`);
    target.style.setProperty('--spotlight-y', `${y}px`);
    target.style.setProperty('--spotlight-opacity', '1');
  }

  function handleMouseLeave(e: MouseEvent) {
    const target = e.currentTarget as HTMLElement;
    target.style.setProperty('--spotlight-opacity', '0');
  }

  function attachSpotlight() {
    const elements = document.querySelectorAll<HTMLElement>('.spotlight');
    elements.forEach(el => {
      el.removeEventListener('mousemove', handleMouseMove);
      el.removeEventListener('mouseleave', handleMouseLeave);
      el.addEventListener('mousemove', handleMouseMove);
      el.addEventListener('mouseleave', handleMouseLeave);
    });
  }

  // Attach on load and on DOM changes
  attachSpotlight();
  const observer = new MutationObserver(attachSpotlight);
  observer.observe(document.body, { childList: true, subtree: true });
}
