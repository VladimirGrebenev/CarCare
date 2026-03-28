import { render } from '@testing-library/svelte/svelte5';
import Home from '../src/routes/+page.svelte';

test('renders dashboard root container', () => {
  const { container } = render(Home);
  expect(container.querySelector('.dashboard-root')).toBeTruthy();
});
