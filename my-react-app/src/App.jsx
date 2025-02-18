import React, { useState } from "react";
import Login from "./components/Login";
import Register from "./components/Register";
import Wallet from "./components/Wallet";
import { Button, Container, Paper, Typography } from "@mui/material";

const App = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [view, setView] = useState("home"); // "home", "login", "register"

  const handleRegisterSuccess = () => {
    setView("login"); // Redirect to login after successful registration
  };

  if (isLoggedIn) {
    return <Wallet setIsLoggedIn={setIsLoggedIn} />;
  }

  return (
    <Container maxWidth="sm">
      <Paper elevation={3} sx={{ p: 3, mt: 3, textAlign: "center" }}>
        {view === "home" && (
          <>
            <Typography variant="h4" gutterBottom>
              Welcome to Blockchain Wallet
            </Typography>
            <Button variant="contained" color="primary" fullWidth sx={{ mb: 2 }} onClick={() => setView("login")}>
              Login
            </Button>
            <Button variant="contained" color="secondary" fullWidth onClick={() => setView("register")}>
              Register
            </Button>
          </>
        )}
        {view === "login" && <Login setIsLoggedIn={setIsLoggedIn} />}
        {view === "register" && <Register setIsLoggedIn={setIsLoggedIn} onRegisterSuccess={handleRegisterSuccess} />}
      </Paper>
    </Container>
  );
};

export default App;