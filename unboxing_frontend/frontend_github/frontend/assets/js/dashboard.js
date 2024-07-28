document.addEventListener("DOMContentLoaded", function () {
  // Check if the user is authenticated
  const authToken = localStorage.getItem("authToken");
  if (!authToken) {
    // Redirect to login page if not authenticated
    window.location.href = "login.html";
  }

  // Logout functionality
  const logoutButton = document.getElementById("logoutButton");
  if (logoutButton) {
    logoutButton.addEventListener("click", function (event) {
      event.preventDefault();
      localStorage.removeItem("authToken");
      window.location.href = "login.html";
    });
  }

  // Highlight active menu link
  const navLinks = document.querySelectorAll(".nav-links a");
  navLinks.forEach((link) => {
    if (link.href === window.location.href) {
      link.classList.add("active");
    }
  });
});
