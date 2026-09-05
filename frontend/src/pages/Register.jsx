import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import client, { apiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";

export default function Register() {
  const { loginWithToken } = useAuth();
  const navigate = useNavigate();

  const [form, setForm] = useState({ username: "", email: "", password: "", confirm: "" });
  const [errors, setErrors] = useState({});
  const [serverError, setServerError] = useState("");
  const [busy, setBusy] = useState(false);

  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.value }));

  function validate() {
    const errs = {};
    if (form.username.trim().length < 3) errs.username = "At least 3 characters";
    if (!/^\S+@\S+\.\S+$/.test(form.email)) errs.email = "Enter a valid email";
    if (form.password.length < 8) errs.password = "At least 8 characters";
    if (form.confirm !== form.password) errs.confirm = "Passwords do not match";
    return errs;
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setServerError("");
    const errs = validate();
    setErrors(errs);
    if (Object.keys(errs).length > 0) return;

    setBusy(true);
    try {
      await client.post("/users/register", {
        username: form.username.trim(),
        email: form.email.trim(),
        password: form.password,
      });
      // Auto-login after successful registration
      const { data } = await client.post("/users/login", {
        email: form.email.trim(),
        password: form.password,
      });
      loginWithToken(data.token, { email: form.email.trim(), username: form.username.trim() });
      navigate("/", { replace: true });
    } catch (err) {
      setServerError(apiError(err));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="page auth-page">
      <form className="card auth-card" onSubmit={handleSubmit} noValidate>
        <h1>Create account</h1>
        <p className="muted">Join CineDeck — it takes a minute.</p>

        <label className="field">
          <span>Username</span>
          <input value={form.username} onChange={set("username")} placeholder="filmfan99" />
          {errors.username && <em className="field-error">{errors.username}</em>}
        </label>

        <label className="field">
          <span>Email</span>
          <input type="email" value={form.email} onChange={set("email")} placeholder="you@example.com" />
          {errors.email && <em className="field-error">{errors.email}</em>}
        </label>

        <label className="field">
          <span>Password</span>
          <input
            type="password"
            value={form.password}
            onChange={set("password")}
            placeholder="min. 8 characters"
          />
          {errors.password && <em className="field-error">{errors.password}</em>}
        </label>

        <label className="field">
          <span>Confirm password</span>
          <input type="password" value={form.confirm} onChange={set("confirm")} placeholder="repeat password" />
          {errors.confirm && <em className="field-error">{errors.confirm}</em>}
        </label>

        {serverError && <div className="error-inline">{serverError}</div>}

        <button type="submit" className="btn btn-primary btn-block" disabled={busy}>
          {busy ? "Creating…" : "Sign up"}
        </button>

        <p className="muted center">
          Already registered? <Link to="/login">Log in</Link>
        </p>
      </form>
    </div>
  );
}
