import { logout } from "./auth";

interface FetchOptions extends RequestInit{
    skipAuth? : boolean
}

const isNetworkError = (error: unknown): boolean => {
  return !window.navigator.onLine || error instanceof Error;
};

const refreshToken = async () => {
    try {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/refresh`, {
            method: 'GET',
            credentials: 'include',

        });

        if (!response.ok) {
            throw new Error('Failed to refresh token');
        }

        
        return true;
    } catch (error) {
        return error
    }
};

export const fetchWithAuth = async (url: string, options: FetchOptions = {}): Promise<Response> => {
    const { skipAuth=false, ...fetchOptions } = options;
    const finalOptions: RequestInit = {
        ...fetchOptions,
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
            ...fetchOptions.headers,
        },
    };

    try {
        const response = await fetch(url, finalOptions);

        if (response.status === 401 && !skipAuth) {
            const refreshResponse = await refreshToken();

            if (refreshResponse instanceof Error) {
                if (isNetworkError(refreshResponse)) {
                    throw new Error('Network error during token refresh');
                } else {
                    await logout();
                    throw new Error('Authentication failed');
                }
            }

            // Retry the original request with the new token
            return fetch(url, finalOptions);
        }

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return response;
    }
    catch (error) {
        if (error instanceof Error) {
            if (error.message === 'Authentication failed') {
              throw error;
            } else if (error.message === 'Network error during token refresh') {
              throw error;
            }
          }
          throw new Error('Network request failed');
    }
}

export const http = {
    get: async (url: string, options: FetchOptions = {}) => {
        const response = await fetchWithAuth(url, { method: 'GET', ...options });
        return response.json();
    },
    post: async (url: string, data: any, options: FetchOptions = {}) => {
        const response = await fetchWithAuth(url, { method: 'POST', body: JSON.stringify(data), ...options });
        return response.json();
    },
    put: async (url: string, data: any, options: FetchOptions = {}) => {
        const response = await fetchWithAuth(url, { method: 'PUT', body: JSON.stringify(data), ...options });
        return response.json();
    },
    delete: async (url: string, options: FetchOptions = {}) => {
        const response = await fetchWithAuth(url, { method: 'DELETE', ...options });
        return response.json();
    },
    patch: async (url: string, data: any, options: FetchOptions = {}) => {
        const response = await fetchWithAuth(url, { method: 'PATCH', body: JSON.stringify(data), ...options });
        return response.json();
    }
}