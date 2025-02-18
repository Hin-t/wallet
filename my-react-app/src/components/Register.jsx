import React, { useState } from "react";
import { register } from "../api/walletApi";
import { TextField, Button, Paper, Typography, Container } from "@mui/material";

const Register = ({ setIsLoggedIn }) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const handleRegister = async () => {
    if (!username || !password || !confirmPassword) {
      setError("All fields are required.");
      return;
    }
    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }
    setError("");

    try {
      const response = await register(username, password);
      // Assuming the backend returns a 'success' status in its response
      if (response.data.success) {
        setMessage("Registration successful! Please log in.");
        // You can redirect to login or perform further actions here
      } else {
        setMessage(response.data.message || "Registration failed. Try a different username.");
      }
    } catch (error) {
      setMessage("Error registering user. Please try again.");
      setError(error.response?.data?.message || "An error occurred. Please try again.");
    }
  };

  return (
    <Container maxWidth="sm">
      <Paper elevation={3} sx={{ p: 3, mt: 3, textAlign: "center" }}>
        <Typography variant="h4" gutterBottom>
          Register
        </Typography>
        <TextField 
          label="Username" 
          fullWidth 
          value={username} 
          onChange={(e) => setUsername(e.target.value)} 
          sx={{ mb: 2 }} 
        />
        <TextField 
          label="Password" 
          type="password" 
          fullWidth 
          value={password} 
          onChange={(e) => setPassword(e.target.value)} 
          sx={{ mb: 2 }} 
        />
        <TextField 
          label="Confirm Password" 
          type="password" 
          fullWidth 
          value={confirmPassword} 
          onChange={(e) => setConfirmPassword(e.target.value)} 
          sx={{ mb: 3 }} 
        />
        <Button variant="contained" color="primary" fullWidth onClick={handleRegister}>
          Register
        </Button>
        {error && <Typography variant="body1" color="error" sx={{ mt: 2 }}>{error}</Typography>}
        {message && <Typography variant="body1" sx={{ mt: 2 }}>{message}</Typography>}
      </Paper>
    </Container>
  );
};

export default Register;
