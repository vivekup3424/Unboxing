document.addEventListener("DOMContentLoaded", function () {
  // Existing code for greeting message
  const greeting = document.createElement("p");
  const currentTime = new Date().getHours();

  console.log("Current Time:", currentTime); // Debugging: Log the current time

  if (currentTime < 12) {
    greeting.textContent = "Good morning! Welcome to your dashboard.";
  } else if (currentTime < 18) {
    greeting.textContent = "Good afternoon! Ready to manage your tasks?";
  } else {
    greeting.textContent = "Good evening! Let's wrap up your day.";
  }

  console.log("Greeting Message:", greeting.textContent); // Debugging: Log the greeting message

  const heroElement = document.querySelector(".hero");
  if (heroElement) {
    heroElement.appendChild(greeting);
    console.log("Greeting appended to .hero"); // Debugging: Confirm greeting appended
  } else {
    console.error("Error: .hero element not found"); // Debugging: Log if .hero is not found
  }

  // Registration form handling
  const registerForm = document.getElementById("registerForm");
  const messageElement = document.getElementById("message");

  if (!registerForm) {
    console.error("Error: Register form element not found"); // Debugging: Log if register form not found
    return; // Exit if the form is not found
  }

  if (!messageElement) {
    console.error("Error: Message element not found"); // Debugging: Log if message element not found
  }

  registerForm.addEventListener("submit", async function (event) {
    event.preventDefault();

    const formData = new FormData(registerForm);
    const data = {
      name: formData.get("name"),
      email: formData.get("email"),
      password: formData.get("password"),
      "secret-key": formData.get("secret-key"),
    };

    console.log("Form Data:", data); // Debugging: Log form data

    try {
      const response = await fetch("http://localhost:4000/admin/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: new URLSearchParams(data),
        mode: "cors",
      });

      console.log("Response:", response); // Debugging: Log the response object

      if (response.ok) {
        console.log("Registration successful, redirecting to login page"); // Debugging: Log success
        window.location.href = "login.html"; // Redirect to login page
      } else {
        const errorText = await response.text();
        console.error("Registration failed:", errorText); // Debugging: Log error text
        showMessage(errorText);
      }
    } catch (error) {
      showMessage("An error occurred. Please try again.");
      console.error("Error:", error); // Debugging: Log any caught errors
    }
  });

  function showMessage(message) {
    messageElement.textContent = message;
    messageElement.style.display = "block";
    console.log("Message displayed:", message); // Debugging: Log the displayed message
  }
});
