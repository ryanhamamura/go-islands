import React, { useState, useEffect } from 'react';

function UserProfile({ userId }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Simulate fetching user data
    const fetchUser = async () => {
      setLoading(true);
      try {
        // In a real app, this would be an API call
        // await fetch(`/api/users/${userId}`)
        
        // Simulated data
        const mockUser = {
          id: userId,
          name: 'John Doe',
          email: 'john@example.com',
          role: 'Developer',
          joined: '2023-01-01'
        };
        
        setTimeout(() => {
          setUser(mockUser);
          setLoading(false);
        }, 500);
      } catch (error) {
        console.error('Error fetching user:', error);
        setLoading(false);
      }
    };
    
    fetchUser();
  }, [userId]);

  if (loading) return (
    <div className="island-component user-profile loading" aria-busy="true" aria-live="polite">
      <div className="loading-indicator" role="status">
        <span className="sr-only">Loading user profile...</span>
        Loading user profile...
      </div>
    </div>
  );
  
  if (!user) return (
    <div className="island-component user-profile error" role="alert">
      <p className="error-message">User not found</p>
    </div>
  );

  return (
    <div className="island-component user-profile" role="region" aria-labelledby="profile-heading">
      <h3 id="profile-heading">User Profile</h3>
      <div className="profile-card">
        <h4>{user.name}</h4>
        <dl className="profile-details">
          <dt>Email</dt>
          <dd>{user.email}</dd>
          <dt>Role</dt>
          <dd>{user.role}</dd>
          <dt>Joined</dt>
          <dd>{user.joined}</dd>
        </dl>
        <button 
          onClick={() => alert(`Contact ${user.name}`)}
          className="contact-button"
          aria-label={`Contact ${user.name}`}
        >
          Contact
        </button>
      </div>
    </div>
  );
}

// This allows the component to be used standalone
if (import.meta.url.includes('/userProfile.')) {
  import.meta.hot?.accept();
  
  const island = document.querySelector('[data-component="UserProfile"]');
  const props = JSON.parse(island?.dataset.props || '{}');
  
  // Mount directly when loaded as a separate chunk
  if (island) {
    import('react-dom/client').then(({ createRoot }) => {
      createRoot(island).render(<UserProfile {...props} />);
    });
  }
}

export default UserProfile;