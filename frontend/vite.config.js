import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  const isProduction = mode === 'production';
  
  return {
    plugins: [
      react()
    ],
    build: {
      outDir: '../backend/static',
      emptyOutDir: true,
      manifest: true,
      minify: isProduction,
      sourcemap: !isProduction,
      rollupOptions: {
        input: {
          main: './src/main.jsx',
          'islands-client': './src/islands-client.js',
          counter: './src/islands/Counter.jsx',
          userProfile: './src/islands/UserProfile.jsx'
        },
        output: {
          manualChunks: (id) => {
            // Create a vendors chunk for node_modules
            if (id.includes('node_modules')) {
              return 'vendors';
            }
          }
        }
      },
      // Optimize chunks
      chunkSizeWarningLimit: 600,
      cssCodeSplit: true
    },
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true
        }
      },
      // For better dev experience
      hmr: {
        overlay: true
      }
    },
    // Improve performance
    optimizeDeps: {
      include: ['react', 'react-dom']
    }
  };
});