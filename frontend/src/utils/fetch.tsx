import { isNetworkError, logout,refreshToken } from "./auth";

interface FetchOptions extends RequestInit{
    skipAuth? : boolean
}

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
            try {
                await refreshToken();
                // Retry the original request with the new token
                return fetch(url, finalOptions);
            } catch (refreshError) {
                if (isNetworkError(refreshError)) {
                    throw new Error('Network error during token refresh');
                } else {
                    await logout();
                    window.location.href = '/'; // Redirect to login page
                    throw new Error('Authentication failed');
                }
            }
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
        if (!response.ok) {
            console.error('Error fetching data:', response.statusText);
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to fetch data');
        }
        return response.json();
    },
    post: async (url: string, data: any, options: FetchOptions = {}) => {
        const response = await fetchWithAuth(url, { method: 'POST', body: JSON.stringify(data), ...options });
        if (!response.ok) {
            console.error('Error fetching data:', response.statusText);
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to fetch data');
        }
        return response.json();
    },
    put: async (url: string, data: any, options: FetchOptions = {}) => {
        const response = await fetchWithAuth(url, { method: 'PUT', body: JSON.stringify(data), ...options });
        if (!response.ok) {
            console.error('Error fetching data:', response.statusText);
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to fetch data');
        }
        return response.json();
    },
    delete: async (url: string, options: FetchOptions = {}) => {
        const response = await fetchWithAuth(url, { method: 'DELETE', ...options });
        if (!response.ok) {
            console.error('Error fetching data:', response.statusText);
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to fetch data');
        }
        return response.json();
    },
    patch: async (url: string, data: any, options: FetchOptions = {}) => {
        const response = await fetchWithAuth(url, { method: 'PATCH', body: JSON.stringify(data), ...options });
        if (!response.ok) {
            console.error('Error fetching data:', response.statusText);
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to fetch data');
        }
        return response.json();
    }
}