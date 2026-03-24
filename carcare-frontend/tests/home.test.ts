import { render, screen } from '@testing-library/svelte/svelte5';
import Home from '../src/routes/+page.svelte';

test('renders welcome message', () => {
  render(Home);
  expect(screen.getByText(/CarCare/i)).toBeTruthy();
});
