import { createFileRoute, useNavigate } from '@tanstack/react-router';
import React, { useState } from 'react';
import { auth } from '../utils/auth'; // Import the auth object
import { http } from '@/utils/fetch';

export const Route = createFileRoute('/login')({
  component: RouteComponent,
});

function RouteComponent() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate({ from: '/login' });

  //   const categories = async ():Promise<void>  => {
  //   try {
      
  //     const response = await fetch(`${import.meta.env.VITE_API_URL}/categories/`, {
  //       method: 'POST',
  //       credentials: 'include', // Important for cookies
  //       headers: {
  //         'Content-Type': 'application/json',
  //       },
  //       body: JSON.stringify({ name: "Groceries10" }),
  //     });
      
  //     if (!response.ok) {
  //       console.log('Categories API error:', response.status, response.statusText);
  //     }
  //     return;
  //   } catch (error) {
  //     console.log('Error fetching categories: ', error);
  //     throw new Error('API error: ' + (error instanceof Error ? error.message : 'An unexpected error occurred'));
  //   }
  // };

  const categories = async ():Promise<void>  => {
    try {
      const response = await http.get(`${import.meta.env.VITE_API_URL}/categories/`);
        if (!response.ok) {
        console.log('Categories API error:', response.status, response.statusText);
        }
      return;
    } catch (error) {
      console.log('Error fetching categories: ', error);
      throw new Error('API error: ' + (error instanceof Error ? error.message : 'An unexpected error occurred'));
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null); // Clear previous errors
    try {
      await auth.login(email, password);
      await categories()
      // Navigate to a protected route on successful login, e.g., dashboard
      navigate({ to: '/' }); // Adjust the target route as needed
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unknown error occurred');
    }
  };

  return (
    <div className="p-4 max-w-md mx-auto">
      <h2 className="text-xl font-semibold mb-4">Login</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="email" className="block text-sm font-medium text-gray-700">
            Email
          </label>
          <input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          />
        </div>
        <div>
          <label htmlFor="password"className="block text-sm font-medium text-gray-700">
            Password
          </label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          />
        </div>
        {error && <p className="text-red-500 text-sm">{error}</p>}
        <button
          type="submit"
          className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Login
        </button>
      </form>
    </div>
  );
}
