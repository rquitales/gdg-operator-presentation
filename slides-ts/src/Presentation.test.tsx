import React from 'react';
import { render, screen } from '@testing-library/react';
import Presentation from './Presentation';

test('renders presentation title', () => {
  render(<Presentation />);
  const headerElement = screen.getByText(/Welcome to Spectacle/i);
  expect(headerElement).toBeInTheDocument();
});
