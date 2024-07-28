document.addEventListener("DOMContentLoaded", function () {
  const loginForm = document.getElementById("loginForm");
  const messageElement = document.getElementById("message");
  const formButton = document.getElementById("formButton");
  if (!loginForm) {
    console.error("Error: Login form element not found"); // Debugging: Log if register form not found
    return; // Exit if the form is not found
  }

  if (!messageElement) {
    console.error("Error: Message element not found"); // Debugging: Log if message element not found
  }
  formButton.addEventListener("click", async function (event) {
    event.preventDefault(); // Prevent the form from submitting the traditional way
    console.log("submit event triggered");

    const formData = new FormData(loginForm);
    const data = {
      email: formData.get("email"),
      password: formData.get("password"),
    };

    console.log("Sending login data:", data); // Debugging: Log the data being sent

    try {
      const response = await fetch(
        "http://localhost:4000/tokens/authentication", // Ensure this URL is correct
        {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded", // Correct Content-Type
          },
          body: new URLSearchParams(data), // Ensure data is being sent
          mode: "cors",
        }
      );

      console.log("Response status:", response.status); // Debugging: Log response status

      if (response.ok) {
        const jsonResponse = await response.json();
        console.log("Authentication successful:", jsonResponse);

        // Store the authentication token in localStorage or sessionStorage
        localStorage.setItem("authToken", jsonResponse.authentication_token);

        // Redirect to the dashboard or another page
        window.location.href = "dashboard.html";
      } else {
        const errorText = await response.text();
        showMessage("Login failed: " + errorText);
        console.error("Error response:", errorText); // Debugging: Log error response
      }
    } catch (error) {
      showMessage("An error occurred. Please try again.");
      console.error("Error:", error);
    }
  });

  function showMessage(message) {
    messageElement.textContent = message;
    messageElement.style.display = "block";
  }
});
