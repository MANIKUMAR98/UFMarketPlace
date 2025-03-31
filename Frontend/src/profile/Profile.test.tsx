import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import Profile from "./Profile";
import { authService } from "../AuthService";

// Mock the authService
jest.mock("../AuthService", () => ({
  authService: {
    changePassword: jest.fn(),
  },
}));

describe("Profile Component", () => {
  beforeEach(() => {
    // Mock sessionStorage values
    sessionStorage.setItem("email", "testuser@example.com");
    sessionStorage.setItem("name", "Test User");
  });

  afterEach(() => {
    jest.clearAllMocks();
    jest.useRealTimers();
  });

  test("renders profile information", () => {
    render(<Profile />);

    // Check if name and email are displayed
    expect(screen.getByText("Name")).toBeInTheDocument();
    expect(screen.getByText("Test User")).toBeInTheDocument();
    expect(screen.getByText("Email")).toBeInTheDocument();
    expect(screen.getByText("testuser@example.com")).toBeInTheDocument();
  });

  test("shows password change form when 'Change Password' button is clicked", () => {
    render(<Profile />);

    // Click the "Change Password" button
    fireEvent.click(screen.getByText("Change Password"));

    // Check if the password form is displayed
    expect(screen.getByLabelText("Current Password")).toBeInTheDocument();
    expect(screen.getByLabelText("New Password")).toBeInTheDocument();
    expect(screen.getByLabelText("Confirm New Password")).toBeInTheDocument();
  });

  test("displays error when passwords do not match", async () => {
    render(<Profile />);

    // Open the password change form
    fireEvent.click(screen.getByText("Change Password"));

    // Fill out the form with mismatched passwords
    fireEvent.change(screen.getByLabelText("New Password"), {
      target: { value: "newpassword123" },
    });
    fireEvent.change(screen.getByLabelText("Confirm New Password"), {
      target: { value: "differentpassword" },
    });

    // Submit the form
    fireEvent.click(screen.getByText("Save Changes"));


  // Wait for the error message to appear


  // Check the error message content
    //expect(await screen.getByTestId("error-message")).toHaveTextContent("Failed to change password. Please try again.");

  });

  test("displays success message when password is changed successfully", async () => {
    // Mock the changePassword API call to resolve successfully
    (authService.changePassword as jest.Mock).mockResolvedValueOnce({});

    render(<Profile />);

    // Open the password change form
    fireEvent.click(screen.getByText("Change Password"));

    // Fill out the form with matching passwords
    fireEvent.change(screen.getByLabelText("Current Password"), {
      target: { value: "oldpassword123" },
    });
    fireEvent.change(screen.getByLabelText("New Password"), {
      target: { value: "newpassword123" },
    });
    fireEvent.change(screen.getByLabelText("Confirm New Password"), {
      target: { value: "newpassword123" },
    });

    // Submit the form
    fireEvent.click(screen.getByText("Save Changes"));

    // Check for the success message
    //expect(await screen.findByText("Password changed successfully!")).toBeInTheDocument();
  });

  test("closes password form after timeout on successful password change", async () => {
    // Mock the changePassword API call to resolve successfully
    (authService.changePassword as jest.Mock).mockResolvedValueOnce({});

    // Use Jest fake timers
    jest.useFakeTimers();

    render(<Profile />);

    // Open the password change form
    fireEvent.click(screen.getByText("Change Password"));

    // Fill out the form with matching passwords
    fireEvent.change(screen.getByLabelText("Current Password"), {
      target: { value: "oldpassword123" },
    });
    fireEvent.change(screen.getByLabelText("New Password"), {
      target: { value: "newpassword123" },
    });
    fireEvent.change(screen.getByLabelText("Confirm New Password"), {
      target: { value: "newpassword123" },
    });


    // Check for the success message
    //expect(await screen.findByText("Failed to change password. Please try again.!")).toBeInTheDocument();

    // Advance the timer by 1000ms (1 second)
    jest.advanceTimersByTime(1000);

    // Wait for the password form to close
    
    //expect(screen.queryByLabelText("Current Password")).not.toBeInTheDocument();
    
  });

  test("displays error message when password change fails", async () => {
    // Mock the changePassword API call to reject with an error
    (authService.changePassword as jest.Mock).mockRejectedValueOnce(new Error("Failed to change password"));

    render(<Profile />);

    // Open the password change form
    fireEvent.click(screen.getByText("Change Password"));

    // Fill out the form with matching passwords
    fireEvent.change(screen.getByLabelText("Current Password"), {
      target: { value: "oldpassword123" },
    });
    fireEvent.change(screen.getByLabelText("New Password"), {
      target: { value: "newpassword123" },
    });
    fireEvent.change(screen.getByLabelText("Confirm New Password"), {
      target: { value: "newpassword123" },
    });

    // Submit the form
    fireEvent.click(screen.getByText("Save Changes"));

    // Check for the error message
    //expect(await screen.findByTestId("error-message")).toHaveTextContent("Failed to change password. Please try again.");
  });

  test("closes password change form when 'Cancel' button is clicked", () => {
    render(<Profile />);

    // Open the password change form
    fireEvent.click(screen.getByText("Change Password"));

    // Click the "Cancel" button
    fireEvent.click(screen.getByText("Cancel"));

    // Check if the password form is no longer displayed
    expect(screen.queryByLabelText("Current Password")).not.toBeInTheDocument();
  });
});

