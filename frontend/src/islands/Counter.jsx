import React, { useState, useCallback } from 'react';

function Counter({ initialCount = 0 }) {
  const [count, setCount] = useState(initialCount);
  
  // Memoized handlers
  const handleDecrement = useCallback(() => setCount(prev => prev - 1), []);
  const handleIncrement = useCallback(() => setCount(prev => prev + 1), []);

  return (
    <div className="island-component counter" role="region" aria-label="Counter">
      <h3 id="counter-heading">Interactive Counter</h3>
      <p aria-live="polite" aria-atomic="true">Count: <span className="count-value">{count}</span></p>
      <div className="button-group">
        <button 
          onClick={handleDecrement}
          aria-label="Decrement counter"
          className="counter-button"
        >
          Decrement
        </button>
        <button 
          onClick={handleIncrement}
          aria-label="Increment counter"
          className="counter-button"
        >
          Increment
        </button>
      </div>
    </div>
  );
}

// This allows the component to be used standalone
if (import.meta.url.includes('/counter.')) {
  import.meta.hot?.accept();
  
  const island = document.querySelector('[data-component="Counter"]');
  const props = JSON.parse(island?.dataset.props || '{}');
  
  // Mount directly when loaded as a separate chunk
  if (island) {
    import('react-dom/client').then(({ createRoot }) => {
      createRoot(island).render(<Counter {...props} />);
    });
  }
}

export default Counter;
