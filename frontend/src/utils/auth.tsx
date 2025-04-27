import { redirect } from "@tanstack/react-router";

export const isNetworkError = (error: unknown): boolean => {
    return !window.navigator.onLine && error instanceof Error;
  };

export const login = async (email: string, password: string): Promise<void> => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/login`, {
      method: "POST",
      credentials: "include", // Important for cookies
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });
    if (!response.ok) {
      throw new Error("Login failed: " + response.statusText);
    }
    return;
  } catch (error) {
    throw new Error(
      "Login error: " +
        (error instanceof Error
          ? error.message
          : "An unexpected error occurred")
    );
  }
};

export const logout = async (): Promise<void> => {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_API_URL}/auth/logout`,
      {
        method: "DELETE",
        credentials: "include", // Important for cookies
        headers: {
          "Content-Type": "application/json",
        },
      }
    );
    if (!response.ok) {
      throw new Error("Logout failed: " + response.statusText);
    }
    return;
  } catch (error) {
    throw new Error(
      "Logout error: " +
        (error instanceof Error
          ? error.message
          : "An unexpected error occurred")
    );
  }
};

export const refreshToken = async (): Promise<void> => {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_API_URL}/auth/refresh-token`,
      {
        method: "GET",
        credentials: "include",
      }
    );

    if (!response.ok) {
      throw new Error("Failed to refresh token");
    }
    return;
  } catch (error) {
    throw error;
  }
};

export type Auth = {
  status: "authenticated" | "unauthenticated" | "loading";
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<void>;
};

export const auth: Auth = {
  status: "loading",
  login: async (email, password) => {
    try {
      await login(email, password);
      auth.status = "authenticated";
    } catch (error) {
      auth.status = "unauthenticated";
      throw error;
    }
  },
  logout: async () => {
    try {
      await logout();
      auth.status = "unauthenticated";
      console.log("Logout successful");
      // Redirect to the login page or home page after logout
      redirect({ to: "/login" });
    } catch (error) {
      throw error;
    }
  },
  refresh: async () => {
    try {
      await refreshToken();
      auth.status = "authenticated";
    } catch (error) {
      if (isNetworkError(error)) {
        throw new Error("Network error during token refresh");
      } else {
        auth.status = "unauthenticated";
        throw redirect({ to: "/login" });
      }
    }
  },
};
