import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface User {
    id: string;
    username: string;
    email: string;
    role: string;
}

export interface AuthState {
    token: string | null;
    user: User | null;
    isAuthenticated: boolean;
}

const initialState: AuthState = {
    token: null,
    user: null,
    isAuthenticated: false
};

// Load from localStorage if available
const storedAuth = browser ? localStorage.getItem('auth') : null;
const initialAuth: AuthState = storedAuth ? JSON.parse(storedAuth) : initialState;

export const auth = writable<AuthState>(initialAuth);

// Subscribe to changes and update localStorage
if (browser) {
    auth.subscribe((value) => {
        if (value.isAuthenticated) {
            localStorage.setItem('auth', JSON.stringify(value));
        } else {
            localStorage.removeItem('auth');
        }
    });
}

export const login = (token: string, user: User) => {
    auth.set({
        token,
        user,
        isAuthenticated: true
    });
};

export const logout = () => {
    auth.set(initialState);
};
