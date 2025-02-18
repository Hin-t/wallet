import axios from "axios";

const BACKEND_API = "http://localhost:8081";

// Function to handle user login
export const login = async (username, password) => {
  try {
    const response = await axios.post(`${BACKEND_API}/user/login`, { username, password });
    return response;
  } catch (error) {
    throw new Error("Login failed");
  }
};

// Function to create a new account
export const createAccount = async () => {
  try {
    const response = await axios.post(`${BACKEND_API}/user/createAccount`);
    return response;
  } catch (error) {
    throw new Error("Account creation failed");
  }
};

// Function to query balance of an address
export const queryBalance = async (address) => {
  try {
    const response = await axios.get(`${BACKEND_API}/transaction/queryBalance?address=${address}`);
    return response;
  } catch (error) {
    throw new Error("Error querying balance");
  }
};

// Function to send a transaction
export const sendTransaction = async (from, to, amount) => {
  try {
    const response = await axios.post(`${BACKEND_API}/transaction/send`, { from, to, amount });
    return response;
  } catch (error) {
    throw new Error("Error sending transaction");
  }
};

// Function to register 
export const register = async (username, password) => {
    try {
      const response = await axios.post(`${BACKEND_API}/wallet/register`, { username, password });
      return response; // Make sure this response contains the expected data structure
    } catch (error) {
      throw error; // Handle error in the component
    }
  };