import { render, screen } from "@testing-library/react";
import Footer from "./Footer";

describe("Footer Component", () => {
  test("renders copyright text", () => {
    render(<Footer />);
    expect(screen.getByText(/© 2025 UFMarketPlace. All rights reserved./i)).toBeInTheDocument();
  });

  test("renders privacy policy and terms of service links", () => {
    render(<Footer />);
    expect(screen.getByRole("link", { name: /privacy policy/i })).toHaveAttribute("href", "/privacy-policy");
    expect(screen.getByRole("link", { name: /terms of service/i })).toHaveAttribute("href", "/terms-of-service");
  });

  test("renders social media icons with correct links", () => {
    render(<Footer />);
    expect(screen.getByLabelText(/facebook/i)).toHaveAttribute("href", "https://facebook.com");
    expect(screen.getByLabelText(/twitter/i)).toHaveAttribute("href", "https://twitter.com");
    expect(screen.getByLabelText(/instagram/i)).toHaveAttribute("href", "https://instagram.com");
  });

  test("social media icons have correct aria-labels", () => {
    render(<Footer />);
    expect(screen.getByLabelText(/facebook/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/twitter/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/instagram/i)).toBeInTheDocument();
  });
});