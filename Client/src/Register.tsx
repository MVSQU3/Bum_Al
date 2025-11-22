import React, { useState } from "react";
import { useAuthStore, type RegisterData } from "./store/userAuthStore";
import { Link } from "react-router-dom";

const Register = () => {
  const { register } = useAuthStore();
  const [formData, setFormData] = useState<RegisterData>({
    fullName: "",
    email: "",
    password: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    register(formData);
  };

  return (
    <>
      <div className="hero bg-base-200 min-h-screen">
        <div className="hero-content flex-col lg:flex-row-reverse">
          <div className="text-center lg:text-left">
            <h1 className="text-5xl font-bold">Sign Up now!</h1>
            <p className="py-6">
              Provident cupiditate voluptatem et in. Quaerat fugiat ut assumenda
              excepturi exercitationem quasi. In deleniti eaque aut repudiandae
              et a id nisi.
            </p>
          </div>
          <form onSubmit={handleSubmit} className="w-1/2">
            <div className="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
              <div className="card-body">
                <fieldset className="fieldset">
                  <label className="label">Email</label>
                  <input
                    type="text"
                    name="fullName"
                    value={formData.fullName || ""}
                    onChange={handleChange}
                    className="input"
                    placeholder="Full name"
                  />
                  <label className="label">Email</label>
                  <input
                    type="email"
                    name="email"
                    value={formData.email || ""}
                    onChange={handleChange}
                    className="input"
                    placeholder="Email"
                  />
                  <label className="label">Password</label>
                  <input
                    type="password"
                    name="password"
                    value={formData.password || ""}
                    className="input"
                    placeholder="Password"
                    onChange={handleChange}
                  />
                  <div>
                    <Link className="link link-hover" to={"#"}>
                      Forgot password?
                    </Link>
                  </div>
                  <div>
                    <Link className="link link-hover" to={"/login"}>
                      Vous avez un compte ? Connectez vous
                    </Link>
                  </div>
                  <button className="btn btn-neutral mt-4" type="submit">
                    Login
                  </button>
                </fieldset>
              </div>
            </div>
          </form>
        </div>
      </div>
    </>
  );
};

export default Register;
