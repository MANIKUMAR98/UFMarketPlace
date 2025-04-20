import { useState, FormEvent, useEffect } from "react";
import { motion } from "framer-motion";
import "./Profile.css";
import Header from "../header/Header";
import { authService } from "../AuthService";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

interface User {
  name: string;
  email: string;
  phone?: string;
  address?: string;
}

const Profile = () => {
  const [showPasswordForm, setShowPasswordForm] = useState(false);
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [showDeleteConfirmation, setShowDeleteConfirmation] = useState(false);
  const [deleteError, setDeleteError] = useState("");
  const [deleteSuccess, setDeleteSuccess] = useState("");
  const [userEmail, setUserEmail] = useState("");
  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [address, setAddress] = useState("");
  const [isEditingPhone, setIsEditingPhone] = useState(false);
  const [isEditingAddress, setIsEditingAddress] = useState(false);

  useEffect(() => {
    const fetchUserProfile = async () => {
      const userProfile = await authService.getUserProfile();
      setUserEmail(userProfile.email);
      setName(userProfile.name);
      setPhone(userProfile.phone || "");
      setAddress(userProfile.address || "");
    };
    fetchUserProfile();
  }, []);

  const handlePasswordChange = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    if (newPassword !== confirmPassword) {
      return setError("New passwords do not match");
    }

    if (newPassword.length < 8) {
      return setError("Password must be at least 8 characters");
    }

    try {
      await authService.changePassword({
        name: name,
        email: userEmail,
        password: newPassword,
      });

      setSuccess("Password changed successfully!");
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");

      setTimeout(() => setShowPasswordForm(false), 1000);
    } catch (err) {
      setError("Failed to change password. Please try again.");
    }
  };

  const handleSave = async (field: "phone" | "address") => {
    try {
      const updates: Partial<User> =
        field === "phone" ? { phone } : { address };
      await authService.updateUserProfile({
        phone: phone,
        address: address,
      });
      field === "phone" ? setIsEditingPhone(false) : setIsEditingAddress(false);
    } catch (err) {
      alert(`Failed to update ${field}`);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="profile-container"
    >
      <h1 className="profile-title">Profile</h1>
      <div className="profile-content">
        <div className="profile-field">
          <label className="profile-label">Name</label>
          <p className="profile-value">{name}</p>
        </div>

        <div className="profile-field">
          <label className="profile-label">Email</label>
          <p className="profile-value">{userEmail}</p>
        </div>

        <div className="profile-field">
          <label className="profile-label">Phone</label>
          {isEditingPhone ? (
            <>
              <input
                type="text"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                className="form-input"
              />
              <div className="form-buttons">
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  className="password-submit-btn"
                  onClick={() => handleSave("phone")}
                >
                  Save
                </motion.button>
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  className="cancel-btn"
                  onClick={() => setIsEditingPhone(false)}
                >
                  Cancel
                </motion.button>
              </div>
            </>
          ) : (
            <div className="value-with-icon">
              <p className="profile-value">{phone}</p>
              <motion.button
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.95 }}
                className="icon-btn"
                onClick={() => setIsEditingPhone(true)}
              >
                <i className="fa-solid fa-pencil"></i>
              </motion.button>
            </div>
          )}
        </div>

        <div className="profile-field">
          <label className="profile-label">Address</label>
          {isEditingAddress ? (
            <>
              <input
                type="text"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                className="form-input"
              />
              <div className="form-buttons">
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  className="password-submit-btn"
                  onClick={() => handleSave("address")}
                >
                  Save
                </motion.button>
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  className="cancel-btn"
                  onClick={() => setIsEditingAddress(false)}
                >
                  Cancel
                </motion.button>
              </div>
            </>
          ) : (
            <div className="value-with-icon">
              <p className="profile-value">{address}</p>
              <motion.button
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.95 }}
                className="icon-btn"
                onClick={() => setIsEditingAddress(true)}
              >
                <i className="fa-solid fa-pencil"></i>
              </motion.button>
            </div>
          )}
        </div>

        {!showPasswordForm && (
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="change-password-btn"
            onClick={() => {
              setShowPasswordForm(true);
              setError("");
              setSuccess("");
            }}
          >
            Change Password
          </motion.button>
        )}

        {showPasswordForm && (
          <motion.form
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            onSubmit={handlePasswordChange}
            className="password-form"
          >
            <div className="form-group">
              <label htmlFor="currentPassword" className="form-label">
                Current Password
              </label>
              <input
                type="password"
                id="currentPassword"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
                className="form-input"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="newPassword" className="form-label">
                New Password
              </label>
              <input
                type="password"
                id="newPassword"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                className="form-input"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="confirmPassword" className="form-label">
                Confirm New Password
              </label>
              <input
                type="password"
                id="confirmPassword"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                className="form-input"
                required
              />
            </div>

            {(error || success) && (
              <div className="message-container">
                {error && <p className="error-message">{error}</p>}
                {success && <p className="success-message">{success}</p>}
              </div>
            )}

            <div className="form-buttons">
              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                type="submit"
                className="password-submit-btn"
              >
                Save Changes
              </motion.button>

              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                type="button"
                onClick={() => {
                  setShowPasswordForm(false);
                  setError("");
                  setSuccess("");
                }}
                className="cancel-btn"
              >
                Cancel
              </motion.button>
            </div>
          </motion.form>
        )}

        <motion.button
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          className="delete-account-btn"
          onClick={() => setShowDeleteConfirmation(true)}
        >
          Delete Account
        </motion.button>
      </div>

      {showDeleteConfirmation && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="delete-confirmation-modal"
        >
          <div className="delete-confirmation-content">
            <h3>Delete Account</h3>
            <p>
              Are you sure you want to delete your account? All your data will
              be permanently removed.
            </p>

            {(deleteError || deleteSuccess) && (
              <div className="message-container">
                {deleteError && <p className="error-message">{deleteError}</p>}
                {deleteSuccess && (
                  <p className="success-message">{deleteSuccess}</p>
                )}
              </div>
            )}

            <div className="confirmation-buttons">
              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                className="confirm-delete-btn"
                onClick={async () => {
                  try {
                    await authService.deleteUser();
                    setDeleteSuccess(
                      "Account deleted successfully. Redirecting..."
                    );
                    setTimeout(() => {
                      sessionStorage.clear();
                      window.location.href = "/login";
                    }, 1000);
                  } catch (err) {
                    setDeleteError(
                      "Failed to delete account. Please try again."
                    );
                  }
                }}
              >
                Delete
              </motion.button>

              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                className="cancel-delete-btn"
                onClick={() => {
                  setShowDeleteConfirmation(false);
                  setDeleteError("");
                  setDeleteSuccess("");
                }}
              >
                Cancel
              </motion.button>
            </div>
          </div>
        </motion.div>
      )}
    </motion.div>
  );
};

export default Profile;
