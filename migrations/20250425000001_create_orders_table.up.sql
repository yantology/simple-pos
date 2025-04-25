-- filepath: c:\Users\LENOVO\Documents\golang\simple-pos\migrations\20250425000001_create_orders_table.up.sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    total NUMERIC(12, 2) NOT NULL,
    product JSONB NOT NULL, -- Store product information as JSONB to retain even if products are deleted
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create trigger on orders table
CREATE TRIGGER update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
