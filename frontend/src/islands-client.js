/**
 * Island architecture client script
 * 
 * This script is responsible for hydrating React islands
 * embedded within server-rendered HTML pages.
 */
import React from 'react';
import { hydrateRoot } from 'react-dom/client';

// Track which islands are hydrated for monitoring
const hydratedIslands = new Set();

// Performance metrics
const metrics = {
  startTime: performance.now(),
  componentsHydrated: 0,
  errors: 0
};

// Error boundary component for island hydration
class IslandErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }

  componentDidCatch(error, errorInfo) {
    console.error(`Error hydrating island ${this.props.name}:`, error, errorInfo);
    metrics.errors++;
    
    // In production, you might want to log this to an error tracking service
    if (process.env.NODE_ENV === 'production') {
      // sendToErrorTracking(error, errorInfo);
    }
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="island-error" role="alert">
          <p>Something went wrong loading this component.</p>
          {process.env.NODE_ENV !== 'production' && (
            <pre>{this.state.error?.message}</pre>
          )}
        </div>
      );
    }

    return this.props.children;
  }
}

// Dynamically load island components only when needed
async function loadIslandComponent(name) {
  try {
    switch (name) {
      case 'Counter':
        return (await import('./islands/Counter.jsx')).default;
      case 'UserProfile':
        return (await import('./islands/UserProfile.jsx')).default;
      default:
        console.error(`Unknown island component: ${name}`);
        return null;
    }
  } catch (error) {
    console.error(`Failed to load component ${name}:`, error);
    metrics.errors++;
    return null;
  }
}

// Process all islands on the page
document.querySelectorAll('[data-island]').forEach(async (island) => {
  const componentName = island.dataset.component;
  const props = JSON.parse(island.dataset.props || '{}');
  
  // Skip already hydrated islands
  if (hydratedIslands.has(island)) {
    return;
  }
  
  try {
    const Component = await loadIslandComponent(componentName);
    
    if (Component) {
      // Mark as hydrated
      hydratedIslands.add(island);
      
      // Hydrate with error boundary
      hydrateRoot(
        island, 
        <IslandErrorBoundary name={componentName}>
          <Component {...props} />
        </IslandErrorBoundary>
      );
      
      metrics.componentsHydrated++;
    }
  } catch (error) {
    console.error(`Error hydrating component ${componentName}:`, error);
    metrics.errors++;
  }
});

// Report performance metrics
window.addEventListener('load', () => {
  const hydrationTime = performance.now() - metrics.startTime;
  console.debug('Island hydration metrics:', {
    ...metrics,
    hydrationTime: `${hydrationTime.toFixed(2)}ms`
  });
});