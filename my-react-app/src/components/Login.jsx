import React, { useState } from "react";
import { login } from "../api/walletApi";
import { TextField, Button, Paper, Typography, Container } from "@mui/material";

const Login = ({ setIsLoggedIn }) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const handleLogin = async () => {
    if (username.trim() === "" || password.trim() === "") {
      setError("Username and password cannot be empty.");
      return;
    }
    setError("");
    try {
      const response = await login(username, password);
      if (response.data.success) {
        setIsLoggedIn(true);
        setMessage("Login successful!");
      } else {
        setMessage("Invalid credentials.");
      }
    } catch (error) {
      setMessage("Error logging in.");
    }
  };

  return (
    <Container maxWidth="sm">
      <Paper elevation={3} style={{ padding: "20px", marginTop: "20px", textAlign: "center" }}>
        <Typography variant="h4" gutterBottom>
          Login to Wallet
        </Typography>
        <TextField
          label="Username"
          fullWidth
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          error={!!error}
          style={{ marginBottom: "10px" }}
        />
        <TextField
          label="Password"
          type="password"
          fullWidth
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          error={!!error}
          style={{ marginBottom: "20px" }}
        />
        <Button variant="contained" color="primary" fullWidth onClick={handleLogin}>
          Login
        </Button>
        <Typography variant="body1" color="error" style={{ marginTop: "10px" }}>
          {error}
        </Typography>
        <Typography variant="body1" style={{ marginTop: "10px" }}>
          {message}
        </Typography>
      </Paper>
    </Container>
  );
};

export default Login;
