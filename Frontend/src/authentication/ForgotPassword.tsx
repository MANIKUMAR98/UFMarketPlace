import { useState, FormEvent } from "react";
import { motion } from "framer-motion";
import { authService } from "../AuthService";
import { useNavigate } from "react-router-dom";
import "./ForgotPassword.css";

const ForgotPassword = () => {
  const [email, setEmail] = useState("");
  const [otp, setOtp] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [step, setStep] = useState("email"); // "email" or "reset"
  const navigate = useNavigate();

  const handleSendOtp = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    if (!email.includes("ufl.edu")) {
      return setError("Only UF emails are allowed");
    }

    setIsLoading(true);

    try {
      sessionStorage.setItem("email", email);
      await authService.sendEmailVerificationCode();
      setSuccess("OTP sent to your email!");
      setTimeout(() => {
        setSuccess("");
        setStep("reset");
      }, 1500);
    } catch (err) {
      setError("Failed to send OTP. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleResetPassword = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    if (password !== confirmPassword) {
      return setError("Passwords do not match");
    }

    if (password.length < 8) {
      return setError("Password must be at least 8 characters long");
    }

    setIsLoading(true);

    try {
      await authService.resetPassword(email, otp, password);
      setSuccess("Password reset successfully!");
      setTimeout(() => navigate("/login"), 1500);
    } catch (err) {
      setError(
        "Failed to reset password. Please check your OTP and try again."
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="auth-form"
      >
        <h2>Forgot Password</h2>

        {step === "email" ? (
          <form className="pass-form" onSubmit={handleSendOtp}>
            <div className="input-container">
              <input
                type="email"
                placeholder="UF Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            {error && <div className="error-message">{error}</div>}
            {success && <div className="success-message">{success}</div>}

            <button type="submit" disabled={isLoading} className="submit-btn">
              {isLoading ? "Sending..." : "Send OTP"}
            </button>

            <button
              type="button"
              onClick={() => navigate("/login")}
              className="link-button"
            >
              Back to Login
            </button>
          </form>
        ) : (
          <motion.form
            className="pass-form"
            onSubmit={handleResetPassword}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
          >
            <div className="input-container">
              <input
                type="text"
                placeholder="Enter OTP"
                value={otp}
                onChange={(e) => setOtp(e.target.value)}
                required
              />
            </div>

            <div className="input-container">
              <input
                type="password"
                placeholder="New Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>

            <div className="input-container">
              <input
                type="password"
                placeholder="Confirm New Password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
            </div>

            {error && <div className="error-message">{error}</div>}
            {success && <div className="success-message">{success}</div>}

            <button type="submit" disabled={isLoading} className="submit-btn">
              {isLoading ? "Resetting..." : "Reset Password"}
            </button>

            <button
              type="button"
              onClick={() => setStep("email")}
              className="link-button"
            >
              Back to Email
            </button>
          </motion.form>
        )}
      </motion.div>
    </div>
  );
};

export default ForgotPassword;
