
export const login = async (email: string, password: string):Promise<void>  => {
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/login`, {
        method: 'POST',
        credentials: 'include', // Important for cookies
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });
      if (!response.ok) {
        throw new Error('Login failed: ' + response.statusText);
      }
      return ;
    } catch (error) {
      throw new Error('Login error: ' + (error instanceof Error ? error.message : 'An unexpected error occurred'));
    }
  };

export const logout = async (): Promise<void> => {
    try {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/logout`, {
            method: 'POST',
            credentials: 'include', // Important for cookies
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (!response.ok) {
            throw new Error('Logout failed: ' + response.statusText);
        }

        window.location.href = '/'; // Redirect to the home page or login page after logout
        return ;
    } catch (error) {
        throw new Error('Logout error: ' + (error instanceof Error ? error.message : 'An unexpected error occurred'));
    }
}

export type Auth = {
    status : 'authenticated' | 'unauthenticated' | 'loading';
    login: (email: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
}

export const auth: Auth = {
    status: 'loading',
    login: async (email, password) => {
        try {
            await login(email, password);
            auth.status = 'authenticated';
        } catch (error) {
            auth.status = 'unauthenticated';
            throw error;
        }
    },
    logout: async () => {
        try {
            await logout();
            auth.status = 'unauthenticated';
        } catch (error) {
            throw error;
        }
    }
}
