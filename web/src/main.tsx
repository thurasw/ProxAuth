import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { NextUIProvider } from '@nextui-org/react';
import { AuthProvider } from './context/AuthProvider.tsx';

(window as any).global = globalThis;

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <NextUIProvider>
            <AuthProvider>
                <App />
            </AuthProvider>
        </NextUIProvider>
    </React.StrictMode>
)
