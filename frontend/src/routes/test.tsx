import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/test')({
  component: RouteComponent,
})

function RouteComponent() {
  const categories = async ():Promise<void>  => {
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/categories`, {
        method: 'GET',
        credentials: 'include', // Important for cookies
        headers: {
          'Content-Type': 'application/json',
        },
        
      });
      if (!response.ok) {
        
      }
      return ;
    } catch (error) {
      console.log('Error fetching categories: ', error);
      throw new Error('Login error: ' + (error instanceof Error ? error.message : 'An unexpected error occurred'));
    }
  };

  return <div>
    <button 
   
    onClick={async () => {
      const response = await categories();
      console.log(response);
    }}>
    Test``
  </button></div>;
}
