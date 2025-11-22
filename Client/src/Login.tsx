import React, { useState } from "react";
import { useAuthStore } from "./store/userAuthStore";
import { useNavigate } from "react-router-dom";

const Login = () => {
  const { login } = useAuthStore();
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };
  // const navigate = useNavigate();
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    login(formData);
  };

  return (
    <>
      <div>Login</div>
      <form onSubmit={handleSubmit} className="flex flex-col">
        <input
          type="email"
          value={formData.email}
          placeholder="email"
          className="input"
          name="email"
          onChange={handleChange}
        />
        <input
          type="password"
          value={formData.password}
          name="password"
          placeholder="password"
          className="input"
          onChange={handleChange}
        />
        <button type="submit" className="btn btn-accent w-1/3">
          Login
        </button>
      </form>
    </>
  );
};

export default Login;
