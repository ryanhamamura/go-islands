// This is the main entry point for the client-side hydration
import React from 'react';
import { hydrateRoot } from 'react-dom/client';

// Import island components
import Counter from './islands/Counter';
import UserProfile from './islands/UserProfile';

// Component registry
const components = {
  Counter,
  UserProfile
};

// Find and hydrate all island components on the page
document.querySelectorAll('[data-island]').forEach(island => {
  const componentName = island.dataset.component;
  const props = JSON.parse(island.dataset.props || '{}');
  
  if (components[componentName]) {
    hydrateRoot(island, React.createElement(components[componentName], props));
  }
});
